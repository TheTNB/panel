export interface Module {
  name: string
  path: string
  comment: string
  read_only: boolean
  auth_user: string
  secret: string
  hosts_allow: string
}
