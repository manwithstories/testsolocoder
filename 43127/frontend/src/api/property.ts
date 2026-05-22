import request from '@/utils/request'
import type { Property, Facility, PaginationParams, PaginationResult, ApiResponse } from '@/types'

export const getProperties = (params: PaginationParams & {
  region?: string
  layout?: string
  minRent?: number
  maxRent?: number
  sortBy?: string
  sortOrder?: string
}) => {
  return request.get<any, ApiResponse<PaginationResult<Property>>>('/properties', { params })
}

export const getProperty = (id: number) => {
  return request.get<any, ApiResponse<Property>>(`/properties/${id}`)
}

export const createProperty = (data: any) => {
  return request.post<any, ApiResponse<Property>>('/properties', data)
}

export const updateProperty = (id: number, data: any) => {
  return request.put<any, ApiResponse<null>>(`/properties/${id}`, data)
}

export const deleteProperty = (id: number) => {
  return request.delete<any, ApiResponse<null>>(`/properties/${id}`)
}

export const updatePropertyStatus = (id: number, status: number) => {
  return request.put<any, ApiResponse<null>>(`/properties/${id}/status`, { status })
}

export const uploadImage = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<any, ApiResponse<{ url: string }>>('/properties/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export const getFacilities = () => {
  return request.get<any, ApiResponse<Facility[]>>('/facilities')
}

export const getMyProperties = (params: PaginationParams) => {
  return request.get<any, ApiResponse<PaginationResult<Property>>>('/properties/my/list', { params })
}
