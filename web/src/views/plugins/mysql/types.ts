export interface Database {
  name: string
}

export interface User {
  user: string
  host: string
  grants: string
}

export interface Backup {
  name: string
  size: string
}
