import type { PluginOption } from 'vite'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'
import vueDevTools from 'vite-plugin-vue-devtools'

import { setupStaticCopyPlugin } from './copy'
import { setupHtmlPlugin } from './html'
import unplugins from './unplugin'

export function setupVitePlugins(viteEnv: ViteEnv): PluginOption[] {
  return [vue(), vueDevTools(), ...unplugins, unocss(), setupStaticCopyPlugin(), setupHtmlPlugin(viteEnv)]
}
