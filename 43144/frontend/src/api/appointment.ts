import { apiGet, apiPost, apiPut } from './client'
import {
  Appointment,
  CreateAppointmentRequest,
  UpdateAppointmentRequest,
  AppointmentListQuery,
  ApiResponse,
} from '../types'

export const listAppointments = (params?: AppointmentListQuery): Promise<ApiResponse<any>> => {
  return apiGet('/appointments', params)
}

export const getAppointment = (id: number): Promise<ApiResponse<Appointment>> => {
  return apiGet<Appointment>(`/appointments/${id}`)
}

export const createAppointment = (data: CreateAppointmentRequest): Promise<ApiResponse<Appointment>> => {
  return apiPost<Appointment>('/appointments', data)
}

export const updateAppointment = (id: number, data: UpdateAppointmentRequest): Promise<ApiResponse<Appointment>> => {
  return apiPut<Appointment>(`/appointments/${id}`, data)
}

export const cancelAppointment = (id: number, reason?: string): Promise<ApiResponse<Appointment>> => {
  return apiPut<Appointment>(`/appointments/${id}/cancel`, { reason })
}

export const rescheduleAppointment = (id: number, data: UpdateAppointmentRequest): Promise<ApiResponse<Appointment>> => {
  return apiPut<Appointment>(`/appointments/${id}/reschedule`, data)
}

export const confirmAppointment = (id: number): Promise<ApiResponse<Appointment>> => {
  return apiPut<Appointment>(`/appointments/${id}/confirm`)
}

export const completeAppointment = (id: number): Promise<ApiResponse<Appointment>> => {
  return apiPut<Appointment>(`/appointments/${id}/complete`)
}
