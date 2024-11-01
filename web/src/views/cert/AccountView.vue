<script setup lang="ts">
import {
  type MessageReactive,
  NButton,
  NDataTable,
  NInput,
  NPopconfirm,
  NSpace,
  NTag
} from 'naive-ui'

import cert from '@/api/panel/cert'
import type { Account } from '@/views/cert/types'

const props = defineProps({
  caProviders: {
    type: Array<any>,
    required: true
  },
  algorithms: {
    type: Array<any>,
    required: true
  }
})

const { caProviders, algorithms } = toRefs(props)

let messageReactive: MessageReactive | null = null

const updateAccountModel = ref<any>({
  hmac_encoded: '',
  email: '',
  kid: '',
  key_type: 'P256',
  ca: 'googlecn'
})
const updateAccountModal = ref(false)
const updateAccount = ref<any>()

const columns: any = [
  {
    title: '邮箱',
    key: 'email',
    minWidth: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: 'CA',
    key: 'ca',
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
            return caProviders.value?.find((item: any) => item.value === row.ca)?.label
          }
        }
      )
    }
  },
  {
    title: '密钥类型',
    key: 'key_type',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true }
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
              updateAccount.value = row.id
              updateAccountModel.value.email = row.email
              updateAccountModel.value.hmac_encoded = row.hmac_encoded
              updateAccountModel.value.kid = row.kid
              updateAccountModel.value.key_type = row.key_type
              updateAccountModel.value.ca = row.ca
              updateAccountModal.value = true
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
              await cert.accountDelete(row.id)
              window.$message.success('删除成功')
              onPageChange(1)
            }
          },
          {
            default: () => {
              return '确定删除账号吗？'
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
const data = ref<Account[]>([] as Account[])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200]
})

const onPageChange = (page: number) => {
  pagination.page = page
  getAccountList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const getAccountList = async (page: number, limit: number) => {
  const { data } = await cert.accounts(page, limit)
  return data
}

const handleUpdateAccount = async () => {
  messageReactive = window.$message.loading('正在向 CA 注册账号，请耐心等待', {
    duration: 0
  })
  await cert.accountUpdate(updateAccount.value, updateAccountModel.value)
  messageReactive.destroy()
  window.$message.success('更新成功')
  updateAccountModal.value = false
  onPageChange(1)
  updateAccountModel.value.email = ''
  updateAccountModel.value.hmac_encoded = ''
  updateAccountModel.value.kid = ''
}

onMounted(() => {
  onPageChange(pagination.page)
  window.$bus.on('cert:refresh-account', () => {
    onPageChange(pagination.page)
  })
})

onUnmounted(() => {
  window.$bus.off('cert:refresh-account')
})
</script>

<template>
  <n-space vertical size="large">
    <n-data-table
      striped
      remote
      :scroll-x="1000"
      :loading="false"
      :columns="columns"
      :data="data"
      :row-key="(row: any) => row.id"
      :pagination="pagination"
      @update:page="onPageChange"
      @update:page-size="onPageSizeChange"
    />
  </n-space>
  <n-modal
    v-model:show="updateAccountModal"
    preset="card"
    title="修改账号"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-alert type="info"> Google 和 SSL.com 需要先去官网获得 KID 和 HMAC 并填入 </n-alert>
      <n-alert type="warning">
        境内无法使用 Google，其他 CA 视网络情况而定，建议使用 GoogleCN 或 Let's Encrypt
      </n-alert>
      <n-form :model="updateAccountModel">
        <n-form-item path="ca" label="CA">
          <n-select
            v-model:value="updateAccountModel.ca"
            placeholder="选择 CA"
            clearable
            :options="caProviders"
          />
        </n-form-item>
        <n-form-item path="key_type" label="密钥类型">
          <n-select
            v-model:value="updateAccountModel.key_type"
            placeholder="选择密钥类型"
            clearable
            :options="algorithms"
          />
        </n-form-item>
        <n-form-item path="email" label="邮箱">
          <n-input
            v-model:value="updateAccountModel.email"
            type="text"
            @keydown.enter.prevent
            placeholder="输入邮箱地址"
          />
        </n-form-item>
        <n-form-item path="kid" label="KID">
          <n-input
            v-model:value="updateAccountModel.kid"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 KID"
          />
        </n-form-item>
        <n-form-item path="hmac_encoded" label="HMAC">
          <n-input
            v-model:value="updateAccountModel.hmac_encoded"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 HMAC"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleUpdateAccount">提交</n-button>
    </n-space>
  </n-modal>
</template>

<style scoped lang="scss"></style>
