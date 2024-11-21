import { resolveResError } from '@/utils/http/helpers'
import { createAlova, Method } from 'alova'
import adapterFetch from 'alova/fetch'
import VueHook from 'alova/vue'
import axios from 'axios'
import { reqReject, reqResolve, resReject, resResolve } from './interceptors'

export function createAxios(options = {}) {
  const defaultOptions = {
    adapter: 'fetch',
    timeout: 0
  }
  const service = axios.create({
    ...defaultOptions,
    ...options
  })
  service.interceptors.request.use(reqResolve, reqReject)
  service.interceptors.response.use(resResolve, resReject)
  return service
}

export const request = createAxios({
  baseURL: import.meta.env.VITE_BASE_API
})

export const http = createAlova({
  id: 'panel',
  cacheFor: null,
  statesHook: VueHook,
  requestAdapter: adapterFetch(),
  baseURL: import.meta.env.VITE_BASE_API,
  responded: {
    onSuccess: async (response: Response, method: Method) => {
      const ct = response.headers.get('Content-Type')
      const json =
        ct && ct.includes('application/json')
          ? await response.json()
          : { code: response.status, message: await response.text() }
      const { status, statusText } = response
      const { meta } = method
      if (status !== 200) {
        const code = json?.code ?? status
        const message = resolveResError(code, json?.message ?? statusText)
        const noAlert = meta?.noAlert
        if (!noAlert) {
          if (code === 422) {
            window.$message.error(message)
          } else if (code !== 401) {
            window.$dialog.error({
              title: '接口响应异常',
              content: message,
              maskClosable: false
            })
          }
        }
        throw new Error(message)
      }
      return json.data
    },
    onError: (error: any, method: Method) => {
      const { code, message } = error
      const { meta } = method
      const errorMessage = resolveResError(code, message)
      const noAlert = meta?.noAlert

      if (!noAlert) {
        window.$dialog.error({
          title: '接口请求失败',
          content: errorMessage,
          maskClosable: false
        })
      }

      throw error
    }
  }
})
