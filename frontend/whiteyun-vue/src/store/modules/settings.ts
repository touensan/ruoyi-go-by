import defaultSettings from '@/settings'
import { useDark } from '@vueuse/core'
import { useDynamicTitle } from '@/utils/dynamicTitle'
import { handleThemeStyle } from '@/utils/theme'

const isDark = useDark({
  storageKey: 'apiauth-theme',
  valueDark: 'dark',
  valueLight: 'light'
})

const { sideTheme, showSettings, navType, tagsView, tagsViewPersist, tagsIcon, tagsViewStyle, fixedHeader, sidebarLogo, dynamicTitle, footerVisible, footerContent } = defaultSettings

const storageSetting = JSON.parse(localStorage.getItem('layout-setting') || '{}') || {}

interface SettingsState {
  title: string
  theme: string
  sideTheme: string
  showSettings: boolean
  navType: number
  tagsView: boolean
  tagsViewPersist: boolean
  tagsViewStyle: string
  tagsIcon: boolean
  fixedHeader: boolean
  sidebarLogo: boolean
  dynamicTitle: boolean
  footerVisible: boolean
  footerContent: string
  isDark: boolean
}

const useSettingsStore = defineStore(
  'settings',
  {
    state: (): SettingsState => ({
      title: '',
      theme: !storageSetting.theme || storageSetting.theme.toUpperCase() === '#409EFF'
        ? '#05a9e8'
        : storageSetting.theme,
      sideTheme: storageSetting.sideTheme || sideTheme,
      showSettings: showSettings,
      navType: storageSetting.navType === undefined ? navType : storageSetting.navType,
      tagsView: storageSetting.tagsView === undefined ? tagsView : storageSetting.tagsView,
      tagsViewPersist: storageSetting.tagsViewPersist === undefined ? tagsViewPersist : storageSetting.tagsViewPersist,
      tagsIcon: storageSetting.tagsIcon === undefined ? tagsIcon : storageSetting.tagsIcon,
      tagsViewStyle: storageSetting.tagsViewStyle === undefined ? tagsViewStyle : storageSetting.tagsViewStyle,
      fixedHeader: storageSetting.fixedHeader === undefined ? fixedHeader : storageSetting.fixedHeader,
      sidebarLogo: storageSetting.sidebarLogo === undefined ? sidebarLogo : storageSetting.sidebarLogo,
      dynamicTitle: storageSetting.dynamicTitle === undefined ? dynamicTitle : storageSetting.dynamicTitle,
      footerVisible: storageSetting.footerVisible === undefined ? footerVisible : storageSetting.footerVisible,
      footerContent: footerContent,
      isDark: isDark.value
    }),
    actions: {
      // 修改布局设置
      changeSetting(data: { key: string; value: any }) {
        const { key, value } = data
        if (this.hasOwnProperty(key)) {
          (this as any)[key] = value
        }
      },
      // 设置网页标题
      setTitle(title: string) {
        this.title = title
        useDynamicTitle()
      },
      // 与官网共用 apiauth-theme，保证同一浏览器中的前后台主题连续。
      setTheme(dark: boolean) {
        this.isDark = dark
        isDark.value = dark
        nextTick(() => {
          handleThemeStyle(this.theme)
        })
      },
      toggleTheme() {
        this.setTheme(!this.isDark)
      }
    }
  })

export default useSettingsStore
