export interface CronTask {
  id: number
  name: string
  status: boolean
  type: string
  time: string
  shell: string
  log: string
  created_at: string
  updated_at: string
}
