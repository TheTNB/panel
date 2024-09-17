<script setup lang="ts">
import { NButton } from 'naive-ui'
import toolbox from '@/api/plugins/toolbox'
import Editor from '@guolao/vue-monaco-editor'

const currentTab = ref('dns')
const dns1 = ref('')
const dns2 = ref('')
const swap = ref(0)
const swapFree = ref('')
const swapUsed = ref('')
const swapTotal = ref('')
const hosts = ref('')
const timezone = ref('')
const timezones = ref<any[]>([])
const rootPassword = ref('')

const getDNS = async () => {
  await toolbox.dns().then((res: any) => {
    dns1.value = res.data[0]
    dns2.value = res.data[1]
  })
}

const getSwap = async () => {
  await toolbox.swap().then((res: any) => {
    swap.value = res.data.size
    swapFree.value = res.data.free
    swapUsed.value = res.data.used
    swapTotal.value = res.data.total
  })
}

const getHosts = async () => {
  toolbox.hosts().then((res: any) => {
    hosts.value = res.data
  })
}

const getTimezone = async () => {
  toolbox.timezone().then((res: any) => {
    timezone.value = res.data.timezone
    timezones.value = res.data.timezones
  })
}

const handleSaveDNS = async () => {
  await toolbox.setDns(dns1.value, dns2.value)
  window.$message.success('保存成功')
}

const handleSaveSwap = async () => {
  await toolbox.setSwap(swap.value)
  window.$message.success('保存成功')
}

const handleSaveHosts = async () => {
  await toolbox.setHosts(hosts.value)
  window.$message.success('保存成功')
}

const handleSaveTimezone = async () => {
  await toolbox.setTimezone(timezone.value)
  window.$message.success('保存成功')
}

const handleSaveRootPassword = async () => {
  await toolbox.setRootPassword(rootPassword.value)
  window.$message.success('保存成功')
}

onMounted(() => {
  getDNS()
  getSwap()
  getHosts()
  getTimezone()
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button v-if="currentTab == 'dns'" class="ml-16" type="primary" @click="handleSaveDNS">
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button v-if="currentTab == 'swap'" class="ml-16" type="primary" @click="handleSaveSwap">
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button v-if="currentTab == 'hosts'" class="ml-16" type="primary" @click="handleSaveHosts">
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        保存
      </n-button>
      <n-button
        v-if="currentTab == 'timezone'"
        class="ml-16"
        type="primary"
        @click="handleSaveTimezone"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        修改
      </n-button>
      <n-button
        v-if="currentTab == 'root-password'"
        class="ml-16"
        type="primary"
        @click="handleSaveRootPassword"
      >
        <TheIcon :size="18" class="mr-5" icon="material-symbols:save-outline" />
        修改
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="dns" tab="DNS">
        <n-form>
          <n-form-item label="DNS1">
            <n-input v-model:value="dns1" />
          </n-form-item>
          <n-form-item label="DNS2">
            <n-input v-model:value="dns2" />
          </n-form-item>
        </n-form>
      </n-tab-pane>
      <n-tab-pane name="swap" tab="SWAP">
        <n-space vertical>
          <n-alert type="info">
            总共 {{ swapTotal }}，已使用 {{ swapUsed }}，剩余 {{ swapFree }}
          </n-alert>
          <n-form>
            <n-form-item label="SWAP大小">
              <n-input-number v-model:value="swap" />
              MB
            </n-form-item>
          </n-form>
        </n-space>
      </n-tab-pane>
      <n-tab-pane name="hosts" tab="Host">
        <n-space vertical>
          <n-alert type="warning">
            此处修改的是系统 Hosts 文件，如果你不了解这是干什么的，请不要随意修改！
          </n-alert>
          <Editor
            v-model:value="hosts"
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
      <n-tab-pane name="timezone" tab="时区">
        <n-form>
          <n-form-item label="选择时区">
            <n-select v-model:value="timezone" placeholder="请选择时区" :options="timezones" />
          </n-form-item>
        </n-form>
      </n-tab-pane>
      <n-tab-pane name="root-password" tab="Root 密码">
        <n-form>
          <n-form-item label="Root 密码">
            <n-input v-model:value="rootPassword" />
          </n-form-item>
        </n-form>
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
