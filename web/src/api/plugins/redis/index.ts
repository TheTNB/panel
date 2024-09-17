import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/plugins/redis/load'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/plugins/redis/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/redis/config', { config })
}
