import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 获取配置
  config: (name: string): Promise<AxiosResponse<any>> =>
    request.get('/apps/frp/config', { params: { name } }),
  // 保存配置
  saveConfig: (name: string, config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/frp/config', { name, config })
}
