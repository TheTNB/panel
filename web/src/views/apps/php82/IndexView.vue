<script setup lang="ts">
defineOptions({
  name: 'apps-php82-index'
})

import Editor from '@guolao/vue-monaco-editor'
import { NButton, NDataTable, NPopconfirm } from 'naive-ui'

import php from '@/api/apps/php'
import systemctl from '@/api/panel/systemctl'
import { renderIcon } from '@/utils'

const currentTab = ref('status')
const version = Number(82)
const status = ref(false)
const isEnabled = ref(false)
const config = ref('')
const fpmConfig = ref('')
const errorLog = ref('')
const slowLog = ref('')

const statusType = computed(() => {
  return status.value ? 'success' : 'error'
})
const statusStr = computed(() => {
  return status.value ? '正常运行中' : '已停止运行'
})

const extensionColumns: any = [
  {
    title: '拓展名',
    key: 'name',
    minWidth: 250,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '描述',
    key: 'description',
    resizable: true,
    minWidth: 250,
    ellipsis: { tooltip: true }
  },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    align: 'center',
    hideInExcel: true,
    render(row: any) {
      return [
        !row.installed
          ? h(
              NButton,
              {
                size: 'small',
                type: 'primary',
                secondary: true,
                onClick: () => handleInstallExtension(row.slug)
              },
              {
                default: () => '安装',
                icon: renderIcon('material-symbols:download-rounded', { size: 14 })
              }
            )
          : null,
        row.installed
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleUninstallExtension(row.slug)
              },
              {
                default: () => {
                  return '确定卸载' + row.name + '吗？'
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'error'
                    },
                    {
                      default: () => '删除',
                      icon: renderIcon('material-symbols:delete-outline', { size: 14 })
                    }
                  )
                }
              }
            )
          : null
      ]
    }
  }
]

const loadColumns: any = [
  {
    title: '属性',
    key: 'name',
    minWidth: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: '当前值',
    key: 'value',
    minWidth: 200,
    ellipsis: { tooltip: true }
  }
]

const extensions = ref<any[]>([])
const load = ref<any[]>([])

const getLoad = async () => {
  const { data } = await php.load(version)
  return data
}

const getExtensions = async () => {
  const { data } = await php.extensions(version)
  return data
}

const getStatus = async () => {
  await systemctl.status(`php-fpm-${version}`).then((res: any) => {
    status.value = res.data
  })
}

const getIsEnabled = async () => {
  await systemctl.isEnabled(`php-fpm-${version}`).then((res: any) => {
    isEnabled.value = res.data
  })
}

const getErrorLog = async () => {
  php.errorLog(version).then((res: any) => {
    errorLog.value = res.data
  })
}

const getSlowLog = async () => {
  php.slowLog(version).then((res: any) => {
    slowLog.value = res.data
  })
}

const getConfig = async () => {
  php.config(version).then((res: any) => {
    config.value = res.data
  })
}

const getFPMConfig = async () => {
  php.fpmConfig(version).then((res: any) => {
    fpmConfig.value = res.data
  })
}

const handleSaveConfig = async () => {
  await php.saveConfig(version, config.value)
  window.$message.success('保存成功')
}

const handleSaveFPMConfig = async () => {
  await php.saveFPMConfig(version, fpmConfig.value)
  window.$message.success('保存成功')
  await getFPMConfig()
}

const handleClearErrorLog = async () => {
  await php.clearErrorLog(version)
  window.$message.success('清空成功')
}

const handleClearSlowLog = async () => {
  await php.clearSlowLog(version)
  window.$message.success('清空成功')
}

const handleIsEnabled = async () => {
  if (isEnabled.value) {
    await systemctl.enable(`php-fpm-${version}`)
    window.$message.success('开启自启动成功')
  } else {
    await systemctl.disable(`php-fpm-${version}`)
    window.$message.success('禁用自启动成功')
  }
  await getIsEnabled()
}

