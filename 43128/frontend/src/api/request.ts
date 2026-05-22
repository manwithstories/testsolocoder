import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'
import type { ApiResp } from '@/types'

const http: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

http.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (err) => Promise.reject(err),
)

http.interceptors.response.use(
  (res: AxiosResponse<ApiResp<any>>) => {
    const data = res.data
    if (data && data.code !== undefined && data.code !== 0) {
      ElMessage.error(data.message || '请求失败')
      return Promise.reject(new Error(data.message))
    }
    return res
  },
  (err) => {
    const status = err?.response?.status
    if (status === 401) {
      localStorage.removeItem('token')
      router.push('/login')
    } else if (status >= 500) {
      ElMessage.error('服务器错误')
    } else {
      ElMessage.error(err?.response?.data?.message || '请求失败')
    }
    return Promise.reject(err)
  },
)

export default http

export function get<T = any>(url: string, params?: any) {
  return http.get<ApiResp<T>>(url, { params }).then((r) => r.data)
}
export function post<T = any>(url: string, data?: any) {
  return http.post<ApiResp<T>>(url, data).then((r) => r.data)
}
export function put<T = any>(url: string, data?: any) {
  return http.put<ApiResp<T>>(url, data).then((r) => r.data)
}
export function del<T = any>(url: string) {
  return http.delete<ApiResp<T>>(url).then((r) => r.data)
}
