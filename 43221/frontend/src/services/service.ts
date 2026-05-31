import { http } from './request'
import {
  Service,
  Schedule,
  PaginatedResponse,
  CreateServiceRequest,
  CreateScheduleRequest,
  BatchCreateScheduleRequest,
} from '@/types'

export const serviceApi = {
  getAll: (params?: { page?: number; page_size?: number; service_type?: string }) =>
    http.get<PaginatedResponse<Service>>('/services', { params }),
  getById: (id: string) =>
    http.get<Service>(`/services/${id}`),
  getProfessionalServices: (params?: { page?: number; page_size?: number }) =>
    http.get<PaginatedResponse<Service>>('/services/professional/list', { params }),
  create: (data: CreateServiceRequest) =>
    http.post<Service>('/services', data),
  update: (id: string, data: any) =>
    http.put<Service>(`/services/${id}`, data),
  delete: (id: string) =>
    http.delete(`/services/${id}`),
  getSchedules: (serviceId: string, params?: { date?: string; only_available?: boolean }) =>
    http.get<Schedule[]>(`/services/${serviceId}/schedules`, { params }),
  createSchedule: (data: CreateScheduleRequest) =>
    http.post<Schedule>('/services/schedules', data),
  batchCreateSchedules: (data: BatchCreateScheduleRequest) =>
    http.post('/services/schedules/batch', data),
  deleteSchedules: (serviceId: string, params: { start_date: string; end_date: string }) =>
    http.delete(`/services/${serviceId}/schedules`, { params }),
}

export { Service, Schedule }
