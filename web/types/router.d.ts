import { RouteRecordRaw } from 'vue-router'

interface Meta {
  title?: string
  icon?: string
  order?: number
  role?: Array<string>
  requireAuth?: boolean
}

interface RouteItem {
  name: string
  path: string
  redirect?: string
  isHidden?: boolean
  meta?: Meta
  children?: RoutesType
}

type RouteType = RouteRecordRaw & RouteItem

type RoutesType = Array<RouteType>

/** 前端导入的路由模块 */
type RouteModule = Record<string, { default: RouteType }>
