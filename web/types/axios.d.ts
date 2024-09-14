import type { InternalAxiosRequestConfig } from 'axios'

interface RequestConfig extends InternalAxiosRequestConfig {
  /** 接口是否需要错误提醒 */
  noNeedTip?: boolean
}

interface ErrorResolveResponse {
  code?: number | string
  message: string
  data?: any
}
