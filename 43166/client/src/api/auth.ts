import request from '@/utils/request'
import {
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  User
} from '@/types'

export const authApi = {
  login: (data: LoginRequest) => {
    return request.post<any, LoginResponse>('/auth/login', data)
  },

  register: (data: RegisterRequest) => {
    return request.post<any, User>('/auth/register', data)
  },

  logout: () => {
    return request.post<any, null>('/auth/logout')
  },

  getProfile: () => {
    return request.get<any, User>('/user/profile')
  },

  updateProfile: (data: Partial<User>) => {
    return request.put<any, null>('/user/profile', data)
  },

  changePassword: (oldPassword: string, newPassword: string) => {
    return request.put<any, null>('/user/password', { oldPassword, newPassword })
  }
}
