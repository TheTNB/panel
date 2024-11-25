import { createHtmlPlugin } from 'vite-plugin-html'

export function setupHtmlPlugin(viteEnv: ViteEnv) {
  const { VITE_APP_TITLE } = viteEnv
  return createHtmlPlugin({
    minify: true,
    inject: {
      data: {
        title: VITE_APP_TITLE
      }
    }
  })
}
