import { get, post, put } from '@/utils/request'
import type { ApiResponse, Application, ApplyRequest, ApplicationHistory, PaginatedResponse } from '@/types'

export function apply(data: ApplyRequest) {
  return post<ApiResponse<Application>>('/applications', data)
}

export function getApplication(id: number) {
  return get<ApiResponse<Application>>(`/applications/${id}`)
}

export function listMyApplications(params?: { page?: number; page_size?: number; status?: string }) {
  return get<ApiResponse<PaginatedResponse<Application>>>('/applications/my', { params })
}

export function listJobApplications(jobId: number, params?: { page?: number; page_size?: number; status?: string }) {
  return get<ApiResponse<PaginatedResponse<Application>>>(`/company/applications/job/${jobId}`, { params })
}

export function listCompanyApplications(params?: {
  page?: number
  page_size?: number
  status?: string
  keyword?: string
}) {
  return get<ApiResponse<PaginatedResponse<Application>>>('/company/applications', { params })
}

export function updateApplicationStatus(id: number, status: string, changeReason?: string) {
  return put<ApiResponse<Application>>(`/applications/${id}/status`, { status, change_reason: changeReason })
}

export function bulkUpdateStatus(ids: number[], status: string, reason?: string) {
  return put<ApiResponse<{ updated_count: number }>>('/company/applications/bulk/status', { ids, status, reason })
}

export function withdrawApplication(id: number) {
  return put<ApiResponse<any>>(`/applications/${id}/withdraw`)
}

export function getApplicationHistory(id: number) {
  return get<ApiResponse<ApplicationHistory[]>>(`/applications/${id}/history`)
}

export function getStatusCount(jobId: number) {
  return get<ApiResponse<Record<string, number>>>(`/company/applications/job/${jobId}/status-count`)
}
