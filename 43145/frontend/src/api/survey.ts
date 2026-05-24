import { get, post, put, del } from '@/utils/request'
import type { Survey, PaginatedData } from '@/types'

export const surveyApi = {
  list(params: { page?: number; page_size?: number; status?: number; category?: string; keyword?: string; sort_by?: string; sort_order?: string }) {
    return get<PaginatedData<Survey>>('/surveys', { params })
  },

  listAll(params: { page?: number; page_size?: number; status?: number; category?: string; keyword?: string; sort_by?: string; sort_order?: string }) {
    return get<PaginatedData<Survey>>('/surveys/all', { params })
  },

  getById(id: number) {
    return get<Survey>(`/surveys/${id}`)
  },

  create(data: any) {
    return post<Survey>('/surveys', data)
  },

  update(id: number, data: any) {
    return put(`/surveys/${id}`, data)
  },

  remove(id: number) {
    return del(`/surveys/${id}`)
  },

  publish(id: number, data?: { start_time?: string; end_time?: string }) {
    return post(`/surveys/${id}/publish`, data)
  },

  close(id: number) {
    return post(`/surveys/${id}/close`)
  },

  copy(id: number, data?: { title?: string }) {
    return post<Survey>(`/surveys/${id}/copy`, data)
  }
}
