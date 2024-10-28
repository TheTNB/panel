import type { Router } from 'vue-router'

import { trans } from '@/i18n/i18n'
import { useThemeStore } from '@/store'

export function createPageTitleGuard(router: Router) {
  const themeStore = useThemeStore()
  router.afterEach((to) => {
    const pageTitle = String(to.meta.title)
    if (pageTitle) document.title = `${trans(pageTitle)} | ${themeStore.name}`
    else document.title = themeStore.name
  })
}
