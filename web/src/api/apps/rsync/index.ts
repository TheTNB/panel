import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/apps/rsync/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/rsync/config', { config }),
  // 模块列表
  modules: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/rsync/modules', { params: { page, limit } }),
  // 添加模块
  addModule: (module: any): Promise<AxiosResponse<any>> =>
    request.post('/apps/rsync/modules', module),
  // 删除模块
  deleteModule: (name: string): Promise<AxiosResponse<any>> =>
    request.delete('/apps/rsync/modules/' + name),
  // 更新模块
  updateModule: (name: string, module: any): Promise<AxiosResponse<any>> =>
    request.post('/apps/rsync/modules/' + name, module)
}
