import request from './request'
import { Bill, WithdrawRequest, PaginatedResponse, CreateWithdrawParams } from '@/types'

export const billApi = {
  getList: (params?: {
    page?: number
    page_size?: number
    type?: string
    start_date?: string
    end_date?: string
  }) => {
    return request.get<any, PaginatedResponse<Bill> & {
      income_total: number
      withdraw_total: number
      penalty_total: number
    }>('/bills', { params })
  },

  getBalance: () => {
    return request.get<any, {
      balance: number
      total_income: number
      pending_income: number
      order_count: number
    }>('/bills/balance')
  },

  getIncomeSummary: (params?: { month?: string }) => {
    return request.get<any, {
      month: string
      total_income: number
      order_count: number
      total_platform_fee: number
    }>('/bills/income-summary', { params })
  },
}

export const withdrawApi = {
  getList: (params?: { page?: number; page_size?: number; status?: string }) => {
    return request.get<any, PaginatedResponse<WithdrawRequest>>('/withdraws', { params })
  },

  create: (params: CreateWithdrawParams) => {
    return request.post<any, WithdrawRequest>('/withdraws', params)
  },

  handle: (id: number, params: { approved: boolean; remark?: string; transfer_no?: string }) => {
    return request.put<any, any>(`/withdraws/${id}/handle`, params)
  },
}
