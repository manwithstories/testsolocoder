import request from '@/utils/request'
import type { ApiResponse, Payment, MemberCard, CustomerPackage, PageResult, PageParams } from '@/types'

export const getPayments = (params: PageParams & { pay_method?: string; start_date?: string; end_date?: string }) => {
  return request.get<ApiResponse<PageResult<Payment>>>('/payments', { params })
}

export const getPayment = (id: number) => {
  return request.get<ApiResponse<Payment>>(`/payments/${id}`)
}

export const createPayment = (data: {
  appointment_id: number
  customer_id: number
  pay_method: string
  amount?: number
  points_used?: number
  card_id?: number
  package_id?: number
}) => {
  return request.post<ApiResponse<Payment>>('/payments', data)
}

export const createMemberCard = (data: { customer_id: number; card_type?: string; balance?: number; discount?: number }) => {
  return request.post<ApiResponse<MemberCard>>('/payments/member-card', data)
}

export const getMemberCards = (id: number) => {
  return request.get<ApiResponse<MemberCard[]>>(`/payments/member-card/${id}`)
}

export const rechargeCard = (id: number, data: { amount: number }) => {
  return request.post<ApiResponse<null>>(`/payments/member-card/${id}/recharge`, data)
}

export const purchasePackage = (data: { customer_id: number; service_id: number }) => {
  return request.post<ApiResponse<CustomerPackage>>('/payments/package', data)
}

export const getCustomerPackages = (id: number) => {
  return request.get<ApiResponse<CustomerPackage[]>>(`/payments/package/${id}`)
}
