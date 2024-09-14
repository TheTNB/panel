<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm } from 'naive-ui'
import supervisor from '@/api/plugins/supervisor'
import service from '@/api/panel/system/service'
import { renderIcon } from '@/utils'
import type { Process } from '@/views/plugins/supervisor/types'
import Editor from '@guolao/vue-monaco-editor'

const currentTab = ref('status')
const serviceName = ref('supervisor')
const status = ref(false)
const isEnabled = ref(false)
const config = ref('')
const log = ref('')
const processLog = ref('')

const addProcessModal = ref(false)
const addProcessModel = ref({
  name: '',
  user: 'www',
  path: '',
  command: '',
  num: 1
})

const editProcessModal = ref(false)
const editProcessModel = ref({
  process: '',
  config: ''
})

const processLogModal = ref(false)

const statusType = computed(() => {
  return status.value ? 'success' : 'error'
})
const statusStr = computed(() => {
  return status.value ? '正常运行中' : '已停止运行'
})

const processColumns: any = [
  {
    title: '名称',
    key: 'name',
    fixed: 'left',
    resizable: true,
    ellipsis: { tooltip: true }
  },
  { title: '状态', key: 'status', width: 100, resizable: true, ellipsis: { tooltip: true } },
  { title: 'PID', key: 'pid', width: 100, resizable: true, ellipsis: { tooltip: true } },
  { title: '运行时间', key: 'uptime', width: 150, resizable: true, ellipsis: { tooltip: true } },
  {
    title: '操作',
    key: 'actions',
    width: 500,
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
            onClick: () => handleShowProcessLog(row)
          },
          {
            default: () => '日志',
            icon: renderIcon('material-symbols:visibility', { size: 14 })
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'info',
            style: 'margin-left: 15px',
            onClick: () => handleEditProcess(row.name)
          },
          {
            default: () => '配置',
            icon: renderIcon('material-symbols:settings-outline', { size: 14 })
          }
        ),
        row.status != 'RUNNING'
          ? h(
              NButton,
              {
                size: 'small',
                type: 'primary',
                secondary: true,
                style: 'margin-left: 15px',
                onClick: () => handleProcessStart(row.name)
              },
              {
                default: () => '启动',
                icon: renderIcon('material-symbols:play-arrow-outline', { size: 18 })
              }
            )
          : null,
        row.status == 'RUNNING'
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleProcessStop(row.name)
              },
              {
                default: () => {
                  return '确定停止进程' + row.name + '吗？'
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'warning',
                      style: 'margin-left: 15px'
                    },
                    {
                      default: () => '停止',
                      icon: renderIcon('material-symbols:stop-outline', { size: 18 })
                    }
                  )
                }
              }
            )
          : null,
        row.status == 'RUNNING'
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleProcessRestart(row.name)
              },
              {
                default: () => {
                  return '确定重启进程' + row.name + '吗？'
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'primary',
                      style: 'margin-left: 15px'
                    },
                    {
                      default: () => '重启',
                      icon: renderIcon('material-symbols:replay', { size: 18 })
                    }
                  )
                }
              }
            )
          : null,
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleProcessDelete(row.name)
          },
          {
            default: () => {
              return '确定删除进程' + row.name + '吗？'
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  style: 'margin-left: 15px'
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

const processes = ref<Process[]>([])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
})

const getProcesses = async (page: number, limit: number) => {
  const { data } = await supervisor.processes(page, limit)
  return data
}

