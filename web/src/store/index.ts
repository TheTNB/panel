import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import type { App } from 'vue'

export async function setupStore(app: App) {
  const pinia = createPinia()
  pinia.use(piniaPluginPersistedstate)
  app.use(pinia)
}

export * from './modules'
