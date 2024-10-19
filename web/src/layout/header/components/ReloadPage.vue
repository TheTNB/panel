<script lang="ts" setup>
import { useTabStore } from '@/store'

const tabStore = useTabStore()

const handleReloadPage = () => {
  tabStore.reloadTab(tabStore.active)
}
</script>

<template>
  <n-icon mr-20 cursor-pointer size="18" @click="handleReloadPage">
    <icon-mdi-refresh />
  </n-icon>
</template>
