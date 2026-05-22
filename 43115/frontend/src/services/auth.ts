import request from './request'
import { User, LoginParams, RegisterParams, Address, PaginatedResponse } from '@/types'

export const authApi = {
  login: (params: LoginParams) => {
    return request.post<any, { token: string; user_info: User }>('/auth/login', params)
  },

  register: (params: RegisterParams) => {
    return request.post<any, { token: string; user_info: User }>('/auth/register', params)
  },

  logout: () => {
    return request.post<any, any>('/auth/logout')
  },

  getCurrentUser: () => {
    return request.get<any, User>('/auth/me')
  },

  updateProfile: (params: Partial<User>) => {
    return request.put<any, any>('/auth/profile', params)
  },

  changePassword: (params: { old_password: string; new_password: string }) => {
    return request.put<any, any>('/auth/password', params)
  },
}

export const userApi = {
  getList: (params: { page?: number; page_size?: number; role?: string; status?: string; keyword?: string }) => {
    return request.get<any, PaginatedResponse<User>>('/admin/users', { params })
  },

  getDetail: (id: number) => {
    return request.get<any, User>(`/admin/users/${id}`)
  },

  toggleStatus: (id: number) => {
    return request.put<any, any>(`/admin/users/${id}/toggle-status`)
  },

  reviewCertification: (id: number, params: { approved: boolean; reject_reason?: string }) => {
    return request.put<any, any>(`/admin/users/${id}/review-certification`, params)
  },
}

export const addressApi = {
  getList: () => {
    return request.get<any, Address[]>('/addresses')
  },

  create: (params: Omit<Address, 'id' | 'user_id' | 'created_at' | 'updated_at'>) => {
    return request.post<any, Address>('/addresses', params)
  },

  update: (id: number, params: Partial<Address>) => {
    return request.put<any, any>(`/addresses/${id}`, params)
  },

  delete: (id: number) => {
    return request.delete<any, any>(`/addresses/${id}`)
  },

  setDefault: (id: number) => {
    return request.put<any, any>(`/addresses/${id}/default`)
  },
}

export const certificationApi = {
  submit: (params: {
    real_name: string
    id_card: string
    certification_desc?: string
    certifications: Array<{
      cert_type: string
      cert_name: string
      cert_number?: string
      cert_image?: string
    }>
  }) => {
    return request.post<any, any>('/certification', params)
  },
}
