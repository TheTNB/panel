<script lang="ts" setup>
import ReloadPage from '@/layout/header/components/ReloadPage.vue'
import AppTab from '@/layout/tab/IndexView.vue'
import FullScreen from './components/FullScreen.vue'
import ThemeMode from './components/ThemeMode.vue'
import UserAvatar from './components/UserAvatar.vue'

import MenuCollapse from '@/layout/header/components/MenuCollapse.vue'
import ThemeSetting from '@/layout/header/components/ThemeSetting.vue'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<template>
  <div w-full flex items-center justify-between>
    <MenuCollapse v-if="themeStore.isMobile" />
    <section v-if="!themeStore.isMobile && themeStore.tab.visible" w-0 flex-1 pr-12>
      <AppTab />
    </section>
    <span v-if="!themeStore.isMobile && themeStore.tab.visible" mx-6 opacity-20>|</span>
    <div ml-auto flex flex-shrink-0 items-center px-12>
      <ReloadPage />
      <FullScreen />
      <ThemeMode />
      <ThemeSetting />
      <UserAvatar />
    </div>
  </div>
</template>
