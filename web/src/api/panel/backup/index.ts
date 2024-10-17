import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取备份列表
  list: (type: string, page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get(`/backup/${type}`, { params: { page, limit } }),
  // 创建备份
  create: (type: string, target: string, path: string): Promise<AxiosResponse<any>> =>
    request.post(`/backup/${type}`, { target, path }),
  // 上传备份
  upload: (type: string, formData: FormData): Promise<AxiosResponse<any>> => {
    return request.post(`/backup/${type}/upload`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
  // 删除备份
  delete: (type: string, file: string): Promise<AxiosResponse<any>> =>
    request.delete(`/backup/${type}/delete`, { data: { file } }),
  // 恢复备份
  restore: (type: string, file: string, target: string): Promise<AxiosResponse<any>> =>
    request.post(`/backup/${type}/restore`, { file, target })
}
