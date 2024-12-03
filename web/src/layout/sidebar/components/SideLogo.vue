<script lang="ts" setup>
import logoImg from '@/assets/images/logo.png'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
const logo = computed(() => themeStore.logo || logoImg)
</script>

<template>
  <router-link class="h-60 f-c-c" to="/">
    <n-image :src="logo" height="32" preview-disabled />
    <h2
      v-show="!themeStore.sider.collapsed"
      class="ml-10 max-w-140 flex-shrink-0 text-18 font-bold"
    >
      {{ themeStore.name }}
    </h2>
  </router-link>
</template>
