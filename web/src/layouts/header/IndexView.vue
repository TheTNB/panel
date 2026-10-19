<script lang="ts" setup>
import MenuCollapse from '@/layouts/header/components/MenuCollapse.vue'
import ReloadPage from '@/layouts/header/components/ReloadPage.vue'
import StickyNote from '@/layouts/header/components/StickyNote.vue'
import TaskStatus from '@/layouts/header/components/TaskStatus.vue'
import AppTab from '@/layouts/tab/IndexView.vue'
import { translateTitle } from '@/locales/menu'
import { useThemeStore } from '@/stores'

import FullScreen from './components/FullScreen.vue'
import ThemeMode from './components/ThemeMode.vue'
import UserAvatar from './components/UserAvatar.vue'

const themeStore = useThemeStore()
const route = useRoute()

// 手机上没有标签栏，靠标题告诉用户当前在哪个页面
const title = computed(() =>
  typeof route.meta.title === 'string' ? translateTitle(route.meta.title) : '',
)
</script>

<template>
  <div class="flex w-full items-center justify-between">
    <template v-if="themeStore.isMobile">
      <menu-collapse />
      <span class="text-base font-medium ml-3 min-w-0 truncate">{{ title }}</span>
    </template>
    <section v-if="!themeStore.isMobile && themeStore.tab.visible" class="pr-3 flex-1 w-0">
      <app-tab />
    </section>
    <span
      v-if="!themeStore.isMobile && themeStore.tab.visible"
      class="text-text-tertiary mx-2 opacity-40"
      >|</span
    >
    <div class="ml-auto pl-2 flex flex-shrink-0 gap-2 items-center md:gap-5">
      <task-status />
      <sticky-note />
      <template v-if="!themeStore.isMobile">
        <reload-page />
        <full-screen />
      </template>
      <theme-mode />
      <user-avatar />
    </div>
  </div>
</template>
