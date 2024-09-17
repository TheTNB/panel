import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // DNS
  dns: (): Promise<AxiosResponse<any>> => request.get('/plugins/toolbox/dns'),
  // 设置 DNS
  setDns: (dns1: string, dns2: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/toolbox/dns', { dns1, dns2 }),
  // SWAP
  swap: (): Promise<AxiosResponse<any>> => request.get('/plugins/toolbox/swap'),
  // 设置 SWAP
  setSwap: (size: number): Promise<AxiosResponse<any>> =>
    request.post('/plugins/toolbox/swap', { size }),
  // 时区
  timezone: (): Promise<AxiosResponse<any>> => request.get('/plugins/toolbox/timezone'),
  // 设置时区
  setTimezone: (timezone: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/toolbox/timezone', { timezone }),
  // Hosts
  hosts: (): Promise<AxiosResponse<any>> => request.get('/plugins/toolbox/hosts'),
  // 设置 Hosts
  setHosts: (hosts: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/toolbox/hosts', { hosts }),
  // 设置 Root 密码
  setRootPassword: (password: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/toolbox/rootPassword', { password })
}
