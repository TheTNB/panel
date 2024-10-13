import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 服务名称
  service: (): Promise<AxiosResponse<any>> => request.get('/apps/supervisor/service'),
  // 获取错误日志
  log: (): Promise<AxiosResponse<any>> => request.get('/apps/supervisor/log'),
  // 清空错误日志
  clearLog: (): Promise<AxiosResponse<any>> => request.post('/apps/supervisor/clearLog'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/apps/supervisor/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/supervisor/config', { config }),
  // 进程列表
  processes: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/supervisor/processes', { params: { page, limit } }),
  // 进程启动
  startProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/supervisor/processes/${process}/start`, {}),
  // 进程停止
  stopProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/supervisor/processes/${process}/stop`, {}),
  // 进程重启
  restartProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/supervisor/processes/${process}/restart`, {}),
  // 进程日志
  processLog: (process: string): Promise<AxiosResponse<any>> =>
    request.get(`/apps/supervisor/processes/${process}/log`),
  // 清空进程日志
  clearProcessLog: (process: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/supervisor/processes/${process}/clearLog`, {}),
  // 进程配置
  processConfig: (process: string): Promise<AxiosResponse<any>> =>
    request.get(`/apps/supervisor/processes/${process}`),
  // 保存进程配置
  saveProcessConfig: (process: string, config: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/supervisor/processes/${process}`, { config }),
  // 创建进程
  createProcess: (process: any): Promise<AxiosResponse<any>> =>
    request.post('/apps/supervisor/processes', process),
  // 删除进程
  deleteProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.delete(`/apps/supervisor/processes/${process}`)
}
