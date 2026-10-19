<script setup lang="ts">
import { useThemeVars } from 'naive-ui'
import { useGettext } from 'vue3-gettext'

import { useThemeStore } from '@/stores'

import { useWindowStack } from './composables/useWindowStack'

const { $gettext } = useGettext()
const themeVars = useThemeVars()
const themeStore = useThemeStore()

const props = withDefaults(
  defineProps<{
    title?: string
    minWidth?: number
    minHeight?: number
    defaultWidth?: number
    defaultHeight?: number
    beforeClose?: () => Promise<boolean> | boolean // 关闭前的确认回调，返回 true 继续关闭，false 取消关闭
    closeOnOverlay?: boolean // 点击遮罩层是否最小化窗口，默认 true
  }>(),
  {
    title: '',
    minWidth: 400,
    minHeight: 300,
    defaultWidth: 800,
    defaultHeight: 600,
    closeOnOverlay: true,
  },
)

const show = defineModel<boolean>('show', { default: false })
const minimized = defineModel<boolean>('minimized', { default: false })

const { zIndex, dock } = useWindowStack(computed(() => show.value && !minimized.value))

// 窗口状态
const isMaximized = ref(false)
// 最大化或手机全屏时位置尺寸固定，不能拖拽和缩放
const locked = computed(() => isMaximized.value || themeStore.isMobile)

// 窗口位置和大小
const position = ref({ x: 0, y: 0 })
const size = ref({ width: props.defaultWidth, height: props.defaultHeight })

// 最大化前的状态（用于恢复）
const beforeMaximize = ref({ x: 0, y: 0, width: 0, height: 0 })

// 拖拽状态
const isDragging = ref(false)
const dragStart = ref({ x: 0, y: 0 })

// 调整大小状态
const isResizing = ref(false)
const resizeDirection = ref('')
const resizeStart = ref({ x: 0, y: 0, width: 0, height: 0, posX: 0, posY: 0 })

// 窗口样式
const windowStyle = computed(() => ({
  ...(themeStore.isMobile
    ? { left: '0', top: '0', width: '100%', height: '100%' }
    : {
        left: position.value.x + 'px',
        top: position.value.y + 'px',
        width: size.value.width + 'px',
        height: size.value.height + 'px',
      }),
  zIndex: zIndex.value,
  background: themeVars.value.cardColor,
  '--border-color': themeVars.value.borderColor,
  '--text-color-1': themeVars.value.textColor1,
  '--text-color-2': themeVars.value.textColor2,
  '--text-color-3': themeVars.value.textColor3,
  '--hover-color': themeVars.value.buttonColor2Hover,
  '--primary-color': themeVars.value.primaryColor,
  '--border-radius': themeVars.value.borderRadius,
}))

// 按默认尺寸居中，视口放不下时先缩到视口大小，保证坐标不为负、标题栏按钮可见
function initPosition() {
  const width = Math.min(props.defaultWidth, window.innerWidth)
  const height = Math.min(props.defaultHeight, window.innerHeight)
  size.value = { width, height }
  position.value = {
    x: (window.innerWidth - width) / 2,
    y: (window.innerHeight - height) / 2,
  }
}

// 开始拖拽
function startDrag(e: MouseEvent) {
  if (locked.value) return
  isDragging.value = true
  dragStart.value = {
    x: e.clientX - position.value.x,
    y: e.clientY - position.value.y,
  }
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
}

function onDrag(e: MouseEvent) {
  if (!isDragging.value) return
  position.value = {
    x: Math.max(0, Math.min(window.innerWidth - size.value.width, e.clientX - dragStart.value.x)),
    y: Math.max(0, Math.min(window.innerHeight - size.value.height, e.clientY - dragStart.value.y)),
  }
}

function stopDrag() {
  isDragging.value = false
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
}

// 开始调整大小
function startResize(e: MouseEvent, direction: string) {
  if (locked.value) return
  e.preventDefault()
  e.stopPropagation()
  isResizing.value = true
  resizeDirection.value = direction
  resizeStart.value = {
    x: e.clientX,
    y: e.clientY,
    width: size.value.width,
    height: size.value.height,
    posX: position.value.x,
    posY: position.value.y,
  }
  document.addEventListener('mousemove', onResize)
  document.addEventListener('mouseup', stopResize)
}

