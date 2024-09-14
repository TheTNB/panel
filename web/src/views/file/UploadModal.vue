<script setup lang="ts">
import type { UploadCustomRequestOptions } from 'naive-ui'
import * as api from '@/api/panel/file'
import EventBus from '@/utils/event'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const path = defineModel<string>('path', { type: String, required: true })
const upload = ref<any>(null)

const uploadRequest = ({ file, onFinish, onError, onProgress }: UploadCustomRequestOptions) => {
  const formData = new FormData()
  formData.append('file', file.file as File)
  api.default
    .upload(`${path.value}/${file.name}`, formData, onProgress)
    .then(() => {
      window.$message.success(`上传 ${file.name} 成功`)
      EventBus.emit('file:refresh')
      onFinish()
    })
    .catch(() => {
      onError()
    })
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="上传"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-upload
      ref="upload"
      directory-dnd
      multiple
      action="/api/panel/file/upload"
      :custom-request="uploadRequest"
    >
      <n-upload-dragger>
        <div style="margin-bottom: 12px">
          <the-icon :size="48" class="mr-5" icon="bi:arrow-up-square" />
        </div>
        <NText style="font-size: 16px"> 点击或者拖动文件到该区域来上传</NText>
        <NP depth="3" style="margin: 8px 0 0 0"> 不支持断点续传，大文件建议使用 FTP 上传 </NP>
      </n-upload-dragger>
    </n-upload>
  </n-modal>
</template>

<style scoped lang="scss"></style>
