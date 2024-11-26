<script setup lang="ts">
import database from '@/api/panel/database'
import { NButton, NInput } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const createModel = ref({
  name: '',
  type: 'mysql',
  host: '127.0.0.1',
  port: 3306,
  username: '',
  password: '',
  remark: ''
})

const databaseType = [
  { label: 'MySQL', value: 'mysql' },
  { label: 'PostgreSQL', value: 'postgresql' }
]

watch(
  () => createModel.value.type,
  (value) => {
    if (value === 'mysql') {
      createModel.value.port = 3306
    } else if (value === 'postgresql') {
      createModel.value.port = 5432
    }
  }
)

const handleCreate = () => {
  useRequest(() => database.serverCreate(createModel.value)).onSuccess(() => {
    show.value = false
    window.$message.success('添加成功')
    window.$bus.emit('database:refresh')
  })
}
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="添加服务器"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="show = false"
  >
    <n-form :model="createModel">
      <n-form-item path="name" label="名称">
        <n-input
          v-model:value="createModel.name"
          type="text"
          @keydown.enter.prevent
          placeholder="输入数据库服务器名称"
        />
      </n-form-item>
      <n-form-item path="type" label="类型">
        <n-select
          v-model:value="createModel.type"
          @keydown.enter.prevent
          placeholder="选择数据库类型"
          :options="databaseType"
        />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="15">
          <n-form-item path="host" label="主机">
            <n-input
              v-model:value="createModel.host"
              type="text"
              @keydown.enter.prevent
              placeholder="输入数据库服务器主机"
            />
          </n-form-item>
        </n-col>
        <n-col :span="2"></n-col>
        <n-col :span="7">
          <n-form-item path="port" label="端口">
            <n-input-number
              w-full
              v-model:value="createModel.port"
              @keydown.enter.prevent
              placeholder="输入数据库服务器端口"
            />
          </n-form-item>
        </n-col>
      </n-row>
      <n-form-item path="username" label="用户名">
        <n-input
          v-model:value="createModel.username"
          type="text"
          @keydown.enter.prevent
          placeholder="输入数据库服务器用户名"
        />
      </n-form-item>
      <n-form-item path="password" label="密码">
        <n-input
          v-model:value="createModel.password"
          type="password"
          @keydown.enter.prevent
          placeholder="输入数据库服务器密码"
        />
      </n-form-item>
      <n-form-item path="remark" label="备注">
        <n-input
          v-model:value="createModel.remark"
          type="textarea"
          @keydown.enter.prevent
          placeholder="输入数据库服务器备注"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block @click="handleCreate">提交</n-button>
  </n-modal>
</template>

<style scoped lang="scss"></style>
