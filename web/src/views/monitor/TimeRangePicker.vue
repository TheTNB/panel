<script setup lang="ts">
import { useGettext } from 'vue3-gettext'

import { useResponsive } from '@/components/system/composables/useResponsive'

export type TimePreset = 'yesterday' | 'today' | 'week' | 'custom'

defineProps<{ preset: TimePreset }>()
const emit = defineEmits<{ change: [preset: TimePreset, customRange?: [number, number]] }>()

const { $gettext } = useGettext()
const { isMobile } = useResponsive()

const presets = computed<{ label: string; value: TimePreset }[]>(() => [
  { label: $gettext('Yesterday'), value: 'yesterday' },
  { label: $gettext('Today'), value: 'today' },
  { label: $gettext('Last 7 Days'), value: 'week' },
  { label: $gettext('Custom'), value: 'custom' },
])

const customShow = ref(false)
const tempRange = ref<[number, number] | null>(null)

// 下拉选中自定义时只弹日期面板，确认后才生效，下拉值随 preset 回显
const handleSelect = (value: TimePreset) => {
  if (value === 'custom') customShow.value = true
  else emit('change', value)
}

const confirmCustom = () => {
  if (tempRange.value) emit('change', 'custom', tempRange.value)
  customShow.value = false
}
</script>

<template>
  <n-popover
    v-if="isMobile"
    :show="customShow"
    trigger="manual"
    placement="bottom-end"
    :show-arrow="false"
    @clickoutside="customShow = false"
  >
    <template #trigger>
      <n-select
        size="small"
        class="!w-32"
        :value="preset"
        :options="presets"
        @update:value="handleSelect"
      />
    </template>
    <!-- 双月面板约 580px，手机上限宽并横向滚动，保证当前月完整可见 -->
    <div class="max-w-[calc(100vw-56px)] overflow-x-auto">
      <n-date-picker
        v-model:value="tempRange"
        type="daterange"
        panel
        :default-time="['00:00:00', '23:59:59']"
        :actions="['confirm']"
        @confirm="confirmCustom"
      />
    </div>
  </n-popover>
  <n-button-group v-else size="small">
    <n-button
      v-for="item in presets.slice(0, 3)"
      :key="item.value"
      :type="preset === item.value ? 'primary' : 'default'"
      @click="emit('change', item.value)"
    >
      {{ item.label }}
    </n-button>
    <n-popover v-model:show="customShow" trigger="click" placement="bottom-end" :show-arrow="false">
      <template #trigger>
        <n-button :type="preset === 'custom' ? 'primary' : 'default'">
          {{ $gettext('Custom') }}
        </n-button>
      </template>
      <n-date-picker
        v-model:value="tempRange"
        type="daterange"
        panel
        :default-time="['00:00:00', '23:59:59']"
        :actions="['confirm']"
        @confirm="confirmCustom"
      />
    </n-popover>
  </n-button-group>
</template>
