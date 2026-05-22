import { request } from '@/utils/request'
import { ApiResponse, PaginatedData, Guest, GuestTable } from '@/types'

export const guestApi = {
  create: (weddingId: number, data: Partial<Guest>) => 
    request.post<ApiResponse<Guest>>(`/weddings/${weddingId}/guests`, data),
  
  getList: (weddingId: number, params?: { search?: string; group?: string; rsvp_status?: string; page?: number; page_size?: number }) => 
    request.get<ApiResponse<PaginatedData<Guest>>>(`/weddings/${weddingId}/guests`, { params }),
  
  getById: (weddingId: number, id: number) => 
    request.get<ApiResponse<Guest>>(`/weddings/${weddingId}/guests/${id}`),
  
  update: (weddingId: number, id: number, data: Partial<Guest>) => 
    request.put<ApiResponse<Guest>>(`/weddings/${weddingId}/guests/${id}`, data),
  
  delete: (weddingId: number, id: number) => 
    request.delete<ApiResponse>(`/weddings/${weddingId}/guests/${id}`),
  
  updateRSVP: (weddingId: number, id: number, status: string) => 
    request.put<ApiResponse>(`/weddings/${weddingId}/guests/${id}/rsvp`, { rsvp_status: status }),
  
  import: (weddingId: number, file: File, onProgress?: (progress: number) => void) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.upload(`/weddings/${weddingId}/guests/import`, file, onProgress)
  },
  
  export: (weddingId: number) => 
    request.get(`/weddings/${weddingId}/guests/export`, { responseType: 'blob' }),
  
  createTable: (weddingId: number, data: Partial<GuestTable>) => 
    request.post<ApiResponse<GuestTable>>(`/weddings/${weddingId}/guests/tables`, data),
  
  getTables: (weddingId: number) => 
    request.get<ApiResponse<GuestTable[]>>(`/weddings/${weddingId}/guests/tables`),
  
  updateTable: (weddingId: number, id: number, data: Partial<GuestTable>) => 
    request.put<ApiResponse<GuestTable>>(`/weddings/${weddingId}/guests/tables/${id}`, data),
  
  deleteTable: (weddingId: number, id: number) => 
    request.delete<ApiResponse>(`/weddings/${weddingId}/guests/tables/${id}`),
  
  assignSeat: (weddingId: number, data: { guest_id: number; table_id: number; seat_number: number }) => 
    request.post<ApiResponse>(`/weddings/${weddingId}/guests/seat-assign`, data)
}
