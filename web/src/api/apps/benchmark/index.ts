import type { AxiosResponse } from 'axios'

import { request } from '@/utils'

export default {
  // 运行评分
  test: (name: string, multi: boolean): Promise<AxiosResponse<any>> =>
    request.post('/apps/benchmark/test', { name, multi })
}
