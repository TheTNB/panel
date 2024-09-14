import { defineStore } from 'pinia'
import { toLogin } from '@/utils'
import { usePermissionStore, useTabStore } from '@/store'
import { resetRouter } from '@/router'
import user from '@/api/panel/user'

interface UserInfo {
  id?: string
  username?: string
  role?: Array<string>
}

export const useUserStore = defineStore('user', {
  state() {
    return {
      userInfo: <UserInfo>{}
    }
  },
  getters: {
    userId(): string {
      return this.userInfo.id || ''
    },
    username(): string {
      return this.userInfo.username || ''
    },
    role(): Array<string> {
      return this.userInfo.role || []
    }
  },
  actions: {
    async getUserInfo() {
      try {
        const res: any = await user.info()
        const { id, username, role } = res.data
        this.userInfo = { id, username, role }
        return Promise.resolve(res.data)
      } catch (error) {
        return Promise.reject(error)
      }
    },
    async logout() {
      user.logout().then(() => {
        const { resetTabs } = useTabStore()
        const { resetPermission } = usePermissionStore()
        resetPermission()
        resetTabs()
        resetRouter()
        this.$reset()
        toLogin()
      })
    },
    setUserInfo(userInfo = {}) {
      this.userInfo = { ...this.userInfo, ...userInfo }
    }
  }
})
