<script setup lang="ts">
import database from '@/api/panel/database'
import { NButton, NInput } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const createModel = ref({
  database: '',
  username: '',
  password: '',
  host: 'localhost'
})

const hostType = [
  { label: '本地（localhost）', value: 'localhost' },
  { label: '所有（%）', value: '%' },
  { label: '指定', value: '' }
]

const handleCreate = () => {
  useRequest(() => database.create(createModel.value)).onSuccess(() => {
    show.value = false
    window.$message.success('创建成功')
    window.$bus.emit('database:refresh')
  })
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建数据库"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="show = false"
  >
    <n-form :model="createModel">
      <n-form-item path="database" label="数据库名">
        <n-input
          v-model:value="createModel.database"
          type="text"
          @keydown.enter.prevent
          placeholder="输入数据库名称"
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
    </n-form>
    <n-button type="info" block @click="handleCreate">提交</n-button>
  </n-modal>
</template>

<style scoped lang="scss"></style>
