import request from '@/utils/request'
import type {
  Order,
  ListOrdersParams,
  QuoteOrderParams,
  OrderActionParams,
  PaginationResult
} from '@/types'

export interface OrderListResult extends PaginationResult<Order> {}

export const listOrders = (params: ListOrdersParams) => {
  return request.get<any, OrderListResult>('/orders/', { params })
}

export const getOrder = (id: number | string) => {
  return request.get<any, Order>(`/orders/${id}`)
}

export const quoteOrder = (id: number | string, data: QuoteOrderParams) => {
  return request.post<any, Order>(`/orders/${id}/quote`, data)
}

export const confirmOrder = (id: number | string, data: OrderActionParams = {}) => {
  return request.post<any, Order>(`/orders/${id}/confirm`, data)
}

export const cancelOrder = (id: number | string, data: OrderActionParams = {}) => {
  return request.post<any, Order>(`/orders/${id}/cancel`, data)
}

export const startProduceOrder = (id: number | string, data: OrderActionParams = {}) => {
  return request.post<any, Order>(`/orders/${id}/produce`, data)
}

export const shipOrder = (id: number | string, data: OrderActionParams = {}) => {
  return request.post<any, Order>(`/orders/${id}/ship`, data)
}

export const completeOrder = (id: number | string, data: OrderActionParams = {}) => {
  return request.post<any, Order>(`/orders/${id}/complete`, data)
}
