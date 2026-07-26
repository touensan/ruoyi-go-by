<template>
  <div class="app-container points-page" v-loading="loading">
    <section class="points-head">
      <div><span>预存积分</span><h1>积分与订单</h1><p>充值时 1 元兑换 1 积分，积分到账后用于兑换平台内的服务。</p></div>
      <div class="points-balance"><small>可用积分</small><strong>{{ pointText(overview.pointsMinor) }}</strong></div>
    </section>
    <section class="recharge-card">
      <div><h2>充值积分</h2><p>支付由系统配置中的易支付渠道处理，回调验签成功后自动入账。</p></div>
      <el-button type="primary" size="large" @click="rechargeOpen = true">立即充值</el-button>
    </section>
    <section class="history-card">
      <el-tabs v-model="tab">
        <el-tab-pane label="积分流水" name="ledger" />
        <el-tab-pane label="充值订单" name="orders" />
      </el-tabs>
      <el-table v-if="tab === 'ledger'" :data="ledger">
        <el-table-column label="类型" width="130"><template #default="{ row }">{{ ledgerType(row.entryType) }}</template></el-table-column>
        <el-table-column label="积分变动" width="140"><template #default="{ row }"><strong :class="row.amountMinor >= 0 ? 'positive' : 'negative'">{{ signed(row.amountMinor) }}</strong></template></el-table-column>
        <el-table-column label="变动后积分" width="150"><template #default="{ row }">{{ pointText(row.balanceAfterMinor) }}</template></el-table-column>
        <el-table-column prop="description" label="说明" min-width="220" />
        <el-table-column label="时间" min-width="180"><template #default="{ row }">{{ parseTime(row.createdAt) }}</template></el-table-column>
      </el-table>
      <el-table v-else :data="orders">
        <el-table-column prop="orderNo" label="订单号" min-width="230" />
        <el-table-column label="积分" width="110"><template #default="{ row }">{{ row.points }}</template></el-table-column>
        <el-table-column label="支付方式" width="120"><template #default="{ row }">{{ row.payType === 'alipay' ? '支付宝' : '微信支付' }}</template></el-table-column>
        <el-table-column label="状态" width="110"><template #default="{ row }"><el-tag :type="row.status === 'PAID' ? 'success' : 'warning'">{{ row.status === 'PAID' ? '已到账' : '待支付' }}</el-tag></template></el-table-column>
        <el-table-column label="创建时间" min-width="180"><template #default="{ row }">{{ parseTime(row.createdAt) }}</template></el-table-column>
      </el-table>
    </section>
    <el-dialog v-model="rechargeOpen" title="充值积分" width="500px" :close-on-click-modal="false">
      <p class="recharge-note">充值金额与积分数量按 1:1 计算；积分属于站内权益。</p>
      <div class="presets"><button v-for="value in [10, 50, 100, 300, 500, 1000]" :key="value" :class="{ active: rechargePoints === value }" @click="rechargePoints = value"><strong>{{ value }}</strong><span>积分</span></button></div>
      <el-form label-position="top">
        <el-form-item label="充值积分"><el-input-number v-model="rechargePoints" :min="1" :max="100000" style="width: 100%" /></el-form-item>
        <el-form-item label="支付方式"><el-radio-group v-model="payType"><el-radio-button value="alipay">支付宝</el-radio-button><el-radio-button value="wxpay">微信支付</el-radio-button></el-radio-group></el-form-item>
      </el-form>
      <template #footer><el-button @click="rechargeOpen = false">取消</el-button><el-button type="primary" :loading="recharging" @click="recharge">支付并充值 {{ rechargePoints }} 积分</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="AccountPoints">
import { createAccountRecharge, getAccountOverview, listAccountLedger, listAccountRecharges } from '@/api/platform'
import { parseTime } from '@/utils/ruoyi'
const route = useRoute()
const loading = ref(false), recharging = ref(false), rechargeOpen = ref(false)
const tab = ref('ledger'), rechargePoints = ref(100), payType = ref('alipay')
const overview = ref<any>({}), ledger = ref<any[]>([]), orders = ref<any[]>([])
const pointText = (value: number) => (Number(value || 0) / 100).toLocaleString('zh-CN')
const signed = (value: number) => `${value > 0 ? '+' : ''}${pointText(value)}`
const ledgerType = (value: string) => ({ RECHARGE: '充值', CODE_REDEMPTION: '兑换码', ADJUSTMENT: '人工调整', PURCHASE: '服务兑换' } as Record<string, string>)[value] || value
async function load() {
  loading.value = true
  try {
    const [a, b, c] = await Promise.all([getAccountOverview(), listAccountLedger(), listAccountRecharges()])
    overview.value = a.data || {}; ledger.value = b.data || []; orders.value = c.data || []
  } finally { loading.value = false }
}
async function recharge() {
  recharging.value = true
  try {
    const result = await createAccountRecharge(rechargePoints.value, payType.value)
    if (!/^https?:\/\//i.test(result.data.payInfo || '')) throw new Error('支付地址无效')
    window.open(result.data.payInfo, '_blank', 'noopener,noreferrer')
    rechargeOpen.value = false; await load()
  } finally { recharging.value = false }
}
if (route.query.payment === 'success') ElMessage.success('支付结果已返回，积分已自动入账')
load()
</script>

<style scoped lang="scss">
.points-head, .recharge-card, .history-card { border: 1px solid var(--admin-line); border-radius: 15px; background: var(--admin-surface); }
.points-head, .recharge-card { display: flex; align-items: center; justify-content: space-between; gap: 24px; padding: 26px 28px; }
.points-head span { color: var(--el-color-primary); font-size: 13px; font-weight: 700; }.points-head h1 { margin: 6px 0; font-size: 38px; }.points-head p, .recharge-card p, .recharge-note { margin: 0; color: var(--el-text-color-secondary); }
.points-balance { min-width: 180px; padding: 16px; border-radius: 12px; background: var(--admin-surface-soft); }.points-balance small, .points-balance strong { display: block; }.points-balance strong { margin-top: 6px; font-size: 30px; }
.recharge-card { margin-top: 16px; }.recharge-card h2 { margin: 0 0 7px; }.history-card { margin-top: 16px; padding: 4px 20px 20px; }.positive { color: var(--el-color-success); }.negative { color: var(--el-color-danger); }
.presets { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin: 18px 0; }.presets button { padding: 13px; border: 1px solid var(--admin-line); border-radius: 10px; color: inherit; background: var(--admin-surface); }.presets button.active { border-color: var(--el-color-primary); color: var(--el-color-primary); }.presets strong, .presets span { display: block; }
@media (max-width: 600px) { .points-head, .recharge-card { align-items: stretch; flex-direction: column; }.points-balance { min-width: 0; }.presets { grid-template-columns: repeat(2, 1fr); } }
</style>
