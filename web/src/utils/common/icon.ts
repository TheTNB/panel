import { h } from 'vue'
import { Icon } from '@iconify/vue'
import { NIcon } from 'naive-ui'

interface Props {
  size?: number
  color?: string
  class?: string
}

export function renderIcon(icon: string, props: Props = { size: 12 }) {
  return () => h(NIcon, props, { default: () => h(Icon, { icon }) })
}
