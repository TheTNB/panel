const proxyConfigMappings: Record<ProxyType, ProxyConfig[]> = {
  dev: [
    {
      prefix: '/api/ws',
      target: 'ws://localhost:8888/api/ws',
      secure: false
    },
    {
      prefix: '/api',
      target: 'http://localhost:8080/api'
    }
  ],
  test: [
    {
      prefix: '/api',
      target: 'http://localhost:8080/api'
    }
  ],
  prod: [
    {
      prefix: '/api',
      target: 'http://localhost:8080/api'
    }
  ]
}

export function getProxyConfigs(envType: ProxyType = 'dev'): ProxyConfig[] {
  return proxyConfigMappings[envType]
}
