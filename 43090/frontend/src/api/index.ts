import { request } from '@/utils/request'
import type {
  User,
  AuctionItem,
  AuctionSession,
  Bid,
  AutoBid,
  Order,
  Payment,
  Review,
  Notification,
  Statistics,
  PageResponse,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
} from '@/types'

export const authApi = {
  login: (data: LoginRequest) => request.post<LoginResponse>('/auth/login', data),
  register: (data: RegisterRequest) => request.post<User>('/auth/register', data),
  getUserInfo: () => request.get<User>('/user/info'),
  updateUserInfo: (data: Partial<User>) => request.put('/user/info', data),
  changePassword: (data: { old_password: string; new_password: string }) =>
    request.post('/user/change-password', data),
}

export const itemApi = {
  getList: (params: any) => request.get<PageResponse<AuctionItem>>('/items', { params }),
  getDetail: (id: number) => request.get<AuctionItem>(`/items/${id}`),
  create: (data: any) => request.post<AuctionItem>('/items', data),
  update: (id: number, data: any) => request.put(`/items/${id}`, data),
  delete: (id: number) => request.delete(`/items/${id}`),
  online: (id: number) => request.post(`/items/${id}/online`),
  offline: (id: number) => request.post(`/items/${id}/offline`),
  uploadImages: (id: number, formData: FormData) =>
    request.post(`/items/${id}/images`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),
  getMyItems: (params: any) => request.get<PageResponse<AuctionItem>>('/my/items', { params }),
  getBidHistory: (id: number, params: any) =>
    request.get<PageResponse<Bid>>(`/items/${id}/bids`, { params }),
  getCurrentBid: (id: number) => request.get<Bid>(`/items/${id}/current-bid`),
}

export const bidApi = {
  placeBid: (id: number, data: { amount: number; max_auto_bid?: number }) =>
    request.post<Bid>(`/items/${id}/bid`, data),
  getMyBids: (params: any) => request.get<PageResponse<Bid>>('/my/bids', { params }),
  setAutoBid: (data: { auction_item_id: number; max_price: number }) =>
    request.post<AutoBid>('/my/auto-bids', data),
  cancelAutoBid: (id: number) => request.delete(`/my/auto-bids/${id}`),
  getMyAutoBids: () => request.get<AutoBid[]>('/my/auto-bids'),
  getAutoBid: (itemId: number) => request.get<AutoBid>(`/items/${itemId}/auto-bid`),
}

export const sessionApi = {
  getList: (params: any) => request.get<PageResponse<AuctionSession>>('/sessions', { params }),
  getActive: () => request.get<AuctionSession[]>('/sessions/active'),
  getDetail: (id: number) => request.get<AuctionSession>(`/sessions/${id}`),
  create: (data: any) => request.post<AuctionSession>('/sessions', data),
  update: (id: number, data: any) => request.put(`/sessions/${id}`, data),
  addItems: (id: number, data: { auction_item_ids: number[] }) =>
    request.post(`/sessions/${id}/items`, data),
  removeItem: (sessionId: number, itemId: number) =>
    request.delete(`/sessions/${sessionId}/items/${itemId}`),
  start: (id: number) => request.post(`/sessions/${id}/start`),
  end: (id: number) => request.post(`/sessions/${id}/end`),
  cancel: (id: number) => request.post(`/sessions/${id}/cancel`),
}

export const orderApi = {
  create: (data: { auction_item_id: number; shipping_info?: string }) =>
    request.post<Order>('/orders', data),
  getDetail: (id: number) => request.get<Order>(`/orders/${id}`),
  getByNo: (orderNo: string) => request.get<Order>(`/orders/no/${orderNo}`),
  getBuyerOrders: (params: any) => request.get<PageResponse<Order>>('/my/orders/buyer', { params }),
  getSellerOrders: (params: any) =>
    request.get<PageResponse<Order>>('/my/orders/seller', { params }),
  pay: (id: number, data: { method: string }) => request.post<Payment>(`/orders/${id}/pay`, data),
  ship: (id: number, trackingNo: string) =>
    request.post(`/orders/${id}/ship`, null, { params: { tracking_no: trackingNo } }),
  confirmDelivery: (id: number) => request.post(`/orders/${id}/confirm-delivery`),
  complete: (id: number) => request.post(`/orders/${id}/complete`),
  getAllOrders: (params: any) => request.get<PageResponse<Order>>('/admin/orders', { params }),
}

export const notificationApi = {
  getMyNotifications: (params: any) =>
    request.get<PageResponse<Notification>>('/my/notifications', { params }),
  getUnreadCount: () => request.get<{ count: number }>('/my/notifications/unread-count'),
  markAsRead: (data: { notification_ids: number[] }) =>
    request.post('/my/notifications/mark-read', data),
  markAllAsRead: () => request.post('/my/notifications/mark-all-read'),
}

export const reviewApi = {
  create: (data: { order_id: number; rating: number; content?: string }) =>
    request.post<Review>('/reviews', data),
  getOrderReviews: (orderId: number) => request.get<Review[]>(`/reviews/order/${orderId}`),
  getUserReviews: (userId: number, params: any) =>
    request.get<PageResponse<Review>>(`/reviews/user/${userId}`, { params }),
  getUserRating: (userId: number) =>
    request.get<{ avg_rating: number; total: number }>(`/reviews/user/${userId}/rating`),
}

export const adminApi = {
  getUsers: (params: any) => request.get<PageResponse<User>>('/admin/users', { params }),
  updateUserStatus: (id: number, status: number) =>
    request.put(`/admin/users/${id}/status`, null, { params: { status } }),
  getStatistics: (params: any) => request.get<Statistics>('/admin/statistics/overall', { params }),
  getMyStatistics: () => request.get('/my/statistics'),
  exportOrders: (params: any) =>
    request.get('/admin/export/orders', { params, responseType: 'blob' }),
  exportBids: (params: any) =>
    request.get('/admin/export/bids', { params, responseType: 'blob' }),
}
