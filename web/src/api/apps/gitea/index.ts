import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/apps/gitea/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/gitea/config', { config })
}
