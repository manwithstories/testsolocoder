import request from './index'
import type { ApiResponse } from './index'

export interface Ticket {
  id: number
  order_id?: number
  customer_id: number
  staff_id?: number
  agent_id?: number
  type: 'complaint' | 'refund'
  title: string
  content: string
  status: string
  escalated: boolean
  created_at: string
  updated_at: string
}

export function createTicket(payload: {
  order_id?: number
  staff_id?: number
  type: string
  title: string
  content: string
}) {
  return request.post<ApiResponse<Ticket>>('/tickets', payload)
}

export function listTickets(params?: Record<string, string>) {
  return request.get<ApiResponse<Ticket[]>>('/tickets', { params })
}

export function assignTicket(id: number, agent_id: number) {
  return request.post<ApiResponse<string>>(`/admin/tickets/${id}/assign`, { agent_id })
}

export function resolveTicket(id: number, payload: { result: string; refund?: boolean }) {
  return request.post<ApiResponse<string>>(`/admin/tickets/${id}/resolve`, payload)
}

export function closeTicket(id: number) {
  return request.post<ApiResponse<string>>(`/admin/tickets/${id}/close`)
}
