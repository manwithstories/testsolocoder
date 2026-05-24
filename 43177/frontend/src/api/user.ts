import { request } from './index'
import type { User, LoginData, ApiResponse, PaginationData, TechnicianProfile, Review } from '@/types'

export const authApi = {
  register: (data: {
    username: string
    password: string
    phone: string
    email?: string
    real_name?: string
    role: string
  }) => request.post<ApiResponse>('/auth/register', data),

  login: (data: { username: string; password: string }) =>
    request.post<ApiResponse<LoginData>>('/auth/login', data)
}

export const userApi = {
  getProfile: () => request.get<ApiResponse<User>>('/user/profile'),

  updateProfile: (data: {
    phone?: string
    email?: string
    real_name?: string
    avatar?: string
    address?: string
    longitude?: number
    latitude?: number
  }) => request.put<ApiResponse>('/user/profile', data),

  submitCertificate: (data: {
    certificate_image: string
    certificate_no: string
    specialty?: string
    experience_years?: number
  }) => request.post<ApiResponse>('/user/certificate', data),

  getTechnicians: (params?: { page?: number; page_size?: number; specialty?: string }) =>
    request.get<ApiResponse<PaginationData<any>>>('/technicians', { params }),

  getTechnicianDetail: (id: number) =>
    request.get<ApiResponse>(`/technicians/${id}`),

  getReviewList: (params?: { technician_id?: number; page?: number; page_size?: number }) =>
    request.get<ApiResponse<PaginationData<Review>>>('/reviews', { params }),

  createReview: (data: { order_id: number; rating: number; content?: string; images?: string }) =>
    request.post<ApiResponse>('/reviews', data),

  replyReview: (id: number, data: { reply: string }) =>
    request.post<ApiResponse>(`/reviews/${id}/reply`, data)
}
