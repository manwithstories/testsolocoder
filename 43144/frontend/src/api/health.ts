import { apiGet, apiPost, apiPut, apiDelete, apiUpload } from './client'
import {
  HealthRecord,
  CreateHealthRecordRequest,
  HealthRecordListQuery,
  HealthReminder,
  ApiResponse,
} from '../types'

export const listHealthRecords = (params?: HealthRecordListQuery): Promise<ApiResponse<any>> => {
  return apiGet('/health/records', params)
}

export const getHealthRecord = (id: number): Promise<ApiResponse<HealthRecord>> => {
  return apiGet<HealthRecord>(`/health/records/${id}`)
}

export const createHealthRecord = (data: CreateHealthRecordRequest): Promise<ApiResponse<HealthRecord>> => {
  return apiPost<HealthRecord>('/health/records', data)
}

export const updateHealthRecord = (id: number, data: any): Promise<ApiResponse<HealthRecord>> => {
  return apiPut<HealthRecord>(`/health/records/${id}`, data)
}

export const deleteHealthRecord = (id: number): Promise<ApiResponse<any>> => {
  return apiDelete(`/health/records/${id}`)
}

export const uploadHealthReport = (id: number, formData: FormData): Promise<ApiResponse<any>> => {
  return apiUpload(`/health/records/${id}/report`, formData)
}

export const getHealthReminders = (petId: number): Promise<ApiResponse<HealthReminder[]>> => {
  return apiGet<HealthReminder[]>(`/health/pets/${petId}/reminders`)
}

export const completeHealthReminder = (id: number): Promise<ApiResponse<any>> => {
  return apiPut(`/health/reminders/${id}/complete`)
}

export const getPetHealthSummary = (petId: number): Promise<ApiResponse<any>> => {
  return apiGet(`/health/pets/${petId}/summary`)
}
