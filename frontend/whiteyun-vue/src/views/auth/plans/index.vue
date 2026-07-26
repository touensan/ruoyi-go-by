<template>
  <div class="app-container">
    <el-form ref="queryRef" :model="query" :inline="true" v-show="showSearch">
      <el-form-item label="应用" prop="applicationPublicId">
        <el-select v-model="query.applicationPublicId" clearable filterable placeholder="全部应用" style="width: 220px">
          <el-option v-for="app in apps" :key="app.publicId" :label="`${app.name} (${app.appCode})`" :value="app.publicId" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态" prop="status"><el-select v-model="query.status" clearable placeholder="全部状态" style="width: 160px"><el-option label="有效" value="ACTIVE" /><el-option label="停用" value="INACTIVE" /></el-select></el-form-item>
      <el-form-item><el-button type="primary" icon="Search" @click="search">搜索</el-button><el-button icon="Refresh" @click="resetQuery">重置</el-button></el-form-item>
    </el-form>
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5"><el-button type="primary" plain icon="Plus" @click="openPlan()" v-hasPermi="['auth:plan:add']">新增套餐</el-button></el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="load" />
    </el-row>
    <el-table v-loading="loading" :data="rows">
      <el-table-column label="应用" prop="applicationCode" width="150" />
      <el-table-column label="套餐编码" prop="planCode" min-width="150" />
      <el-table-column label="名称" prop="name" min-width="160" />
      <el-table-column label="类型" prop="licenseKind" width="110" />
      <el-table-column label="激活 / 并发" width="120"><template #default="{ row }">{{ row.maxActivations }} / {{ row.maxConcurrentSeats }}</template></el-table-column>
      <el-table-column label="允许平台" min-width="190"><template #default="{ row }">{{ row.allowedPlatforms.join('、') || '不限' }}</template></el-table-column>
      <el-table-column label="状态" width="100"><template #default="{ row }"><auth-status-tag :status="row.status" /></template></el-table-column>
      <el-table-column label="操作" width="90"><template #default="{ row }"><el-button link type="primary" @click="openPlan(row)" v-hasPermi="['auth:plan:edit']">修改</el-button></template></el-table-column>
      <template #empty><el-empty description="暂无授权套餐" /></template>
    </el-table>
    <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="load" />

    <el-dialog v-model="planOpen" :title="form.publicId ? '修改授权套餐' : '新增授权套餐'" width="760px" :close-on-click-modal="false">
      <el-form ref="planRef" :model="form" :rules="rules" label-width="125px">
        <el-row :gutter="18">
          <el-col :span="12" :xs="24"><el-form-item label="应用" prop="applicationPublicId"><el-select v-model="form.applicationPublicId" :disabled="!!form.publicId" filterable style="width: 100%"><el-option v-for="app in apps" :key="app.publicId" :label="`${app.name} (${app.appCode})`" :value="app.publicId" /></el-select></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="套餐编码" prop="planCode"><el-input v-model="form.planCode" :disabled="!!form.publicId" /></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="名称" prop="name"><el-input v-model="form.name" /></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="许可证类型"><el-select v-model="form.licenseKind" :disabled="!!form.publicId" style="width: 100%"><el-option label="订阅" value="SUBSCRIPTION" /><el-option label="永久" value="PERPETUAL" /></el-select></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="有效天数"><el-input-number v-model="form.validityDays" :min="1" :disabled="form.licenseKind === 'PERPETUAL'" style="width: 100%" /></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="最低客户端版本"><el-input v-model="form.minimumClientVersion" placeholder="可选" /></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="最大激活数"><el-input-number v-model="form.maxActivations" :min="1" style="width: 100%" /></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="最大并发席位"><el-input-number v-model="form.maxConcurrentSeats" :min="1" style="width: 100%" /></el-form-item></el-col>
          <el-col :span="24"><el-form-item label="允许平台"><el-select v-model="form.allowedPlatforms" multiple allow-create filterable default-first-option style="width: 100%" placeholder="输入平台后回车，例如 linux-amd64" /></el-form-item></el-col>
          <el-col :span="8" :xs="24"><el-form-item label="刷新间隔(秒)"><el-input-number v-model="form.refreshIntervalSeconds" :min="30" style="width: 100%" /></el-form-item></el-col>
          <el-col :span="8" :xs="24"><el-form-item label="离线租约(秒)"><el-input-number v-model="form.leaseOfflineSeconds" :min="60" style="width: 100%" /></el-form-item></el-col>
          <el-col :span="8" :xs="24"><el-form-item label="失败宽限(秒)"><el-input-number v-model="form.failureGraceSeconds" :min="0" style="width: 100%" /></el-form-item></el-col>
          <el-col :span="12" :xs="24"><el-form-item label="离线许可证(秒)"><el-input-number v-model="form.offlineLicenseSeconds" :min="0" style="width: 100%" /></el-form-item></el-col>
          <el-col v-if="form.publicId" :span="12" :xs="24"><el-form-item label="状态"><el-radio-group v-model="form.status"><el-radio value="ACTIVE">有效</el-radio><el-radio value="INACTIVE">停用</el-radio></el-radio-group></el-form-item></el-col>
          <el-col :span="24">
            <el-form-item label="权益 JSON" prop="entitlementsText">
              <el-input v-model="form.entitlementsText" type="textarea" :rows="6" placeholder='[{"featureCode":"exports","value":true}]' />
              <div class="form-tip">仅填写 featureCode 与 value；保存前会进行 JSON 数组校验。</div>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer><el-button @click="planOpen = false">取消</el-button><el-button type="primary" :loading="submitting" @click="savePlan">保存套餐</el-button></template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="AuthPlans">
