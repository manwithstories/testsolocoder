import { http } from './request'
import {
  User,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  PaginatedResponse,
} from '@/types'

export const authApi = {
  login: (data: LoginRequest) =>
    http.post<LoginResponse>('/auth/login', data),
  register: (data: RegisterRequest) =>
    http.post<User>('/auth/register', data),
  refreshToken: () =>
    http.post<{ access_token: string }>('/auth/refresh'),
}

export const userApi = {
  getProfile: () =>
    http.get<User>('/users/profile'),
  updateProfile: (data: any) =>
    http.put<User>('/users/profile', data),
  getUserById: (id: string) =>
    http.get<User>(`/users/${id}`),
  getUsers: (params?: { page?: number; page_size?: number; role?: string }) =>
    http.get<PaginatedResponse<User>>('/users', { params }),
  verifyProfessional: (id: string, data: { status: string; note?: string }) =>
    http.put(`/users/${id}/verify`, data),
  getPendingVerifications: (params?: { page?: number; page_size?: number }) =>
    http.get<PaginatedResponse<User>>('/users/verifications/pending', { params }),
}
