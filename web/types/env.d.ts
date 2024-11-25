type ProxyType = 'dev' | 'test' | 'prod'

interface ViteEnv {
  VITE_PORT: number
  VITE_USE_PROXY?: boolean
  VITE_USE_HASH?: boolean
  VITE_APP_TITLE: string
  VITE_PUBLIC_PATH: string
  VITE_BASE_API: string
  VITE_PROXY_TYPE?: ProxyType
}

interface ProxyConfig {
  /** 匹配代理的前缀，接口地址匹配到此前缀将代理的target地址 */
  prefix: string
  /** 代理目标地址，后端真实接口地址 */
  target: string
  /** 是否校验https证书 */
  secure?: boolean
  /** 是否修改请求头中的host */
  changeOrigin?: boolean
  /** 是否代理websocket */
  ws?: boolean
}
