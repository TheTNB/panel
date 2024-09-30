import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取信息
  info: (): Promise<AxiosResponse<any>> => request.get('/ssh/info'),
  // 保存信息
  saveInfo: (
    host: string,
    port: number,
    user: string,
    password: string
  ): Promise<AxiosResponse<any>> => request.post('/ssh/info', { host, port, user, password })
}
