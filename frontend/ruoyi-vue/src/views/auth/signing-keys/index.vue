<template>
  <div class="app-container">
    <el-alert
      title="私钥边界"
      description="此页面只登记受限密钥目录中的文件引用并展示公钥元数据；不会上传、读取或显示私钥内容。"
      type="warning" :closable="false" show-icon class="mb16"
    />
    <el-form ref="queryRef" :model="query" :inline="true" v-show="showSearch">
      <el-form-item label="用途" prop="purpose">
        <el-select v-model="query.purpose" clearable placeholder="全部用途" style="width: 180px">
          <el-option v-for="item in purposes" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="query.status" clearable placeholder="全部状态" style="width: 180px">
          <el-option v-for="item in ['PENDING','ACTIVE','RETIRED','REVOKED']" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" @click="search">搜索</el-button>
        <el-button icon="Refresh" @click="reset">重置</el-button>
      </el-form-item>
    </el-form>
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5"><el-button type="primary" plain icon="Plus" @click="openRegister" v-hasPermi="['auth:signing-key:add']">登记密钥</el-button></el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="load" />
    </el-row>
    <el-table v-loading="loading" :data="rows">
      <el-table-column label="KID" prop="kid" min-width="180" show-overflow-tooltip />
      <el-table-column label="用途" prop="purpose" width="110" />
      <el-table-column label="算法" prop="algorithm" width="110" />
      <el-table-column label="状态" width="100"><template #default="{ row }"><auth-status-tag :status="row.status" /></template></el-table-column>
      <el-table-column label="签名截止" width="180"><template #default="{ row }">{{ row.signUntil ? parseTime(row.signUntil) : '未设置' }}</template></el-table-column>
      <el-table-column label="验证截止" width="180"><template #default="{ row }">{{ parseTime(row.verifyUntil) }}</template></el-table-column>
      <el-table-column label="公钥" prop="publicKey" min-width="180" show-overflow-tooltip />
      <el-table-column label="操作" width="190" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.status === 'PENDING'" link type="primary" :loading="busy === row.kid" @click="confirmAction(row, 'activate')" v-hasPermi="['auth:signing-key:rotate']">启用轮换</el-button>
          <el-button v-if="row.status === 'ACTIVE'" link type="warning" :loading="busy === row.kid" @click="confirmAction(row, 'retire')" v-hasPermi="['auth:signing-key:rotate']">停止签名</el-button>
          <el-button v-if="row.status !== 'REVOKED'" link type="danger" :loading="busy === row.kid" @click="openRevoke(row)" v-hasPermi="['auth:signing-key:revoke']">紧急吊销</el-button>
        </template>
      </el-table-column>
      <template #empty><el-empty description="暂无签名密钥元数据" /></template>
    </el-table>
    <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="load" />

    <el-dialog v-model="registerOpen" title="登记签名密钥" width="600px" :close-on-click-modal="false">
      <el-alert title="请先由服务器运维将私钥文件安全放入配置的密钥目录。" type="info" :closable="false" class="mb16" />
      <el-form ref="registerRef" :model="form" :rules="rules" label-width="110px">
        <el-form-item label="KID" prop="kid"><el-input v-model="form.kid" maxlength="64" placeholder="例如 lease-2026-q3" /></el-form-item>
        <el-form-item label="用途" prop="purpose"><el-select v-model="form.purpose" style="width: 100%"><el-option v-for="item in purposes" :key="item" :label="item" :value="item" /></el-select></el-form-item>
        <el-form-item label="文件引用" prop="providerRef"><el-input v-model="form.providerRef" autocomplete="off" placeholder="仅文件名，例如 lease-2026-q3.pem" /></el-form-item>
        <el-form-item label="生效时间" prop="notBefore"><el-date-picker v-model="form.notBefore" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZ" style="width: 100%" /></el-form-item>
        <el-form-item label="签名截止"><el-date-picker v-model="form.signUntil" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZ" style="width: 100%" /></el-form-item>
        <el-form-item label="验证截止" prop="verifyUntil"><el-date-picker v-model="form.verifyUntil" type="datetime" value-format="YYYY-MM-DDTHH:mm:ssZ" style="width: 100%" /></el-form-item>
      </el-form>
      <template #footer><el-button @click="registerOpen = false">取消</el-button><el-button type="primary" :loading="submitting" @click="submitRegister">登记为待启用</el-button></template>
    </el-dialog>

    <el-dialog v-model="revokeOpen" title="紧急吊销签名密钥" width="520px" :close-on-click-modal="false">
      <el-alert title="吊销后使用该 KID 的签名将不再可信，此操作不可恢复。" type="error" :closable="false" show-icon class="mb16" />
      <el-form ref="revokeRef" :model="revokeForm" :rules="{ reason: [{ required: true, message: '请输入吊销原因', trigger: 'blur' }] }" label-width="80px">
        <el-form-item label="KID"><el-input :model-value="selected?.kid" disabled /></el-form-item>
        <el-form-item label="原因" prop="reason"><el-input v-model="revokeForm.reason" type="textarea" maxlength="128" show-word-limit /></el-form-item>
      </el-form>
      <template #footer><el-button @click="revokeOpen = false">取消</el-button><el-button type="danger" :loading="!!busy" @click="submitRevoke">确认吊销</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="AuthSigningKeys">
