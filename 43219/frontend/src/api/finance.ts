import request from './index'
import type { ApiResponse } from './index'

export interface Wallet {
  id: number
  user_id: number
  balance: number
  frozen: number
}

export interface Settlement {
  id: number
  order_id: number
  company_id: number
  staff_id: number
  total_amount: number
  company_share: number
  staff_share: number
  status: string
  created_at: string
}

export function requestWithdrawal(payload: { amount: number; account: string }) {
  return request.post<ApiResponse<string>>('/staff-area/withdraw', payload)
}

export function myEarnings() {
  return request.get<ApiResponse<{
    total_earned: number
    withdrawn: number
    balance: number
    settlements: Settlement[]
    withdrawals: any[]
  }>>('/staff-area/earnings')
}

export function walletInfo() {
  return request.get<ApiResponse<Wallet>>('/staff-area/wallet')
}

export function companyMonthly(month?: string) {
  return request.get<ApiResponse<any>>('/company/finance/monthly', { params: { month } })
}

export function exportFinanceCSV() {
  return request.get('/company/finance/export', { responseType: 'blob' })
}
