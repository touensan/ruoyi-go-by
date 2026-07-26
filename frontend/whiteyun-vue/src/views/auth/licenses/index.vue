<template>
  <div class="app-container">
    <el-form ref="queryRef" :model="query" :inline="true" v-show="showSearch">
      <el-form-item label="许可证掩码" prop="keyMask"><el-input v-model="query.keyMask" placeholder="仅支持掩码查询" clearable @keyup.enter="search" /></el-form-item>
      <el-form-item label="客户 ID" prop="customerPublicId"><el-input v-model="query.customerPublicId" placeholder="客户公开 ID" clearable @keyup.enter="search" /></el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="query.status" clearable placeholder="全部状态" style="width: 160px">
          <el-option v-for="item in ['PENDING','ACTIVE','SUSPENDED','EXPIRED','REVOKED','REPLACED']" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item><el-button type="primary" icon="Search" @click="search">搜索</el-button><el-button icon="Refresh" @click="reset">重置</el-button></el-form-item>
    </el-form>
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5"><el-button type="primary" plain icon="Plus" @click="openIssue" v-hasPermi="['auth:license:issue']">签发许可证</el-button></el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="load" />
    </el-row>
    <el-table v-loading="loading" :data="rows">
      <el-table-column label="许可证" prop="keyMask" min-width="180" />
      <el-table-column label="应用" prop="applicationCode" width="150" />
      <el-table-column label="客户 ID" prop="customerPublicId" min-width="210" show-overflow-tooltip />
      <el-table-column label="状态" width="100"><template #default="{ row }"><auth-status-tag :status="row.status" /></template></el-table-column>
      <el-table-column label="生效时间" width="180"><template #default="{ row }">{{ row.startsAt ? parseTime(row.startsAt) : '-' }}</template></el-table-column>
      <el-table-column label="到期时间" width="180"><template #default="{ row }">{{ row.expiresAt ? parseTime(row.expiresAt) : '永久' }}</template></el-table-column>
      <el-table-column label="操作" width="310" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" @click="showDetail(row)" v-hasPermi="['auth:license:query']">详情</el-button>
          <el-button v-if="row.status === 'PENDING'" link type="primary" :loading="busy === row.publicId" @click="simpleAction(row, 'approve')" v-hasPermi="['auth:license:issue']">批准</el-button>
          <el-button v-if="['ACTIVE','SUSPENDED'].includes(row.status)" link type="primary" :loading="busy === row.publicId" @click="simpleAction(row, 'renew')" v-hasPermi="['auth:license:renew']">续期</el-button>
          <el-button v-if="row.status === 'ACTIVE'" link type="warning" :loading="busy === row.publicId" @click="openReason(row, 'suspend')" v-hasPermi="['auth:license:suspend']">暂停</el-button>
          <el-button v-if="row.status === 'SUSPENDED'" link type="success" :loading="busy === row.publicId" @click="simpleAction(row, 'resume')" v-hasPermi="['auth:license:resume']">恢复</el-button>
          <el-dropdown v-if="!['REVOKED','REPLACED'].includes(row.status)" trigger="click">
            <el-button link type="danger">更多<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
            <template #dropdown><el-dropdown-menu>
              <el-dropdown-item @click="openReason(row, 'replace')" v-hasPermi="['auth:license:replace']">替换并生成新许可证</el-dropdown-item>
              <el-dropdown-item divided @click="openReason(row, 'revoke')" v-hasPermi="['auth:license:revoke']">永久吊销</el-dropdown-item>
            </el-dropdown-menu></template>
          </el-dropdown>
        </template>
      </el-table-column>
      <template #empty><el-empty description="暂无许可证" /></template>
    </el-table>
    <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="load" />

    <el-dialog v-model="issueOpen" title="签发许可证" width="560px" :close-on-click-modal="false">
      <el-alert title="提交后完整许可证只显示一次；页面不会把它写入 URL 或持久状态。" type="warning" :closable="false" show-icon class="mb16" />
      <el-form ref="issueRef" :model="issueForm" :rules="issueRules" label-width="100px" autocomplete="off">
        <el-form-item label="客户 ID" prop="customerPublicId"><el-input v-model="issueForm.customerPublicId" autocomplete="off" placeholder="客户公开 ID" /></el-form-item>
        <el-form-item label="授权套餐" prop="planPublicId">
          <el-select v-model="issueForm.planPublicId" filterable style="width: 100%" placeholder="选择有效套餐">
            <el-option v-for="plan in plans" :key="plan.publicId" :label="`${plan.applicationCode} / ${plan.name}`" :value="plan.publicId" />
          </el-select>
        </el-form-item>
        <el-form-item label="生效时间"><el-date-picker v-model="issueForm.startsAt" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZ" style="width: 100%" placeholder="留空则立即生效" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="issueOpen = false">取消</el-button><el-button type="primary" :loading="submitting" @click="submitIssue">确认签发</el-button></template>
    </el-dialog>

    <el-dialog v-model="reasonOpen" :title="actionTitle" width="520px" :close-on-click-modal="false">
      <el-alert :title="actionWarning" type="warning" :closable="false" show-icon class="mb16" />
      <el-form ref="reasonRef" :model="reasonForm" :rules="{ reason: [{ required: true, message: '请输入操作原因', trigger: 'blur' }] }" label-width="80px">
        <el-form-item label="许可证"><el-input :model-value="selected?.keyMask" disabled /></el-form-item>
        <el-form-item label="原因" prop="reason"><el-input v-model="reasonForm.reason" type="textarea" maxlength="128" show-word-limit /></el-form-item>
      </el-form>
      <template #footer><el-button @click="reasonOpen = false">取消</el-button><el-button type="danger" :loading="!!busy" @click="submitReasonAction">确认执行</el-button></template>
    </el-dialog>

    <el-drawer v-model="detailOpen" title="许可证详情" size="520px">
      <el-descriptions v-if="detail" :column="1" border>
        <el-descriptions-item label="许可证">{{ detail.keyMask }}</el-descriptions-item>
        <el-descriptions-item label="公开 ID">{{ detail.publicId }}</el-descriptions-item>
        <el-descriptions-item label="客户 ID">{{ detail.customerPublicId }}</el-descriptions-item>
        <el-descriptions-item label="套餐 ID">{{ detail.planPublicId }}</el-descriptions-item>
        <el-descriptions-item label="状态"><auth-status-tag :status="detail.status" /></el-descriptions-item>
        <el-descriptions-item label="吊销序列">{{ detail.revocationSerial }}</el-descriptions-item>
        <el-descriptions-item label="生效">{{ detail.startsAt ? parseTime(detail.startsAt) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="到期">{{ detail.expiresAt ? parseTime(detail.expiresAt) : '永久' }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
    <one-time-license-dialog v-model="secretOpen" :license-key="secret" @cleared="secret = ''" />
  </div>
</template>

<script setup lang="ts" name="AuthLicenses">
import { actOnAuthLicense, getAuthLicense, issueAuthLicense, listAuthLicenses, listAuthPlans } from '@/api/auth'
import type { AuthLicense, AuthPlan, IssuedAuthLicense } from '@/types/api/auth'
import { parseTime } from '@/utils/ruoyi'
import AuthStatusTag from '../components/AuthStatusTag.vue'
import OneTimeLicenseDialog from '../components/OneTimeLicenseDialog.vue'

const { proxy } = getCurrentInstance() as any
const rows = ref<AuthLicense[]>([])
const plans = ref<AuthPlan[]>([])
const detail = ref<AuthLicense>()
const selected = ref<AuthLicense>()
const total = ref(0)
const loading = ref(false)
const showSearch = ref(true)
const issueOpen = ref(false)
const reasonOpen = ref(false)
const detailOpen = ref(false)
const secretOpen = ref(false)
const secret = ref('')
const submitting = ref(false)
const busy = ref('')
const action = ref('')
const query = reactive({ pageNum: 1, pageSize: 10, keyMask: '', customerPublicId: '', status: '' })
const issueForm = reactive({ customerPublicId: '', planPublicId: '', startsAt: '' })
const reasonForm = reactive({ reason: '' })
const issueRules = {
  customerPublicId: [{ required: true, message: '请输入客户公开 ID', trigger: 'blur' }],
  planPublicId: [{ required: true, message: '请选择套餐', trigger: 'change' }]
}
const actionTitle = computed(() => ({ suspend: '暂停许可证', revoke: '永久吊销许可证', replace: '替换许可证' }[action.value] || '确认操作'))
const actionWarning = computed(() => action.value === 'replace' ? '旧许可证将失效，新许可证只显示一次。' : action.value === 'revoke' ? '吊销不可恢复，所有关联租约将失效。' : '暂停后客户端将无法续租，恢复前请确认业务影响。')

async function load() {
  loading.value = true
  try { const result = await listAuthLicenses(query); rows.value = result.rows; total.value = result.total }
  finally { loading.value = false }
}
function search() { query.pageNum = 1; load() }
function reset() { proxy.resetForm('queryRef'); search() }
async function openIssue() {
  Object.assign(issueForm, { customerPublicId: '', planPublicId: '', startsAt: '' })
  const result = await listAuthPlans({ pageNum: 1, pageSize: 100, status: 'ACTIVE' })
  plans.value = result.rows
  issueOpen.value = true
}
function submitIssue() {
  proxy.$refs.issueRef.validate(async (valid: boolean) => {
    if (!valid || submitting.value) return
    submitting.value = true
    try {
      const result = await issueAuthLicense({ ...issueForm, startsAt: issueForm.startsAt || null })
      secret.value = result.data.licenseKey
      issueOpen.value = false
      secretOpen.value = true
      load()
    } finally { submitting.value = false }
  })
}
async function showDetail(row: AuthLicense) { const result = await getAuthLicense(row.publicId); detail.value = result.data; detailOpen.value = true }
async function simpleAction(row: AuthLicense, name: string) {
  const labels: Record<string, string> = { approve: '批准', renew: '续期', resume: '恢复' }
  try {
    await proxy.$modal.confirm(`是否确认${labels[name]}许可证 ${row.keyMask}？`)
  } catch { return }
  if (busy.value) return
  busy.value = row.publicId
  try {
    await actOnAuthLicense(row.publicId, name)
    proxy.$modal.msgSuccess(`${labels[name]}成功`)
    load()
  } finally { busy.value = '' }
}
function openReason(row: AuthLicense, name: string) { selected.value = row; action.value = name; reasonForm.reason = ''; reasonOpen.value = true }
function submitReasonAction() {
  proxy.$refs.reasonRef.validate(async (valid: boolean) => {
    if (!valid || !selected.value || busy.value) return
    busy.value = selected.value.publicId
    try {
      const result = await actOnAuthLicense(selected.value.publicId, action.value, reasonForm.reason)
      if (action.value === 'replace') {
        const issued = result.data as IssuedAuthLicense
        secret.value = issued.licenseKey
        secretOpen.value = true
      }
      proxy.$modal.msgSuccess('操作成功')
      reasonOpen.value = false
      load()
    } finally { busy.value = '' }
  })
}
load()
</script>