const onPageChange = (page: number) => {
  pagination.page = page
  getProcesses(page, pagination.pageSize).then((res) => {
    processes.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const getStatus = async () => {
  await service.status(serviceName.value).then((res: any) => {
    status.value = res.data
  })
}

const getIsEnabled = async () => {
  await service.isEnabled(serviceName.value).then((res: any) => {
    isEnabled.value = res.data
  })
}

const getLog = async () => {
  supervisor.log().then((res: any) => {
    log.value = res.data
  })
}

const getConfig = async () => {
  supervisor.config().then((res: any) => {
    config.value = res.data
  })
}

const handleSaveConfig = async () => {
  await supervisor.saveConfig(config.value)
  window.$message.success('保存成功')
  await getLog()
}

const handleClearLog = async () => {
  await supervisor.clearLog()
  window.$message.success('清空成功')
}

const handleStart = async () => {
  await service.start(serviceName.value)
  window.$message.success('启动成功')
  await getStatus()
  await getLog()
}

const handleIsEnabled = async () => {
  if (isEnabled.value) {
    await service.enable(serviceName.value)
    window.$message.success('开启自启动成功')
  } else {
    await service.disable(serviceName.value)
    window.$message.success('禁用自启动成功')
  }
  await getIsEnabled()
}

const handleStop = async () => {
  await service.stop(serviceName.value)
  window.$message.success('停止成功')
  await getStatus()
  await getLog()
}

const handleRestart = async () => {
  await service.restart(serviceName.value)
  window.$message.success('重启成功')
  await getStatus()
  await getLog()
}

const handleReload = async () => {
  await service.reload(serviceName.value)
  window.$message.success('重载成功')
  await getStatus()
  await getLog()
}

const handleAddProcess = async () => {
  await supervisor.addProcess(addProcessModel.value)
  window.$message.success('添加成功')
  addProcessModal.value = false
  onPageChange(1)
}

const handleProcessStart = async (name: string) => {
  await supervisor.startProcess(name)
  window.$message.success('启动成功')
}

const handleProcessStop = async (name: string) => {
  await supervisor.stopProcess(name)
  window.$message.success('停止成功')
}

const handleProcessRestart = async (name: string) => {
  await supervisor.restartProcess(name)
  window.$message.success('重启成功')
}

const handleProcessDelete = async (name: string) => {
  await supervisor.deleteProcess(name)
  window.$message.success('删除成功')
  onPageChange(1)
}

const handleShowProcessLog = (row: any) => {
  supervisor.processLog(row.name).then((res) => {
    processLogModal.value = true
    processLog.value = res.data
  })
}

const handleEditProcess = async (name: string) => {
  await getProcessConfig(name)
  editProcessModal.value = true
}

const getProcessConfig = async (name: string) => {
  editProcessModel.value.process = name
  await supervisor.processConfig(name).then((res: any) => {
    editProcessModel.value.config = res.data
  })
}

const handleSaveProcessConfig = async () => {
  await supervisor.saveProcessConfig(editProcessModel.value.process, editProcessModel.value.config)
  window.$message.success('保存成功')
}

let timer: any = null

onMounted(async () => {
  await supervisor.service().then((res: any) => {
    serviceName.value = res.data
  })
  await getStatus()
  await getIsEnabled()
  await getConfig()
  onPageChange(1)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button
        v-if="currentTab == 'config'"
        class="ml-16"
        type="primary"
        @click="handleSaveConfig"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'processes'"
        class="ml-16"
        type="primary"
        @click="addProcessModal = true"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:add" />
        添加进程
      </n-button>
      <n-button v-if="currentTab == 'log'" class="ml-16" type="primary" @click="handleClearLog">
        <TheIcon :size="18" class="mr-5" icon="material-symbols:delete-outline" />
        清空日志
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="status" tab="运行状态">
        <n-space vertical>
          <n-card title="运行状态" rounded-10>
            <template #header-extra>
              <n-switch v-model:value="isEnabled" @update:value="handleIsEnabled">
                <template #checked> 自启动开</template>
                <template #unchecked> 自启动关</template>
              </n-switch>
            </template>
            <n-space vertical>
              <n-alert :type="statusType">
                {{ statusStr }}
              </n-alert>
              <n-space>
                <n-button type="success" @click="handleStart">
                  <TheIcon
                    :size="24"
                    class="mr-5"
                    icon="material-symbols:play-arrow-outline-rounded"
                  />
                  启动
                </n-button>
                <n-popconfirm @positive-click="handleStop">
                  <template #trigger>
                    <n-button type="error">
                      <TheIcon
                        :size="24"
                        class="mr-5"
                        icon="material-symbols:stop-outline-rounded"
                      />
                      停止
                    </n-button>
                  </template>
                  停止 Supervisor 会导致 Supervisor 管理的所有进程被杀死，确定要停止吗？
                </n-popconfirm>
                <n-button type="warning" @click="handleRestart">
                  <TheIcon :size="18" class="mr-5" icon="material-symbols:replay-rounded" />
                  重启
                </n-button>
                <n-button type="primary" @click="handleReload">
                  <TheIcon :size="20" class="mr-5" icon="material-symbols:refresh-rounded" />
                  重载
                </n-button>
              </n-space>
            </n-space>
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="processes" tab="进程管理">
        <n-card title="进程列表" :segmented="true" rounded-10>
          <n-data-table
            striped
            remote
            :loading="false"
            :columns="processColumns"
            :data="processes"
            :row-key="(row: any) => row.name"
            @update:page="onPageChange"
            @update:page-size="onPageSizeChange"
          />
        </n-card>
      </n-tab-pane>
      <n-tab-pane name="config" tab="主配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 Supervisor 主配置文件，如果你不了解各参数的含义，请不要随意修改！
          </n-alert>
          <Editor
            v-model:value="config"
            language="ini"
            theme="vs-dark"
            height="60vh"
            mt-8
            :options="{
              automaticLayout: true,
              formatOnType: true,
              formatOnPaste: true
            }"
          />
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="log" tab="日志">
        <Editor
          v-model:value="log"
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
      </n-tab-pane>
    </n-tabs>
  </common-page>
  <n-modal
    v-model:show="addProcessModal"
    preset="card"
    title="添加进程"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="addProcessModal = false"
  >
    <n-form :model="addProcessModel">
      <n-form-item path="name" label="名称">
        <n-input
          v-model:value="addProcessModel.name"
          type="text"
          @keydown.enter.prevent
          placeholder="名称禁止使用中文"
        />
      </n-form-item>
      <n-form-item path="command" label="启动命令">
        <n-input
          v-model:value="addProcessModel.command"
          type="text"
          @keydown.enter.prevent
          placeholder="启动命令中的文件请填写绝对路径"
        />
      </n-form-item>
      <n-form-item path="path" label="运行目录">
        <n-input
          v-model:value="addProcessModel.path"
          type="text"
          @keydown.enter.prevent
          placeholder="运行目录请填写绝对路径"
        />
      </n-form-item>
      <n-form-item path="user" label="启动用户">
        <n-input
          v-model:value="addProcessModel.user"
          type="text"
          @keydown.enter.prevent
          placeholder="一般情况下填写www即可"
        />
      </n-form-item>
      <n-form-item path="num" label="进程数量">
        <n-input-number v-model:value="addProcessModel.num" :min="1" />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button type="info" block @click="handleAddProcess">提交</n-button>
      </n-col>
    </n-row>
  </n-modal>
  <n-modal
    v-model:show="processLogModal"
    preset="card"
    title="进程日志"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <Editor
      v-model:value="processLog"
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
    v-model:show="editProcessModal"
    preset="card"
    title="进程配置"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="handleSaveProcessConfig"
  >
    <Editor
      v-model:value="editProcessModel.config"
      language="ini"
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
