<script setup lang="ts">
import { NButton, NSpace } from 'naive-ui'
import file from '@/api/panel/file'
import EventBus from '@/utils/event'
import { checkName, lastDirectory } from '@/utils/file'
import UploadModal from '@/views/file/UploadModal.vue'
import type { Marked } from '@/views/file/types'

const path = defineModel<string>('path', { type: String, required: true })
const selected = defineModel<string[]>('selected', { type: Array, default: () => [] })
const marked = defineModel<Marked[]>('marked', { type: Array, default: () => [] })
const archive = defineModel<boolean>('archive', { type: Boolean, required: true })
const permission = defineModel<boolean>('permission', { type: Boolean, required: true })

const upload = ref(false)
const newModal = ref(false)
const newModel = ref({
  dir: false,
  path: ''
})

const showNew = (value: string) => {
  newModel.value.dir = value !== 'file'
  newModel.value.path = ''
  newModal.value = true
}

const handleNew = () => {
  if (!checkName(newModel.value.path)) {
    window.$message.error('名称不合法')
    return
  }

  const fullPath = path.value + '/' + newModel.value.path
  file.create(fullPath, newModel.value.dir).then(() => {
    newModal.value = false
    window.$message.success('创建成功')
    EventBus.emit('file:refresh')
  })
}

const handleCopy = () => {
  if (!selected.value.length) {
    window.$message.error('请选择要复制的文件/文件夹')
    return
  }
  marked.value = selected.value.map((path) => ({
    name: lastDirectory(path),
    source: path,
    type: 'copy'
  }))
  window.$message.success('标记成功，请前往目标路径粘贴')
}

const handleMove = () => {
  if (!selected.value.length) {
    window.$message.error('请选择要移动的文件/文件夹')
    return
  }
  marked.value = selected.value.map((path) => ({
    name: lastDirectory(path),
    source: path,
    type: 'move'
  }))
  window.$message.success('标记成功，请前往目标路径粘贴')
}

const handlePaste = async () => {
  if (!marked.value.length) {
    window.$message.error('请先标记需要复制或移动的文件/文件夹')
    return
  }

  for (const { name, source, type } of marked.value) {
    const target = path.value + '/' + name
    if (type === 'copy') {
      await file.copy(source, target).then(() => {
        window.$message.success(`复制 ${source} 到 ${target} 成功`)
        EventBus.emit('file:refresh')
      })
    } else {
      await file.move(source, target).then(() => {
        window.$message.success(`移动 ${source} 到 ${target} 成功`)
        EventBus.emit('file:refresh')
      })
    }
  }

  marked.value = []
}

const bulkDelete = () => {
  if (!selected.value.length) {
    window.$message.error('请选择要删除的文件/文件夹')
    return
  }

  for (const path of selected.value) {
    file.delete(path).then(() => {
      window.$message.success(`删除 ${path} 成功`)
      EventBus.emit('file:refresh')
    })
  }
}
</script>

<template>
  <n-flex>
    <n-popselect
      :options="[
        { label: '文件', value: 'file' },
        { label: '文件夹', value: 'folder' }
      ]"
      @update:value="showNew"
    >
      <n-button type="primary"> 新建 </n-button>
    </n-popselect>
    <n-button @click="upload = true"> 上传 </n-button>
    <n-button style="display: none"> 远程下载 </n-button>
    <div ml-auto>
      <n-flex>
        <n-button v-if="marked.length" secondary type="primary" @click="handlePaste">
          粘贴
        </n-button>
        <n-button-group v-if="selected.length">
          <n-button @click="handleCopy"> 复制 </n-button>
          <n-button @click="handleMove"> 移动 </n-button>
          <n-button @click="archive = true"> 压缩 </n-button>
          <n-button @click="permission = true"> 权限 </n-button>
          <n-popconfirm @positive-click="bulkDelete">
            <template #trigger>
              <n-button>删除</n-button>
            </template>
            确定要批量删除吗？
          </n-popconfirm>
        </n-button-group>
      </n-flex>
    </div>
  </n-flex>
  <n-modal
    v-model:show="newModal"
    preset="card"
    title="新建"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="newModel">
        <n-form-item label="名称">
          <n-input v-model:value="newModel.path" />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleNew">提交</n-button>
    </n-space>
  </n-modal>
  <upload-modal v-model:show="upload" v-model:path="path" />
</template>

<style scoped lang="scss"></style>
