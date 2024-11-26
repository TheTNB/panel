<script setup lang="ts">
import database from '@/api/panel/database'
import { NButton, NInput } from 'naive-ui'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const id = defineModel<number>('id', { type: Number, required: true })
const updateModel = ref({
  password: '',
  privileges: [],
  remark: ''
})

const handleUpdate = () => {
  useRequest(() => database.userUpdate(id.value, updateModel.value)).onSuccess(() => {
    show.value = false
    window.$message.success('修改成功')
    window.$bus.emit('database-user:refresh')
  })
}

watch(
  () => show.value,
  (value) => {
    if (value && id.value) {
      database.userGet(id.value).then((data: any) => {
        updateModel.value.password = data.password
        updateModel.value.privileges = data.privileges
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
    title="修改用户"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="show = false"
  >
    <n-form :model="updateModel">
      <n-form-item path="password" label="密码">
        <n-input
          v-model:value="updateModel.password"
          type="password"
          @keydown.enter.prevent
          placeholder="输入密码"
        />
      </n-form-item>
      <n-form-item path="privileges" label="授权">
        <n-dynamic-input v-model:value="updateModel.privileges" placeholder="输入数据库名" />
      </n-form-item>
      <n-form-item path="remark" label="备注">
        <n-input
          v-model:value="updateModel.remark"
          type="textarea"
          @keydown.enter.prevent
          placeholder="输入数据库用户备注"
        />
      </n-form-item>
    </n-form>
    <n-button type="info" block @click="handleUpdate">提交</n-button>
  </n-modal>
</template>

<style scoped lang="scss"></style>
