import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 面板信息
  panel: (): Promise<Response> => fetch('/api/info/panel'),
  // 面板菜单
  menu: (): Promise<AxiosResponse<any>> => request.get('/info/menu'),
  // 首页应用
  homeApps: (): Promise<AxiosResponse<any>> => request.get('/info/homeApps'),
  // 实时监控
  realtime: (): Promise<AxiosResponse<any>> => request.get('/info/realtime'),
  // 系统信息
  systemInfo: (): Promise<AxiosResponse<any>> => request.get('/info/systemInfo'),
  // 统计信息
  countInfo: (): Promise<AxiosResponse<any>> => request.get('/info/countInfo'),
  // 已安装的数据库和PHP
  installedDbAndPhp: (): Promise<AxiosResponse<any>> => request.get('/info/installedDbAndPhp'),
  // 检查更新
  checkUpdate: (): Promise<AxiosResponse<any>> => request.get('/info/checkUpdate'),
  // 更新日志
  updateInfo: (): Promise<AxiosResponse<any>> => request.get('/info/updateInfo'),
  // 更新面板
  update: (): Promise<AxiosResponse<any>> => request.post('/info/update', null),
  // 重启面板
  restart: (): Promise<AxiosResponse<any>> => request.post('/info/restart')
}
