import request from './request'
import type { MaintenanceRecord, MaintenanceSchedule, CreateMaintenanceRequest, CreateMaintenanceScheduleRequest, PaginatedResponse } from '@/types/maintenance'

export const getMaintenancesApi = (params?: { ship_id?: string; status?: string; maintenance_type?: string; due_soon?: boolean; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<MaintenanceRecord>>('/maintenance', { params })
}

export const getMaintenanceApi = (id: string) => {
  return request.get<MaintenanceRecord>(`/maintenance/${id}`)
}

export const createMaintenanceApi = (data: CreateMaintenanceRequest) => {
  return request.post<MaintenanceRecord>('/maintenance', data)
}

export const updateMaintenanceApi = (id: string, data: Partial<CreateMaintenanceRequest> & { status?: string }) => {
  return request.put<MaintenanceRecord>(`/maintenance/${id}`, data)
}

export const deleteMaintenanceApi = (id: string) => {
  return request.delete(`/maintenance/${id}`)
}

export const getSchedulesApi = (params?: { ship_id?: string; due_soon?: boolean }) => {
  return request.get<MaintenanceSchedule[]>('/maintenance-schedules', { params })
}

export const createScheduleApi = (data: CreateMaintenanceScheduleRequest) => {
  return request.post<MaintenanceSchedule>('/maintenance-schedules', data)
}

export const updateScheduleApi = (id: string, data: Partial<CreateMaintenanceScheduleRequest>) => {
  return request.put<MaintenanceSchedule>(`/maintenance-schedules/${id}`, data)
}

export const deleteScheduleApi = (id: string) => {
  return request.delete(`/maintenance-schedules/${id}`)
}
