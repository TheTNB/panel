import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 负载状态
  load: (): Promise<AxiosResponse<any>> => request.get('/apps/mysql/load'),
  // 获取配置
  config: (): Promise<AxiosResponse<any>> => request.get('/apps/mysql/config'),
  // 保存配置
  saveConfig: (config: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/mysql/config', { config }),
  // 获取错误日志
  errorLog: (): Promise<AxiosResponse<any>> => request.get('/apps/mysql/errorLog'),
  // 清空错误日志
  clearErrorLog: (): Promise<AxiosResponse<any>> => request.post('/apps/mysql/clearErrorLog'),
  // 获取慢查询日志
  slowLog: (): Promise<AxiosResponse<any>> => request.get('/apps/mysql/slowLog'),
  // 清空慢查询日志
  clearSlowLog: (): Promise<AxiosResponse<any>> => request.post('/apps/mysql/clearSlowLog'),
  // 获取 root 密码
  rootPassword: (): Promise<AxiosResponse<any>> => request.get('/apps/mysql/rootPassword'),
  // 修改 root 密码
  setRootPassword: (password: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/mysql/rootPassword', { password }),
  // 数据库列表
  databases: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/mysql/databases', { params: { page, limit } }),
  // 创建数据库
  addDatabase: (database: any): Promise<AxiosResponse<any>> =>
    request.post('/apps/mysql/databases', database),
  // 删除数据库
  deleteDatabase: (database: string): Promise<AxiosResponse<any>> =>
    request.delete('/apps/mysql/databases', { params: { database } }),
  // 备份列表
  backups: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/mysql/backups', { params: { page, limit } }),
  // 创建备份
  createBackup: (database: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/mysql/backups', { database }),
  // 上传备份
  uploadBackup: (backup: any): Promise<AxiosResponse<any>> =>
    request.put('/apps/mysql/backups', backup),
  // 删除备份
  deleteBackup: (name: string): Promise<AxiosResponse<any>> =>
    request.delete('/apps/mysql/backups', { params: { name } }),
  // 还原备份
  restoreBackup: (backup: string, database: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/mysql/backups/restore', { backup, database }),
  // 用户列表
  users: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/apps/mysql/users', { params: { page, limit } }),
  // 创建用户
  addUser: (user: any): Promise<AxiosResponse<any>> => request.post('/apps/mysql/users', user),
  // 删除用户
  deleteUser: (user: string): Promise<AxiosResponse<any>> =>
    request.delete('/apps/mysql/users', { params: { user } }),
  // 设置用户密码
  setUserPassword: (user: string, password: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/mysql/users/password', { user, password }),
  // 设置用户权限
  setUserPrivileges: (user: string, database: string): Promise<AxiosResponse<any>> =>
    request.post('/apps/mysql/users/privileges', { user, database })
}
