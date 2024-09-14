<script setup lang="ts">
import CertView from '@/views/cert/CertView.vue'
import UserView from '@/views/cert/UserView.vue'
import DNSView from '@/views/cert/DNSView.vue'

const currentTab = ref('cert')
</script>

<template>
  <common-page show-footer>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="cert" tab="证书列表">
        <CertView />
      </n-tab-pane>
      <n-tab-pane name="user" tab="账号列表">
        <UserView />
      </n-tab-pane>
      <n-tab-pane name="dns" tab="DNS 列表">
        <DNSView />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>
