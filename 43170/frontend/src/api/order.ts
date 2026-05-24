import api from './index'
import type { ApiResponse, Order, CreateOrderRequest, Review, CreateReviewRequest, Settlement } from '@/types'

export const orderApi = {
  createOrder: (data: CreateOrderRequest) => {
    return api.post<any, ApiResponse<Order>>('/orders/', data)
  },

  getOrder: (id: number) => {
    return api.get<any, ApiResponse<Order>>(`/orders/${id}`)
  },

  getMyOrders: (status?: string) => {
    const params = status ? { status } : {}
    return api.get<any, ApiResponse<Order[]>>('/orders/', { params })
  },

  confirmOrder: (id: number) => {
    return api.put<any, ApiResponse<null>>(`/orders/${id}/confirm`)
  },

  rejectOrder: (id: number, reason: string) => {
    return api.put<any, ApiResponse<null>>(`/orders/${id}/reject`, { reason })
  },

  startRental: (id: number) => {
    return api.put<any, ApiResponse<null>>(`/orders/${id}/start`)
  },

  completeOrder: (id: number) => {
    return api.put<any, ApiResponse<null>>(`/orders/${id}/complete`)
  },

  cancelOrder: (id: number) => {
    return api.put<any, ApiResponse<null>>(`/orders/${id}/cancel`)
  },

  getOrderStatusHistory: (id: number) => {
    return api.get<any, ApiResponse<any>>(`/orders/${id}/status-history`)
  }
}

export const reviewApi = {
  createReview: (data: CreateReviewRequest) => {
    return api.post<any, ApiResponse<Review>>('/reviews/', data)
  },

  getEquipmentReviews: (equipmentId: number, page = 1, pageSize = 10) => {
    return api.get<any, ApiResponse<any>>(`/reviews/equipment/${equipmentId}`, {
      params: { page, pageSize }
    })
  },

  getUserReviews: (userId: number) => {
    return api.get<any, ApiResponse<Review[]>>(`/reviews/user/${userId}`)
  },

  getMyReviews: () => {
    return api.get<any, ApiResponse<Review[]>>('/reviews/my')
  }
}

export const settlementApi = {
  createSettlement: (data: any) => {
    return api.post<any, ApiResponse<Settlement>>('/settlements/', data)
  },

  getSettlement: (id: number) => {
    return api.get<any, ApiResponse<Settlement>>(`/settlements/${id}`)
  },

  getOrderSettlement: (orderId: number) => {
    return api.get<any, ApiResponse<Settlement>>(`/settlements/order/${orderId}`)
  },

  getMySettlements: () => {
    return api.get<any, ApiResponse<Settlement[]>>('/settlements/')
  }
}

export const exportApi = {
  exportRentalRecords: (format: string = 'csv') => {
    return api.get<any, any>('/export/rentals', {
      params: { format },
      responseType: 'blob'
    })
  },

  exportRevenueReport: (format: string = 'csv', startDate?: string, endDate?: string) => {
    return api.get<any, any>('/export/revenue', {
      params: { format, startDate, endDate },
      responseType: 'blob'
    })
  }
}
