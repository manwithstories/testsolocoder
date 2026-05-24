import { request } from './index'
import type { ApiResponse, Order, OrderLog, PaginationData } from '@/types'

export const orderApi = {
  createOrder: (data: {
    service_item_id: number
    title: string
    description?: string
    images?: string
    address: string
    longitude?: number
    latitude?: number
    contact_name: string
    contact_phone: string
    appointment_time?: string
    urgent_level?: number
  }) => request.post<ApiResponse>('/orders', data),

  getOrders: (params?: { page?: number; page_size?: number; status?: string }) =>
    request.get<ApiResponse<PaginationData<Order>>>('/orders', { params }),

  getOrderDetail: (id: number) =>
    request.get<ApiResponse>(`/orders/${id}`),

  acceptOrder: (id: number, data?: { quoted_price?: number }) =>
    request.post<ApiResponse>(`/orders/${id}/accept`, data),

  arriveAtSite: (id: number) =>
    request.post<ApiResponse>(`/orders/${id}/arrive`),

  startRepair: (id: number) =>
    request.post<ApiResponse>(`/orders/${id}/start`),

  completeOrder: (id: number, data?: { final_price?: number; note?: string }) =>
    request.post<ApiResponse>(`/orders/${id}/complete`, data),

  cancelOrder: (id: number, data: { cancel_reason: string }) =>
    request.post<ApiResponse>(`/orders/${id}/cancel`, data),

  requestRefund: (id: number, data: { refund_reason: string; refund_amount?: number }) =>
    request.post<ApiResponse>(`/orders/${id}/refund`, data)
}

export const adminOrderApi = {
  getAllOrders: (params?: { page?: number; page_size?: number; status?: string }) =>
    request.get<ApiResponse<PaginationData<Order>>>('/admin/orders', { params }),

  getOrdersByStatus: (status: string, params?: { page?: number; page_size?: number }) =>
    request.get<ApiResponse<PaginationData<Order>>>(`/admin/orders/status/${status}`, { params }),

  getRefundList: (params?: { page?: number; page_size?: number }) =>
    request.get<ApiResponse<PaginationData<Order>>>('/admin/refunds', { params }),

  approveRefund: (id: number) =>
    request.post<ApiResponse>(`/admin/refunds/${id}/approve`),

  rejectRefund: (id: number, data: { reason: string }) =>
    request.post<ApiResponse>(`/admin/refunds/${id}/reject`, data),

  getOrderLogs: (params?: { order_id?: number; page?: number; page_size?: number }) =>
    request.get<ApiResponse<PaginationData<OrderLog>>>('/admin/order-logs', { params })
}
