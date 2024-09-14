<script setup lang="ts">
import { NButton, NPopconfirm } from 'naive-ui'
import Editor from '@guolao/vue-monaco-editor'
import podman from '@/api/plugins/podman'
import service from '@/api/panel/system/service'

const currentTab = ref('status')
const status = ref(false)
const isEnabled = ref(false)
const registryConfig = ref('')
const storageConfig = ref('')

const statusStr = computed(() => {
  return status.value ? '正常运行中' : '已停止运行'
})

const getStatus = async () => {
  await service.status('podman').then((res: any) => {
    status.value = res.data
  })
}

const getIsEnabled = async () => {
  await service.isEnabled('podman').then((res: any) => {
    isEnabled.value = res.data
  })
}

const getConfig = async () => {
  podman.registryConfig().then((res: any) => {
    registryConfig.value = res.data
  })
  podman.storageConfig().then((res: any) => {
    storageConfig.value = res.data
  })
}

const handleSaveRegistryConfig = async () => {
  await podman.saveRegistryConfig(registryConfig.value)
  window.$message.success('保存成功')
}

const handleSaveStorageConfig = async () => {
  await podman.saveStorageConfig(storageConfig.value)
  window.$message.success('保存成功')
}

const handleStart = async () => {
  await service.start('podman')
  window.$message.success('启动成功')
  await getStatus()
}

const handleStop = async () => {
  await service.stop('podman')
  window.$message.success('停止成功')
  await getStatus()
}

const handleRestart = async () => {
  await service.restart('podman')
  window.$message.success('重启成功')
  await getStatus()
}

const handleIsEnabled = async () => {
  if (isEnabled.value) {
    await service.enable('podman')
    window.$message.success('开启自启动成功')
  } else {
    await service.disable('podman')
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
        v-if="currentTab == 'registryConfig'"
        class="ml-16"
        type="primary"
        @click="handleSaveRegistryConfig"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'storageConfig'"
        class="ml-16"
        type="primary"
        @click="handleSaveStorageConfig"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="status" tab="运行状态">
        <n-flex vertical>
          <n-alert type="info">
            Podman 是一个无守护进程的容器管理工具，处于停止状态为正常现象且不会影响使用！
          </n-alert>
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
                      <TheIcon
                        :size="24"
                        class="mr-5"
                        icon="material-symbols:stop-outline-rounded"
                      />
                      停止
                    </n-button>
                  </template>
                  确定要停止 Podman 吗？
                </n-popconfirm>
                <n-button type="warning" @click="handleRestart">
                  <TheIcon :size="18" class="mr-5" icon="material-symbols:replay-rounded" />
                  重启
                </n-button>
              </n-space>
            </n-space>
          </n-card>
        </n-flex>
      </n-tab-pane>
      <n-tab-pane name="registryConfig" tab="注册表配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 Podman 注册表配置文件（/etc/containers/registries.conf）
          </n-alert>
          <Editor
            v-model:value="registryConfig"
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
      <n-tab-pane name="storageConfig" tab="存储配置">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是 Podman 存储配置文件（/etc/containers/storage.conf）
          </n-alert>
          <Editor
            v-model:value="storageConfig"
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
