<template>
  <div class="app-container">
    <el-alert
      title="运行态页面仅展示公开标识和设备摘要；凭据哈希、租约载荷及签名不会返回浏览器。"
      type="info" :closable="false" show-icon class="mb16"
    />
    <el-form ref="queryRef" :model="query" :inline="true" v-show="showSearch">
      <el-form-item label="搜索" prop="search">
        <el-input v-model="query.search" placeholder="公开 ID、掩码或设备摘要" clearable @keyup.enter="search" />
      </el-form-item>
      <el-form-item label="许可证 ID" prop="licensePublicId">
        <el-input v-model="query.licensePublicId" placeholder="可选" clearable @keyup.enter="search" />
      </el-form-item>
      <el-form-item :label="tab === 'anomalies' ? '级别' : '状态'" prop="status">
        <el-input v-model="query.status" :placeholder="tab === 'anomalies' ? '例如 HIGH' : '例如 ACTIVE'" clearable @keyup.enter="search" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" @click="search">搜索</el-button>
        <el-button icon="Refresh" @click="reset">重置</el-button>
      </el-form-item>
    </el-form>
    <el-row class="mb8"><right-toolbar v-model:showSearch="showSearch" @queryTable="load" /></el-row>
    <el-tabs v-model="tab" @tab-change="changeTab">
      <el-tab-pane label="安装" name="installations" />
      <el-tab-pane label="激活" name="activations" />
      <el-tab-pane label="租约" name="leases" />
      <el-tab-pane label="停用记录" name="deactivations" />
      <el-tab-pane label="异常" name="anomalies" />
    </el-tabs>

    <el-table v-loading="loading" :data="rows">
      <template v-if="tab === 'installations'">
        <el-table-column label="安装 ID" prop="publicId" min-width="220" show-overflow-tooltip />
        <el-table-column label="应用" prop="applicationCode" width="150" />
        <el-table-column label="设备摘要" prop="bindingDisplay" min-width="180" show-overflow-tooltip />
        <el-table-column label="平台 / 版本" min-width="180"><template #default="{ row }">{{ row.platform }} / {{ row.clientVersion }}</template></el-table-column>
        <el-table-column label="最后在线" width="180"><template #default="{ row }">{{ parseTime(row.lastSeenAt) }}</template></el-table-column>
      </template>
      <template v-else-if="tab === 'activations'">
        <el-table-column label="激活 ID" prop="publicId" min-width="210" show-overflow-tooltip />
        <el-table-column label="许可证" prop="licenseKeyMask" width="170" />
        <el-table-column label="设备摘要" prop="bindingDisplay" min-width="180" show-overflow-tooltip />
        <el-table-column label="激活时间" width="180"><template #default="{ row }">{{ parseTime(row.activatedAt) }}</template></el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'ACTIVE'" link type="danger"
              :loading="busy === row.publicId" @click="openDeactivate(row)"
              v-hasPermi="['auth:activation:deactivate']"
            >停用</el-button>
          </template>
        </el-table-column>
      </template>
      <template v-else-if="tab === 'leases'">
        <el-table-column label="租约 ID" prop="publicId" min-width="210" show-overflow-tooltip />
        <el-table-column label="许可证" prop="licenseKeyMask" width="170" />
        <el-table-column label="序列" prop="serial" width="100" />
        <el-table-column label="签名 KID" prop="signingKeyKid" min-width="150" show-overflow-tooltip />
        <el-table-column label="离线截止" width="180"><template #default="{ row }">{{ parseTime(row.offlineUntil) }}</template></el-table-column>
      </template>
      <template v-else-if="tab === 'deactivations'">
        <el-table-column label="停用 ID" prop="publicId" min-width="210" show-overflow-tooltip />
        <el-table-column label="许可证" prop="licenseKeyMask" width="170" />
        <el-table-column label="激活 ID" prop="activationPublicId" min-width="210" show-overflow-tooltip />
        <el-table-column label="原因" prop="reason" min-width="180" show-overflow-tooltip />
        <el-table-column label="操作者" width="150"><template #default="{ row }">{{ row.actorRef || row.actorType }}</template></el-table-column>
        <el-table-column label="停用时间" width="180"><template #default="{ row }">{{ parseTime(row.deactivatedAt) }}</template></el-table-column>
      </template>
      <template v-else>
        <el-table-column label="级别" width="90"><template #default="{ row }"><el-tag :type="row.severity === 'HIGH' ? 'danger' : 'warning'">{{ row.severity === 'HIGH' ? '高' : '中' }}</el-tag></template></el-table-column>
        <el-table-column label="异常代码" prop="anomalyCode" min-width="260" show-overflow-tooltip />
        <el-table-column label="对象" min-width="230" show-overflow-tooltip><template #default="{ row }">{{ row.subjectType }} · {{ row.subjectPublicId }}</template></el-table-column>
        <el-table-column label="说明" prop="summary" min-width="240" />
        <el-table-column label="检测时间" width="180"><template #default="{ row }">{{ parseTime(row.detectedAt) }}</template></el-table-column>
      </template>
      <el-table-column v-if="tab !== 'anomalies'" label="状态" prop="status" width="100">
        <template #default="{ row }"><auth-status-tag :status="row.status || (tab === 'deactivations' ? 'DEACTIVATED' : '')" /></template>
      </el-table-column>
      <template #empty><el-empty description="暂无运行态数据" /></template>
    </el-table>
    <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="load" />

    <el-dialog v-model="deactivateOpen" title="确认停用激活" width="520px" :close-on-click-modal="false">
      <el-alert title="停用后当前租约会立即撤销，客户端需重新激活。" type="warning" :closable="false" show-icon class="mb16" />
      <el-form ref="deactivateRef" :model="deactivateForm" :rules="{ reason: [{ required: true, message: '请输入停用原因', trigger: 'blur' }] }" label-width="80px">
        <el-form-item label="激活 ID"><el-input :model-value="selected?.publicId" disabled /></el-form-item>
        <el-form-item label="停用原因" prop="reason"><el-input v-model="deactivateForm.reason" type="textarea" maxlength="128" show-word-limit /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="deactivateOpen = false">取消</el-button>
        <el-button type="danger" :loading="!!busy" @click="confirmDeactivate">确认停用</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="AuthOperations">
