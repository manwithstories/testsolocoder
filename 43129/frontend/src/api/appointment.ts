import request from '@/utils/request'
import type { ApiResponse, Appointment, PageResult, PageParams, AvailableSlot } from '@/types'

export const getAppointments = (params: PageParams & {
  customer_id?: number
  technician_id?: number
  service_id?: number
  status?: string
  start_date?: string
  end_date?: string
}) => {
  return request.get<ApiResponse<PageResult<Appointment>>>('/appointments', { params })
}

export const getMyAppointments = (params: PageParams) => {
  return request.get<ApiResponse<PageResult<Appointment>>>('/appointments/my', { params })
}

export const getAppointment = (id: number) => {
  return request.get<ApiResponse<Appointment>>(`/appointments/${id}`)
}

export const createAppointment = (data: {
  customer_id: number
  technician_id: number
  service_id: number
  package_id?: number
  appointment_date: string
  start_time: string
  remark?: string
}) => {
  return request.post<ApiResponse<Appointment>>('/appointments', data)
}

export const cancelAppointment = (data: { id: number; cancel_reason?: string }) => {
  return request.post<ApiResponse<null>>('/appointments/cancel', data)
}

export const rescheduleAppointment = (data: { id: number; appointment_date: string; start_time: string }) => {
  return request.post<ApiResponse<Appointment>>('/appointments/reschedule', data)
}

export const completeAppointment = (id: number) => {
  return request.post<ApiResponse<null>>(`/appointments/${id}/complete`)
}

export const getAvailableSlots = (id: number, params: { date: string; duration?: number }) => {
  return request.get<ApiResponse<AvailableSlot>>(`/available/${id}/slots`, { params })
}
