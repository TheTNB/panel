<script setup lang="ts">
import { useThemeStore } from '@/stores'

type Size = 'sm' | 'md' | 'lg' | 'full'

const props = withDefaults(defineProps<{ size?: Size }>(), { size: 'md' })

const themeStore = useThemeStore()

const widths: Record<Size, { width: string; maxWidth?: string }> = {
  sm: { width: '640px' },
  md: { width: '60vw', maxWidth: '960px' },
  lg: { width: '70vw', maxWidth: '1080px' },
  full: { width: '90vw' },
}

// n-card 自带 width: 100%，用行内 style 才能稳定覆盖；调用方的 style（如高度）由 Vue 自动合并
const widthStyle = computed(() =>
  themeStore.isMobile ? { width: 'calc(100vw - 24px)' } : widths[props.size],
)
</script>

<template>
  <n-modal preset="card" size="huge" :bordered="false" :segmented="false" :style="widthStyle">
    <template v-for="(_, name) in $slots" #[name]="slotProps">
      <slot :name="name" v-bind="slotProps ?? {}" />
    </template>
  </n-modal>
</template>
