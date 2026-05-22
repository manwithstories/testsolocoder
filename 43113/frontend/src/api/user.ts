import { request } from '@/utils/request'
import type {
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  User,
  PageResult,
  PointLog
} from '@/types'

export const authApi = {
  login: (data: LoginRequest) => {
    return request.post<LoginResponse>('/auth/login', data)
  },

  register: (data: RegisterRequest) => {
    return request.post<LoginResponse>('/auth/register', data)
  }
}

export const userApi = {
  getProfile: () => {
    return request.get<User>('/user/profile')
  },

  updateProfile: (data: Partial<User>) => {
    return request.put<User>('/user/profile', data)
  },

  getUserById: (id: number) => {
    return request.get<User>(`/public/users/${id}`)
  },

  getUserList: (params: { page?: number; pageSize?: number; keyword?: string }) => {
    return request.get<PageResult<User>>('/admin/users', { params })
  },

  updateUserStatus: (id: number, status: string) => {
    return request.put(`/admin/users/${id}/status`, { status })
  },

  getPointLogs: (params: { page?: number; pageSize?: number }) => {
    return request.get<PageResult<PointLog>>('/user/points/logs', { params })
  },

  applyExpert: (data: { field: string; description: string }) => {
    return request.post('/user/expert/apply', data)
  },

  getExpertApplications: (params: { page?: number; pageSize?: number; status?: string }) => {
    return request.get('/admin/expert/applications', { params })
  },

  reviewExpertApplication: (id: number, data: { status: string; remark?: string }) => {
    return request.put(`/admin/expert/applications/${id}`, data)
  }
}
