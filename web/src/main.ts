import '@/styles/reset.css'
import '@/styles/index.scss'
import 'uno.css'

import { createApp } from 'vue'
import App from './App.vue'
import { setupStore, useThemeStore } from './store'
import { setupRouter } from './router'
import { setupI18n } from '@/i18n/i18n'
import { setupNaiveDiscreteApi } from './utils'
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor'
import info from '@/api/panel/info'

async function setupApp() {
  const app = createApp(App)
  app.use(VueMonacoEditorPlugin, {
    paths: {
      vs: 'https://cdnjs.admincdn.com/monaco-editor/0.48.0/min/vs'
    },
    'vs/nls': {
      availableLanguages: { '*': 'zh-cn' }
    }
  })
  await setupStore(app)
  await setupNaiveDiscreteApi()
  await setupPanel().then(() => {
    setupI18n(app)
  })
  await setupRouter(app)
  app.mount('#app')
}

const title = ref('')

const setupPanel = async () => {
  const themeStore = useThemeStore()
  await info
    .panel()
    .then((response) => response.json())
    .then((data) => {
      title.value = data.data.name || import.meta.env.VITE_APP_TITLE
      themeStore.setLanguage(data.data.language || 'zh_CN')
    })
    .catch((err) => {
      console.error(err)
    })
}

setupApp()

export { title }
