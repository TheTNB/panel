import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取设置
  list: (): Promise<AxiosResponse<any>> => request.get('/panel/setting/list'),
  // 保存设置
  update: (settings: any): Promise<AxiosResponse<any>> =>
    request.post('/panel/setting/update', settings),
  // 获取HTTPS设置
  getHttps: (): Promise<AxiosResponse<any>> => request.get('/panel/setting/https'),
  // 保存HTTPS设置
  updateHttps: (https: any): Promise<AxiosResponse<any>> =>
    request.post('/panel/setting/https', https)
}
