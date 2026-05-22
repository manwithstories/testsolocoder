import request from './request'
import { Review, Complaint, PaginatedResponse, CreateReviewParams, CreateComplaintParams } from '@/types'

export const reviewApi = {
  getList: (params?: {
    page?: number
    page_size?: number
    service_item_id?: number
    provider_id?: number
    min_rating?: number
  }) => {
    return request.get<any, PaginatedResponse<Review>>('/reviews', { params })
  },

  getDetail: (id: number) => {
    return request.get<any, Review>(`/reviews/${id}`)
  },

  create: (params: CreateReviewParams) => {
    return request.post<any, Review>('/reviews', params)
  },

  reply: (id: number, params: { content: string }) => {
    return request.put<any, any>(`/reviews/${id}/reply`, params)
  },
}

export const complaintApi = {
  getList: (params?: { page?: number; page_size?: number; status?: string }) => {
    return request.get<any, PaginatedResponse<Complaint>>('/complaints', { params })
  },

  create: (params: CreateComplaintParams) => {
    return request.post<any, Complaint>('/complaints', params)
  },

  handle: (id: number, params: { status: string; result: string }) => {
    return request.put<any, any>(`/complaints/${id}/handle`, params)
  },
}
