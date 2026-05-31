import request from './index'
import type { ApiResponse } from './index'

export interface LoginPayload {
  username: string
  password: string
}

export interface RegisterPayload {
  username: string
  password: string
  email: string
  phone: string
  role: 'customer' | 'staff' | 'company'
}

export interface AuthUser {
  id: number
  username: string
  role: string
  token: string
}

export function login(payload: LoginPayload) {
  return request.post<ApiResponse<AuthUser>>('/auth/login', payload)
}

export function register(payload: RegisterPayload) {
  return request.post<ApiResponse<AuthUser>>('/auth/register', payload)
}

export function logout() {
  return request.post<ApiResponse<null>>('/auth/logout')
}
