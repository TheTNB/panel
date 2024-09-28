export interface Jail {
  name: string
  enabled: boolean
  log_path: string
  max_retry: number
  find_time: number
  ban_time: number
}
