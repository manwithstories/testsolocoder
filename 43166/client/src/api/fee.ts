import request from '@/utils/request'
import {
  ApplicationFee,
  FeeStandard,
  DiscountPolicy,
  CalculateFeeRequest,
  PayFeeRequest,
  PaginatedData
} from '@/types'

export const feeApi = {
  calculate: (data: CalculateFeeRequest) => {
    return request.post<any, ApplicationFee>('/fees/calculate', data)
  },

  create: (data: {
    applicationId: number
    companyType: string
    capital: number
    discountCode?: string
  }) => {
    return request.post<any, ApplicationFee>('/fees', data)
  },

  pay: (data: PayFeeRequest) => {
    return request.post<any, ApplicationFee>('/fees/pay', data)
  },

  getByApplicationId: (applicationId: number) => {
    return request.get<any, ApplicationFee>(`/fees/${applicationId}`)
  },

  getList: (params?: {
    page?: number
    pageSize?: number
    status?: string
  }) => {
    return request.get<any, PaginatedData<ApplicationFee>>('/fees', { params })
  },

  getStandards: () => {
    return request.get<any, FeeStandard[]>('/admin/fee-standards')
  },

  createStandard: (data: Partial<FeeStandard>) => {
    return request.post<any, FeeStandard>('/admin/fee-standards', data)
  },

  updateStandard: (id: number, data: Partial<FeeStandard>) => {
    return request.put<any, null>(`/admin/fee-standards/${id}`, data)
  },

  getDiscounts: () => {
    return request.get<any, DiscountPolicy[]>('/admin/discounts')
  },

  createDiscount: (data: Partial<DiscountPolicy>) => {
    return request.post<any, DiscountPolicy>('/admin/discounts', data)
  },

  updateDiscount: (id: number, data: Partial<DiscountPolicy>) => {
    return request.put<any, null>(`/admin/discounts/${id}`, data)
  },

  deleteDiscount: (id: number) => {
    return request.delete<any, null>(`/admin/discounts/${id}`)
  }
}
