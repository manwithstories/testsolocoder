import request from '@/utils/request'
import { LoginRequest, LoginResponse, User } from '@/types'

export const login = (data: LoginRequest) => {
  return request.post<LoginResponse>('/auth/login', data)
}

export const register = (data: any) => {
  return request.post('/auth/register', data)
}

export const getProfile = () => {
  return request.get<User>('/auth/profile')
}

export const logout = () => {
  return request.post('/auth/logout')
}

export const updateProfile = (data: Partial<User>) => {
  return request.put<User>('/auth/profile', data)
}

export const changePassword = (data: { oldPassword: string; newPassword: string }) => {
  return request.post('/auth/change-password', data)
}
