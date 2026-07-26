<template>
  <div class="app-container home">
    <section class="welcome-panel">
      <div class="welcome-copy">
        <p class="product-code">RuoYi-Go BY · v1.3.0</p>
        <h1>你好，{{ displayName }}</h1>
        <p class="product-summary">
          这里是你的工作台。账户资料与业务服务始终属于用户空间，管理员会在同一套界面中获得额外的管理入口。
        </p>
        <div class="welcome-actions">
          <el-button type="primary" @click="openProfile('userinfo')">完善个人资料</el-button>
          <el-button @click="openProfile('resetPwd')">账户安全</el-button>
        </div>
      </div>

      <div class="identity-card">
        <span>当前身份</span>
        <strong>{{ identityLabel }}</strong>
        <el-tag :type="isAdministrator ? 'primary' : 'info'" effect="light">
          {{ isAdministrator ? '用户空间 + 管理权限' : '用户空间' }}
        </el-tag>
      </div>
    </section>

    <section class="workspace-section">
      <div class="section-heading">
        <h2>我的账户</h2>
        <p>所有账号共享的基础入口，普通用户不会看到系统管理功能。</p>
      </div>

      <div class="action-grid">
        <button class="action-card" type="button" @click="openProfile('userinfo')">
          <span class="action-icon"><svg-icon icon-class="user" /></span>
          <span class="action-copy">
            <strong>个人资料</strong>
            <small>维护昵称、邮箱、手机号和头像</small>
          </span>
          <span class="action-arrow">→</span>
        </button>

        <button class="action-card" type="button" @click="openProfile('resetPwd')">
          <span class="action-icon"><svg-icon icon-class="password" /></span>
          <span class="action-copy">
            <strong>账户安全</strong>
            <small>更新登录密码并检查账户信息</small>
          </span>
          <span class="action-arrow">→</span>
        </button>
      </div>
    </section>

    <section v-if="isAdministrator" class="workspace-section">
      <div class="section-heading">
        <h2>管理入口</h2>
        <p>管理能力叠加在用户工作台之上，不建立另一套割裂的后台。</p>
      </div>

      <div class="action-grid admin-grid">
        <button
          v-for="entry in adminEntries"
          :key="entry.path"
          class="action-card"
          type="button"
          @click="router.push(entry.path)"
        >
          <span class="action-icon"><svg-icon :icon-class="entry.icon" /></span>
          <span class="action-copy">
            <strong>{{ entry.title }}</strong>
            <small>{{ entry.description }}</small>
          </span>
          <span class="action-arrow">→</span>
        </button>
      </div>
    </section>

    <section class="architecture-note">
      <div>
        <strong>统一身份架构</strong>
        <p>后续项目把面向客户的功能加入“我的账户”，把运营能力加入管理菜单，即可沿用同一套权限模型。</p>
      </div>
      <el-tag effect="plain">v1.3.0</el-tag>
    </section>
  </div>
</template>

<script setup lang="ts">
import useUserStore from '@/store/modules/user'

const router = useRouter()
const userStore = useUserStore()

const displayName = computed(() => userStore.nickName || userStore.name || '用户')
const isAdministrator = computed(() =>
  userStore.permissions.includes('*:*:*') || userStore.roles.includes('admin')
)
const identityLabel = computed(() => isAdministrator.value ? '管理员' : '普通用户')

const adminEntries = [
  {
    title: '用户管理',
    description: '管理用户、状态与角色分配',
    icon: 'user',
    path: '/system/user'
  },
  {
    title: '角色与权限',
    description: '配置管理角色与菜单权限',
    icon: 'peoples',
    path: '/system/role'
  },
  {
    title: '系统配置',
    description: '维护站点、支付与邮件配置',
    icon: 'system',
    path: '/system/settings'
  }
]

function openProfile(activeTab: 'userinfo' | 'resetPwd'): void {
  if (activeTab === 'userinfo') {
    router.push({ name: 'UserProfile' })
    return
  }

  router.push({ name: 'Profile', params: { activeTab } })
}
</script>

<style scoped lang="scss">
.home {
  color: var(--el-text-color-primary);
}

