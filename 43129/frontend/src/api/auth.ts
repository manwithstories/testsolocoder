import request from '@/utils/request'
import type { ApiResponse, User, Customer, PageResult, PageParams } from '@/types'

export const login = (data: { phone: string; password: string }) => {
  return request.post<ApiResponse<{ token: string; user: User }>>('/auth/login', data)
}

export const register = (data: { phone: string; password: string; nickname?: string; role?: string }) => {
  return request.post<ApiResponse<{ token: string; user: User }>>('/auth/register', data)
}

export const getProfile = () => {
  return request.get<ApiResponse<User>>('/auth/profile')
}

export const getCustomers = (params: PageParams & { keyword?: string }) => {
  return request.get<ApiResponse<PageResult<Customer>>>('/customers', { params })
}

export const getCustomer = (id: number) => {
  return request.get<ApiResponse<Customer>>(`/customers/${id}`)
}

export const getMyCustomer = () => {
  return request.get<ApiResponse<Customer>>('/customers/my')
}

export const createCustomer = (data: any) => {
  return request.post<ApiResponse<Customer>>('/customers', data)
}

export const updateCustomer = (id: number, data: any) => {
  return request.put<ApiResponse<Customer>>(`/customers/${id}`, data)
}

export const updateMyCustomer = (data: any) => {
  return request.put<ApiResponse<Customer>>('/customers/my', data)
}
