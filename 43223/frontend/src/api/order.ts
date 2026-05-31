import { request } from './client'
import {
  Order,
  CartItem,
  PaginatedData,
  CreateOrderRequest,
} from '@/types'

export const orderApi = {
  list: (params?: { page?: number; page_size?: number; status?: string; order_no?: string }) =>
    request.get<PaginatedData<Order>>('/orders', { params }),

  get: (id: number) =>
    request.get<Order>(`/orders/${id}`),

  create: (data: CreateOrderRequest) =>
    request.post('/orders', data),

  updateStatus: (id: number, status: string) =>
    request.patch(`/orders/${id}/status`, { status }),

  cancel: (id: number) =>
    request.post(`/orders/${id}/cancel`),

  pay: (orderId: number, method: string) =>
    request.post('/orders/pay', { order_id: orderId, method }),
}

export const cartApi = {
  get: () =>
    request.get<{ items: CartItem[]; total_amount: number; total_count: number }>('/cart'),

  add: (productId: number, quantity: number) =>
    request.post('/cart', { product_id: productId, quantity }),

  update: (id: number, quantity: number) =>
    request.put(`/cart/${id}`, { quantity }),

  remove: (id: number) =>
    request.delete(`/cart/${id}`),

  clear: () =>
    request.delete('/cart'),
}
