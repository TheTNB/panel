<script setup lang="ts">
defineOptions({
  name: 'apps-mysql-index'
})

import Editor from '@guolao/vue-monaco-editor'
import { NButton, NDataTable, NInput, NPopconfirm } from 'naive-ui'

import mysql from '@/api/apps/mysql'
import systemctl from '@/api/panel/systemctl'

const currentTab = ref('status')
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

const load = ref<any[]>([])

const getLoad = async () => {
  const { data } = await mysql.load()
  return data
}

const getIsEnabled = async () => {
  await systemctl.isEnabled('mysqld').then((res: any) => {
    isEnabled.value = res.data
  })
}

const getStatus = async () => {
  await systemctl.status('mysqld').then((res: any) => {
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
  window.$message.success('清空成功')
}

const handleClearSlowLog = async () => {
  await mysql.clearSlowLog()
  window.$message.success('清空成功')
}

const handleIsEnabled = async () => {
  if (isEnabled.value) {
    await systemctl.enable('mysqld')
    window.$message.success('开启自启动成功')
  } else {
    await systemctl.disable('mysqld')
    window.$message.success('禁用自启动成功')
  }
  await getIsEnabled()
}

const handleStart = async () => {
  await systemctl.start('mysqld')
  window.$message.success('启动成功')
  await getStatus()
}

const handleStop = async () => {
  await systemctl.stop('mysqld')
  window.$message.success('停止成功')
  await getStatus()
}

const handleRestart = async () => {
  await systemctl.restart('mysqld')
  window.$message.success('重启成功')
  await getStatus()
}

const handleReload = async () => {
  await systemctl.reload('mysqld')
  window.$message.success('重载成功')
  await getStatus()
}

const handleSetRootPassword = async () => {
  await mysql.setRootPassword(rootPassword.value)
  window.$message.success('修改成功')
}

onMounted(() => {
  getStatus()
  getIsEnabled()
  getRootPassword()
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
        v-if="currentTab == 'error-log'"
        class="ml-16"
        type="primary"
        @click="handleClearErrorLog"
      >
        <TheIcon :size="18" icon="material-symbols:delete-outline" />
        清空日志
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
                  停止 MySQL 会导致使用 MySQL 的网站无法访问，确定要停止吗？
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
          <n-card title="Root 密码" rounded-10>
            <n-space vertical>
              <n-input v-model:value="rootPassword"></n-input>
              <n-button type="primary" @click="handleSetRootPassword">保存修改</n-button>
            </n-space>
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="config" tab="修改配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 MySQL 主配置文件，如果您不了解各参数的含义，请不要随意修改！
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
      <n-tab-pane name="slow-log" tab="慢查询日志">
        <realtime-log :path="slowLog" />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
