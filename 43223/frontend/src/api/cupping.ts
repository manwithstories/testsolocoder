import { request } from './client'
import {
  CuppingScore,
  PaginatedData,
  CreateCuppingScoreRequest,
} from '@/types'

export const cuppingApi = {
  list: (params?: { page?: number; page_size?: number; product_id?: string; user_id?: string }) =>
    request.get<PaginatedData<CuppingScore>>('/cupping/scores', { params }),

  get: (id: number) =>
    request.get<CuppingScore>(`/cupping/scores/${id}`),

  create: (data: CreateCuppingScoreRequest) =>
    request.post<CuppingScore>('/cupping/scores', data),

  update: (id: number, data: CreateCuppingScoreRequest) =>
    request.put(`/cupping/scores/${id}`, data),

  delete: (id: number) =>
    request.delete(`/cupping/scores/${id}`),

  getProductStats: (productId: string) =>
    request.get('/cupping/stats', { params: { product_id: productId } }),

  getTrend: (params?: { product_id?: string; days?: number }) =>
    request.get('/cupping/trend', { params }),

  getMyScores: () =>
    request.get<CuppingScore[]>('/cupping/my-scores'),
}
