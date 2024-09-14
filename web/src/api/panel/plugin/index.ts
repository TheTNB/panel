import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取插件列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/plugin/list', { params: { page, limit } }),
  // 安装插件
  install: (slug: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/plugin/install', { slug }),
  // 卸载插件
  uninstall: (slug: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/plugin/uninstall', { slug }),
  // 更新插件
  update: (slug: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/plugin/update', { slug }),
  // 设置首页显示
  updateShow: (slug: string, show: boolean): Promise<AxiosResponse<any>> =>
    request.post('/panel/plugin/updateShow', { slug, show }),
  // 插件是否已安装
  isInstalled: (slug: string): Promise<AxiosResponse<any>> =>
    request.get('/panel/plugin/isInstalled', { params: { slug } })
}
