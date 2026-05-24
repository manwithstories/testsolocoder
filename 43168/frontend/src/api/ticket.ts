import request from '@/utils/request'
import type { PaginationParams, PaginationResult } from '@/types'

export interface Ticket {
  id: number
  title: string
  description: string
  type: number
  status: number
  priority: number
  ownerId?: number
  ownerName?: string
  handlerId?: number
  handlerName?: string
  relatedOrderId?: number
  relatedOrderNo?: string
  createdAt: string
  updatedAt: string
  closedAt?: string
  reply?: string
}

export interface ListTicketParams extends PaginationParams {
  title?: string
  status?: number
  type?: number
  priority?: number
  keyword?: string
}

export interface TicketFormData {
  id?: number
  title: string
  description: string
  type: number
  priority?: number
  relatedOrderId?: number
}

export const listTickets = (params: ListTicketParams) => {
  return request.get<any, PaginationResult<Ticket>>('/tickets', { params })
}

export const getTicket = (id: number | string) => {
  return request.get<any, Ticket>(`/tickets/${id}`)
}

export const createTicket = (data: TicketFormData) => {
  return request.post<any, Ticket>('/tickets', data)
}

export const updateTicket = (id: number | string, data: Partial<TicketFormData> & { status?: number; reply?: string }) => {
  return request.put<any, Ticket>(`/tickets/${id}`, data)
}

export const deleteTicket = (id: number | string) => {
  return request.delete<any, void>(`/tickets/${id}`)
}

export const assignTicket = (id: number | string, handlerId: number) => {
  return request.patch<any, void>(`/tickets/${id}/assign`, { handlerId })
}

export const closeTicket = (id: number | string, reply: string) => {
  return request.patch<any, void>(`/tickets/${id}/close`, { reply })
}
