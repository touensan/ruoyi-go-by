<template>
  <div class="login">
    <header class="login-header">
      <a href="/" class="login-brand" aria-label="返回 APIAuth 官网">APIAuth</a>
      <nav class="login-nav" aria-label="辅助导航">
        <a href="/">官网首页</a>
        <a href="/docs/">接入文档</a>
        <div class="theme-segmented" role="group" aria-label="界面主题">
          <button type="button" :class="{ active: !settingsStore.isDark }" @click="settingsStore.setTheme(false)">日间</button>
          <button type="button" :class="{ active: settingsStore.isDark }" @click="settingsStore.setTheme(true)">夜间</button>
        </div>
      </nav>
    </header>

    <main class="login-shell">
      <section class="login-story">
        <div class="story-copy">
          <span class="release-badge"><i></i> 授权服务运行中</span>
          <p class="eyebrow">APIAuth 授权控制中心</p>
          <h1>软件授权<br>一处管理</h1>
          <p class="story-summary">创建应用与授权方案，签发许可证，并管理激活、续期、暂停与吊销。</p>
          <div class="capability-list" aria-label="产品能力">
            <span>在线激活</span>
            <span>离线授权</span>
            <span>三态验权</span>
            <span>Ed25519</span>
          </div>
        </div>
        <div class="cube-visual" aria-hidden="true">
          <img :src="cubeVisualUrl" alt="">
          <div class="visual-status">
            <span><i></i>许可证运营台</span>
            <strong>运行正常</strong>
          </div>
        </div>
      </section>

      <section class="login-panel">
        <div class="login-panel-heading">
          <p class="eyebrow">管理后台</p>
          <h2>登录控制台</h2>
          <p>使用已开通的管理账号继续。</p>
        </div>
        <el-form ref="loginRef" :model="loginForm" :rules="loginRules" class="login-form">
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              type="text"
              size="large"
              auto-complete="username"
              placeholder="账号"
            >
              <template #prefix><svg-icon icon-class="user" class="el-input__icon input-icon" /></template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              size="large"
              auto-complete="current-password"
              placeholder="密码"
              show-password
              @keyup.enter="handleLogin"
            >
              <template #prefix><svg-icon icon-class="password" class="el-input__icon input-icon" /></template>
            </el-input>
          </el-form-item>
          <el-form-item prop="code" v-if="captchaEnabled" class="captcha-row">
            <el-input
              v-model="loginForm.code"
              size="large"
              auto-complete="off"
              placeholder="验证码"
              @keyup.enter="handleLogin"
            >
              <template #prefix><svg-icon icon-class="validCode" class="el-input__icon input-icon" /></template>
            </el-input>
            <button type="button" class="login-code" aria-label="刷新验证码" @click="getCode">
              <img :src="codeUrl" class="login-code-img" alt="验证码"/>
            </button>
          </el-form-item>
          <el-checkbox v-model="loginForm.rememberMe" class="remember-account">记住账号</el-checkbox>
          <el-form-item class="submit-row">
            <el-button :loading="loading" size="large" type="primary" @click.prevent="handleLogin">
              <span v-if="!loading">进入控制台</span>
              <span v-else>正在登录...</span>
            </el-button>
            <router-link v-if="register" class="link-type" :to="'/register'">立即注册</router-link>
          </el-form-item>
        </el-form>
        <p class="login-note">账号由平台管理员开通；许可证用户无需登录此处。</p>
      </section>
    </main>

    <footer class="el-login-footer">
      <span>{{ footerContent }}</span>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { getCodeImg } from "@/api/login"
import Cookies from "js-cookie"
import useUserStore from '@/store/modules/user'
import useSettingsStore from '@/store/modules/settings'
import defaultSettings from '@/settings'
import cubeVisualUrl from '../../../site-resend-preview/static/cube-fallback.jpg'
import type { CaptchaInfoResult } from '@/types/api/login'
import type { LoginForm } from '@/types/api/login'

const footerContent = defaultSettings.footerContent
const userStore = useUserStore()
const settingsStore = useSettingsStore()
const route = useRoute()
const router = useRouter()
const { proxy } = getCurrentInstance() as any

const loginForm = ref<LoginForm>({
  username: "",
  password: "",
  rememberMe: false,
  code: "",
  uuid: ""
})

const loginRules = {
  username: [{ required: true, trigger: "blur", message: "请输入您的账号" }],
  password: [{ required: true, trigger: "blur", message: "请输入您的密码" }],
  code: [{ required: true, trigger: "change", message: "请输入验证码" }]
}

