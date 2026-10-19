import type { GlobalThemeOverrides } from 'naive-ui'

import { breakpoint, darkSemantic, lightSemantic, typography, type SemanticPalette } from './tokens'

const buildOverrides = (sem: SemanticPalette): GlobalThemeOverrides => ({
  common: {
    fontFamily: typography.fontFamily.base,
    fontFamilyMono: typography.fontFamily.mono,
    primaryColor: sem.brand,
    primaryColorHover: sem.brandHover,
    primaryColorPressed: sem.brandPressed,
    primaryColorSuppl: sem.brandHover,
  },
})

export const lightThemeOverrides = buildOverrides(lightSemantic)
export const darkThemeOverrides = buildOverrides(darkSemantic)

// 手机端收窄卡片与对话框内边距，preset="card" 的弹窗也会合并全局 Card 覆盖
export const mobileThemeOverrides: GlobalThemeOverrides = {
  Card: {
    paddingSmall: '8px 12px',
    paddingMedium: '12px 16px',
    paddingLarge: '12px 16px',
    paddingHuge: '16px',
  },
  Dialog: {
    padding: '16px',
    // 关闭按钮是绝对定位，右边距要随 padding 一起收窄才能与内容右沿对齐
    closeMargin: '20px 14px 0 0',
  },
}

// naive 栅格断点与 uno 断点同源，naive 的键名是 s/m/l/xl/xxl
const px = (value: string) => Number.parseInt(value)
export const naiveBreakpoints = {
  xs: 0,
  s: px(breakpoint.sm),
  m: px(breakpoint.md),
  l: px(breakpoint.lg),
  xl: px(breakpoint.xl),
  xxl: px(breakpoint['2xl']),
}
