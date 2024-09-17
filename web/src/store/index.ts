import { createPinia } from 'pinia'
import type { App } from 'vue'

export async function setupStore(app: App) {
  app.use(createPinia())
}

export * from './modules'
