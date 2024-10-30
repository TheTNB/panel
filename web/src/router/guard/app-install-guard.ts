import app from '@/api/panel/app'
import type { Router } from 'vue-router'

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
      await app.isInstalled('nginx').then((res) => {
        if (!res.data.installed) {
          window.$message.error(`Web 服务器 ${res.data.name} 未安装`)
          return router.push({ name: 'app-index' })
        }
      })
    }
    // 容器
    if (to.path.startsWith('/container')) {
      let flag = false
      await app.isInstalled('docker').then((res) => {
        if (res.data.installed) {
          flag = true
        }
      })
      if (!flag) {
        await app.isInstalled('podman').then((res) => {
          if (res.data.installed) {
            flag = true
          }
        })
      }
      if (!flag) {
        window.$message.error(`容器引擎 Docker / Podman 未安装`)
        return router.push({ name: 'app-index' })
      }
    }
  })
}
