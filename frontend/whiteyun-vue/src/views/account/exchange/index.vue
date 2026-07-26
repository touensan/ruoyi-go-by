<template>
  <div class="app-container exchange-page" v-loading="loading">
    <section class="exchange-head">
      <div><span>任务权益</span><h1>兑换中心</h1><p>完成任务领取可转让兑换码，也可以输入已有兑换码充值积分。</p></div>
      <el-tag :type="center.enabled ? 'success' : 'info'" round>{{ center.enabled ? '服务已开放' : '暂未开放' }}</el-tag>
    </section>
    <section class="redeem-card">
      <div><h2>兑换积分</h2><p>兑换码只可使用一次，兑换成功后积分立即进入当前账户。</p></div>
      <div class="redeem-form"><el-input v-model="redeemCode" size="large" placeholder="R1POINT-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX" @keyup.enter="redeem" /><el-button type="primary" size="large" :loading="redeeming" @click="redeem">立即兑换</el-button></div>
    </section>
    <section class="tasks">
      <div class="section-title"><div><span>赚取积分</span><h2>可参与任务</h2></div><small>雨云是可选提供商，项目管理员可以独立开关。</small></div>
      <el-empty v-if="!center.tasks?.length" description="当前没有任务" />
      <div v-else class="task-grid">
        <article v-for="task in center.tasks" :key="task.publicId">
          <div class="task-meta"><el-tag effect="plain">{{ task.provider === 'RAINYUN' ? '雨云任务' : '平台任务' }}</el-tag><strong>+{{ task.rewardPoints }} 积分</strong></div>
          <h3>{{ task.name }}</h3><p>{{ task.summary }}</p><div class="requirements">{{ task.requirements }}</div>
          <div class="task-actions">
            <el-button v-if="task.actionUrl" tag="a" :href="task.actionUrl" target="_blank">前往完成</el-button>
            <el-button type="primary" :disabled="!task.available || !!task.claimStatus" @click="openClaim(task)">{{ task.claimStatus ? statusLabel(task.claimStatus) : task.available ? '提交核验' : '等待开放' }}</el-button>
          </div>
        </article>
      </div>
    </section>
    <section v-if="center.claims?.length" class="history">
      <h2>我的任务记录</h2>
      <el-table :data="center.claims"><el-table-column prop="taskName" label="任务" min-width="180" /><el-table-column prop="providerSubject" label="核验账号" min-width="140" /><el-table-column label="奖励" width="110"><template #default="{ row }">{{ row.rewardPoints }} 积分</template></el-table-column><el-table-column label="状态" width="120"><template #default="{ row }"><el-tag>{{ statusLabel(row.status) }}</el-tag></template></el-table-column><el-table-column prop="codeMask" label="兑换码" min-width="220" /></el-table>
    </section>
    <el-dialog v-model="claimOpen" :title="selectedTask?.name || '提交任务'" width="500px">
      <el-alert v-if="selectedTask?.provider === 'RAINYUN'" title="请输入通过活动入口注册的雨云 UID，系统将通过雨云接口核验下级关系。" type="info" :closable="false" />
      <el-form label-position="top" class="claim-form"><el-form-item :label="selectedTask?.provider === 'RAINYUN' ? '雨云 UID' : '任务凭证'"><el-input v-model="providerSubject" /></el-form-item></el-form>
      <template #footer><el-button @click="claimOpen = false">取消</el-button><el-button type="primary" :loading="claiming" @click="submit">提交核验</el-button></template>
    </el-dialog>
    <el-dialog v-model="codeOpen" title="请保存兑换码" width="560px">
      <el-alert title="完整兑换码只显示一次，可以自行兑换或转交他人。" type="warning" :closable="false" />
      <el-input class="code-result" :model-value="issuedCode" readonly><template #append><el-button @click="copyCode">复制</el-button></template></el-input>
      <template #footer><el-button @click="closeCode">我已保存</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts" name="AccountExchange">
import { getExchangeCenter, redeemExchangeCode, submitRewardClaim } from '@/api/platform'
const loading = ref(false), redeeming = ref(false), claiming = ref(false), claimOpen = ref(false), codeOpen = ref(false)
const center = ref<any>({ enabled: false, tasks: [], claims: [] }), redeemCode = ref(''), providerSubject = ref(''), issuedCode = ref(''), selectedTask = ref<any>()
const statusLabel = (value: string) => ({ PENDING: '等待审核', ISSUED: '已发放', REJECTED: '未通过' } as Record<string, string>)[value] || value
async function load() { loading.value = true; try { center.value = (await getExchangeCenter()).data || center.value } finally { loading.value = false } }
function openClaim(task: any) { selectedTask.value = task; providerSubject.value = ''; claimOpen.value = true }
async function submit() {
  if (!providerSubject.value.trim()) return
  claiming.value = true
  try {
    const result = await submitRewardClaim(selectedTask.value.publicId, providerSubject.value.trim())
    claimOpen.value = false
    if (result.data?.code) { issuedCode.value = result.data.code; codeOpen.value = true }
    ElMessage.success(result.data?.code ? '核验通过，兑换码已生成' : '任务已提交'); await load()
  } finally { claiming.value = false }
}
async function redeem() {
  if (!redeemCode.value.trim()) return
  redeeming.value = true
  try { const result = await redeemExchangeCode(redeemCode.value.trim()); ElMessage.success(`${result.data.points} 积分已到账`); redeemCode.value = ''; await load() }
  finally { redeeming.value = false }
}
async function copyCode() { await navigator.clipboard.writeText(issuedCode.value); ElMessage.success('兑换码已复制') }
function closeCode() { issuedCode.value = ''; codeOpen.value = false }
load()
</script>
<style scoped lang="scss">
.exchange-head, .redeem-card, .task-grid article, .history { border: 1px solid var(--admin-line); border-radius: 15px; background: var(--admin-surface); }
.exchange-head, .redeem-card { display: flex; justify-content: space-between; align-items: center; gap: 24px; padding: 26px 28px; }.exchange-head span, .section-title span { color: var(--el-color-primary); font-size: 13px; font-weight: 700; }.exchange-head h1 { margin: 6px 0; font-size: 40px; }.exchange-head p, .redeem-card p, .task-grid p, .section-title small { color: var(--el-text-color-secondary); }
.redeem-card { margin-top: 16px; }.redeem-card h2 { margin: 0 0 6px; }.redeem-form { display: flex; gap: 10px; min-width: 48%; }.tasks { margin-top: 22px; }.section-title { display: flex; justify-content: space-between; align-items: end; margin-bottom: 14px; }.section-title h2 { margin: 4px 0 0; }
.task-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }.task-grid article { display: flex; flex-direction: column; padding: 21px; }.task-grid h3 { margin: 17px 0 7px; }.task-grid p { min-height: 44px; line-height: 1.6; }.task-meta, .task-actions { display: flex; justify-content: space-between; align-items: center; gap: 10px; }.task-meta strong { color: var(--el-color-primary); }.requirements { margin: 14px 0; padding: 12px; border-radius: 10px; background: var(--admin-surface-soft); color: var(--el-text-color-secondary); line-height: 1.55; }.task-actions { margin-top: auto; justify-content: flex-end; }.history { margin-top: 18px; padding: 20px; }.claim-form, .code-result { margin-top: 18px; }
@media (max-width: 900px) { .task-grid { grid-template-columns: 1fr; }.redeem-card { align-items: stretch; flex-direction: column; }.redeem-form { min-width: 0; } }
@media (max-width: 560px) { .exchange-head, .section-title { align-items: flex-start; flex-direction: column; }.redeem-form { flex-direction: column; } }
</style>
