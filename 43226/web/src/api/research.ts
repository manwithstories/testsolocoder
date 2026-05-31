import request from '@/utils/request'
import type {
  ResearchApplication,
  ApiResponse,
  PageResult
} from '@/types'

export const createResearchApplication = (data: {
  collection_id: number
  purpose: string
  institution: string
}) => {
  return request.post<any, ApiResponse<ResearchApplication>>('/research/applications', data)
}

export const listMyResearchApplications = (params?: {
  page?: number
  page_size?: number
  status?: string
}) => {
  return request.get<any, ApiResponse<PageResult<ResearchApplication>>>('/research/applications/my', { params })
}

export const getResearchApplication = (id: number) => {
  return request.get<any, ApiResponse<ResearchApplication>>(`/research/applications/${id}`)
}

export const cancelResearchApplication = (id: number) => {
  return request.put<any, ApiResponse<void>>(`/research/applications/${id}/cancel`)
}
