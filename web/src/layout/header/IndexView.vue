<script lang="ts" setup>
import ReloadPage from '@/layout/header/components/ReloadPage.vue'
import BreadCrumb from './components/BreadCrumb.vue'
import FullScreen from './components/FullScreen.vue'
import MenuCollapse from './components/MenuCollapse.vue'
import ThemeMode from './components/ThemeMode.vue'
import UserAvatar from './components/UserAvatar.vue'
</script>

<template>
  <div flex items-center>
    <MenuCollapse />
    <BreadCrumb ml-15 hidden sm:block />
  </div>
  <div ml-auto flex items-center>
    <ReloadPage />
    <ThemeMode />
    <FullScreen />
    <UserAvatar />
  </div>
</template>
