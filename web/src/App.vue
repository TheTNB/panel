<script lang="ts" setup>
import AppProvider from '@/components/common/AppProvider.vue'
</script>

<template>
  <AppProvider>
    <router-view v-slot="{ Component }">
      <component :is="Component" />
    </router-view>
  </AppProvider>
</template>
