import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 登录
  login: (username: string, password: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/user/login', {
      username,
      password
    }),
  // 登出
  logout: (): Promise<AxiosResponse<any>> => request.post('/panel/user/logout'),
  // 是否登录
  isLogin: (): Promise<AxiosResponse<any>> => request.get('/panel/user/isLogin'),
  // 获取用户信息
  info: (): Promise<AxiosResponse<any>> => request.get('/panel/user/info')
}
