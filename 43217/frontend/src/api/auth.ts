import { request } from '@/utils/request'
import type { LoginRequest, LoginResponse, RegisterRequest, User } from '@/types'

export const login = (username: string, password: string): Promise<LoginResponse> => {
  return request.post<LoginResponse>('/auth/login', { username, password } as LoginRequest)
}

export const register = (data: RegisterRequest): Promise<User> => {
  return request.post<User>('/auth/register', data)
}

export const getUserInfo = (): Promise<User> => {
  return request.get<User>('/auth/userinfo')
}

export const logout = (): Promise<void> => {
  return request.post<void>('/auth/logout')
}
