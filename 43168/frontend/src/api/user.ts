import request from '@/utils/request'

export type UserRole = 'manufacturer' | 'designer' | 'owner'

export interface LoginParams {
  username: string
  password: string
  role?: UserRole
}

export interface RegisterParams {
  username: string
  password: string
  role: UserRole
  nickname?: string
  phone?: string
  email?: string
}

export interface UpdateProfileParams {
  nickname?: string
  avatar?: string
  phone?: string
  email?: string
}

export interface ChangePasswordParams {
  old_password: string
  new_password: string
}

export interface ListUsersParams {
  page?: number
  page_size?: number
  keyword?: string
  role?: UserRole
  status?: number
}

export interface UserInfo {
  id: number
  username: string
  role: UserRole
  nickname: string
  avatar: string
  phone: string
  email: string
  status: number
  created_at: string
  updated_at: string
}

export interface LoginResult {
  token: string
  user: UserInfo
}

export interface ListResult<T> {
  total: number
  page: number
  page_size: number
  list: T[]
}

export const login = (data: LoginParams) => {
  return request.post<any, LoginResult>('/auth/login', data)
}

export const register = (data: RegisterParams) => {
  return request.post<any, LoginResult>('/auth/register', data)
}

export const getUserInfo = () => {
  return request.get<any, UserInfo>('/users/profile')
}

export const updateProfile = (data: UpdateProfileParams) => {
  return request.put<any, UserInfo>('/users/profile', data)
}

export const changePassword = (data: ChangePasswordParams) => {
  return request.put<any, void>('/users/password', data)
}

export const listUsers = (params: ListUsersParams) => {
  return request.get<any, ListResult<UserInfo>>('/users/', { params })
}

export const getUserById = (id: number | string) => {
  return request.get<any, UserInfo>(`/users/${id}`)
}
