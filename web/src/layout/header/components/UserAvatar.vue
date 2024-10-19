<script lang="ts" setup>
import user from '@/api/panel/user'
import { router } from '@/router'
import { useUserStore } from '@/store'
import { renderIcon } from '@/utils'

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

const handleSelect = (key: string) => {
  if (key === 'logout') {
    window.$dialog.info({
      content: '确认退出？',
      title: '提示',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick() {
        user.logout().then(() => {
          userStore.logout()
        })
        window.$message.success('已退出登录!')
      }
    })
  }
  if (key === 'changePassword') {
    router.push({ name: 'setting-index' })
  }
}

const username = computed(() => {
  if (userStore.username !== '') {
    return userStore.username
  }
  return '未知'
})
</script>

<template>
  <n-dropdown :options="options" @select="handleSelect">
    <div flex cursor-pointer items-center>
      <span>{{ username }}</span>
    </div>
  </n-dropdown>
</template>
