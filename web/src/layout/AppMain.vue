<script lang="ts" setup>
import { useTabStore } from '@/store'

const tabStore = useTabStore()

const keepAliveNames = computed(() => {
  return tabStore.tabs.filter((item) => item.keepAlive).map((item) => item.name)
})
</script>

<template>
  <router-view v-slot="{ Component, route }">
    <keep-alive :include="keepAliveNames">
      <component :is="Component" v-if="!tabStore.reloading" :key="route.path" />
    </keep-alive>
  </router-view>
</template>
