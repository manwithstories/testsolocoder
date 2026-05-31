import { request } from '@/utils/request'
import type { Agency, Package, PackageItem, TimeSlot, PaginationParams, PaginationResponse } from '@/types'

export const registerAgency = (data: any): Promise<{ agency: Agency; admin_user: any }> => {
  return request.post('/agencies/register', data)
}

export const getAgencies = (): Promise<Agency[]> => {
  return request.get('/agencies')
}

export const getAgency = (id: number): Promise<Agency> => {
  return request.get(`/agencies/${id}`)
}

export const updateAgency = (data: any): Promise<void> => {
  return request.put('/agencies', data)
}

export const getPackages = (params: { page: number; page_size: number; keyword?: string }): Promise<PaginationResponse<Package>> => {
  return request.get('/packages', { params })
}

export const getPackage = (id: number): Promise<Package> => {
  return request.get(`/packages/${id}`)
}

export const getHotPackages = (limit: number = 10): Promise<Package[]> => {
  return request.get('/packages/hot', { params: { limit } })
}

export const getAgencyPackages = (params: PaginationParams): Promise<PaginationResponse<Package>> => {
  return request.get('/agency/packages', { params })
}

export const createPackage = (data: any): Promise<Package> => {
  return request.post('/packages', data)
}

export const updatePackage = (id: number, data: any): Promise<void> => {
  return request.put(`/packages/${id}`, data)
}

export const updatePackagePrice = (id: number, price: number): Promise<void> => {
  return request.patch(`/packages/${id}/price`, { price })
}

export const updatePackageStatus = (id: number, status: number): Promise<void> => {
  return request.patch(`/packages/${id}/status`, { status })
}

export const getPackageTimeSlots = (packageId: number): Promise<TimeSlot[]> => {
  return request.get(`/packages/${packageId}/timeslots`)
}

export const createTimeSlot = (data: any): Promise<void> => {
  return request.post('/timeslots', data)
}
