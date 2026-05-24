import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { ApiResponse } from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response: AxiosResponse) => {
    const data = response.data as ApiResponse;
    if (data.code === 0) {
      return response;
    }
    return Promise.reject(new Error(data.message || '请求失败'));
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const apiRequest = {
  get: <T = any>(url: string, config?: AxiosRequestConfig) =>
    api.get<ApiResponse<T>>(url, config).then((res) => res.data.data),

  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    api.post<ApiResponse<T>>(url, data, config).then((res) => res.data.data),

  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig) =>
    api.put<ApiResponse<T>>(url, data, config).then((res) => res.data.data),

  delete: <T = any>(url: string, config?: AxiosRequestConfig) =>
    api.delete<ApiResponse<T>>(url, config).then((res) => res.data.data),
};

export default api;
