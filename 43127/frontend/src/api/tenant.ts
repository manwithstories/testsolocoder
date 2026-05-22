import request from '@/utils/request'
import type { Tenant, Appointment, Contract, PaginationParams, PaginationResult, ApiResponse } from '@/types'

export const getTenants = (params: PaginationParams) => {
  return request.get<any, ApiResponse<PaginationResult<Tenant>>>('/tenants', { params })
}

export const getTenant = (id: number) => {
  return request.get<any, ApiResponse<Tenant>>(`/tenants/${id}`)
}

export const createTenant = (data: any) => {
  return request.post<any, ApiResponse<Tenant>>('/tenants', data)
}

export const updateTenant = (id: number, data: any) => {
  return request.put<any, ApiResponse<null>>(`/tenants/${id}`, data)
}

export const deleteTenant = (id: number) => {
  return request.delete<any, ApiResponse<null>>(`/tenants/${id}`)
}

export const getAppointments = (params: PaginationParams & { status?: number }) => {
  return request.get<any, ApiResponse<PaginationResult<Appointment>>>('/appointments', { params })
}

export const createAppointment = (data: any) => {
  return request.post<any, ApiResponse<Appointment>>('/appointments', data)
}

export const updateAppointmentStatus = (id: number, status: number) => {
  return request.put<any, ApiResponse<null>>(`/appointments/${id}/status`, { status })
}

export const getContracts = (params: PaginationParams & { status?: number }) => {
  return request.get<any, ApiResponse<PaginationResult<Contract>>>('/contracts', { params })
}

export const getContract = (id: number) => {
  return request.get<any, ApiResponse<Contract>>(`/contracts/${id}`)
}

export const createContract = (data: any) => {
  return request.post<any, ApiResponse<Contract>>('/contracts', data)
}

export const updateContract = (id: number, data: any) => {
  return request.put<any, ApiResponse<null>>(`/contracts/${id}`, data)
}

export const updateContractStatus = (id: number, status: number) => {
  return request.put<any, ApiResponse<null>>(`/contracts/${id}/status`, { status })
}

export const getExpiringContracts = () => {
  return request.get<any, ApiResponse<{ expiring: string[]; count: number }>>('/contracts/expiring')
}
