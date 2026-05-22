import request from '@/utils/request'
import type { RentRecord, RepairOrder, UtilityFee, Notice, PaginationParams, PaginationResult, ApiResponse } from '@/types'

export const getRentBills = (params: PaginationParams & { month?: string; status?: number }) => {
  return request.get<any, ApiResponse<PaginationResult<RentRecord> & { totalAmount: number; paidAmount: number }>>('/rent/bills', { params })
}

export const getRentBill = (id: number) => {
  return request.get<any, ApiResponse<RentRecord>>(`/rent/bills/${id}`)
}

export const generateRentBills = (month?: string) => {
  return request.post<any, ApiResponse<{ generated: number; month: string }>>('/rent/generate', { month })
}

export const payRentBill = (id: number, data: { paidAmount: number; remark?: string }) => {
  return request.post<any, ApiResponse<null>>(`/rent/bills/${id}/pay`, data)
}

export const calculateLateFee = () => {
  return request.post<any, ApiResponse<{ processed: number }>>('/rent/calculate-late-fee')
}

export const getRepairs = (params: PaginationParams & { status?: number; category?: string; priority?: number }) => {
  return request.get<any, ApiResponse<PaginationResult<RepairOrder>>>('/repairs', { params })
}

export const getRepair = (id: number) => {
  return request.get<any, ApiResponse<RepairOrder>>(`/repairs/${id}`)
}

export const getMyRepairs = (params: PaginationParams) => {
  return request.get<any, ApiResponse<PaginationResult<RepairOrder>>>('/repairs/my', { params })
}

export const createRepair = (data: any) => {
  return request.post<any, ApiResponse<RepairOrder>>('/repairs', data)
}

export const assignRepair = (id: number, handlerId: number) => {
  return request.put<any, ApiResponse<null>>(`/repairs/${id}/assign`, { handlerId })
}

export const updateRepairStatus = (id: number, data: { status: number; processNote?: string }) => {
  return request.put<any, ApiResponse<null>>(`/repairs/${id}/status`, data)
}

export const getFees = (params: PaginationParams & { month?: string; type?: string; status?: number }) => {
  return request.get<any, ApiResponse<PaginationResult<UtilityFee> & { totalAmount: number; paidAmount: number }>>('/fees', { params })
}

export const createFee = (data: any) => {
  return request.post<any, ApiResponse<UtilityFee>>('/fees', data)
}

export const updateFee = (id: number, data: any) => {
  return request.put<any, ApiResponse<null>>(`/fees/${id}`, data)
}

export const deleteFee = (id: number) => {
  return request.delete<any, ApiResponse<null>>(`/fees/${id}`)
}

export const payFee = (id: number) => {
  return request.post<any, ApiResponse<null>>(`/fees/${id}/pay`)
}

export const batchGenerateFees = (data: any) => {
  return request.post<any, ApiResponse<{ generated: number }>>('/fees/batch', data)
}

export const getNotices = (params: PaginationParams & { type?: number; building?: string }) => {
  return request.get<any, ApiResponse<PaginationResult<Notice>>>('/notices', { params })
}

export const getNotice = (id: number) => {
  return request.get<any, ApiResponse<Notice>>(`/notices/${id}`)
}

export const createNotice = (data: any) => {
  return request.post<any, ApiResponse<Notice>>('/notices', data)
}

export const updateNotice = (id: number, data: any) => {
  return request.put<any, ApiResponse<null>>(`/notices/${id}`, data)
}

export const updateNoticeFields = (id: number, data: any) => {
  return request.patch<any, ApiResponse<null>>(`/notices/${id}`, data)
}

export const deleteNotice = (id: number) => {
  return request.delete<any, ApiResponse<null>>(`/notices/${id}`)
}

export const updateNoticeStatus = (id: number, status: number) => {
  return request.put<any, ApiResponse<null>>(`/notices/${id}/status`, { status })
}

export const getOverview = () => {
  return request.get<any, ApiResponse<{
    totalProperties: number
    rentedProperties: number
    occupancyRate: string
    totalIncome: number
    pendingRepairs: number
    activeContracts: number
  }>>('/stats/overview')
}

export const getOccupancyTrend = () => {
  return request.get<any, ApiResponse<{ months: string[]; rates: number[] }>>('/stats/occupancy-trend')
}

export const getIncomeTrend = () => {
  return request.get<any, ApiResponse<{ months: string[]; incomes: number[] }>>('/stats/income-trend')
}

export const getRepairStats = () => {
  return request.get<any, ApiResponse<{ byCategory: Record<string, number>; byStatus: Record<string, number>; total: number }>>('/stats/repair-stats')
}

export const exportRentRecords = (month?: string) => {
  return request.get('/stats/export/rent', { params: { month }, responseType: 'blob' })
}

export const exportRepairOrders = () => {
  return request.get('/stats/export/repairs', { responseType: 'blob' })
}
