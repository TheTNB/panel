import { decrypto, encrypto } from '@/utils'

export function setSession(key: string, value: unknown) {
  const json = encrypto(value)
  sessionStorage.setItem(key, json)
}

export function getSession<T>(key: string) {
  const json = sessionStorage.getItem(key)
  let data: T | null = null
  if (json) {
    data = decrypto(json)
  }
  return data
}

export function removeSession(key: string) {
  window.sessionStorage.removeItem(key)
}

export function clearSession() {
  window.sessionStorage.clear()
}
