<script setup lang="ts">
import cert from '@/api/panel/cert'
import { NButton, NInput, NSpace } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })

const props = defineProps({
  dnsProviders: Array<any>
})

const { dnsProviders } = toRefs(props)

const model = ref<any>({
  data: {
    ak: '',
    sk: ''
  },
  type: 'aliyun',
  name: ''
})

const handleCreateDNS = async () => {
  await cert.dnsCreate(model.value)
  window.$message.success('创建成功')
  show.value = false
  model.value.data.ak = ''
  model.value.data.sk = ''
  model.value.name = ''
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建 DNS"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="model">
        <n-form-item path="name" label="备注名称">
          <n-input
            v-model:value="model.name"
            type="text"
            @keydown.enter.prevent
            placeholder="输入备注名称"
          />
        </n-form-item>
        <n-form-item path="type" label="DNS">
          <n-select
            v-model:value="model.type"
            placeholder="选择 DNS"
            clearable
            :options="dnsProviders"
          />
        </n-form-item>
        <n-form-item v-if="model.type == 'aliyun'" path="ak" label="Access Key">
          <n-input
            v-model:value="model.data.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Access Key"
          />
        </n-form-item>
        <n-form-item v-if="model.type == 'aliyun'" path="sk" label="Secret Key">
          <n-input
            v-model:value="model.data.sk"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Secret Key"
          />
        </n-form-item>
        <n-form-item v-if="model.type == 'tencent'" path="ak" label="SecretId">
          <n-input
            v-model:value="model.data.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入腾讯云 SecretId"
          />
        </n-form-item>
        <n-form-item v-if="model.type == 'tencent'" path="sk" label="SecretKey">
          <n-input
            v-model:value="model.data.sk"
            type="text"
            @keydown.enter.prevent
            placeholder="输入腾讯云 SecretKey"
          />
        </n-form-item>
        <n-form-item v-if="model.type == 'huawei'" path="ak" label="AccessKeyId">
          <n-input
            v-model:value="model.data.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入华为云 AccessKeyId"
          />
        </n-form-item>
        <n-form-item v-if="model.type == 'huawei'" path="sk" label="SecretAccessKey">
          <n-input
            v-model:value="model.data.sk"
            type="text"
            @keydown.enter.prevent
            placeholder="输入华为云 SecretAccessKey"
          />
        </n-form-item>
        <n-form-item v-if="model.type == 'cloudflare'" path="ak" label="API Key">
          <n-input
            v-model:value="model.data.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 Cloudflare API Key"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleCreateDNS">提交</n-button>
    </n-space>
  </n-modal>
</template>

<style scoped lang="scss"></style>
