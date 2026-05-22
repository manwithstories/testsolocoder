import { request } from '@/utils/request'
import { ApiResponse, PaginatedData, Vendor } from '@/types'

export const vendorApi = {
  create: (data: Partial<Vendor>) => 
    request.post<ApiResponse<Vendor>>('/vendors', data),
  
  getList: (params?: { search?: string; category?: string; wedding_id?: number; page?: number; page_size?: number }) => 
    request.get<ApiResponse<PaginatedData<Vendor>>>('/vendors', { params }),
  
  getById: (id: number) => 
    request.get<ApiResponse<{ vendor: Vendor; reviews: any[] }>>(`/vendors/${id}`),
  
  update: (id: number, data: Partial<Vendor>) => 
    request.put<ApiResponse<Vendor>>(`/vendors/${id}`, data),
  
  delete: (id: number) => 
    request.delete<ApiResponse>(`/vendors/${id}`),
  
  addReview: (id: number, data: { rating: number; content?: string }) => 
    request.post<ApiResponse>(`/vendors/${id}/reviews`, data),
  
  getCategories: () => 
    request.get<ApiResponse<string[]>>('/vendors/categories')
}
