<template>
  <div class="app-container">
    <el-form ref="queryRef" :model="query" :inline="true" v-show="showSearch">
      <el-form-item label="应用编码" prop="appCode"><el-input v-model="query.appCode" clearable placeholder="例如 desktop-client" @keyup.enter="search" /></el-form-item>
      <el-form-item label="状态" prop="status"><el-select v-model="query.status" clearable placeholder="全部状态" style="width: 160px"><el-option label="有效" value="ACTIVE" /><el-option label="停用" value="INACTIVE" /></el-select></el-form-item>
      <el-form-item><el-button type="primary" icon="Search" @click="search">搜索</el-button><el-button icon="Refresh" @click="resetQuery">重置</el-button></el-form-item>
    </el-form>
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5"><el-button type="primary" plain icon="Plus" @click="openApp()" v-hasPermi="['auth:application:add']">新增应用</el-button></el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="load" />
    </el-row>
    <el-table v-loading="loading" :data="rows">
      <el-table-column label="应用编码" prop="appCode" min-width="170" />
      <el-table-column label="名称" prop="name" min-width="180" />
      <el-table-column label="状态" width="100"><template #default="{ row }"><auth-status-tag :status="row.status" /></template></el-table-column>
      <el-table-column label="更新时间" width="180"><template #default="{ row }">{{ parseTime(row.updatedAt) }}</template></el-table-column>
      <el-table-column label="操作" width="210">
        <template #default="{ row }">
          <el-button link type="primary" @click="openCatalog(row)" v-hasPermi="['auth:application:query']">版本与权益</el-button>
          <el-button link type="primary" @click="openApp(row)" v-hasPermi="['auth:application:edit']">修改</el-button>
        </template>
      </el-table-column>
      <template #empty><el-empty description="暂无授权应用" /></template>
    </el-table>
    <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="load" />

    <el-dialog v-model="appOpen" :title="appForm.publicId ? '修改应用' : '新增应用'" width="520px">
      <el-form ref="appRef" :model="appForm" :rules="appRules" label-width="90px">
        <el-form-item label="应用编码" prop="appCode"><el-input v-model="appForm.appCode" :disabled="!!appForm.publicId" maxlength="64" /></el-form-item>
        <el-form-item label="应用名称" prop="name"><el-input v-model="appForm.name" maxlength="128" /></el-form-item>
        <el-form-item v-if="appForm.publicId" label="状态"><el-radio-group v-model="appForm.status"><el-radio value="ACTIVE">有效</el-radio><el-radio value="INACTIVE">停用</el-radio></el-radio-group></el-form-item>
      </el-form>
      <template #footer><el-button @click="appOpen = false">取消</el-button><el-button type="primary" :loading="submitting" @click="saveApp">保存</el-button></template>
    </el-dialog>

    <el-drawer v-model="catalogOpen" :title="`${current?.name || ''} · 版本与权益`" size="760px">
      <el-tabs v-model="catalogTab">
        <el-tab-pane label="版本" name="versions">
          <el-button type="primary" plain icon="Plus" class="mb8" @click="openVersion()" v-hasPermi="['auth:version:edit']">新增版本</el-button>
          <el-table :data="versions">
            <el-table-column label="版本" prop="version" />
            <el-table-column label="渠道" prop="channel" />
            <el-table-column label="最低协议" prop="minimumProtocolVersion" width="100" />
            <el-table-column label="状态" width="100"><template #default="{ row }"><auth-status-tag :status="row.status" /></template></el-table-column>
            <el-table-column label="操作" width="80"><template #default="{ row }"><el-button link type="primary" @click="openVersion(row)" v-hasPermi="['auth:version:edit']">修改</el-button></template></el-table-column>
            <template #empty><el-empty description="暂无版本" /></template>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="权益定义" name="features">
          <el-button type="primary" plain icon="Plus" class="mb8" @click="openFeature()" v-hasPermi="['auth:feature:edit']">新增权益</el-button>
          <el-table :data="features">
            <el-table-column label="权益编码" prop="featureCode" min-width="150" />
            <el-table-column label="名称" prop="name" />
            <el-table-column label="值类型" prop="valueType" width="110" />
            <el-table-column label="状态" width="100"><template #default="{ row }"><auth-status-tag :status="row.status" /></template></el-table-column>
            <el-table-column label="操作" width="80"><template #default="{ row }"><el-button link type="primary" @click="openFeature(row)" v-hasPermi="['auth:feature:edit']">修改</el-button></template></el-table-column>
            <template #empty><el-empty description="暂无权益定义" /></template>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-drawer>

    <el-dialog v-model="versionOpen" :title="versionForm.publicId ? '修改版本' : '新增版本'" width="520px">
      <el-form ref="versionRef" :model="versionForm" :rules="{ version: [{ required: true, message: '请输入版本', trigger: 'blur' }] }" label-width="100px">
        <el-form-item label="版本" prop="version"><el-input v-model="versionForm.version" :disabled="!!versionForm.publicId" /></el-form-item>
        <el-form-item label="渠道"><el-input v-model="versionForm.channel" :disabled="!!versionForm.publicId" placeholder="stable" /></el-form-item>
        <el-form-item label="最低协议"><el-input-number v-model="versionForm.minimumProtocolVersion" :min="1" /></el-form-item>
        <el-form-item label="状态"><el-select v-model="versionForm.status" style="width: 100%"><el-option label="有效" value="ACTIVE" /><el-option label="停用" value="INACTIVE" /></el-select></el-form-item>
      </el-form>
      <template #footer><el-button @click="versionOpen = false">取消</el-button><el-button type="primary" :loading="submitting" @click="saveVersion">保存</el-button></template>
    </el-dialog>

    <el-dialog v-model="featureOpen" :title="featureForm.publicId ? '修改权益' : '新增权益'" width="520px">
      <el-form ref="featureRef" :model="featureForm" :rules="{ featureCode: [{ required: true, message: '请输入权益编码', trigger: 'blur' }], name: [{ required: true, message: '请输入名称', trigger: 'blur' }] }" label-width="90px">
        <el-form-item label="权益编码" prop="featureCode"><el-input v-model="featureForm.featureCode" :disabled="!!featureForm.publicId" /></el-form-item>
        <el-form-item label="名称" prop="name"><el-input v-model="featureForm.name" /></el-form-item>
        <el-form-item label="值类型"><el-select v-model="featureForm.valueType" :disabled="!!featureForm.publicId" style="width: 100%"><el-option v-for="item in ['BOOLEAN','INTEGER','STRING','JSON']" :key="item" :label="item" :value="item" /></el-select></el-form-item>
        <el-form-item v-if="featureForm.publicId" label="状态"><el-select v-model="featureForm.status" style="width: 100%"><el-option label="有效" value="ACTIVE" /><el-option label="停用" value="INACTIVE" /></el-select></el-form-item>
      </el-form>
      <template #footer><el-button @click="featureOpen = false">取消</el-button><el-button type="primary" :loading="submitting" @click="saveFeature">保存</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="AuthCatalog">
