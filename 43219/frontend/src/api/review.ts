import request from './index'
import type { ApiResponse } from './index'

export interface Review {
  id: number
  order_id: number
  staff_id: number
  customer_id: number
  rating: number
  content?: string
  images?: string
  created_at: string
  staff?: any
}

export function createReview(payload: {
  order_id: number
  staff_id: number
  rating: number
  content?: string
  images?: string
}) {
  return request.post<ApiResponse<Review>>('/reviews', payload)
}

export function listReviews(params?: Record<string, string | number>) {
  return request.get<ApiResponse<Review[]>>('/reviews', { params })
}

export function myReviews() {
  return request.get<ApiResponse<Review[]>>('/reviews/mine')
}
