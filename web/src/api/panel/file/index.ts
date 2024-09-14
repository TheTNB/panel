import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // 创建文件/文件夹
  create: (path: string, dir: boolean): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/create', { path, dir }),
  // 获取文件内容
  content: (path: string): Promise<AxiosResponse<any>> =>
    request.get('/panel/file/content', { params: { path } }),
  // 保存文件
  save: (path: string, content: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/save', { path, content }),
  // 删除文件
  delete: (path: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/delete', { path }),
  // 上传文件
  upload: (path: string, formData: FormData, onProgress: any): Promise<AxiosResponse<any>> => {
    formData.append('path', path)
    return request.post('/panel/file/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent: any) => {
        onProgress({ percent: Math.ceil((progressEvent.loaded / progressEvent.total) * 100) })
      }
    })
  },
  // 移动文件
  move: (source: string, target: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/move', { source, target }),
  // 复制文件
  copy: (source: string, target: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/copy', { source, target }),
  // 远程下载
  remoteDownload: (url: string, path: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/remoteDownload', { url, path }),
  // 获取文件信息
  info: (path: string): Promise<AxiosResponse<any>> =>
    request.get('/panel/file/info', { params: { path } }),
  // 修改文件权限
  permission: (
    path: string,
    mode: string,
    owner: string,
    group: string
  ): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/permission', { path, mode, owner, group }),
  // 压缩文件
  archive: (paths: string[], file: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/archive', { paths, file }),
  // 解压文件
  unArchive: (file: string, path: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/unArchive', { file, path }),
  // 搜索文件
  search: (keyword: string): Promise<AxiosResponse<any>> =>
    request.post('/panel/file/search', { keyword }),
  // 获取文件列表
  list: (path: string, page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/file/list', { params: { path, page, limit } })
}
