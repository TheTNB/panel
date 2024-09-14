import dayjs from 'dayjs'

type Time = undefined | string | Date

/** 格式化时间，默认格式：YYYY-MM-DD HH:mm:ss */
export function formatDateTime(time: Time, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(time).format(format)
}

/** 格式化日期，默认格式：YYYY-MM-DD */
export function formatDate(date: Time = undefined, format = 'YYYY-MM-DD') {
  return formatDateTime(date, format)
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
