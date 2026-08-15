const menuButton = document.querySelector('.docs-menu')
const sidebar = document.querySelector('.docs-sidebar')
const searchInput = document.querySelector('.search-box input')
const navLinks = [...document.querySelectorAll('.docs-nav a')]
const sections = [...document.querySelectorAll('.doc-section')]

menuButton?.addEventListener('click', () => {
  const open = sidebar.classList.toggle('is-open')
  menuButton.setAttribute('aria-expanded', String(open))
  menuButton.setAttribute('aria-label', open ? '关闭文档导航' : '打开文档导航')
})

navLinks.forEach((link) => {
  link.addEventListener('click', () => {
    sidebar.classList.remove('is-open')
    menuButton?.setAttribute('aria-expanded', 'false')
    menuButton?.setAttribute('aria-label', '打开文档导航')
  })
})

const updateSearch = () => {
  const query = searchInput.value.trim().toLowerCase()
  navLinks.forEach((link) => {
    const target = document.querySelector(link.getAttribute('href'))
    const searchable = `${link.textContent} ${target?.dataset.search || ''}`.toLowerCase()
    link.classList.toggle('is-hidden', Boolean(query) && !searchable.includes(query))
  })
}

searchInput?.addEventListener('input', updateSearch)

document.addEventListener('keydown', (event) => {
  if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
    event.preventDefault()
    searchInput?.focus()
  }
})

const activeObserver = new IntersectionObserver(
  (entries) => {
    const visible = entries.find((entry) => entry.isIntersecting)
    if (!visible) return
    navLinks.forEach((link) =>
      link.classList.toggle('active', link.getAttribute('href') === `#${visible.target.id}`)
    )
  },
  { rootMargin: '-20% 0px -70% 0px' }
)

sections.forEach((section) => activeObserver.observe(section))

document.querySelectorAll('.copy-code').forEach((button) => {
  button.addEventListener('click', async () => {
    const code = button.parentElement?.querySelector('code')?.textContent || ''
    try {
      await navigator.clipboard.writeText(code)
      button.textContent = '已复制'
      window.setTimeout(() => {
        button.textContent = '复制'
      }, 1600)
    } catch {
      button.textContent = '请手动复制'
    }
  })
})
