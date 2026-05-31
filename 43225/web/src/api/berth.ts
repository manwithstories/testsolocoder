import request from './request'
import type { Berth, Dock, BerthReservation, CreateBerthRequest, CreateReservationRequest, PaginatedResponse } from '@/types/berth'

export const getDocksApi = (params?: { city?: string; country?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<Dock>>('/docks', { params })
}

export const getDockApi = (id: string) => {
  return request.get<Dock>(`/docks/${id}`)
}

export const createDockApi = (data: Partial<Dock>) => {
  return request.post<Dock>('/admin/docks', data)
}

export const updateDockApi = (id: string, data: Partial<Dock>) => {
  return request.put<Dock>(`/admin/docks/${id}`, data)
}

export const deleteDockApi = (id: string) => {
  return request.delete(`/admin/docks/${id}`)
}

export const getBerthsApi = (params?: { dock_id?: string; berth_type?: string; status?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<Berth>>('/berths', { params })
}

export const getBerthApi = (id: string) => {
  return request.get<Berth>(`/berths/${id}`)
}

export const createBerthApi = (data: CreateBerthRequest) => {
  return request.post<Berth>('/admin/berths', data)
}

export const updateBerthApi = (id: string, data: Partial<CreateBerthRequest>) => {
  return request.put<Berth>(`/admin/berths/${id}`, data)
}

export const deleteBerthApi = (id: string) => {
  return request.delete(`/admin/berths/${id}`)
}

export const checkBerthAvailabilityApi = (params: { berth_id: string; start_time: string; end_time: string }) => {
  return request.get('/berths/availability', { params })
}

export const getReservationsApi = (params?: { status?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<BerthReservation>>('/berth-reservations', { params })
}

export const createReservationApi = (data: CreateReservationRequest) => {
  return request.post<BerthReservation>('/berth-reservations', data)
}

export const cancelReservationApi = (id: string) => {
  return request.put(`/berth-reservations/${id}/cancel`)
}

export const recordWaterLevelApi = (dockId: string, data: { height: number; unit?: string; recorded_at: string }) => {
  return request.post(`/admin/docks/${dockId}/water-levels`, data)
}

export const getWaterLevelsApi = (dockId: string, params?: { start_date?: string; end_date?: string }) => {
  return request.get(`/docks/${dockId}/water-levels`, { params })
}
