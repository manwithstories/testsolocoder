import request from './request'
import { Order, OrderInvitation, PaginatedResponse, CreateOrderParams } from '@/types'

export const orderApi = {
  getList: (params?: { page?: number; page_size?: number; status?: string }) => {
    return request.get<any, PaginatedResponse<Order>>('/orders', { params })
  },

  getDetail: (id: number) => {
    return request.get<any, Order>(`/orders/${id}`)
  },

  create: (params: CreateOrderParams) => {
    return request.post<any, Order>('/orders', params)
  },

  cancel: (id: number, params: { reason: string }) => {
    return request.put<any, any>(`/orders/${id}/cancel`, params)
  },

  start: (id: number, params?: { longitude?: number; latitude?: number; location?: string }) => {
    return request.put<any, any>(`/orders/${id}/start`, params)
  },

  complete: (id: number, params?: { location?: string }) => {
    return request.put<any, any>(`/orders/${id}/complete`, params)
  },
}

export const invitationApi = {
  getMyInvitations: (params?: { page?: number; page_size?: number; status?: string }) => {
    return request.get<any, PaginatedResponse<OrderInvitation>>('/invitations', { params })
  },

  respond: (id: number, params: { accepted: boolean; reject_reason?: string }) => {
    return request.put<any, any>(`/invitations/${id}/respond`, params)
  },
}
