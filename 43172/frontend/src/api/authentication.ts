import { request } from './index'
import type { Authentication, Review, ApiResponse } from '@/types'

export const authServiceApi = {
  listAuthentications(params?: {
    page?: number
    page_size?: number
    status?: string
  }) {
    return request.get<{ list: Authentication[]; total: number }>('/authentications', { params })
  },

  getAuthentication(id: number) {
    return request.get<Authentication>(`/authentications/${id}`)
  },

  getAuthenticationByOrder(orderId: number) {
    return request.get<Authentication>(`/authentications/order/${orderId}`)
  },

  createAuthentication(data: { order_id: number }) {
    return request.post<Authentication>('/authentications', data)
  },

  acceptAuthentication(id: number) {
    return request.post<Authentication>(`/authentications/${id}/accept`)
  },

  completeAuthentication(id: number, data: {
    result: string
    report_file?: string
    report_content?: string
    authenticator_notes?: string
  }) {
    return request.post<Authentication>(`/authentications/${id}/complete`, data)
  },

  rejectAuthentication(id: number, data: { reason: string }) {
    return request.post<Authentication>(`/authentications/${id}/reject`, data)
  },

  cancelAuthentication(id: number) {
    return request.post(`/authentications/${id}/cancel`)
  },

  downloadReport(id: number) {
    return `/api/v1/authentications/${id}/report/download`
  }
}

export const reviewApi = {
  listReviews(params?: {
    page?: number
    page_size?: number
    reviewee_id?: number
    reviewer_id?: number
    min_rating?: number
  }) {
    return request.get<{ list: Review[]; total: number }>('/reviews', { params })
  },

  getReview(id: number) {
    return request.get<Review>(`/reviews/${id}`)
  },

  createReview(data: {
    order_id: number
    reviewee_id: number
    rating: number
    content?: string
    images?: string
    is_anonymous?: boolean
  }) {
    return request.post<Review>('/reviews', data)
  },

  getUserAverageRating(userId: number) {
    return request.get<{ user_id: number; average_rating: number }>(`/reviews/user/${userId}/rating`)
  }
}

export default authServiceApi
