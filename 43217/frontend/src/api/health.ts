import { request } from '@/utils/request'
import type { HealthRecord, AbnormalItem, RecheckReminder, PaginationParams, PaginationResponse } from '@/types'

export const getHealthRecords = (employeeId: number, params: PaginationParams): Promise<PaginationResponse<HealthRecord>> => {
  return request.get(`/employees/${employeeId}/health-records`, { params })
}

export const getHealthRecordByYear = (employeeId: number, year: number): Promise<HealthRecord> => {
  return request.get(`/employees/${employeeId}/health-records/${year}`)
}

export const getAllHealthRecords = (employeeId: number): Promise<HealthRecord[]> => {
  return request.get(`/employees/${employeeId}/health-records/all`)
}

export const getTrendData = (employeeId: number): Promise<any> => {
  return request.get(`/employees/${employeeId}/health-trend`)
}

export const getHealthSummary = (employeeId: number): Promise<any> => {
  return request.get(`/employees/${employeeId}/health-summary`)
}

export const createAbnormalItem = (data: any): Promise<AbnormalItem> => {
  return request.post('/abnormal-items', data)
}

export const getAbnormalItems = (employeeId: number, params: PaginationParams): Promise<PaginationResponse<AbnormalItem>> => {
  return request.get(`/employees/${employeeId}/abnormal-items`, { params })
}

export const getAllAbnormalItems = (employeeId: number): Promise<AbnormalItem[]> => {
  return request.get(`/employees/${employeeId}/abnormal-items/all`)
}

export const setRecheckDate = (data: any): Promise<void> => {
  return request.post('/abnormal-items/set-recheck', data)
}

export const updateRecheckStatus = (data: any): Promise<void> => {
  return request.put('/abnormal-items/recheck-status', data)
}

export const getNeedRecheckItems = (employeeId: number): Promise<AbnormalItem[]> => {
  return request.get(`/employees/${employeeId}/need-recheck`)
}

export const getReminders = (employeeId: number, params: PaginationParams): Promise<PaginationResponse<RecheckReminder>> => {
  return request.get(`/employees/${employeeId}/reminders`, { params })
}

export const getUnreadReminders = (employeeId: number): Promise<RecheckReminder[]> => {
  return request.get(`/employees/${employeeId}/reminders/unread`)
}

export const markReminderAsRead = (id: number): Promise<void> => {
  return request.put(`/reminders/${id}/read`)
}
