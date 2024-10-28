<script lang="ts" setup>
const { isFullscreen, toggle } = useFullscreen()
</script>

<template>
  <n-tooltip trigger="hover">
    <template #trigger>
      <n-icon mr-20 cursor-pointer size="20" @click="toggle">
        <icon-ant-design:fullscreen-exit-outlined v-if="isFullscreen" />
        <icon-ant-design:fullscreen-outlined v-else />
      </n-icon>
    </template>
    全屏显示
  </n-tooltip>
</template>
