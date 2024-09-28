import type { Router } from 'vue-router'
import app from '@/api/panel/app'

export function createAppInstallGuard(router: Router) {
  router.beforeEach(async (to) => {
    const slug = to.path.split('/').pop()
    if (to.path.startsWith('/apps/') && slug) {
      await app.isInstalled(slug).then((res) => {
        if (!res.data.installed) {
          window.$message.error(`应用 ${res.data.name} 未安装`)
          return router.push({ name: 'app-index' })
        }
      })
    }

    // 网站
    if (to.path.startsWith('/website')) {
      await app.isInstalled('openresty').then((res) => {
        if (!res.data.installed) {
          window.$message.error(`Web 服务器 ${res.data.name} 未安装`)
          return router.push({ name: 'app-index' })
        }
      })
    }
    // 容器
    if (to.path.startsWith('/container')) {
      await app.isInstalled('podman').then((res) => {
        if (!res.data.installed) {
          window.$message.error(`容器引擎 ${res.data.name} 未安装`)
          return router.push({ name: 'app-index' })
        }
      })
    }
  })
}
