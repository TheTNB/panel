import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // DNS
  dns: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/dns'),
  // 设置 DNS
  updateDns: (dns1: string, dns2: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/dns', { dns1, dns2 }),
  // SWAP
  swap: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/swap'),
  // 设置 SWAP
  updateSwap: (size: number): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/swap', { size }),
  // 时区
  timezone: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/timezone'),
  // 设置时区
  updateTimezone: (timezone: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/timezone', { timezone }),
  // 设置时间
  updateTime: (time: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/time', { time }),
  // 同步时间
  syncTime: (): Promise<AxiosResponse<any>> => request.post('/apps/toolbox/syncTime'),
  // 主机名
  hostname: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/hostname'),
  // Hosts
  hosts: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/hosts'),
  // 设置主机名
  updateHostname: (hostname: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/hostname', { hostname }),
  // 设置 Hosts
  updateHosts: (hosts: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/hosts', { hosts }),
  // 设置 Root 密码
  updateRootPassword: (password: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/rootPassword', { password })
}
