import { get, post, put } from '@/utils/request'
import type { ApiResponse, Interview, ScheduleInterviewRequest, PaginatedResponse } from '@/types'

export function scheduleInterview(data: ScheduleInterviewRequest) {
  return post<ApiResponse<Interview>>('/company/interviews', data)
}

export function getInterview(id: number) {
  return get<ApiResponse<Interview>>(`/interviews/${id}`)
}

export function listCompanyInterviews(params?: { page?: number; page_size?: number; status?: string }) {
  return get<ApiResponse<PaginatedResponse<Interview>>>('/company/interviews', { params })
}

export function listMyInterviews(params?: { page?: number; page_size?: number; status?: string }) {
  return get<ApiResponse<PaginatedResponse<Interview>>>('/interviews/my', { params })
}

export function updateInterview(id: number, data: any) {
  return put<ApiResponse<Interview>>(`/company/interviews/${id}`, data)
}

export function acceptInterview(id: number) {
  return put<ApiResponse<any>>(`/interviews/${id}/accept`)
}

export function rejectInterview(id: number) {
  return put<ApiResponse<any>>(`/interviews/${id}/reject`)
}

export function completeInterview(id: number, feedback: string, rating: number) {
  return put<ApiResponse<any>>(`/company/interviews/${id}/complete`, { feedback, rating })
}

export function cancelInterview(id: number) {
  return put<ApiResponse<any>>(`/company/interviews/${id}/cancel`)
}
