<script setup lang="ts">
import { NButton, NInput } from 'naive-ui'
import { useGettext } from 'vue3-gettext'

import database from '@/api/panel/database'

const { $gettext } = useGettext()
const show = defineModel<boolean>('show', { type: Boolean, required: true })
const id = defineModel<number>('id', { type: Number, required: true })
const serverType = ref('')
const updateModel = ref({
  name: '',
  host: '127.0.0.1',
  port: 3306,
  username: '',
  password: '',
  remark: '',
})

const loading = ref(false)

const handleUpdate = () => {
  loading.value = true
  useRequest(() => database.serverUpdate(id.value, updateModel.value))
    .onSuccess(() => {
      show.value = false
      window.$message.success($gettext('Modified successfully'))
      window.$bus.emit('database-user:refresh')
      window.$bus.emit('database-server:refresh')
    })
    .onComplete(() => {
      loading.value = false
    })
}

watch(
  () => show.value,
  (value) => {
    if (value && id.value) {
      useRequest(database.serverGet(id.value)).onSuccess(({ data }: { data: any }) => {
        serverType.value = data.type
        updateModel.value.name = data.name
        updateModel.value.host = data.host
        updateModel.value.port = data.port
        updateModel.value.username = data.username
        updateModel.value.password = data.password
        updateModel.value.remark = data.remark
      })
    }
  },
)
</script>

<template>
  <app-modal v-model:show="show" :title="$gettext('Modify Server')">
    <n-form :model="updateModel">
      <n-form-item path="name" :label="$gettext('Name')">
        <n-input
          v-model:value="updateModel.name"
          type="text"
          @keydown.enter.prevent
          :placeholder="$gettext('Enter database server name')"
        />
      </n-form-item>
      <n-grid cols="1 m:24" responsive="screen">
        <n-form-item-gi :span="15" path="host" :label="$gettext('Host')">
          <n-input
            v-model:value="updateModel.host"
            type="text"
            @keydown.enter.prevent
            :placeholder="$gettext('Enter database server host')"
          />
        </n-form-item-gi>
        <n-gi :span="2" />
        <n-form-item-gi :span="7" path="port" :label="$gettext('Port')">
          <n-input-number
            w-full
            v-model:value="updateModel.port"
            @keydown.enter.prevent
            :placeholder="$gettext('Enter database server port')"
          />
        </n-form-item-gi>
      </n-grid>
      <n-form-item v-if="serverType !== 'redis'" path="username" :label="$gettext('Username')">
        <n-input
          v-model:value="updateModel.username"
          type="text"
          @keydown.enter.prevent
          :placeholder="$gettext('Enter database server username')"
        />
      </n-form-item>
      <n-form-item path="password" :label="$gettext('Password')">
        <n-input
          v-model:value="updateModel.password"
          type="password"
          show-password-on="click"
          @keydown.enter.prevent
          :placeholder="$gettext('Enter database server password')"
        />
      </n-form-item>
      <n-form-item path="remark" :label="$gettext('Comment')">
        <n-input
          v-model:value="updateModel.remark"
          type="textarea"
          @keydown.enter.prevent
          :placeholder="$gettext('Enter database server comment')"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block :loading="loading" :disabled="loading" @click="handleUpdate">{{
      $gettext('Submit')
    }}</n-button>
  </app-modal>
</template>

<style scoped lang="scss"></style>
