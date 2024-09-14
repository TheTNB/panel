import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 服务名称
  service: (): Promise<AxiosResponse<any>> => request.get('/plugins/supervisor/service'),
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/plugins/supervisor/load'),
  // 获取错误日志
  log: (): Promise<AxiosResponse<any>> => request.get('/plugins/supervisor/log'),
  // 清空错误日志
  clearLog: (): Promise<AxiosResponse<any>> => request.post('/plugins/supervisor/clearLog'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/plugins/supervisor/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/config', { config }),
  // 进程列表
  processes: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/plugins/supervisor/processes', { params: { page, limit } }),
  // 进程启动
  startProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/startProcess', { process }),
  // 进程停止
  stopProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/stopProcess', { process }),
  // 进程重启
  restartProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/restartProcess', { process }),
  // 进程日志
  processLog: (process: string): Promise<AxiosResponse<any>> =>
    request.get('/plugins/supervisor/processLog', { params: { process } }),
  // 清空进程日志
  clearProcessLog: (process: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/clearProcessLog', { process }),
  // 进程配置
  processConfig: (process: string): Promise<AxiosResponse<any>> =>
    request.get('/plugins/supervisor/processConfig', { params: { process } }),
  // 保存进程配置
  saveProcessConfig: (process: string, config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/processConfig', { process, config }),
  // 添加进程
  addProcess: (process: any): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/addProcess', process),
  // 删除进程
  deleteProcess: (process: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/supervisor/deleteProcess', { process })
}
