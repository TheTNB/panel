import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取任务列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/cron/list', { params: { page, limit } }),
  // 获取任务脚本
  script: (id: number): Promise<AxiosResponse<any>> => request.get('/panel/cron/' + id),
  // 添加任务
  add: (task: any): Promise<AxiosResponse<any>> => request.post('/panel/cron/add', task),
  // 修改任务
  update: (id: number, name: string, time: string, script: string): Promise<AxiosResponse<any>> =>
    request.put('/panel/cron/' + id, { name, time, script }),
  // 删除任务
  delete: (id: number): Promise<AxiosResponse<any>> => request.delete('/panel/cron/' + id),
  // 获取任务日志
  log: (id: number): Promise<AxiosResponse<any>> => request.get('/panel/cron/log/' + id),
  // 修改任务状态
  status: (id: number, status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/panel/cron/status', { id, status })
}
