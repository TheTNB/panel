import { useThemeStore } from '@/store'
import mitt from 'mitt'
import * as NaiveUI from 'naive-ui'

export async function setupNaiveDiscreteApi() {
  const themeStore = useThemeStore()
  const configProviderProps = computed(() => ({
    theme: themeStore.naiveTheme,
    themeOverrides: themeStore.naiveThemeOverrides
  }))
  const { message, dialog, notification, loadingBar } = NaiveUI.createDiscreteApi(
    ['message', 'dialog', 'notification', 'loadingBar'],
    { configProviderProps }
  )

  window.$loadingBar = loadingBar
  window.$notification = notification
  window.$message = message
  window.$dialog = dialog
  window.$bus = mitt()
}
