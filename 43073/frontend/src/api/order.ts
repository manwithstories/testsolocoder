import { request } from './request'

export interface OrderTicket {
  ticketTypeId: number
  quantity: number
}

export interface OrderData {
  activityId: number
  tickets: OrderTicket[]
  couponCode?: string
  remark?: string
}

export const createOrder = (data: OrderData) => {
  return request({
    url: '/orders',
    method: 'post',
    data
  })
}

export const getOrderList = (params: any) => {
  return request({
    url: '/orders',
    method: 'get',
    params
  })
}

export const getOrder = (id: number) => {
  return request({
    url: `/orders/${id}`,
    method: 'get'
  })
}

export const payOrder = (id: number) => {
  return request({
    url: `/orders/${id}/pay`,
    method: 'post'
  })
}

export const cancelOrder = (id: number) => {
  return request({
    url: `/orders/${id}/cancel`,
    method: 'post'
  })
}
