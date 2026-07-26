import playwright from '../../site/node_modules/playwright/index.js'
import { execFileSync } from 'node:child_process'
import { chmodSync, existsSync, mkdirSync, readFileSync, unlinkSync, writeFileSync } from 'node:fs'
import { resolve } from 'node:path'

const baseUrl = process.env.APIAUTH_ADMIN_URL || 'https://apiauth.whiteyun.com/admin/'
const { chromium } = playwright
const credentialsFile = process.env.APIAUTH_ADMIN_CREDENTIALS || '/root/apiauth-admin-credentials.txt'
const outputRoot = resolve(process.argv[2] || 'artifacts/admin-theme-qa')
const mode = process.argv.includes('--dark') ? 'dark' : 'light'
const viewportName = process.argv.find((arg) => arg.startsWith('--viewport='))?.split('=')[1] || 'desktop'
const viewports = {
  mobile: { width: 390, height: 844 },
  tablet: { width: 768, height: 1024 },
  desktop: { width: 1440, height: 900 },
  issue: { width: 2048, height: 1152 },
  user: { width: 1463, height: 780 },
  report: { width: 2560, height: 1364 },
  reference: { width: 2520, height: 1416 }
}
const viewport = viewports[viewportName]
const deviceScaleFactor = Number(process.argv.find((arg) => arg.startsWith('--dpr='))?.split('=')[1] || 1)
const captchaCodeFile = process.env.APIAUTH_CAPTCHA_CODE_FILE
const mockBackend = process.env.APIAUTH_QA_MOCK === '1'
const navType = Number(process.argv.find((arg) => arg.startsWith('--nav='))?.split('=')[1] || 1)
const menuStyle = process.argv.find((arg) => arg.startsWith('--menu='))?.split('=')[1] || 'theme-light'
const shouldScanPages = process.argv.includes('--scan-pages')
const shouldExerciseControls = process.argv.includes('--exercise-controls')
const shouldExerciseHome = process.argv.includes('--exercise-home')
const targetPath = process.argv.find((arg) => arg.startsWith('--path='))?.split('=')[1] || 'system/settings'
const shouldOpenLayout = !process.argv.includes('--skip-layout')

if (!viewport) {
  throw new Error(`未知视口：${viewportName}`)
}
if (!existsSync(credentialsFile)) {
  throw new Error('未找到后台验收凭据文件')
}

mkdirSync(outputRoot, { recursive: true })
const credentials = readFileSync(credentialsFile, 'utf8')
const username = credentials.match(/用户名\s*[：:]\s*(\S+)/)?.[1]
const password = credentials.match(/密码\s*[：:]\s*(\S+)/)?.[1]

if (!username || !password) {
  throw new Error('后台验收凭据格式不正确')
}

function recognizeCaptcha(imagePath) {
  const variants = [
    ['--psm', '7'],
    ['--psm', '8']
  ]
  for (const args of variants) {
    try {
      const text = execFileSync('tesseract', [imagePath, 'stdout', ...args], {
        encoding: 'utf8',
        stdio: ['ignore', 'pipe', 'ignore']
      })
      const candidate = text.replace(/[^a-z0-9]/gi, '').slice(0, 6)
      if (candidate.length >= 4) return candidate
    } catch {
      // 尝试下一种识别模式。
    }
  }
  return ''
}

