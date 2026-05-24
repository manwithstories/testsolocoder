import request from '@/utils/request'
import {
  User,
  CreateAgentRequest,
  UpdateAgentProfileRequest,
  PaginatedData
} from '@/types'

export const agentApi = {
  getList: (params?: {
    page?: number
    pageSize?: number
    status?: string
    keyword?: string
  }) => {
    return request.get<any, PaginatedData<User>>('/admin/agents', { params })
  },

  getAvailable: () => {
    return request.get<any, User[]>('/admin/agents/available')
  },

  getById: (id: number) => {
    return request.get<any, User>(`/admin/agents/${id}`)
  },

  create: (data: CreateAgentRequest) => {
    return request.post<any, User>('/admin/agents', data)
  },

  updateProfile: (id: number, data: UpdateAgentProfileRequest) => {
    return request.put<any, null>(`/admin/agents/${id}`, data)
  },

  delete: (id: number) => {
    return request.delete<any, null>(`/admin/agents/${id}`)
  },

  updateSchedule: (id: number, startTime: string, endTime: string) => {
    return request.post<any, null>(`/admin/agents/${id}/schedule`, { startTime, endTime })
  },

  updateMaxApps: (id: number, maxApps: number) => {
    return request.put<any, null>(`/admin/agents/${id}/max-apps`, { maxApps })
  },

  getStats: (id: number) => {
    return request.get<any, Record<string, any>>(`/admin/agents/${id}/stats`)
  },

  getPerformance: (id: number, params?: { startDate?: string; endDate?: string }) => {
    return request.get<any, Record<string, any>>(`/admin/agents/${id}/performance`, { params })
  },

  autoAssign: (applicationId: number) => {
    return request.post<any, User>(`/admin/agents/auto-assign/${applicationId}`)
  }
}
