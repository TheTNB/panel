import type { PaginationProps } from 'naive-ui'

import { useThemeStore } from '@/stores'

/**
 * 表格分页配置：手机用简洁模式，去掉每页条数与跳转
 */
export function useTablePagination(
  page: MaybeRefOrGetter<number>,
  pageSize: MaybeRefOrGetter<number>,
  total: MaybeRefOrGetter<number | undefined>,
  pageSizes = [20, 50, 100, 200],
) {
  const themeStore = useThemeStore()

  return computed<PaginationProps>(() => ({
    page: toValue(page),
    pageSize: toValue(pageSize),
    itemCount: toValue(total) ?? 0,
    pageSizes,
    simple: themeStore.isMobile,
    showQuickJumper: !themeStore.isMobile,
    showSizePicker: !themeStore.isMobile,
  }))
}
