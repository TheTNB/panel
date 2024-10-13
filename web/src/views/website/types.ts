export interface Website {
  id: number
  name: string
  status: boolean
  path: string
  php: number
  ssl: boolean
  remark: string
  created_at: string
  updated_at: string
}

export interface WebsiteSetting {
  id: number
  name: string
  ports: number[]
  ssl_ports: number[]
  quic_ports: number[]
  domains: string[]
  root: string
  path: string
  index: string
  php: number
  open_basedir: boolean
  ssl: boolean
  ssl_certificate: string
  ssl_certificate_key: string
  ssl_not_before: string
  ssl_not_after: string
  ssl_dns_names: string[]
  ssl_issuer: string
  ssl_ocsp_server: string[]
  http_redirect: boolean
  hsts: boolean
  ocsp: boolean
  rewrite: string
  raw: string
  log: string
}
