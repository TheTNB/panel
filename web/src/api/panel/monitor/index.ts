import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 开关
  setting: (): Promise<AxiosResponse<any>> => request.get('/monitor/setting'),
  // 保存天数
  updateSetting: (enabled: boolean, days: number): Promise<AxiosResponse<any>> =>
    request.post('/monitor/setting', { enabled, days }),
  // 清空监控记录
  clear: (): Promise<AxiosResponse<any>> => request.post('/monitor/clear'),
  // 监控记录
  list: (start: number, end: number): Promise<AxiosResponse<any>> =>
    request.get('/monitor/list', { params: { start, end } })
}
