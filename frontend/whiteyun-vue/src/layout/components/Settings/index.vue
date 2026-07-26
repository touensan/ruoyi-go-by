<template>
  <el-drawer v-model="showSettings" :withHeader="false" :lock-scroll="false" direction="rtl" :size="drawerSize" :class="sideTheme">
    <section class="settings-group" aria-labelledby="navigation-title">
      <div class="setting-drawer-title">
        <h3 id="navigation-title" class="drawer-title">菜单导航设置</h3>
        <p>桌面端使用所选布局，手机端自动收纳为侧滑菜单。</p>
      </div>
      <div class="nav-wrap" role="radiogroup" aria-label="菜单导航布局">
        <button
          v-for="option in navigationOptions"
          :key="option.value"
          type="button"
          class="layout-choice"
          :class="{ activeItem: navType === option.value }"
          :aria-pressed="navType === option.value"
          @click="handleNavType(option.value)"
        >
          <span class="layout-preview" :class="`layout-preview--${option.preview}`" aria-hidden="true">
            <i class="preview-header"></i>
            <i class="preview-sidebar"></i>
            <i class="preview-content"></i>
          </span>
          <span>{{ option.label }}</span>
          <el-icon v-if="navType === option.value" class="choice-check"><Check /></el-icon>
        </button>
      </div>
    </section>

    <section class="settings-group" aria-labelledby="menu-style-title">
      <div class="setting-drawer-title">
        <h3 id="menu-style-title" class="drawer-title">菜单风格设置</h3>
        <p>{{ menuStyleDescription }}</p>
      </div>
      <div class="menu-style-grid" role="radiogroup" aria-label="菜单颜色风格">
        <button
          v-for="option in menuStyleOptions"
          :key="option.value"
          type="button"
          class="menu-style-choice"
          :class="{ activeItem: sideTheme === option.value }"
          :aria-pressed="sideTheme === option.value"
          @click="handleTheme(option.value)"
        >
          <span class="menu-style-swatch" :class="option.value" aria-hidden="true"></span>
          <span class="menu-style-copy">
            <strong>{{ option.label }}</strong>
            <small>{{ option.description }}</small>
          </span>
          <el-icon v-if="sideTheme === option.value" class="choice-check"><Check /></el-icon>
        </button>
      </div>
    </section>

    <div class="drawer-item">
      <span>主题颜色</span>
      <span class="comp-style">
        <el-color-picker v-model="theme" :predefine="predefineColors" @change="themeChange"/>
      </span>
    </div>
    <el-divider />

    <h3 class="drawer-title">系统布局配置</h3>

    <div class="drawer-item">
      <span>开启页签</span>
      <span class="comp-style">
        <el-switch v-model="settingsStore.tagsView" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>持久化标签页</span>
      <span class="comp-style">
        <el-switch v-model="settingsStore.tagsViewPersist" :disabled="!settingsStore.tagsView" @change="tagsViewPersistChange" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>显示页签图标</span>
      <span class="comp-style">
        <el-switch v-model="settingsStore.tagsIcon" :disabled="!settingsStore.tagsView" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>标签页样式</span>
      <span class="comp-style">
        <el-radio-group v-model="settingsStore.tagsViewStyle" :disabled="!settingsStore.tagsView" size="small">
          <el-radio-button label="card">卡片</el-radio-button>
          <el-radio-button label="chrome">谷歌</el-radio-button>
        </el-radio-group>
      </span>
    </div>

    <div class="drawer-item">
      <span>固定 Header</span>
      <span class="comp-style">
        <el-switch v-model="settingsStore.fixedHeader" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>显示 Logo</span>
      <span class="comp-style">
        <el-switch v-model="settingsStore.sidebarLogo" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>动态标题</span>
      <span class="comp-style">
        <el-switch v-model="settingsStore.dynamicTitle" @change="dynamicTitleChange" class="drawer-switch" />
      </span>
    </div>

    <div class="drawer-item">
      <span>底部版权</span>
      <span class="comp-style">
        <el-switch v-model="settingsStore.footerVisible" class="drawer-switch" />
      </span>
    </div>

    <el-divider />

    <el-button type="primary" plain icon="DocumentAdd" @click="saveSetting">保存配置</el-button>
    <el-button plain icon="Refresh" @click="resetSetting">重置配置</el-button>
  </el-drawer>

</template>

<script setup lang="ts">
import useAppStore from '@/store/modules/app'
import useSettingsStore from '@/store/modules/settings'
import usePermissionStore from '@/store/modules/permission'
import { handleThemeStyle } from '@/utils/theme'
import { useWindowSize } from '@vueuse/core'

