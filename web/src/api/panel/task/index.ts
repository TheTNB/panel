import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取状态
  status: (): Promise<AxiosResponse<any>> => request.get('/panel/task/status'),
  // 获取任务列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/task/list', { params: { page, limit } }),
  // 获取任务日志
  log: (id: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/task/log', { params: { id } }),
  // 删除任务
  delete: (id: number): Promise<AxiosResponse<any>> => request.post('/panel/task/delete', { id })
}
