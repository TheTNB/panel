<script setup lang="ts">
defineOptions({
  name: 'firewall-index'
})

import ForwardView from '@/views/firewall/ForwardView.vue'
import IpRuleView from '@/views/firewall/IpRuleView.vue'
import RuleView from '@/views/firewall/RuleView.vue'
import SettingView from '@/views/firewall/SettingView.vue'

const currentTab = ref('rule')
</script>

<template>
  <common-page show-footer>
    <n-tabs v-model:value="currentTab" type="line" animated size="large">
      <n-tab-pane name="rule" tab="端口规则">
        <rule-view />
      </n-tab-pane>
      <n-tab-pane name="ip-rule" tab="IP 规则">
        <ip-rule-view />
      </n-tab-pane>
      <n-tab-pane name="forward" tab="端口转发">
        <forward-view />
      </n-tab-pane>
      <n-tab-pane name="setting" tab="设置">
        <setting-view />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>

<style scoped lang="scss"></style>
