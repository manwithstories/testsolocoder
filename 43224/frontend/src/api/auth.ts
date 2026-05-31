import request from '@/utils/request'

export interface LoginData {
  username: string
  password: string
}

export interface RegisterData {
  username: string
  password: string
  email: string
  phone?: string
  real_name?: string
  role: string
}

export function login(data: LoginData) {
  return request.post('/auth/login', data)
}

export function register(data: RegisterData) {
  return request.post('/auth/register', data)
}

export function getCurrentUser() {
  return request.get('/auth/me')
}

export function updateProfile(data: any) {
  return request.put('/auth/profile', data)
}

export function changePassword(data: any) {
  return request.put('/auth/password', data)
}

export function listUsers(params?: any) {
  return request.get('/users', { params })
}

export function updateUserStatus(id: number, status: string) {
  return request.put(`/users/${id}/status`, { status })
}
