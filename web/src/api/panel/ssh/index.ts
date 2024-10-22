import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取主机列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/ssh', { params: { page, limit } }),
  // 获取主机信息
  get: (id: number): Promise<AxiosResponse<any>> => request.get(`/ssh/${id}`),
  // 创建主机
  create: (req: any): Promise<AxiosResponse<any>> => request.post('/ssh', req),
  // 修改主机
  update: (id: number, req: any): Promise<AxiosResponse<any>> => request.put(`/ssh/${id}`, req),
  // 删除主机
  delete: (id: number): Promise<AxiosResponse<any>> => request.delete(`/ssh/${id}`)
}
