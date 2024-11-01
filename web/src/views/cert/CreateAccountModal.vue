<script setup lang="ts">
import cert from '@/api/panel/cert'
import type { MessageReactive } from 'naive-ui'
import { NButton, NInput, NSpace } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })

const props = defineProps({
  caProviders: {
    type: Array<any>,
    required: true
  },
  algorithms: {
    type: Array<any>,
    required: true
  }
})

const { caProviders, algorithms } = toRefs(props)

let messageReactive: MessageReactive | null = null

const model = ref<any>({
  hmac_encoded: '',
  email: '',
  kid: '',
  key_type: 'P256',
  ca: 'googlecn'
})

const showEAB = computed(() => {
  return model.value.ca === 'google' || model.value.ca === 'sslcom'
})

const handleCreateAccount = async () => {
  messageReactive = window.$message.loading('正在向 CA 注册账号，请耐心等待', {
    duration: 0
  })
  cert
    .accountCreate(model.value)
    .then(() => {
      show.value = false
      window.$message.success('创建成功')
      model.value.email = ''
      model.value.hmac_encoded = ''
      model.value.kid = ''
    })
    .finally(() => {
      messageReactive?.destroy()
      window.$bus.emit('cert:refresh-account')
      window.$bus.emit('cert:refresh-async')
    })
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建账号"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-alert type="info"> Google 和 SSL.com 需要先去官网获得 KID 和 HMAC 并填入 </n-alert>
      <n-alert type="warning">
        境内无法使用 Google，其他 CA 视网络情况而定，建议使用 GoogleCN 或 Let's Encrypt
      </n-alert>
      <n-form :model="model">
        <n-form-item path="ca" label="CA">
          <n-select
            v-model:value="model.ca"
            placeholder="选择 CA"
            clearable
            :options="caProviders"
          />
        </n-form-item>
        <n-form-item path="key_type" label="密钥类型">
          <n-select
            v-model:value="model.key_type"
            placeholder="选择密钥类型"
            clearable
            :options="algorithms"
          />
        </n-form-item>
        <n-form-item path="email" label="邮箱">
          <n-input
            v-model:value="model.email"
            type="text"
            @keydown.enter.prevent
            placeholder="输入邮箱地址"
          />
        </n-form-item>
        <n-form-item v-if="showEAB" path="kid" label="KID">
          <n-input
            v-model:value="model.kid"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 KID"
          />
        </n-form-item>
        <n-form-item v-if="showEAB" path="hmac_encoded" label="HMAC">
          <n-input
            v-model:value="model.hmac_encoded"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 HMAC"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleCreateAccount">提交</n-button>
    </n-space>
  </n-modal>
</template>

<style scoped lang="scss"></style>
