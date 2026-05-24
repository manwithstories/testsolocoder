import { post, get, put } from '@/utils/request'
import type { ApiResponse, LoginRequest, RegisterRequest, User, ApplicantProfile, Company, PaginatedResponse } from '@/types'

export function login(data: LoginRequest) {
  return post<ApiResponse<{ user: User; token: string }>>('/auth/login', data)
}

export function register(data: RegisterRequest) {
  return post<ApiResponse<{ user: User; token: string }>>('/auth/register', data)
}

export function getProfile() {
  return get<ApiResponse<User>>('/users/profile')
}

export function updateProfile(data: ApplicantProfile) {
  return put<ApiResponse<ApplicantProfile>>('/users/profile', data)
}

export function getCompany() {
  return get<ApiResponse<Company>>('/users/company')
}

export function updateCompany(data: Company) {
  return put<ApiResponse<Company>>('/users/company', data)
}

export function listUsers(params?: { page?: number; page_size?: number; role?: string; status?: string }) {
  return get<ApiResponse<PaginatedResponse<User>>>('/admin/users', { params })
}

export function updateUserStatus(id: number, status: string) {
  return put<ApiResponse<any>>(`/admin/users/${id}/status`, { status })
}
