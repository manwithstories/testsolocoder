import request from './index'
import type { ApiResponse } from './index'

export interface ServiceItem {
  id: number
  company_id: number
  name: string
  category: string
  desc: string
  min_price: number
  max_price: number
  duration: number
  skills: string
  status: string
  created_at: string
  updated_at: string
}

export interface CreateServicePayload {
  name: string
  category: string
  desc: string
  min_price: number
  max_price: number
  duration: number
  skills?: string
}

export function listServices(params?: Record<string, string | number>) {
  return request.get<ApiResponse<ServiceItem[]>>('/services', { params })
}

export function getService(id: number) {
  return request.get<ApiResponse<ServiceItem>>(`/services/${id}`)
}

export function createService(payload: CreateServicePayload) {
  return request.post<ApiResponse<ServiceItem>>('/company/services', payload)
}

export function updateService(id: number, payload: CreateServicePayload) {
  return request.put<ApiResponse<ServiceItem>>(`/company/services/${id}`, payload)
}

export function deleteService(id: number) {
  return request.delete<ApiResponse<null>>(`/company/services/${id}`)
}

export function myServices() {
  return request.get<ApiResponse<ServiceItem[]>>('/company/my-services')
}
