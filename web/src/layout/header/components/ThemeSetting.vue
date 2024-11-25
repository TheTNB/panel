<script setup lang="ts">
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<template>
  <n-tooltip trigger="hover">
    <template #trigger>
      <n-color-picker
        class="mr-16 h-24 w-24"
        :show-alpha="false"
        v-model:value="themeStore.primaryColor"
        :render-label="() => ''"
      />
    </template>
    设置主题色
  </n-tooltip>
</template>
