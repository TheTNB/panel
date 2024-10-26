<script setup lang="ts">
import cert from '@/api/panel/cert'
import { NButton, NSpace } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })

const model = ref<any>({
  cert: '',
  key: ''
})

const handleSubmit = async () => {
  await cert.certUpload(model.value)
  show.value = false
  window.$message.success('创建成功')
  model.value.cert = ''
  model.value.key = ''
  window.$bus.emit('cert:refresh-cert')
  window.$bus.emit('cert:refresh-async')
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="上传证书"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="model">
        <n-form-item label="证书">
          <n-input
            v-model:value="model.cert"
            type="textarea"
            placeholder="输入 PEM 证书文件的内容"
          />
        </n-form-item>
        <n-form-item label="私钥">
          <n-input
            v-model:value="model.key"
            type="textarea"
            placeholder="输入 KEY 私钥文件的内容"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleSubmit">提交</n-button>
    </n-space>
  </n-modal>
</template>

<style scoped lang="scss"></style>
