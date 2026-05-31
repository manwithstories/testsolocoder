import { request } from './client'
import {
  RoastingRecord,
  PaginatedData,
  CreateRoastingRecordRequest,
} from '@/types'

export const roastingApi = {
  list: (params?: {
    page?: number
    page_size?: number
    product_id?: string
    roaster_id?: string
    batch_number?: string
  }) => request.get<PaginatedData<RoastingRecord>>('/roasting/records', { params }),

  get: (id: number) =>
    request.get<RoastingRecord>(`/roasting/records/${id}`),

  create: (data: CreateRoastingRecordRequest) =>
    request.post<RoastingRecord>('/roasting/records', data),

  update: (id: number, data: CreateRoastingRecordRequest) =>
    request.put<RoastingRecord>(`/roasting/records/${id}`, data),

  delete: (id: number) =>
    request.delete(`/roasting/records/${id}`),

  compare: (recordIds: number[]) =>
    request.post('/roasting/compare', { record_ids: recordIds }),

  getStats: (params?: { roaster_id?: string; product_id?: string }) =>
    request.get('/roasting/stats', { params }),
}
