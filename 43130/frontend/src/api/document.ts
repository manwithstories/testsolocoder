import { request } from '@/utils/request'
import { ApiResponse, Document } from '@/types'

export const documentApi = {
  getList: (weddingId: number, params?: { category?: string }) => 
    request.get<ApiResponse<Document[]>>(`/weddings/${weddingId}/documents`, { params }),
  
  getById: (weddingId: number, id: number) => 
    request.get<ApiResponse<{ document: Document; versions: Document[] }>>(`/weddings/${weddingId}/documents/${id}`),
  
  upload: (weddingId: number, file: File, data?: { category?: string; notes?: string }, onProgress?: (progress: number) => void) => {
    const formData = new FormData()
    formData.append('file', file)
    if (data?.category) formData.append('category', data.category)
    if (data?.notes) formData.append('notes', data.notes)
    return request.upload(`/weddings/${weddingId}/documents/upload`, file, onProgress)
  },
  
  download: (weddingId: number, id: number) => 
    request.get(`/weddings/${weddingId}/documents/${id}/download`, { responseType: 'blob' }),
  
  update: (weddingId: number, id: number, data: Partial<Document>) => 
    request.put<ApiResponse<Document>>(`/weddings/${weddingId}/documents/${id}`, data),
  
  delete: (weddingId: number, id: number) => 
    request.delete<ApiResponse>(`/weddings/${weddingId}/documents/${id}`),
  
  uploadVersion: (weddingId: number, id: number, file: File, onProgress?: (progress: number) => void) => 
    request.upload(`/weddings/${weddingId}/documents/${id}/version`, file, onProgress),
  
  getCategories: (weddingId: number) => 
    request.get<ApiResponse<string[]>>(`/weddings/${weddingId}/documents/categories`)
}
