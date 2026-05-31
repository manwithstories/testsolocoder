import request from '@/utils/request'
import type {
  Reservation,
  VisitRecord,
  ApiResponse,
  PageResult
} from '@/types'

export const createReservation = (data: {
  exhibition_id: number
  time_slot_id: number
  visitor_count: number
  guide_type: string
}) => {
  return request.post<any, ApiResponse<Reservation>>('/reservations', data)
}

export const confirmReservation = (id: number) => {
  return request.put<any, ApiResponse<void>>(`/reservations/${id}/confirm`)
}

export const cancelReservation = (id: number, reason?: string) => {
  return request.put<any, ApiResponse<void>>(`/reservations/${id}/cancel`, { reason })
}

export const rescheduleReservation = (id: number, new_time_slot_id: number) => {
  return request.put<any, ApiResponse<void>>(`/reservations/${id}/reschedule`, { new_time_slot_id })
}

export const getReservation = (id: number) => {
  return request.get<any, ApiResponse<Reservation>>(`/reservations/${id}`)
}

export const getReservationByQRCode = (qr_code: string) => {
  return request.get<any, ApiResponse<Reservation>>(`/reservations/qr/${qr_code}`)
}

export const listMyReservations = (params?: { page?: number; page_size?: number }) => {
  return request.get<any, ApiResponse<PageResult<Reservation>>>('/reservations/my', { params })
}

export const listExhibitionReservations = (exhibitionId: number, params?: {
  page?: number
  page_size?: number
  status?: string
}) => {
  return request.get<any, ApiResponse<PageResult<Reservation>>>(`/reservations/exhibition/${exhibitionId}`, { params })
}

export const checkIn = (qr_code: string) => {
  return request.post<any, ApiResponse<void>>('/reservations/check-in', { qr_code })
}

export const checkOut = (id: number) => {
  return request.put<any, ApiResponse<void>>(`/reservations/${id}/check-out`)
}

export const rateVisit = (id: number, data: { rating: number; comment?: string; favorite?: boolean }) => {
  return request.put<any, ApiResponse<void>>(`/reservations/${id}/rate`, data)
}

export const listVisitRecords = (params?: { page?: number; page_size?: number }) => {
  return request.get<any, ApiResponse<PageResult<VisitRecord>>>('/visits/records', { params })
}

export const getVisitStats = () => {
  return request.get<any, ApiResponse<any>>('/visits/stats')
}

export const getReservationStatus = (exhibitionId: number) => {
  return request.get<any, ApiResponse<any>>(`/reservations/status/${exhibitionId}`)
}
