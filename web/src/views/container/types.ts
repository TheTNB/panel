export interface ContainerList {
  command: string
  created: string
  id: string
  image: string
  image_id: string
  labels: {
    [key: string]: string
  }
  name: string
  ports: [
    {
      IP: string
      PrivatePort: number
      PublicPort: number
      Type: string
    }
  ]
  state: string
  status: string
}

export interface ImageList {
  id: string
  created: string
  containers: number
  size: string
  labels: {
    [key: string]: string
  }
  repo_tags: string[]
  repo_digests: string[]
}

export interface NetworkList {
  id: string
  name: string
  driver: string
  ipv6: boolean
  scope: string
  internal: boolean
  attachable: boolean
  ingress: boolean
  labels: {
    [key: string]: string
  }
  options: {
    [key: string]: string
  }
  ipam: {
    driver: string
    config: {
      subnet: string
      gateway: string
      ip_range: string
      aux_address: {
        [key: string]: string
      }
    }[]
    options: {
      [key: string]: string
    }
  }
}

export interface VolumeList {
  id: string
  created: string
  driver: string
  mount: string
  labels: {
    [key: string]: string
  }
  options: {
    [key: string]: string
  }
  scope: string
  status: {
    [key: string]: string
  }
  usage: {
    ref_count: number
    size: string
  }
}
