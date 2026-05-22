import request from './request'
import type {
  Booking, CreateBookingRequest, PriceCalculation,
  Order, Review, CreateReviewRequest,
  MaintenancePlan, Message, DashboardStats,
  PromoCode, PricingRule,
  PageResult, ApiResponse
} from '@/types'
import { userApi } from './user'
import { carApi } from './car'
import { storeApi } from './store'

export { userApi, carApi, storeApi }

export const bookingApi = {
  createBooking: (data: CreateBookingRequest) =>
    request.post<any, ApiResponse<Booking>>('/bookings', data),

  getBookings: (params?: {
    page?: number; page_size?: number; user_id?: number;
    status?: string; car_id?: number; start_date?: string; end_date?: string
  }) =>
    request.get<any, ApiResponse<PageResult<Booking>>>('/bookings', { params }),

  getMyBookings: (params?: { page?: number; page_size?: number }) =>
    request.get<any, ApiResponse<PageResult<Booking>>>('/bookings', { params }),

  getBookingById: (id: number) =>
    request.get<any, ApiResponse<Booking>>(`/bookings/${id}`),

  getBookingByNo: (no: string) =>
    request.get<any, ApiResponse<Booking>>(`/bookings/no/${no}`),

  confirmBooking: (id: number) =>
    request.put<any, ApiResponse<null>>(`/admin/bookings/${id}/confirm`),

  cancelBooking: (id: number, reason?: string) =>
    request.put<any, ApiResponse<null>>(`/bookings/${id}/cancel`, { reason }),

  completeBooking: (id: number) =>
    request.put<any, ApiResponse<null>>(`/admin/bookings/${id}/complete`),

  checkAvailability: (params: { car_id: number; pickup_time: string; return_time: string }) =>
    request.get<any, ApiResponse<{ available: boolean }>>('/bookings/check-availability', { params }),

  calculatePrice: (params: {
    car_id: number; pickup_time: string; return_time: string; promo_code?: string
  }) =>
    request.get<any, ApiResponse<PriceCalculation>>('/bookings/calculate-price', { params })
}

export const orderApi = {
  getOrders: (params?: {
    page?: number; page_size?: number; user_id?: number;
    status?: string; start_date?: string; end_date?: string
  }) =>
    request.get<any, ApiResponse<PageResult<Order>>>('/admin/orders', { params }),

  getMyOrders: (params?: { page?: number; page_size?: number }) =>
    request.get<any, ApiResponse<PageResult<Order>>>('/orders', { params }),

  getOrderById: (id: number) =>
    request.get<any, ApiResponse<Order>>(`/orders/${id}`),

  getOrderByNo: (no: string) =>
    request.get<any, ApiResponse<Order>>(`/orders/no/${no}`),

  updateOrderStatus: (id: number, status: string) =>
    request.put<any, ApiResponse<null>>(`/admin/orders/${id}/status`, { status }),

  refundOrder: (id: number, reason?: string) =>
    request.put<any, ApiResponse<null>>(`/admin/orders/${id}/refund`, { reason }),

  exportOrders: (params?: { status?: string; start_date?: string; end_date?: string }) =>
    request.get<any, ApiResponse<{ file_path: string }>>('/admin/orders/export', { params })
}

