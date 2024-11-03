<script setup lang="ts">
import backup from '@/api/panel/backup'
import { renderIcon } from '@/utils'
import type { MessageReactive } from 'naive-ui'
import { NButton, NInput, NPopconfirm } from 'naive-ui'

import { formatDateTime } from '@/utils'
import type { Backup } from './types'

const type = defineModel<string>('type', { type: String, required: true })

let messageReactive: MessageReactive | null = null

const restoreModal = ref(false)
const restoreModel = ref({
  file: '',
  target: ''
})

const columns: any = [
  {
    title: '文件名',
    key: 'name',
    minWidth: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '大小',
    key: 'size',
    width: 160,
    ellipsis: { tooltip: true }
  },
  {
    title: '更新日期',
    key: 'time',
    width: 200,
    ellipsis: { tooltip: true },
    render(row: any) {
      return formatDateTime(row.time)
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    align: 'center',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'warning',
            secondary: true,
            onClick: () => {
              restoreModel.value.file = row.path
              restoreModal.value = true
            }
          },
          {
            default: () => '恢复',
            icon: renderIcon('material-symbols:settings-backup-restore-rounded', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row.name)
          },
          {
            default: () => {
              return '确定删除备份吗？'
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  style: 'margin-left: 15px;'
                },
                {
                  default: () => '删除',
                  icon: renderIcon('material-symbols:delete-outline', { size: 14 })
                }
              )
            }
          }
        )
      ]
    }
  }
]

const data = ref<Backup[]>([])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200]
})

const getList = async (page: number, limit: number) => {
  const { data } = await backup.list(type.value, page, limit)
  return data
}

const onPageChange = (page: number) => {
  pagination.page = page
  getList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const handleRestore = async () => {
  messageReactive = window.$message.loading('恢复中...', {
    duration: 0
  })
  await backup.restore(type.value, restoreModel.value.file, restoreModel.value.target).then(() => {
    messageReactive?.destroy()
    window.$message.success('恢复成功')
    onPageChange(pagination.page)
  })
}

const handleDelete = async (file: string) => {
  await backup.delete(type.value, file).then(() => {
    window.$message.success('删除成功')
    onPageChange(pagination.page)
  })
}

onMounted(() => {
  onPageChange(pagination.page)
  window.$bus.on('backup:refresh', () => {
    onPageChange(pagination.page)
  })
})

onUnmounted(() => {
  window.$bus.off('backup:refresh')
})
</script>

<template>
  <n-data-table
    striped
    remote
    :scroll-x="1000"
    :loading="false"
    :columns="columns"
    :data="data"
    :row-key="(row: any) => row.name"
    @update:page="onPageChange"
    @update:page-size="onPageSizeChange"
  />
  <n-modal
    v-model:show="restoreModal"
    preset="card"
    title="恢复备份"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="restoreModal = false"
  >
    <n-form :model="restoreModel">
      <n-form-item path="name" label="恢复目标">
        <n-input v-model:value="restoreModel.target" type="text" @keydown.enter.prevent />
      </n-form-item>
    </n-form>
    <n-button type="info" block @click="handleRestore">提交</n-button>
  </n-modal>
</template>

<style scoped lang="scss"></style>
