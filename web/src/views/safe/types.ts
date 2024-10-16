export interface FirewallRule {
  port_start: number
  port_end: number
  protocols: string[]
  address: string
  strategy: string
}
