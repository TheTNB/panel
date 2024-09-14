import path from 'node:path'

/**
 * * 项目根路径
 * @descrition 结尾不带/
 */
export function getRootPath() {
  return path.resolve(process.cwd())
}

/**
 * * 项目src路径
 * @param srcName src目录名称(默认: "src")
 * @descrition 结尾不带斜杠
 */
export function getSrcPath(srcName = 'src') {
  return path.resolve(getRootPath(), srcName)
}

/**
 * * 转换env配置
 * @param envOptions
 * @descrition boolean和数字类型转换
 */
export function convertEnv(envOptions: Record<string, any>): ViteEnv {
  const result: any = {}
  if (!envOptions) return result

  for (const envKey in envOptions) {
    let envVal = envOptions[envKey]
    if (['true', 'false'].includes(envVal)) envVal = envVal === 'true'

    if (['VITE_PORT'].includes(envKey)) envVal = +envVal

    result[envKey] = envVal
  }
  return result
}
