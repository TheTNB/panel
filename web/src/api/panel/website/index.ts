import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/website', { params: { page, limit } }),
  // 创建
  create: (data: any): Promise<AxiosResponse<any>> => request.post('/website', data),
  // 删除
  delete: (id: number, path: boolean, db: boolean): Promise<AxiosResponse<any>> =>
    request.delete(`/website/${id}`, { data: { path, db } }),
  // 获取默认配置
  defaultConfig: (): Promise<AxiosResponse<any>> => request.get('/website/defaultConfig'),
  // 保存默认配置
  saveDefaultConfig: (index: string, stop: string): Promise<AxiosResponse<any>> =>
    request.post('/website/defaultConfig', { index, stop }),
  // 网站配置
  config: (id: number): Promise<AxiosResponse<any>> => request.get('/website/' + id),
  // 保存网站配置
  saveConfig: (id: number, data: any): Promise<AxiosResponse<any>> =>
    request.put('/website/' + id, data),
  // 清空日志
  clearLog: (id: number): Promise<AxiosResponse<any>> => request.delete('/website/' + id + '/log'),
  // 更新备注
  updateRemark: (id: number, remark: string): Promise<AxiosResponse<any>> =>
    request.post('/website/' + id + '/updateRemark', { remark }),
  // 重置配置
  resetConfig: (id: number): Promise<AxiosResponse<any>> =>
    request.post('/website/' + id + '/resetConfig'),
  // 修改状态
  status: (id: number, status: boolean): Promise<AxiosResponse<any>> =>
    request.post('/website/' + id + '/status', { status })
}
