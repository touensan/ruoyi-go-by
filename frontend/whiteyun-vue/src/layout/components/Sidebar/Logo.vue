<template>
  <div class="sidebar-logo-container" :class="{ 'collapse': collapse }">
    <transition name="sidebarLogoFade">
      <router-link v-if="collapse" key="collapse" class="sidebar-logo-link" to="/">
        <span class="sidebar-monogram" :aria-label="appTitle">{{ monogram }}</span>
      </router-link>
      <router-link v-else key="expand" class="sidebar-logo-link" to="/">
        <span class="sidebar-wordmark">{{ appTitle }}</span>
        <span class="sidebar-console">控制台</span>
      </router-link>
    </transition>
  </div>
</template>

<script setup lang="ts">
const appTitle = import.meta.env.VITE_APP_TITLE || 'Whiteyun Vue'
const monogram = appTitle.trim().charAt(0).toUpperCase() || 'W'

defineProps({
  collapse: {
    type: Boolean,
    required: true
  }
})

</script>

<style lang="scss" scoped>
.sidebarLogoFade-enter-active {
  transition: opacity 1.5s;
}

.sidebarLogoFade-enter,
.sidebarLogoFade-leave-to {
  opacity: 0;
}

.sidebar-logo-container {
  position: relative;
  height: 50px;
  line-height: 50px;
  background: var(--sidebar-bg);
  overflow: hidden;
  border-bottom: 1px solid var(--admin-line);

  & .sidebar-logo-link {
    height: 100%;
    width: 100%;

    display: flex;
    align-items: center;
    gap: 9px;
    padding: 0 18px;

    .sidebar-monogram {
      display: grid;
      place-items: center;
      width: 30px;
      height: 30px;
      border: 1px solid var(--admin-line-strong);
      border-radius: 9px;
      color: var(--sidebar-logo-text);
      font-size: 16px;
      font-weight: 700;
    }

    .sidebar-wordmark {
      min-width: 0;
      overflow: hidden;
      color: var(--sidebar-logo-text);
      font-size: 18px;
      font-weight: 650;
      letter-spacing: -0.04em;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .sidebar-console {
      margin-left: auto;
      padding-left: 8px;
      border-left: 1px solid var(--admin-line);
      color: var(--sidebar-text);
      font-size: 11px;
      letter-spacing: .08em;
    }
  }

  &.collapse {
    .sidebar-logo-link {
      justify-content: center;
      padding: 0;
    }
  }
}
</style>
