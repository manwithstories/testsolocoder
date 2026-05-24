import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios'
import { ApiResponse } from '../types'

const apiClient: AxiosInstance = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error)
  }
)

export const apiGet = async <T = any>(url: string, params?: any): Promise<ApiResponse<T>> => {
  return apiClient.get(url, { params })
}

export const apiPost = async <T = any>(url: string, data?: any): Promise<ApiResponse<T>> => {
  return apiClient.post(url, data)
}

export const apiPut = async <T = any>(url: string, data?: any): Promise<ApiResponse<T>> => {
  return apiClient.put(url, data)
}

export const apiDelete = async <T = any>(url: string): Promise<ApiResponse<T>> => {
  return apiClient.delete(url)
}

export const apiUpload = async <T = any>(url: string, formData: FormData): Promise<ApiResponse<T>> => {
  return apiClient.post(url, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

export const apiDownload = async (url: string, params?: any): Promise<Blob> => {
  const response = await axios.get(url, {
    baseURL: '/api',
    params,
    responseType: 'blob',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  })
  return response.data
}

export default apiClient
