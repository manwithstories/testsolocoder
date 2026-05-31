import { request } from '@/utils/request'
import type { Appointment, Report, ReportItem, PaginationParams, PaginationResponse } from '@/types'

export const createAppointment = (data: any): Promise<Appointment> => {
  return request.post('/appointments', data)
}

export const getAppointment = (id: number): Promise<Appointment> => {
  return request.get(`/appointments/${id}`)
}

export const getEmployeeAppointments = (employeeId: number, params: PaginationParams): Promise<PaginationResponse<Appointment>> => {
  return request.get(`/employees/${employeeId}/appointments`, { params })
}

export const getCompanyAppointments = (params: PaginationParams): Promise<PaginationResponse<Appointment>> => {
  return request.get('/company/appointments', { params })
}

export const getAgencyAppointments = (params: PaginationParams): Promise<PaginationResponse<Appointment>> => {
  return request.get('/agency/appointments', { params })
}

export const rescheduleAppointment = (id: number, data: any): Promise<void> => {
  return request.put(`/appointments/${id}/reschedule`, data)
}

export const cancelAppointment = (id: number, reason: string): Promise<void> => {
  return request.put(`/appointments/${id}/cancel`, { reason })
}

export const completeAppointment = (id: number): Promise<void> => {
  return request.patch(`/appointments/${id}/complete`)
}

export const getEmployeeAppointmentStatus = (employeeId: number): Promise<any> => {
  return request.get(`/employees/${employeeId}/appointment-status`)
}

export const checkQuota = (employeeId: number, packageId: number): Promise<{ can_book: boolean }> => {
  return request.get(`/employees/${employeeId}/check-quota`, { params: { package_id: packageId } })
}

export const createReport = (data: any): Promise<Report> => {
  return request.post('/reports', data)
}

export const getReport = (id: number): Promise<Report> => {
  return request.get(`/reports/${id}`)
}

export const getReportByAppointment = (appointmentId: number): Promise<Report> => {
  return request.get(`/appointments/${appointmentId}/report`)
}

export const getEmployeeReports = (employeeId: number, params: PaginationParams): Promise<PaginationResponse<Report>> => {
  return request.get(`/employees/${employeeId}/reports`, { params })
}

export const getCompanyReports = (params: PaginationParams): Promise<PaginationResponse<Report>> => {
  return request.get('/company/reports', { params })
}

export const uploadReportFile = (formData: FormData): Promise<{ file_url: string; file_name: string }> => {
  return request.post('/reports/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

export const downloadReportFile = (reportId: number): Promise<void> => {
  return request.get(`/reports/${reportId}/download`, {
    responseType: 'blob'
  })
}

export const getAbnormalReports = (employeeId: number): Promise<Report[]> => {
  return request.get(`/employees/${employeeId}/abnormal-reports`)
}

export const getAbnormalItems = (employeeId: number): Promise<ReportItem[]> => {
  return request.get(`/employees/${employeeId}/abnormal-items`)
}
