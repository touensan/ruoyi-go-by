<template>
  <div class="app-container points-admin">
    <section class="admin-head"><div><span>积分运营</span><h1>积分管理</h1><p>查看用户积分账户，并通过可审计流水进行人工增减。</p></div><el-button type="primary" icon="EditPen" @click="openAdjust()" v-hasPermi="['platform:points:adjust']">调整积分</el-button></section>
    <section class="admin-panel">
      <div class="filters"><el-input v-model="query.username" clearable placeholder="搜索用户名" @keyup.enter="load" /><el-button type="primary" icon="Search" @click="load">查询</el-button></div>
      <el-table v-loading="loading" :data="rows">
        <el-table-column label="用户" min-width="220"><template #default="{ row }"><strong>{{ row.nickname || row.username }}</strong><small>@{{ row.username }}</small></template></el-table-column>
        <el-table-column label="可用积分" min-width="160"><template #default="{ row }"><b>{{ pointText(row.pointsMinor) }}</b></template></el-table-column>
        <el-table-column label="状态" width="110"><template #default="{ row }"><el-tag type="success">{{ row.status }}</el-tag></template></el-table-column>
        <el-table-column label="更新时间" min-width="180"><template #default="{ row }">{{ parseTime(row.updatedAt) }}</template></el-table-column>
        <el-table-column label="操作" width="100"><template #default="{ row }"><el-button link type="primary" @click="openAdjust(row)">调整</el-button></template></el-table-column>
      </el-table>
      <pagination v-show="total > 0" :total="total" v-model:page="query.pageNum" v-model:limit="query.pageSize" @pagination="load" />
    </section>
    <el-dialog v-model="adjustOpen" title="调整用户积分" width="480px">
      <el-form label-position="top"><el-form-item label="用户"><el-select v-model="form.systemUserId" filterable style="width: 100%"><el-option v-for="item in rows" :key="item.systemUserId" :label="`${item.nickname || item.username} (@${item.username})`" :value="item.systemUserId" /></el-select></el-form-item><el-form-item label="变动积分"><el-input-number v-model="form.points" :min="-1000000" :max="1000000" style="width: 100%" /></el-form-item><el-form-item label="调整原因"><el-input v-model="form.reason" type="textarea" :rows="3" /></el-form-item></el-form>
      <template #footer><el-button @click="adjustOpen = false">取消</el-button><el-button type="primary" :loading="saving" @click="submitAdjust">确认调整</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts" name="PlatformPoints">
import { adjustPoints, listPointAccounts } from '@/api/platform'
import { parseTime } from '@/utils/ruoyi'
const loading = ref(false), saving = ref(false), adjustOpen = ref(false), total = ref(0), rows = ref<any[]>([])
const query = reactive({ pageNum: 1, pageSize: 10, username: '', status: '' })
const form = reactive<any>({ systemUserId: undefined, points: 0, reason: '' })
const pointText = (value: number) => (Number(value || 0) / 100).toLocaleString('zh-CN')
async function load() { loading.value = true; try { const result = await listPointAccounts(query); rows.value = result.rows || []; total.value = result.total || 0 } finally { loading.value = false } }
function openAdjust(row?: any) { Object.assign(form, { systemUserId: row?.systemUserId, points: 0, reason: '' }); adjustOpen.value = true }
async function submitAdjust() { if (!form.systemUserId || !form.points || !form.reason.trim()) return; saving.value = true; try { await adjustPoints(form); ElMessage.success('积分已调整'); adjustOpen.value = false; await load() } finally { saving.value = false } }
load()
</script>
<style scoped lang="scss">
.admin-head, .admin-panel { border: 1px solid var(--admin-line); border-radius: 14px; background: var(--admin-surface); }.admin-head { display: flex; align-items: center; justify-content: space-between; gap: 20px; padding: 24px 26px; }.admin-head span { color: var(--el-color-primary); font-size: 13px; font-weight: 700; }.admin-head h1 { margin: 5px 0; }.admin-head p { margin: 0; color: var(--el-text-color-secondary); }.admin-panel { margin-top: 16px; padding: 20px; }.filters { display: flex; justify-content: flex-end; gap: 10px; margin-bottom: 14px; }.filters .el-input { width: 240px; }.el-table small { display: block; color: var(--el-text-color-secondary); margin-top: 3px; }
</style>
