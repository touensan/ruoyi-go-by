<template>
  <div class="app-container account-page" v-loading="loading">
    <section class="account-hero">
      <div><span>用户工作台</span><h1>控制台</h1><p>积分、充值和任务权益都属于当前登录账号；管理员也使用同一套用户工作台。</p></div>
      <div class="balance"><small>可用积分</small><strong>{{ points(overview.pointsMinor) }}</strong><span>积分</span></div>
    </section>
    <section class="metric-grid">
      <div><span>可用积分</span><strong>{{ points(overview.pointsMinor) }}</strong><small>可兑换平台服务</small></div>
      <div><span>可参与任务</span><strong>{{ overview.taskCount || 0 }}</strong><small>完成任务领取兑换码</small></div>
      <div><span>任务记录</span><strong>{{ overview.claimCount || 0 }}</strong><small>等待核验与已发放</small></div>
      <div><span>充值订单</span><strong>{{ overview.rechargeCount || 0 }}</strong><small>人民币充值 1:1 获得积分</small></div>
    </section>
    <section class="quick-grid">
      <button @click="$router.push('/account/marketplace')"><strong>服务市场</strong><small>项目可扩展自己的商品与服务</small></button>
      <button @click="$router.push('/account/exchange')"><strong>兑换中心</strong><small>做任务、领卡密、兑换积分</small></button>
      <button @click="$router.push('/account/points')"><strong>积分与订单</strong><small>充值、订单与积分流水</small></button>
      <button @click="$router.push('/account/profile')"><strong>个人资料</strong><small>维护头像、账号与密码</small></button>
    </section>
  </div>
</template>

<script setup lang="ts" name="AccountConsole">
import { getAccountOverview } from '@/api/platform'
const loading = ref(false)
const overview = ref<any>({})
const points = (value: number) => (Number(value || 0) / 100).toLocaleString('zh-CN')
async function load() {
  loading.value = true
  try { overview.value = (await getAccountOverview()).data || {} }
  finally { loading.value = false }
}
load()
</script>

<style scoped lang="scss">
.account-hero, .metric-grid > div, .quick-grid button { border: 1px solid var(--admin-line); border-radius: 15px; background: var(--admin-surface); }
.account-hero { display: flex; justify-content: space-between; align-items: end; gap: 24px; padding: 28px 30px; }
.account-hero > div:first-child span { color: var(--el-color-primary); font-weight: 700; font-size: 13px; }
.account-hero h1 { margin: 6px 0 8px; font-size: clamp(30px, 4vw, 44px); }
.account-hero p, .balance small, .balance span, .metric-grid span, .metric-grid small, .quick-grid small { color: var(--el-text-color-secondary); }
.balance { min-width: 190px; padding: 16px 18px; border-radius: 12px; background: var(--admin-surface-soft); }
.balance > * { display: block; }.balance strong { margin: 5px 0; font-size: 30px; }
.metric-grid, .quick-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-top: 16px; }
.metric-grid > div { padding: 20px; }.metric-grid span, .metric-grid small { display: block; }.metric-grid strong { display: block; margin: 10px 0 4px; font-size: 28px; }
.quick-grid { grid-template-columns: repeat(2, 1fr); }
.quick-grid button { padding: 22px; text-align: left; color: inherit; cursor: pointer; }.quick-grid strong, .quick-grid small { display: block; }.quick-grid small { margin-top: 7px; }
@media (max-width: 900px) { .metric-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 560px) { .account-hero { align-items: stretch; flex-direction: column; }.balance { min-width: 0; }.metric-grid, .quick-grid { grid-template-columns: 1fr; } }
</style>
