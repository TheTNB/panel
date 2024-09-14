import { request } from '@/utils'
import type { AxiosResponse } from 'axios'

export default {
  // CA 供应商列表
  caProviders: (): Promise<AxiosResponse<any>> => request.get('/panel/cert/caProviders'),
  // DNS 供应商列表
  dnsProviders: (): Promise<AxiosResponse<any>> => request.get('/panel/cert/dnsProviders'),
  // 证书算法列表
  algorithms: (): Promise<AxiosResponse<any>> => request.get('/panel/cert/algorithms'),
  // ACME 用户列表
  users: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/cert/users', { params: { page, limit } }),
  // ACME 用户详情
  userInfo: (id: number): Promise<AxiosResponse<any>> => request.get(`/panel/cert/users/${id}`),
  // ACME 用户添加
  userAdd: (data: any): Promise<AxiosResponse<any>> => request.post('/panel/cert/users', data),
  // ACME 用户更新
  userUpdate: (id: number, data: any): Promise<AxiosResponse<any>> =>
    request.put(`/panel/cert/users/${id}`, data),
  // ACME 用户删除
  userDelete: (id: number): Promise<AxiosResponse<any>> =>
    request.delete(`/panel/cert/users/${id}`),
  // DNS 记录列表
  dns: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/cert/dns', { params: { page, limit } }),
  // DNS 记录详情
  dnsInfo: (id: number): Promise<AxiosResponse<any>> => request.get(`/panel/cert/dns/${id}`),
  // DNS 记录添加
  dnsAdd: (data: any): Promise<AxiosResponse<any>> => request.post('/panel/cert/dns', data),
  // DNS 记录更新
  dnsUpdate: (id: number, data: any): Promise<AxiosResponse<any>> =>
    request.put(`/panel/cert/dns/${id}`, data),
  // DNS 记录删除
  dnsDelete: (id: number): Promise<AxiosResponse<any>> => request.delete(`/panel/cert/dns/${id}`),
  // 证书列表
  certs: (page: number, limit: number): Promise<AxiosResponse<any>> =>
    request.get('/panel/cert/certs', { params: { page, limit } }),
  // 证书详情
  certInfo: (id: number): Promise<AxiosResponse<any>> => request.get(`/panel/cert/certs/${id}`),
  // 证书添加
  certAdd: (data: any): Promise<AxiosResponse<any>> => request.post('/panel/cert/certs', data),
  // 证书更新
  certUpdate: (id: number, data: any): Promise<AxiosResponse<any>> =>
    request.put(`/panel/cert/certs/${id}`, data),
  // 证书删除
  certDelete: (id: number): Promise<AxiosResponse<any>> =>
    request.delete(`/panel/cert/certs/${id}`),
  // 签发
  obtain: (id: number): Promise<AxiosResponse<any>> => request.post(`/panel/cert/obtain`, { id }),
  // 续签
  renew: (id: number): Promise<AxiosResponse<any>> => request.post(`/panel/cert/renew`, { id }),
  // 获取 DNS 记录
  manualDNS: (id: number): Promise<AxiosResponse<any>> =>
    request.post(`/panel/cert/manualDNS`, { id }),
  // 部署
  deploy: (id: number, website_id: number): Promise<AxiosResponse<any>> =>
    request.post(`/panel/cert/deploy`, { id, website_id })
}
