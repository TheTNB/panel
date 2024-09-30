import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // DNS
  dns: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/dns'),
  // 设置 DNS
  setDns: (dns1: string, dns2: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/dns', { dns1, dns2 }),
  // SWAP
  swap: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/swap'),
  // 设置 SWAP
  setSwap: (size: number): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/swap', { size }),
  // 时区
  timezone: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/timezone'),
  // 设置时区
  setTimezone: (timezone: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/timezone', { timezone }),
  // Hosts
  hosts: (): Promise<AxiosResponse<any>> => request.get('/apps/toolbox/hosts'),
  // 设置 Hosts
  setHosts: (hosts: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/hosts', { hosts }),
  // 设置 Root 密码
  setRootPassword: (password: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/toolbox/rootPassword', { password })
}
