function addThemeColorCssVars() {
  const key = '__THEME_COLOR__'
  const defaultColor = '#66CCFF'
  const themeColor = window.localStorage.getItem(key) || defaultColor
  const cssVars = `--primary-color: ${themeColor}`
  document.documentElement.style.cssText = cssVars
}

addThemeColorCssVars()
