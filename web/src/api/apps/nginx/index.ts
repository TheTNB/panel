import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/apps/nginx/load'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/apps/nginx/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/nginx/config', { config }),
  // 获取错误日志
  errorLog: (): Promise<AxiosResponse<any>> => request.get('/apps/nginx/errorLog'),
  // 清空错误日志
  clearErrorLog: (): Promise<AxiosResponse<any>> => request.post('/apps/nginx/clearErrorLog')
}
