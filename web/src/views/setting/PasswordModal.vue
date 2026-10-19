<script setup lang="ts">
import { NButton, NInput } from 'naive-ui'
import { useGettext } from 'vue3-gettext'

import user from '@/api/panel/user'

const { $gettext } = useGettext()
const show = defineModel<boolean>('show', { type: Boolean, required: true })
const id = defineModel<number>('id', { type: Number, required: true })
const model = ref({
  password: '',
})

const loading = ref(false)

// 弹窗按用户复用，每次打开时清空上次输入的密码
watch(show, (value) => {
  if (value) {
    model.value.password = ''
  }
})

const handleUpdate = () => {
  loading.value = true
  useRequest(() => user.updatePassword(id.value, model.value.password))
    .onSuccess(() => {
      show.value = false
      window.$message.success($gettext('Updated successfully'))
      window.$bus.emit('user:refresh')
    })
    .onComplete(() => {
      loading.value = false
    })
}
</script>

<template>
  <app-modal v-model:show="show" :title="$gettext('Change Password')">
    <n-form :model="model">
      <n-form-item path="password" :label="$gettext('Password')">
        <n-input
          v-model:value="model.password"
          type="password"
          show-password-on="click"
          @keydown.enter.prevent
          :placeholder="$gettext('Enter user password')"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block :loading="loading" :disabled="loading" @click="handleUpdate">{{
      $gettext('Submit')
    }}</n-button>
  </app-modal>
</template>

<style scoped lang="scss"></style>
