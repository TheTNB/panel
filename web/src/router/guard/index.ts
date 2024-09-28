import type { Router } from 'vue-router'
import { createPageLoadingGuard } from './page-loading-guard'
import { createPageTitleGuard } from './page-title-guard'
import { createAppInstallGuard } from './app-install-guard'

export function setupRouterGuard(router: Router) {
  createPageLoadingGuard(router)
  createPageTitleGuard(router)
  createAppInstallGuard(router)
}
