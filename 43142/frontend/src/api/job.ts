import { get, post, put, del } from '@/utils/request'
import type { ApiResponse, Job, CreateJobRequest, PaginatedResponse } from '@/types'

export function listJobs(params?: {
  page?: number
  page_size?: number
  company_id?: number
  status?: string
  keyword?: string
}) {
  return get<ApiResponse<PaginatedResponse<Job>>>('/jobs', { params })
}

export function getJob(id: number) {
  return get<ApiResponse<Job>>(`/jobs/${id}`)
}

export function createJob(data: CreateJobRequest) {
  return post<ApiResponse<Job>>('/company/jobs', data)
}

export function updateJob(id: number, data: Partial<CreateJobRequest> & { status?: string }) {
  return put<ApiResponse<Job>>(`/company/jobs/${id}`, data)
}

export function deleteJob(id: number) {
  return del<ApiResponse<any>>(`/company/jobs/${id}`)
}

export function listMyJobs(params?: { page?: number; page_size?: number; status?: string }) {
  return get<ApiResponse<PaginatedResponse<Job>>>('/company/jobs', { params })
}

export function publishJob(id: number) {
  return put<ApiResponse<any>>(`/company/jobs/${id}/publish`)
}

export function closeJob(id: number) {
  return put<ApiResponse<any>>(`/company/jobs/${id}/close`)
}

export function bulkImportJobs(jobs: any[]) {
  return post<ApiResponse<{ success_count: number; errors: string[] }>>('/company/jobs/bulk/import', { jobs })
}

export function bulkDeleteJobs(ids: number[]) {
  return del<ApiResponse<any>>('/company/jobs/bulk', { data: { ids } })
}

export function exportJobs() {
  return get<ApiResponse<any[]>>('/company/jobs/export')
}

export function getJobViewStats(id: number, days?: number) {
  return get<ApiResponse<{ days: number; views: number }>>(`/jobs/${id}/stats/views`, { params: { days } })
}
