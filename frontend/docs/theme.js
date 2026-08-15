(() => {
  const storageKey = 'ruoyi-go-by-theme'
  const root = document.documentElement
  const systemTheme = window.matchMedia('(prefers-color-scheme: dark)')

  const readStoredTheme = () => {
    try {
      const value = window.localStorage.getItem(storageKey)
      return value === 'light' || value === 'dark' ? value : null
    } catch {
      return null
    }
  }

  const updateThemeColor = (theme) => {
    let meta = document.querySelector('meta[name="theme-color"]')
    if (!meta) {
      meta = document.createElement('meta')
      meta.setAttribute('name', 'theme-color')
      document.head.appendChild(meta)
    }
    meta.setAttribute('content', theme === 'dark' ? '#000000' : '#f7f7f6')
  }

  const updateControls = (theme) => {
    document.querySelectorAll('[data-theme-option]').forEach((button) => {
      const active = button.dataset.themeOption === theme
      button.setAttribute('aria-pressed', String(active))
      button.classList.toggle('is-active', active)
    })
  }

  const applyTheme = (theme, persist = false) => {
    root.dataset.theme = theme
    root.style.colorScheme = theme
    updateThemeColor(theme)
    updateControls(theme)

    if (persist) {
      try {
        window.localStorage.setItem(storageKey, theme)
      } catch {
        // Storage may be disabled; the active page still keeps the selected theme.
      }
    }
  }

  applyTheme(readStoredTheme() || (systemTheme.matches ? 'dark' : 'light'))

  document.addEventListener('DOMContentLoaded', () => {
    updateControls(root.dataset.theme)
    document.querySelectorAll('[data-theme-option]').forEach((button) => {
      button.addEventListener('click', () => applyTheme(button.dataset.themeOption, true))
    })
  })

  const handleSystemTheme = (event) => {
    if (!readStoredTheme()) applyTheme(event.matches ? 'dark' : 'light')
  }

  if (typeof systemTheme.addEventListener === 'function') {
    systemTheme.addEventListener('change', handleSystemTheme)
  } else {
    systemTheme.addListener(handleSystemTheme)
  }

  window.addEventListener('storage', (event) => {
    if (event.key !== storageKey) return
    const theme = event.newValue === 'light' || event.newValue === 'dark'
      ? event.newValue
      : (systemTheme.matches ? 'dark' : 'light')
    applyTheme(theme)
  })
})()
