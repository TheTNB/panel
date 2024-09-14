import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取防火墙状态
  firewallStatus: (): Promise<AxiosResponse<any>> => request.get('/panel/safe/firewallStatus'),
  // 设置防火墙状态
  setFirewallStatus: (status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/panel/safe/firewallStatus', { status }),
  // 获取防火墙规则
  firewallRules: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/safe/firewallRules', { params: { page, limit } }),
  // 添加防火墙规则
  addFirewallRule: (port: string, protocol: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/safe/firewallRules', { port, protocol }),
  // 删除防火墙规则
  deleteFirewallRule: (port: string, protocol: string): Promise<AxiosResponse<any>> =>
    request.delete('/panel/safe/firewallRules', { data: { port, protocol } }),
  // 获取SSH状态
  sshStatus: (): Promise<AxiosResponse<any>> => request.get('/panel/safe/sshStatus'),
  // 设置SSH状态
  setSshStatus: (status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/panel/safe/sshStatus', { status }),
  // 获取SSH端口
  sshPort: (): Promise<AxiosResponse<any>> => request.get('/panel/safe/sshPort'),
  // 设置SSH端口
  setSshPort: (port: number): Promise<AxiosResponse<any>> =>
    request.post('/panel/safe/sshPort', { port }),
  // 获取Ping状态
  pingStatus: (): Promise<AxiosResponse<any>> => request.get('/panel/safe/pingStatus'),
  // 设置Ping状态
  setPingStatus: (status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/panel/safe/pingStatus', { status })
}
