import request from './request'
import type { Ship, CreateShipRequest, SearchShipRequest, PaginatedResponse, ApiResponse } from '@/types/ship'

export const getShipsApi = (params?: SearchShipRequest) => {
  return request.get<PaginatedResponse<Ship>>('/ships', { params })
}

export const getShipApi = (id: string) => {
  return request.get<Ship>(`/ships/${id}`)
}

export const getMyShipsApi = () => {
  return request.get<Ship[]>('/my-ships')
}

export const createShipApi = (data: CreateShipRequest) => {
  return request.post<Ship>('/ships', data)
}

export const updateShipApi = (id: string, data: Partial<CreateShipRequest>) => {
  return request.put<Ship>(`/ships/${id}`, data)
}

export const deleteShipApi = (id: string) => {
  return request.delete(`/ships/${id}`)
}

export const uploadShipImageApi = (id: string, formData: FormData) => {
  return request.post(`/ships/${id}/images`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export const deleteShipImageApi = (shipId: string, imageId: string) => {
  return request.delete(`/ships/${shipId}/images/${imageId}`)
}
