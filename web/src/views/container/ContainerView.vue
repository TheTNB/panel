<script setup lang="ts">
import { NButton, NDataTable, NSwitch, NDropdown, NInput } from 'naive-ui'
import type { ContainerList } from '@/views/container/types'
import container from '@/api/panel/container'
import Editor from '@guolao/vue-monaco-editor'
import ContainerCreate from '@/views/container/ContainerCreate.vue'

const data = ref<ContainerList[]>([] as ContainerList[])

const logModal = ref(false)
const logs = ref('')
const renameModal = ref(false)
const renameModel = ref({
  id: '',
  name: ''
})

const containerCreateModal = ref(false)
const selectedRowKeys = ref<any>([])

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const columns: any = [
  { type: 'selection', fixed: 'left' },
  { title: '容器名', key: 'name', width: 150, resizable: true, ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'state',
    width: 100,
    resizable: true,
    render(row: any) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.state === 'running',
        onUpdateValue: (value: boolean) => {
          if (value) {
            handleStart(row.id)
          } else {
            handleStop(row.id)
          }
        }
      })
    }
  },
  { title: '镜像', key: 'image', width: 300, resizable: true, ellipsis: { tooltip: true } },
  {
    title: '端口（主机->容器）',
    key: 'ports',
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return row.ports
        .map((port: any) => {
          return `${port.IP ? port.IP + ':' : ''}${port.PublicPort}->${port.PrivatePort}/${port.Type}`
        })
        .join(', ')
    }
  },
  {
    title: '运行状态',
    key: 'status',
    width: 300,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '操作',
    key: 'actions',
    width: 250,
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
            default: () => '日志'
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'success',
            style: 'margin-left: 15px;',
            onClick: () => {
              renameModel.value.id = row.id
              renameModel.value.name = row.name
              renameModal.value = true
            }
          },
          {
            default: () => '重命名'
          }
        ),
        h(
          NDropdown,
          {
            options: [
              {
                label: '启动',
                key: 'start',
                disabled: row.state === 'running'
              },
              {
                label: '停止',
                key: 'stop',
                disabled: row.state !== 'running'
              },
              {
                label: '重启',
                key: 'restart',
                disabled: row.state !== 'running'
              },
              {
                label: '强制停止',
                key: 'forceStop',
                disabled: row.state !== 'running'
              },
              {
                label: '暂停',
                key: 'pause',
                disabled: row.state !== 'running'
              },
              {
                label: '恢复',
                key: 'unpause',
                disabled: row.state === 'running'
              },
              {
                label: '删除',
                key: 'delete'
              }
            ],
            onSelect: (key: string) => {
              switch (key) {
                case 'start':
                  handleStart(row.id)
                  break
                case 'stop':
                  handleStop(row.id)
                  break
                case 'restart':
                  handleRestart(row.id)
                  break
                case 'forceStop':
                  handleForceStop(row.id)
                  break
                case 'pause':
                  handlePause(row.id)
                  break
                case 'unpause':
                  handleUnpause(row.id)
                  break
                case 'delete':
                  handleDelete(row.id)
                  break
              }
            }
          },
          {
            default: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'primary',
                  style: 'margin-left: 15px;'
                },
                {
                  default: () => '更多'
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

const onPageChange = (page: number) => {
  pagination.page = page
  getContainerList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const getContainerList = async (page: number, pageSize: number) => {
  const { data } = await container.containerList(page, pageSize)
  return data
}

const handleShowLog = async (row: any) => {
  container.containerLogs(row.id).then((res) => {
    logs.value = res.data
    logModal.value = true
  })
}

const handleRename = () => {
  container.containerRename(renameModel.value.id, renameModel.value.name).then(() => {
    window.$message.success('重命名成功')
    renameModal.value = false
    onPageChange(pagination.page)
  })
}

const handleStart = (id: string) => {
  container.containerStart(id).then(() => {
    window.$message.success('启动成功')
    onPageChange(pagination.page)
  })
}

const handleStop = (id: string) => {
  container.containerStop(id).then(() => {
    window.$message.success('停止成功')
    onPageChange(pagination.page)
  })
}

const handleRestart = (id: string) => {
  container.containerRestart(id).then(() => {
    window.$message.success('重启成功')
    onPageChange(pagination.page)
  })
}

const handleForceStop = (id: string) => {
  container.containerKill(id).then(() => {
    window.$message.success('强制停止成功')
    onPageChange(pagination.page)
  })
}

const handlePause = (id: string) => {
  container.containerPause(id).then(() => {
    window.$message.success('暂停成功')
    onPageChange(pagination.page)
  })
}

const handleUnpause = (id: string) => {
  container.containerUnpause(id).then(() => {
    window.$message.success('恢复成功')
    onPageChange(pagination.page)
  })
}

const handleDelete = (id: string) => {
  container.containerRemove(id).then(() => {
    window.$message.success('删除成功')
    onPageChange(pagination.page)
  })
}

const handlePrune = () => {
  container.containerPrune().then(() => {
    window.$message.success('清理成功')
    onPageChange(pagination.page)
  })
}

const bulkStart = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要启动的容器')
    return
  }

  for (const id of selectedRowKeys.value) {
    await container.containerStart(id).then(() => {
      let container = data.value.find((item) => item.id === id)
      window.$message.success(`${container?.name} 启动成功`)
    })
  }

  onPageChange(pagination.page)
}

const bulkStop = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要停止的容器')
    return
  }

  for (const id of selectedRowKeys.value) {
    await container.containerStop(id).then(() => {
      let container = data.value.find((item) => item.id === id)
      window.$message.success(`${container?.name} 停止成功`)
    })
  }

  onPageChange(pagination.page)
}

