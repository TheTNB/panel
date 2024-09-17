import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 保护列表
  jails: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/plugins/fail2ban/jails', { params: { page, limit } }),
  // 添加保护
  add: (data: any): Promise<AxiosResponse<any>> => request.post('/plugins/fail2ban/jails', data),
  // 删除保护
  delete: (name: string): Promise<AxiosResponse<any>> =>
    request.delete('/plugins/fail2ban/jails', { params: { name } }),
  // 封禁列表
  jail: (name: string): Promise<AxiosResponse<any>> =>
    request.get('/plugins/fail2ban/jails/' + name),
  // 解封 IP
  unban: (name: string, ip: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/fail2ban/unban', { name, ip }),
  // 获取白名单
  whitelist: (): Promise<AxiosResponse<any>> => request.get('/plugins/fail2ban/whiteList'),
  // 设置白名单
  setWhitelist: (ip: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/fail2ban/whiteList', { ip })
}
