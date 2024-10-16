import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取防火墙状态
  firewallStatus: (): Promise<AxiosResponse<any>> => request.get('/firewall/status'),
  // 设置防火墙状态
  setFirewallStatus: (status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/firewall/status', { status }),
  // 获取防火墙规则
  firewallRules: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/firewall/rule', { params: { page, limit } }),
  // 添加防火墙规则
  addFirewallRule: (port: number, protocol: string): Promise<AxiosResponse<any>> =>
    request.post('/firewall/rule', { port, protocol }),
  // 删除防火墙规则
  deleteFirewallRule: (port: number, protocol: string): Promise<AxiosResponse<any>> =>
    request.delete('/firewall/rule', { params: { port, protocol } }),
  // 获取SSH
  ssh: (): Promise<AxiosResponse<any>> => request.get('/safe/ssh'),
  // 设置SSH
  setSsh: (status: boolean, port: number): Promise<AxiosResponse<any>> =>
    request.post('/safe/ssh', { status, port }),
  // 获取Ping状态
  pingStatus: (): Promise<AxiosResponse<any>> => request.get('/safe/ping'),
  // 设置Ping状态
  setPingStatus: (status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/safe/ping', { status })
}
