import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取SSH
  ssh: (): Promise<AxiosResponse<any>> => request.get('/safe/ssh'),
  // 设置SSH
  setSsh: (status: boolean, port: number): Promise<AxiosResponse<any>> =>
    request.post('/safe/ssh', { status, port }),
  // 获取Ping状态
  pingStatus: (): Promise<AxiosResponse<any>> => request.get('/safe/ping'),
  // 设置Ping状态
  setPingStatus: (status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/safe/ping', { status })
}
