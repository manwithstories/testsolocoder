import { request } from './client'
import {
  Product,
  PaginatedData,
  CreateProductRequest,
  UpdateProductRequest,
} from '@/types'

export const productApi = {
  list: (params?: {
    page?: number
    page_size?: number
    status?: string
    keyword?: string
    origin?: string
    roast_level?: string
    process_method?: string
    min_price?: string
    max_price?: string
    sort_by?: string
    sort_order?: string
    roaster_id?: string
  }) => request.get<PaginatedData<Product>>('/products', { params }),

  get: (id: number) =>
    request.get<Product>(`/products/${id}`),

  create: (data: CreateProductRequest) =>
    request.post<Product>('/products', data),

  update: (id: number, data: UpdateProductRequest) =>
    request.put<Product>(`/products/${id}`, data),

  updateStatus: (id: number, status: string) =>
    request.patch(`/products/${id}/status`, { status }),

  delete: (id: number) =>
    request.delete(`/products/${id}`),

  uploadImage: (id: number, formData: FormData) =>
    request.post(`/products/${id}/images`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  deleteImage: (productId: number, imageId: number) =>
    request.delete(`/products/${productId}/images/${imageId}`),

  importCSV: (formData: FormData) =>
    request.post('/products/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),

  getOrigins: () =>
    request.get<string[]>('/origins'),

  getRoastLevels: () =>
    request.get<string[]>('/roast-levels'),

  getProcessMethods: () =>
    request.get<string[]>('/process-methods'),
}
