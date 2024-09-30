import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 开关
  switch: (monitor: boolean): Promise<AxiosResponse<any>> =>
    request.post('/monitor/switch', { monitor }),
  // 保存天数
  saveDays: (days: number): Promise<AxiosResponse<any>> =>
    request.post('/monitor/saveDays', { days }),
  // 清空监控记录
  clear: (): Promise<AxiosResponse<any>> => request.post('/monitor/clear'),
  // 监控记录
  list: (start: number, end: number): Promise<AxiosResponse<any>> =>
    request.get('/monitor/list', { params: { start, end } }),
  // 开关和天数
  switchAndDays: (): Promise<AxiosResponse<any>> => request.get('/monitor/switchAndDays')
}
