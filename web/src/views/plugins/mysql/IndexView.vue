<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm } from 'naive-ui'
import mysql from '@/api/plugins/mysql'
import service from '@/api/panel/system/service'
import { generateRandomString, renderIcon } from '@/utils'
import type { Backup, Database, User } from '@/views/plugins/mysql/types'
import type { UploadFileInfo, MessageReactive } from 'naive-ui'
import Editor from '@guolao/vue-monaco-editor'

let messageReactive: MessageReactive | null = null

const currentTab = ref('status')
const currentDatabase = ref('')
const status = ref(false)
const isEnabled = ref(false)
const config = ref('')
const errorLog = ref('')
const slowLog = ref('')
const rootPassword = ref('')

const statusType = computed(() => {
  return status.value ? 'success' : 'error'
})
const statusStr = computed(() => {
  return status.value ? '正常运行中' : '已停止运行'
})

const addDatabaseModel = ref({
  database: '',
  user: '',
  password: generateRandomString(16)
})
const addUserModel = ref({
  database: '',
  user: '',
  password: generateRandomString(16)
})
const changePasswordModel = ref({
  user: '',
  password: generateRandomString(16)
})
const changePrivilegesModel = ref({
  user: '',
  database: ''
})

const addDatabaseModal = ref(false)
const addUserModal = ref(false)
const changePasswordModal = ref(false)
const changePrivilegesModal = ref(false)
const backupModal = ref(false)

