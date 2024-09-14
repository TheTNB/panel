<script setup lang="ts">
import { NButton, NPopconfirm } from 'naive-ui'
import Editor from '@guolao/vue-monaco-editor'
import frp from '@/api/plugins/frp'
import service from '@/api/panel/system/service'

const currentTab = ref('frps')
const status = ref({
  frpc: false,
  frps: false
})
const isEnabled = ref({
  frpc: false,
  frps: false
})
const config = ref({
  frpc: '',
  frps: ''
})

const statusStr = computed(() => {
  return {
    frpc: status.value.frpc ? '正常运行中' : '已停止运行',
    frps: status.value.frps ? '正常运行中' : '已停止运行'
  }
})

const getStatus = async () => {
  await service.status('frps').then((res: any) => {
    status.value.frps = res.data
  })
  await service.status('frpc').then((res: any) => {
    status.value.frpc = res.data
  })
}

const getIsEnabled = async () => {
  await service.isEnabled('frps').then((res: any) => {
    isEnabled.value.frps = res.data
  })
  await service.isEnabled('frpc').then((res: any) => {
    isEnabled.value.frpc = res.data
  })
}

const getConfig = async () => {
  frp.config('frps').then((res: any) => {
    config.value.frps = res.data
  })
  frp.config('frpc').then((res: any) => {
    config.value.frpc = res.data
  })
}

const handleSaveConfig = async (service: string) => {
  await frp.saveConfig(service, config.value[service as keyof typeof config.value])
  window.$message.success('保存成功')
}

const handleStart = async (name: string) => {
  await service.start(name)
  window.$message.success('启动成功')
  await getStatus()
}

const handleStop = async (name: string) => {
  await service.stop(name)
  window.$message.success('停止成功')
  await getStatus()
}

const handleRestart = async (name: string) => {
  await service.restart(name)
  window.$message.success('重启成功')
  await getStatus()
}

const handleIsEnabled = async (name: string) => {
  if (isEnabled.value[name as keyof typeof isEnabled.value]) {
    await service.enable(name)
    window.$message.success('开启自启动成功')
  } else {
    await service.disable(name)
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
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="frps" tab="Frps">
        <n-space vertical>
          <n-card title="运行状态" rounded-10>
            <template #header-extra>
              <n-switch v-model:value="isEnabled.frps" @update:value="handleIsEnabled('frps')">
                <template #checked> 自启动开 </template>
                <template #unchecked> 自启动关 </template>
              </n-switch>
            </template>
            <n-space vertical>
              <n-alert :type="status.frps ? 'success' : 'error'">
                {{ statusStr.frps }}
              </n-alert>
              <n-space>
                <n-button type="success" @click="handleStart('frps')">
                  <TheIcon
                    :size="24"
                    class="mr-5"
                    icon="material-symbols:play-arrow-outline-rounded"
                  />
                  启动
                </n-button>
                <n-popconfirm @positive-click="handleStop('frps')">
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
                  确定要停止 Frps 吗？
                </n-popconfirm>
                <n-button type="warning" @click="handleRestart('frps')">
                  <TheIcon :size="18" class="mr-5" icon="material-symbols:replay-rounded" />
                  重启
                </n-button>
              </n-space>
            </n-space>
          </n-card>
          <n-card title="修改配置" rounded-10>
            <template #header-extra>
              <n-button type="primary" @click="handleSaveConfig('frps')">
                <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline-rounded" />
                保存
              </n-button>
            </template>
            <Editor
              v-model:value="config.frps"
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
          </n-card>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="frpc" tab="Frpc">
        <n-space vertical>
          <n-card title="运行状态" rounded-10>
            <template #header-extra>
              <n-switch v-model:value="isEnabled.frpc" @update:value="handleIsEnabled('frpc')">
                <template #checked> 自启动开 </template>
                <template #unchecked> 自启动关 </template>
              </n-switch>
            </template>
            <n-space vertical>
              <n-alert :type="status.frpc ? 'success' : 'error'">
                {{ statusStr.frpc }}
              </n-alert>
              <n-space>
                <n-button type="success" @click="handleStart('frpc')">
                  <TheIcon
                    :size="24"
                    class="mr-5"
                    icon="material-symbols:play-arrow-outline-rounded"
                  />
                  启动
                </n-button>
                <n-popconfirm @positive-click="handleStop('frpc')">
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
                  确定要停止 Frpc 吗？
                </n-popconfirm>
                <n-button type="warning" @click="handleRestart('frpc')">
                  <TheIcon :size="18" class="mr-5" icon="material-symbols:replay-rounded" />
                  重启
                </n-button>
              </n-space>
            </n-space>
          </n-card>
          <n-card title="修改配置" rounded-10>
            <template #header-extra>
              <n-button type="primary" @click="handleSaveConfig('frpc')">
                <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline-rounded" />
                保存
              </n-button>
            </template>
            <Editor
              v-model:value="config.frpc"
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
          </n-card>
        </n-space>
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
