import request from '@/utils/request'
import { User, PageParams, PageResult } from '@/types'

export const getUserList = (params: PageParams) => {
  return request.get<PageResult<User>>('/users', { params })
}

export const getUserById = (id: number) => {
  return request.get<User>(`/users/${id}`)
}

export const createUser = (data: Omit<User, 'id' | 'createdAt' | 'updatedAt'>) => {
  return request.post<User>('/users', data)
}

export const updateUser = (id: number, data: Partial<User>) => {
  return request.put<User>(`/users/${id}`, data)
}

export const deleteUser = (id: number) => {
  return request.delete(`/users/${id}`)
}

export const updateUserStatus = (id: number, status: boolean) => {
  return request.patch(`/users/${id}/status`, { status })
}

export const resetPassword = (id: number) => {
  return request.post(`/users/${id}/reset-password`)
}
