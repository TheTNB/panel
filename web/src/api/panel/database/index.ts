import { http } from '@/utils'

export default {
  // 获取数据库列表
  list: (page: number, limit: number) => http.Get(`/database`, { params: { page, limit } }),
  // 创建数据库
  create: (data: any) => http.Post(`/database`, data),
  // 删除数据库
  delete: (server_id: number, name: string) => http.Delete(`/database`, { server_id, name }),
  // 获取数据库服务器列表
  serverList: (page: number, limit: number) =>
    http.Get('/databaseServer', { params: { page, limit } }),
  // 创建数据库服务器
  serverCreate: (data: any) => http.Post('/databaseServer', data),
  // 更新数据库服务器
  serverUpdate: (id: number, data: any) => http.Put(`/databaseServer/${id}`, data),
  // 删除数据库服务器
  serverDelete: (id: number) => http.Delete(`/databaseServer/${id}`),
  // 同步数据库
  serverSync: (id: number) => http.Post(`/databaseServer/${id}/sync`),
  // 获取数据库用户列表
  userList: (page: number, limit: number) => http.Get('/databaseUser', { params: { page, limit } }),
  // 创建数据库用户
  userCreate: (data: any) => http.Post('/databaseUser', data),
  // 更新数据库用户
  userUpdate: (id: number, data: any) => http.Put(`/databaseUser/${id}`, data),
  // 删除数据库用户
  userDelete: (id: number) => http.Delete(`/databaseUser/${id}`)
}
