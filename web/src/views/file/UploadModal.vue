<script setup lang="ts">
import type { UploadCustomRequestOptions } from 'naive-ui'

import api from '@/api/panel/file'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const path = defineModel<string>('path', { type: String, required: true })
const upload = ref<any>(null)

const uploadRequest = ({ file, onFinish, onError, onProgress }: UploadCustomRequestOptions) => {
  const formData = new FormData()
  formData.append('file', file.file as File)
  api
    .upload(`${path.value}/${file.name}`, formData, onProgress)
    .then(() => {
      window.$message.success(`上传 ${file.name} 成功`)
      window.$bus.emit('file:refresh')
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
    <n-flex vertical>
      <n-alert type="info">若上传报网络错误，请开启面板 HTTPS 后重试</n-alert>
      <n-upload
        ref="upload"
        multiple
        directory-dnd
        action="/api/panel/file/upload"
        :custom-request="uploadRequest"
      >
        <n-upload-dragger>
          <div style="margin-bottom: 12px">
            <the-icon :size="48" icon="bi:arrow-up-square" />
          </div>
          <NText text-18> 点击或者拖动文件到该区域来上传</NText>
          <NP depth="3" m-10> 大文件建议使用 SFTP 上传 </NP>
        </n-upload-dragger>
      </n-upload>
    </n-flex>
  </n-modal>
</template>

<style scoped lang="scss"></style>
