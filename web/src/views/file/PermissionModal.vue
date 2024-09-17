<script setup lang="ts">
import { NButton, NInput } from 'naive-ui'
import file from '@/api/panel/file'
import EventBus from '@/utils/event'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const selected = defineModel<string[]>('selected', { type: Array, required: true })
const mode = ref('755')
const owner = ref('www')
const group = ref('www')

const checkbox = ref({
  owner: ['read', 'write', 'execute'],
  group: ['read', 'execute'],
  other: ['read', 'execute']
})

const handlePermission = async () => {
  for (const path of selected.value) {
    await file
      .permission(path, `0${mode.value}`, owner.value, group.value)
      .then(() => {
        window.$message.success(`修改 ${path} 成功`)
        show.value = false
        selected.value = []
      })
      .catch(() => {
        window.$message.error(`修改 ${path} 失败`)
      })
  }
  EventBus.emit('file:refresh')
}

const calculateOctal = (permissions: string[]) => {
  let octal = 0
  if (permissions.includes('read')) octal += 4
  if (permissions.includes('write')) octal += 2
  if (permissions.includes('execute')) octal += 1
  return octal
}

const calculateMode = () => {
  const owner = calculateOctal(checkbox.value.owner)
  const group = calculateOctal(checkbox.value.group)
  const other = calculateOctal(checkbox.value.other)

  mode.value = `${owner}${group}${other}`
}

const updateCheckboxes = () => {
  const permissions = mode.value.split('').map(Number)
  checkbox.value.owner = permissions[0] & 4 ? ['read'] : []
  if (permissions[0] & 2) checkbox.value.owner.push('write')
  if (permissions[0] & 1) checkbox.value.owner.push('execute')

  checkbox.value.group = permissions[1] & 4 ? ['read'] : []
  if (permissions[1] & 2) checkbox.value.group.push('write')
  if (permissions[1] & 1) checkbox.value.group.push('execute')

  checkbox.value.other = permissions[2] & 4 ? ['read'] : []
  if (permissions[2] & 2) checkbox.value.other.push('write')
  if (permissions[2] & 1) checkbox.value.other.push('execute')
}

const title = computed(() => {
  return selected.value.length > 1 ? '批量修改权限' : `修改权限 -  ${selected.value[0]}`
})

watch(mode, updateCheckboxes, { immediate: true })
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="title"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-flex vertical>
      <n-form>
        <n-row :gutter="[0, 24]">
          <n-col :span="8">
            <n-form-item label="所有者">
              <n-checkbox-group v-model:value="checkbox.owner" @update:value="calculateMode">
                <n-checkbox value="read" label="读取" />
                <n-checkbox value="write" label="写入" />
                <n-checkbox value="execute" label="执行" />
              </n-checkbox-group>
            </n-form-item>
          </n-col>
          <n-col :span="8">
            <n-form-item label="用户组">
              <n-checkbox-group v-model:value="checkbox.group" @update:value="calculateMode">
                <n-checkbox value="read" label="读取" />
                <n-checkbox value="write" label="写入" />
                <n-checkbox value="execute" label="执行" />
              </n-checkbox-group>
            </n-form-item>
          </n-col>
          <n-col :span="8">
            <n-form-item label="其他">
              <n-checkbox-group v-model:value="checkbox.other" @update:value="calculateMode">
                <n-checkbox value="read" label="读取" />
                <n-checkbox value="write" label="写入" />
                <n-checkbox value="execute" label="执行" />
              </n-checkbox-group>
            </n-form-item>
          </n-col>
        </n-row>
        <n-form-item label="权限">
          <n-input v-model:value="mode" />
        </n-form-item>
        <n-form-item label="所有者">
          <n-input v-model:value="owner" />
        </n-form-item>
        <n-form-item label="用户组">
          <n-input v-model:value="group" />
        </n-form-item>
      </n-form>
      <n-button type="primary" @click="handlePermission"> 修改</n-button>
    </n-flex>
  </n-modal>
</template>

<style scoped lang="scss"></style>
