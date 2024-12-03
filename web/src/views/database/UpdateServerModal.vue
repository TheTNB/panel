<script setup lang="ts">
import database from '@/api/panel/database'
import { NButton, NInput } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const id = defineModel<number>('id', { type: Number, required: true })
const updateModel = ref({
  name: '',
  host: '127.0.0.1',
  port: 3306,
  username: '',
  password: '',
  remark: ''
})

const handleUpdate = () => {
  useRequest(() => database.serverUpdate(id.value, updateModel.value)).onSuccess(() => {
    show.value = false
    window.$message.success('修改成功')
    window.$bus.emit('database-user:refresh')
  })
}

watch(
  () => show.value,
  (value) => {
    if (value && id.value) {
      database.serverGet(id.value).then((data: any) => {
        updateModel.value.name = data.name
        updateModel.value.host = data.host
        updateModel.value.port = data.port
        updateModel.value.username = data.username
        updateModel.value.password = data.password
        updateModel.value.remark = data.remark
      })
    }
  }
)
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="修改服务器"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="show = false"
  >
    <n-form :model="updateModel">
      <n-form-item path="name" label="名称">
        <n-input
          v-model:value="updateModel.name"
          type="text"
          @keydown.enter.prevent
          placeholder="输入数据库服务器名称"
        />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="15">
          <n-form-item path="host" label="主机">
            <n-input
              v-model:value="updateModel.host"
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
              v-model:value="updateModel.port"
              @keydown.enter.prevent
              placeholder="输入数据库服务器端口"
            />
          </n-form-item>
        </n-col>
      </n-row>
      <n-form-item path="username" label="用户名">
        <n-input
          v-model:value="updateModel.username"
          type="text"
          @keydown.enter.prevent
          placeholder="输入数据库服务器用户名"
        />
      </n-form-item>
      <n-form-item path="password" label="密码">
        <n-input
          v-model:value="updateModel.password"
          type="password"
          show-password-on="click"
          @keydown.enter.prevent
          placeholder="输入数据库服务器密码"
        />
      </n-form-item>
      <n-form-item path="remark" label="备注">
        <n-input
          v-model:value="updateModel.remark"
          type="textarea"
          @keydown.enter.prevent
          placeholder="输入数据库服务器备注"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block @click="handleUpdate">提交</n-button>
  </n-modal>
</template>

<style scoped lang="scss"></style>
