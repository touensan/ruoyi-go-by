<template>
  <div class="app-container">
    <el-form ref="queryRef" :model="query" :inline="true" v-show="showSearch">
      <el-form-item label="动作" prop="action">
        <el-input v-model="query.action" placeholder="例如 LICENSE_REVOKED" clearable @keyup.enter="search" />
      </el-form-item>
      <el-form-item label="目标 ID" prop="targetPublicId">
        <el-input v-model="query.targetPublicId" placeholder="目标公开 ID" clearable @keyup.enter="search" />
      </el-form-item>
      <el-form-item label="请求 ID" prop="requestId">
        <el-input v-model="query.requestId" placeholder="幂等请求 ID" clearable @keyup.enter="search" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" icon="Search" @click="search">搜索</el-button>
        <el-button icon="Refresh" @click="reset">重置</el-button>
      </el-form-item>
    </el-form>
    <el-row class="mb8"><right-toolbar v-model:showSearch="showSearch" @queryTable="load" /></el-row>
    <el-table v-loading="loading" :data="rows">
      <el-table-column label="时间" width="180"><template #default="{ row }">{{ parseTime(row.occurredAt) }}</template></el-table-column>
      <el-table-column label="动作" prop="action" min-width="210" show-overflow-tooltip />
      <el-table-column label="操作者" min-width="150"><template #default="{ row }">{{ row.actorType }} · {{ row.actorRef || '-' }}</template></el-table-column>
      <el-table-column label="目标" min-width="220" show-overflow-tooltip><template #default="{ row }">{{ row.targetType }} · {{ row.targetPublicId || '-' }}</template></el-table-column>
      <el-table-column label="请求 ID" prop="requestId" min-width="210" show-overflow-tooltip />
      <el-table-column label="结果" width="100"><template #default="{ row }"><auth-status-tag :status="row.result" /></template></el-table-column>
      <template #empty><el-empty description="暂无符合条件的审计事件" /></template>
    </el-table>
    <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="load" />
  </div>
</template>

<script setup lang="ts" name="AuthAudit">
import { listAuthAudit } from '@/api/auth'
import type { AuthAuditEvent } from '@/types/api/auth'
import { parseTime } from '@/utils/ruoyi'
import AuthStatusTag from '../components/AuthStatusTag.vue'

const { proxy } = getCurrentInstance() as any
const loading = ref(false)
const showSearch = ref(true)
const rows = ref<AuthAuditEvent[]>([])
const total = ref(0)
const query = reactive({ pageNum: 1, pageSize: 10, action: '', targetPublicId: '', requestId: '' })

async function load() {
  loading.value = true
  try {
    const result = await listAuthAudit(query)
    rows.value = result.rows
    total.value = result.total
  } finally { loading.value = false }
}
function search() { query.pageNum = 1; load() }
function reset() { proxy.resetForm('queryRef'); search() }
load()
</script>
