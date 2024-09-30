import { Icon } from '@iconify/vue'
import { NIcon } from 'naive-ui'
import { h } from 'vue'

interface Props {
  size?: number
  color?: string
  class?: string
}

export function renderIcon(icon: string, props: Props = { size: 12 }) {
  return () => h(NIcon, props, { default: () => h(Icon, { icon }) })
}
