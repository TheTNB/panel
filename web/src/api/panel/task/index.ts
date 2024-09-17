import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取状态
  status: (): Promise<AxiosResponse<any>> => request.get('/task/status'),
  // 获取任务列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/task/list', { params: { page, limit } }),
  // 获取任务日志
  log: (id: number): Promise<AxiosResponse<any>> =>
    request.get('/task/log', { params: { id } }),
  // 删除任务
  delete: (id: number): Promise<AxiosResponse<any>> => request.post('/task/delete', { id })
}
