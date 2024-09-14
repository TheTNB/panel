/** 侧边栏 */
interface Sider {
  width: number
  /** 折叠时的宽度 */
  collapsedWidth: number
  /** 是否折叠 */
  collapsed: boolean
}

/** 头部样式 */
interface Header {
  /** 是否显示 */
  visible: boolean
  /** 头部高度 */
  height: number
}

/** 标多页签样式 */
interface Tab {
  /** 是否显示 */
  visible: boolean
  /** 头部高度 */
  height: number
}

interface OtherColor {
  /** 信息 */
  info: string
  /** 成功 */
  success: string
  /** 警告 */
  warning: string
  /** 错误 */
  error: string
}

declare namespace Theme {
  interface Setting {
    isMobile: boolean
    darkMode: boolean
    sider: Sider
    header: Header
    tab: Tab
    /** 主题颜色 */
    primaryColor: string
    otherColor: OtherColor
    /** 语言 */
    language: string
  }
}
