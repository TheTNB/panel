export interface Cert {
  id: number
  user_id: number
  website_id: number
  dns_id: number
  type: string
  domains: string[]
  auto_renew: boolean
  cert_url: string
  cert: string
  key: string
  created_at: string
  updated_at: string
  website: Website
  dns: DNS
  user: User
}

export interface Website {
  id: number
  name: string
  status: boolean
  path: string
  php: string
  ssl: boolean
  remark: string
  created_at: string
  updated_at: string
}

export interface DNS {
  id: number
  type: string
  name: string
  data: {
    id: string
    token: string
    access_key: string
    secret_key: string
    api_key: string
  }
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  email: string
  ca: string
  kid: string
  hmac_encoded: string
  private_key: string
  key_type: string
  created_at: string
  updated_at: string
}
