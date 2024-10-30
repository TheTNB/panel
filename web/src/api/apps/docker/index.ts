import { http } from '@/utils'

export const getConfig = () => http.Get('/apps/docker/config')
export const updateConfig = (config: string) => http.Post('/apps/docker/config', { config })
