<script lang="ts" setup>
import SideSetting from '@/layout/sidebar/components/SideSetting.vue'
import SideLogo from './components/SideLogo.vue'
import SideMenu from './components/SideMenu.vue'
</script>

<template>
  <div h-full flex flex-col>
    <side-logo />
    <side-menu />
    <side-setting mt-auto />
  </div>
</template>