const codeUrl = ref("")
const loading = ref(false)
// 验证码开关
const captchaEnabled = ref(true)
// 注册开关
const register = ref(false)
const redirect = ref<string | undefined>(undefined)

watch(route, (newRoute: any) => {
    redirect.value = (newRoute.query && newRoute.query.redirect) as string | undefined
}, { immediate: true })

function handleLogin(): void {
  proxy.$refs.loginRef.validate((valid: boolean) => {
    if (valid) {
      loading.value = true
      // 仅记住用户名；密码不得进入浏览器持久化存储。
      if (loginForm.value.rememberMe) {
        Cookies.set("username", loginForm.value.username, { expires: 30 })
        Cookies.set("rememberMe", loginForm.value.rememberMe, { expires: 30 })
        Cookies.remove("password")
      } else {
        // 否则移除
        Cookies.remove("username")
        Cookies.remove("password")
        Cookies.remove("rememberMe")
      }
      // 调用action的登录方法
      userStore.login(loginForm.value).then(() => {
        const query = route.query
        const otherQueryParams = Object.keys(query).reduce((acc: Record<string, any>, cur) => {
          if (cur !== "redirect") {
            acc[cur] = query[cur]
          }
          return acc
        }, {})
        router.push({ path: redirect.value || "/", query: otherQueryParams })
      }).catch(() => {
        loading.value = false
        // 重新获取验证码
        if (captchaEnabled.value) {
          getCode()
        }
      })
    }
  })
}

function getCode(): void {
  getCodeImg().then(res => {
    captchaEnabled.value = res.captchaEnabled === undefined ? true : res.captchaEnabled
    if (captchaEnabled.value) {
      codeUrl.value = "data:image/gif;base64," + res.img
      loginForm.value.uuid = res.uuid
    }
  })
}

function getCookie(): void {
  const username = Cookies.get("username")
  const rememberMe = Cookies.get("rememberMe")
  // 清理旧版本使用可逆前端密钥保存的密码 Cookie。
  Cookies.remove("password")
  loginForm.value = {
    username: username === undefined ? loginForm.value.username : username,
    password: "",
    rememberMe: rememberMe === undefined ? false : Boolean(rememberMe),
    code: "",
    uuid: ""
  }
}

getCode()
getCookie()
</script>

<style lang='scss' scoped>
.login {
  min-height: 100%;
  color: var(--admin-text);
  background:
    radial-gradient(circle at 72% 24%, var(--admin-accent-glow), transparent 32%),
    var(--admin-bg);
  overflow: auto;
}

.login-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: min(1240px, calc(100% - 48px));
  height: 78px;
  margin: 0 auto;
}

.login-brand {
  color: var(--admin-text);
  font-size: 25px;
  font-weight: 650;
  letter-spacing: -.05em;
}

.login-nav {
  display: flex;
  align-items: center;
  gap: 24px;
  color: var(--admin-muted);
  font-size: 14px;

  a:hover {
    color: var(--admin-text);
  }
}

.theme-segmented {
  display: flex;
  gap: 2px;
  padding: 3px;
  border: 1px solid var(--admin-line);
  border-radius: 11px;
  background: var(--admin-surface-soft);

  button {
    height: 30px;
    padding: 0 12px;
    border: 0;
    border-radius: 8px;
    color: var(--admin-muted);
    background: transparent;
    font: inherit;
    cursor: pointer;
  }

  button.active {
    color: var(--admin-text);
    background: var(--admin-surface);
    box-shadow: 0 1px 3px var(--admin-shadow);
  }
}

.login-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(360px, .8fr);
  width: min(1240px, calc(100% - 48px));
  min-height: calc(100vh - 138px);
  margin: 0 auto;
  border: 1px solid var(--admin-line);
  border-radius: 24px;
  background: var(--admin-surface);
  overflow: hidden;
}

.login-story {
  position: relative;
  display: grid;
  align-items: center;
  min-height: 660px;
  padding: 70px;
  background: var(--admin-hero-bg);
  overflow: hidden;
}

.story-copy {
  position: relative;
  z-index: 2;
  max-width: 560px;
}

.release-badge {
  display: inline-flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 42px;
  padding: 8px 12px;
  border: 1px solid var(--admin-line);
  border-radius: 999px;
  color: var(--admin-muted);
  background: var(--admin-surface-soft);
  font-size: 13px;

  i,
  .visual-status i {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: #42d392;
    box-shadow: 0 0 12px rgba(66, 211, 146, .65);
  }
}

.eyebrow {
  margin: 0 0 14px;
  color: var(--admin-primary);
  font-size: 12px;
  font-weight: 650;
  letter-spacing: .12em;
}

