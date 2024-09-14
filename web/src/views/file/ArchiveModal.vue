<script setup lang="ts">
import { NButton, NInput } from 'naive-ui'
import { generateRandomString, getBase } from '@/utils'
import * as api from '@/api/panel/file'
import EventBus from '@/utils/event'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const path = defineModel<string>('path', { type: String, required: true })
const selected = defineModel<string[]>('selected', { type: Array, default: () => [] })
const format = ref('.zip')
const loading = ref(false)

const generateName = () => {
  return selected.value.length > 0
    ? `${getBase(selected.value[0])}-${generateRandomString(6)}${format.value}`
    : `${path.value}/${generateRandomString(8)}${format.value}`
}

const file = ref(generateName())

const ensureExtension = (extension: string) => {
  if (!file.value.endsWith(extension)) {
    file.value = `${getBase(file.value)}${extension}`
  }
}

const handleArchive = async () => {
  ensureExtension(format.value)
  loading.value = true
  const message = window.$message.loading('正在压缩中...', {
    duration: 0
  })
  await api.default
    .archive(selected.value, file.value)
    .then(() => {
      window.$message.success('压缩成功')
      show.value = false
      selected.value = []
    })
    .catch(() => {
      window.$message.error('压缩失败')
    })
  message?.destroy()
  loading.value = false
  EventBus.emit('file:refresh')
}

onMounted(() => {
  watch(
    selected,
    () => {
      file.value = generateName()
    },
    { immediate: true }
  )
})
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="压缩"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-flex vertical>
      <n-form>
        <n-form-item label="待压缩">
          <n-dynamic-input v-model:value="selected" :min="1" />
        </n-form-item>
        <n-form-item label="压缩为">
          <n-input v-model:value="file" />
        </n-form-item>
        <n-form-item label="格式">
          <n-select
            v-model:value="format"
            :options="[
              { label: '.zip', value: '.zip' },
              { label: '.gz', value: '.gz' },
              { label: '.xz', value: '.xz' },
              { label: '.bz2', value: '.bz2' },
              { label: '.tar', value: '.tar' },
              { label: '.tar.gz', value: '.tar.gz' },
              { label: '.tar.bz2', value: '.tar.bz2' }
            ]"
            @update:value="ensureExtension"
          />
        </n-form-item>
      </n-form>
      <n-button :loading="loading" :disabled="loading" type="primary" @click="handleArchive">
        压缩
      </n-button>
    </n-flex>
  </n-modal>
</template>

<style scoped lang="scss"></style>
