<script lang="ts" setup>
import { useUserStore } from '@/store'
import { renderIcon } from '@/utils'
import { router } from '@/router'

const userStore = useUserStore()

const options = [
  {
    label: '修改密码',
    key: 'changePassword',
    icon: renderIcon('mdi:key', { size: 14 })
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon('mdi:exit-to-app', { size: 14 })
  }
]

function handleSelect(key: string) {
  if (key === 'logout') {
    window.$dialog.info({
      content: '确认退出？',
      title: '提示',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick() {
        userStore.logout()
        window.$message.success('已退出登录!')
      }
    })
  }
  if (key === 'changePassword') {
    router.push({ name: 'setting-index' })
  }
}
</script>

<template>
  <n-dropdown :options="options" @select="handleSelect">
    <div flex cursor-pointer items-center>
      <span hidden sm:block>{{ userStore.username }}</span>
    </div>
  </n-dropdown>
</template>
