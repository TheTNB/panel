import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取配置
  config: (service: string): Promise<AxiosResponse<any>> =>
    request.get('/plugins/frp/config', { params: { service } }),
  // 保存配置
  saveConfig: (service: string, config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/frp/config', { service, config })
}
