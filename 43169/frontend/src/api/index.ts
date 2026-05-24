import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { message } from 'antd'
import { useAuthStore } from '@/store/authStore'

const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code !== 0) {
      message.error(res.message || '请求失败')
      if (res.code === 401) {
        useAuthStore.getState().logout()
        window.location.href = '/login'
      }
      return Promise.reject(new Error(res.message))
    }
    return res
  },
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
      window.location.href = '/login'
    } else {
      message.error(error.response?.data?.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export const apiGet = <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> =>
  api.get(url, config).then((res) => res.data)

export const apiPost = <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> =>
  api.post(url, data, config).then((res) => res.data)

export const apiPut = <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> =>
  api.put(url, data, config).then((res) => res.data)

export const apiDelete = <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> =>
  api.delete(url, config).then((res) => res.data)

export const apiUpload = <T = any>(url: string, file: File, fieldName = 'file'): Promise<T> => {
  const formData = new FormData()
  formData.append(fieldName, file)
  return api.post(url, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }).then((res) => res.data)
}

export default api
