import request from './index'
import type { ApiResponse } from './index'

export interface OrderItem {
  id: number
  booking_id: number
  customer_id: number
  staff_id: number
  company_id: number
  total_amount: number
  status: string
  report_text?: string
  report_images?: string
  reported_at?: string
  confirmed_at?: string
  paid_at?: string
  created_at: string
  updated_at: string
  booking?: any
}

export function listOrders(params?: Record<string, string>) {
  return request.get<ApiResponse<OrderItem[]>>('/orders', { params })
}

export function getOrder(id: number) {
  return request.get<ApiResponse<OrderItem>>(`/orders/${id}`)
}

export function submitReport(id: number, payload: { report_text: string; report_images?: string }) {
  return request.post<ApiResponse<OrderItem>>(`/orders/${id}/report`, payload)
}

export function confirmOrder(id: number) {
  return request.post<ApiResponse<OrderItem>>(`/orders/${id}/confirm`)
}

export function requestRefund(id: number, payload?: { reason?: string }) {
  return request.post<ApiResponse<any>>(`/orders/${id}/refund`, payload || {})
}
