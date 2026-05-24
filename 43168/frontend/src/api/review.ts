import request from '@/utils/request'
import type { PaginationParams, PaginationResult } from '@/types'

export interface Review {
  id: number
  orderId?: number
  orderNo?: string
  productId?: number
  productName?: string
  productScore: number
  serviceScore: number
  content: string
  images?: string[]
  ownerId?: number
  ownerName?: string
  status: number
  reply?: string
  replyAt?: string
  createdAt: string
  updatedAt: string
}

export interface ListReviewParams extends PaginationParams {
  orderNo?: string
  status?: number
  keyword?: string
}

export interface ReviewFormData {
  id?: number
  orderId?: number
  productId?: number
  productScore: number
  serviceScore: number
  content: string
  images?: string[]
}

export const listReviews = (params: ListReviewParams) => {
  return request.get<any, PaginationResult<Review>>('/reviews', { params })
}

export const getReview = (id: number | string) => {
  return request.get<any, Review>(`/reviews/${id}`)
}

export const createReview = (data: ReviewFormData) => {
  return request.post<any, Review>('/reviews', data)
}

export const updateReview = (id: number | string, data: ReviewFormData) => {
  return request.put<any, Review>(`/reviews/${id}`, data)
}

export const deleteReview = (id: number | string) => {
  return request.delete<any, void>(`/reviews/${id}`)
}

export const replyReview = (id: number | string, reply: string) => {
  return request.post<any, void>(`/reviews/${id}/reply`, { reply })
}
