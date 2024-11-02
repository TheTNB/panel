import { http } from '@/utils'

// 负载状态
export const getLoad = () => http.Get('/apps/memcached/load')
// 获取配置
export const getConfig = () => http.Get('/apps/memcached/config')
// 保存配置
export const updateConfig = (config: string) => http.Post('/apps/memcached/config', { config })
