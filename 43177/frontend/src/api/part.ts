import { request } from './index'
import type { ApiResponse, Part, PartRequest, PaginationData, PartRequestItem, PartUsage } from '@/types'

export const partApi = {
  getParts: (params?: { page?: number; page_size?: number; category?: string; low_stock?: string }) =>
    request.get<ApiResponse<PaginationData<Part>>>('/admin/parts', { params }),

  getPartDetail: (id: number) =>
    request.get<ApiResponse<Part>>(`/admin/parts/${id}`),

  createPart: (data: {
    name: string
    code: string
    category?: string
    description?: string
    price: number
    stock?: number
    min_stock?: number
    image?: string
  }) => request.post<ApiResponse<Part>>('/admin/parts', data),

  updatePart: (id: number, data: {
    name?: string
    category?: string
    description?: string
    price?: number
    stock?: number
    min_stock?: number
    image?: string
    status?: boolean
  }) => request.put<ApiResponse>(`/admin/parts/${id}`, data),

  deletePart: (id: number) =>
    request.delete<ApiResponse>(`/admin/parts/${id}`),

  createPartRequest: (data: {
    items: { part_id: number; quantity: number }[]
    remark?: string
  }) => request.post<ApiResponse>('/part-requests', data),

  getPartRequests: (params?: { page?: number; page_size?: number; status?: string }) =>
    request.get<ApiResponse<PaginationData<PartRequest>>>('/part-requests', { params }),

  getPartRequestDetail: (id: number) =>
    request.get<ApiResponse<{ request: PartRequest; items: PartRequestItem[] }>>(`/part-requests/${id}`),

  approvePartRequest: (id: number) =>
    request.post<ApiResponse>(`/admin/part-requests/${id}/approve`),

  rejectPartRequest: (id: number, data: { remark: string }) =>
    request.post<ApiResponse>(`/admin/part-requests/${id}/reject`, data),

  shipPartRequest: (id: number) =>
    request.post<ApiResponse>(`/admin/part-requests/${id}/ship`),

  receivePartRequest: (id: number) =>
    request.post<ApiResponse>(`/part-requests/${id}/receive`),

  usePart: (data: { order_id: number; part_id: number; quantity: number }) =>
    request.post<ApiResponse>('/parts/use', data)
}
