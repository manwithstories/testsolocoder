import request from './request'
import { ServiceCategory, ServiceItem, ServiceArea, PaginatedResponse } from '@/types'

export const serviceApi = {
  getCategories: () => {
    return request.get<any, ServiceCategory[]>('/services/categories')
  },

  createCategory: (params: { name: string; description?: string; icon?: string; sort_order?: number }) => {
    return request.post<any, ServiceCategory>('/admin/service-categories', params)
  },

  updateCategory: (id: number, params: Partial<ServiceCategory>) => {
    return request.put<any, any>(`/admin/service-categories/${id}`, params)
  },

  deleteCategory: (id: number) => {
    return request.delete<any, any>(`/admin/service-categories/${id}`)
  },

  getAreas: () => {
    return request.get<any, ServiceArea[]>('/services/areas')
  },

  createArea: (params: { province: string; city: string; district: string }) => {
    return request.post<any, ServiceArea>('/admin/service-areas', params)
  },

  getList: (params?: {
    page?: number
    page_size?: number
    category_id?: number
    min_price?: number
    max_price?: number
    min_rating?: number
    keyword?: string
  }) => {
    return request.get<any, PaginatedResponse<ServiceItem>>('/services', { params })
  },

  getDetail: (id: number) => {
    return request.get<any, ServiceItem>(`/services/${id}`)
  },

  create: (params: {
    category_id: number
    name: string
    description?: string
    images?: string
    base_price: number
    price_unit?: string
    min_duration?: number
    max_duration?: number
    service_area_ids?: number[]
  }) => {
    return request.post<any, ServiceItem>('/services', params)
  },

  update: (id: number, params: Partial<ServiceItem> & { service_area_ids?: number[] }) => {
    return request.put<any, any>(`/services/${id}`, params)
  },

  delete: (id: number) => {
    return request.delete<any, any>(`/services/${id}`)
  },

  getMyServices: () => {
    return request.get<any, ServiceItem[]>('/services/mine/list')
  },
}
