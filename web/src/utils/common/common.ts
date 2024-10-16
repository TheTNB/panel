import { DateTime, Duration } from 'luxon'

/** 格式化时间，默认格式：yyyy-MM-dd HH:mm:ss */
export function formatDateTime(time: any, format = 'yyyy-MM-dd HH:mm:ss'): string {
  const dateTime = time ? DateTime.fromJSDate(new Date(time)) : DateTime.now()
  return dateTime.toFormat(format)
}

/** 格式化日期，默认格式：yyyy-MM-dd */
export function formatDate(date: any, format = 'yyyy-MM-dd') {
  return formatDateTime(date, format)
}

/** 格式化持续时间，转为 x天x小时x分钟x秒 */
export function formatDuration(seconds: number) {
  const duration = Duration.fromObject({ seconds }).shiftTo('days', 'hours', 'minutes', 'seconds')
  const days = duration.days
  const hours = duration.hours
  const minutes = duration.minutes
  const secs = duration.seconds

  return `${days}天${hours}小时${minutes}分钟${secs}秒`
}

/** 转时间戳 */
export function toTimestamp(time: any) {
  return DateTime.fromJSDate(new Date(time)).toSeconds()
}

/** 生成随机字符串 */
export function generateRandomString(length: number) {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    const randomIndex = Math.floor(Math.random() * characters.length)
    result += characters[randomIndex]
  }
  return result
}
