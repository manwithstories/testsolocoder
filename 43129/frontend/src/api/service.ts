import request from '@/utils/request'
import type { ApiResponse, Service, PageResult, PageParams } from '@/types'

export const getServices = (params: PageParams & { category?: string; is_package?: boolean }) => {
  return request.get<ApiResponse<PageResult<Service>>>('/services', { params })
}

export const getAllServices = () => {
  return request.get<ApiResponse<Service[]>>('/services/all')
}

export const getService = (id: number) => {
  return request.get<ApiResponse<Service>>(`/services/${id}`)
}

export const createService = (data: any) => {
  return request.post<ApiResponse<Service>>('/services', data)
}

export const updateService = (id: number, data: any) => {
  return request.put<ApiResponse<Service>>(`/services/${id}`, data)
}

export const deleteService = (id: number) => {
  return request.delete<ApiResponse<null>>(`/services/${id}`)
}

export const addPackageService = (data: { service_id: number; child_service_id: number; count?: number }) => {
  return request.post<ApiResponse<null>>('/services/package', data)
}

export const getPackageServices = (id: number) => {
  return request.get<ApiResponse<any[]>>(`/services/${id}/package-services`)
}

export const deletePackageServices = (id: number) => {
  return request.delete<ApiResponse<null>>(`/services/${id}/package-services`)
}
