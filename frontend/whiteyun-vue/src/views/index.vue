<template>
  <div class="app-container home">
    <section class="intro-panel">
      <div class="intro-copy">
        <p class="product-code">Whiteyun Vue · {{ productVersion }}</p>
        <h1>白云 R1-Go</h1>
        <p class="product-summary">
          以 Go 生态和 RuoYi 权限模型为基础，配套持续维护的白云风格 Vue 3 管理前端。
        </p>
        <div class="intro-actions">
          <el-button class="intro-primary-action" type="primary" @click="openProject">
            <span>访问项目仓库</span>
            <el-icon><ArrowRight /></el-icon>
          </el-button>
          <el-tag class="intro-version-tag" effect="plain">当前版本 {{ productVersion }}</el-tag>
        </div>
      </div>

      <div class="runtime-box">
        <div>
          <span>后端基线</span>
          <strong>Go + Gin</strong>
        </div>
        <div>
          <span>数据库基线</span>
          <strong>MySQL 5.7+</strong>
        </div>
        <div>
          <span>演示站点</span>
          <el-link :href="demoOrigin" target="_blank" type="primary">r1-go.whiteyun.com</el-link>
        </div>
      </div>
    </section>

    <section class="stack-section">
      <div class="section-heading">
        <h2>技术栈</h2>
        <p>保留成熟的用户、角色、菜单、配置、日志与代码生成能力，并升级主题和响应式体验。</p>
      </div>

      <el-row :gutter="16">
        <el-col v-for="group in stackGroups" :key="group.title" :xs="24" :md="8">
          <el-card shadow="never" class="stack-card">
            <template #header>
              <span>{{ group.title }}</span>
            </template>
            <div class="stack-list">
              <el-tag v-for="item in group.items" :key="item" effect="plain">{{ item }}</el-tag>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </section>

    <section class="changelog-section">
      <div class="section-heading">
        <h2>更新日志</h2>
        <p>记录 RuoYi-Go BY 与 Whiteyun Vue 的真实版本演进。</p>
      </div>

      <el-scrollbar class="changelog-scroll" max-height="260px">
        <el-collapse v-model="activeLog" accordion class="changelog-collapse">
          <el-collapse-item v-for="log in changelog" :key="log.version" :name="log.version">
            <template #title>
              <div class="changelog-title-row">
                <span class="log-version">{{ log.version }}</span>
                <span class="log-date">{{ log.date }}</span>
                <span class="log-title">{{ log.title }}</span>
                <span class="log-status" :class="{ 'is-current': log.current }">{{ log.status }}</span>
              </div>
            </template>
            <div class="changelog-detail">
              <p>{{ log.summary }}</p>
              <ul>
                <li v-for="item in log.items" :key="item">{{ item }}</li>
              </ul>
            </div>
          </el-collapse-item>
        </el-collapse>
      </el-scrollbar>
    </section>

    <section class="credits">
      <span>开源底座致谢：</span>
      <el-link href="https://github.com/touensan/ruoyi-go-by" target="_blank" type="primary">
        touensan/ruoyi-go-by
      </el-link>
      <span>、</span>
      <el-link href="https://gitcode.com/yangzongzhuan/RuoYi-Vue3/tree/typescript" target="_blank" type="primary">
        RuoYi-Vue3 TypeScript
      </el-link>
    </section>
  </div>
</template>

<script setup lang="ts">
const productVersion = '1.2.0'
const projectUrl = 'https://github.com/touensan/ruoyi-go-by'
const demoOrigin = 'https://r1-go.whiteyun.com/'
const activeLog = ref<string>('')

const stackGroups = [
  {
    title: '后端技术',
    items: ['Go 1.25+', 'Gin 1.12', 'GORM 1.31', 'MySQL 5.7+', '兼容 MySQL 8', 'Redis 5/6/7', 'JWT v5']
  },
  {
    title: '前端技术',
    items: ['Vue 3.5', 'Vite 6', 'TypeScript', 'Element Plus', 'Pinia', 'Vue Router 4', 'Axios']
  },
  {
    title: '部署运行',
    items: ['Nginx 反向代理', '宝塔面板', '单二进制服务', '静态资源发布', '本地动态配置', '备份与回档']
  }
]

const changelog = [
  {
    version: 'v1.2.0',
    date: '2026-07-26',
    title: 'Whiteyun Vue 成为默认前端',
    status: '当前版本',
    current: true,
    summary: '新增独立维护的 Whiteyun Vue 前端，同时保留冻结的 RuoYi Vue 兼容版本。',
    items: [
      '重做管理端日间与夜间配色，统一侧边栏、顶部导航、标签页、表单和弹层。',
      '三种导航布局均可随主题正确切换，桌面、平板和手机端保持可用。',
      '登录页、后台首页、系统配置与代码生成页面共享同一套视觉变量。',
      '默认构建入口切换到 frontend/whiteyun-vue，历史路径继续兼容旧脚本。',
      'frontend/ruoyi-vue 仅作旧视觉兼容保留，不再接收常规功能与界面更新。'
    ]
  },
  {
    version: 'v1.1.2',
    date: '2026-06-14',
    title: '系统配置与工程底座',
    status: '历史版本',
    current: false,
    summary: '补齐站点、支付和邮箱配置，并整理后端运行、构建与部署基础。',
    items: [
      '系统配置集中管理站点信息、支付参数和邮箱参数。',
      '后端基于 Go、Gin、GORM、MySQL、Redis 与 JWT。',
      '管理端继续兼容动态菜单、RBAC、文件上传和代码生成。'
    ]
  }
]

