import request from './request'
import type { VoyageLog, CreateVoyageLogRequest, PaginatedResponse } from '@/types/voyage'

export const getVoyageLogsApi = (params?: { rental_id?: string; ship_id?: string; start_date?: string; end_date?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<VoyageLog>>('/voyage-logs', { params })
}

export const getVoyageLogApi = (id: string) => {
  return request.get<VoyageLog>(`/voyage-logs/${id}`)
}

export const createVoyageLogApi = (data: CreateVoyageLogRequest) => {
  return request.post<VoyageLog>('/voyage-logs', data)
}

export const updateVoyageLogApi = (id: string, data: Partial<CreateVoyageLogRequest>) => {
  return request.put<VoyageLog>(`/voyage-logs/${id}`, data)
}

export const deleteVoyageLogApi = (id: string) => {
  return request.delete(`/voyage-logs/${id}`)
}

export const exportVoyageLogsApi = (params: { rental_id: string; start_date?: string; end_date?: string; format: 'pdf' | 'csv' }) => {
  return request.get('/voyage-logs/export', { params, responseType: 'blob' })
}