export const reviewApi = {
  createReview: (data: CreateReviewRequest) =>
    request.post<any, ApiResponse<Review>>('/reviews', data),

  getCarReviews: (carId: number, params?: { page?: number; page_size?: number }) =>
    request.get<any, ApiResponse<PageResult<Review>>>(`/cars/${carId}/reviews`, { params }),

  getMyReviews: (params?: { page?: number; page_size?: number }) =>
    request.get<any, ApiResponse<PageResult<Review>>>('/my-reviews', { params }),

  getAllReviews: (params?: {
    page?: number; page_size?: number; car_id?: number; min_rating?: number
  }) =>
    request.get<any, ApiResponse<PageResult<Review>>>('/admin/reviews', { params }),

  updateReview: (id: number, data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>(`/reviews/${id}`, data),

  deleteReview: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/reviews/${id}`),

  toggleReviewHidden: (id: number, isHidden: boolean) =>
    request.put<any, ApiResponse<null>>(`/admin/reviews/${id}/hidden`, { is_hidden: isHidden }),

  likeReview: (id: number) =>
    request.post<any, ApiResponse<null>>(`/reviews/${id}/like`)
}

export const maintenanceApi = {
  createMaintenance: (data: Record<string, any>) =>
    request.post<any, ApiResponse<MaintenancePlan>>('/admin/maintenance', data),

  getMaintenance: (params?: {
    page?: number; page_size?: number; car_id?: number; status?: string
  }) =>
    request.get<any, ApiResponse<PageResult<MaintenancePlan>>>('/admin/maintenance', { params }),

  getMaintenanceById: (id: number) =>
    request.get<any, ApiResponse<MaintenancePlan>>(`/admin/maintenance/${id}`),

  getCarMaintenance: (carId: number, params?: { page?: number; page_size?: number }) =>
    request.get<any, ApiResponse<PageResult<MaintenancePlan>>>(`/admin/cars/${carId}/maintenance`, { params }),

  updateMaintenance: (id: number, data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>(`/admin/maintenance/${id}`, data),

  startMaintenance: (id: number) =>
    request.put<any, ApiResponse<null>>(`/admin/maintenance/${id}/start`),

  completeMaintenance: (id: number) =>
    request.put<any, ApiResponse<null>>(`/admin/maintenance/${id}/complete`),

  cancelMaintenance: (id: number) =>
    request.put<any, ApiResponse<null>>(`/admin/maintenance/${id}/cancel`),

  deleteMaintenance: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/admin/maintenance/${id}`),

  getUpcomingMaintenance: () =>
    request.get<any, ApiResponse<MaintenancePlan[]>>('/admin/maintenance/upcoming')
}

export const messageApi = {
  getMessages: (params?: { page?: number; page_size?: number }) =>
    request.get<any, ApiResponse<PageResult<Message>>>('/messages', { params }),

  getUnreadCount: () =>
    request.get<any, ApiResponse<{ count: number }>>('/messages/unread-count'),

  markAsRead: (id: number) =>
    request.put<any, ApiResponse<null>>(`/messages/${id}/read`),

  markAllAsRead: () =>
    request.put<any, ApiResponse<null>>('/messages/read-all'),

  deleteMessage: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/messages/${id}`)
}

export const statsApi = {
  getDashboardStats: () =>
    request.get<any, ApiResponse<DashboardStats>>('/admin/stats/dashboard'),

  getRevenueStats: (params?: { start_date?: string; end_date?: string }) =>
    request.get<any, ApiResponse<Record<string, any>>>('/admin/stats/revenue', { params })
}

export const promoApi = {
  createPromo: (data: Record<string, any>) =>
    request.post<any, ApiResponse<PromoCode>>('/admin/promos', data),

  getPromos: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    request.get<any, ApiResponse<PageResult<PromoCode>>>('/admin/promos', { params }),

  getPromoById: (id: number) =>
    request.get<any, ApiResponse<PromoCode>>(`/admin/promos/${id}`),

  getPromoByCode: (code: string) =>
    request.get<any, ApiResponse<PromoCode>>(`/promos/${code}`),

  updatePromo: (id: number, data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>(`/admin/promos/${id}`, data),

  deletePromo: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/admin/promos/${id}`)
}

export const pricingApi = {
  createRule: (data: Record<string, any>) =>
    request.post<any, ApiResponse<PricingRule>>('/admin/pricing-rules', data),

  getRules: (params?: { page?: number; page_size?: number; type?: string }) =>
    request.get<any, ApiResponse<PageResult<PricingRule>>>('/admin/pricing-rules', { params }),

  getRuleById: (id: number) =>
    request.get<any, ApiResponse<PricingRule>>(`/admin/pricing-rules/${id}`),

  updateRule: (id: number, data: Record<string, any>) =>
    request.put<any, ApiResponse<null>>(`/admin/pricing-rules/${id}`, data),

  deleteRule: (id: number) =>
    request.delete<any, ApiResponse<null>>(`/admin/pricing-rules/${id}`),

  toggleRuleActive: (id: number) =>
    request.put<any, ApiResponse<null>>(`/admin/pricing-rules/${id}/toggle`)
}
