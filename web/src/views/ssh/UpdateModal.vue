<script setup lang="ts">
import ssh from '@/api/panel/ssh'
import { NInput } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const id = defineModel<number>('id', { type: Number, required: true })
const loading = ref(false)

const model = ref({
  name: '',
  host: '127.0.0.1',
  port: 22,
  auth_method: 'password',
  user: 'root',
  password: '',
  key: '',
  remark: ''
})

const handleSubmit = async () => {
  loading.value = true
  await ssh
    .update(id.value, model.value)
    .then(() => {
      window.$message.success('更新成功')
      id.value = 0
      loading.value = false
      show.value = false
      window.$bus.emit('ssh:refresh')
    })
    .catch(() => {
      loading.value = false
    })
}

watch(show, () => {
  if (id.value > 0) {
    ssh.get(id.value).then((res) => {
      model.value.name = res.data.name
      model.value.host = res.data.host
      model.value.port = res.data.port
      model.value.auth_method = res.data.config.auth_method
      model.value.user = res.data.config.user
      model.value.password = res.data.config.password
      model.value.key = res.data.config.key
      model.value.remark = res.data.remark
    })
  }
})
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="创建主机"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form>
      <n-form-item label="名称">
        <n-input v-model:value="model.name" placeholder="127.0.0.1" />
      </n-form-item>
      <n-row :gutter="[0, 24]" pt-20>
        <n-col :span="15">
          <n-form-item label="主机">
            <n-input v-model:value="model.host" placeholder="127.0.0.1" />
          </n-form-item>
        </n-col>
        <n-col :span="2"> </n-col>
        <n-col :span="7">
          <n-form-item label="端口">
            <n-input-number v-model:value="model.port" :min="1" :max="65535" />
          </n-form-item>
        </n-col>
      </n-row>
      <n-form-item label="认证方式">
        <n-select
          v-model:value="model.auth_method"
          :options="[
            { label: '密码', value: 'password' },
            { label: '私钥', value: 'publickey' }
          ]"
        >
        </n-select>
      </n-form-item>
      <n-form-item v-if="model.auth_method == 'password'" label="用户名">
        <n-input v-model:value="model.user" placeholder="root" />
      </n-form-item>
      <n-form-item v-if="model.auth_method == 'password'" label="密码">
        <n-input v-model:value="model.password" type="password" show-password-on="click" />
      </n-form-item>
      <n-form-item v-if="model.auth_method == 'publickey'" label="私钥">
        <n-input v-model:value="model.key" type="textarea" />
      </n-form-item>
      <n-form-item label="备注">
        <n-input v-model:value="model.remark" type="textarea" />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]" pt-20>
      <n-col :span="24">
        <n-button type="info" block :loading="loading" @click="handleSubmit"> 提交 </n-button>
      </n-col>
    </n-row>
  </n-modal>
</template>

<style scoped lang="scss"></style>
