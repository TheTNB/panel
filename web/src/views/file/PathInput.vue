<script setup lang="ts">
import type { InputInst } from 'naive-ui'
import EventBus from '@/utils/event'
import { onUnmounted } from 'vue'
import { checkPath } from '@/utils/file'

const path = defineModel<string>('path', { type: String, required: true })
const isInput = ref(false)
const pathInput = ref<InputInst | null>(null)
const input = ref('www')

const history: string[] = []
let current = -1

const handleInput = () => {
  isInput.value = true
  nextTick(() => {
    pathInput.value?.focus()
  })
}

const handleBlur = () => {
  if (!checkPath(input.value)) {
    window.$message.error('路径不合法')
    return
  }

  isInput.value = false
  path.value = '/' + input.value
  handlePushHistory(path.value)
}

const handleRefresh = () => {
  EventBus.emit('file:refresh')
}

const handleUp = () => {
  const count = splitPath(path.value, '/').length
  setPath(count - 2)
}

const handleBack = () => {
  if (current > 0) {
    current--
    path.value = history[current]
    input.value = path.value.slice(1)
  }
}

const handleForward = () => {
  if (current < history.length - 1) {
    current++
    path.value = history[current]
    input.value = path.value.slice(1)
  }
}

const splitPath = (str: string, delimiter: string) => {
  if (str === delimiter || str === '') {
    return []
  }
  return str.split(delimiter).slice(1)
}

const setPath = (index: number) => {
  const newPath = splitPath(path.value, '/')
    .slice(0, index + 1)
    .join('/')
  path.value = '/' + newPath
  input.value = newPath
  handlePushHistory(path.value)
}

const handlePushHistory = (path: string) => {
  // 防止在前进后退时重复添加
  if (current != history.length - 1) {
    return
  }

  history.splice(current + 1)
  history.push(path)
  current = history.length - 1
}

const handleSearch = () => {
  window.$message.info('搜索功能暂未实现')
}

watch(
  path,
  (value) => {
    input.value = value.slice(1)
  },
  { immediate: true }
)

onMounted(() => {
  EventBus.on('push-history', handlePushHistory)
})

onUnmounted(() => {
  EventBus.off('push-history', handlePushHistory)
})
</script>

<template>
  <n-flex>
    <n-button @click="handleBack">
      <icon-bi-arrow-left />
    </n-button>
    <n-button @click="handleForward">
      <icon-bi-arrow-right />
    </n-button>
    <n-button @click="handleUp">
      <icon-bi-arrow-up />
    </n-button>
    <n-button @click="handleRefresh">
      <icon-bi-arrow-clockwise />
    </n-button>
    <n-input-group flex-1>
      <n-tag size="large" v-if="!isInput" flex-1 @click="handleInput">
        <n-breadcrumb separator=">">
          <n-breadcrumb-item @click.stop="setPath(-1)"> 根目录 </n-breadcrumb-item>
          <n-breadcrumb-item
            v-for="(item, index) in splitPath(path, '/')"
            :key="index"
            @click.stop="setPath(index)"
          >
            {{ item }}
          </n-breadcrumb-item>
        </n-breadcrumb>
      </n-tag>
      <n-input-group-label v-if="isInput">/</n-input-group-label>
      <n-input
        ref="pathInput"
        v-model:value="input"
        v-if="isInput"
        @keyup.enter="handleBlur"
        @blur="handleBlur"
      />
    </n-input-group>
    <n-input-group w-400>
      <n-input placeholder="请输入搜索内容">
        <template #suffix>
          <n-checkbox> 包含子目录 </n-checkbox>
        </template>
      </n-input>
      <n-button type="primary" @click="handleSearch">
        <icon-bi-search />
      </n-button>
    </n-input-group>
  </n-flex>
</template>

<style scoped lang="scss"></style>
