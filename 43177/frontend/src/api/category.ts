import { request } from './index'
import type { ApiResponse, Category, ServiceItem } from '@/types'

export const categoryApi = {
  getCategories: () => request.get<ApiResponse<Category[]>>('/categories'),

  getCategoryDetail: (id: number) =>
    request.get<ApiResponse<Category>>(`/categories/${id}`),

  createCategory: (data: { name: string; code: string; icon?: string; sort?: number }) =>
    request.post<ApiResponse<Category>>('/admin/categories', data),

  updateCategory: (id: number, data: { name?: string; icon?: string; sort?: number; status?: boolean }) =>
    request.put<ApiResponse>(`/admin/categories/${id}`, data),

  deleteCategory: (id: number) =>
    request.delete<ApiResponse>(`/admin/categories/${id}`)
}

export const serviceItemApi = {
  getServiceItems: (params?: { category_id?: number }) =>
    request.get<ApiResponse<ServiceItem[]>>('/service-items', { params }),

  getServiceItemDetail: (id: number) =>
    request.get<ApiResponse<ServiceItem>>(`/service-items/${id}`),

  createServiceItem: (data: {
    category_id: number
    name: string
    description?: string
    min_price: number
    max_price: number
    estimated_time: number
    image?: string
    sort?: number
  }) => request.post<ApiResponse<ServiceItem>>('/admin/service-items', data),

  updateServiceItem: (id: number, data: {
    name?: string
    description?: string
    min_price?: number
    max_price?: number
    estimated_time?: number
    image?: string
    sort?: number
    status?: boolean
  }) => request.put<ApiResponse>(`/admin/service-items/${id}`, data),

  deleteServiceItem: (id: number) =>
    request.delete<ApiResponse>(`/admin/service-items/${id}`)
}