function openProject(): void {
  window.open(projectUrl, '_blank', 'noopener,noreferrer')
}
</script>

<style scoped lang="scss">
.home {
  color: var(--el-text-color-primary);
}

.intro-panel {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 340px;
  gap: 24px;
  padding: 28px 32px;
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
  font-weight: 600;
}

.intro-copy h1 {
  margin: 0;
  font-size: clamp(32px, 3.5vw, 48px);
  line-height: 1.25;
  font-weight: 650;
}

.product-summary {
  max-width: 760px;
  margin: 16px 0 0;
  color: var(--el-text-color-regular);
  font-size: 16px;
  line-height: 1.8;
}

.intro-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin-top: 22px;
}

.intro-primary-action {
  min-height: 40px;
  padding: 0 17px;
  border-radius: 10px;
  font-weight: 600;
  letter-spacing: .01em;

  :deep(.el-icon) {
    margin-left: 2px;
    font-size: 14px;
  }
}

.intro-version-tag {
  --el-tag-bg-color: var(--admin-surface-soft);
  --el-tag-border-color: var(--admin-line-strong);
  --el-tag-text-color: var(--admin-muted);
  height: 32px;
  padding: 0 11px;
  border-color: var(--admin-line-strong) !important;
  border-radius: 8px;
  color: var(--admin-muted) !important;
  background: var(--admin-surface-soft) !important;
}

.runtime-box {
  display: grid;
  gap: 14px;
  padding: 18px;
  border: 1px solid var(--admin-line);
  border-radius: 13px;
  background: var(--admin-surface-soft);
}

.runtime-box div {
  min-width: 0;
}

.runtime-box span {
  display: block;
  margin-bottom: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.runtime-box strong {
  display: block;
  font-size: 17px;
  font-weight: 650;
}

.runtime-box :deep(.el-link__inner) {
  word-break: break-all;
}

.stack-section {
  margin-top: 20px;
}

.changelog-section {
  margin-top: 20px;
}

.section-heading {
  margin-bottom: 14px;
}

.section-heading h2 {
  margin: 0;
  font-size: 22px;
  font-weight: 650;
}

.section-heading p {
  margin: 8px 0 0;
  color: var(--el-text-color-secondary);
  line-height: 1.7;
}

.stack-card {
  height: 100%;
  border-radius: 13px;
}

.stack-card :deep(.el-card__header) {
  font-weight: 650;
}

.stack-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.changelog-scroll {
  border: 1px solid var(--el-border-color-light);
  border-radius: 13px;
  background: var(--el-bg-color);
}

.changelog-collapse {
  border: 0;
}

.changelog-collapse :deep(.el-collapse-item__wrap) {
  border-bottom: 0;
}

.changelog-collapse :deep(.el-collapse-item__header) {
  height: 42px;
  padding: 0 12px 0 16px;
  border-bottom-color: var(--el-border-color-lighter);
  color: var(--el-text-color-primary);
  line-height: 42px;
  transition: background-color 0.2s ease;
}

.changelog-collapse :deep(.el-collapse-item__header:hover) {
  background: var(--el-fill-color-lighter);
}

.changelog-collapse :deep(.el-collapse-item__header.is-active) {
  background: var(--el-color-primary-light-9);
  border-bottom-color: var(--el-color-primary-light-7);
}

.changelog-collapse :deep(.el-collapse-item__arrow) {
  margin-left: 10px;
  color: var(--el-text-color-secondary);
}

.changelog-collapse :deep(.el-collapse-item__content) {
  padding: 10px 16px 14px 16px;
}

.changelog-title-row {
  display: grid;
  grid-template-columns: 76px 104px minmax(0, 1fr) 72px;
  align-items: center;
  gap: 10px;
  width: 100%;
  min-width: 0;
}

.log-version {
  color: var(--el-text-color-primary);
  font-weight: 650;
}

.log-date {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.log-title {
  min-width: 0;
  overflow: hidden;
  color: var(--el-text-color-regular);
  font-size: 14px;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.log-status {
  justify-self: end;
  width: 64px;
  height: 24px;
  border: 1px solid var(--el-border-color);
  border-radius: 7px;
  color: var(--el-text-color-secondary);
  background: var(--el-fill-color-lighter);
  font-size: 12px;
  line-height: 22px;
  text-align: center;
  white-space: nowrap;
}

.log-status.is-current {
  border-color: var(--el-color-success-light-5);
  color: var(--el-color-success);
  background: var(--el-color-success-light-9);
}

.changelog-detail {
  padding-left: 180px;
}

.changelog-detail p {
  margin: 0 0 8px;
  color: var(--el-text-color-secondary);
  line-height: 1.6;
}

.changelog-detail ul {
  display: grid;
  grid-template-columns: repeat(2, minmax(260px, 1fr));
  gap: 6px 18px;
  margin: 0;
  padding-left: 18px;
  color: var(--el-text-color-regular);
  line-height: 1.6;
}

.credits {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  margin-top: 20px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

@media (max-width: 900px) {
  .intro-panel {
    grid-template-columns: 1fr;
    padding: 22px;
  }

  .changelog-title-row {
    grid-template-columns: 72px 96px minmax(0, 1fr);
  }

  .log-status {
    display: none;
  }

  .changelog-detail {
    padding-left: 0;
  }

  .changelog-detail ul {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 520px) {
  .intro-copy h1 {
    font-size: 26px;
  }

  .product-summary {
    font-size: 15px;
  }

  .changelog-title-row {
    grid-template-columns: 70px minmax(0, 1fr);
    gap: 8px;
  }

  .log-date {
    display: none;
  }
}
</style>
