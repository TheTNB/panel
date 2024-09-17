import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取信息
  info: (): Promise<AxiosResponse<any>> => request.get('/plugins/phpmyadmin/info'),
  // 设置端口
  port: (port: number): Promise<AxiosResponse<any>> =>
    request.post('/plugins/phpmyadmin/port', { port }),
  // 获取配置
  getConfig: (): Promise<AxiosResponse<any>> => request.get('/plugins/phpmyadmin/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/phpmyadmin/config', { config })
}
