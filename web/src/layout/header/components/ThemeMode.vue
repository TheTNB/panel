<script lang="ts" setup>
import { useThemeStore } from '@/store'

const theme = useThemeStore()
</script>

<template>
  <n-icon mr-20 cursor-pointer size="18" @click="theme.toggleDarkMode">
    <icon-mdi-moon-waning-crescent v-if="theme.darkMode" />
    <icon-mdi-white-balance-sunny v-else />
  </n-icon>
</template>
