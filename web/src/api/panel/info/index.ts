import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 面板信息
  panel: (): Promise<Response> => fetch('/api/info/panel'),
  // 面板菜单
  menu: (): Promise<AxiosResponse<any>> => request.get('/info/menu'),
  // 首页插件
  homePlugins: (): Promise<AxiosResponse<any>> => request.get('/info/homePlugins'),
  // 实时监控
  nowMonitor: (): Promise<AxiosResponse<any>> => request.get('/info/nowMonitor'),
  // 系统信息
  systemInfo: (): Promise<AxiosResponse<any>> => request.get('/info/systemInfo'),
  // 统计信息
  countInfo: (): Promise<AxiosResponse<any>> => request.get('/info/countInfo'),
  // 已安装的数据库和PHP
  installedDbAndPhp: (): Promise<AxiosResponse<any>> =>
    request.get('/info/installedDbAndPhp'),
  // 检查更新
  checkUpdate: (): Promise<AxiosResponse<any>> => request.get('/info/checkUpdate'),
  // 更新日志
  updateInfo: (): Promise<AxiosResponse<any>> => request.get('/info/updateInfo'),
  // 更新面板
  update: (): Promise<AxiosResponse<any>> => request.post('/info/update', null),
  // 重启面板
  restart: (): Promise<AxiosResponse<any>> => request.post('/info/restart')
}
