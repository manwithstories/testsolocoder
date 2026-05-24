import request from '@/utils/request'
import {
  OverviewStats,
  StatusDistribution,
  CompanyTypeDistribution,
  AgentPerformance,
  TimeSeriesData,
  ExportTask,
  PaginatedData
} from '@/types'

export const statisticsApi = {
  getOverview: (params?: { startDate?: string; endDate?: string }) => {
    return request.get<any, OverviewStats>('/admin/statistics/overview', { params })
  },

  getStatusDistribution: () => {
    return request.get<any, StatusDistribution[]>('/admin/statistics/status-distribution')
  },

  getCompanyTypeDistribution: () => {
    return request.get<any, CompanyTypeDistribution[]>('/admin/statistics/company-type-distribution')
  },

  getAgentPerformance: (params?: { startDate?: string; endDate?: string }) => {
    return request.get<any, AgentPerformance[]>('/admin/statistics/agent-performance', { params })
  },

  getApplicationTimeSeries: (params: {
    startDate: string
    endDate: string
    interval?: string
  }) => {
    return request.get<any, TimeSeriesData[]>('/admin/statistics/application-time-series', { params })
  },

  getRevenueStats: (params: {
    startDate: string
    endDate: string
    interval?: string
  }) => {
    return request.get<any, TimeSeriesData[]>('/admin/statistics/revenue', { params })
  }
}

export const exportApi = {
  createTask: (data: {
    type: string
    fields?: string[]
    startDate?: string
    endDate?: string
    conditions?: Record<string, string>
  }) => {
    return request.post<any, ExportTask>('/exports', data)
  },

  getTasks: (params?: {
    page?: number
    pageSize?: number
    status?: string
  }) => {
    return request.get<any, PaginatedData<ExportTask>>('/exports', { params })
  },

  getTask: (id: number) => {
    return request.get<any, ExportTask>(`/exports/${id}`)
  },

  download: (id: number) => {
    return request.get(`/exports/${id}/download`, {
      responseType: 'blob'
    })
  }
}
