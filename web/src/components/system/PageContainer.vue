<script lang="ts" setup>
import AppFooter from '@/components/common/AppFooter.vue'

import LoadingSkeleton from './LoadingSkeleton.vue'
import PageHeader from './PageHeader.vue'

interface Props {
  title?: string
  description?: string
  showFooter?: boolean
  flex?: boolean
  loading?: boolean
  bare?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: undefined,
  description: undefined,
  showFooter: false,
  flex: false,
  loading: false,
  bare: false,
})

const slots = defineSlots<{
  header?: () => any
  tabs?: () => any
  actions?: () => any
  default?: () => any
  footer?: () => any
}>()

const hasHeader = computed(() => !!(props.title || slots.header || slots.actions))

// 手机屏幕小，带页脚的页面让内容区至少撑满一屏，页脚滚动后才出现
// 用 flex-basis 而不是 min-height：显式 min-height 会取消 flex 子项按内容算的最小高度，内容会溢出卡片
const bodyClass = computed(() => {
  if (props.flex) return 'flex-1 min-h-0'
  return props.showFooter || slots.footer ? 'flex-1 lt-md:basis-full lt-md:shrink-0' : 'flex-1'
})
</script>

<template>
  <transition appear mode="out-in" name="fade-slide">
    <section class="cus-scroll-y page-container">
      <header v-if="hasHeader || slots.tabs">
        <slot v-if="slots.header" name="header" />
        <PageHeader v-else-if="title" :title="title" :description="description">
          <template v-if="slots.actions" #actions>
            <slot name="actions" />
          </template>
        </PageHeader>
        <n-card v-if="slots.tabs" size="small">
          <slot name="tabs" />
        </n-card>
      </header>

      <div v-if="bare" class="flex flex-col gap-1" :class="bodyClass">
        <LoadingSkeleton v-if="loading" type="card" />
        <slot v-else />
      </div>
      <n-card
        v-else
        :class="bodyClass"
        :content-class="flex ? 'flex flex-col min-h-0 h-full' : undefined"
      >
        <LoadingSkeleton v-if="loading" type="card" />
        <slot v-else />
      </n-card>

      <slot v-if="slots.footer" name="footer" />
      <!-- flex 页面内容区固定占满视口、没有页面滚动，手机上页脚无处可滚到，直接隐藏 -->
      <AppFooter v-else-if="showFooter" class="mt-auto pt-5" :class="{ 'lt-md:hidden': flex }" />
    </section>
  </transition>
</template>
