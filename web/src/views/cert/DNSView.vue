<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm, NSpace, NTag } from 'naive-ui'
import cert from '@/api/panel/cert'
import type { DNS } from '@/views/cert/types'

const addDNSModel = ref<any>({
  data: {
    token: '',
    id: '',
    access_key: '',
    api_key: '',
    secret_key: ''
  },
  type: 'dnspod',
  name: ''
})
const updateDNSModel = ref<any>({
  data: {
    token: '',
    id: '',
    access_key: '',
    api_key: '',
    secret_key: ''
  },
  type: 'dnspod',
  name: ''
})
const addDNSModal = ref(false)
const updateDNSModal = ref(false)
const updateDNS = ref<any>()

const dnsProviders = ref<any>([])

const dnsColumns: any = [
  { title: '备注名称', key: 'name', resizable: true, ellipsis: { tooltip: true } },
  {
    title: '类型',
    key: 'type',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(
        NTag,
        {
          type: 'info',
          bordered: false
        },
        {
          default: () => {
            switch (row.type) {
              case 'dnspod':
                return 'DnsPod'
              case 'tencent':
                return '腾讯云'
              case 'aliyun':
                return '阿里云'
              case 'cloudflare':
                return 'Cloudflare'
              default:
                return '未知'
            }
          }
        }
      )
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            onClick: () => {
              updateDNS.value = row.id
              updateDNSModel.value.data.token = row.dns_param.token
              updateDNSModel.value.data.id = row.dns_param.id
              updateDNSModel.value.data.access_key = row.dns_param.access_key
              updateDNSModel.value.data.api_key = row.dns_param.api_key
              updateDNSModel.value.data.secret_key = row.dns_param.secret_key
              updateDNSModel.value.type = row.type
              updateDNSModel.value.name = row.name
              updateDNSModal.value = true
            }
          },
          {
            default: () => '修改'
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: async () => {
              await cert.dnsDelete(row.id)
              window.$message.success('删除成功')
              onDnsPageChange(1)
            }
          },
          {
            default: () => {
              return '确定删除 DNS 吗？'
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
const dnsData = ref<DNS[]>([] as DNS[])

const dnsPagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 10,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
})

const onDnsPageChange = (page: number) => {
  dnsPagination.page = page
  getDnsList(page, dnsPagination.pageSize).then((res) => {
    dnsData.value = res.items
    dnsPagination.itemCount = res.total
    dnsPagination.pageCount = res.total / dnsPagination.pageSize + 1
  })
}

const onDnsPageSizeChange = (pageSize: number) => {
  dnsPagination.pageSize = pageSize
  onDnsPageChange(1)
}

const getDnsList = async (page: number, limit: number) => {
  const { data } = await cert.dns(page, limit)
  return data
}

const handleAddDNS = async () => {
  await cert.dnsAdd(addDNSModel.value)
  window.$message.success('添加成功')
  addDNSModal.value = false
  onDnsPageChange(1)
  addDNSModel.value.data.token = ''
  addDNSModel.value.data.id = ''
  addDNSModel.value.data.access_key = ''
  addDNSModel.value.data.api_key = ''
  addDNSModel.value.data.secret_key = ''
  addDNSModel.value.name = ''
}

const handleUpdateDNS = async () => {
  await cert.dnsUpdate(updateDNS.value, updateDNSModel.value)
  window.$message.success('更新成功')
  updateDNSModal.value = false
  onDnsPageChange(1)
  updateDNSModel.value.data.token = ''
  updateDNSModel.value.data.id = ''
  updateDNSModel.value.data.access_key = ''
  updateDNSModel.value.data.api_key = ''
  updateDNSModel.value.data.secret_key = ''
  updateDNSModel.value.name = ''
}

onMounted(async () => {
  cert.dnsProviders().then((res) => {
    for (const item of res.data) {
      dnsProviders.value.push({
        label: item.name,
        value: item.dns
      })
    }
  })
  onDnsPageChange(1)
})
</script>

<template>
  <n-space vertical size="large">
    <n-card rounded-10>
      <n-space>
        <n-button type="primary" @click="addDNSModal = true"> 添加 DNS </n-button>
      </n-space>
    </n-card>
    <n-data-table
      striped
      remote
      :loading="false"
      :scroll-x="1200"
      :columns="dnsColumns"
      :data="dnsData"
      :row-key="(row: any) => row.id"
      :pagination="dnsPagination"
      @update:page="onDnsPageChange"
      @update:page-size="onDnsPageSizeChange"
    />
  </n-space>
  <n-modal
    v-model:show="addDNSModal"
    preset="card"
    title="添加 DNS"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="addDNSModel">
        <n-form-item path="name" label="备注名称">
          <n-input
            v-model:value="addDNSModel.name"
            type="text"
            @keydown.enter.prevent
            placeholder="输入备注名称"
          />
        </n-form-item>
        <n-form-item path="type" label="DNS">
          <n-select
            v-model:value="addDNSModel.type"
            placeholder="选择 DNS"
            clearable
            :options="dnsProviders"
          />
        </n-form-item>
        <n-form-item v-if="addDNSModel.type == 'dnspod'" path="id" label="ID">
          <n-input
            v-model:value="addDNSModel.data.id"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 DnsPod ID"
          />
        </n-form-item>
        <n-form-item v-if="addDNSModel.type == 'dnspod'" path="token" label="Token">
          <n-input
            v-model:value="addDNSModel.data.token"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 DnsPod Token"
          />
        </n-form-item>
        <n-form-item v-if="addDNSModel.type == 'tencent'" path="access_key" label="SecretId">
          <n-input
            v-model:value="addDNSModel.data.access_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入腾讯云 SecretId"
          />
        </n-form-item>
        <n-form-item v-if="addDNSModel.type == 'tencent'" path="secret_key" label="SecretKey">
          <n-input
            v-model:value="addDNSModel.data.secret_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入腾讯云 SecretKey"
          />
        </n-form-item>
        <n-form-item v-if="addDNSModel.type == 'aliyun'" path="access_key" label="Access Key">
          <n-input
            v-model:value="addDNSModel.data.access_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Access Key"
          />
        </n-form-item>
        <n-form-item v-if="addDNSModel.type == 'aliyun'" path="secret_key" label="Secret Key">
          <n-input
            v-model:value="addDNSModel.data.secret_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Secret Key"
          />
        </n-form-item>
        <n-form-item v-if="addDNSModel.type == 'cloudflare'" path="api_key" label="API Key">
          <n-input
            v-model:value="addDNSModel.data.api_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 Cloudflare API Key"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleAddDNS">提交</n-button>
    </n-space>
  </n-modal>
  <n-modal
    v-model:show="updateDNSModal"
    preset="card"
    title="修改 DNS"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-form :model="updateDNSModel">
        <n-form-item path="name" label="备注名称">
          <n-input
            v-model:value="updateDNSModel.name"
            type="text"
            @keydown.enter.prevent
            placeholder="输入备注名称"
          />
        </n-form-item>
        <n-form-item path="type" label="DNS">
          <n-select
            v-model:value="updateDNSModel.type"
            placeholder="选择 DNS"
            clearable
            :options="dnsProviders"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'dnspod'" path="id" label="ID">
          <n-input
            v-model:value="updateDNSModel.data.id"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 DnsPod ID"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'dnspod'" path="token" label="Token">
          <n-input
            v-model:value="updateDNSModel.data.token"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 DnsPod Token"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'aliyun'" path="access_key" label="Access Key">
          <n-input
            v-model:value="updateDNSModel.data.access_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Access Key"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'aliyun'" path="secret_key" label="Secret Key">
          <n-input
            v-model:value="updateDNSModel.data.secret_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Secret Key"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'cloudflare'" path="api_key" label="API Key">
          <n-input
            v-model:value="updateDNSModel.data.api_key"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 Cloudflare API Key"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleUpdateDNS">提交</n-button>
    </n-space>
  </n-modal>
</template>

<style scoped lang="scss"></style>
