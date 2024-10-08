import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取配置
  config: (service: string): Promise<AxiosResponse<any>> =>
    request.get('/apps/frp/config', { params: { service } }),
  // 保存配置
  saveConfig: (service: string, config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/frp/config', { service, config })
}
