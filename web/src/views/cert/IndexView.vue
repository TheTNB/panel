<script setup lang="ts">
defineOptions({
  name: 'cert-index'
})

import { NButton } from 'naive-ui'

import app from '@/api/panel/app'
import cert from '@/api/panel/cert'
import website from '@/api/panel/website'
import AccountView from '@/views/cert/AccountView.vue'
import CertView from '@/views/cert/CertView.vue'
import CreateAccountModal from '@/views/cert/CreateAccountModal.vue'
import CreateCertModal from '@/views/cert/CreateCertModal.vue'
import CreateDnsModal from '@/views/cert/CreateDnsModal.vue'
import DnsView from '@/views/cert/DnsView.vue'

const currentTab = ref('cert')

const createDNS = ref(false)
const createCert = ref(false)
const createAccount = ref(false)

const algorithms = ref<any>([])
const websites = ref<any>([])
const dns = ref<any>([])
const accounts = ref<any>([])
const dnsProviders = ref<any>([])
const caProviders = ref<any>([])

const getAsyncData = async () => {
  const { data: algorithmData } = await cert.algorithms()
  algorithms.value = algorithmData

  websites.value = []
  websites.value.push({
    label: '无',
    value: 0
  })
  app.isInstalled('nginx').then(async (res) => {
    if (res.data.installed) {
      const { data: websiteData } = await website.list(1, 10000)
      for (const item of websiteData.items) {
        websites.value.push({
          label: item.name,
          value: item.id
        })
      }
    }
  })

  const { data: dnsData } = await cert.dns(1, 10000)
  dns.value = []
  dns.value.push({
    label: '无',
    value: 0
  })
  for (const item of dnsData.items) {
    dns.value.push({
      label: item.name,
      value: item.id
    })
  }

  const { data: accountData } = await cert.accounts(1, 10000)
  accounts.value = []
  for (const item of accountData.items) {
    accounts.value.push({
      label: item.email,
      value: item.id
    })
  }

  const { data: dnsProviderData } = await cert.dnsProviders()
  dnsProviders.value = dnsProviderData

  const { data: caProviderData } = await cert.caProviders()
  caProviders.value = caProviderData
}

onMounted(() => {
  getAsyncData()
  window.$bus.on('cert:refresh-async', getAsyncData)
})

onUnmounted(() => {
  window.$bus.off('cert:refresh-async')
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button v-if="currentTab == 'cert'" type="primary" @click="createCert = true">
        <TheIcon :size="18" icon="material-symbols:add" />
        创建证书
      </n-button>
      <n-button v-if="currentTab == 'user'" type="primary" @click="createAccount = true">
        <TheIcon :size="18" icon="material-symbols:add" />
        创建账号
      </n-button>
      <n-button v-if="currentTab == 'dns'" type="primary" @click="createDNS = true">
        <TheIcon :size="18" icon="material-symbols:add" />
        创建 DNS
      </n-button>
    </template>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="cert" tab="证书列表">
        <cert-view
          :accounts="accounts"
          :algorithms="algorithms"
          :websites="websites"
          :dns="dns"
          :ca-providers="caProviders"
        />
      </n-tab-pane>
      <n-tab-pane name="user" tab="账号列表">
        <account-view :ca-providers="caProviders" :algorithms="algorithms" />
      </n-tab-pane>
      <n-tab-pane name="dns" tab="DNS 列表">
        <dns-view :dns-providers="dnsProviders" />
      </n-tab-pane>
    </n-tabs>
  </common-page>
  <create-cert-modal
    v-model:show="createCert"
    :accounts="accounts"
    :algorithms="algorithms"
    :websites="websites"
    :dns="dns"
  />
  <create-dns-modal v-model:show="createDNS" :dns-providers="dnsProviders" />
  <create-account-modal
    v-model:show="createAccount"
    :ca-providers="caProviders"
    :algorithms="algorithms"
  />
</template>
