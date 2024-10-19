<script lang="ts" setup>
import AppMain from './AppMain.vue'
import AppHeader from './header/IndexView.vue'
import SideBar from './sidebar/IndexView.vue'

import { useThemeStore } from '@/store'

const themeStore = useThemeStore()
</script>

<template>
  <n-layout has-sider wh-full>
    <n-layout-sider
      v-if="!themeStore.isMobile"
      :collapsed="themeStore.sider.collapsed"
      :collapsed-width="themeStore.sider.collapsedWidth"
      :native-scrollbar="false"
      :width="themeStore.sider.width"
      bordered
      collapse-mode="width"
    >
      <SideBar />
    </n-layout-sider>
    <n-drawer
      v-else
      :auto-focus="false"
      :show="!themeStore.sider.collapsed"
      :width="themeStore.sider.width"
      display-directive="show"
      placement="left"
      @mask-click="themeStore.setCollapsed(true)"
    >
      <SideBar />
    </n-drawer>

    <article flex-col flex-1 overflow-hidden>
      <header
        :style="`height: ${themeStore.header.height}px`"
        dark="bg-dark border-0"
        flex
        items-center
        border-b
        bg-white
        px-15
        bc-eee
      >
        <AppHeader />
      </header>
      <section bg="#f5f6fb" flex-1 overflow-hidden dark:bg-hex-101014>
        <AppMain />
      </section>
    </article>
  </n-layout>
</template>
