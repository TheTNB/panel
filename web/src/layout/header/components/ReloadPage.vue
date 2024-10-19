<script lang="ts" setup>
import { useTabStore } from '@/store'

const tabStore = useTabStore()

const handleReloadPage = () => {
  tabStore.reloadTab(tabStore.active)
}
</script>

<template>
  <n-tooltip trigger="hover">
    <template #trigger>
      <n-icon mr-20 cursor-pointer size="20" @click="handleReloadPage">
        <icon-mdi-refresh />
      </n-icon>
    </template>
    刷新标签页
  </n-tooltip>
</template>
