<script setup lang="ts">
import { NButton } from 'naive-ui'

defineOptions({
  name: 'task-index'
})

import TheIcon from '@/components/custom/TheIcon.vue'
import CreateModal from '@/views/task/CreateModal.vue'
import CronView from '@/views/task/CronView.vue'
import TaskView from '@/views/task/TaskView.vue'

const current = ref('cron')

const create = ref(false)
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button v-if="current == 'cron'" type="primary" @click="create = true">
        <TheIcon :size="18" icon="material-symbols:add" />
        创建任务
      </n-button>
    </template>
    <n-tabs v-model:value="current" type="line" animated size="large">
      <n-tab-pane name="cron" tab="计划任务">
        <cron-view />
      </n-tab-pane>
      <n-tab-pane name="task" tab="面板任务">
        <task-view />
      </n-tab-pane>
    </n-tabs>
  </common-page>
  <create-modal v-model:show="create" />
</template>

<style scoped lang="scss"></style>
