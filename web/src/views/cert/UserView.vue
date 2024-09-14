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
import type { User } from '@/views/cert/types'

let messageReactive: MessageReactive | null = null
const addUserModel = ref<any>({
  hmac_encoded: '',
  email: '',
  kid: '',
  key_type: 'P256',
  ca: 'letsencrypt'
})
const updateUserModel = ref<any>({
  hmac_encoded: '',
  email: '',
  kid: '',
  key_type: 'P256',
  ca: 'letsencrypt'
})
const addUserModal = ref(false)
const updateUserModal = ref(false)
const updateUser = ref<any>()

const caProviders = ref<any>([])
const algorithms = ref<any>([])

const userColumns: any = [
  { title: '邮箱', key: 'email', resizable: true, ellipsis: { tooltip: true } },
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
            switch (row.ca) {
              case 'letsencrypt':
                return "Let's Encrypt"
              case 'zerossl':
                return 'ZeroSSL'
              case 'sslcom':
                return 'SSL.com'
              case 'buypass':
                return 'Buypass'
              case 'google':
                return 'Google'
              default:
                return '未知'
            }
          }
        }
      )
    }
  },
  { title: '密钥类型', key: 'key_type', width: 150, resizable: true, ellipsis: { tooltip: true } },
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
              updateUser.value = row.id
              updateUserModel.value.email = row.email
              updateUserModel.value.hmac_encoded = row.hmac_encoded
              updateUserModel.value.kid = row.kid
              updateUserModel.value.key_type = row.key_type
              updateUserModel.value.ca = row.ca
              updateUserModal.value = true
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
              await cert.userDelete(row.id)
              window.$message.success('删除成功')
              onUserPageChange(1)
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
const userData = ref<User[]>([] as User[])

const userPagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 10,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
})

const onUserPageChange = (page: number) => {
  userPagination.page = page
  getUserList(page, userPagination.pageSize).then((res) => {
    userData.value = res.items
    userPagination.itemCount = res.total
    userPagination.pageCount = res.total / userPagination.pageSize + 1
  })
}

const onUserPageSizeChange = (pageSize: number) => {
  userPagination.pageSize = pageSize
  onUserPageChange(1)
}

const getUserList = async (page: number, limit: number) => {
  const { data } = await cert.users(page, limit)
  return data
}

const handleAddUser = async () => {
  messageReactive = window.$message.loading('正在向 CA 注册账号，请耐心等待', {
    duration: 0
  })
  await cert.userAdd(addUserModel.value)
  messageReactive.destroy()
  window.$message.success('添加成功')
  addUserModal.value = false
  onUserPageChange(1)
  addUserModel.value.email = ''
  addUserModel.value.hmac_encoded = ''
  addUserModel.value.kid = ''
}

const handleUpdateUser = async () => {
  messageReactive = window.$message.loading('正在向 CA 注册账号，请耐心等待', {
    duration: 0
  })
  await cert.userUpdate(updateUser.value, updateUserModel.value)
  messageReactive.destroy()
  window.$message.success('更新成功')
  updateUserModal.value = false
  onUserPageChange(1)
  updateUserModel.value.email = ''
  updateUserModel.value.hmac_encoded = ''
  updateUserModel.value.kid = ''
}

onMounted(() => {
  cert.caProviders().then((res) => {
    for (const item of res.data) {
      caProviders.value.push({
        label: item.name,
        value: item.ca
      })
    }
  })
  cert.algorithms().then((res) => {
    for (const item of res.data) {
      algorithms.value.push({
        label: item.name,
        value: item.key
      })
    }
  })
  onUserPageChange(1)
})
</script>

<template>
  <n-space vertical size="large">
    <n-card rounded-10>
      <n-space>
        <n-button type="primary" @click="addUserModal = true"> 添加账号 </n-button>
      </n-space>
    </n-card>
    <n-data-table
      striped
      remote
      :loading="false"
      :scroll-x="1200"
      :columns="userColumns"
      :data="userData"
      :row-key="(row: any) => row.id"
      :pagination="userPagination"
      @update:page="onUserPageChange"
      @update:page-size="onUserPageSizeChange"
    />
  </n-space>
  <n-modal
    v-model:show="addUserModal"
    preset="card"
    title="添加账号"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-space vertical>
      <n-alert type="info"> Google 和 SSL.com 需要先去官网获得 KID 和 HMAC 并填入 </n-alert>
      <n-alert type="warning">
        境内无法使用 Google CA，其他 CA 视网络情况而定，建议使用 Let's Encrypt
      </n-alert>
      <n-form :model="addUserModel">
        <n-form-item path="ca" label="CA">
          <n-select
            v-model:value="addUserModel.ca"
            placeholder="选择 CA"
            clearable
            :options="caProviders"
          />
        </n-form-item>
        <n-form-item path="key_type" label="密钥类型">
          <n-select
            v-model:value="addUserModel.key_type"
            placeholder="选择密钥类型"
            clearable
            :options="algorithms"
          />
        </n-form-item>
        <n-form-item path="email" label="邮箱">
          <n-input
            v-model:value="addUserModel.email"
            type="text"
            @keydown.enter.prevent
            placeholder="输入邮箱地址"
          />
        </n-form-item>
        <n-form-item path="kid" label="KID">
          <n-input
            v-model:value="addUserModel.kid"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 KID"
          />
        </n-form-item>
        <n-form-item path="hmac_encoded" label="HMAC">
          <n-input
            v-model:value="addUserModel.hmac_encoded"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 HMAC"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleAddUser">提交</n-button>
    </n-space>
  </n-modal>
  <n-modal
    v-model:show="updateUserModal"
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
        境内无法使用 Google CA，其他 CA 视网络情况而定，建议使用 Let's Encrypt
      </n-alert>
      <n-form :model="updateUserModel">
        <n-form-item path="ca" label="CA">
          <n-select
            v-model:value="updateUserModel.ca"
            placeholder="选择 CA"
            clearable
            :options="caProviders"
          />
        </n-form-item>
        <n-form-item path="key_type" label="密钥类型">
          <n-select
            v-model:value="updateUserModel.key_type"
            placeholder="选择密钥类型"
            clearable
            :options="algorithms"
          />
        </n-form-item>
        <n-form-item path="email" label="邮箱">
          <n-input
            v-model:value="updateUserModel.email"
            type="text"
            @keydown.enter.prevent
            placeholder="输入邮箱地址"
          />
        </n-form-item>
        <n-form-item path="kid" label="KID">
          <n-input
            v-model:value="updateUserModel.kid"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 KID"
          />
        </n-form-item>
        <n-form-item path="hmac_encoded" label="HMAC">
          <n-input
            v-model:value="updateUserModel.hmac_encoded"
            type="text"
            @keydown.enter.prevent
            placeholder="输入 HMAC"
          />
        </n-form-item>
      </n-form>
      <n-button type="info" block @click="handleUpdateUser">提交</n-button>
    </n-space>
  </n-modal>
</template>

<style scoped lang="scss"></style>
