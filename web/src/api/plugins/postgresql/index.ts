import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/plugins/postgresql/load'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/plugins/postgresql/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/postgresql/config', { config }),
  // 获取用户配置
  userConfig: (): Promise<AxiosResponse<any>> => request.get('/plugins/postgresql/userConfig'),
  // 保存配置
  saveUserConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/postgresql/userConfig', { config }),
  // 获取日志
  log: (): Promise<AxiosResponse<any>> => request.get('/plugins/postgresql/log'),
  // 清空错误日志
  clearLog: (): Promise<AxiosResponse<any>> => request.post('/plugins/postgresql/clearLog'),
  // 数据库列表
  databases: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/plugins/postgresql/databases', { params: { page, limit } }),
  // 创建数据库
  addDatabase: (database: any): Promise<AxiosResponse<any>> =>
    request.post('/plugins/postgresql/databases', database),
  // 删除数据库
  deleteDatabase: (database: string): Promise<AxiosResponse<any>> =>
    request.delete('/plugins/postgresql/databases', { params: { database } }),
  // 备份列表
  backups: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/plugins/postgresql/backups', { params: { page, limit } }),
  // 创建备份
  createBackup: (database: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/postgresql/backups', { database }),
  // 上传备份
  uploadBackup: (backup: any): Promise<AxiosResponse<any>> =>
    request.put('/plugins/postgresql/backups', backup),
  // 删除备份
  deleteBackup: (name: string): Promise<AxiosResponse<any>> =>
    request.delete('/plugins/postgresql/backups', { params: { name } }),
  // 还原备份
  restoreBackup: (backup: string, database: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/postgresql/backups/restore', { backup, database }),
  // 角色列表
  roles: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/plugins/postgresql/roles', { params: { page, limit } }),
  // 创建角色
  addRole: (user: any): Promise<AxiosResponse<any>> =>
    request.post('/plugins/postgresql/roles', user),
  // 删除角色
  deleteRole: (user: string): Promise<AxiosResponse<any>> =>
    request.delete('/plugins/postgresql/roles', { params: { user } }),
  // 设置角色密码
  setRolePassword: (user: string, password: string): Promise<AxiosResponse<any>> =>
    request.post('/plugins/postgresql/roles/password', { user, password })
}
