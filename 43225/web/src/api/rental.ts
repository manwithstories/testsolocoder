import request from './request'
import type { Rental, CreateRentalRequest, UpdateRentalStatusRequest, PaginatedResponse } from '@/types/rental'

export const getRentalsApi = (params?: { status?: string; ship_id?: string; start_date?: string; end_date?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<Rental>>('/rentals', { params })
}

export const getRentalApi = (id: string) => {
  return request.get<Rental>(`/rentals/${id}`)
}

export const getMyRentalsApi = () => {
  return request.get<Rental[]>('/my-rentals')
}

export const createRentalApi = (data: CreateRentalRequest) => {
  return request.post<Rental>('/rentals', data)
}

export const updateRentalStatusApi = (id: string, data: UpdateRentalStatusRequest) => {
  return request.put<Rental>(`/rentals/${id}/status`, data)
}

export const cancelRentalApi = (id: string, reason?: string) => {
  return request.delete(`/rentals/${id}`, { data: { reason } })
}
