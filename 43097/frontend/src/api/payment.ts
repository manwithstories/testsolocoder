import request from '@/utils/request'
import { Payment, PageParams, PageResult, PaymentStatus, PaymentMethod } from '@/types'

export const getPaymentList = (params: PageParams & { status?: PaymentStatus; method?: PaymentMethod; startDate?: string; endDate?: string }) => {
  return request.get<PageResult<Payment>>('/payments', { params })
}

export const getPaymentById = (id: number) => {
  return request.get<Payment>(`/payments/${id}`)
}

export const createPayment = (data: Omit<Payment, 'id' | 'paymentNo' | 'createdAt' | 'updatedAt'>) => {
  return request.post<Payment>('/payments', data)
}

export const updatePayment = (id: number, data: Partial<Payment>) => {
  return request.put<Payment>(`/payments/${id}`, data)
}

export const refundPayment = (id: number, data?: { amount?: number; reason?: string }) => {
  return request.post(`/payments/${id}/refund`, data)
}

export const getPaymentsByBookingId = (bookingId: number) => {
  return request.get<Payment[]>(`/payments/booking/${bookingId}`)
}

export const getPaymentsByCheckInId = (checkInId: number) => {
  return request.get<Payment[]>(`/payments/checkin/${checkInId}`)
}

export const getPaymentStatistics = (startDate: string, endDate: string) => {
  return request.get('/payments/statistics', { params: { startDate, endDate } })
}
