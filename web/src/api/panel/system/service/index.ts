import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 服务状态
  status: (service: string): Promise<AxiosResponse<any>> =>
    request.get('/panel/system/service/status', { params: { service } }),
  // 是否启用服务
  isEnabled: (service: string): Promise<AxiosResponse<any>> =>
    request.get('/panel/system/service/isEnabled', { params: { service } }),
  // 启用服务
  enable: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/system/service/enable', { service }),
  // 禁用服务
  disable: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/system/service/disable', { service }),
  // 重启服务
  restart: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/system/service/restart', { service }),
  // 重载服务
  reload: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/system/service/reload', { service }),
  // 启动服务
  start: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/system/service/start', { service }),
  // 停止服务
  stop: (service: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/system/service/stop', { service })
}
