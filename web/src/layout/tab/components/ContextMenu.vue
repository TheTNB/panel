<script lang="ts" setup>
import { useTabStore } from '@/store'
import { renderIcon } from '@/utils'

interface Props {
  show?: boolean
  keepAlive?: boolean
  currentPath?: string
  x: number
  y: number
}

const props = withDefaults(defineProps<Props>(), {
  show: false,
  keepAlive: false,
  currentPath: ''
})

const emit = defineEmits(['update:show'])

const tabStore = useTabStore()

const options = computed(() => [
  {
    label: '关闭',
    key: 'close',
    disabled: tabStore.tabs.length <= 1,
    icon: renderIcon('mdi:close', { size: 14 })
  },
  {
    label: '重载',
    key: 'reload',
    disabled: props.currentPath !== tabStore.active,
    icon: renderIcon('mdi:refresh', { size: 14 })
  },
  {
    label: '固定',
    key: 'pin',
    disabled: props.keepAlive,
    icon: renderIcon('mdi:pin', { size: 14 })
  },
  {
    label: '取消固定',
    key: 'unpin',
    disabled: !props.keepAlive,
    icon: renderIcon('mdi:pin-off', { size: 14 })
  },
  {
    label: '关闭其他',
    key: 'close-other',
    disabled: tabStore.tabs.length <= 1,
    icon: renderIcon('mdi:arrow-expand-horizontal', { size: 14 })
  },
  {
    label: '关闭左侧',
    key: 'close-left',
    disabled: tabStore.tabs.length <= 1 || props.currentPath === tabStore.tabs[0].path,
    icon: renderIcon('mdi:arrow-expand-left', { size: 14 })
  },
  {
    label: '关闭右侧',
    key: 'close-right',
    disabled:
      tabStore.tabs.length <= 1 ||
      props.currentPath === tabStore.tabs[tabStore.tabs.length - 1].path,
    icon: renderIcon('mdi:arrow-expand-right', { size: 14 })
  }
])

const dropdownShow = computed({
  get() {
    return props.show
  },
  set(show) {
    emit('update:show', show)
  }
})

const actionMap = new Map([
  [
    'close',
    () => {
      tabStore.removeTab(props.currentPath)
    }
  ],
  [
    'reload',
    () => {
      tabStore.reloadTab(props.currentPath)
    }
  ],
  [
    'pin',
    () => {
      tabStore.pinTab(props.currentPath)
    }
  ],
  [
    'unpin',
    () => {
      tabStore.unpinTab(props.currentPath)
    }
  ],
  [
    'close-other',
    () => {
      tabStore.removeOther(props.currentPath)
    }
  ],
  [
    'close-left',
    () => {
      tabStore.removeLeft(props.currentPath)
    }
  ],
  [
    'close-right',
    () => {
      tabStore.removeRight(props.currentPath)
    }
  ]
])

function handleHideDropdown() {
  dropdownShow.value = false
}

function handleSelect(key: string) {
  const actionFn = actionMap.get(key)
  actionFn && actionFn()
  handleHideDropdown()
}
</script>

<template>
  <n-dropdown
    :options="options"
    :show="dropdownShow"
    :x="x"
    :y="y"
    placement="bottom-start"
    @clickoutside="handleHideDropdown"
    @select="handleSelect"
  />
</template>
