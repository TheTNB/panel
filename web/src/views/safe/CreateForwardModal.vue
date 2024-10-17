<script setup lang="ts">
import firewall from '@/api/panel/firewall'
import { NButton } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const loading = ref(false)

const protocols = [
  {
    label: 'TCP',
    value: 'tcp'
  },
  {
    label: 'UDP',
    value: 'udp'
  },
  {
    label: 'TCP/UDP',
    value: 'tcp/udp'
  }
]

const createModel = ref({
  protocol: 'tcp',
  port: 8080,
  target_ip: '127.0.0.1',
  target_port: 80
})

const handleCreate = async () => {
  await firewall.createForward(createModel.value).then(() => {
    window.$message.success(`创建成功`)
  })
  createModel.value = {
    protocol: 'tcp',
    port: 8080,
    target_ip: '127.0.0.1',
    target_port: 80
  }
  show.value = false
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建转发"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="show = false"
  >
    <n-form :model="createModel">
      <n-form-item path="protocols" label="传输协议">
        <n-select v-model:value="createModel.protocol" :options="protocols" />
      </n-form-item>
      <n-form-item path="address" label="目标 IP">
        <n-input v-model:value="createModel.target_ip" placeholder="127.0.0.1" />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="12">
          <n-form-item path="address" label="源端口">
            <n-input-number
              v-model:value="createModel.port"
              :min="1"
              :max="65535"
              placeholder="8080"
            />
          </n-form-item>
        </n-col>
        <n-col :span="12">
          <n-form-item path="address" label="目标端口">
            <n-input-number
              v-model:value="createModel.target_port"
              :min="1"
              :max="65535"
              placeholder="80"
            />
          </n-form-item>
        </n-col>
      </n-row>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button type="info" :loading="loading" :disabled="loading" block @click="handleCreate">
          提交
        </n-button>
      </n-col>
    </n-row>
  </n-modal>
</template>

<style scoped lang="scss"></style>
