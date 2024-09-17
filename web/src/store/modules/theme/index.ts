import {
  dateEnUS,
  dateZhCN,
  enUS,
  type GlobalThemeOverrides,
  type NDateLocale,
  type NLocale,
  zhCN
} from 'naive-ui'
import { darkTheme } from 'naive-ui'
import type { BuiltInGlobalTheme } from 'naive-ui/es/themes/interface'
import { defineStore } from 'pinia'
import { getNaiveThemeOverrides, initThemeSettings } from './helpers'

type ThemeState = Theme.Setting

const locales: Record<string, { locale: NLocale; dateLocale: NDateLocale }> = {
  zh_CN: { locale: zhCN, dateLocale: dateZhCN },
  en: { locale: enUS, dateLocale: dateEnUS }
}

export const useThemeStore = defineStore('theme-store', {
  state: (): ThemeState => initThemeSettings(),
  getters: {
    naiveThemeOverrides(): GlobalThemeOverrides {
      return getNaiveThemeOverrides({
        primary: this.primaryColor,
        ...this.otherColor
      })
    },
    naiveTheme(): BuiltInGlobalTheme | undefined {
      return this.darkMode ? darkTheme : undefined
    },
    naiveLocale(): NLocale {
      return locales[this.language].locale
    },
    naiveDateLocale(): NDateLocale {
      return locales[this.language].dateLocale
    }
  },
  actions: {
    setIsMobile(isMobile: boolean) {
      this.isMobile = isMobile
    },
    /** 设置暗黑模式 */
    setDarkMode(darkMode: boolean) {
      this.darkMode = darkMode
    },
    /** 切换/关闭 暗黑模式 */
    toggleDarkMode() {
      this.darkMode = !this.darkMode
    },
    /** 切换/关闭 折叠侧边栏 */
    toggleCollapsed() {
      this.sider.collapsed = !this.sider.collapsed
    },
    /** 设置 折叠侧边栏 */
    setCollapsed(collapsed: boolean) {
      this.sider.collapsed = collapsed
    },
    /** 设置主题色 */
    setPrimaryColor(color: string) {
      this.primaryColor = color
    },
    /** 设置语言 */
    setLanguage(language: string) {
      this.language = language
    }
  }
})
