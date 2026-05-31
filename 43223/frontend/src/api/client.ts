import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ApiResponse } from '@/types'

const API_BASE_URL = '/api'

const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const data = response.data
    if (data.code === 0) {
      return response
    }
    return Promise.reject(new Error(data.message || '请求失败'))
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    const message = error.response?.data?.message || error.message || '网络错误'
    return Promise.reject(new Error(message))
  }
)

export const request = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) =>
    api.get<any, ApiResponse<T>>(url, config).then((res) => res.data),
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    api.post<any, ApiResponse<T>>(url, data, config).then((res) => res.data),
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    api.put<any, ApiResponse<T>>(url, data, config).then((res) => res.data),
  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    api.patch<any, ApiResponse<T>>(url, data, config).then((res) => res.data),
  delete: <T = any>(url: string, config?: AxiosRequestConfig) =>
    api.delete<any, ApiResponse<T>>(url, config).then((res) => res.data),
}

export default api