import {
  addAuthApplication, addAuthFeature, addAuthVersion, listAuthApplications,
  listAuthFeatures, listAuthVersions, updateAuthApplication, updateAuthFeature, updateAuthVersion
} from '@/api/auth'
import type { AuthApplication, AuthFeature, AuthVersion } from '@/types/api/auth'
import { parseTime } from '@/utils/ruoyi'
import AuthStatusTag from '../components/AuthStatusTag.vue'

const { proxy } = getCurrentInstance() as any
const rows = ref<AuthApplication[]>([])
const versions = ref<AuthVersion[]>([])
const features = ref<AuthFeature[]>([])
const current = ref<AuthApplication>()
const total = ref(0)
const loading = ref(false)
const submitting = ref(false)
const showSearch = ref(true)
const appOpen = ref(false)
const catalogOpen = ref(false)
const versionOpen = ref(false)
const featureOpen = ref(false)
const catalogTab = ref('versions')
const query = reactive({ pageNum: 1, pageSize: 10, appCode: '', status: '' })
const appForm = reactive<any>({})
const versionForm = reactive<any>({})
const featureForm = reactive<any>({})
const appRules = { appCode: [{ required: true, message: '请输入应用编码', trigger: 'blur' }], name: [{ required: true, message: '请输入应用名称', trigger: 'blur' }] }

