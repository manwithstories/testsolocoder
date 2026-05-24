import request from '@/utils/request'

export interface SalesTrendItem {
  date: string
  amount: number
  count: number
}

export interface CustomerProfile {
  label: string
  value: number
}

export interface StatisticsOverview {
  totalOrders: number
  totalAmount: number
  totalDesigns: number
  totalTickets: number
  avgOrderAmount: number
  avgReviewScore: number
}

export interface StatisticsQuery {
  startDate?: string
  endDate?: string
  type?: 'daily' | 'weekly' | 'monthly'
}

export const getOverview = (params?: StatisticsQuery) => {
  return request.get<any, StatisticsOverview>('/statistics/overview', { params })
}

export const getSalesTrend = (params?: StatisticsQuery) => {
  return request.get<any, SalesTrendItem[]>('/statistics/sales-trend', { params })
}

export const getCustomerProfile = () => {
  return request.get<any, { areaDistribution: CustomerProfile[]; houseTypeDistribution: CustomerProfile[]; budgetDistribution: CustomerProfile[] }>('/statistics/customer-profile')
}

export const exportExcel = (params?: StatisticsQuery) => {
  return request.get<any, Blob>('/statistics/export/excel', { params, responseType: 'blob' })
}

export const exportPdf = (params?: StatisticsQuery) => {
  return request.get<any, Blob>('/statistics/export/pdf', { params, responseType: 'blob' })
}
