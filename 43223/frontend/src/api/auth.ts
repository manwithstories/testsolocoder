import { request } from './client'
import {
  User,
  LoginRequest,
  RegisterRequest,
  PaginatedData,
  UpdateUserRequest,
} from '@/types'

export const authApi = {
  register: (data: RegisterRequest) =>
    request.post<{ user: User; token: string }>('/auth/register', data),

  login: (data: LoginRequest) =>
    request.post<{ user: User; token: string }>('/auth/login', data),

  getProfile: () =>
    request.get<User>('/auth/profile'),

  updateProfile: (data: Partial<UpdateUserRequest>) =>
    request.put<User>('/auth/profile', data),
}

export const userApi = {
  list: (params?: { page?: number; page_size?: number; role?: string; status?: string; keyword?: string }) =>
    request.get<PaginatedData<User>>('/users', { params }),

  get: (id: number) =>
    request.get<User>(`/users/${id}`),

  updateStatus: (id: number, status: string) =>
    request.patch(`/users/${id}/status`, { status }),

  updateRole: (id: number, role: string) =>
    request.patch(`/users/${id}/role`, { role }),

  delete: (id: number) =>
    request.delete(`/users/${id}`),

  getActivity: (params?: { days?: number }) =>
    request.get('/users/activity', { params }),
}