const databaseColumns: any = [
  { title: '库名', key: 'name', fixed: 'left', resizable: true, ellipsis: { tooltip: true } },
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
            type: 'warning',
            secondary: true,
            onClick: () => {
              currentDatabase.value = row.name
              backupModal.value = true
            }
          },
          {
            default: () => '备份',
            icon: renderIcon('material-symbols:save-outline', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteDatabase(row.name)
          },
          {
            default: () => {
              return '确定删除数据库吗？'
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

const userColumns: any = [
  { title: '用户名', key: 'user', fixed: 'left', resizable: true, ellipsis: { tooltip: true } },
  { title: '主机', key: 'host', resizable: true, ellipsis: { tooltip: true } },
  { title: '权限', key: 'grants', width: 350, resizable: true, ellipsis: { tooltip: true } },
  {
    title: '操作',
    key: 'actions',
    width: 300,
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
            onClick: () => showChangePasswordModal(row.user)
          },
          {
            default: () => '改密',
            icon: renderIcon('majesticons:key-line', { size: 14 })
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'warning',
            secondary: true,
            style: 'margin-left: 15px;',
            onClick: () => showChangePrivilegesModal(row.user)
          },
          {
            default: () => '权限',
            icon: renderIcon('majesticons:lock-line', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteUser(row.user)
          },
          {
            default: () => {
              return '确定删除用户吗？'
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

const loadColumns: any = [
  { title: '属性', key: 'name', fixed: 'left', resizable: true, ellipsis: { tooltip: true } },
  { title: '当前值', key: 'value', width: 200, ellipsis: { tooltip: true } }
]

const backupColumns: any = [
  { title: '文件名', key: 'name', fixed: 'left', resizable: true, ellipsis: { tooltip: true } },
  { title: '大小', key: 'size', width: 200, ellipsis: { tooltip: true } },
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
            type: 'warning',
            secondary: true,
            onClick: () => handleRestoreBackup(row)
          },
          {
            default: () => '恢复',
            icon: renderIcon('material-symbols:settings-backup-restore-rounded', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteBackup(row.name)
          },
          {
            default: () => {
              return '确定删除备份吗？'
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

const databases = ref<Database[]>([] as Database[])
const users = ref<User[]>([] as User[])
const backup = ref<Backup[]>([])
const load = ref<any[]>([])

const databasePagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 10,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
})

const userPagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 10,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
})

const backupPagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 10,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100]
})

const getLoad = async () => {
  const { data } = await mysql.load()
  return data
}

const getDatabaseList = async (page: number, limit: number) => {
  const { data } = await mysql.databases(page, limit)
  return data
}

const getUserList = async (page: number, limit: number) => {
  const { data } = await mysql.users(page, limit)
  return data
}

const getBackupList = async (page: number, limit: number) => {
  const { data } = await mysql.backups(page, limit)
  return data
}

const onDatabasePageChange = (page: number) => {
  databasePagination.page = page
  getDatabaseList(page, databasePagination.pageSize).then((res) => {
    databases.value = res.items
    databasePagination.itemCount = res.total
    databasePagination.pageCount = res.total / databasePagination.pageSize + 1
  })
}

const onUserPageChange = (page: number) => {
  userPagination.page = page
  getUserList(page, userPagination.pageSize).then((res) => {
    users.value = res.items
    userPagination.itemCount = res.total
    userPagination.pageCount = res.total / userPagination.pageSize + 1
  })
}

const onBackupPageChange = (page: number) => {
  backupPagination.page = page
  getBackupList(page, backupPagination.pageSize).then((res) => {
    backup.value = res.items
    backupPagination.itemCount = res.total
    backupPagination.pageCount = res.total / backupPagination.pageSize + 1
  })
}

const onDatabasePageSizeChange = (pageSize: number) => {
  databasePagination.pageSize = pageSize
  onDatabasePageChange(1)
}

const onUserPageSizeChange = (pageSize: number) => {
  userPagination.pageSize = pageSize
  onUserPageChange(1)
}

const onBackupPageSizeChange = (pageSize: number) => {
  backupPagination.pageSize = pageSize
  onBackupPageChange(1)
}

const handleDeleteDatabase = async (name: string) => {
  mysql.deleteDatabase(name).then(() => {
    window.$message.success('删除成功')
    onDatabasePageChange(databasePagination.page)
  })
}

const handleDeleteUser = async (user: string) => {
  mysql.deleteUser(user).then(() => {
    window.$message.success('删除成功')
    onUserPageChange(userPagination.page)
  })
}

const showChangePasswordModal = (user: string) => {
  changePasswordModel.value.user = user
  changePasswordModal.value = true
}

const showChangePrivilegesModal = (user: string) => {
  changePrivilegesModel.value.user = user
  changePrivilegesModal.value = true
}

const getIsEnabled = async () => {
  await service.isEnabled('mysqld').then((res: any) => {
    isEnabled.value = res.data
  })
}

const getStatus = async () => {
  await service.status('mysqld').then((res: any) => {
    status.value = res.data
  })
}

const getRootPassword = async () => {
  await mysql.rootPassword().then((res: any) => {
    rootPassword.value = res.data
  })
}

const getErrorLog = async () => {
  const { data } = await mysql.errorLog()
  return data
}

const getSlowLog = async () => {
  const { data } = await mysql.slowLog()
  return data
}

const getConfig = async () => {
  const { data } = await mysql.config()
  return data
}

const handleSaveConfig = async () => {
  await mysql.saveConfig(config.value)
  window.$message.success('保存成功')
}

const handleClearErrorLog = async () => {
  await mysql.clearErrorLog()
  getErrorLog().then((res) => {
    errorLog.value = res
  })
  window.$message.success('清空成功')
}

const handleClearSlowLog = async () => {
  await mysql.clearSlowLog()
  getSlowLog().then((res) => {
    slowLog.value = res
  })
  window.$message.success('清空成功')
}

const handleIsEnabled = async () => {
  if (isEnabled.value) {
    await service.enable('mysqld')
    window.$message.success('开启自启动成功')
  } else {
    await service.disable('mysqld')
    window.$message.success('禁用自启动成功')
  }
  await getIsEnabled()
}

const handleStart = async () => {
  await service.start('mysqld')
  window.$message.success('启动成功')
  await getStatus()
}

const handleStop = async () => {
  await service.stop('mysqld')
  window.$message.success('停止成功')
  await getStatus()
}

const handleRestart = async () => {
  await service.restart('mysqld')
  window.$message.success('重启成功')
  await getStatus()
}

const handleReload = async () => {
  await service.reload('mysqld')
  window.$message.success('重载成功')
  await getStatus()
}

const handleSetRootPassword = async () => {
  await mysql.setRootPassword(rootPassword.value)
  window.$message.success('修改成功')
}

const handleAddDatabase = async () => {
  mysql.addDatabase(addDatabaseModel.value).then(() => {
    window.$message.success('添加成功')
    addDatabaseModal.value = false
    addDatabaseModel.value = {
      database: '',
      user: '',
      password: generateRandomString(16)
    }
    onDatabasePageChange(databasePagination.page)
    onUserPageChange(userPagination.page)
  })
}

const handleAddUser = async () => {
  mysql.addUser(addUserModel.value).then(() => {
    window.$message.success('添加成功')
    addUserModal.value = false
    addDatabaseModel.value = {
      user: '',
      password: generateRandomString(16),
      database: ''
    }
    onUserPageChange(userPagination.page)
  })
}

const handleChangePassword = async () => {
  mysql
    .setUserPassword(changePasswordModel.value.user, changePasswordModel.value.password)
    .then(() => {
      window.$message.success('修改成功')
      changePasswordModal.value = false
      changePasswordModel.value = {
        user: '',
        password: generateRandomString(16)
      }
      onUserPageChange(userPagination.page)
    })
}

const handleChangePrivileges = async () => {
  mysql
    .setUserPrivileges(changePrivilegesModel.value.user, changePrivilegesModel.value.database)
    .then(() => {
      window.$message.success('修改成功')
      changePrivilegesModal.value = false
      changePrivilegesModel.value = {
        user: '',
        database: ''
      }
      onUserPageChange(userPagination.page)
    })
}

const handleUploadBackup = async (files: UploadFileInfo[]) => {
  messageReactive = window.$message.loading('上传中...', {
    duration: 0
  })
  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    const formData = new FormData()
    formData.append('file', file.file as Blob, file.name)
    await mysql.uploadBackup(formData).then(() => {
      messageReactive?.destroy()
      window.$message.success('上传成功')
      onBackupPageChange(backupPagination.page)
    })
  }
}

const handleCreateBackup = async () => {
  messageReactive = window.$message.loading('创建中...', {
    duration: 0
  })
  await mysql.createBackup(currentDatabase.value).then(() => {
    messageReactive?.destroy()
    window.$message.success('创建成功')
    onBackupPageChange(backupPagination.page)
  })
}

const handleRestoreBackup = async (row: any) => {
  messageReactive = window.$message.loading('恢复中...', {
    duration: 0
  })
  await mysql.restoreBackup(row.name, currentDatabase.value).then(() => {
    messageReactive?.destroy()
    window.$message.success('恢复成功')
    onBackupPageChange(backupPagination.page)
  })
}

const handleDeleteBackup = async (name: string) => {
  await mysql.deleteBackup(name).then(() => {
    window.$message.success('删除成功')
    onBackupPageChange(backupPagination.page)
  })
}

onMounted(() => {
  getStatus()
  getIsEnabled()
  getRootPassword()
  onDatabasePageChange(databasePagination.page)
  onUserPageChange(userPagination.page)
  onBackupPageChange(backupPagination.page)
  getLoad().then((res) => {
    load.value = res
  })
  getErrorLog().then((res) => {
    errorLog.value = res
  })
  getSlowLog().then((res) => {
    slowLog.value = res
  })
  getConfig().then((res) => {
    config.value = res
  })
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-space v-if="currentTab == 'manage'">
        <n-button class="ml-16" type="info" @click="addUserModal = true">
          <TheIcon :size="18" class="mr-5" icon="material-symbols:add" />
          新建用户
        </n-button>
        <n-button class="ml-16" type="primary" @click="addDatabaseModal = true">
          <TheIcon :size="18" class="mr-5" icon="material-symbols:add" />
          新建数据库
        </n-button>
      </n-space>
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
        v-if="currentTab == 'error-log'"
        class="ml-16"
        type="primary"
        @click="handleClearErrorLog"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:delete-outline" />
        清空日志
      </n-button>
      <n-button
        v-if="currentTab == 'slow-log'"
        class="ml-16"
        type="primary"
        @click="handleClearSlowLog"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:delete-outline" />
        清空慢日志
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
                  停止 MySQL 会导致使用 MySQL 的网站无法访问，确定要停止吗？
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
          <n-card title="Root 密码" rounded-10>
            <n-space vertical>
              <n-input v-model:value="rootPassword"></n-input>
              <n-button type="primary" @click="handleSetRootPassword">保存修改</n-button>
            </n-space>
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="manage" tab="管理">
        <n-space vertical>
          <n-card title="数据库" :segmented="true" rounded-10>
            <n-data-table
              striped
              remote
              :loading="false"
              :columns="databaseColumns"
              :data="databases"
              :row-key="(row: any) => row.name"
              @update:page="onDatabasePageChange"
              @update:page-size="onDatabasePageSizeChange"
            />
          </n-card>
          <n-card title="用户" :segmented="true" rounded-10>
            <n-data-table
              striped
              remote
              :loading="false"
              :columns="userColumns"
              :data="users"
              :row-key="(row: any) => row.user"
              @update:page="onUserPageChange"
              @update:page-size="onUserPageSizeChange"
            />
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="config" tab="修改配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 MySQL 主配置文件，如果你不了解各参数的含义，请不要随意修改！
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
      <n-tab-pane name="load" tab="负载状态">
        <n-data-table striped remote :loading="false" :columns="loadColumns" :data="load" />
      </n-tab-pane>
      <n-tab-pane name="error-log" tab="错误日志">
        <Editor
          v-model:value="errorLog"
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
      <n-tab-pane name="slow-log" tab="慢查询日志">
        <Editor
          v-model:value="slowLog"
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
  <n-modal v-model:show="addDatabaseModal" title="新建数据库">
    <n-card
      closable
      @close="() => (addDatabaseModal = false)"
      title="新建数据库"
      style="width: 60vw"
    >
      <n-form :model="addDatabaseModel">
        <n-form-item path="database" label="数据库名">
          <n-input
            v-model:value="addDatabaseModel.database"
            type="text"
            @keydown.enter.prevent
            placeholder="输入数据库名"
          />
        </n-form-item>
        <n-form-item path="user" label="用户名">
          <n-input
            v-model:value="addDatabaseModel.user"
            type="text"
            @keydown.enter.prevent
            placeholder="输入用户名"
          />
        </n-form-item>
        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="addDatabaseModel.password"
            type="text"
            @keydown.enter.prevent
            placeholder="建议使用生成器生成随机密码"
          />
        </n-form-item>
      </n-form>
      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <n-button type="info" block @click="handleAddDatabase">提交</n-button>
        </n-col>
      </n-row>
    </n-card>
  </n-modal>
  <n-modal v-model:show="addUserModal" title="新建用户">
    <n-card closable @close="() => (addUserModal = false)" title="新建用户" style="width: 60vw">
      <n-form :model="addUserModel">
        <n-form-item path="user" label="用户名">
          <n-input
            v-model:value="addUserModel.user"
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
        <n-form-item path="database" label="数据库名">
          <n-input
            v-model:value="addUserModel.database"
            type="text"
            @keydown.enter.prevent
            placeholder="输入授权给该用户的数据库名"
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
  <n-modal v-model:show="backupModal">
    <n-card
      closable
      @close="() => (backupModal = false)"
      :title="'备份管理 - ' + currentDatabase"
      style="width: 60vw"
    >
      <n-space vertical>
        <n-space>
          <n-button type="primary" @click="handleCreateBackup">创建备份</n-button>
          <n-upload
            accept=".sql,.zip,tar.gz,.tar,.rar,.bz2"
            :default-upload="false"
            :show-file-list="false"
            @update:file-list="handleUploadBackup"
          >
            <n-button>上传备份</n-button>
          </n-upload>
        </n-space>
        <n-data-table
          striped
          remote
          :loading="false"
          :columns="backupColumns"
          :data="backup"
          :row-key="(row: any) => row.name"
          @update:page="onBackupPageChange"
          @update:page-size="onBackupPageSizeChange"
        />
      </n-space>
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
  <n-modal v-model:show="changePrivilegesModal">
    <n-card
      closable
      @close="() => (changePrivilegesModal = false)"
      title="修改权限"
      style="width: 60vw"
    >
      <n-form :model="changePrivilegesModel">
        <n-form-item path="database" label="数据库名">
          <n-input
            v-model:value="changePrivilegesModel.database"
            type="text"
            @keydown.enter.prevent
            placeholder="输入授权给该用户的数据库名"
          />
        </n-form-item>
      </n-form>
      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <n-button type="info" block @click="handleChangePrivileges">提交</n-button>
        </n-col>
      </n-row>
    </n-card>
  </n-modal>
</template>