const handleStart = async () => {
  await systemctl.start(`php-fpm-${version}`)
  window.$message.success('启动成功')
  await getStatus()
}

const handleStop = async () => {
  await systemctl.stop(`php-fpm-${version}`)
  window.$message.success('停止成功')
  await getStatus()
}

const handleRestart = async () => {
  await systemctl.restart(`php-fpm-${version}`)
  window.$message.success('重启成功')
  await getStatus()
}

const handleReload = async () => {
  await systemctl.reload(`php-fpm-${version}`)
  window.$message.success('重载成功')
  await getStatus()
}

const handleInstallExtension = async (slug: string) => {
  await php.installExtension(version, slug).then(() => {
    window.$message.success('任务已提交，请稍后查看任务进度')
  })
}

const handleUninstallExtension = async (name: string) => {
  await php.uninstallExtension(version, name).then(() => {
    window.$message.success('任务已提交，请稍后查看任务进度')
  })
}

onMounted(() => {
  getStatus()
  getIsEnabled()
  getExtensions().then((res) => {
    extensions.value = res
  })
  getLoad().then((res) => {
    load.value = res
  })
  getErrorLog()
  getSlowLog()
  getConfig()
  getFPMConfig()
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button
        v-if="currentTab == 'config'"
        class="ml-16"
        type="primary"
        @click="handleSaveConfig"
      >
        <TheIcon :size="18" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'fpm-config'"
        class="ml-16"
        type="primary"
        @click="handleSaveFPMConfig"
      >
        <TheIcon :size="18" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'error-log'"
        class="ml-16"
        type="primary"
        @click="handleClearErrorLog"
      >
        <TheIcon :size="18" icon="material-symbols:delete-outline" />
        清空错误日志
      </n-button>
      <n-button
        v-if="currentTab == 'slow-log'"
        class="ml-16"
        type="primary"
        @click="handleClearSlowLog"
      >
        <TheIcon :size="18" icon="material-symbols:delete-outline" />
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
                  <TheIcon :size="24" icon="material-symbols:play-arrow-outline-rounded" />
                  启动
                </n-button>
                <n-popconfirm @positive-click="handleStop">
                  <template #trigger>
                    <n-button type="error">
                      <TheIcon :size="24" icon="material-symbols:stop-outline-rounded" />
                      停止
                    </n-button>
                  </template>
                  停止 PHP {{ version }} 会导致使用 PHP {{ version }} 的网站无法访问，确定要停止吗？
                </n-popconfirm>
                <n-button type="warning" @click="handleRestart">
                  <TheIcon :size="18" icon="material-symbols:replay-rounded" />
                  重启
                </n-button>
                <n-button type="primary" @click="handleReload">
                  <TheIcon :size="20" icon="material-symbols:refresh-rounded" />
                  重载
                </n-button>
              </n-space>
            </n-space>
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="extensions" tab="拓展管理">
        <n-card title="拓展列表" :segmented="true" rounded-10>
          <n-data-table
            striped
            remote
            :scroll-x="1000"
            :loading="false"
            :columns="extensionColumns"
            :data="extensions"
            :row-key="(row: any) => row.slug"
          />
        </n-card>
      </n-tab-pane>
      <n-tab-pane name="config" tab="主配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 PHP {{ version }} 主配置文件，如果您不了解各参数的含义，请不要随意修改！
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
      <n-tab-pane name="fpm-config" tab="FPM 配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 PHP {{ version }} FPM 配置文件，如果您不了解各参数的含义，请不要随意修改！
          </n-alert>
          <Editor
            v-model:value="fpmConfig"
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
        <n-data-table
          striped
          remote
          :scroll-x="400"
          :loading="false"
          :columns="loadColumns"
          :data="load"
        />
      </n-tab-pane>
      <n-tab-pane name="error-log" tab="错误日志">
        <realtime-log :path="errorLog" />
      </n-tab-pane>
      <n-tab-pane name="slow-log" tab="慢日志">
        <realtime-log :path="slowLog" />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
