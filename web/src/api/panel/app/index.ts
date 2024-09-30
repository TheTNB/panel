import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取应用列表
  list: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/app/list', { params: { page, limit } }),
  // 安装应用
  install: (slug: string): Promise<AxiosResponse<any>> => request.post('/app/install', { slug }),
  // 卸载应用
  uninstall: (slug: string): Promise<AxiosResponse<any>> =>
    request.post('/app/uninstall', { slug }),
  // 更新应用
  update: (slug: string): Promise<AxiosResponse<any>> => request.post('/app/update', { slug }),
  // 设置首页显示
  updateShow: (slug: string, show: boolean): Promise<AxiosResponse<any>> =>
    request.post('/app/updateShow', { slug, show }),
  // 应用是否已安装
  isInstalled: (slug: string): Promise<AxiosResponse<any>> =>
    request.get('/app/isInstalled', { params: { slug } })
}