async function login(page) {
  await page.goto(baseUrl, { waitUntil: 'networkidle' })
  if (!page.url().includes('/login')) return
  await page.locator('#app').waitFor({ state: 'attached' })

  for (let attempt = 1; attempt <= 8; attempt += 1) {
    await page.getByPlaceholder('账号').fill(username)
    await page.getByPlaceholder('密码').fill(password)
    const captcha = page.locator('.login-code-img')
    if (await captcha.count()) {
      const rawPath = resolve(outputRoot, '.captcha.png')
      const preparedPath = resolve(outputRoot, '.captcha-prepared.png')
      await captcha.screenshot({ path: rawPath })
      try {
        execFileSync('convert', [
          rawPath,
          '-colorspace', 'Gray',
          '-resize', '400%',
          '-contrast-stretch', '2%x2%',
          '-threshold', '58%',
          preparedPath
        ], { stdio: 'ignore' })
      } catch {
        // ImageMagick 不可用时直接识别原图。
      }
      let code = recognizeCaptcha(existsSync(preparedPath) ? preparedPath : rawPath)
      if (captchaCodeFile) {
        process.stdout.write(`CAPTCHA_READY ${attempt} ${rawPath}\n`)
        for (let waitCount = 0; waitCount < 100; waitCount += 1) {
          if (existsSync(captchaCodeFile)) {
            code = readFileSync(captchaCodeFile, 'utf8').replace(/\s+/g, '')
            unlinkSync(captchaCodeFile)
            break
          }
          await new Promise((resolveWait) => setTimeout(resolveWait, 500))
        }
      }
      if (code.length < 4) {
        await page.locator('.login-code').click()
        await page.waitForTimeout(300)
        continue
      }
      await page.getByPlaceholder('验证码').fill(code)
    }
    await page.getByRole('button', { name: '进入控制台' }).click()
    await page.waitForTimeout(1400)
    if (!page.url().includes('/login')) return
    const feedback = await page.locator('.el-message__content').allTextContents()
    if (feedback.length) {
      process.stdout.write(`LOGIN_FEEDBACK ${feedback.join(' / ')}\n`)
    }
    if (await page.locator('.login-code').count()) {
      await page.locator('.login-code').click()
      await page.waitForTimeout(300)
    }
  }
  throw new Error('后台验证码连续识别失败，未能建立验收会话')
}

async function collectThemeDefects(page) {
  return page.evaluate(() => {
    const parseRgb = (value) => {
      const match = value.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/)
      return match ? match.slice(1).map(Number) : null
    }
    const luminance = (rgb) => {
      if (!rgb) return 0
      const channels = rgb.map((value) => {
        const channel = value / 255
        return channel <= 0.03928 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4
      })
      return channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722
    }
    const selectorFor = (element) => {
      if (element.id) return `#${element.id}`
      const classes = [...element.classList].slice(0, 3)
      return `${element.tagName.toLowerCase()}${classes.map((name) => `.${name}`).join('')}`
    }
    return [...document.querySelectorAll('body *')]
      .map((element) => {
        const rect = element.getBoundingClientRect()
        const style = getComputedStyle(element)
        return {
          selector: selectorFor(element),
          area: Math.round(rect.width * rect.height),
          background: style.backgroundColor,
          text: style.color,
          backgroundLuminance: luminance(parseRgb(style.backgroundColor)),
          textLuminance: luminance(parseRgb(style.color))
        }
      })
      .filter((item) => item.area > 3000 && item.backgroundLuminance > 0.72)
      .sort((a, b) => b.area - a.area)
      .slice(0, 40)
  })
}

