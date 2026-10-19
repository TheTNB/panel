import { useThemeStore } from '@/stores'

export type Breakpoint = 'mobile' | 'tablet' | 'desktop' | 'wide'

// 全站共用一个窗口宽度监听，AppProvider 据此写入 themeStore.isMobile
export const useWindowWidth = createSharedComposable(() => useWindowSize().width)

/**
 * 响应式断点 hook
 * - mobile  : <  768（与 uno 的 md 断点对齐）
 * - tablet  : 768 - 1023
 * - desktop : 1024 - 1439
 * - wide    : >= 1440
 */
export function useResponsive() {
  const themeStore = useThemeStore()
  const width = useWindowWidth()

  const isMobile = computed(() => themeStore.isMobile)

  const breakpoint = computed<Breakpoint>(() => {
    if (isMobile.value) return 'mobile'
    if (width.value < 1024) return 'tablet'
    if (width.value < 1440) return 'desktop'
    return 'wide'
  })

  const isTablet = computed(() => breakpoint.value === 'tablet')
  const isDesktop = computed(() => breakpoint.value === 'desktop' || breakpoint.value === 'wide')

  return { width, breakpoint, isMobile, isTablet, isDesktop }
}
