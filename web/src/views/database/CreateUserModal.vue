<script setup lang="ts">
import database from '@/api/panel/database'
import { NButton, NInput } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const createModel = ref({
  server_id: null,
  username: '',
  password: '',
  host: 'localhost',
  privileges: [],
  remark: ''
})

const servers = ref<{ label: string; value: string }[]>([])

const hostType = [
  { label: '本地（localhost）', value: 'localhost' },
  { label: '所有（%）', value: '%' },
  { label: '指定', value: '' }
]

const handleCreate = () => {
  useRequest(() => database.userCreate(createModel.value)).onSuccess(() => {
    show.value = false
    window.$message.success('创建成功')
    window.$bus.emit('database-user:refresh')
  })
}

onMounted(() => {
  database.serverList(1, 10000).then((data: any) => {
    for (const server of data.items) {
      servers.value.push({
        label: server.name,
        value: server.id
      })
    }
  })
})
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建用户"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="show = false"
  >
    <n-form :model="createModel">
      <n-form-item path="server_id" label="服务器">
        <n-select
          v-model:value="createModel.server_id"
          @keydown.enter.prevent
          placeholder="选择服务器"
          :options="servers"
        />
      </n-form-item>
      <n-form-item path="username" label="用户名">
        <n-input
          v-model:value="createModel.username"
          type="text"
          @keydown.enter.prevent
          placeholder="输入用户名"
        />
      </n-form-item>
      <n-form-item path="password" label="密码">
        <n-input
          v-model:value="createModel.password"
          type="password"
          @keydown.enter.prevent
          placeholder="输入密码"
        />
      </n-form-item>
      <n-form-item path="host-select" label="主机">
        <n-select
          v-model:value="createModel.host"
          @keydown.enter.prevent
          placeholder="选择主机"
          :options="hostType"
        />
      </n-form-item>
      <n-form-item v-if="createModel.host === ''" path="host" label="指定主机">
        <n-input
          v-model:value="createModel.host"
          type="text"
          @keydown.enter.prevent
          placeholder="输入受支持的主机地址"
        />
      </n-form-item>
      <n-form-item path="privileges" label="授权">
        <n-dynamic-input v-model:value="createModel.privileges" placeholder="输入数据库名" />
      </n-form-item>
      <n-form-item path="remark" label="备注">
        <n-input
          v-model:value="createModel.remark"
          type="textarea"
          @keydown.enter.prevent
          placeholder="输入数据库用户备注"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block @click="handleCreate">提交</n-button>
  </n-modal>
</template>

<style scoped lang="scss"></style>
