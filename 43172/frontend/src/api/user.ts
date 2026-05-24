import { request } from './index'
import type { User, AuthenticatorProfile, ApiResponse } from '@/types'

export const userApi = {
  getUser(id: number) {
    return request.get<User>(`/users/${id}`)
  },

  listUsers(params?: {
    page?: number
    page_size?: number
    role?: string
    status?: string
  }) {
    return request.get<{ list: User[]; total: number }>('/users', { params })
  },

  listAuthenticators(params?: {
    page?: number
    page_size?: number
    status?: string
  }) {
    return request.get<{ list: AuthenticatorProfile[]; total: number }>('/admin/authenticators', { params })
  },

  approveAuthenticator(id: number) {
    return request.post(`/admin/authenticators/${id}/approve`)
  },

  rejectAuthenticator(id: number, data: { reason: string }) {
    return request.post(`/admin/authenticators/${id}/reject`, data)
  }
}

export default userApi
