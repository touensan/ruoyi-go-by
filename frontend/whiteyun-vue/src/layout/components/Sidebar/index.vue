<template>
  <div :class="['sidebar-theme-wrapper', {'has-logo':showLogo}, sideTheme]" class="sidebar-container">
    <logo v-if="showLogo" :collapse="isCollapse" />
    <el-scrollbar wrap-class="scrollbar-wrapper">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :background-color="getMenuBackground"
        :text-color="getMenuTextColor"
        :unique-opened="true"
        :active-text-color="theme"
        :collapse-transition="false"
        mode="vertical"
        :class="sideTheme"
      >
        <template v-for="(menuRoute, index) in sidebarRouters" :key="menuRoute.path + index">
          <template v-if="isAccountSection(menuRoute)">
            <li class="account-section-divider" :class="{ collapsed: isCollapse }">
              <span v-if="!isCollapse">我的服务</span>
            </li>
            <sidebar-item
              v-for="child in visibleAccountChildren(menuRoute)"
              :key="child.path"
              :item="accountMenuItem(child)"
              :base-path="accountMenuPath(menuRoute.path, child.path)"
            />
          </template>
          <sidebar-item v-else :item="menuRoute" :base-path="menuRoute.path" />
        </template>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import Logo from './Logo.vue'
import SidebarItem from './SidebarItem.vue'
import useAppStore from '@/store/modules/app'
import useSettingsStore from '@/store/modules/settings'
import usePermissionStore from '@/store/modules/permission'

const route = useRoute()
const appStore = useAppStore()
const settingsStore = useSettingsStore()
const permissionStore = usePermissionStore()

const sidebarRouters = computed(() => permissionStore.sidebarRouters)
const showLogo = computed(() => settingsStore.sidebarLogo)
const sideTheme = computed(() => settingsStore.sideTheme)
const theme = computed(() => settingsStore.theme)
const isCollapse = computed(() => !appStore.sidebar.opened)
const isAccountSection = (menuRoute: any) => String(menuRoute.path || '').replace(/^\/|\/$/g, '') === 'account'
const visibleAccountChildren = (menuRoute: any) => (menuRoute.children || []).filter((child: any) => !child.hidden)
const accountMenuPath = (parentPath: string, childPath: string) => {
  const parent = String(parentPath || '').replace(/\/$/g, '')
  const child = String(childPath || '').replace(/^\/+/g, '')
  return `${parent}/${child}`.replace(/\/+/g, '/')
}
const accountMenuItem = (child: any) => ({ ...child, path: '', children: undefined })

// 获取菜单背景色
const getMenuBackground = computed(() => {
  return 'var(--sidebar-bg)'
})

// 获取菜单文字颜色
const getMenuTextColor = computed(() => {
  return 'var(--sidebar-text)'
})

const activeMenu = computed(() => {
  const { meta, path } = route
  if (meta.activeMenu) {
    return meta.activeMenu
  }
  return path
})
</script>

<style lang="scss" scoped>
.sidebar-container {
  background-color: v-bind(getMenuBackground);

  .scrollbar-wrapper {
    background-color: v-bind(getMenuBackground);
  }

  .el-menu {
    border: none;
    height: 100%;
    width: 100% !important;

    .el-menu-item, .el-sub-menu__title {
      &:hover {
        background-color: var(--menu-hover, rgba(0, 0, 0, 0.06)) !important;
      }
    }

    .el-menu-item {
      color: v-bind(getMenuTextColor);

      &.is-active {
        color: var(--menu-active-text, #409eff);
        background-color: var(--menu-hover, rgba(0, 0, 0, 0.06)) !important;
      }
    }

    .el-sub-menu__title {
      color: v-bind(getMenuTextColor);
    }
  }

  .account-section-divider {
    position: relative;
    height: 40px;
    margin: 14px 18px 2px;
    border-top: 1px solid var(--admin-line);
    list-style: none;

    span {
      position: absolute;
      top: 10px;
      left: 0;
      color: var(--sidebar-text);
      font-size: 11px;
      font-weight: 650;
      letter-spacing: .08em;
      opacity: .7;
    }

    &.collapsed {
      height: 18px;
      margin: 12px 12px 0;
    }
  }
}
</style>
