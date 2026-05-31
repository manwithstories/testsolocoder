import request from './request'
import type { LoginRequest, RegisterRequest, UpdateProfileRequest, LoginResponse, User, ApiResponse, PaginatedResponse } from '@/types/user'

export const loginApi = (data: LoginRequest) => {
  return request.post<LoginResponse>('/auth/login', data)
}

export const registerApi = (data: RegisterRequest) => {
  return request.post<ApiResponse>('/auth/register', data)
}

export const getProfileApi = () => {
  return request.get<User>('/auth/profile')
}

export const updateProfileApi = (data: UpdateProfileRequest) => {
  return request.put<User>('/auth/profile', data)
}

export const changePasswordApi = (data: { current_password: string; new_password: string }) => {
  return request.put('/auth/password', data)
}

export const getUsersApi = (params?: { role?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<User>>('/admin/users', { params })
}

export const getUserByIdApi = (id: string) => {
  return request.get<User>(`/admin/users/${id}`)
}

export const toggleUserStatusApi = (id: string) => {
  return request.put(`/admin/users/${id}/toggle-status`)
}
