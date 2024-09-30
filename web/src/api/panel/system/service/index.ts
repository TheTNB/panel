import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 服务状态
  status: (service: string): Promise<AxiosResponse<any>> =>
    request.get('/system/service/status', { params: { service } }),
  // 是否启用服务
  isEnabled: (service: string): Promise<AxiosResponse<any>> =>
    request.get('/system/service/isEnabled', { params: { service } }),
  // 启用服务
  enable: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/system/service/enable', { service }),
  // 禁用服务
  disable: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/system/service/disable', { service }),
  // 重启服务
  restart: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/system/service/restart', { service }),
  // 重载服务
  reload: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/system/service/reload', { service }),
  // 启动服务
  start: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/system/service/start', { service }),
  // 停止服务
  stop: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/system/service/stop', { service })
}