function onResize(e: MouseEvent) {
  if (!isResizing.value) return

  const deltaX = e.clientX - resizeStart.value.x
  const deltaY = e.clientY - resizeStart.value.y
  const dir = resizeDirection.value

  let newWidth = resizeStart.value.width
  let newHeight = resizeStart.value.height
  let newX = resizeStart.value.posX
  let newY = resizeStart.value.posY

  // 右边
  if (dir.includes('e')) {
    newWidth = Math.max(props.minWidth, resizeStart.value.width + deltaX)
  }
  // 左边
  if (dir.includes('w')) {
    const maxDelta = resizeStart.value.width - props.minWidth
    const actualDelta = Math.min(deltaX, maxDelta)
    newWidth = resizeStart.value.width - actualDelta
    newX = resizeStart.value.posX + actualDelta
  }
  // 下边
  if (dir.includes('s')) {
    newHeight = Math.max(props.minHeight, resizeStart.value.height + deltaY)
  }
  // 上边
  if (dir.includes('n')) {
    const maxDelta = resizeStart.value.height - props.minHeight
    const actualDelta = Math.min(deltaY, maxDelta)
    newHeight = resizeStart.value.height - actualDelta
    newY = resizeStart.value.posY + actualDelta
  }

  // 限制在窗口内
  newX = Math.max(0, newX)
  newY = Math.max(0, newY)
  newWidth = Math.min(newWidth, window.innerWidth - newX)
  newHeight = Math.min(newHeight, window.innerHeight - newY)

  size.value = { width: newWidth, height: newHeight }
  position.value = { x: newX, y: newY }
}

function stopResize() {
  isResizing.value = false
  document.removeEventListener('mousemove', onResize)
  document.removeEventListener('mouseup', stopResize)
}

// 最大化/还原
function toggleMaximize() {
  if (minimized.value) {
    minimized.value = false
    return
  }

  if (isMaximized.value) {
    // 还原
    position.value = { x: beforeMaximize.value.x, y: beforeMaximize.value.y }
    size.value = { width: beforeMaximize.value.width, height: beforeMaximize.value.height }
    isMaximized.value = false
  } else {
    // 最大化
    beforeMaximize.value = {
      x: position.value.x,
      y: position.value.y,
      width: size.value.width,
      height: size.value.height,
    }
    position.value = { x: 0, y: 0 }
    size.value = { width: window.innerWidth, height: window.innerHeight }
    isMaximized.value = true
  }
}

// 最小化
function minimize() {
  minimized.value = true
}

// 从最小化恢复
function restore() {
  minimized.value = false
}

// 处理遮罩层点击
function handleOverlayClick() {
  if (props.closeOnOverlay) {
    minimize()
  }
}

// 关闭
async function close() {
  // 如果提供了 beforeClose 回调，先执行它
  if (props.beforeClose) {
    const result = await props.beforeClose()
    if (!result) return // 如果返回 false，取消关闭
  }
  show.value = false
}

// 双击标题栏最大化/还原，手机固定全屏不响应
function onTitleDoubleClick() {
  if (themeStore.isMobile) return
  toggleMaximize()
}

// 监听显示状态
watch(show, (newShow) => {
  if (newShow) {
    minimized.value = false
    isMaximized.value = false
    initPosition()
  }
})

// 手机全屏期间不维护尺寸，回到桌面时最大化的铺满视口，其余按默认尺寸重新居中
// 调用方的默认尺寸可能也随视口计算，且它的 resize 监听在判定 isMobile 的监听之后才跑，
// 放到动画帧里取，此时本帧所有 resize 监听和由此触发的重渲染都已完成
watch(
  () => themeStore.isMobile,
  (mobile) => {
    if (mobile) return
    requestAnimationFrame(() => (isMaximized.value ? handleWindowResize() : initPosition()))
  },
)

// 浏览器缩放或旋转屏幕后把窗口收回视口内
function handleWindowResize() {
  if (themeStore.isMobile) return
  if (isMaximized.value) {
    size.value = { width: window.innerWidth, height: window.innerHeight }
    return
  }
  const width = Math.min(size.value.width, window.innerWidth)
  const height = Math.min(size.value.height, window.innerHeight)
  size.value = { width, height }
  position.value = {
    x: Math.min(position.value.x, window.innerWidth - width),
    y: Math.min(position.value.y, window.innerHeight - height),
  }
}

onMounted(() => {
  // 打开状态可能是持久化恢复的，不会经过 show 的 watch
  if (show.value) initPosition()
  window.addEventListener('resize', handleWindowResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleWindowResize)
})
</script>

