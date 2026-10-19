<script lang="ts" setup>
import { useGettext } from 'vue3-gettext'

import home from '@/api/panel/home'

interface HealthIssue {
  key: string
  level: 'error' | 'warning'
  message: string
  since: string
}

const { $gettext } = useGettext()

const { data } = useAutoRequest(() => home.health(), { initialData: [] })

// 前端根据后端稳定 key 选择本地化文案，避免依赖后端错误详情
const titleFor = (issue: HealthIssue): string => {
  switch (issue.key) {
    case 'database:panel':
      return $gettext('Main database write failed, panel may not work properly')
    case 'database:stat':
      return $gettext('Website statistics database write failed, statistics may be inaccurate')
    case 'database:scan':
      return $gettext('Scan events database write failed, scan awareness may be inaccurate')
    default:
      return $gettext('Panel encountered a fault: %{key}', { key: issue.key })
  }
}

const hintFor = (issue: HealthIssue): string => {
  switch (issue.key) {
    case 'database:panel':
    case 'database:stat':
    case 'database:scan':
      return $gettext(
        'Try running "acepanel fix" on the server to detect and rebuild the broken database file.',
      )
    default:
      return ''
  }
}

// 一次性显示所有问题，最严重的靠前
const sorted = computed<HealthIssue[]>(() => {
  const order: Record<string, number> = { error: 0, warning: 1 }
  const list = (data.value ?? []) as HealthIssue[]
  return [...list].sort((a, b) => (order[a.level] ?? 9) - (order[b.level] ?? 9))
})
</script>

<template>
  <div v-if="sorted.length" class="border-b border-border-default flex flex-col gap-1">
    <n-alert
      v-for="issue in sorted"
      :key="issue.key"
      :type="issue.level"
      :show-icon="true"
      :bordered="false"
      class="rounded-none"
    >
      <div class="flex flex-col">
        <span class="font-medium">{{ titleFor(issue) }}</span>
        <span v-if="hintFor(issue)" class="text-xs opacity-80">{{ hintFor(issue) }}</span>
        <span v-if="issue.message" class="text-xs font-mono opacity-60 break-all">
          {{ issue.message }}
        </span>
      </div>
    </n-alert>
  </div>
</template>
