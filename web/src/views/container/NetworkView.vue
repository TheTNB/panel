<script setup lang="ts">
import { NButton, NDataTable, NFlex, NInput, NPopconfirm, NTag } from 'naive-ui'
import type { NetworkList } from '@/views/container/types'
import container from '@/api/panel/container'

const createModel = ref({
  name: '',
  driver: 'bridge',
  ipv4: {
    enabled: true,
    subnet: '',
    gateway: '',
    ip_range: ''
  },
  ipv6: {
    enabled: false,
    subnet: '',
    gateway: '',
    ip_range: ''
  },
  options: [],
  labels: []
})

const options = [
  { label: 'bridge', value: 'bridge' },
  { label: 'host', value: 'host' },
  { label: 'overlay', value: 'overlay' },
  { label: 'macvlan', value: 'macvlan' },
  { label: 'ipvlan', value: 'ipvlan' },
  { label: 'none', value: 'none' }
]

const createModal = ref(false)
const loading = ref(false)

const data = ref<NetworkList[]>([] as NetworkList[])
const selectedRowKeys = ref<any>([])

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const columns: any = [
  { type: 'selection', fixed: 'left' },
  { title: '名称', key: 'name', width: 150, resizable: true, ellipsis: { tooltip: true } },
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
    title: '子网',
    key: 'subnet',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any): any {
      return h(NFlex, null, {
        default: () =>
          row.ipam.config.map((tag: any) =>
            h(NTag, null, {
              default: () => tag.subnet
            })
          )
      })
    }
  },
  {
    title: '网关',
    key: 'gateway',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any): any {
      return h(NFlex, null, {
        default: () =>
          row.ipam.config.map((tag: any) =>
            h(NTag, null, {
              default: () => tag.gateway
            })
          )
      })
    }
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
  getNetworkList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const getNetworkList = async (page: number, pageSize: number) => {
  const { data } = await container.networkList(page, pageSize)
  return data
}

const handleDelete = async (row: any) => {
  container.networkRemove(row.id).then(() => {
    window.$message.success('删除成功')
    onPageChange(pagination.page)
  })
}

const handlePrune = () => {
  container.networkPrune().then(() => {
    window.$message.success('清理成功')
    onPageChange(pagination.page)
  })
}

const handleCreate = () => {
  loading.value = true
  container
    .networkCreate(createModel.value)
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
        <n-button type="primary" @click="createModal = true">创建网络</n-button>
        <n-button type="primary" @click="handlePrune" ghost>清理网络</n-button>
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
    title="创建网络"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-form :model="createModel">
      <n-form-item path="name" label="网络名">
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
      <n-form-item path="ipv4" label="IPV4">
        <n-switch v-model:value="createModel.ipv4.enabled" />
      </n-form-item>
      <n-form-item v-if="createModel.ipv4.enabled" path="subnet" label="子网">
        <n-input
          v-model:value="createModel.ipv4.subnet"
          type="text"
          @keydown.enter.prevent
          placeholder="172.16.10.0/24"
        />
      </n-form-item>
      <n-form-item v-if="createModel.ipv4.enabled" path="gateway" label="网关">
        <n-input
          v-model:value="createModel.ipv4.gateway"
          type="text"
          @keydown.enter.prevent
          placeholder="172.16.10.254"
        />
      </n-form-item>
      <n-form-item v-if="createModel.ipv4.enabled" path="ip_range" label="IP范围">
        <n-input
          v-model:value="createModel.ipv4.ip_range"
          type="text"
          @keydown.enter.prevent
          placeholder="172.16.10.0/24"
        />
      </n-form-item>
      <n-form-item path="ipv6" label="IPV6">
        <n-switch v-model:value="createModel.ipv6.enabled" />
      </n-form-item>
      <n-form-item v-if="createModel.ipv6.enabled" path="subnet" label="子网">
        <n-input
          v-model:value="createModel.ipv6.subnet"
          type="text"
          @keydown.enter.prevent
          placeholder="2408:400e::/48"
        />
      </n-form-item>
      <n-form-item v-if="createModel.ipv6.enabled" path="gateway" label="网关">
        <n-input
          v-model:value="createModel.ipv6.gateway"
          type="text"
          @keydown.enter.prevent
          placeholder="2408:400e::1"
        />
      </n-form-item>
      <n-form-item v-if="createModel.ipv6.enabled" path="ip_range" label="IP范围">
        <n-input
          v-model:value="createModel.ipv6.ip_range"
          type="text"
          @keydown.enter.prevent
          placeholder="2408:400e::/64"
        />
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
