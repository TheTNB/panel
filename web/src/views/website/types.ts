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
  name: string
  ports: string[]
  ssl_ports: string[]
  quic_ports: string[]
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
  waf: boolean
  waf_mode: string
  waf_cc_deny: string
  waf_cache: string
  rewrite: string
  raw: string
  log: string
}

export interface Backup {
  name: string
  size: string
}
