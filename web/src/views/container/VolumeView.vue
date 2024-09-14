<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm } from 'naive-ui'
import type { VolumeList } from '@/views/container/types'
import container from '@/api/panel/container'

const createModel = ref({
  name: '',
  driver: 'local',
  options: [],
  labels: []
})

const options = [{ label: 'local', value: 'local' }]

const createModal = ref(false)
const loading = ref(false)

const data = ref<VolumeList[]>([] as VolumeList[])
const selectedRowKeys = ref<any>([])

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const columns: any = [
  { type: 'selection', fixed: 'left' },
  { title: '名称', key: 'id', width: 150, resizable: true, ellipsis: { tooltip: true } },
  {
    title: '驱动',
    key: 'driver',
    width: 100,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '范围',
    key: 'scope',
    width: 100,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '挂载点',
    key: 'mount',
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '创建时间',
    key: 'created',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NPopconfirm,
          {
            onPositiveClick: async () => {
              await handleDelete(row)
            }
          },
          {
            default: () => {
              return '确定删除吗？'
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'error'
                },
                {
                  default: () => '删除'
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
  getVolumeList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const getVolumeList = async (page: number, pageSize: number) => {
  const { data } = await container.volumeList(page, pageSize)
  return data
}

const handleDelete = async (row: any) => {
  container.volumeRemove(row.id).then(() => {
    window.$message.success('删除成功')
    onPageChange(pagination.page)
  })
}

const handlePrune = () => {
  container.volumePrune().then(() => {
    window.$message.success('清理成功')
    onPageChange(pagination.page)
  })
}

const handleCreate = () => {
  loading.value = true
  container
    .volumeCreate(createModel.value)
    .then(() => {
      window.$message.success('创建成功')
      onPageChange(pagination.page)
    })
    .finally(() => {
      loading.value = false
      createModal.value = false
    })
}

onMounted(() => {
  onPageChange(pagination.page)
})
</script>

<template>
  <n-space vertical size="large">
    <n-card rounded-10>
      <n-space>
        <n-button type="primary" @click="createModal = true">创建卷</n-button>
        <n-button type="primary" @click="handlePrune" ghost>清理卷</n-button>
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
    v-model:show="createModal"
    preset="card"
    title="创建卷"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form :model="createModel">
      <n-form-item path="name" label="卷名">
        <n-input v-model:value="createModel.name" type="text" @keydown.enter.prevent />
      </n-form-item>
      <n-form-item path="driver" label="驱动">
        <n-select
          :options="options"
          v-model:value="createModel.driver"
          type="text"
          @keydown.enter.prevent
        >
        </n-select>
      </n-form-item>
      <n-form-item path="env" label="标签">
        <n-dynamic-input
          v-model:value="createModel.labels"
          preset="pair"
          key-placeholder="标签名"
          value-placeholder="标签值"
        />
      </n-form-item>
      <n-form-item path="env" label="选项">
        <n-dynamic-input
          v-model:value="createModel.options"
          preset="pair"
          key-placeholder="选项名"
          value-placeholder="选项值"
        />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button type="info" :loading="loading" :disabled="loading" block @click="handleCreate">
          提交
        </n-button>
      </n-col>
    </n-row>
  </n-modal>
</template>
