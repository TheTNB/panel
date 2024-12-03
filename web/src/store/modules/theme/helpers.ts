import { addColorAlpha, getColorPalette } from '@/utils'
import type { GlobalThemeOverrides } from 'naive-ui'
import themeSetting from '~/settings/theme.json'

type ColorType = 'primary' | 'info' | 'success' | 'warning' | 'error'
type ColorScene = '' | 'Suppl' | 'Hover' | 'Pressed' | 'Active'
type ColorKey = `${ColorType}Color${ColorScene}`
type ThemeColor = Partial<Record<ColorKey, string>>

interface ColorAction {
  scene: ColorScene
  handler: (color: string) => string
}

/** 初始化主题配置 */
export function defaultSettings(): Theme.Setting {
  const isMobile = themeSetting.isMobile || false
  const darkMode = themeSetting.darkMode || false
  const sider = themeSetting.sider || {
    width: 160,
    collapsedWidth: 64,
    collapsed: false
  }
  const header = themeSetting.header || { visible: true, height: 60 }
  const tab = themeSetting.tab || { visible: true, height: 50 }
  const primaryColor = themeSetting.primaryColor || '#00BFFF'
  const otherColor = themeSetting.otherColor || {
    info: '#0099ad',
    success: '#52c41a',
    warning: '#faad14',
    error: '#f5222d'
  }
  const locale = themeSetting.locale || 'zh_CN'
  const name = themeSetting.name || import.meta.env.VITE_APP_TITLE
  const logo = ''
  return { isMobile, darkMode, sider, header, tab, primaryColor, otherColor, locale, name, logo }
}

/** 获取naive的主题颜色 */
export function getNaiveThemeOverrides(colors: Record<ColorType, string>): GlobalThemeOverrides {
  const { primary, info, success, warning, error } = colors

  const themeColors = getThemeColors([
    ['primary', primary],
    ['info', info],
    ['success', success],
    ['warning', warning],
    ['error', error]
  ])

  const colorLoading = primary

  return {
    common: {
      ...themeColors
    },
    LoadingBar: {
      colorLoading
    }
  }
}

/** 获取主题颜色的各种场景对应的颜色 */
function getThemeColors(colors: [ColorType, string][]) {
  const colorActions: ColorAction[] = [
    { scene: '', handler: (color) => color },
    { scene: 'Suppl', handler: (color) => color },
    { scene: 'Hover', handler: (color) => getColorPalette(color, 5) },
    { scene: 'Pressed', handler: (color) => getColorPalette(color, 7) },
    { scene: 'Active', handler: (color) => addColorAlpha(color, 0.1) }
  ]

  const themeColor: ThemeColor = {}

  colors.forEach((color) => {
    colorActions.forEach((action) => {
      const [colorType, colorValue] = color
      const colorKey: ColorKey = `${colorType}Color${action.scene}`
      themeColor[colorKey] = action.handler(colorValue)
    })
  })

  return themeColor
}
