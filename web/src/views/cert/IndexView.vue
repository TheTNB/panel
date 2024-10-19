<script setup lang="ts">
defineOptions({
  name: 'cert-index'
})

import AccountView from '@/views/cert/AccountView.vue'
import CertView from '@/views/cert/CertView.vue'
import DnsView from '@/views/cert/DnsView.vue'

const currentTab = ref('cert')
</script>

<template>
  <common-page show-footer>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="cert" tab="证书列表">
        <cert-view />
      </n-tab-pane>
      <n-tab-pane name="user" tab="账号列表">
        <account-view />
      </n-tab-pane>
      <n-tab-pane name="dns" tab="DNS 列表">
        <dns-view />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
