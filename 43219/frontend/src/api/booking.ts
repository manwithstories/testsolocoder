import request from './index'
import type { ApiResponse } from './index'

export interface Booking {
  id: number
  customer_id: number
  service_id: number
  staff_id?: number
  start_at: string
  end_at: string
  address: string
  remark?: string
  price: number
  status: string
  reschedule_count: number
  need_review: boolean
  created_at: string
  updated_at: string
  service?: any
  staff?: any
  customer?: any
}

export interface CreateBookingPayload {
  service_id: number
  start_at: string
  end_at: string
  address: string
  remark?: string
  staff_id?: number
  price?: number
}

export function createBooking(payload: CreateBookingPayload) {
  return request.post<ApiResponse<Booking>>('/bookings', payload)
}

export function listBookings(params?: Record<string, string>) {
  return request.get<ApiResponse<Booking[]>>('/bookings', { params })
}

export function getBooking(id: number) {
  return request.get<ApiResponse<Booking>>(`/bookings/${id}`)
}

export function confirmBooking(id: number) {
  return request.post<ApiResponse<Booking>>(`/bookings/${id}/confirm`)
}

export function rescheduleBooking(id: number, payload: { start_at: string; end_at: string }) {
  return request.post<ApiResponse<Booking>>(`/bookings/${id}/reschedule`, payload)
}

export function cancelBooking(id: number) {
  return request.post<ApiResponse<string>>(`/bookings/${id}/cancel`)
}

export function reviewReschedule(id: number, payload: { approve: boolean; start_at?: string; end_at?: string }) {
  return request.post<ApiResponse<Booking>>(`/company/bookings/${id}/review-reschedule`, payload)
}
