import type { FormProps } from 'naive-ui'

import { useThemeStore } from '@/stores'

type FormLayout = Pick<FormProps, 'labelPlacement' | 'labelWidth' | 'inline'>

/**
 * 表单标签布局：桌面左侧标签，手机顶部标签
 * inline 只在传 true 时才写入，避免盖掉模板上直接写的 inline 属性
 */
export function useFormLayout(labelWidth: number | 'auto' = 120, inline = false) {
  const themeStore = useThemeStore()

  return computed<FormLayout>(() => {
    const layout: FormLayout = themeStore.isMobile
      ? { labelPlacement: 'top' }
      : { labelPlacement: 'left', labelWidth }
    return inline ? { ...layout, inline: !themeStore.isMobile } : layout
  })
}
