import request from './request'
import type { Transaction, Settlement, FinancialSummary, PaginatedResponse } from '@/types/finance'

export const getTransactionsApi = (params?: { transaction_type?: string; status?: string; start_date?: string; end_date?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<Transaction>>('/transactions', { params })
}

export const getTransactionApi = (id: string) => {
  return request.get<Transaction>(`/transactions/${id}`)
}

export const createTransactionApi = (data: any) => {
  return request.post<Transaction>('/transactions', data)
}

export const updateTransactionStatusApi = (id: string, status: string) => {
  return request.put<Transaction>(`/transactions/${id}/status`, { status })
}

export const getSettlementsApi = (params?: { status?: string; page?: number; page_size?: number }) => {
  return request.get<PaginatedResponse<Settlement>>('/settlements', { params })
}

export const createSettlementApi = (data: { month: number; year: number }) => {
  return request.post<Settlement>('/settlements', data)
}

export const getFinancialSummaryApi = (params?: { start_date?: string; end_date?: string }) => {
  return request.get<FinancialSummary>('/financial-summary', { params })
}

export const exportFinancialReportApi = (params: { start_date: string; end_date: string; user_id?: string; transaction_type?: string; format: 'pdf' | 'csv' }) => {
  return request.get('/financial-report/export', { params, responseType: 'blob' })
}

export const exportMonthlyReportApi = (params: { year: number; month: number; user_id?: string; format: 'pdf' | 'csv' }) => {
  return request.get('/monthly-report/export', { params, responseType: 'blob' })
}
