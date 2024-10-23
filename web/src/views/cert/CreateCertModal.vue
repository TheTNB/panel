<script setup lang="ts">
import cert from '@/api/panel/cert'
import { NButton, NSpace } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })

const props = defineProps({
  algorithms: Array<any>,
  websites: Array<any>,
  accounts: Array<any>,
  dns: Array<any>
})

const { algorithms, websites, accounts, dns } = toRefs(props)

const model = ref<any>({
  domains: [],
  dns_id: 0,
  type: 'P256',
  account_id: null,
  website_id: 0,
  auto_renew: true
})

const handleCreateCert = async () => {
  await cert.certCreate(model.value)
  show.value = false
  window.$message.success('创建成功')
  model.value.domains = []
  model.value.dns_id = 0
  model.value.type = 'P256'
  model.value.account_id = 0
  model.value.website_id = 0
  model.value.auto_renew = true
  window.$bus.emit('cert:refresh-cert')
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建证书"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-alert type="info">
        可以通过选择网站 / DNS 中的任意一项来自动签发和部署证书，也可以手动输入域名并设置 DNS
        解析来签发证书
      </n-alert>
      <n-form :model="model">
        <n-form-item label="域名">
          <n-dynamic-input
            v-model:value="model.domains"
            placeholder="example.com"
            :min="1"
            show-sort-button
          />
        </n-form-item>
        <n-form-item path="type" label="密钥类型">
          <n-select
            v-model:value="model.type"
            placeholder="选择密钥类型"
            clearable
            :options="algorithms"
          />
        </n-form-item>
        <n-form-item path="website_id" label="网站">
          <n-select
            v-model:value="model.website_id"
            placeholder="选择用于部署证书的网站"
            clearable
            :options="websites"
          />
        </n-form-item>
        <n-form-item path="account_id" label="账号">
          <n-select
            v-model:value="model.account_id"
            placeholder="选择用于签发证书的账号"
            clearable
            :options="accounts"
          />
        </n-form-item>
        <n-form-item path="account_id" label="DNS">
          <n-select
            v-model:value="model.dns_id"
            placeholder="选择用于签发证书的DNS"
            clearable
            :options="dns"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleCreateCert">提交</n-button>
    </n-space>
  </n-modal>
</template>

<style scoped lang="scss"></style>