import {
  deactivateAuthActivation, listAuthActivations, listAuthAnomalies,
  listAuthDeactivations, listAuthInstallations, listAuthLeases
} from '@/api/auth'
import type { AuthActivation, AuthRuntimeQuery } from '@/types/api/auth'
import { parseTime } from '@/utils/ruoyi'
import AuthStatusTag from '../components/AuthStatusTag.vue'

const { proxy } = getCurrentInstance() as any
const tab = ref('installations')
const rows = ref<any[]>([])
const total = ref(0)
const loading = ref(false)
const showSearch = ref(true)
const busy = ref('')
const deactivateOpen = ref(false)
const selected = ref<AuthActivation>()
const deactivateForm = reactive({ reason: '' })
const query = reactive<AuthRuntimeQuery>({ pageNum: 1, pageSize: 10, search: '', licensePublicId: '', status: '' })
const loaders: Record<string, (query: AuthRuntimeQuery) => Promise<any>> = {
  installations: listAuthInstallations, activations: listAuthActivations,
  leases: listAuthLeases, deactivations: listAuthDeactivations, anomalies: listAuthAnomalies
}

async function load() {
  loading.value = true
  try {
    const result = await loaders[tab.value](query)
    rows.value = result.rows
    total.value = result.total
  } finally { loading.value = false }
}
function search() { query.pageNum = 1; load() }
function reset() { proxy.resetForm('queryRef'); search() }
function changeTab() { query.pageNum = 1; load() }
function openDeactivate(row: AuthActivation) {
  selected.value = row
  deactivateForm.reason = ''
  deactivateOpen.value = true
}
function confirmDeactivate() {
  proxy.$refs.deactivateRef.validate(async (valid: boolean) => {
    if (!valid || !selected.value || busy.value) return
    busy.value = selected.value.publicId
    try {
      await deactivateAuthActivation(selected.value.publicId, deactivateForm.reason)
      proxy.$modal.msgSuccess('激活已停用，相关租约已撤销')
      deactivateOpen.value = false
      load()
    } finally { busy.value = '' }
  })
}
load()
</script>