const bulkRestart = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要重启的容器')
    return
  }

  for (const id of selectedRowKeys.value) {
    await container.containerRestart(id).then(() => {
      let container = data.value.find((item) => item.id === id)
      window.$message.success(`${container?.name} 重启成功`)
    })
  }

  onPageChange(pagination.page)
}

const bulkForceStop = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要强制停止的容器')
    return
  }

  for (const id of selectedRowKeys.value) {
    await container.containerKill(id).then(() => {
      let container = data.value.find((item) => item.id === id)
      window.$message.success(`${container?.name} 强制停止成功`)
    })
  }

  onPageChange(pagination.page)
}

const bulkDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要删除的容器')
    return
  }

  for (const id of selectedRowKeys.value) {
    await container.containerRemove(id).then(() => {
      let container = data.value.find((item) => item.id === id)
      window.$message.success(`${container?.name} 删除成功`)
    })
  }

  onPageChange(pagination.page)
}

const bulkPause = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要暂停的容器')
    return
  }

  for (const id of selectedRowKeys.value) {
    await container.containerPause(id).then(() => {
      let container = data.value.find((item) => item.id === id)
      window.$message.success(`${container?.name} 暂停成功`)
    })
  }

  onPageChange(pagination.page)
}

const bulkUnpause = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要恢复的容器')
    return
  }

  for (const id of selectedRowKeys.value) {
    await container.containerUnpause(id).then(() => {
      let container = data.value.find((item) => item.id === id)
      window.$message.success(`${container?.name} 恢复成功`)
    })
  }

  onPageChange(pagination.page)
}

const closeContainerCreateModal = () => {
  containerCreateModal.value = false
  onPageChange(pagination.page)
}

onMounted(() => {
  onPageChange(pagination.page)
})
</script>

<template>
  <n-space vertical size="large">
    <n-card rounded-10>
      <n-space>
        <n-button type="primary" @click="containerCreateModal = true">创建容器</n-button>
        <n-button type="primary" @click="handlePrune" ghost>清理容器</n-button>
        <n-button-group>
          <n-button @click="bulkStart">启动</n-button>
          <n-button @click="bulkStop">停止</n-button>
          <n-button @click="bulkRestart">重启</n-button>
          <n-button @click="bulkForceStop">强制停止</n-button>
          <n-button @click="bulkPause">暂停</n-button>
          <n-button @click="bulkUnpause">恢复</n-button>
          <n-button @click="bulkDelete">删除</n-button>
        </n-button-group>
      </n-space>
    </n-card>
    <n-card rounded-10>
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
        @update:checked-row-keys="onChecked"
      />
    </n-card>
  </n-space>
  <n-modal
    v-model:show="logModal"
    preset="card"
    title="日志"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <Editor
      v-model:value="logs"
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
    v-model:show="renameModal"
    preset="card"
    title="重命名"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form :model="renameModel">
      <n-form-item path="name" label="新名称">
        <n-input
          v-model:value="renameModel.name"
          type="text"
          @keydown.enter.prevent
          placeholder="输入新名称"
        />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button type="info" block @click="handleRename">提交</n-button>
      </n-col>
    </n-row>
  </n-modal>
  <ContainerCreate :show="containerCreateModal" @close="closeContainerCreateModal" />
</template>
