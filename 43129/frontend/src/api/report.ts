import request from '@/utils/request'
import type { ApiResponse, ReportData, RevenueReport, TechnicianPerformance, ServiceRank } from '@/types'

export const getRevenueReport = (params: { start_date: string; end_date: string }) => {
  return request.get<ApiResponse<RevenueReport>>('/reports/revenue', { params })
}

export const getTechnicianPerformance = (params: { start_date: string; end_date: string }) => {
  return request.get<ApiResponse<TechnicianPerformance[]>>('/reports/technician-performance', { params })
}

export const getServiceRanking = (params: { start_date: string; end_date: string }) => {
  return request.get<ApiResponse<ServiceRank[]>>('/reports/service-ranking', { params })
}

export const getFullReport = (params: { start_date: string; end_date: string }) => {
  return request.get<ApiResponse<ReportData>>('/reports/full', { params })
}

export const exportExcel = (params: { start_date: string; end_date: string }) => {
  return request.get('/reports/export/excel', { 
    params, 
    responseType: 'blob' as const 
  })
}

export const exportPDF = (params: { start_date: string; end_date: string }) => {
  return request.get('/reports/export/pdf', { 
    params, 
    responseType: 'blob' as const 
  })
}
