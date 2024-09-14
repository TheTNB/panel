import type { Router } from 'vue-router'
import plugin from '@/api/panel/plugin'

export function createPluginInstallGuard(router: Router) {
  router.beforeEach(async (to) => {
    const slug = to.path.split('/').pop()
    if (to.path.startsWith('/plugins/') && slug) {
      await plugin.isInstalled(slug).then((res) => {
        if (!res.data.installed) {
          window.$message.error(`插件 ${res.data.name} 未安装`)
          return router.push({ name: 'plugin-index' })
        }
      })
    }

    // 网站
    if (to.path.startsWith('/website')) {
      await plugin.isInstalled('openresty').then((res) => {
        if (!res.data.installed) {
          window.$message.error(`Web 服务器 ${res.data.name} 未安装`)
          return router.push({ name: 'plugin-index' })
        }
      })
    }
    // 容器
    if (to.path.startsWith('/container')) {
      await plugin.isInstalled('podman').then((res) => {
        if (!res.data.installed) {
          window.$message.error(`容器引擎 ${res.data.name} 未安装`)
          return router.push({ name: 'plugin-index' })
        }
      })
    }
  })
}
