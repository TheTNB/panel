<script lang="ts" setup>
import { useAppStore } from '@/store'

const appStore = useAppStore()

const handleReloadPage = () => {
  appStore.reloadPage()
}
</script>

<template>
  <n-icon mr-20 cursor-pointer size="18" @click="handleReloadPage">
    <icon-mdi-refresh />
  </n-icon>
</template>
