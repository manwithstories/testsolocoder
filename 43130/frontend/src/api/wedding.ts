import { request } from '@/utils/request'
import { ApiResponse, PaginatedData, Wedding } from '@/types'

export const weddingApi = {
  create: (data: Partial<Wedding>) => 
    request.post<ApiResponse<Wedding>>('/weddings', data),
  
  getList: (params?: { search?: string; status?: string; page?: number; page_size?: number }) => 
    request.get<ApiResponse<PaginatedData<Wedding>>>('/weddings', { params }),
  
  getById: (id: number) => 
    request.get<ApiResponse<Wedding>>(`/weddings/${id}`),
  
  update: (id: number, data: Partial<Wedding>) => 
    request.put<ApiResponse<Wedding>>(`/weddings/${id}`, data),
  
  delete: (id: number) => 
    request.delete<ApiResponse>(`/weddings/${id}`),
  
  updateStatus: (id: number, status: string) => 
    request.put<ApiResponse>(`/weddings/${id}/status`, { status })
}
