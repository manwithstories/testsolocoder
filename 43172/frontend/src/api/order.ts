import { request } from './index'
import type { Order, ApiResponse } from '@/types'

export const orderApi = {
  listOrders(params?: {
    page?: number
    page_size?: number
    status?: string
  }) {
    return request.get<{ list: Order[]; total: number }>('/orders', { params })
  },

  getOrder(id: number) {
    return request.get<Order>(`/orders/${id}`)
  },

  getOrderByNumber(orderNumber: string) {
    return request.get<Order>(`/orders/number/${orderNumber}`)
  },

  createOrder(data: {
    product_id: number
    shipping_address: string
    need_auth?: boolean
    remark?: string
  }) {
    return request.post<Order>('/orders', data)
  },

  payOrder(id: number, data: { payment_method: string }) {
    return request.post<Order>(`/orders/${id}/pay`, data)
  },

  shipOrder(id: number, data: { tracking_number: string }) {
    return request.post<Order>(`/orders/${id}/ship`, data)
  },

  confirmDelivery(id: number) {
    return request.post<Order>(`/orders/${id}/confirm`)
  },

  cancelOrder(id: number, data?: { reason?: string }) {
    return request.post<Order>(`/orders/${id}/cancel`, data)
  }
}

export default orderApi
