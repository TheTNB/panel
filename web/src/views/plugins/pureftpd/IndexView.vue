<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm } from 'naive-ui'
import { generateRandomString, renderIcon } from '@/utils'
import type { User } from '@/views/plugins/pureftpd/types'
import pureftpd from '@/api/plugins/pureftpd'
import service from '@/api/panel/system/service'

const currentTab = ref('status')
const status = ref(false)
const isEnabled = ref(false)
const port = ref(0)
const addUserModal = ref(false)
const changePasswordModal = ref(false)

const statusType = computed(() => {
  return status.value ? 'success' : 'error'
})
const statusStr = computed(() => {
  return status.value ? '正常运行中' : '已停止运行'
})

const addUserModel = ref({
  username: '',
  password: generateRandomString(16),
  path: ''
})

const changePasswordModel = ref({
  username: '',
  password: generateRandomString(16)
})

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
})

const userColumns: any = [
  {
    title: '用户名',
    key: 'username',
    fixed: 'left',
    width: 250,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  { title: '路径', key: 'path', resizable: true, ellipsis: { tooltip: true } },
  {
    title: '操作',
    key: 'actions',
    width: 240,
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
            secondary: true,
            onClick: () => {
              changePasswordModel.value.username = row.username
              changePasswordModel.value.password = generateRandomString(16)
              changePasswordModal.value = true
            }
          },
          {
            default: () => '改密',
            icon: renderIcon('material-symbols:key-outline', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteUser(row.username)
          },
          {
            default: () => {
              return '确定删除用户' + row.username + '吗？'
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

const users = ref<User[]>([] as User[])

const getStatus = async () => {
  await service.status('pure-ftpd').then((res: any) => {
    status.value = res.data
  })
}

const getIsEnabled = async () => {
  await service.isEnabled('pure-ftpd').then((res: any) => {
    isEnabled.value = res.data
  })
}

const getPort = async () => {
  await pureftpd.port().then((res: any) => {
    port.value = res.data
  })
}

const handleSavePort = async () => {
  await pureftpd.setPort(port.value)
  window.$message.success('保存成功')
}

const handleStart = async () => {
  await service.start('pure-ftpd')
  window.$message.success('启动成功')
  await getStatus()
}

const handleIsEnabled = async () => {
  if (isEnabled.value) {
    await service.enable('pure-ftpd')
    window.$message.success('开启自启动成功')
  } else {
    await service.disable('pure-ftpd')
    window.$message.success('禁用自启动成功')
  }
  await getIsEnabled()
}

const handleStop = async () => {
  await service.stop('pure-ftpd')
  window.$message.success('停止成功')
  await getStatus()
}

const handleRestart = async () => {
  await service.restart('pure-ftpd')
  window.$message.success('重启成功')
  await getStatus()
}

const getUsers = async (page: number, limit: number) => {
  const { data } = await pureftpd.list(page, limit)
  return data
}

const onPageChange = (page: number) => {
  pagination.page = page
  getUsers(page, pagination.pageSize).then((res) => {
    users.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const handleAddUser = async () => {
  await pureftpd.add(
    addUserModel.value.username,
    addUserModel.value.password,
    addUserModel.value.path
  )
  window.$message.success('添加成功')
  onPageChange(1)
  addUserModal.value = false
  addUserModel.value.username = ''
  addUserModel.value.password = generateRandomString(16)
  addUserModel.value.path = ''
}

const handleChangePassword = async () => {
  await pureftpd.changePassword(
    changePasswordModel.value.username,
    changePasswordModel.value.password
  )
  window.$message.success('修改成功')
  onPageChange(1)
  changePasswordModal.value = false
}

const handleDeleteUser = async (username: string) => {
  await pureftpd.delete(username)
  window.$message.success('删除成功')
  onPageChange(1)
}

onMounted(() => {
  getStatus()
  getIsEnabled()
  getPort()
  onPageChange(1)
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button v-if="currentTab == 'status'" class="ml-16" type="primary" @click="handleSavePort">
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'users'"
        class="ml-16"
        type="primary"
        @click="addUserModal = true"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:add" />
        添加用户
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="status" tab="运行状态">
        <n-space vertical>
          <n-card title="运行状态" rounded-10>
            <template #header-extra>
              <n-switch v-model:value="isEnabled" @update:value="handleIsEnabled">
                <template #checked> 自启动开 </template>
                <template #unchecked> 自启动关 </template>
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
                  停止 Pure-Ftpd 会导致无法使用 FTP 服务，确定要停止吗？
                </n-popconfirm>
                <n-button type="warning" @click="handleRestart">
                  <TheIcon :size="18" class="mr-5" icon="material-symbols:replay-rounded" />
                  重启
                </n-button>
              </n-space>
            </n-space>
          </n-card>
          <n-card title="端口设置" rounded-10>
            <n-input-number v-model:value="port" min="1" />
            修改 Pure-Ftpd 监听端口
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="users" tab="用户管理">
        <n-card title="用户列表" :segmented="true" rounded-10>
          <n-data-table
            striped
            remote
            :loading="false"
            :columns="userColumns"
            :data="users"
            :row-key="(row: any) => row.username"
            @update:page="onPageChange"
            @update:page-size="onPageSizeChange"
          />
        </n-card>
      </n-tab-pane>
    </n-tabs>
  </common-page>
  <n-modal v-model:show="addUserModal" title="新建用户">
    <n-card closable @close="() => (addUserModal = false)" title="新建用户" style="width: 60vw">
      <n-form :model="addUserModel">
        <n-form-item path="username" label="用户名">
          <n-input
            v-model:value="addUserModel.username"
            type="text"
            @keydown.enter.prevent
            placeholder="输入用户名"
          />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="addUserModel.password"
            type="text"
            @keydown.enter.prevent
            placeholder="建议使用生成器生成随机密码"
          />
        </n-form-item>
        <n-form-item path="path" label="目录">
          <n-input
            v-model:value="addUserModel.path"
            type="text"
            @keydown.enter.prevent
            placeholder="输入授权给该用户的目录"
          />
        </n-form-item>
      </n-form>
      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <n-button type="info" block @click="handleAddUser">提交</n-button>
        </n-col>
      </n-row>
    </n-card>
  </n-modal>
  <n-modal v-model:show="changePasswordModal">
    <n-card
      closable
      @close="() => (changePasswordModal = false)"
      title="修改密码"
      style="width: 60vw"
    >
      <n-form :model="changePasswordModel">
        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="changePasswordModel.password"
            type="text"
            @keydown.enter.prevent
            placeholder="建议使用生成器生成随机密码"
          />
        </n-form-item>
      </n-form>
      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <n-button type="info" block @click="handleChangePassword">提交</n-button>
        </n-col>
      </n-row>
    </n-card>
  </n-modal>
</template>
