import type { PluginOption } from 'vite'
import vue from '@vitejs/plugin-vue'
import unocss from 'unocss/vite'

import unplugins from './unplugin'
import { setupHtmlPlugin } from './html'
import { setupStaticCopyPlugin } from './copy'

export function setupVitePlugins(viteEnv: ViteEnv): PluginOption[] {
  return [vue(), ...unplugins, unocss(), setupStaticCopyPlugin(), setupHtmlPlugin(viteEnv)]
}
