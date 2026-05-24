import request from '@/utils/request'
import {
  Application,
  CreateApplicationRequest,
  UpdateApplicationRequest,
  PaginatedData
} from '@/types'

export const applicationApi = {
  create: (data: CreateApplicationRequest) => {
    return request.post<any, Application>('/applications', data)
  },

  getList: (params?: {
    page?: number
    pageSize?: number
    status?: string
    keyword?: string
  }) => {
    return request.get<any, PaginatedData<Application>>('/applications', { params })
  },

  getById: (id: number) => {
    return request.get<any, Application>(`/applications/${id}`)
  },

  update: (id: number, data: UpdateApplicationRequest) => {
    return request.put<any, null>(`/applications/${id}`, data)
  },

  submit: (id: number) => {
    return request.post<any, null>(`/applications/${id}/submit`)
  },

  cancel: (id: number) => {
    return request.post<any, null>(`/applications/${id}/cancel`)
  },

  uploadMaterials: (id: number, field: string, file: File) => {
    const formData = new FormData()
    formData.append('field', field)
    formData.append('file', file)
    return request.post<any, { filePath: string; fileUrl: string }>(
      `/applications/${id}/upload`,
      formData,
      { headers: { 'Content-Type': 'multipart/form-data' } }
    )
  },

  review: (id: number, approved: boolean, comments: string) => {
    return request.post<any, null>(`/applications/${id}/review`, { approved, comments })
  },

  assignAgent: (id: number, agentId: number) => {
    return request.post<any, null>(`/applications/${id}/assign`, { agentId })
  }
}
