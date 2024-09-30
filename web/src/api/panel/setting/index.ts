import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取设置
  list: (): Promise<AxiosResponse<any>> => request.get('/setting'),
  // 保存设置
  update: (settings: any): Promise<AxiosResponse<any>> => request.post('/setting', settings),
  // 获取HTTPS设置
  getHttps: (): Promise<AxiosResponse<any>> => request.get('/setting/https'),
  // 保存HTTPS设置
  updateHttps: (https: any): Promise<AxiosResponse<any>> => request.post('/setting/https', https)
}
