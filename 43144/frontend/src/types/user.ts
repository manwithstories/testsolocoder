import { RescueStation } from './rescue'

export interface User {
  id: number
  email: string
  name: string
  phone?: string
  role: string
  rescue_id?: number
  rescue?: RescueStation
  address?: string
  is_verified: boolean
  avatar?: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
  phone?: string
  role: string
  rescue_name?: string
}

export interface LoginResponse {
  token: string
  user_id: number
  email: string
  name: string
  role: string
  rescue_id?: number
}

export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}
