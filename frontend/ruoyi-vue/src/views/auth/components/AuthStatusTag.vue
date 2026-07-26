<template>
  <el-tag :type="tagType" effect="light">{{ statusLabel }}</el-tag>
</template>

<script setup lang="ts">
const props = defineProps<{ status?: string }>()

const tagType = computed(() => {
  if (['ACTIVE', 'ISSUED', 'SUCCESS'].includes(props.status || '')) return 'success'
  if (['PENDING', 'SUSPENDED', 'SUPERSEDED'].includes(props.status || '')) return 'warning'
  if (['REVOKED', 'DEACTIVATED', 'DISABLED', 'FAILED'].includes(props.status || '')) return 'danger'
  return 'info'
})

const labels: Record<string, string> = {
  ACTIVE: '有效', INACTIVE: '未启用', PENDING: '待处理', SUSPENDED: '已暂停',
  REVOKED: '已吊销', DEACTIVATED: '已停用', RETIRED: '已退役',
  ISSUED: '已签发', SUPERSEDED: '已替代', SUCCESS: '成功', FAILED: '失败'
}
const statusLabel = computed(() => labels[props.status || ''] || props.status || '-')
</script>
