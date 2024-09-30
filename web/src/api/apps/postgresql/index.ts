import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/apps/postgresql/load'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/apps/postgresql/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/postgresql/config', { config }),
  // 获取用户配置
  userConfig: (): Promise<AxiosResponse<any>> => request.get('/apps/postgresql/userConfig'),
  // 保存配置
  saveUserConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/postgresql/userConfig', { config }),
  // 获取日志
  log: (): Promise<AxiosResponse<any>> => request.get('/apps/postgresql/log'),
  // 清空错误日志
  clearLog: (): Promise<AxiosResponse<any>> => request.post('/apps/postgresql/clearLog'),
  // 数据库列表
  databases: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/postgresql/databases', { params: { page, limit } }),
  // 创建数据库
  addDatabase: (database: any): Promise<AxiosResponse<any>> =>
    request.post('/apps/postgresql/databases', database),
  // 删除数据库
  deleteDatabase: (database: string): Promise<AxiosResponse<any>> =>
    request.delete('/apps/postgresql/databases', { params: { database } }),
  // 备份列表
  backups: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/postgresql/backups', { params: { page, limit } }),
  // 创建备份
  createBackup: (database: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/postgresql/backups', { database }),
  // 上传备份
  uploadBackup: (backup: any): Promise<AxiosResponse<any>> =>
    request.put('/apps/postgresql/backups', backup),
  // 删除备份
  deleteBackup: (name: string): Promise<AxiosResponse<any>> =>
    request.delete('/apps/postgresql/backups', { params: { name } }),
  // 还原备份
  restoreBackup: (backup: string, database: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/postgresql/backups/restore', { backup, database }),
  // 角色列表
  roles: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/postgresql/roles', { params: { page, limit } }),
  // 创建角色
  addRole: (user: any): Promise<AxiosResponse<any>> => request.post('/apps/postgresql/roles', user),
  // 删除角色
  deleteRole: (user: string): Promise<AxiosResponse<any>> =>
    request.delete('/apps/postgresql/roles', { params: { user } }),
  // 设置角色密码
  setRolePassword: (user: string, password: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/postgresql/roles/password', { user, password })
}
