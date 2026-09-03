<template>
  <div class="app-container home">
    <section class="intro-panel">
      <div class="intro-copy">
        <p class="product-code">ruoyi-go-by</p>
        <h1>RuoYi-Go BY</h1>
        <p class="product-summary">
          这是一个基于 Go 语言和若依框架开发的 RuoYi 风格开源项目。
        </p>
        <div class="intro-actions">
          <el-button type="primary" @click="goTarget(officialSite)">访问官网</el-button>
          <el-tag type="success" effect="plain">当前版本 v{{ productVersion }}</el-tag>
        </div>
      </div>

      <div class="runtime-box">
        <div>
          <span>项目名称</span>
          <strong>RuoYi-Go BY</strong>
        </div>
        <div>
          <span>英文代号</span>
          <strong>ruoyi-go-by</strong>
        </div>
        <div>
          <span>官网地址</span>
          <el-link :href="officialSite" target="_blank" type="primary">{{ officialSite }}</el-link>
        </div>
      </div>
    </section>

    <section class="stack-section">
      <div class="section-heading">
        <h2>技术栈</h2>
        <p>后端以 Go 生态为核心，前端保留 Vue3 TypeScript 管理端，部署面向常见宝塔/Nginx/MySQL/Redis 环境。</p>
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
        <p>记录 RuoYi-Go BY 的版本演进，只保留本项目的真实变更。</p>
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
      <span>致谢：</span>
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
const productVersion = '1.1.2'
const officialSite = 'https://github.com/touensan/ruoyi-go-by'
const activeLog = ref<string>('')

const stackGroups = [
  {
    title: '后端技术',
    items: ['Go 1.25+', 'Gin 1.12', 'GORM 1.31', 'MySQL 5.7+', 'Redis 5/6/7', 'JWT v5', 'excelize']
  },
  {
    title: '前端技术',
    items: ['Vue 3.5', 'Vite 6', 'TypeScript', 'Element Plus', 'Pinia', 'Vue Router 4', 'Axios']
  },
  {
    title: '部署运行',
    items: ['Nginx 反向代理', '宝塔面板', '单二进制服务', '静态资源内置发布', '本地动态配置', '备份归档目录']
  }
]

const changelog = [
  {
    version: 'v1.1.2',
    date: '2026-06-14',
    title: '系统配置与商业底座能力',
    status: '当前版本',
    current: true,
    summary: '新增系统配置入口，补齐站点、支付和邮箱配置能力，并优化支付配置页的后台表单体验。',
    items: [
      '系统配置已加入系统管理侧边栏第一项，使用站点配置、支付配置、邮箱配置三个横向标签页。',
      '站点配置支持标题、Logo、图标、描述、关键词、前端头部代码、备案、客服邮箱和版权信息。',
      '支付配置支持易支付 V1/V2，支付方式通过复选框启用支付宝支付和微信支付。',
      '支付功能测试真实请求上游商户、下单、查询接口，V2 会额外尝试关闭测试订单。',
      '支付发起和测试结果不依赖浏览器弹出新窗口，便于后续业务页使用当前页二维码或当前标签页跳转。',
      '邮箱配置支持 SMTP、QQ 邮箱、Gmail 和自定义服务器，邮箱测试会真实发送邮件。',
      '支付配置表单改为紧凑选项条、全宽商户字段和居中动作区，提升 RuoYi 管理端一致性。'
    ]
  },
  {
    version: 'v1.1.1',
    date: '2026-06-14',
    title: '底座准备与品牌清理',
    status: '历史版本',
    current: false,
    summary: '完成 RuoYi-Go BY 当前底座准备、技术栈主流化、运行兼容和接口补齐。',
    items: [
      '完成 RuoYi-Go BY 项目目录、源代码、运行配置和文档整理。',
      '接入 ruoyi-go 后端和 RuoYi-Vue3-ts 管理端，发布入口统一为 /admin。',
      '后端依赖升级到主流长期维护版本，覆盖 Gin、GORM、go-redis、JWT、excelize 等核心库。',
      '按 MySQL 5.7 和 Redis 5/6/7 做向下兼容，降低常见宝塔环境部署门槛。',
      '补齐前端管理端需要的公告、监控、缓存、在线用户、上传下载、任务和代码生成兼容接口。',
      '完成 RuoYi-Go BY 首页、项目入口、技术栈、致谢链接和默认数据清理。'
    ]
  }
]

function goTarget(url: string): void {
  window.open(url, '_blank')
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
