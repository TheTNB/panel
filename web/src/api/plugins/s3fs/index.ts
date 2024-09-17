import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/plugins/s3fs/list', { params: { page, limit } }),
  // 添加
  add: (data: any): Promise<AxiosResponse<any>> => request.post('/plugins/s3fs/add', data),
  // 删除
  delete: (id: number): Promise<AxiosResponse<any>> => request.post('/plugins/s3fs/delete', { id })
}
