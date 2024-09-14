<script lang="ts" setup>
import { renderIcon } from '@/utils'
import type { Meta } from '~/types/router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const generator: any = (routerMap: any) => {
  return routerMap.map((item: any) => {
    const currentMenu = {
      ...item,
      label: t(String(item.meta.title)),
      key: item.path,
      disabled: item.path === '/',
      icon: getIcon(item.meta)
    }
    // 是否有子菜单，并递归处理
    if (item.children && item.children.length > 0) {
      currentMenu.children = generator(item.children, currentMenu)
    }
    return currentMenu
  })
}

const breadcrumbList = computed(() => {
  return generator(route.matched)
})

function handleBreadClick(path: string) {
  if (path === route.path) return
  router.push(path)
}

function getIcon(meta?: Meta, size = 16) {
  if (meta?.icon) return renderIcon(meta.icon, { size })
  return ''
}
</script>

<template>
  <n-breadcrumb>
    <template v-for="routeItem in breadcrumbList" :key="routeItem.name">
      <n-breadcrumb-item v-if="routeItem.meta.title">
        <n-dropdown
          v-if="routeItem.children.length"
          :options="routeItem.children"
          @select="handleBreadClick"
        >
          <span class="link-text">
            <component :is="routeItem.icon" v-if="routeItem.icon" />
            {{ $t(routeItem.meta.title) }}
          </span>
        </n-dropdown>
        <span v-else class="link-text">
          <component :is="routeItem.icon" v-if="routeItem.icon" />
          {{ $t(routeItem.meta.title) }}
        </span>
      </n-breadcrumb-item>
    </template>
  </n-breadcrumb>
</template>
