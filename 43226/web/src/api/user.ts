import request from '@/utils/request'
import type {
  User,
  ApiResponse,
  PageResult
} from '@/types'

export interface UserQuery {
  page?: number
  page_size?: number
  keyword?: string
  role?: string
  status?: number
}

export const listUsers = (params?: UserQuery) => {
  return request.get<any, ApiResponse<PageResult<User>>>('/users', { params })
}

export const getUser = (id: number) => {
  return request.get<any, ApiResponse<User>>(`/users/${id}`)
}

export const createUser = (data: Partial<User> & { password?: string }) => {
  return request.post<any, ApiResponse<User>>('/users', data)
}

export const updateUser = (id: number, data: Partial<User>) => {
  return request.put<any, ApiResponse<void>>(`/users/${id}`, data)
}

export const deleteUser = (id: number) => {
  return request.delete<any, ApiResponse<void>>(`/users/${id}`)
}

export const updateUserRole = (id: number, role: string) => {
  return request.put<any, ApiResponse<void>>(`/users/${id}/role`, { role })
}

export const updateUserStatus = (id: number, status: number) => {
  return request.put<any, ApiResponse<void>>(`/users/${id}/status`, { status })
}

export const listGuides = () => {
  return request.get<any, ApiResponse<User[]>>('/users/guides/list')
}
