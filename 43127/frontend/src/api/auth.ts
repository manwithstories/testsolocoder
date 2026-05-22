import request from '@/utils/request'
import type { User, ApiResponse, PaginationParams, PaginationResult } from '@/types'

export const login = (data: { username: string; password: string }) => {
  return request.post<any, ApiResponse<{ token: string; user: User }>>('/auth/login', data)
}

export const register = (data: {
  username: string
  password: string
  realName: string
  phone: string
  email?: string
  role: string
}) => {
  return request.post<any, ApiResponse<null>>('/auth/register', data)
}

export const getProfile = () => {
  return request.get<any, ApiResponse<User>>('/auth/profile')
}

export const updateProfile = (data: {
  realName?: string
  phone?: string
  email?: string
  avatar?: string
}) => {
  return request.put<any, ApiResponse<null>>('/auth/profile', data)
}

export const getUsers = (params: PaginationParams) => {
  return request.get<any, ApiResponse<PaginationResult<User>>>('/users', { params })
}

export const updateUserStatus = (id: number, status: number) => {
  return request.put<any, ApiResponse<null>>(`/users/${id}/status`, { status })
}

export const deleteUser = (id: number) => {
  return request.delete<any, ApiResponse<null>>(`/users/${id}`)
}
