<template>
  <div class="app-container home">
    <section class="intro-panel">
      <div class="intro-copy">
        <p class="product-code">COMMERCIAL LICENSE CENTER</p>
        <h1>APIAuth</h1>
        <p class="product-summary">
          面向 WordPress 插件、Java 投票程序与 DNS 分发程序的统一商业授权管理后台。
        </p>
        <div class="intro-actions">
          <el-button type="primary" @click="goAuthorization">进入授权管理</el-button>
          <el-tag type="success" effect="plain">当前版本 {{ productVersion }}</el-tag>
        </div>
      </div>

      <div class="runtime-box">
        <div>
          <span>运行环境</span>
          <strong>生产试运行</strong>
        </div>
        <div>
          <span>公开协议</span>
          <strong>/api/license/v1</strong>
        </div>
        <div>
          <span>管理域名</span>
          <el-link :href="adminOrigin" target="_blank" type="primary">apiauth.whiteyun.com</el-link>
        </div>
      </div>
    </section>

    <section class="stack-section">
      <div class="section-heading">
        <h2>技术栈</h2>
        <p>后端负责许可证、签名租约与公开协议，管理端负责商业目录、许可证、运行态、密钥和审计。</p>
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
        <p>记录 APIAuth 的版本演进，只保留本项目的真实变更。</p>
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
const router = useRouter()
const productVersion = '1.0.0-rc.1'
const adminOrigin = 'https://apiauth.whiteyun.com/admin/'
const activeLog = ref<string>('')

const stackGroups = [
  {
    title: '后端技术',
    items: ['Go 1.26.5', 'Gin 1.12', 'GORM 1.31', 'MySQL 5.7+', '兼容 MySQL 8', 'Redis 7', 'Ed25519', 'JWT v5']
  },
  {
    title: '前端技术',
    items: ['Vue 3.5', 'Vite 6', 'TypeScript', 'Element Plus', 'Pinia', 'Vue Router 4', 'Axios']
  },
  {
    title: '部署运行',
    items: ['宝塔与 Nginx', 'systemd', '单二进制服务', '静态资源内置', '本机数据服务', '签名备份与恢复']
  }
]

const changelog = [
  {
    version: '1.0.0-rc.1',
    date: '2026-07-25',
    title: '生产试运行候选',
    status: '当前候选',
    current: true,
    summary: '统一授权后台、公开协议、三语言 SDK、离线签名与生产运维能力进入试运行验收。',
    items: [
      '管理应用、版本、权益、套餐、许可证、运行态、签名密钥和授权审计。',
      '公开授权 API 支持激活、刷新、解绑、公钥集、离线请求与更新检查。',
      'Go、PHP/WordPress、Java SDK 共享 Ed25519 签名租约与三态验证语义。',
      '许可证与安装凭据只存摘要，完整许可证仅在受控签发或替换响应中显示一次。',
      'Nginx、systemd、MySQL 8、Redis 7、备份恢复与监控基线已部署；程序以 MySQL 5.7+ 为兼容基准。'
    ]
  },
  {
    version: 'foundation',
    date: '2026-07-25',
    title: '授权平台工程底座',
    status: '已替代',
    current: false,
    summary: '在 ruoyi-go-by 与 RuoYi-Vue3 TypeScript 基础上建立 APIAuth 独立商业授权领域。',
    items: [
      '管理 JWT/RBAC 与客户端许可证身份域保持隔离。',
      'MySQL 版本化 migration 是数据库结构唯一权威，生产禁用自动迁移。',
      '应用、套餐、许可证与签名租约的关键写入保持事务、幂等与审计一致。'
    ]
  }
]

function goAuthorization(): void {
  router.push('/auth/catalog')
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
  border-radius: 8px;
  background: var(--el-bg-color);
}

.product-code {
  margin: 0 0 8px;
  color: var(--el-color-primary);
  font-size: 13px;
  font-weight: 600;
}

.intro-copy h1 {
  margin: 0;
  font-size: 32px;
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

.runtime-box {
  display: grid;
  gap: 14px;
  padding: 18px;
  border-radius: 8px;
  background: var(--el-fill-color-lighter);
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
  border-radius: 8px;
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
  border-radius: 8px;
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
  border-radius: 4px;
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
