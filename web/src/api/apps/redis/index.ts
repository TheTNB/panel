import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/apps/redis/load'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/apps/redis/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/redis/config', { config })
}
