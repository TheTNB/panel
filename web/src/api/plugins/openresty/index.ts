import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/plugins/openresty/load'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/plugins/openresty/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/openresty/config', { config }),
  // 获取错误日志
  errorLog: (): Promise<AxiosResponse<any>> => request.get('/plugins/openresty/errorLog'),
  // 清空错误日志
  clearErrorLog: (): Promise<AxiosResponse<any>> => request.post('/plugins/openresty/clearErrorLog')
}
