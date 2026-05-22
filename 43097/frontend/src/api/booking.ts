import request from '@/utils/request'
import { Booking, PageParams, PageResult, BookingStatus } from '@/types'

export const getBookingList = (params: PageParams & { status?: BookingStatus; startDate?: string; endDate?: string }) => {
  return request.get<PageResult<Booking>>('/bookings', { params })
}

export const getBookingById = (id: number) => {
  return request.get<Booking>(`/bookings/${id}`)
}

export const createBooking = (data: Omit<Booking, 'id' | 'bookingNo' | 'createdAt' | 'updatedAt'>) => {
  return request.post<Booking>('/bookings', data)
}

export const updateBooking = (id: number, data: Partial<Booking>) => {
  return request.put<Booking>(`/bookings/${id}`, data)
}

export const deleteBooking = (id: number) => {
  return request.delete(`/bookings/${id}`)
}

export const updateBookingStatus = (id: number, status: BookingStatus) => {
  return request.patch(`/bookings/${id}/status`, { status })
}

export const confirmBooking = (id: number) => {
  return request.post(`/bookings/${id}/confirm`)
}

export const cancelBooking = (id: number, reason?: string) => {
  return request.post(`/bookings/${id}/cancel`, { reason })
}

export const checkInFromBooking = (id: number, data: any) => {
  return request.post(`/bookings/${id}/checkin`, data)
}

export const getBookingByNo = (bookingNo: string) => {
  return request.get<Booking>(`/bookings/no/${bookingNo}`)
}
