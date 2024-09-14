<script setup lang="ts">
import type { Task } from '@/views/task/types'
import { NButton, NDataTable, NPopconfirm } from 'naive-ui'
import { renderIcon } from '@/utils'
import task from '@/api/panel/task'
import Editor from '@guolao/vue-monaco-editor'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const taskLogModal = ref(false)
const taskLog = ref('')

const autoRefresh = ref(false)
const currentTaskId = ref(0)

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: t('taskIndex.columns.name'),
    key: 'name',
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('taskIndex.columns.status'),
    key: 'status',
    width: 100,
    ellipsis: { tooltip: true },
    render(row: any) {
      return row.status === 'finished'
        ? t('taskIndex.options.status.finished')
        : row.status === 'waiting'
          ? t('taskIndex.options.status.waiting')
          : row.status === 'failed'
            ? t('taskIndex.options.status.failed')
            : t('taskIndex.options.status.running')
    }
  },
  {
    title: t('taskIndex.columns.createdAt'),
    key: 'created_at',
    width: 160,
    ellipsis: { tooltip: true }
  },
  {
    title: t('taskIndex.columns.updatedAt'),
    key: 'updated_at',
    width: 160,
    ellipsis: { tooltip: true }
  },
  {
    title: t('taskIndex.columns.actions'),
    key: 'actions',
    width: 200,
    align: 'center',
    fixed: 'right',
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
                  handleShowLog(row.id)
                  currentTaskId.value = row.id
                  taskLogModal.value = true
                  autoRefresh.value = true
                }
              },
              {
                default: () => t('taskIndex.buttons.log'),
                icon: renderIcon('material-symbols:visibility', { size: 14 })
              }
            )
          : null,
        row.status != 'waiting' && row.status != 'running'
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleDelete(row.id),
                onNegativeClick: () => {
                  window.$message.info(t('taskIndex.buttons.undelete'))
                }
              },
              {
                default: () => {
                  return t('taskIndex.confirm.delete')
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
                      default: () => t('taskIndex.buttons.delete'),
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
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
})

const handleDelete = (id: number) => {
  task.delete(id).then(() => {
    window.$message.success(t('taskIndex.alerts.delete'))
    onPageChange(pagination.page)
  })
}

const handleShowLog = (id: number) => {
  task
    .log(id)
    .then((res) => {
      taskLog.value = res.data
    })
    .catch(() => {
      autoRefresh.value = false
    })
}

const getTaskList = async (page: number, limit: number) => {
  const { data } = await task.list(page, limit)
  return data
}

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const onPageChange = (page: number) => {
  pagination.page = page
  getTaskList(page, pagination.pageSize).then((res) => {
    tasks.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

let timer: any = null
const setAutoRefreshTimer = () => {
  timer = setInterval(() => {
    handleShowLog(currentTaskId.value)
  }, 2000)
}

watch(
  () => autoRefresh.value,
  (value) => {
    if (value) {
      setAutoRefreshTimer()
    } else {
      clearInterval(timer)
    }
  },
  { immediate: true }
)

onMounted(() => {
  onPageChange(pagination.page)
})
onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <common-page show-footer>
    <n-flex vertical>
      <n-alert type="info">若日志无法加载，请关闭广告拦截插件！</n-alert>
      <n-data-table
        striped
        remote
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
    :title="$t('taskIndex.logModal.title')"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="
      () => {
        autoRefresh = false
        taskLogModal = false
      }
    "
    @mask-click="
      () => {
        autoRefresh = false
        taskLogModal = false
      }
    "
  >
    <template #header-extra>
      <n-switch v-model:value="autoRefresh" style="margin-right: 10px">
        <template #checked>{{ $t('taskIndex.logModal.autoRefresh.on') }}</template>
        <template #unchecked>{{ $t('taskIndex.logModal.autoRefresh.off') }}</template>
      </n-switch>
    </template>
    <Editor
      v-model:value="taskLog"
      language="shell"
      theme="vs-dark"
      height="60vh"
      mt-8
      :options="{
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true,
        readOnly: true
      }"
    />
  </n-modal>
</template>
