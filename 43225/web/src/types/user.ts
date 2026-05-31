export type UserRole = 'admin' | 'owner' | 'tenant'

export interface User {
  id: string
  username: string
  email: string
  role: UserRole
  full_name: string
  phone?: string
  avatar_url?: string
  company?: string
  address?: string
  city?: string
  country?: string
  timezone: string
  is_active: boolean
  is_email_verified: boolean
  last_login_at?: string
  created_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  role: UserRole
  full_name: string
  phone?: string
  company?: string
}

export interface UpdateProfileRequest {
  full_name?: string
  phone?: string
  company?: string
  address?: string
  city?: string
  country?: string
  timezone?: string
  avatar_url?: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
  error?: string
}

export interface PaginatedResponse<T> {
  code: number
  message: string
  data: T[]
  total: number
  page: number
  page_size: number
}
