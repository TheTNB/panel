import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取防火墙状态
  status: (): Promise<AxiosResponse<any>> => request.get('/firewall/status'),
  // 设置防火墙状态
  updateStatus: (status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/firewall/status', { status }),
  // 获取防火墙规则
  rules: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/firewall/rule', { params: { page, limit } }),
  // 创建防火墙规则
  createRule: (rule: any): Promise<AxiosResponse<any>> => request.post('/firewall/rule', rule),
  // 删除防火墙规则
  deleteRule: (rule: any): Promise<AxiosResponse<any>> =>
    request.delete('/firewall/rule', { data: rule }),
  // 获取防火墙IP规则
  ipRules: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/firewall/ipRule', { params: { page, limit } }),
  // 创建防火墙IP规则
  createIpRule: (rule: any): Promise<AxiosResponse<any>> => request.post('/firewall/ipRule', rule),
  // 删除防火墙IP规则
  deleteIpRule: (rule: any): Promise<AxiosResponse<any>> =>
    request.delete('/firewall/ipRule', { data: rule }),
  // 获取防火墙转发规则
  forwards: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/firewall/forward', { params: { page, limit } }),
  // 创建防火墙转发规则
  createForward: (rule: any): Promise<AxiosResponse<any>> =>
    request.post('/firewall/forward', rule),
  // 删除防火墙转发规则
  deleteForward: (rule: any): Promise<AxiosResponse<any>> =>
    request.delete('/firewall/forward', { data: rule })
}
