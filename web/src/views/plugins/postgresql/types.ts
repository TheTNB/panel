export interface Database {
  name: string
}

export interface Role {
  role: string
  attributes: string[]
}

export interface Backup {
  name: string
  size: string
}
