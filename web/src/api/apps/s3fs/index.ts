import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/s3fs/list', { params: { page, limit } }),
  // 添加
  add: (data: any): Promise<AxiosResponse<any>> => request.post('/apps/s3fs/add', data),
  // 删除
  delete: (id: number): Promise<AxiosResponse<any>> => request.post('/apps/s3fs/delete', { id })
}
