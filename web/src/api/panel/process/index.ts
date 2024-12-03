import { http } from '@/utils'

export default {
  // 获取进程列表
  list: (page: number, limit: number) => http.Get(`/process`, { params: { page, limit } }),
  // 杀死进程
  kill: (pid: number) => http.Post(`/process/kill`, { pid })
}
