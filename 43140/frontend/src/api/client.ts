import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'

const API_BASE_URL = '/api/v1'

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
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

export const api = {
  get: <T>(url: string, config?: AxiosRequestConfig) =>
    apiClient.get<unknown, T>(url, config),
  post: <T>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
    apiClient.post<unknown, T>(url, data, config),
  put: <T>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
    apiClient.put<unknown, T>(url, data, config),
  delete: <T>(url: string, config?: AxiosRequestConfig) =>
    apiClient.delete<unknown, T>(url, config),
  upload: <T>(url: string, formData: FormData) =>
    apiClient.post<unknown, T>(url, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),
}

export default apiClient
