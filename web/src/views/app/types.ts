export interface App {
  name: string
  description: string
  slug: string
  channels: Channel[]
  installed: boolean
  installed_channel: string
  installed_version: string
  update_exist: boolean
  show: boolean
}

export interface Channel {
  slug: string
  name: string
  panel: string
  subs: Sub[]
}

export interface Sub {
  log: string
  version: string
}
