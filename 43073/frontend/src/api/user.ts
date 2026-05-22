import { request } from './request'

export interface LoginData {
  username: string
  password: string
}

export interface RegisterData {
  username: string
  password: string
  email?: string
  phone?: string
}

export interface LoginResponse {
  token: string
  userId: number
  username: string
  role: string
}

export const login = (data: LoginData) => {
  return request<LoginResponse>({
    url: '/auth/login',
    method: 'post',
    data
  })
}

export const register = (data: RegisterData) => {
  return request({
    url: '/auth/register',
    method: 'post',
    data
  })
}

export const getCurrentUser = () => {
  return request({
    url: '/users/me',
    method: 'get'
  })
}

export const getUserList = (params: any) => {
  return request({
    url: '/users',
    method: 'get',
    params
  })
}

export const updateUser = (id: number, data: any) => {
  return request({
    url: `/users/${id}`,
    method: 'put',
    data
  })
}

export const deleteUser = (id: number) => {
  return request({
    url: `/users/${id}`,
    method: 'delete'
  })
}
