import { request } from './request'

export interface CouponData {
  type: string
  value: number
  minAmount: number
  totalCount: number
  startTime?: string
  endTime?: string
}

export const createCoupon = (data: CouponData) => {
  return request({
    url: '/coupons',
    method: 'post',
    data
  })
}

export const getCouponList = (params: any) => {
  return request({
    url: '/coupons',
    method: 'get',
    params
  })
}

export const deleteCoupon = (id: number) => {
  return request({
    url: `/coupons/${id}`,
    method: 'delete'
  })
}
