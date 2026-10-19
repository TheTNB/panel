import { h } from 'vue'

import { $gettext } from '@/utils/gettext'

export function notifyWebsiteCreated(name: string, warnings?: string[]) {
  if (!warnings?.length) {
    window.$message.success($gettext('Website %{ name } created successfully', { name }))
    return
  }
  window.$dialog.warning({
    title: $gettext('Website %{ name } created, but some steps failed', { name }),
    content: () => h('div', { style: { whiteSpace: 'pre-line' } }, warnings.join('\n')),
    maskClosable: false,
  })
}
