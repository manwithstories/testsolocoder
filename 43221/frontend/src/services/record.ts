import { http } from './request'
import {
  ConsultRecord,
  Review,
  PaginatedResponse,
  CreateConsultRecordRequest,
  CreateReviewRequest,
} from '@/types'

export const recordApi = {
  createConsultRecord: (data: CreateConsultRecordRequest) =>
    http.post<ConsultRecord>('/records/consult', data),
  getConsultRecordById: (id: string) =>
    http.get<ConsultRecord>(`/records/${id}`),
  getClientConsultRecords: (params?: { page?: number; page_size?: number }) =>
    http.get<PaginatedResponse<ConsultRecord>>('/records/client/list', { params }),
  getProfessionalConsultRecords: (params?: { page?: number; page_size?: number }) =>
    http.get<PaginatedResponse<ConsultRecord>>('/records/professional/list', { params }),
  createReview: (data: CreateReviewRequest) =>
    http.post<Review>('/records/review', data),
  getProfessionalReviews: (params?: { page?: number; page_size?: number; status?: string }) =>
    http.get<PaginatedResponse<Review>>('/reviews/professional/list', { params }),
  getServiceReviews: (serviceId: string, params?: { page?: number; page_size?: number }) =>
    http.get<PaginatedResponse<Review>>(`/services/${serviceId}/reviews`, { params }),
  getPendingReviews: (params?: { page?: number; page_size?: number }) =>
    http.get<PaginatedResponse<Review>>('/reviews/pending', { params }),
  updateReviewStatus: (data: { review_id: string; status: string; reject_reason?: string }) =>
    http.put('/reviews/status', data),
  getProfessionalReviewStats: (professionalId: string) =>
    http.get<{ average_rating: number; total_reviews: number }>(`/reviews/professional/stats/${professionalId}`),
}

export { ConsultRecord, Review }
