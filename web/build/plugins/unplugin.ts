import AutoImport from 'unplugin-auto-import/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'

/**
 * * unplugin-icons应用，自动引入iconify图标
 * usage: https://github.com/antfu/unplugin-icons
 * 图标库: https://icones.js.org/
 */
import IconsResolver from 'unplugin-icons/resolver'
import Icons from 'unplugin-icons/vite'

export default [
  AutoImport({
    imports: ['vue', 'vue-router'],
    dts: 'types/auto-imports.d.ts',
    eslintrc: {
      enabled: true
    }
  }),
  Icons({
    compiler: 'vue3',
    scale: 1,
    defaultClass: 'inline-block'
  }),
  Components({
    resolvers: [
      NaiveUiResolver(),
      IconsResolver({ customCollections: ['custom'], prefix: 'icon' })
    ],
    dts: 'types/components.d.ts'
  })
]
