import { request } from '@/utils/request'
import type { Billing, Transaction, PaginationParams, PaginationResponse } from '@/types'

export const generateMonthlyBilling = (data: { company_id: number; agency_id: number; period: string }): Promise<Billing> => {
  return request.post('/billings/generate', data)
}

export const getBilling = (id: number): Promise<Billing> => {
  return request.get(`/billings/${id}`)
}

export const getCompanyBillings = (params: PaginationParams): Promise<PaginationResponse<Billing>> => {
  return request.get('/company/billings', { params })
}

export const getAgencyBillings = (params: PaginationParams): Promise<PaginationResponse<Billing>> => {
  return request.get('/agency/billings', { params })
}

export const payBilling = (data: { billing_id: number; payment_method: string }): Promise<void> => {
  return request.post('/billings/pay', data)
}

export const recharge = (data: { company_id: number; amount: number; payment_method: string }): Promise<void> => {
  return request.post('/company/recharge', data)
}

export const getTransactions = (params: PaginationParams): Promise<PaginationResponse<Transaction>> => {
  return request.get('/company/transactions', { params })
}

export const getCompanyBalance = (): Promise<any> => {
  return request.get('/company/balance')
}
