import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/plugins/gitea/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/gitea/config', { config })
}
