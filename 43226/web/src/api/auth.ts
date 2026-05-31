import request from '@/utils/request'
import type {
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  User,
  ApiResponse,
  PageResult
} from '@/types'

export const login = (data: LoginRequest) => {
  return request.post<any, ApiResponse<LoginResponse>>('/auth/login', data)
}

export const register = (data: RegisterRequest) => {
  return request.post<any, ApiResponse<LoginResponse>>('/auth/register', data)
}

export const getProfile = () => {
  return request.get<any, ApiResponse<User>>('/users/me')
}

export const updateProfile = (data: Partial<User>) => {
  return request.put<any, ApiResponse<void>>('/users/me', data)
}

export const listUsers = (params?: { page?: number; page_size?: number; keyword?: string }) => {
  return request.get<any, ApiResponse<PageResult<User>>>('/users', { params })
}

export const getUserById = (id: number) => {
  return request.get<any, ApiResponse<User>>(`/users/${id}`)
}

export const listGuides = () => {
  return request.get<any, ApiResponse<User[]>>('/users/guides/list')
}
