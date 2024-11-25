<script setup lang="ts">
import file from '@/api/panel/file'
import { checkName, lastDirectory } from '@/utils/file'
import UploadModal from '@/views/file/UploadModal.vue'
import type { Marked } from '@/views/file/types'

const path = defineModel<string>('path', { type: String, required: true })
const selected = defineModel<string[]>('selected', { type: Array, default: () => [] })
const marked = defineModel<Marked[]>('marked', { type: Array, default: () => [] })
const markedType = defineModel<string>('markedType', { type: String, required: true })
const compress = defineModel<boolean>('compress', { type: Boolean, required: true })
const permission = defineModel<boolean>('permission', { type: Boolean, required: true })

const upload = ref(false)
const create = ref(false)
const createModel = ref({
  dir: false,
  path: ''
})
const download = ref(false)
const downloadModel = ref({
  path: '',
  url: ''
})

const showCreate = (value: string) => {
  createModel.value.dir = value !== 'file'
  createModel.value.path = ''
  create.value = true
}

const handleCreate = () => {
  if (!checkName(createModel.value.path)) {
    window.$message.error('名称不合法')
    return
  }

  const fullPath = path.value + '/' + createModel.value.path
  file.create(fullPath, createModel.value.dir).then(() => {
    create.value = false
    window.$message.success('创建成功')
    window.$bus.emit('file:refresh')
  })
}

const handleDownload = () => {
  if (!checkName(downloadModel.value.path)) {
    window.$message.error('名称不合法')
    return
  }

  file
    .remoteDownload(path.value + '/' + downloadModel.value.path, downloadModel.value.url)
    .then(() => {
      download.value = false
      window.$message.success('下载任务创建成功')
      window.$bus.emit('file:refresh')
    })
}

const handleCopy = () => {
  if (!selected.value.length) {
    window.$message.error('请选择要复制的文件/文件夹')
    return
  }
  markedType.value = 'copy'
  marked.value = selected.value.map((path) => ({
    name: lastDirectory(path),
    source: path,
    force: false
  }))
  selected.value = []
  window.$message.success('标记成功，请前往目标路径粘贴')
}

const handleMove = () => {
  if (!selected.value.length) {
    window.$message.error('请选择要移动的文件/文件夹')
    return
  }
  markedType.value = 'move'
  marked.value = selected.value.map((path) => ({
    name: lastDirectory(path),
    source: path,
    force: false
  }))
  selected.value = []
  window.$message.success('标记成功，请前往目标路径粘贴')
}

const handleCancel = () => {
  marked.value = []
}

const handlePaste = async () => {
  if (!marked.value.length) {
    window.$message.error('请先标记需要复制或移动的文件/文件夹')
    return
  }

  // 查重
  let flag = false
  const paths = marked.value.map((item) => {
    return {
      name: item.name,
      source: item.source,
      target: path.value + '/' + item.name,
      force: false
    }
  })
  const sources = paths.map((item: any) => item.target)
  await file.exist(sources).then(async (res) => {
    for (let i = 0; i < res.data.length; i++) {
      if (res.data[i]) {
        flag = true
        paths[i].force = true
      }
    }
    if (flag) {
      window.$dialog.warning({
        title: '警告',
        content: `存在同名项
      ${paths
        .filter((item) => item.force)
        .map((item) => item.name)
        .join(', ')} 是否覆盖？`,
        positiveText: '覆盖',
        negativeText: '取消',
        onPositiveClick: async () => {
          if (markedType.value == 'copy') {
            await file.copy(paths).then(() => {
              window.$message.success('复制成功')
            })
          } else {
            await file.move(paths).then(() => {
              window.$message.success('移动成功')
            })
          }
          marked.value = []
          window.$bus.emit('file:refresh')
        },
        onNegativeClick: () => {
          marked.value = []
          window.$message.info('已取消')
        }
      })
    } else {
      if (markedType.value == 'copy') {
        await file.copy(paths).then(() => {
          window.$message.success('复制成功')
        })
      } else {
        await file.move(paths).then(() => {
          window.$message.success('移动成功')
        })
      }
      marked.value = []
      window.$bus.emit('file:refresh')
    }
  })
}

const bulkDelete = async () => {
  if (!selected.value.length) {
    window.$message.error('请选择要删除的文件/文件夹')
    return
  }

  for (const path of selected.value) {
    await file.delete(path).then(() => {
      window.$message.success(`删除 ${path} 成功`)
      window.$bus.emit('file:refresh')
    })
  }

  selected.value = []
}

// 自动填充下载文件名
watch(
  () => downloadModel.value.url,
  (newUrl) => {
    if (!newUrl) return
    try {
      const url = new URL(newUrl)
      const path = url.pathname.split('/').pop()
      if (path) {
        downloadModel.value.path = decodeURIComponent(path)
      }
    } catch (error) {
      /* empty */
    }
  }
)
</script>

<template>
  <n-flex>
    <n-popselect
      :options="[
        { label: '文件', value: 'file' },
        { label: '文件夹', value: 'folder' }
      ]"
      @update:value="showCreate"
    >
      <n-button type="primary"> 创建 </n-button>
    </n-popselect>
    <n-button @click="upload = true"> 上传 </n-button>
    <n-button @click="download = true"> 远程下载 </n-button>
    <div ml-auto>
      <n-flex>
        <n-button v-if="marked.length" secondary type="error" @click="handleCancel">
          取消
        </n-button>
        <n-button v-if="marked.length" secondary type="primary" @click="handlePaste">
          粘贴
        </n-button>
        <n-button-group v-if="selected.length">
          <n-button @click="handleCopy"> 复制 </n-button>
          <n-button @click="handleMove"> 移动 </n-button>
          <n-button @click="compress = true"> 压缩 </n-button>
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
    v-model:show="create"
    preset="card"
    title="创建"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="createModel">
        <n-form-item label="名称">
          <n-input v-model:value="createModel.path" />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleCreate">提交</n-button>
    </n-space>
  </n-modal>
  <n-modal
    v-model:show="download"
    preset="card"
    title="远程下载"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="downloadModel">
        <n-form-item label="下载链接">
          <n-input :input-props="{ type: 'url' }" v-model:value="downloadModel.url" />
        </n-form-item>
        <n-form-item label="保存文件名">
          <n-input v-model:value="downloadModel.path" />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleDownload">提交</n-button>
    </n-space>
  </n-modal>
  <upload-modal v-model:show="upload" v-model:path="path" />
</template>

<style scoped lang="scss"></style>
