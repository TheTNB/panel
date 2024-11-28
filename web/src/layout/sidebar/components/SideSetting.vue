<script lang="ts" setup>
import TheIcon from '@/components/custom/TheIcon.vue'
import MenuCollapse from '@/layout/header/components/MenuCollapse.vue'
import { useThemeStore } from '@/store'

const themeStore = useThemeStore()

const settingModal = ref(false)
</script>

<template>
  <div h-40 flex justify-between px-20>
    <menu-collapse />
    <n-tooltip trigger="hover">
      <template #trigger>
        <the-icon
          v-show="!themeStore.sider.collapsed"
          :size="22"
          icon="material-symbols:settings"
          @click="settingModal = true"
        />
      </template>
      菜单设置
    </n-tooltip>
    <n-modal
      v-model:show="settingModal"
      preset="card"
      title="菜单设置"
      style="width: 60vw"
      size="huge"
      :bordered="false"
      :segmented="false"
      @close="settingModal = false"
      @mask-click="settingModal = false"
    >
      <n-form> </n-form>
    </n-modal>
  </div>
</template>
