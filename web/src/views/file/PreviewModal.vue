<script setup lang="ts">
import file from '@/api/panel/file'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const path = defineModel<string>('path', { type: String, required: true })

const mime = ref('')
const content = ref('')
const img = computed(() => {
  return `data:${mime.value};base64,${content.value}`
})

watch(
  () => path.value,
  () => {
    content.value = ''
    file.content(path.value).then((res) => {
      mime.value = res.data.mime
      content.value = res.data.content
    })
  }
)
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="'预览 - ' + path"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-image width="100%" :src="img" preview-disabled :show-toolbar="false">
      <template #placeholder>
        <n-spin />
      </template>
    </n-image>
  </n-modal>
</template>

<style scoped lang="scss"></style>