const { proxy } = getCurrentInstance()
const appStore = useAppStore()
const settingsStore = useSettingsStore()
const permissionStore = usePermissionStore()
const showSettings = ref<boolean>(false)
const navType = ref<number>(settingsStore.navType)
const theme = ref<string>(settingsStore.theme)
const sideTheme = ref<string>(settingsStore.sideTheme)
const tagsViewPersist = ref(settingsStore.tagsViewPersist)
const storeSettings = computed(() => settingsStore)
const predefineColors = ref<string[]>(["#409EFF", "#ff4500", "#ff8c00", "#ffd700", "#90ee90", "#00ced1", "#1e90ff", "#c71585"])
const { width } = useWindowSize()
const drawerSize = computed(() => width.value <= 520 ? `${Math.min(width.value * 0.92, 380)}px` : '380px')
const navigationOptions = [
  { value: 1, label: '左侧导航', preview: 'left' },
  { value: 2, label: '混合导航', preview: 'mixed' },
  { value: 3, label: '顶部导航', preview: 'top' }
]
const menuStyleOptions = computed(() => [
  {
    value: 'theme-dark',
    label: '深色菜单',
    description: settingsStore.isDark ? '纯黑层级，更强对比' : '深色菜单，浅色内容区'
  },
  {
    value: 'theme-light',
    label: settingsStore.isDark ? '柔和深色' : '浅色菜单',
    description: settingsStore.isDark ? '深灰层级，降低视觉压迫' : '浅色菜单，与内容区连贯'
  }
])
const menuStyleDescription = computed(() => settingsStore.isDark
  ? '夜间模式下两种菜单都会自动映射为完整深色体系。'
  : '日间模式可选择深色菜单或浅色菜单。'
)

/** 是否需要dynamicTitle */
function dynamicTitleChange(): void {
  useSettingsStore().setTitle(useSettingsStore().title)
}

function tagsViewPersistChange(val: boolean): void {
  settingsStore.tagsViewPersist = val
  tagsViewPersist.value = val
}

function themeChange(val: string): void {
  settingsStore.theme = val
  handleThemeStyle(val)
}

function handleTheme(val: string): void {
  settingsStore.changeSetting({ key: 'sideTheme', value: val })
  sideTheme.value = val
}

function handleNavType(val: number): void {
  settingsStore.navType = val
  navType.value = val
}

function applyNavigationMode(val: number): void {
  if (appStore.device === 'mobile') {
    appStore.closeSideBar({ withoutAnimation: true })
    appStore.toggleSideBarHide(false)
    permissionStore.setSidebarRouters(permissionStore.defaultRoutes)
    return
  }
  if (val === 1) {
    appStore.sidebar.opened = true
    appStore.toggleSideBarHide(false)
  }
  if (val === 2) {
    appStore.sidebar.opened = true
    appStore.toggleSideBarHide(false)
  }
  if (val === 3) {
    appStore.sidebar.opened = false
    appStore.toggleSideBarHide(true)
  }
  if ([1, 3].includes(val)) {
    permissionStore.setSidebarRouters(permissionStore.defaultRoutes)
  }
}

/** 菜单导航设置 */
watch([navType, () => appStore.device], ([val]) => {
  applyNavigationMode(val as number)
}, { immediate: true })

function saveSetting(): void {
  proxy.$modal.loading("正在保存到本地，请稍候...")
  if (!tagsViewPersist.value) {
    proxy.$cache.local.remove('tags-view-visited')
  }
  const layoutSetting = {
    "navType": storeSettings.value.navType,
    "tagsView": storeSettings.value.tagsView,
    "tagsIcon": storeSettings.value.tagsIcon,
    "tagsViewStyle": storeSettings.value.tagsViewStyle,
    "tagsViewPersist": storeSettings.value.tagsViewPersist,
    "fixedHeader": storeSettings.value.fixedHeader,
    "sidebarLogo": storeSettings.value.sidebarLogo,
    "dynamicTitle": storeSettings.value.dynamicTitle,
    "footerVisible": storeSettings.value.footerVisible,
    "sideTheme": storeSettings.value.sideTheme,
    "theme": storeSettings.value.theme
  }
  localStorage.setItem("layout-setting", JSON.stringify(layoutSetting))
  setTimeout(proxy.$modal.closeLoading(), 1000)
}

function resetSetting(): void {
  proxy.$cache.local.remove('tags-view-visited')
  proxy.$modal.loading("正在清除设置缓存并刷新，请稍候...")
  localStorage.removeItem("layout-setting")
  setTimeout(() => {
    window.location.reload()
  }, 1000)
}

function openSetting(): void {
  showSettings.value = true
}

defineExpose({
  openSetting
})
</script>

