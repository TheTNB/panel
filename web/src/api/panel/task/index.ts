import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取状态
  status: (): Promise<AxiosResponse<any>> => request.get('/task/status'),
  // 获取任务列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/task', { params: { page, limit } }),
  // 获取任务
  get: (id: number): Promise<AxiosResponse<any>> => request.get('/task/' + id),
  // 删除任务
  delete: (id: number): Promise<AxiosResponse<any>> => request.delete('/task/' + id)
}
