<script setup lang="ts">
import { NButton, NPopconfirm } from 'naive-ui'
import Editor from '@guolao/vue-monaco-editor'
import gitea from '@/api/plugins/gitea'
import service from '@/api/panel/system/service'

const currentTab = ref('status')
const status = ref(false)
const isEnabled = ref(false)
const config = ref('')

const statusStr = computed(() => {
  return status.value ? '正常运行中' : '已停止运行'
})

const getStatus = async () => {
  await service.status('gitea').then((res: any) => {
    status.value = res.data
  })
}

const getIsEnabled = async () => {
  await service.isEnabled('gitea').then((res: any) => {
    isEnabled.value = res.data
  })
}

const getConfig = async () => {
  gitea.config().then((res: any) => {
    config.value = res.data
  })
}

const handleSaveConfig = async () => {
  await gitea.saveConfig(config.value)
  window.$message.success('保存成功')
}

const handleStart = async () => {
  await service.start('gitea')
  window.$message.success('启动成功')
  await getStatus()
}

const handleStop = async () => {
  await service.stop('gitea')
  window.$message.success('停止成功')
  await getStatus()
}

const handleRestart = async () => {
  await service.restart('gitea')
  window.$message.success('重启成功')
  await getStatus()
}

const handleIsEnabled = async () => {
  if (isEnabled.value) {
    await service.enable('gitea')
    window.$message.success('开启自启动成功')
  } else {
    await service.disable('gitea')
    window.$message.success('禁用自启动成功')
  }
  await getIsEnabled()
}

onMounted(() => {
  getStatus()
  getIsEnabled()
  getConfig()
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
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="status" tab="运行状态">
        <n-card title="运行状态" rounded-10>
          <template #header-extra>
            <n-switch v-model:value="isEnabled" @update:value="handleIsEnabled">
              <template #checked> 自启动开 </template>
              <template #unchecked> 自启动关 </template>
            </n-switch>
          </template>
          <n-space vertical>
            <n-alert :type="status ? 'success' : 'error'">
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
                    <TheIcon :size="24" class="mr-5" icon="material-symbols:stop-outline-rounded" />
                    停止
                  </n-button>
                </template>
                确定要停止 Gitea 吗？
              </n-popconfirm>
              <n-button type="warning" @click="handleRestart">
                <TheIcon :size="18" class="mr-5" icon="material-symbols:replay-rounded" />
                重启
              </n-button>
            </n-space>
          </n-space>
        </n-card>
      </n-tab-pane>
      <n-tab-pane name="config" tab="修改配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 Gitea 配置文件，如果你不了解各参数的含义，请不要随意修改！
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
    </n-tabs>
  </common-page>
</template>
