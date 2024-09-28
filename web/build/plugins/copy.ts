import { viteStaticCopy } from 'vite-plugin-static-copy'

export function setupStaticCopyPlugin() {
  return viteStaticCopy({
    targets: [
      {
        src: 'node_modules/monaco-editor/min/vs',
        dest: 'assets'
      }
    ]
  })
}
