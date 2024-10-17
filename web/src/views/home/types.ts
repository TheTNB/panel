interface CpuInfoStat {
  cpu: number
  vendorId: string
  family: string
  model: string
  stepping: number
  physicalId: string
  coreId: string
  cores: number
  modelName: string
  mhz: number
  cacheSize: number
  flags: string[]
  microcode: string
}

interface LoadAvgStat {
  load1: number
  load5: number
  load15: number
}

interface HostInfoStat {
  hostname: string
  uptime: number
  bootTime: number
  procs: number
  os: string
  platform: string
  platformFamily: string
  platformVersion: string
  kernelVersion: string
  kernelArch: string
  virtualizationSystem: string
  virtualizationRole: string
  hostid: string
}

interface VirtualMemoryStat {
  total: number
  available: number
  used: number
  usedPercent: number
  free: number
  active: number
  inactive: number
  wired: number
  laundry: number
  buffers: number
  cached: number
  writeback: number
  dirty: number
  writebacktmp: number
  shared: number
  slab: number
  sreclaimable: number
  sunreclaim: number
  pagetables: number
  swapcached: number
  commitlimit: number
  committedas: number
  hightotal: number
  highfree: number
  lowtotal: number
  lowfree: number
  swaptotal: number
  swapfree: number
  mapped: number
  vmalloctotal: number
  vmallocused: number
  vmallocchunk: number
  hugepagestotal: number
  hugepagesfree: number
  hugepagesize: number
}

interface SwapMemoryStat {
  total: number
  used: number
  free: number
  usedPercent: number
  sin: number
  sout: number
  pgin: number
  pgout: number
  pgfault: number
  pgmajfault: number
}

interface IOCountersStat {
  name: string
  bytesSent: number
  bytesRecv: number
  packetsSent: number
  packetsRecv: number
  errin: number
  errout: number
  dropin: number
  dropout: number
  fifoin: number
  fifoout: number
}

interface DiskIOCountersStat {
  readCount: number
  mergedReadCount: number
  writeCount: number
  mergedWriteCount: number
  readBytes: number
  writeBytes: number
  readTime: number
  writeTime: number
  iopsInProgress: number
  ioTime: number
  weightedIO: number
  name: string
  serialNumber: string
  label: string
}

interface PartitionStat {
  device: string
  mountpoint: string
  fstype: string
  opts: string
}

interface DiskUsageStat {
  path: string
  fstype: string
  total: number
  free: number
  used: number
  usedPercent: number
  inodesTotal: number
  inodesUsed: number
  inodesFree: number
  inodesUsedPercent: number
}

export interface Realtime {
  cpus: CpuInfoStat[]
  percent: number
  percents: number[]
  load: LoadAvgStat
  host: HostInfoStat
  mem: VirtualMemoryStat
  swap: SwapMemoryStat
  net: IOCountersStat[]
  disk_io: DiskIOCountersStat[]
  disk: PartitionStat[]
  disk_usage: DiskUsageStat[]
}

export interface SystemInfo {
  procs: number
  hostname: string
  panel_version: string
  kernel_arch: string
  kernel_version: string
  os_name: string
  boot_time: number
  uptime: number
  nets: any[]
  disks: any[]
}

export interface CountInfo {
  website: number
  database: number
  ftp: number
  cron: number
}

export interface HomeApp {
  description: string
  icon: string
  name: string
  slug: string
  version: string
}

export interface VersionDownload {
  url: string
  arch: string
  checksum: string
}

export interface Version {
  created_at: string
  updated_at: string
  type: string
  version: string
  description: string
  downloads: VersionDownload[]
}
