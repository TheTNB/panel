<script setup lang="ts">
defineOptions({
  name: 'task-index'
})

import { LogInst, NButton, NDataTable, NPopconfirm } from 'naive-ui'

import task from '@/api/panel/task'
import ws from '@/api/ws'
import { formatDateTime, renderIcon } from '@/utils'
import type { Task } from '@/views/task/types'

const taskLogModal = ref(false)
const taskLog = ref('')
const logRef = ref<LogInst | null>(null)
let logWs: WebSocket | null = null

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: '任务名',
    key: 'name',
    minWidth: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '状态',
    key: 'status',
    width: 150,
    ellipsis: { tooltip: true },
    render(row: any) {
      return row.status === 'finished'
        ? '已完成'
        : row.status === 'waiting'
          ? '等待中'
          : row.status === 'failed'
            ? '已失败'
            : '运行中'
    }
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 200,
    ellipsis: { tooltip: true },
    render(row: any): string {
      return formatDateTime(row.created_at)
    }
  },
  {
    title: '完成时间',
    key: 'updated_at',
    width: 200,
    ellipsis: { tooltip: true },
    render(row: any): string {
      return formatDateTime(row.updated_at)
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
        row.status != 'waiting'
          ? h(
              NButton,
              {
                size: 'small',
                type: 'warning',
                secondary: true,
                onClick: () => {
                  handleShowLog(row.log)
                  taskLogModal.value = true
                }
              },
              {
                default: () => '日志',
                icon: renderIcon('material-symbols:visibility', { size: 14 })
              }
            )
          : null,
        row.status != 'waiting' && row.status != 'running'
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleDelete(row.id)
              },
              {
                default: () => {
                  return '确定要删除吗？'
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
          : null
      ]
    }
  }
]

const tasks = ref<Task[]>([] as Task[])

const selectedRowKeys = ref<any>([])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200]
})

const handleDelete = (id: number) => {
  task.delete(id).then(() => {
    window.$message.success('删除成功')
    onPageChange(pagination.page)
  })
}

const handleShowLog = (path: string) => {
  const cmd = `tail -f ${path}`
  ws.exec(cmd)
    .then((ws: WebSocket) => {
      logWs = ws
      taskLogModal.value = true
      ws.onmessage = (event) => {
        taskLog.value += event.data + '\n'
        const lines = taskLog.value.split('\n')
        if (lines.length > 2000) {
          taskLog.value = lines.slice(lines.length - 2000).join('\n')
        }
      }
    })
    .catch(() => {
      window.$message.error('获取日志流失败')
    })
}

const handleCloseLog = () => {
  if (logWs) {
    logWs.close()
  }
  taskLogModal.value = false
  taskLog.value = ''
}

const fetchTaskList = async (page: number, limit: number) => {
  const { data } = await task.list(page, limit)
  return data
}

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const onPageChange = (page: number) => {
  pagination.page = page
  fetchTaskList(page, pagination.pageSize).then((res) => {
    tasks.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

onMounted(() => {
  onPageChange(pagination.page)
})

watchEffect(() => {
  if (taskLog.value) {
    nextTick(() => {
      logRef.value?.scrollTo({ position: 'bottom', silent: true })
    })
  }
})
</script>

<template>
  <common-page show-footer>
    <n-flex vertical>
      <n-alert type="info">若日志无法加载，请关闭广告拦截应用！</n-alert>
      <n-data-table
        striped
        remote
        :scroll-x="1000"
        :loading="false"
        :columns="columns"
        :data="tasks"
        :row-key="(row: any) => row.id"
        :pagination="pagination"
        @update:checked-row-keys="onChecked"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </n-flex>
  </common-page>
  <n-modal
    v-model:show="taskLogModal"
    preset="card"
    title="任务日志"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="handleCloseLog"
    @mask-click="handleCloseLog"
  >
    <n-log ref="logRef" :log="taskLog" trim :rows="40" />
  </n-modal>
</template>
