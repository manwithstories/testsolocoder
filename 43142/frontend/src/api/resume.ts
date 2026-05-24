import { get, post, put, del, upload } from '@/utils/request'
import type { ApiResponse, Resume, CreateResumeRequest, PaginatedResponse } from '@/types'

export function createResume(data: CreateResumeRequest) {
  return post<ApiResponse<Resume>>('/resumes', data)
}

export function getResume(id: number) {
  return get<ApiResponse<Resume>>(`/resumes/${id}`)
}

export function updateResume(id: number, data: CreateResumeRequest) {
  return put<ApiResponse<Resume>>(`/resumes/${id}`, data)
}

export function deleteResume(id: number) {
  return del<ApiResponse<any>>(`/resumes/${id}`)
}

export function listResumes() {
  return get<ApiResponse<Resume[]>>('/resumes')
}

export function setDefaultResume(id: number) {
  return put<ApiResponse<any>>(`/resumes/${id}/default`)
}

export function uploadResumeFile(id: number, file: File, onProgress?: (progress: number) => void) {
  return upload<ApiResponse<Resume>>(`/resumes/${id}/upload`, file, onProgress)
}

export function getDefaultResume() {
  return get<ApiResponse<Resume>>('/resumes/default')
}

export function searchResumes(params?: {
  page?: number
  page_size?: number
  keyword?: string
  skills?: string[]
}) {
  return get<ApiResponse<PaginatedResponse<Resume>>>('/search/resumes', { params })
}
