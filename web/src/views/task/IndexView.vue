<script setup lang="ts">
defineOptions({
  name: 'task-index'
})

import CronView from '@/views/task/CronView.vue'
import TaskView from '@/views/task/TaskView.vue'

const current = ref('cron')
</script>

<template>
  <common-page show-footer>
    <n-tabs v-model:value="current" type="line" animated size="large">
      <n-tab-pane name="cron" tab="计划任务">
        <cron-view />
      </n-tab-pane>
      <n-tab-pane name="task" tab="后台任务">
        <task-view />
      </n-tab-pane>
    </n-tabs>
  </common-page>
</template>

<style scoped lang="scss"></style>
