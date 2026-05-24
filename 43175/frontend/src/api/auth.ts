import request from './request'

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  expireAt: string
  userId: number
  email: string
  username: string
  avatar?: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
  phone?: string
}

export function login(data: LoginRequest): Promise<LoginResponse> {
  return request.post('/auth/login', data)
}

export function register(data: RegisterRequest): Promise<any> {
  return request.post('/auth/register', data)
}

export function refreshToken(token: string): Promise<LoginResponse> {
  return request.post('/auth/refresh', { token })
}

export function getProfile(): Promise<any> {
  return request.get('/user/profile')
}

export function updateProfile(data: any): Promise<any> {
  return request.put('/user/profile', data)
}

export function changePassword(data: any): Promise<any> {
  return request.put('/user/password', data)
}