<style lang='scss' scoped>
.theme-dark {
  --preview-menu-bg: #070707;
  --preview-menu-soft-bg: #15181b;
}

.theme-light {
  --preview-menu-bg: var(--admin-surface-soft);
  --preview-menu-soft-bg: var(--admin-surface-soft);
}

:global(html.dark) .theme-light {
  --preview-menu-bg: #14171a;
  --preview-menu-soft-bg: #14171a;
}

.settings-group {
  margin-bottom: 24px;
}

.setting-drawer-title {
  margin-bottom: 14px;
  color: var(--el-text-color-primary, rgba(0, 0, 0, 0.85));
  line-height: 22px;
  font-weight: bold;

  .drawer-title {
    margin: 0;
    font-size: 14px;
  }

  p {
    margin: 4px 0 0;
    color: var(--admin-muted);
    font-size: 12px;
    font-weight: 400;
    line-height: 18px;
  }
}

.drawer-item {
  color: var(--el-text-color-regular, rgba(0, 0, 0, 0.65));
  padding: 12px 0;
  font-size: 14px;

  .comp-style {
    float: right;
    margin: -3px 8px 0px 0px;
  }
}

.nav-wrap {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.layout-choice,
.menu-style-choice {
  position: relative;
  border: 1px solid var(--admin-line);
  color: var(--admin-text-soft);
  background: var(--admin-surface);
  cursor: pointer;
  transition: border-color .2s ease, background .2s ease, color .2s ease;

  &:hover {
    border-color: var(--admin-line-strong);
    background: var(--admin-surface-hover);
  }

  &.activeItem {
    border-color: var(--admin-primary);
    color: var(--admin-text);
    box-shadow: 0 0 0 1px var(--admin-primary) inset;
  }

  &:focus-visible {
    outline: 2px solid var(--admin-primary);
    outline-offset: 2px;
  }
}

.layout-choice {
  min-width: 0;
  padding: 8px;
  border-radius: 10px;
  font-size: 12px;
  line-height: 18px;
}

.layout-preview {
  position: relative;
  display: block;
  height: 48px;
  margin-bottom: 7px;
  overflow: hidden;
  border: 1px solid var(--admin-line);
  border-radius: 7px;
  background: var(--admin-bg);

  i {
    position: absolute;
    display: block;
  }

  .preview-header,
  .preview-sidebar {
    background: var(--preview-menu-bg);
  }

  .preview-content {
    background: var(--admin-surface);
  }

  .preview-content {
    inset: 13px 4px 4px 22px;
    border-radius: 3px;
  }

  &--left {
    .preview-header {
      inset: 4px 4px auto 22px;
      height: 6px;
      border-radius: 3px;
    }

    .preview-sidebar {
      inset: 0 auto 0 0;
      width: 17px;
    }
  }

  &--mixed {
    .preview-header {
      inset: 0 0 auto 0;
      height: 10px;
    }

    .preview-sidebar {
      inset: 10px auto 0 0;
      width: 17px;
    }
  }

  &--top {
    .preview-header {
      inset: 0 0 auto 0;
      height: 10px;
    }

    .preview-sidebar {
      display: none;
    }

    .preview-content {
      inset: 14px 4px 4px;
    }
  }
}

.menu-style-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}

.menu-style-choice {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 10px;
  text-align: left;
}

.menu-style-swatch {
  flex: 0 0 34px;
  width: 34px;
  height: 38px;
  border: 1px solid var(--admin-line-strong);
  border-radius: 8px;

  &.theme-dark {
    background: #070707;
  }

  &.theme-light {
    background: var(--preview-menu-soft-bg);
  }
}

.menu-style-copy {
  min-width: 0;

  strong,
  small {
    display: block;
  }

  strong {
    color: var(--admin-text);
    font-size: 12px;
    line-height: 18px;
  }

  small {
    margin-top: 2px;
    color: var(--admin-muted);
    font-size: 10px;
    line-height: 15px;
  }
}

.choice-check {
  position: absolute;
  top: 7px;
  right: 7px;
  display: grid;
  width: 18px;
  height: 18px;
  place-items: center;
  border-radius: 50%;
  color: #07151b;
  background: var(--admin-primary);
  font-size: 11px;
}

@media (max-width: 520px) {
  .nav-wrap {
    grid-template-columns: 1fr;
  }

  .layout-choice {
    display: grid;
    grid-template-columns: 74px 1fr;
    align-items: center;
    gap: 10px;
    text-align: left;
  }

  .layout-preview {
    width: 74px;
    margin-bottom: 0;
  }

  .menu-style-grid {
    grid-template-columns: 1fr;
  }
}
</style>
