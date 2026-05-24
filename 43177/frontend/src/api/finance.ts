import { request } from './index'
import type { ApiResponse, WithdrawRequest, Transaction, MonthlyReport, PaginationData } from '@/types'

export const financeApi = {
  getBalance: () => request.get<ApiResponse<{ balance: number; pending_withdraw: number; available_balance: number }>>('/finance/balance'),

  createWithdraw: (data: { amount: number; bank_account: string; bank_name: string }) =>
    request.post<ApiResponse>('/finance/withdraw', data),

  getWithdraws: (params?: { page?: number; page_size?: number; status?: string }) =>
    request.get<ApiResponse<PaginationData<WithdrawRequest>>>('/finance/withdraws', { params }),

  getWithdrawDetail: (id: number) =>
    request.get<ApiResponse<WithdrawRequest>>(`/finance/withdraws/${id}`),

  getTransactions: (params?: { page?: number; page_size?: number; type?: string }) =>
    request.get<ApiResponse<PaginationData<Transaction>>>('/finance/transactions', { params })
}

export const adminFinanceApi = {
  getWithdraws: (params?: { page?: number; page_size?: number; status?: string }) =>
    request.get<ApiResponse<PaginationData<WithdrawRequest>>>('/admin/withdraws', { params }),

  approveWithdraw: (id: number) =>
    request.post<ApiResponse>(`/admin/withdraws/${id}/approve`),

  rejectWithdraw: (id: number, data: { remark: string }) =>
    request.post<ApiResponse>(`/admin/withdraws/${id}/reject`, data),

  completeWithdraw: (id: number) =>
    request.post<ApiResponse>(`/admin/withdraws/${id}/complete`),

  getMonthlyReport: (params?: { month?: string }) =>
    request.get<ApiResponse<MonthlyReport>>('/admin/finance/report', { params }),

  getTechnicianPerformance: (params?: { month?: string }) =>
    request.get<ApiResponse<any[]>>('/admin/finance/performance', { params }),

  settleIncome: () =>
    request.post<ApiResponse>('/admin/finance/settle')
}
