interface Load {
  load1: number[]
  load5: number[]
  load15: number[]
}

interface Cpu {
  percent: string[]
}

interface Mem {
  total: string
  available: string[]
  used: string[]
}

interface Swap {
  total: string
  used: string[]
  free: string[]
}

interface Network {
  sent: string[]
  recv: string[]
  tx: string[]
  rx: string[]
}

export interface MonitorData {
  times: string[]
  load: Load
  cpu: Cpu
  mem: Mem
  swap: Swap
  net: Network
}