async function collectNavigationContrast(page) {
  return page.evaluate(() => {
    const parseColor = (value) => {
      const match = value.match(/rgba?\(([\d.]+)[,\s]+([\d.]+)[,\s]+([\d.]+)(?:[,\s/]+([\d.]+))?/)
      return match
        ? { rgb: match.slice(1, 4).map(Number), alpha: match[4] === undefined ? 1 : Number(match[4]) }
        : null
    }
    const composite = (foreground, background) => foreground.rgb.map(
      (value, index) => Math.round(value * foreground.alpha + background[index] * (1 - foreground.alpha))
    )
    const effectiveBackground = (element) => {
      const layers = []
      let current = element
      while (current) {
        const color = parseColor(getComputedStyle(current).backgroundColor)
        if (color && color.alpha > 0) layers.push(color)
        current = current.parentElement
      }
      return layers.reverse().reduce((background, layer) => composite(layer, background), [255, 255, 255])
    }
    const luminance = (rgb) => {
      const channels = rgb.map((value) => {
        const channel = value / 255
        return channel <= 0.03928 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4
      })
      return channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722
    }
    const contrast = (foreground, background) => {
      const light = Math.max(luminance(foreground), luminance(background))
      const dark = Math.min(luminance(foreground), luminance(background))
      return (light + 0.05) / (dark + 0.05)
    }
    const selectorFor = (element) => {
      if (element.id) return `#${element.id}`
      return `${element.tagName.toLowerCase()}${[...element.classList].slice(0, 3).map((name) => `.${name}`).join('')}`
    }
    const selectors = [
      '.navbar .sidebar-wordmark',
      '.navbar .sidebar-console',
      '.navbar .hamburger-container',
      '.navbar .breadcrumb-container',
      '.navbar .topmenu-container > .el-menu-item',
      '.navbar .topmenu-container > .el-sub-menu > .el-sub-menu__title',
      '.navbar .topbar-menu .el-menu-item',
      '.navbar .topbar-menu .el-sub-menu__title',
      '.navbar .right-menu .right-menu-item',
      '.navbar .theme-segmented button'
    ]
    const elements = [...new Set(selectors.flatMap((selector) => [...document.querySelectorAll(selector)]))]
    const samples = elements
      .filter((element) => {
        const rect = element.getBoundingClientRect()
        const style = getComputedStyle(element)
        return rect.width > 0 && rect.height > 0 && style.visibility !== 'hidden' && Number(style.opacity) > 0
      })
      .map((element) => {
        const style = getComputedStyle(element)
        const foreground = parseColor(style.color)?.rgb || [0, 0, 0]
        const background = effectiveBackground(element)
        return {
          selector: selectorFor(element),
          text: (element.textContent || '').trim().replace(/\s+/g, ' ').slice(0, 80),
          foreground: `rgb(${foreground.join(', ')})`,
          background: `rgb(${background.join(', ')})`,
          contrast: Number(contrast(foreground, background).toFixed(2))
        }
      })
    return {
      samples,
      defects: samples.filter((sample) => sample.contrast < 4.5)
    }
  })
}

const browser = await chromium.launch({ headless: true })
const context = await browser.newContext({
  viewport,
  deviceScaleFactor,
  colorScheme: mode
})
await context.addInitScript(({ theme, nav, menu }) => {
  localStorage.setItem('apiauth-theme', theme)
  const current = JSON.parse(localStorage.getItem('layout-setting') || '{}')
  localStorage.setItem('layout-setting', JSON.stringify({
    ...current,
    navType: nav,
    sideTheme: menu
  }))
}, { theme: mode, nav: navType, menu: menuStyle })
const page = await context.newPage()
const pageErrors = []
const consoleErrors = []
const failedRequests = []
page.on('pageerror', (error) => pageErrors.push(error.message))
page.on('console', (message) => {
  if (message.type() === 'error') consoleErrors.push(message.text())
})
page.on('requestfailed', (request) => {
  failedRequests.push({
    url: request.url(),
    error: request.failure()?.errorText || 'unknown'
  })
})

try {
  if (mockBackend) {
    await context.addCookies([{
      name: 'Admin-Token',
      value: 'theme-qa-token',
      url: baseUrl
    }])
    await page.route('**/api/**', async (route) => {
      const url = new URL(route.request().url())
      const path = url.pathname.replace(/^\/api/, '')
      if (path === '/getInfo') {
        await route.fulfill({
          contentType: 'application/json',
          body: JSON.stringify({
            code: 200,
            user: { userId: 1, userName: 'theme_qa', nickName: 'APIAuth 管理员', avatar: '' },
            roles: ['admin'],
            permissions: ['*:*:*'],
            isDefaultModifyPwd: false,
            isPasswordExpired: false
          })
        })
        return
      }
      if (path === '/getRouters') {
        await route.fulfill({
          contentType: 'application/json',
          body: JSON.stringify({
            code: 200,
            data: [
              {
                path: '/system',
                component: 'Layout',
                name: 'System',
                alwaysShow: true,
                meta: { title: '系统管理', icon: 'system' },
                children: [
                  {
                    path: 'settings',
                    component: 'system/settings/index',
                    name: 'SystemSettings',
                    meta: { title: '系统配置', icon: 'edit' }
                  },
                  {
                    path: 'user',
                    component: 'system/user/index',
                    name: 'User',
                    meta: { title: '用户管理', icon: 'user' }
                  },
                  {
                    path: 'role',
                    component: 'system/role/index',
                    name: 'Role',
                    meta: { title: '角色管理', icon: 'peoples' }
                  },
                  {
                    path: 'menu',
                    component: 'system/menu/index',
                    name: 'Menu',
                    meta: { title: '菜单管理', icon: 'tree-table' }
                  },
                  {
                    path: 'dict',
                    component: 'system/dict/index',
                    name: 'Dict',
                    meta: { title: '字典管理', icon: 'dict' }
                  }
                ]
              },
              {
                path: '/auth',
                component: 'Layout',
                name: 'Authorization',
                alwaysShow: true,
                meta: { title: '授权管理', icon: 'lock' },
                children: [
                  {
                    path: 'licenses',
                    component: 'auth/licenses/index',
                    name: 'Licenses',
                    meta: { title: '许可证', icon: 'list' }
                  },
                  {
                    path: 'catalog',
                    component: 'auth/catalog/index',
                    name: 'Catalog',
                    meta: { title: '应用与版本', icon: 'build' }
                  },
                  {
                    path: 'operations',
                    component: 'auth/operations/index',
                    name: 'Operations',
                    meta: { title: '运行态', icon: 'monitor' }
                  },
                  {
                    path: 'signing-keys',
                    component: 'auth/signing-keys/index',
                    name: 'SigningKeys',
                    meta: { title: '签名密钥', icon: 'key' }
                  },
                  {
                    path: 'audit',
                    component: 'auth/audit/index',
                    name: 'AuthorizationAudit',
                    meta: { title: '授权审计', icon: 'log' }
                  }
                ]
              }
            ]
          })
        })
        return
      }
      if (path === '/system/setting') {
        await route.fulfill({
          contentType: 'application/json',
          body: JSON.stringify({
            code: 200,
            data: {
              site: {
                title: 'APIAuth',
                logo: '',
                favicon: '/admin/favicon.ico',
                description: '多应用商业授权管理平台',
                keywords: 'APIAuth,许可证,授权管理',
                frontendHeadCode: '',
                siteUrl: 'https://apiauth.whiteyun.com',
                icpNo: '',
                publicSecurityNo: '',
                customerServiceEmail: '',
                copyright: 'Copyright © 2026 APIAuth. All Rights Reserved.',
                defaultLanguage: 'zh-CN',
                enableSeo: true
              },
              generated: { notifyUrl: '', returnUrl: '' }
            }
          })
        })
        return
      }
      await route.fulfill({
        contentType: 'application/json',
        body: JSON.stringify({ code: 200, data: [], rows: [], total: 0 })
      })
    })
  } else {
    await login(page)
  }
  await context.storageState({ path: resolve(outputRoot, '.storage-state.json') })
  chmodSync(resolve(outputRoot, '.storage-state.json'), 0o600)

  const targetUrl = new URL(targetPath, baseUrl).toString()
  await page.goto(targetUrl, { waitUntil: 'networkidle' })
  await page.waitForTimeout(1400)

  const label = `${viewportName}-${mode}-nav${navType}-${menuStyle}`
  const captureName = targetPath === 'system/settings'
    ? 'system-settings'
    : targetPath.replace(/^\/+|\/+$/g, '').replace(/[^a-z0-9]+/gi, '-') || 'home'
  await page.screenshot({
    path: resolve(outputRoot, `${label}-${captureName}.png`),
    fullPage: true
  })
  const defects = mode === 'dark' ? await collectThemeDefects(page) : []
  const navigationContrast = await collectNavigationContrast(page)
  const computed = await page.evaluate(() => {
    const sample = (selector) => {
      const element = document.querySelector(selector)
      if (!element) return null
      const style = getComputedStyle(element)
      return {
        background: style.backgroundColor,
        color: style.color,
        opacity: style.opacity,
        visibility: style.visibility
      }
    }
    return {
      htmlClass: document.documentElement.className,
      body: sample('body'),
      settings: sample('.system-settings'),
      settingsTabs: sample('.settings-tabs'),
      input: sample('.el-input__wrapper'),
      formLabel: sample('.el-form-item__label'),
      loadingMasks: [...document.querySelectorAll('.el-loading-mask')].map((element) => {
        const style = getComputedStyle(element)
        return { display: style.display, visibility: style.visibility, opacity: style.opacity }
      }),
      horizontalOverflow: document.documentElement.scrollWidth - window.innerWidth
    }
  })
  const pageAudits = []

  if (shouldScanPages && mode === 'dark') {
    const scanPaths = [
      'system/user',
      'system/role',
      'system/menu',
      'system/dict',
      'auth/licenses',
      'auth/catalog',
      'auth/operations',
      'auth/signing-keys',
      'auth/audit'
    ]
    for (const scanPath of scanPaths) {
      process.stdout.write(`SCAN_PAGE ${scanPath}\n`)
      await page.goto(new URL(scanPath, baseUrl).toString(), { waitUntil: 'networkidle' })
      await page.waitForTimeout(900)
      const scanDefects = await collectThemeDefects(page)
      const horizontalOverflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth)
      pageAudits.push({
        path: scanPath,
        defects: scanDefects,
        horizontalOverflow
      })
    }
    await page.goto(targetUrl, { waitUntil: 'networkidle' })
    await page.waitForTimeout(900)
  }

  const interactionAudits = []
  if (shouldOpenLayout) {
    await page.locator('.avatar-container').hover()
    const layoutSettingItem = page.getByText('布局设置', { exact: true })
    await layoutSettingItem.waitFor({ state: 'visible' })
    await layoutSettingItem.click()
    await page.locator('.el-drawer.open').waitFor({ state: 'visible' })
    await page.waitForTimeout(300)
    await page.screenshot({
      path: resolve(outputRoot, `${label}-layout-settings.png`),
      fullPage: false
    })
    if (shouldExerciseControls) {
      const navigationChoices = page.locator('.layout-choice')
      const menuChoices = page.locator('.menu-style-choice')
      for (let menuIndex = 0; menuIndex < 2; menuIndex += 1) {
        await menuChoices.nth(menuIndex).click()
        for (let navigationIndex = 0; navigationIndex < 3; navigationIndex += 1) {
          await navigationChoices.nth(navigationIndex).click()
          await page.waitForTimeout(120)
          interactionAudits.push({
            menuIndex,
            navigationIndex,
            state: await page.evaluate(() => ({
              navbarClass: document.querySelector('.navbar')?.className || '',
              menuTheme: document.documentElement.dataset.adminMenuTheme,
              navType: document.documentElement.dataset.adminNavType
            })),
            navigationContrast: await collectNavigationContrast(page)
          })
        }
      }
    }
  }
  let homeActionAudit = null
  if (shouldExerciseHome) {
    await page.getByRole('button', { name: '进入授权管理' }).click()
    await page.waitForTimeout(500)
    homeActionAudit = {
      url: page.url(),
      pathname: new URL(page.url()).pathname
    }
  }
  writeFileSync(
    resolve(outputRoot, `${label}-theme-defects.json`),
    `${JSON.stringify({
      url: page.url(),
      viewport,
      deviceScaleFactor,
      mode,
      navType,
      menuStyle,
      targetPath,
      computed,
      navigationContrast,
      interactionAudits,
      homeActionAudit,
      defects,
      pageAudits,
      runtime: { pageErrors, consoleErrors, failedRequests }
    }, null, 2)}\n`
  )
  const scanDefectCount = pageAudits.reduce((total, audit) => total + audit.defects.length, 0)
  const interactionDefectCount = interactionAudits.reduce(
    (total, audit) => total + audit.navigationContrast.defects.length,
    0
  )
  const totalDefects = defects.length + scanDefectCount + navigationContrast.defects.length + interactionDefectCount
  process.stdout.write(`QA_CAPTURE_OK ${label} ${totalDefects}\n`)
} finally {
  await browser.close()
}
