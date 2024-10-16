<script setup lang="ts">
import safe from '@/api/panel/safe'
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

const families = [
  {
    label: 'IPv4',
    value: 'ipv4'
  },
  {
    label: 'IPv6',
    value: 'ipv6'
  }
]

const strategies = [
  {
    label: '接受',
    value: 'accept'
  },
  {
    label: '丢弃',
    value: 'drop'
  },
  {
    label: '拒绝',
    value: 'reject'
  }
]

const directions = [
  {
    label: '传入',
    value: 'in'
  },
  {
    label: '传出',
    value: 'out'
  }
]

const createModel = ref({
  family: 'ipv4',
  protocol: 'tcp',
  port_start: 80,
  port_end: 80,
  address: '',
  strategy: 'accept',
  direction: 'in'
})

const handleCreate = async () => {
  await safe.createFirewallRule(createModel.value).then(() => {
    window.$message.success('创建成功')
    createModel.value = {
      family: 'ipv4',
      protocol: 'tcp',
      port_start: 80,
      port_end: 80,
      address: '',
      strategy: 'accept',
      direction: 'in'
    }
  })
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建规则"
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
      <n-form-item path="family" label="网络协议">
        <n-select v-model:value="createModel.family" :options="families" />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="12">
          <n-form-item path="port_start" label="起始端口">
            <n-input-number
              v-model:value="createModel.port_start"
              :min="1"
              :max="65535"
              placeholder="80"
            />
          </n-form-item>
        </n-col>
        <n-col :span="12">
          <n-form-item path="port_end" label="结束端口">
            <n-input-number
              v-model:value="createModel.port_end"
              :min="1"
              :max="65535"
              placeholder="80"
            />
          </n-form-item>
        </n-col>
      </n-row>
      <n-form-item path="address" label="目标">
        <n-input
          v-model:value="createModel.address"
          placeholder="可选输入 IP 或 IP 段：127.0.0.1 或 172.16.0.0/24（多个以英文逗号隔开）"
        />
      </n-form-item>
      <n-form-item path="strategy" label="策略">
        <n-select v-model:value="createModel.strategy" :options="strategies" />
      </n-form-item>
      <n-form-item path="strategy" label="方向">
        <n-select v-model:value="createModel.direction" :options="directions" />
      </n-form-item>
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