.story-copy h1 {
  margin: 0;
  color: var(--admin-text);
  font-size: clamp(56px, 6.2vw, 90px);
  font-weight: 520;
  letter-spacing: -.075em;
  line-height: .96;
}

.story-summary {
  max-width: 520px;
  margin: 34px 0 0;
  color: var(--admin-muted);
  font-size: 18px;
  line-height: 1.75;
}

.capability-list {
  display: flex;
  flex-wrap: wrap;
  gap: 9px;
  margin-top: 34px;

  span {
    padding: 8px 12px;
    border: 1px solid var(--admin-line);
    border-radius: 9px;
    color: var(--admin-muted);
    background: var(--admin-surface-soft);
    font-size: 13px;
  }
}

.cube-visual {
  position: absolute;
  right: -9%;
  bottom: -16%;
  width: min(64%, 580px);
  opacity: var(--admin-cube-opacity);
  filter: var(--admin-cube-filter);
  transform: rotate(-2deg);

  img {
    display: block;
    width: 100%;
    mix-blend-mode: var(--admin-cube-blend);
  }
}

.visual-status {
  position: absolute;
  top: 16%;
  left: -14%;
  display: flex;
  justify-content: space-between;
  width: 300px;
  padding: 15px 18px;
  border: 1px solid var(--admin-line-strong);
  border-radius: 13px;
  color: var(--admin-text);
  background: var(--admin-glass);
  backdrop-filter: blur(16px);
  font-size: 12px;

  span {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  strong {
    color: #42d392;
    font-weight: 550;
  }
}

.login-panel {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: clamp(42px, 6vw, 86px);
  border-left: 1px solid var(--admin-line);
  background: var(--admin-surface);
}

.login-panel-heading {
  margin-bottom: 34px;

  h2 {
    margin: 0;
    color: var(--admin-text);
    font-size: 34px;
    font-weight: 580;
    letter-spacing: -.045em;
  }

  > p:last-child {
    margin: 12px 0 0;
    color: var(--admin-muted);
    line-height: 1.6;
  }
}

.login-form {
  width: 100%;

  :deep(.el-input__wrapper) {
    min-height: 48px;
  }

  .input-icon {
    width: 15px;
    height: 18px;
  }
}

.captcha-row :deep(.el-form-item__content) {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 118px;
  gap: 10px;
}

.login-code {
  height: 48px;
  padding: 3px;
  border: 1px solid var(--admin-line);
  border-radius: 10px;
  background: var(--admin-surface-soft);
  cursor: pointer;
  overflow: hidden;
}

.el-login-footer {
  width: min(1240px, calc(100% - 48px));
  min-height: 60px;
  margin: 0 auto;
  padding: 20px 0;
  text-align: center;
  color: var(--admin-muted);
  font-size: 12px;
}

.login-code-img {
  display: block;
  width: 100%;
  height: 40px;
  object-fit: cover;
  border-radius: 7px;
}

.remember-account {
  margin: 2px 0 24px;
}

.submit-row {
  margin-bottom: 0;

  :deep(.el-form-item__content),
  :deep(.el-button) {
    width: 100%;
  }
}

.login-note {
  margin: 20px 0 0;
  padding-top: 20px;
  border-top: 1px solid var(--admin-line);
  color: var(--admin-muted);
  font-size: 12px;
  line-height: 1.6;
}

@media (max-width: 980px) {
  .login-shell {
    grid-template-columns: 1fr;
  }

  .login-story {
    min-height: 400px;
    padding: 48px;
  }

  .story-copy h1 {
    font-size: 60px;
  }

  .story-summary {
    max-width: 480px;
  }

  .cube-visual {
    right: -8%;
    bottom: -42%;
    width: 62%;
  }

  .login-panel {
    border-top: 1px solid var(--admin-line);
    border-left: 0;
  }
}

@media (max-width: 640px) {
  .login-header {
    width: calc(100% - 28px);
    height: 68px;
  }

  .login-nav > a {
    display: none;
  }

  .login-shell {
    width: calc(100% - 28px);
    border-radius: 18px;
  }

  .login-story {
    min-height: 328px;
    padding: 34px 26px;
  }

  .release-badge {
    margin-bottom: 24px;
  }

  .story-copy h1 {
    font-size: 48px;
  }

  .story-summary {
    margin-top: 22px;
    font-size: 15px;
  }

  .capability-list {
    margin-top: 22px;
  }

  .cube-visual,
  .visual-status {
    display: none;
  }

  .login-panel {
    padding: 34px 26px;
  }

  .login-panel-heading h2 {
    font-size: 29px;
  }

  .el-login-footer {
    width: calc(100% - 28px);
  }
}
</style>
