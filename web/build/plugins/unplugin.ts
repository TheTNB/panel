import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'

/**
 * * unplugin-icons插件，自动引入iconify图标
 * usage: https://github.com/antfu/unplugin-icons
 * 图标库: https://icones.js.org/
 */
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'

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
      IconsResolver({ customCollections: ['custom'], componentPrefix: 'icon' })
    ],
    dts: 'types/components.d.ts'
  })
]
