import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 获取注册表配置
  registryConfig: (): Promise<AxiosResponse<any>> => request.get('/apps/podman/registryConfig'),
  // 保存注册表配置
  saveRegistryConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/podman/registryConfig', { config }),
  // 获取存储配置
  storageConfig: (): Promise<AxiosResponse<any>> => request.get('/apps/podman/storageConfig'),
  // 保存存储配置
  saveStorageConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/podman/storageConfig', { config })
}
