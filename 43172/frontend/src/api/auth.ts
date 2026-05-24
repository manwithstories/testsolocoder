import { request } from './index'
import type { User, LoginResponse, ApiResponse } from '@/types'

export const authApi = {
  register(data: {
    username: string
    email: string
    password: string
    phone?: string
    real_name?: string
    role: string
  }) {
    return request.post<User>('/auth/register', data)
  },

  registerAuthenticator(data: {
    username: string
    email: string
    password: string
    phone?: string
    real_name?: string
    role: string
    license_number: string
    license_file: string
    certifications?: string
    specialties?: string
  }) {
    return request.post<User>('/auth/register/authenticator', data)
  },

  login(data: {
    username: string
    password: string
  }) {
    return request.post<LoginResponse>('/auth/login', data)
  },

  getProfile() {
    return request.get<User>('/users/profile')
  },

  updateProfile(data: Record<string, any>) {
    return request.put<User>('/users/profile', data)
  }
}

export default authApi
