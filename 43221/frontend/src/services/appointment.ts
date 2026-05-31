import { http } from './request'
import {
  Appointment,
  Payment,
  PaginatedResponse,
  CreateAppointmentRequest,
} from '@/types'

export const appointmentApi = {
  create: (data: CreateAppointmentRequest) =>
    http.post<{ appointment: Appointment; payment: Payment }>('/appointments', data),
  getById: (id: string) =>
    http.get<Appointment>(`/appointments/${id}`),
  getClientAppointments: (params?: { page?: number; page_size?: number; status?: string }) =>
    http.get<PaginatedResponse<Appointment>>('/appointments/client/list', { params }),
  getProfessionalAppointments: (params?: { page?: number; page_size?: number; status?: string }) =>
    http.get<PaginatedResponse<Appointment>>('/appointments/professional/list', { params }),
  confirm: (id: string, data: { appointment_id: string }) =>
    http.put<Appointment>(`/appointments/${id}/confirm`, data),
  cancel: (id: string, data: { appointment_id: string; reason: string }) =>
    http.put<Appointment>(`/appointments/${id}/cancel`, data),
  complete: (id: string) =>
    http.put<Appointment>(`/appointments/${id}/complete`),
  pay: (id: string, data: { transaction_id: string }) =>
    http.put(`/appointments/${id}/pay`, data),
  refund: (id: string, data: { reason: string }) =>
    http.put(`/appointments/${id}/refund`, data),
}

export { Appointment, Payment }
