<template>
  <div class="app-container exchange-admin">
    <section class="admin-head"><div><span>任务与卡密</span><h1>兑换中心</h1><p>雨云只是可选任务提供商；核心任务、卡密和积分结构保持通用。</p></div><el-button type="primary" icon="Tickets" @click="generateOpen = true">生成兑换码</el-button></section>
    <section class="admin-panel">
      <el-tabs v-model="tab" @tab-change="loadTab"><el-tab-pane label="任务" name="tasks" /><el-tab-pane label="审核" name="claims" /><el-tab-pane label="兑换码" name="codes" /><el-tab-pane label="集成设置" name="settings" /></el-tabs>
      <template v-if="tab === 'tasks'">
        <div class="toolbar"><el-button type="primary" icon="Plus" @click="openTask()">新增任务</el-button></div>
        <el-table :data="tasks"><el-table-column label="任务" min-width="220"><template #default="{ row }"><strong>{{ row.name }}</strong><small>{{ row.taskCode }}</small></template></el-table-column><el-table-column label="提供商" width="110"><template #default="{ row }">{{ row.provider === 'RAINYUN' ? '雨云' : '人工' }}</template></el-table-column><el-table-column label="奖励" width="110"><template #default="{ row }">{{ row.rewardPoints }} 积分</template></el-table-column><el-table-column label="状态" width="100"><template #default="{ row }"><el-tag :type="row.status === 'ACTIVE' ? 'success' : 'info'">{{ row.status }}</el-tag></template></el-table-column><el-table-column label="操作" width="90"><template #default="{ row }"><el-button link type="primary" @click="openTask(row)">编辑</el-button></template></el-table-column></el-table>
      </template>
      <template v-else-if="tab === 'claims'">
        <el-table :data="claims"><el-table-column label="用户" width="130"><template #default="{ row }">@{{ row.username }}</template></el-table-column><el-table-column prop="taskName" label="任务" min-width="180" /><el-table-column prop="providerSubject" label="核验账号" min-width="140" /><el-table-column prop="status" label="状态" width="110" /><el-table-column label="操作" width="150"><template #default="{ row }"><template v-if="row.status === 'PENDING'"><el-button link type="success" @click="review(row, true)">通过</el-button><el-button link type="danger" @click="review(row, false)">拒绝</el-button></template></template></el-table-column></el-table>
        <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="loadClaims" />
      </template>
      <template v-else-if="tab === 'codes'">
        <el-table :data="codes"><el-table-column prop="codeMask" label="兑换码掩码" min-width="250" /><el-table-column prop="rewardPoints" label="积分" width="100" /><el-table-column prop="sourceType" label="来源" width="100" /><el-table-column prop="status" label="状态" width="110" /><el-table-column prop="ownerUsername" label="领取用户" min-width="130" /><el-table-column prop="redeemedBy" label="兑换用户" min-width="130" /><el-table-column prop="description" label="说明" min-width="180" /></el-table>
        <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="loadCodes" />
      </template>
      <template v-else>
        <el-form label-position="top" class="settings-form">
          <div class="switch-row"><div><strong>开放兑换中心</strong><p>控制用户端任务与兑换入口。</p></div><el-switch v-model="settings.enabled" /></div>
          <div class="switch-row"><div><strong>启用雨云任务</strong><p>关闭时不影响人工任务和兑换码。</p></div><el-switch v-model="settings.rainyunEnabled" /></div>
          <div class="grid"><el-form-item label="雨云 API 地址"><el-input v-model="settings.rainyunApiBaseUrl" /></el-form-item><el-form-item label="雨云活动入口"><el-input v-model="settings.rainyunInviteUrl" /></el-form-item></div>
          <el-form-item label="雨云 X-Api-Key"><el-input v-model="settings.rainyunApiKey" type="password" show-password :placeholder="settings.rainyunKeyConfigured ? '已配置，留空保持不变' : '请输入密钥'" /></el-form-item>
          <el-checkbox v-model="settings.clearRainyunKey">清除已保存密钥</el-checkbox>
          <div class="setting-actions"><el-button :disabled="!settings.rainyunEnabled" @click="testConnection">检测接口</el-button><el-button type="primary" :loading="saving" @click="saveSettings">保存设置</el-button></div>
        </el-form>
      </template>
    </section>
    <el-dialog v-model="taskOpen" :title="taskForm.publicId ? '编辑任务' : '新增任务'" width="660px">
      <el-form label-position="top"><div class="grid"><el-form-item label="任务编码"><el-input v-model="taskForm.taskCode" /></el-form-item><el-form-item label="名称"><el-input v-model="taskForm.name" /></el-form-item><el-form-item label="提供商"><el-select v-model="taskForm.provider"><el-option label="人工" value="MANUAL" /><el-option label="雨云" value="RAINYUN" /></el-select></el-form-item><el-form-item label="核验方式"><el-select v-model="taskForm.verifyMode"><el-option label="人工审核" value="MANUAL" /><el-option label="雨云下级关系" value="RAINYUN_SUBUSER" /></el-select></el-form-item><el-form-item label="奖励积分"><el-input-number v-model="taskForm.rewardPoints" :min="1" /></el-form-item><el-form-item label="状态"><el-radio-group v-model="taskForm.status"><el-radio-button value="ACTIVE">上架</el-radio-button><el-radio-button value="INACTIVE">下架</el-radio-button></el-radio-group></el-form-item></div><el-form-item label="简介"><el-input v-model="taskForm.summary" type="textarea" /></el-form-item><el-form-item label="活动入口"><el-input v-model="taskForm.actionUrl" /></el-form-item><el-form-item label="完成要求"><el-input v-model="taskForm.requirements" type="textarea" /></el-form-item></el-form>
      <template #footer><el-button @click="taskOpen = false">取消</el-button><el-button type="primary" :loading="saving" @click="saveTask">保存</el-button></template>
    </el-dialog>
    <el-dialog v-model="generateOpen" title="生成兑换码" width="500px">
      <el-form label-position="top"><div class="grid"><el-form-item label="每个兑换码积分"><el-input-number v-model="generateForm.rewardPoints" :min="1" /></el-form-item><el-form-item label="数量"><el-input-number v-model="generateForm.count" :min="1" :max="100" /></el-form-item></div><el-form-item label="说明"><el-input v-model="generateForm.description" /></el-form-item></el-form>
      <template #footer><el-button @click="generateOpen = false">取消</el-button><el-button type="primary" :loading="saving" @click="generate">生成</el-button></template>
    </el-dialog>
    <el-dialog v-model="resultOpen" title="请保存兑换码" width="650px"><el-alert title="完整兑换码只显示一次；服务端仅保存不可逆摘要。" type="warning" :closable="false" /><el-input class="code-result" :model-value="generated.join('\n')" type="textarea" :rows="Math.min(12, generated.length + 1)" readonly /><template #footer><el-button @click="copyGenerated">复制全部</el-button><el-button type="primary" @click="resultOpen = false">我已保存</el-button></template></el-dialog>
  </div>
