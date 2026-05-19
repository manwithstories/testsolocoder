import request from './request'
import type { Member, Membership, LoginRequest, LoginResponse, ApiResponse, PaginatedResponse } from '@/types'

export const memberApi = {
  register: (data: { name: string; phone: string; email?: string; password: string; gender?: string }) => {
    return request.post<any, ApiResponse<Member>>('/members/register', data)
  },

  login: (data: LoginRequest) => {
    return request.post<any, ApiResponse<LoginResponse>>('/members/login', data)
  },

  getList: (params: { page?: number; page_size?: number; keyword?: string }) => {
    return request.get<any, PaginatedResponse<Member[]>>('/members', { params })
  },

  getById: (id: number) => {
    return request.get<any, ApiResponse<Member>>(`/members/${id}`)
  },

  update: (id: number, data: Partial<Member>) => {
    return request.put<any, ApiResponse>(`/members/${id}`, data)
  },

  delete: (id: number) => {
    return request.delete<any, ApiResponse>(`/members/${id}`)
  },

  updateStatus: (id: number, status: number) => {
    return request.patch<any, ApiResponse>(`/members/${id}/status`, { status })
  }
}

export const membershipApi = {
  create: (data: { member_id: number; type: string; start_date?: string; price?: number; auto_renew?: boolean }) => {
    return request.post<any, ApiResponse<Membership>>('/memberships', data)
  },

  getList: (params: { page?: number; page_size?: number; member_id?: number }) => {
    return request.get<any, PaginatedResponse<Membership[]>>('/memberships', { params })
  },

  getById: (id: number) => {
    return request.get<any, ApiResponse<Membership>>(`/memberships/${id}`)
  },

  getByMemberId: (memberId: number) => {
    return request.get<any, ApiResponse<Membership>>(`/memberships/member/${memberId}`)
  },

  renew: (id: number, newType: string) => {
    return request.post<any, ApiResponse<Membership>>(`/memberships/${id}/renew`, { new_type: newType })
  },

  upgrade: (memberId: number, data: { new_type: string; price: number }) => {
    return request.post<any, ApiResponse<Membership>>(`/memberships/${memberId}/upgrade`, data)
  },

  checkValidity: (memberId: number) => {
    return request.get<any, ApiResponse<{ valid: boolean }>>(`/memberships/${memberId}/validity`)
  }
}
