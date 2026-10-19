import type { Ref } from 'vue'

// 编辑器、任务队列等可拖拽窗口会同时存在：
// 后展开的在上层，最小化后停靠在右下角横向排开，有窗口展开时文件列表不响应快捷键
// 层级放在 Naive 从 2000 起的弹出层之下，窗口里的确认框、下拉菜单始终在最上面
const BASE_Z = 1900

const opened = ref<symbol[]>([])

export const hasOpenWindow = computed(() => opened.value.length > 0)

const dock = document.createElement('div')
Object.assign(dock.style, {
  position: 'fixed',
  right: '16px',
  bottom: '16px',
  zIndex: String(BASE_Z + 99),
  display: 'flex',
  flexDirection: 'row-reverse',
  flexWrap: 'wrap',
  maxWidth: 'calc(100vw - 32px)',
  gap: '8px',
})
document.body.appendChild(dock)

export function useWindowStack(open: Ref<boolean>) {
  const id = Symbol()
  const leave = () => (opened.value = opened.value.filter((item) => item !== id))

  watch(
    open,
    (value) => {
      leave()
      if (value) opened.value.push(id)
    },
    { immediate: true },
  )
  onBeforeUnmount(leave)

  const zIndex = computed(() => BASE_Z + 2 * Math.max(0, opened.value.indexOf(id)) + 1)

  return { zIndex, dock }
}