</template>
<script setup lang="ts" name="PlatformExchange">
import { createTask, generateCodes, getPlatformSettings, listClaims, listCodes, listTasks, reviewClaim, savePlatformSettings, testRainyun, updateTask } from '@/api/platform'
const tab = ref('tasks'), saving = ref(false), taskOpen = ref(false), generateOpen = ref(false), resultOpen = ref(false), total = ref(0)
const tasks = ref<any[]>([]), claims = ref<any[]>([]), codes = ref<any[]>([]), generated = ref<string[]>([])
const query = reactive({ pageNum: 1, pageSize: 10, username: '', status: '' })
const settings = reactive<any>({ enabled: true, rainyunEnabled: false, rainyunApiBaseUrl: 'https://api.rainyun.com', rainyunInviteUrl: '', rainyunApiKey: '', rainyunKeyConfigured: false, clearRainyunKey: false })
const taskForm = reactive<any>({ publicId: '', taskCode: '', name: '', summary: '', provider: 'MANUAL', verifyMode: 'MANUAL', rewardPoints: 10, actionUrl: '', requirements: '', displayOrder: 100, status: 'ACTIVE' })
const generateForm = reactive<any>({ rewardPoints: 20, count: 1, description: '平台积分兑换码', expiresAt: null })
async function loadTasks() { tasks.value = (await listTasks()).data || [] }
async function loadClaims() { const result = await listClaims(query); claims.value = result.rows || []; total.value = result.total || 0 }
async function loadCodes() { const result = await listCodes(query); codes.value = result.rows || []; total.value = result.total || 0 }
async function loadSettings() { Object.assign(settings, (await getPlatformSettings()).data || {}, { rainyunApiKey: '', clearRainyunKey: false }) }
async function loadTab() { query.pageNum = 1; total.value = 0; if (tab.value === 'tasks') await loadTasks(); if (tab.value === 'claims') await loadClaims(); if (tab.value === 'codes') await loadCodes(); if (tab.value === 'settings') await loadSettings() }
function openTask(row?: any) { Object.assign(taskForm, row || { publicId: '', taskCode: '', name: '', summary: '', provider: 'MANUAL', verifyMode: 'MANUAL', rewardPoints: 10, actionUrl: '', requirements: '', displayOrder: 100, status: 'ACTIVE' }); taskOpen.value = true }
async function saveTask() { saving.value = true; try { const payload = { ...taskForm }; delete payload.publicId; taskForm.publicId ? await updateTask(taskForm.publicId, payload) : await createTask(payload); taskOpen.value = false; ElMessage.success('任务已保存'); await loadTasks() } finally { saving.value = false } }
async function review(row: any, approved: boolean) { const { value } = await ElMessageBox.prompt('填写审核备注', approved ? '通过任务' : '拒绝任务'); const result = await reviewClaim(row.publicId, { approved, note: value || '' }); if (result.data?.code) { generated.value = [result.data.code]; resultOpen.value = true }; await loadClaims() }
async function generate() { saving.value = true; try { const result = await generateCodes(generateForm); generated.value = (result.data || []).map((item: any) => item.code); generateOpen.value = false; resultOpen.value = true; await loadCodes() } finally { saving.value = false } }
async function saveSettings() { saving.value = true; try { Object.assign(settings, (await savePlatformSettings(settings)).data || {}, { rainyunApiKey: '', clearRainyunKey: false }); ElMessage.success('设置已保存') } finally { saving.value = false } }
async function testConnection() { const result = await testRainyun(); ElMessage.success(`雨云接口正常：${result.data.account}`) }
async function copyGenerated() { await navigator.clipboard.writeText(generated.value.join('\n')); ElMessage.success('已复制') }
Promise.all([loadTasks(), loadClaims(), loadCodes(), loadSettings()])
</script>
<style scoped lang="scss">
.admin-head, .admin-panel { border: 1px solid var(--admin-line); border-radius: 14px; background: var(--admin-surface); }.admin-head { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 24px 26px; }.admin-head span { color: var(--el-color-primary); font-size: 13px; font-weight: 700; }.admin-head h1 { margin: 5px 0; }.admin-head p, .switch-row p { margin: 0; color: var(--el-text-color-secondary); }.admin-panel { margin-top: 16px; padding: 4px 20px 20px; }.toolbar, .setting-actions { display: flex; justify-content: flex-end; gap: 10px; margin-bottom: 14px; }.el-table small { display: block; color: var(--el-text-color-secondary); }.settings-form { max-width: 850px; padding: 10px 4px; }.switch-row { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 18px 0; border-bottom: 1px solid var(--admin-line); }.switch-row strong { display: block; margin-bottom: 5px; }.grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 14px; margin-top: 16px; }.setting-actions { margin-top: 22px; }.code-result { margin-top: 18px; font-family: monospace; }
@media (max-width: 650px) { .admin-head { align-items: flex-start; flex-direction: column; }.grid { grid-template-columns: 1fr; } }
</style>
