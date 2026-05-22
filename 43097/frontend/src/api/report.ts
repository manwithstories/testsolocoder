import request from '@/utils/request'
import { ReportData } from '@/types'

export const getDailyReport = (startDate: string, endDate: string) => {
  return request.get<ReportData[]>('/reports/daily', { params: { startDate, endDate } })
}

export const getMonthlyReport = (year: number, month: number) => {
  return request.get('/reports/monthly', { params: { year, month } })
}

export const getYearlyReport = (year: number) => {
  return request.get('/reports/yearly', { params: { year } })
}

export const getRoomTypeReport = (startDate: string, endDate: string) => {
  return request.get('/reports/room-type', { params: { startDate, endDate } })
}

export const getMemberReport = (startDate: string, endDate: string) => {
  return request.get('/reports/member', { params: { startDate, endDate } })
}

export const getPaymentReport = (startDate: string, endDate: string) => {
  return request.get('/reports/payment', { params: { startDate, endDate } })
}

export const getOccupancyReport = (startDate: string, endDate: string) => {
  return request.get('/reports/occupancy', { params: { startDate, endDate } })
}

export const exportReport = (type: string, startDate: string, endDate: string) => {
  return request.get('/reports/export', {
    params: { type, startDate, endDate },
    responseType: 'blob'
  })
}
