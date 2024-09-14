<script lang="ts" setup>
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<template>
  <n-icon cursor-pointer size="20" @click="themeStore.toggleCollapsed()">
    <icon-mdi:format-indent-increase v-if="themeStore.sider.collapsed" />
    <icon-mdi:format-indent-decrease v-else />
  </n-icon>
</template>
