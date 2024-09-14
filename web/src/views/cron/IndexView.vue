<script setup lang="ts">
import info from '@/api/panel/info'
import website from '@/api/panel/website'
import cron from '@/api/panel/cron'
import { NButton, NDataTable, NInput, NPopconfirm, NSwitch } from 'naive-ui'
import type { CronTask } from '@/views/cron/types'
import { renderIcon } from '@/utils'
import Editor from '@guolao/vue-monaco-editor'

const addModel = ref({
  name: '',
  type: 'shell',
  backup_type: 'website',
  website: '',
  database: '',
  save: 1,
  backup_path: '',
  script: '# 在此输入你要执行的脚本内容',
  time: '* * * * *'
})

const cronSelectModal = ref(false)
const taskLogModal = ref(false)
const editTaskModal = ref(false)

const installedDbAndPhp = ref({
  php: [
    {
      label: '',
      value: ''
    }
  ],
  db: [
    {
      label: '',
      value: ''
    }
  ]
})

const mySQLInstalled = computed(() => {
  return installedDbAndPhp.value.db.find((item) => item.value === 'mysql')
})

const postgreSQLInstalled = computed(() => {
  return installedDbAndPhp.value.db.find((item) => item.value === 'postgresql')
})

const websites = ref<any>([])

const columns: any = [
  { type: 'selection', fixed: 'left' },
  { title: '任务名', key: 'name', width: 150, resizable: true, ellipsis: { tooltip: true } },
  {
    title: '任务类型',
    key: 'type',
    width: 100,
    resizable: true,
    render(row: any) {
      return row.type === 'shell' ? '运行脚本' : row.type === 'backup' ? '备份数据' : '切割日志'
    }
  },
  {
    title: '启用',
    key: 'status',
    width: 60,
    align: 'center',
    resizable: true,
    render(row: any) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.status,
        onUpdateValue: () => handleStatusChange(row)
      })
    }
  },
  { title: '任务周期', key: 'time', width: 100, resizable: true, ellipsis: { tooltip: true } },
  {
    title: '创建时间',
    key: 'created_at',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  { title: '最后更新时间', key: 'updated_at', ellipsis: { tooltip: true } },
  {
    title: '操作',
    key: 'actions',
    width: 280,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'warning',
            secondary: true,
            onClick: () => handleShowLog(row)
          },
          {
            default: () => '日志',
            icon: renderIcon('majesticons:eye-line', { size: 14 })
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 15px;',
            onClick: () => handleEdit(row)
          },
          {
            default: () => '修改',
            icon: renderIcon('material-symbols:edit-outline', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(row.id)
          },
          {
            default: () => {
              return '确定删除任务吗？'
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

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
})

const data = ref<CronTask[]>([] as CronTask[])

const taskLog = ref('')
const editTask = ref({
  id: 0,
  name: '',
  time: '',
  script: ''
})

const getTaskList = async (page: number, limit: number) => {
  const { data } = await cron.list(page, limit)
  return data
}

const getWebsiteList = async (page: number, limit: number) => {
  const { data } = await website.list(page, limit)
  for (const item of data.items) {
    websites.value.push({
      label: item.name,
      value: item.name
    })
  }
  addModel.value.website = websites.value[0]?.value
}

const getPhpAndDb = async () => {
  const { data } = await info.installedDbAndPhp()
  installedDbAndPhp.value = data
}

const onPageChange = (page: number) => {
  pagination.page = page
  getTaskList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const handleStatusChange = async (row: any) => {
  cron.status(row.id, !row.status).then(() => {
    row.status = !row.status
    window.$message.success('修改成功')
  })
}

const handleShowLog = async (row: any) => {
  cron.log(row.id).then((res) => {
    taskLog.value = res.data
    taskLogModal.value = true
  })
}

const handleAdd = async () => {
  await cron.add(addModel.value).then(() => {
    window.$message.success('添加成功')
  })
  onPageChange(pagination.page)
}

const handleEdit = async (row: any) => {
  cron.script(row.id).then((res) => {
    editTask.value.id = row.id
    editTask.value.name = row.name
    editTask.value.time = row.time
    editTask.value.script = res.data
    editTaskModal.value = true
  })
}

const handleDelete = async (id: number) => {
  await cron.delete(id).then(() => {
    window.$message.success('删除成功')
  })
  onPageChange(pagination.page)
}

const saveTaskEdit = async () => {
  cron
    .update(editTask.value.id, editTask.value.name, editTask.value.time, editTask.value.script)
    .then(() => {
      window.$message.success('修改成功')
    })
}

const handleCronSelectUpdate = (value: string) => {
  if (editTaskModal.value) {
    editTask.value.time = value
    return
  }

  addModel.value.time = value
}

onMounted(() => {
  getPhpAndDb()
  getWebsiteList(1, 10000)
  onPageChange(pagination.page)
})
</script>

<template>
  <common-page show-footer>
    <n-space vertical>
      <n-card flex-1 rounded-10 title="添加计划任务">
        <n-space vertical>
          <n-alert type="info">
            面板的计划任务均基于脚本运行，若任务类型满足不了需求，可自行修改对应的脚本。
          </n-alert>
          <n-form>
            <n-form-item label="任务类型">
              <n-select
                v-model:value="addModel.type"
                :options="[
                  { label: '运行脚本', value: 'shell' },
                  { label: '备份数据', value: 'backup' },
                  { label: '切割日志', value: 'cutoff' }
                ]"
              >
              </n-select>
            </n-form-item>
            <n-form-item label="任务名称">
              <n-input v-model:value="addModel.name" placeholder="任务名称" />
            </n-form-item>
            <n-form-item label="任务周期">
              <n-input
                v-model:value="addModel.time"
                placeholder="* * * * *"
                @click="cronSelectModal = true"
              />
            </n-form-item>
            <div v-if="addModel.type === 'shell'">
              <n-text>脚本内容</n-text>
              <Editor
                v-model:value="addModel.script"
                language="shell"
                theme="vs-dark"
                height="40vh"
                mt-8
                :options="{
                  automaticLayout: true,
                  formatOnType: true,
                  formatOnPaste: true
                }"
              />
            </div>
            <n-form-item v-if="addModel.type === 'backup'" label="备份类型">
              <n-radio-group v-model:value="addModel.backup_type">
                <n-radio value="website">网站目录</n-radio>
                <n-radio value="mysql" :disabled="!mySQLInstalled"> MySQL 数据库</n-radio>
                <n-radio value="postgresql" :disabled="!postgreSQLInstalled">
                  PostgreSQL 数据库
                </n-radio>
              </n-radio-group>
            </n-form-item>
            <n-form-item
              v-if="
                (addModel.backup_type === 'website' && addModel.type === 'backup') ||
                addModel.type === 'cutoff'
              "
              label="选择网站"
            >
              <n-select
                v-model:value="addModel.website"
                :options="websites"
                placeholder="选择网站"
              />
            </n-form-item>
            <n-form-item
              v-if="addModel.backup_type !== 'website' && addModel.type === 'backup'"
              label="数据库名"
            >
              <n-input v-model:value="addModel.database" placeholder="数据库名" />
            </n-form-item>
            <n-form-item v-if="addModel.type === 'backup'" label="保存目录">
              <n-input v-model:value="addModel.backup_path" placeholder="保存目录" />
            </n-form-item>
            <n-form-item v-if="addModel.type !== 'shell'" label="保留份数">
              <n-input-number v-model:value="addModel.save" />
            </n-form-item>
          </n-form>
          <n-button type="primary" @click="handleAdd">添加</n-button>
        </n-space>
      </n-card>
      <n-card flex-1 rounded-10 title="计划任务列表">
        <n-data-table
          striped
          remote
          :data="data"
          :columns="columns"
          :row-key="(row: any) => row.id"
          :pagination="pagination"
          :bordered="false"
          :loading="false"
          @update:page="onPageChange"
          @update:page-size="onPageSizeChange"
        />
      </n-card>
    </n-space>
  </common-page>
  <n-modal
    v-model:show="cronSelectModal"
    preset="card"
    title="Cron 表达式生成"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <cron-select @update="handleCronSelectUpdate" />
  </n-modal>
  <n-modal
    v-model:show="taskLogModal"
    preset="card"
    title="任务日志"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <Editor
      v-model:value="taskLog"
      language="ini"
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
  <n-modal
    v-model:show="editTaskModal"
    preset="card"
    title="编辑任务"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="saveTaskEdit"
  >
    <n-form inline>
      <n-form-item label="任务名称">
        <n-input v-model:value="editTask.name" placeholder="任务名称" />
      </n-form-item>
      <n-form-item label="任务周期">
        <n-input
          v-model:value="editTask.time"
          placeholder="* * * * *"
          @click="cronSelectModal = true"
        />
      </n-form-item>
    </n-form>
    <Editor
      v-model:value="editTask.script"
      language="shell"
      theme="vs-dark"
      height="60vh"
      mt-8
      :options="{
        automaticLayout: true,
        formatOnType: true,
        formatOnPaste: true
      }"
    />
  </n-modal>
</template>
