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
  address: [],
  strategy: 'accept',
  direction: 'in'
})

const handleCreate = async () => {
  for (const address of createModel.value.address) {
    await firewall
      .createIpRule({
        ...createModel.value,
        address
      })
      .then(() => {
        window.$message.success(`${address} 创建成功`)
      })
  }
  createModel.value = {
    family: 'ipv4',
    protocol: 'tcp',
    address: [],
    strategy: 'accept',
    direction: 'in'
  }
  show.value = false
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
      <n-form-item path="address" label="IP 地址">
        <n-dynamic-input
          v-model:value="createModel.address"
          show-sort-button
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
    <n-button type="info" block :loading="loading" :disabled="loading" @click="handleCreate">
      提交
    </n-button>
  </n-modal>
</template>

<style scoped lang="scss"></style>
