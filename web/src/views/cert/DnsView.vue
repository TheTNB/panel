<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm, NSpace, NTag } from 'naive-ui'

import cert from '@/api/panel/cert'
import type { DNS } from '@/views/cert/types'

const props = defineProps({
  dnsProviders: Array<any>
})

const { dnsProviders } = toRefs(props)

const updateDNSModel = ref<any>({
  data: {
    ak: '',
    sk: ''
  },
  type: 'aliyun',
  name: ''
})
const updateDNSModal = ref(false)
const updateDNS = ref<any>()

const dnsColumns: any = [
  {
    title: '备注名称',
    key: 'name',
    minWidth: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
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
              case 'aliyun':
                return '阿里云'
              case 'tencent':
                return '腾讯云'
              case 'huawei':
                return '华为云'
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
              updateDNSModel.value.data.ak = row.dns_param.ak
              updateDNSModel.value.data.sk = row.dns_param.sk
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
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200]
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

const handleUpdateDNS = async () => {
  await cert.dnsUpdate(updateDNS.value, updateDNSModel.value)
  window.$message.success('更新成功')
  updateDNSModal.value = false
  onDnsPageChange(1)
  updateDNSModel.value.data.ak = ''
  updateDNSModel.value.data.sk = ''
  updateDNSModel.value.name = ''
}

onMounted(async () => {
  onDnsPageChange(1)
})
</script>

<template>
  <n-space vertical size="large">
    <n-data-table
      striped
      remote
      :scroll-x="1000"
      :loading="false"
      :columns="dnsColumns"
      :data="dnsData"
      :row-key="(row: any) => row.id"
      :pagination="dnsPagination"
      @update:page="onDnsPageChange"
      @update:page-size="onDnsPageSizeChange"
    />
  </n-space>
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

        <n-form-item v-if="updateDNSModel.type == 'aliyun'" path="ak" label="Access Key">
          <n-input
            v-model:value="updateDNSModel.data.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Access Key"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'aliyun'" path="sk" label="Secret Key">
          <n-input
            v-model:value="updateDNSModel.data.sk"
            type="text"
            @keydown.enter.prevent
            placeholder="输入阿里云 Secret Key"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'tencent'" path="ak" label="SecretId">
          <n-input
            v-model:value="updateDNSModel.data.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入腾讯云 SecretId"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'tencent'" path="sk" label="SecretKey">
          <n-input
            v-model:value="updateDNSModel.data.sk"
            type="text"
            @keydown.enter.prevent
            placeholder="输入腾讯云 SecretKey"
          />
        </n-form-item>

        <n-form-item v-if="updateDNSModel.type == 'huawei'" path="ak" label="AccessKeyId">
          <n-input
            v-model:value="updateDNSModel.data.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入华为云 AccessKeyId"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'huawei'" path="sk" label="SecretAccessKey">
          <n-input
            v-model:value="updateDNSModel.data.sk"
            type="text"
            @keydown.enter.prevent
            placeholder="输入华为云 SecretAccessKey"
          />
        </n-form-item>
        <n-form-item v-if="updateDNSModel.type == 'cloudflare'" path="ak" label="API Key">
          <n-input
            v-model:value="updateDNSModel.data.ak"
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
