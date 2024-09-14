<script lang="ts" setup>
import { useAppStore } from '@/store'

const appStore = useAppStore()
</script>

<template>
  <router-view v-slot="{ Component, route }">
    <component :is="Component" v-if="appStore.reloadFlag" :key="route.path" />
  </router-view>
</template>
