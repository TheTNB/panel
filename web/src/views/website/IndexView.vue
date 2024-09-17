<script lang="ts" setup>
import {
  NButton,
  NDataTable,
  NSpace,
  NSwitch,
  NPopconfirm,
  NInput,
  NFlex,
  NCheckbox
} from 'naive-ui'
import website from '@/api/panel/website'
import info from '@/api/panel/info'
import { generateRandomString, isNullOrUndef, renderIcon } from '@/utils'
import type { Backup, Website } from './types'
import type { UploadFileInfo, MessageReactive } from 'naive-ui'
import Editor from '@guolao/vue-monaco-editor'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
let messageReactive: MessageReactive | null = null
const selectedRowKeys = ref<any>([])

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: t('websiteIndex.columns.name'),
    key: 'name',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('websiteIndex.columns.status'),
    key: 'status',
    width: 60,
    align: 'center',
    render(row: any) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.status,
        onUpdateValue: () => handleStatusChange(row)
      })
    }
  },
  {
    title: t('websiteIndex.columns.path'),
    key: 'path',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: 'PHP',
    key: 'php',
    width: 60,
    ellipsis: { tooltip: true },
    render(row: any) {
      return row.php === 0 ? '不使用' : row.php
    }
  },
  {
    title: 'SSL',
    key: 'ssl',
    width: 60,
    align: 'center',
    render(row: any) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.ssl,
        onClick: () => handleEdit(row)
      })
    }
  },
  {
    title: t('websiteIndex.columns.remark'),
    key: 'remark',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true },
    render(row: any) {
      return h(NInput, {
        size: 'small',
        value: row.remark,
        onBlur: () => handleRemark(row),
        onUpdateValue(v) {
          row.remark = v
        }
      })
    }
  },
  {
    title: t('websiteIndex.columns.actions'),
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
              currentWebsite.value = row.id
              backupModal.value = true
            }
          },
          {
            default: () => '备份',
            icon: renderIcon('majesticons:eye-line', { size: 14 })
          }
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            style: 'margin-left: 15px;',
            onClick: () => handleEdit(row)
          },
          {
            default: () => '修改',
            icon: renderIcon('material-symbols:edit-outline', { size: 14 })
          }
        ),
        h(
          NPopconfirm,
          {
            showIcon: false,
            onPositiveClick: () => handleDelete(row.id)
          },
          {
            default: () => {
              return h(
                NFlex,
                {
                  vertical: true
                },
                {
                  default: () => [
                    h('strong', {}, { default: () => `确定删除网站 ${row.name} 吗？` }),
                    h(
                      NCheckbox,
                      {
                        checked: deleteModel.value.path,
                        onUpdateChecked: (v) => (deleteModel.value.path = v)
                      },
                      { default: () => '删除网站目录' }
                    ),
                    h(
                      NCheckbox,
                      {
                        checked: deleteModel.value.db,
                        onUpdateChecked: (v) => (deleteModel.value.db = v)
                      },
                      { default: () => '删除本地同名数据库' }
                    )
                  ]
                }
              )
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

const data = ref<Website[]>([] as Website[])
const backup = ref<Backup[]>([])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
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

const currentWebsite = ref(0)
const addModal = ref(false)
const editDefaultPageModal = ref(false)
const backupModal = ref(false)

const buttonLoading = ref(false)
const buttonDisabled = ref(false)
const addModel = ref({
  name: '',
  domains: [] as Array<string>,
  ports: [] as Array<string>,
  php: '0',
  db: false,
  db_type: '0',
  db_name: '',
  db_user: '',
  db_password: '',
  path: '',
  remark: ''
})
const deleteModel = ref({
  id: 0,
  path: true,
  db: false
})
const editDefaultPageModel = ref({
  index: '',
  stop: ''
})

const installedDbAndPhp = ref({
  php: [
    {
      label: '',
      value: ''
    }
  ],
  db: [
    {
      label: '',
      value: ''
    }
  ]
})

const getPhpAndDb = async () => {
  const { data } = await info.installedDbAndPhp()
  installedDbAndPhp.value = data
}

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

// 修改运行状态
const handleStatusChange = (row: any) => {
  if (isNullOrUndef(row.id)) return

  website.status(row.id, !row.status).then(() => {
    row.status = !row.status
    window.$message.success('已' + (row.status ? '启动' : '停止'))
  })
}

const getWebsiteList = async (page: number, limit: number) => {
  const { data } = await website.list(page, limit)
  return data
}

const getDefaultPage = async () => {
  const { data } = await website.defaultConfig()
  editDefaultPageModel.value = data
}
const getBackupList = async (page: number, limit: number) => {
  const { data } = await website.backupList(page, limit)
  return data
}

const onPageChange = (page: number) => {
  pagination.page = page
  getWebsiteList(page, pagination.pageSize).then((res) => {
    data.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
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

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}
const onBackupPageSizeChange = (pageSize: number) => {
  backupPagination.pageSize = pageSize
  onBackupPageChange(1)
}

const handleRemark = (row: any) => {
  website.updateRemark(row.id, row.remark).then(() => {
    window.$message.success('修改成功')
  })
}

const handleEdit = (row: any) => {
  router.push({
    name: 'website-edit',
    params: {
      id: row.id
    }
  })
}

const handleDelete = async (id: number) => {
  deleteModel.value.id = id
  await website.delete(deleteModel.value).then(() => {
    window.$message.success('删除成功')
    onPageChange(pagination.page)
  })
  deleteModel.value.id = 0
  deleteModel.value.path = true
}

const handleSaveDefaultPage = () => {
  website
    .saveDefaultConfig(editDefaultPageModel.value.index, editDefaultPageModel.value.stop)
    .then(() => {
      window.$message.success('修改成功')
      editDefaultPageModal.value = false
    })
}

const handleAdd = async () => {
  buttonLoading.value = true
  buttonDisabled.value = true
  // 去除空的域名和端口
  addModel.value.domains = addModel.value.domains.filter((item) => item !== '')
  addModel.value.ports = addModel.value.ports.filter((item) => item !== '')
  // 端口为空自动添加 80 端口
  if (addModel.value.ports.length === 0) {
    addModel.value.ports.push('80')
  }
  await website
    .add(addModel.value)
    .then(() => {
      window.$message.success('添加成功')
      getWebsiteList(pagination.page, pagination.pageSize).then((res) => {
        data.value = res.items
        pagination.itemCount = res.total
        pagination.pageCount = res.total / pagination.pageSize + 1
      })
      addModal.value = false
      addModel.value = {
        name: '',
        domains: [] as Array<string>,
        ports: [] as Array<string>,
        php: '0',
        db: false,
        db_type: '0',
        db_name: '',
        db_user: '',
        db_password: '',
        path: '',
        remark: ''
      }
      buttonLoading.value = false
      buttonDisabled.value = false
    })
    .catch(() => {
      buttonLoading.value = false
      buttonDisabled.value = false
    })
}

const batchDelete = async () => {
  if (selectedRowKeys.value.length === 0) {
    window.$message.info('请选择要删除的网站')
    return
  }

  for (const id of selectedRowKeys.value) {
    deleteModel.value.id = id
    deleteModel.value.path = true
    deleteModel.value.db = false
    await website.delete(deleteModel.value).then(() => {
      let site = data.value.find((item) => item.id === id)
      window.$message.success('网站 ' + site?.name + ' 删除成功')
    })
    deleteModel.value.id = 0
  }

  onPageChange(pagination.page)
}

const handleUploadBackup = async (files: UploadFileInfo[]) => {
  messageReactive = window.$message.loading('上传中...', {
    duration: 0
  })
  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    const formData = new FormData()
    formData.append('file', file.file as Blob, file.name)
    await website.uploadBackup(formData).then(() => {
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
  await website.createBackup(currentWebsite.value).then(() => {
    messageReactive?.destroy()
    window.$message.success('创建成功')
    onBackupPageChange(backupPagination.page)
  })
}

const handleRestoreBackup = async (row: any) => {
  messageReactive = window.$message.loading('恢复中...', {
    duration: 0
  })
  await website.restoreBackup(currentWebsite.value, row.name).then(() => {
    messageReactive?.destroy()
    window.$message.success('恢复成功')
    onBackupPageChange(backupPagination.page)
  })
}

const handleDeleteBackup = async (name: string) => {
  await website.deleteBackup(name).then(() => {
    window.$message.success('删除成功')
    onBackupPageChange(backupPagination.page)
  })
}

const formatDbValue = (value: string) => {
  value = value.replace(/\./g, '_')
  value = value.replace(/-/g, '_')
  if (value.length > 16) {
    value = value.substring(0, 16)
  }

  return value
}

onMounted(() => {
  onPageChange(pagination.page)
  getPhpAndDb()
  getDefaultPage()
  onBackupPageChange(backupPagination.page)
})
</script>

<template>
  <common-page show-footer>
    <n-space vertical size="large">
      <n-card rounded-10>
        <n-space>
          <n-button type="primary" @click="addModal = true">
            {{ $t('websiteIndex.create.trigger') }}
          </n-button>
          <n-popconfirm @positive-click="batchDelete">
            <template #trigger>
              <n-button type="error"> 批量删除 </n-button>
            </template>
            这会删除网站目录但不会删除同名数据库，确定删除选中的网站吗？
          </n-popconfirm>
          <n-button type="warning" @click="editDefaultPageModal = true">
            {{ $t('websiteIndex.edit.trigger') }}
          </n-button>
        </n-space>
      </n-card>
      <n-data-table
        striped
        remote
        :loading="false"
        :scroll-x="1200"
        :columns="columns"
        :data="data"
        :row-key="(row: any) => row.id"
        :pagination="pagination"
        @update:checked-row-keys="onChecked"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </n-space>
  </common-page>
  <n-modal
    v-model:show="addModal"
    :title="$t('websiteIndex.create.title')"
    preset="card"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="addModal = false"
  >
    <n-form :model="addModel">
      <n-form-item path="name" :label="$t('websiteIndex.create.fields.name.label')">
        <n-input
          v-model:value="addModel.name"
          type="text"
          @keydown.enter.prevent
          :placeholder="$t('websiteIndex.create.fields.name.placeholder')"
        />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="11">
          <n-form-item :label="$t('websiteIndex.create.fields.domains.label')">
            <n-dynamic-input
              v-model:value="addModel.domains"
              placeholder="example.com"
              :min="1"
              show-sort-button
            />
          </n-form-item>
        </n-col>
        <n-col :span="2"></n-col>
        <n-col :span="11">
          <n-form-item :label="$t('websiteIndex.create.fields.port.label')">
            <n-dynamic-input
              v-model:value="addModel.ports"
              placeholder="80"
              :min="1"
              show-sort-button
            />
          </n-form-item>
        </n-col>
      </n-row>
      <n-row :gutter="[0, 24]">
        <n-col :span="11">
          <n-form-item path="php" :label="$t('websiteIndex.create.fields.phpVersion.label')">
            <n-select
              v-model:value="addModel.php"
              :options="installedDbAndPhp.php"
              :placeholder="$t('websiteIndex.create.fields.phpVersion.placeholder')"
              @keydown.enter.prevent
            >
            </n-select>
          </n-form-item>
        </n-col>
        <n-col :span="2"></n-col>
        <n-col :span="11">
          <n-form-item path="db" :label="$t('websiteIndex.create.fields.db.label')">
            <n-select
              v-model:value="addModel.db_type"
              :options="installedDbAndPhp.db"
              :placeholder="$t('websiteIndex.create.fields.db.placeholder')"
              @keydown.enter.prevent
              @update:value="
                () => {
                  addModel.db = addModel.db_type != '0'
                  addModel.db_name = formatDbValue(addModel.name)
                  addModel.db_user = formatDbValue(addModel.name)
                  addModel.db_password = generateRandomString(16)
                }
              "
            >
            </n-select>
          </n-form-item>
        </n-col>
      </n-row>
      <n-row :gutter="[0, 24]">
        <n-col :span="7">
          <n-form-item
            v-if="addModel.db"
            path="db_name"
            :label="$t('websiteIndex.create.fields.dbName.label')"
          >
            <n-input
              v-model:value="addModel.db_name"
              type="text"
              @keydown.enter.prevent
              :placeholder="$t('websiteIndex.create.fields.dbName.placeholder')"
            />
          </n-form-item>
        </n-col>
        <n-col :span="1"></n-col>
        <n-col :span="7">
          <n-form-item
            v-if="addModel.db"
            path="db_user"
            :label="$t('websiteIndex.create.fields.dbUser.label')"
          >
            <n-input
              v-model:value="addModel.db_user"
              type="text"
              @keydown.enter.prevent
              :placeholder="$t('websiteIndex.create.fields.dbUser.placeholder')"
            />
          </n-form-item>
        </n-col>
        <n-col :span="1"></n-col>
        <n-col :span="8">
          <n-form-item
            v-if="addModel.db"
            path="db_password"
            :label="$t('websiteIndex.create.fields.dbPassword.label')"
          >
            <n-input
              v-model:value="addModel.db_password"
              type="text"
              @keydown.enter.prevent
              :placeholder="$t('websiteIndex.create.fields.dbPassword.placeholder')"
            />
          </n-form-item>
        </n-col>
      </n-row>
      <n-form-item path="path" :label="$t('websiteIndex.create.fields.path.label')">
        <n-input
          v-model:value="addModel.path"
          type="text"
          @keydown.enter.prevent
          :placeholder="$t('websiteIndex.create.fields.path.placeholder')"
        />
      </n-form-item>
      <n-form-item path="remark" :label="$t('websiteIndex.create.fields.remark.label')">
        <n-input
          v-model:value="addModel.remark"
          type="textarea"
          @keydown.enter.prevent
          :placeholder="$t('websiteIndex.create.fields.remark.placeholder')"
        />
      </n-form-item>
    </n-form>
    <n-row :gutter="[0, 24]">
      <n-col :span="24">
        <n-button
          type="info"
          block
          :loading="buttonLoading"
          :disabled="buttonDisabled"
          @click="handleAdd"
        >
          {{ $t('websiteIndex.create.actions.submit') }}
        </n-button>
      </n-col>
    </n-row>
  </n-modal>
  <n-modal
    v-model:show="editDefaultPageModal"
    preset="card"
    title="修改默认页"
    style="width: 80vw"
    size="huge"
    :bordered="false"
    :segmented="false"
    @close="handleSaveDefaultPage"
  >
    <n-tabs type="line" animated>
      <n-tab-pane name="index" tab="默认页">
        <Editor
          v-model:value="editDefaultPageModel.index"
          language="html"
          theme="vs-dark"
          height="60vh"
          mt-8
          :options="{
            automaticLayout: true,
            formatOnType: true,
            formatOnPaste: true
          }"
        />
      </n-tab-pane>
      <n-tab-pane name="stop" tab="停止页">
        <Editor
          v-model:value="editDefaultPageModel.stop"
          language="html"
          theme="vs-dark"
          height="60vh"
          mt-8
          :options="{
            automaticLayout: true,
            formatOnType: true,
            formatOnPaste: true
          }"
        />
      </n-tab-pane>
    </n-tabs>
  </n-modal>
  <n-modal v-model:show="backupModal">
    <n-card closable @close="() => (backupModal = false)" title="备份管理" style="width: 60vw">
      <n-space vertical>
        <n-space>
          <n-button type="primary" @click="handleCreateBackup">创建备份</n-button>
          <n-upload
            accept=".zip,tar.gz,.tar,.rar,.bz2"
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
</template>
