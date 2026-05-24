import { request } from './index'
import type { ApiResponse, User, PaginationData, Review, TechnicianProfile } from '@/types'

export const adminApi = {
  getDashboard: () =>
    request.get<ApiResponse<any>>('/admin/dashboard'),

  getUsers: (params?: { page?: number; page_size?: number; role?: string; status?: string }) =>
    request.get<ApiResponse<PaginationData<User>>>('/admin/users', { params }),

  getUserDetail: (id: number) =>
    request.get<ApiResponse<any>>(`/admin/users/${id}`),

  updateUserStatus: (id: number, data: { status: string }) =>
    request.put<ApiResponse>(`/admin/users/${id}/status`, data),

  getTechnicianVerifyList: (params?: { page?: number; page_size?: number; status?: string }) =>
    request.get<ApiResponse<PaginationData<TechnicianProfile>>>('/admin/technicians/verify', { params }),

  verifyTechnician: (id: number, data: { is_verified: boolean; verify_remark?: string }) =>
    request.post<ApiResponse>(`/admin/technicians/${id}/verify`, data),

  getLowRatingReviews: (params?: { page?: number; page_size?: number }) =>
    request.get<ApiResponse<PaginationData<Review>>>('/admin/reviews/low-rating', { params }),

  interveneReview: (id: number, data: { note: string }) =>
    request.post<ApiResponse>(`/admin/reviews/${id}/intervene`, data)
}
