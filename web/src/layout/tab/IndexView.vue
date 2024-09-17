<script lang="ts" setup>
import ContextMenu from './components/ContextMenu.vue'
import type { TabItem } from '@/store'
import { useTabStore } from '@/store'

const route = useRoute()
const router = useRouter()
const tabStore = useTabStore()

interface ContextMenuOption {
  show: boolean
  x: number
  y: number
  currentPath: string
}

const contextMenuOption = reactive<ContextMenuOption>({
  show: false,
  x: 0,
  y: 0,
  currentPath: ''
})

watch(
  () => route.path,
  () => {
    const { name, fullPath: path } = route
    const title = (route.meta?.title as string) || ''
    tabStore.addTab({ name: name as string, path, title })
  },
  { immediate: true }
)

function handleTagClick(path: string) {
  tabStore.setActiveTab(path)
  router.push(path)
}

function showContextMenu() {
  contextMenuOption.show = true
}

function hideContextMenu() {
  contextMenuOption.show = false
}

function setContextMenu(x: number, y: number, currentPath: string) {
  Object.assign(contextMenuOption, { x, y, currentPath })
}

// 右击菜单
async function handleContextMenu(e: MouseEvent, tabItem: TabItem) {
  const { clientX, clientY } = e
  hideContextMenu()
  setContextMenu(clientX, clientY, tabItem.path)
  await nextTick()
  showContextMenu()
}
</script>

<template>
  <div>
    <n-tabs
      :value="tabStore.activeTab"
      :closable="tabStore.tabs.length > 1"
      type="card"
      @close="(path: string) => tabStore.removeTab(path)"
      bg-white
      dark:bg-dark
    >
      <n-tab
        v-for="item in tabStore.tabs"
        :key="item.path"
        :name="item.path"
        @click="handleTagClick(item.path)"
        @contextmenu.prevent="handleContextMenu($event, item)"
      >
        {{ $t(item.title!) }}
      </n-tab>
    </n-tabs>
    <ContextMenu
      v-model:show="contextMenuOption.show"
      :current-path="contextMenuOption.currentPath"
      :x="contextMenuOption.x"
      :y="contextMenuOption.y"
    />
  </div>
</template>

<style scoped lang="scss">
:deep(.n-tabs) {
  .n-tabs-tab {
    padding-left: 16px;
    height: 36px;
    background: transparent !important;
    border-radius: 4px !important;
    margin-right: 4px;

    &:hover {
      border: 1px solid var(--primary-color) !important;
    }
  }

  .n-tabs-tab--active {
    border: 1px solid var(--primary-color) !important;
    background-color: var(--selected-bg) !important;
  }

  .n-tabs-pad,
  .n-tabs-tab-pad,
  .n-tabs-scroll-padding {
    border: none !important;
  }
}
</style>
