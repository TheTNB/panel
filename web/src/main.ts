import '@/styles/index.scss'
import '@/styles/reset.css'
import '@vue-js-cron/naive-ui/dist/naive-ui.css'
import 'uno.css'

import { createApp } from 'vue'
import App from './App.vue'

import { setupI18n } from '@/i18n/i18n'
import { setupRouter } from './router'
import { setupStore, useThemeStore } from './store'
import { setupNaiveDiscreteApi } from './utils'

import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor'

import dashboard from '@/api/panel/dashboard'
import CronNaivePlugin from '@vue-js-cron/naive-ui'

async function setupApp() {
  const app = createApp(App)
  app.use(VueMonacoEditorPlugin, {
    paths: {
      vs: '/assets/vs'
    },
    'vs/nls': {
      availableLanguages: { '*': 'zh-cn' }
    }
  })
  app.use(CronNaivePlugin)
  await setupStore(app)
  await setupNaiveDiscreteApi()
  await setupPanel().then(() => {
    setupI18n(app)
  })
  await setupRouter(app)
  app.mount('#app')
}

const title = ref(import.meta.env.VITE_APP_TITLE)

const setupPanel = async () => {
  const themeStore = useThemeStore()
  await dashboard
    .panel()
    .then((response) => response.json())
    .then((data) => {
      title.value = data.data.name || import.meta.env.VITE_APP_TITLE
      themeStore.setLocale(data.data.locale || 'zh_CN')
    })
    .catch((err) => {
      console.error(err)
    })
}

setupApp()

export { title }