.welcome-panel {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 280px;
  gap: 24px;
  padding: 30px 32px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 16px;
  background:
    radial-gradient(circle at 82% 18%, var(--admin-accent-glow), transparent 34%),
    var(--el-bg-color);
}

.product-code {
  margin: 0 0 8px;
  color: var(--el-color-primary);
  font-size: 13px;
  font-weight: 650;
}

.welcome-copy h1 {
  margin: 0;
  font-size: clamp(32px, 3.5vw, 48px);
  line-height: 1.25;
  font-weight: 680;
}

.product-summary {
  max-width: 760px;
  margin: 16px 0 0;
  color: var(--el-text-color-regular);
  font-size: 16px;
  line-height: 1.8;
}

.welcome-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 22px;
}

.identity-card {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  gap: 10px;
  padding: 20px;
  border: 1px solid var(--admin-line);
  border-radius: 13px;
  background: var(--admin-surface-soft);
}

.identity-card > span {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.identity-card strong {
  font-size: 24px;
  font-weight: 680;
}

.workspace-section {
  margin-top: 22px;
}

.section-heading {
  margin-bottom: 14px;
}

.section-heading h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 680;
}

.section-heading p {
  margin: 8px 0 0;
  color: var(--el-text-color-secondary);
  line-height: 1.7;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.admin-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.action-card {
  display: grid;
  grid-template-columns: 42px minmax(0, 1fr) auto;
  align-items: center;
  gap: 14px;
  min-width: 0;
  padding: 18px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 13px;
  color: inherit;
  text-align: left;
  background: var(--el-bg-color);
  cursor: pointer;
  transition: border-color .2s ease, box-shadow .2s ease, transform .2s ease;
}

.action-card:hover {
  border-color: var(--el-color-primary-light-5);
  box-shadow: 0 10px 28px rgba(15, 23, 42, .08);
  transform: translateY(-2px);
}

.action-card:focus-visible {
  outline: 3px solid var(--el-color-primary-light-7);
  outline-offset: 2px;
}

.action-icon {
  display: grid;
  width: 42px;
  height: 42px;
  place-items: center;
  border-radius: 11px;
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
  font-size: 20px;
}

.action-copy {
  min-width: 0;
}

.action-copy strong,
.action-copy small {
  display: block;
}

.action-copy strong {
  font-size: 16px;
  font-weight: 650;
}

.action-copy small {
  margin-top: 6px;
  overflow: hidden;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.5;
  text-overflow: ellipsis;
}

.action-arrow {
  color: var(--el-text-color-secondary);
  font-size: 18px;
}

.architecture-note {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  margin-top: 24px;
  padding: 18px 20px;
  border: 1px solid var(--admin-line);
  border-radius: 13px;
  background: var(--admin-surface-soft);
}

.architecture-note strong {
  font-weight: 650;
}

.architecture-note p {
  margin: 6px 0 0;
  color: var(--el-text-color-secondary);
  line-height: 1.6;
}

@media (max-width: 900px) {
  .welcome-panel {
    grid-template-columns: 1fr;
  }

  .identity-card {
    display: grid;
    grid-template-columns: 1fr auto;
  }

  .identity-card strong {
    grid-column: 1;
  }

  .identity-card .el-tag {
    grid-column: 2;
    grid-row: 1 / span 2;
    align-self: center;
  }

  .admin-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .home {
    padding: 12px;
  }

  .welcome-panel {
    gap: 18px;
    padding: 22px 18px;
    border-radius: 14px;
  }

  .welcome-copy h1 {
    font-size: 30px;
  }

  .product-summary {
    font-size: 14px;
  }

  .welcome-actions {
    display: grid;
    grid-template-columns: 1fr;
  }

  .welcome-actions .el-button {
    width: 100%;
    margin-left: 0;
  }

  .action-grid {
    grid-template-columns: 1fr;
  }

  .action-card {
    padding: 16px;
  }

  .architecture-note {
    align-items: flex-start;
  }
}

@media (prefers-reduced-motion: reduce) {
  .action-card {
    transition: none;
  }
}
</style>