import { addAuthPlan, getAuthPlan, listAuthApplications, listAuthPlans, updateAuthPlan } from '@/api/auth'
import type { AuthApplication, AuthPlan } from '@/types/api/auth'
import AuthStatusTag from '../components/AuthStatusTag.vue'

const { proxy } = getCurrentInstance() as any
const rows = ref<AuthPlan[]>([])
const apps = ref<AuthApplication[]>([])
const total = ref(0)
const loading = ref(false)
const submitting = ref(false)
const showSearch = ref(true)
const planOpen = ref(false)
const query = reactive({ pageNum: 1, pageSize: 10, applicationPublicId: '', status: '' })
const emptyForm = () => ({
  publicId: '', applicationPublicId: '', planCode: '', name: '', licenseKind: 'SUBSCRIPTION',
  validityDays: 365 as number | undefined, maxActivations: 1, maxConcurrentSeats: 1,
  allowedPlatforms: [] as string[], minimumClientVersion: '', refreshIntervalSeconds: 3600,
  leaseOfflineSeconds: 86400, failureGraceSeconds: 3600, offlineLicenseSeconds: 0 as number | undefined,
  status: 'ACTIVE', lockVersion: 0, entitlementsText: '[]'
})
const form = reactive(emptyForm())
const rules = {
  applicationPublicId: [{ required: true, message: '请选择应用', trigger: 'change' }],
  planCode: [{ required: true, message: '请输入套餐编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入套餐名称', trigger: 'blur' }],
  entitlementsText: [{ validator: (_: unknown, value: string, callback: (error?: Error) => void) => {
    try { if (!Array.isArray(JSON.parse(value))) throw new Error(); callback() } catch { callback(new Error('请输入合法的 JSON 数组')) }
  }, trigger: 'blur' }]
}

async function loadApps() { const result = await listAuthApplications({ pageNum: 1, pageSize: 100, status: 'ACTIVE' }); apps.value = result.rows }
async function load() { loading.value = true; try { const result = await listAuthPlans(query); rows.value = result.rows; total.value = result.total } finally { loading.value = false } }
function search() { query.pageNum = 1; load() }
function resetQuery() { proxy.resetForm('queryRef'); search() }
async function openPlan(row?: AuthPlan) {
  Object.assign(form, emptyForm())
  if (row) {
    const result = await getAuthPlan(row.publicId)
    const plan = result.data
    Object.assign(form, plan, {
      minimumClientVersion: plan.minimumClientVersion || '',
      offlineLicenseSeconds: plan.offlineLicenseSeconds || 0,
      entitlementsText: JSON.stringify((plan.entitlements || []).map(item => ({ featureCode: item.featureCode, value: item.value })), null, 2)
    })
  }
  planOpen.value = true
}
function savePlan() {
  proxy.$refs.planRef.validate(async (valid: boolean) => {
    if (!valid || submitting.value) return
    submitting.value = true
    try {
      const payload: any = {
        applicationPublicId: form.applicationPublicId, planCode: form.planCode, name: form.name,
        licenseKind: form.licenseKind, validityDays: form.licenseKind === 'PERPETUAL' ? null : form.validityDays,
        maxActivations: form.maxActivations, maxConcurrentSeats: form.maxConcurrentSeats,
        allowedPlatforms: form.allowedPlatforms, minimumClientVersion: form.minimumClientVersion || null,
        refreshIntervalSeconds: form.refreshIntervalSeconds, leaseOfflineSeconds: form.leaseOfflineSeconds,
        failureGraceSeconds: form.failureGraceSeconds, offlineLicenseSeconds: form.offlineLicenseSeconds || null,
        entitlements: JSON.parse(form.entitlementsText)
      }
      if (form.publicId) {
        Object.assign(payload, { status: form.status, lockVersion: form.lockVersion })
        delete payload.applicationPublicId
        delete payload.planCode
        delete payload.licenseKind
        await updateAuthPlan(form.publicId, payload)
      } else await addAuthPlan(payload)
      proxy.$modal.msgSuccess('套餐已保存'); planOpen.value = false; load()
    } finally { submitting.value = false }
  })
}
Promise.all([loadApps(), load()])
</script>
