import request from './request'
import type { Review, CreateReviewRequest, RespondToReviewRequest, PaginatedResponse } from '@/types/review'

export const getReviewsApi = (params: { target_type: 'ship' | 'dock'; target_id: string; min_rating?: number; max_rating?: number; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<Review>>('/reviews', { params })
}

export const getReviewApi = (id: string) => {
  return request.get<Review>(`/reviews/${id}`)
}

export const getMyReviewsApi = () => {
  return request.get<Review[]>('/my-reviews')
}

export const createReviewApi = (data: CreateReviewRequest) => {
  return request.post<Review>('/reviews', data)
}

export const respondToReviewApi = (id: string, data: RespondToReviewRequest) => {
  return request.post<Review>(`/reviews/${id}/respond`, data)
}

export const markReviewHelpfulApi = (id: string) => {
  return request.post(`/reviews/${id}/helpful`)
}

export const deleteReviewApi = (id: string) => {
  return request.delete(`/reviews/${id}`)
}
