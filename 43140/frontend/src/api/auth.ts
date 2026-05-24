import { api } from './client'

export interface User {
  id: number
  email: string
  name: string
  phone?: string
  role: 'admin' | 'company' | 'jobseeker'
  status: string
  avatar?: string
  company?: Company
  jobseeker?: JobSeeker
  created_at: string
  updated_at: string
}

export interface Company {
  id: number
  user_id: number
  company_name: string
  industry?: string
  company_size?: string
  company_type?: string
  address?: string
  website?: string
  logo?: string
  description?: string
  created_at: string
}

export interface JobSeeker {
  id: number
  user_id: number
  birth_date?: string
  gender?: string
  location?: string
  current_title?: string
  current_company?: string
  experience?: string
  education?: string
  created_at: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PaginatedData<T> {
  data: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export const authApi = {
  login: (data: { email: string; password: string }) =>
    api.post<ApiResponse<{ token: string; user: User }>>('/auth/login', data),

  register: (data: {
    email: string
    password: string
    name: string
    phone?: string
    role: 'company' | 'jobseeker'
    company_name?: string
    industry?: string
  }) => api.post<ApiResponse<{ token: string; user: User }>>('/auth/register', data),

  getProfile: () =>
    api.get<ApiResponse<User>>('/auth/profile'),

  updateProfile: (data: { name?: string; phone?: string }) =>
    api.put<ApiResponse<User>>('/auth/profile', data),

  changePassword: (data: { current_password: string; new_password: string }) =>
    api.put<ApiResponse<null>>('/auth/password', data),
}
