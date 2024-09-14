import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/plugins/pureftpd/list', { params: { page, limit } }),
  // 添加
  add: (username: string, password: string, path: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/pureftpd/add', { username, password, path }),
  // 删除
  delete: (username: string): Promise<AxiosResponse<any>> =>
    request.delete('/plugins/pureftpd/delete', { params: { username } }),
  // 修改密码
  changePassword: (username: string, password: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/pureftpd/changePassword', { username, password }),
  // 获取端口
  port: (): Promise<AxiosResponse<any>> => request.get('/plugins/pureftpd/port'),
  // 修改端口
  setPort: (port: number): Promise<AxiosResponse<any>> =>
    request.post('/plugins/pureftpd/port', { port })
}
