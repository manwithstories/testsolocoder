import { apiGet, apiPost, apiPut } from './client'
import { LoginRequest, RegisterRequest, LoginResponse, User, ApiResponse } from '../types'

export const login = (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
  return apiPost<LoginResponse>('/auth/login', data)
}

export const register = (data: RegisterRequest): Promise<ApiResponse<User>> => {
  return apiPost<User>('/auth/register', data)
}

export const getProfile = (): Promise<ApiResponse<User>> => {
  return apiGet<User>('/profile')
}

export const updateProfile = (data: Partial<User>): Promise<ApiResponse<User>> => {
  return apiPut<User>('/profile', data)
}

export const getUsers = (params?: { page?: number; page_size?: number; role?: string }): Promise<ApiResponse<any>> => {
  return apiGet('/users', params)
}

export const getUserById = (id: number): Promise<ApiResponse<User>> => {
  return apiGet<User>(`/users/${id}`)
}

export const verifyUser = (id: number): Promise<ApiResponse<User>> => {
  return apiPut<User>(`/users/${id}/verify`)
}
