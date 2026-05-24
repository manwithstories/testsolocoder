import api from './index'
import type { ApiResponse, LoginRequest, LoginResponse, RegisterRequest, User } from '@/types'

export const authApi = {
  login: (data: LoginRequest) => {
    return api.post<any, ApiResponse<LoginResponse>>('/auth/login', data)
  },

  register: (data: RegisterRequest) => {
    return api.post<any, ApiResponse<User>>('/auth/register', data)
  }
}
