<script setup lang="ts">
import SettingHttps from '@/views/setting/SettingHttps.vue'
import SettingBase from '@/views/setting/SettingBase.vue'

const currentTab = ref('base')
</script>

<template>
  <common-page show-footer>
    <n-tabs v-model:value="currentTab" type="line" animated>
      <n-tab-pane name="base" tab="基本">
        <setting-base />
      </n-tab-pane>
      <n-tab-pane name="https" tab="HTTPS">
        <setting-https />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>

<style scoped lang="scss"></style>
