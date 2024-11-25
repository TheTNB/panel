import { DateTime } from 'luxon'

/**
 * 此处定义的是全局常量，启动或打包后将添加到 window 中
 * https://vitejs.cn/config/#define
 */

// 项目构建时间
const _BUILD_TIME_ = JSON.stringify(DateTime.now().toFormat('yyyy-MM-dd HH:mm:ss'))

export const viteDefine = {
  _BUILD_TIME_
}
