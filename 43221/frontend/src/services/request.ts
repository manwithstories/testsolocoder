import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ApiResponse } from '@/types'

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

const request: AxiosInstance = axios.create({
  baseURL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const res = response.data
    if (res.code !== 0) {
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export interface RequestConfig extends AxiosRequestConfig {}

export const http = {
  get: <T = any>(url: string, config?: RequestConfig) =>
    request.get<ApiResponse<T>>(url, config).then((res) => res.data.data as T),
  post: <T = any>(url: string, data?: any, config?: RequestConfig) =>
    request.post<ApiResponse<T>>(url, data, config).then((res) => res.data.data as T),
  put: <T = any>(url: string, data?: any, config?: RequestConfig) =>
    request.put<ApiResponse<T>>(url, data, config).then((res) => res.data.data as T),
  delete: <T = any>(url: string, config?: RequestConfig) =>
    request.delete<ApiResponse<T>>(url, config).then((res) => res.data.data as T),
}

export default request
