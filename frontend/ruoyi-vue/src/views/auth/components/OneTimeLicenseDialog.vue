<template>
  <el-dialog
    v-model="visible"
    title="许可证已生成"
    width="620px"
    :close-on-click-modal="false"
    destroy-on-close
    @closed="clearSecret"
  >
    <el-alert
      title="该完整许可证仅在此处显示一次"
      description="关闭后系统只保留掩码，无法再次查看。请立即复制到安全的交付渠道。"
      type="warning"
      :closable="false"
      show-icon
      class="mb16"
    />
    <el-input :model-value="licenseKey" readonly type="textarea" :rows="3" autocomplete="off" />
    <template #footer>
      <el-button @click="visible = false">关闭并清除</el-button>
      <el-button type="primary" icon="CopyDocument" @click="copySecret">复制许可证</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
const props = defineProps<{ modelValue: boolean; licenseKey: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: boolean]; cleared: [] }>()
const { proxy } = getCurrentInstance() as any
const visible = computed({
  get: () => props.modelValue,
  set: value => emit('update:modelValue', value)
})

async function copySecret() {
  try {
    await navigator.clipboard.writeText(props.licenseKey)
    proxy.$modal.msgSuccess('已复制，请妥善保管')
  } catch {
    proxy.$modal.msgError('浏览器禁止访问剪贴板，请手动选择复制')
  }
}

function clearSecret() {
  emit('cleared')
}
</script>