<template>
  <Teleport to="body">
    <!-- 遮罩层 -->
    <Transition name="fade">
      <div
        v-if="show && !minimized"
        class="draggable-window-overlay"
        :style="{ zIndex: zIndex - 1 }"
        @click="handleOverlayClick"
      />
    </Transition>

    <!-- 主窗口：最小化时仅隐藏不销毁，保留编辑器的撤销历史、滚动位置和文件树状态 -->
    <Transition name="window">
      <div
        v-if="show"
        v-show="!minimized"
        ref="windowRef"
        class="draggable-window"
        :class="{ maximized: locked, dragging: isDragging, resizing: isResizing }"
        :style="windowStyle"
      >
        <!-- 标题栏 -->
        <div class="draggable-window-header" @mousedown="startDrag" @dblclick="onTitleDoubleClick">
          <span class="draggable-window-title">{{ title }}</span>
          <div class="draggable-window-controls">
            <button
              class="control-btn minimize"
              @click.stop="minimize"
              :title="$gettext('Minimize')"
            >
              <i-mdi-window-minimize />
            </button>
            <button
              v-if="!themeStore.isMobile"
              class="control-btn maximize"
              @click.stop="toggleMaximize"
              :title="isMaximized ? $gettext('Restore') : $gettext('Maximize')"
            >
              <i-mdi-window-restore v-if="isMaximized" />
              <i-mdi-window-maximize v-else />
            </button>
            <button class="control-btn close" @click.stop="close" :title="$gettext('Close')">
              <i-mdi-close />
            </button>
          </div>
        </div>

        <!-- 内容区域 -->
        <div class="draggable-window-content">
          <slot />
        </div>

        <!-- 调整大小的边框 -->
        <template v-if="!locked">
          <div class="resize-handle n" @mousedown="startResize($event, 'n')" />
          <div class="resize-handle s" @mousedown="startResize($event, 's')" />
          <div class="resize-handle e" @mousedown="startResize($event, 'e')" />
          <div class="resize-handle w" @mousedown="startResize($event, 'w')" />
          <div class="resize-handle ne" @mousedown="startResize($event, 'ne')" />
          <div class="resize-handle nw" @mousedown="startResize($event, 'nw')" />
          <div class="resize-handle se" @mousedown="startResize($event, 'se')" />
          <div class="resize-handle sw" @mousedown="startResize($event, 'sw')" />
        </template>
      </div>
    </Transition>
  </Teleport>

  <!-- 最小化后的图标 -->
  <Teleport :to="dock">
    <Transition name="minimize">
      <div
        v-if="show && minimized"
        class="draggable-window-minimized"
        :style="{
          background: themeVars.cardColor,
          color: themeVars.textColor1,
          '--border-radius': themeVars.borderRadius,
        }"
        @click="restore"
      >
        <slot name="icon">
          <i-mdi-file-document-outline />
        </slot>
        <span>{{ title }}</span>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped lang="scss">
.draggable-window-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
}

.draggable-window {
  position: fixed;
  display: flex;
  flex-direction: column;
  border-radius: var(--border-radius);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  overflow: hidden;

  &.maximized {
    border-radius: 0;
  }

  &.dragging,
  &.resizing {
    user-select: none;
  }
}

.draggable-window-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 40px;
  padding: 0 8px 0 16px;
  border-bottom: 1px solid var(--border-color);
  cursor: move;
  flex-shrink: 0;
  user-select: none;

  .maximized & {
    cursor: default;
  }
}

.draggable-window-title {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 500;
  font-size: 14px;
  color: var(--text-color-1);
}

.draggable-window-controls {
  display: flex;
  flex-shrink: 0;
  gap: 4px;
}

.control-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 28px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  color: var(--text-color-2);
  transition: all 0.2s;

  &:hover {
    background: var(--hover-color);
  }

  &.close:hover {
    background: #e81123;
    color: white;
  }
}

.draggable-window-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

// 调整大小的手柄
.resize-handle {
  position: absolute;

  &.n,
  &.s {
    left: 8px;
    right: 8px;
    height: 6px;
    cursor: ns-resize;
  }

  &.e,
  &.w {
    top: 8px;
    bottom: 8px;
    width: 6px;
    cursor: ew-resize;
  }

  &.n {
    top: -3px;
  }
  &.s {
    bottom: -3px;
  }
  &.e {
    right: -3px;
  }
  &.w {
    left: -3px;
  }

  &.ne,
  &.nw,
  &.se,
  &.sw {
    width: 12px;
    height: 12px;
  }

  &.ne {
    top: -3px;
    right: -3px;
    cursor: nesw-resize;
  }
  &.nw {
    top: -3px;
    left: -3px;
    cursor: nwse-resize;
  }
  &.se {
    bottom: -3px;
    right: -3px;
    cursor: nwse-resize;
  }
  &.sw {
    bottom: -3px;
    left: -3px;
    cursor: nesw-resize;
  }
}

// 最小化后的图标
.draggable-window-minimized {
  display: flex;
  align-items: center;
  gap: 8px;
  max-width: 100%;
  padding: 8px 16px;
  border-radius: var(--border-radius);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.25);
  }

  // 图标可能来自调用方插槽，带的是调用方的 scoped 标记
  :deep(svg) {
    flex-shrink: 0;
  }

  span {
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: 13px;
  }
}

// 动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.window-enter-active,
.window-leave-active {
  transition: all 0.2s;
}

.window-enter-from,
.window-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

.minimize-enter-active,
.minimize-leave-active {
  transition: all 0.2s;
}

.minimize-enter-from,
.minimize-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
