export interface Plugin {
  name: string
  description: string
  slug: string
  version: string
  requires: string
  excludes: string
  installed: boolean
  installed_version: string
  show: boolean
}
