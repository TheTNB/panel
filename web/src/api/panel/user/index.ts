import { http } from '@/utils'

export default {
  // 公钥
  key: () => http.Get('/user/key'),
  // 登录
  login: (username: string, password: string) =>
    http.Post('/user/login', {
      username,
      password
    }),
  // 登出
  logout: () => http.Post('/user/logout'),
  // 是否登录
  isLogin: () => http.Get('/user/isLogin'),
  // 获取用户信息
  info: () => http.Get('/user/info')
}
