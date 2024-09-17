import type { AxiosError, AxiosResponse } from 'axios'
import { AxiosRejectError, resolveResError } from './helpers'
import type { RequestConfig } from '~/types/axios'

/** 请求拦截 */
export function reqResolve(config: RequestConfig) {
  return config
}

/** 请求错误拦截 */
export function reqReject(error: AxiosError) {
  return Promise.reject(error)
}

/** 响应拦截 */
export function resResolve(response: AxiosResponse) {
  const { data, status, config, statusText } = response
  if (status !== 200) {
    const code = data?.code ?? status
    const message = resolveResError(code, data?.message ?? statusText)
    const { noNeedTip } = config as RequestConfig

    if (!noNeedTip) {
      if (code == 422) {
        window.$message.error(message)
      } else {
        if (code != 401) {
          window.$dialog.error({
            title: '请求返回异常',
            content: message,
            maskClosable: false
          })
        }
      }
    }

    return Promise.reject(new AxiosRejectError({ code, message, data: data || response }))
  }

  return Promise.resolve(data)
}

/** 响应错误拦截 */
export function resReject(error: AxiosError) {
  if (!error || !error.response) {
    const code = error?.code
    /** 根据code处理对应的操作，并返回处理后的message */
    const message = resolveResError(code, error.message)
    window.$dialog.error({
      title: '请求出现异常',
      content: message,
      maskClosable: false
    })
    return Promise.reject(new AxiosRejectError({ code, message, data: error }))
  }
  const { data, status, config } = error.response
  let { code, message } = data as AxiosRejectError
  code = code ?? status
  message = message ?? error.message
  message = resolveResError(code, message)
  /** 需要错误提醒 */
  const { noNeedTip } = config as RequestConfig

  if (!noNeedTip) {
    if (code == 422) {
      window.$message.error(message)
    } else {
      if (code != 401) {
        window.$dialog.error({
          title: '请求返回异常',
          content: message,
          maskClosable: false
        })
      }
    }
  }

  return Promise.reject(
    new AxiosRejectError({
      code,
      message,
      data: error.response?.data || error.response
    })
  )
}