import { actOnAuthSigningKey, listAuthSigningKeys, registerAuthSigningKey } from '@/api/auth'
import type { AuthSigningKey } from '@/types/api/auth'
import { parseTime } from '@/utils/ruoyi'
import AuthStatusTag from '../components/AuthStatusTag.vue'

const { proxy } = getCurrentInstance() as any
const purposes = ['LEASE', 'KEYSET', 'UPDATE']
const rows = ref<AuthSigningKey[]>([])
const total = ref(0)
const loading = ref(false)
const showSearch = ref(true)
const busy = ref('')
const submitting = ref(false)
const registerOpen = ref(false)
const revokeOpen = ref(false)
const selected = ref<AuthSigningKey>()
const query = reactive({ pageNum: 1, pageSize: 10, purpose: '', status: '' })
const form = reactive({ kid: '', purpose: 'LEASE', providerRef: '', notBefore: '', signUntil: '', verifyUntil: '' })
const revokeForm = reactive({ reason: '' })
const rules = {
  kid: [{ required: true, message: '请输入 KID', trigger: 'blur' }],
  purpose: [{ required: true, message: '请选择用途', trigger: 'change' }],
  providerRef: [{ required: true, message: '请输入文件引用', trigger: 'blur' }],
  notBefore: [{ required: true, message: '请选择生效时间', trigger: 'change' }],
  verifyUntil: [{ required: true, message: '请选择验证截止时间', trigger: 'change' }]
}

async function load() {
  loading.value = true
  try { const result = await listAuthSigningKeys(query); rows.value = result.rows; total.value = result.total }
  finally { loading.value = false }
}
function search() { query.pageNum = 1; load() }
function reset() { proxy.resetForm('queryRef'); search() }
function openRegister() {
  Object.assign(form, { kid: '', purpose: 'LEASE', providerRef: '', notBefore: '', signUntil: '', verifyUntil: '' })
  registerOpen.value = true
}
function submitRegister() {
  proxy.$refs.registerRef.validate(async (valid: boolean) => {
    if (!valid || submitting.value) return
    submitting.value = true
    try {
      await registerAuthSigningKey({ ...form, signUntil: form.signUntil || null })
      proxy.$modal.msgSuccess('密钥已登记，完成复核后再启用')
      registerOpen.value = false
      load()
    } finally { submitting.value = false }
  })
}
async function confirmAction(row: AuthSigningKey, action: 'activate' | 'retire') {
  const message = action === 'activate'
    ? `启用 ${row.kid} 将同时停止同用途旧密钥签名，是否继续？`
    : `停止 ${row.kid} 签名后仍会保留验证窗口，是否继续？`
  try {
    await proxy.$modal.confirm(message)
  } catch { return }
  if (busy.value) return
  busy.value = row.kid
  try {
    await actOnAuthSigningKey(row.kid, action)
    proxy.$modal.msgSuccess('密钥状态已更新')
    load()
  } finally { busy.value = '' }
}
function openRevoke(row: AuthSigningKey) { selected.value = row; revokeForm.reason = ''; revokeOpen.value = true }
function submitRevoke() {
  proxy.$refs.revokeRef.validate(async (valid: boolean) => {
    if (!valid || !selected.value || busy.value) return
    busy.value = selected.value.kid
    try {
      await actOnAuthSigningKey(selected.value.kid, 'revoke', revokeForm.reason)
      proxy.$modal.msgSuccess('密钥已吊销')
      revokeOpen.value = false
      load()
    } finally { busy.value = '' }
  })
}
load()
</script>
