import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 负载状态
  load: (version: number): Promise<AxiosResponse<any>> => request.get(`/apps/php${version}/load`),
  // 获取配置
  config: (version: number): Promise<AxiosResponse<any>> =>
    request.get(`/apps/php${version}/config`),
  // 保存配置
  saveConfig: (version: number, config: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/php${version}/config`, { config }),
  // 获取FPM配置
  fpmConfig: (version: number): Promise<AxiosResponse<any>> =>
    request.get(`/apps/php${version}/fpmConfig`),
  // 保存FPM配置
  saveFPMConfig: (version: number, config: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/php${version}/fpmConfig`, { config }),
  // 获取错误日志
  errorLog: (version: number): Promise<AxiosResponse<any>> =>
    request.get(`/apps/php${version}/errorLog`),
  // 清空错误日志
  clearErrorLog: (version: number): Promise<AxiosResponse<any>> =>
    request.post(`/apps/php${version}/clearErrorLog`),
  // 获取慢日志
  slowLog: (version: number): Promise<AxiosResponse<any>> =>
    request.get(`/apps/php${version}/slowLog`),
  // 清空慢日志
  clearSlowLog: (version: number): Promise<AxiosResponse<any>> =>
    request.post(`/apps/php${version}/clearSlowLog`),
  // 拓展列表
  extensions: (version: number): Promise<AxiosResponse<any>> =>
    request.get(`/apps/php${version}/extensions`),
  // 安装拓展
  installExtension: (version: number, slug: string): Promise<AxiosResponse<any>> =>
    request.post(`/apps/php${version}/extensions`, { slug }),
  // 卸载拓展
  uninstallExtension: (version: number, slug: string): Promise<AxiosResponse<any>> =>
    request.delete(`/apps/php${version}/extensions`, { params: { slug } })
}
