import request from '@/utils/request'
import type {
  GuideSchedule,
  GuideContent,
  ResearchApplication,
  Museum,
  ApiResponse,
  PageResult
} from '@/types'

export const listGuideSchedules = (params?: { start_date?: string; end_date?: string }) => {
  return request.get<any, ApiResponse<GuideSchedule[]>>('/guides/schedules', { params })
}

export const createGuideSchedule = (data: {
  date: string
  start_time: string
  end_time: string
  is_available: boolean
}) => {
  return request.post<any, ApiResponse<GuideSchedule>>('/guides/schedules', data)
}

export const updateGuideSchedule = (id: number, data: {
  date: string
  start_time: string
  end_time: string
  is_available: boolean
}) => {
  return request.put<any, ApiResponse<void>>(`/guides/schedules/${id}`, data)
}

export const deleteGuideSchedule = (id: number) => {
  return request.delete<any, ApiResponse<void>>(`/guides/schedules/${id}`)
}

export const listGuideContents = (params?: {
  collection_id?: number
  exhibition_id?: number
  language?: string
}) => {
  return request.get<any, ApiResponse<GuideContent[]>>('/guides/contents', { params })
}

export const createGuideContent = (data: Partial<GuideContent>) => {
  return request.post<any, ApiResponse<GuideContent>>('/guides/contents', data)
}

export const updateGuideContent = (id: number, data: Partial<GuideContent>) => {
  return request.put<any, ApiResponse<void>>(`/guides/contents/${id}`, data)
}

export const deleteGuideContent = (id: number) => {
  return request.delete<any, ApiResponse<void>>(`/guides/contents/${id}`)
}

export const createResearchApplication = (data: {
  collection_id: number
  purpose: string
  institution: string
}) => {
  return request.post<any, ApiResponse<ResearchApplication>>('/research/applications', data)
}

export const reviewResearchApplication = (id: number, data: {
  status: 'approved' | 'rejected'
  review_comment: string
}) => {
  return request.put<any, ApiResponse<void>>(`/research/applications/${id}/review`, data)
}

export const listMyResearchApplications = (params?: { page?: number; page_size?: number }) => {
  return request.get<any, ApiResponse<PageResult<ResearchApplication>>>('/research/applications/my', { params })
}

export const listResearchApplications = (params?: {
  page?: number
  page_size?: number
  status?: string
}) => {
  return request.get<any, ApiResponse<PageResult<ResearchApplication>>>('/research/applications', { params })
}

export const getResearchApplication = (id: number) => {
  return request.get<any, ApiResponse<ResearchApplication>>(`/research/applications/${id}`)
}

export const getStatistics = (params?: {
  start_date?: string
  end_date?: string
  exhibition_id?: number
  type?: string
}) => {
  return request.get<any, ApiResponse<any>>('/statistics', { params })
}

export const exportStatisticsExcel = (params?: {
  start_date?: string
  end_date?: string
  exhibition_id?: number
}) => {
  return request.get<any, any>('/statistics/export/excel', {
    params,
    responseType: 'blob'
  })
}

export const exportStatisticsPDF = (params?: {
  start_date?: string
  end_date?: string
  exhibition_id?: number
}) => {
  return request.get<any, ApiResponse<any>>('/statistics/export/pdf', { params })
}

export const listMuseums = () => {
  return request.get<any, ApiResponse<Museum[]>>('/museums')
}

export const getMuseum = (id: number) => {
  return request.get<any, ApiResponse<Museum>>(`/museums/${id}`)
}

export const createMuseum = (data: Partial<Museum>) => {
  return request.post<any, ApiResponse<Museum>>('/museums', data)
}

export const updateMuseum = (id: number, data: Partial<Museum>) => {
  return request.put<any, ApiResponse<void>>(`/museums/${id}`, data)
}

export const deleteMuseum = (id: number) => {
  return request.delete<any, ApiResponse<void>>(`/museums/${id}`)
}

export const uploadImage = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<any, ApiResponse<{
    url: string
    filename: string
    size: number
  }>>('/uploads/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
