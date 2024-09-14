<script lang="ts" setup>
import { useThemeStore } from '@/store'
import logo from '@/assets/images/logo.png'
import { title } from '@/main'

const themeStore = useThemeStore()
</script>

<template>
  <router-link class="h-60 f-c-c" to="/">
    <n-image :src="logo" height="42" preview-disabled />
    <h2
      v-show="!themeStore.sider.collapsed"
      class="ml-20 max-w-140 flex-shrink-0 text-18 font-bold"
    >
      {{ title }}
    </h2>
  </router-link>
</template>
