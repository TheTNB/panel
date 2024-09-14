import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 面板信息
  panel: (): Promise<Response> => fetch('/api/panel/info/panel'),
  // 面板菜单
  menu: (): Promise<AxiosResponse<any>> => request.get('/panel/info/menu'),
  // 首页插件
  homePlugins: (): Promise<AxiosResponse<any>> => request.get('/panel/info/homePlugins'),
  // 实时监控
  nowMonitor: (): Promise<AxiosResponse<any>> => request.get('/panel/info/nowMonitor'),
  // 系统信息
  systemInfo: (): Promise<AxiosResponse<any>> => request.get('/panel/info/systemInfo'),
  // 统计信息
  countInfo: (): Promise<AxiosResponse<any>> => request.get('/panel/info/countInfo'),
  // 已安装的数据库和PHP
  installedDbAndPhp: (): Promise<AxiosResponse<any>> =>
    request.get('/panel/info/installedDbAndPhp'),
  // 检查更新
  checkUpdate: (): Promise<AxiosResponse<any>> => request.get('/panel/info/checkUpdate'),
  // 更新日志
  updateInfo: (): Promise<AxiosResponse<any>> => request.get('/panel/info/updateInfo'),
  // 更新面板
  update: (): Promise<AxiosResponse<any>> => request.post('/panel/info/update', null),
  // 重启面板
  restart: (): Promise<AxiosResponse<any>> => request.post('/panel/info/restart')
}