async function load() { loading.value = true; try { const result = await listAuthApplications(query); rows.value = result.rows; total.value = result.total } finally { loading.value = false } }
function search() { query.pageNum = 1; load() }
function resetQuery() { proxy.resetForm('queryRef'); search() }
function openApp(row?: AuthApplication) { Object.assign(appForm, row ? { ...row } : { publicId: '', appCode: '', name: '', status: 'ACTIVE', lockVersion: 0 }); appOpen.value = true }
function saveApp() {
  proxy.$refs.appRef.validate(async (valid: boolean) => {
    if (!valid || submitting.value) return
    submitting.value = true
    try {
      if (appForm.publicId) await updateAuthApplication(appForm.publicId, { name: appForm.name, status: appForm.status, lockVersion: appForm.lockVersion })
      else await addAuthApplication({ appCode: appForm.appCode, name: appForm.name })
      proxy.$modal.msgSuccess('保存成功'); appOpen.value = false; load()
    } finally { submitting.value = false }
  })
}
async function refreshCatalog() {
  if (!current.value) return
  const [versionResult, featureResult] = await Promise.all([listAuthVersions(current.value.publicId), listAuthFeatures(current.value.publicId)])
  versions.value = versionResult.data
  features.value = featureResult.data
}
async function openCatalog(row: AuthApplication) { current.value = row; catalogOpen.value = true; await refreshCatalog() }
function openVersion(row?: AuthVersion) { Object.assign(versionForm, row ? { ...row } : { publicId: '', version: '', channel: 'stable', status: 'ACTIVE', minimumProtocolVersion: 1, releasedAt: null }); versionOpen.value = true }
function saveVersion() {
  proxy.$refs.versionRef.validate(async (valid: boolean) => {
    if (!valid || !current.value || submitting.value) return
    submitting.value = true
    try {
      const payload = { status: versionForm.status, minimumProtocolVersion: versionForm.minimumProtocolVersion, releasedAt: versionForm.releasedAt || null }
      if (versionForm.publicId) await updateAuthVersion(current.value.publicId, versionForm.publicId, payload)
      else await addAuthVersion(current.value.publicId, { ...payload, version: versionForm.version, channel: versionForm.channel })
      versionOpen.value = false; await refreshCatalog(); proxy.$modal.msgSuccess('版本已保存')
    } finally { submitting.value = false }
  })
}
function openFeature(row?: AuthFeature) { Object.assign(featureForm, row ? { ...row } : { publicId: '', featureCode: '', name: '', valueType: 'BOOLEAN', status: 'ACTIVE' }); featureOpen.value = true }
function saveFeature() {
  proxy.$refs.featureRef.validate(async (valid: boolean) => {
    if (!valid || !current.value || submitting.value) return
    submitting.value = true
    try {
      if (featureForm.publicId) await updateAuthFeature(current.value.publicId, featureForm.publicId, { name: featureForm.name, status: featureForm.status })
      else await addAuthFeature(current.value.publicId, { featureCode: featureForm.featureCode, name: featureForm.name, valueType: featureForm.valueType })
      featureOpen.value = false; await refreshCatalog(); proxy.$modal.msgSuccess('权益已保存')
    } finally { submitting.value = false }
  })
}
load()
</script>
