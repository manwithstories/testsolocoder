import { get, post, put } from '@/utils/request'
import type {
  User, LoginResponse, Product, ProductListParams, Order,
  RepairService, RepairOrder, Review, Report, Warranty,
  Notification, Message, WalletLog, Transaction, Favorite,
  DashboardStats, OrderStats, UserStats, ApiResponse, PaginationResponse
} from '@/types'

export const authApi = {
  register: (data: {
    username: string
    password: string
    email?: string
    phone?: string
    role: 'seller' | 'buyer' | 'technician'
    nickname?: string
  }) => post<ApiResponse<{ id: number; username: string; role: string }>>('/auth/register', data),

  login: (data: { username: string; password: string }) =>
    post<ApiResponse<LoginResponse>>('/auth/login', data),

  logout: () => post<ApiResponse<null>>('/auth/logout'),

  refreshToken: (data: { refreshToken: string }) =>
    post<ApiResponse<{ accessToken: string; tokenType: string; expiresIn: number }>>('/auth/refresh', data)
}

export const userApi = {
  getProfile: () => get<ApiResponse<User>>('/user/profile'),

  updateProfile: (data: {
    nickname?: string
    email?: string
    phone?: string
    avatar?: string
  }) => put<ApiResponse<null>>('/user/profile', data),

  changePassword: (data: { oldPassword: string; newPassword: string }) =>
    put<ApiResponse<null>>('/user/password', data),

  submitRealNameAuth: (data: { realName: string; idCard: string }) =>
    post<ApiResponse<null>>('/user/realname-auth', data),

  submitTechnicianCert: (data: { certType: string; certNumber: string; certImage?: string }) =>
    post<ApiResponse<null>>('/user/technician-cert', data),

  getWalletBalance: () => get<ApiResponse<{ balance: number }>>('/user/wallet/balance'),

  recharge: (data: { amount: number; paymentMethod: string }) =>
    post<ApiResponse<Transaction>>('/user/wallet/recharge', data),

  withdraw: (data: { amount: number }) =>
    post<ApiResponse<null>>('/user/wallet/withdraw', data),

  getWalletLogs: (params: { page?: number; pageSize?: number }) =>
    get<PaginationResponse<WalletLog>>('/user/wallet/logs', params),

  getUserStats: () => get<ApiResponse<UserStats>>('/user/stats')
}

export const productApi = {
  getList: (params: ProductListParams) =>
    get<PaginationResponse<Product>>('/products', params),

  getById: (id: number) => get<ApiResponse<Product>>(`/products/${id}`),

  getHot: (limit?: number) => get<ApiResponse<Product[]>>('/products/hot', { limit }),

  getCategories: () => get<ApiResponse<{ categories: string[]; conditions: string[] }>>('/products/categories'),

  create: (data: {
    title: string
    category: string
    brand: string
    model: string
    condition: string
    price: number
    originalPrice?: number
    description: string
    warrantyDays?: number
    images?: string
  }) => post<ApiResponse<Product>>('/products', data),

  update: (id: number, data: any) => put<ApiResponse<null>>(`/products/${id}`, data),

  remove: (id: number) => del<ApiResponse<null>>(`/products/${id}`),

  offShelf: (id: number) => post<ApiResponse<null>>(`/products/${id}/off-shelf`),

  getMyProducts: (params: { page?: number; pageSize?: number; status?: number }) =>
    get<PaginationResponse<Product>>('/my/products', params),

  getPending: (params: { page?: number; pageSize?: number }) =>
    get<PaginationResponse<Product>>('/products/pending', params),

  review: (id: number, data: { approved: boolean; rejectReason?: string }) =>
    post<ApiResponse<null>>(`/products/${id}/review`, data),

  toggleFavorite: (id: number) => post<ApiResponse<{ isFavorited: boolean }>>(`/favorites/${id}`),

  getFavorites: (params: { page?: number; pageSize?: number }) =>
    get<PaginationResponse<Favorite>>('/favorites', params)
}

export const orderApi = {
  create: (data: {
    productId: number
    receiverName: string
    receiverPhone: string
    receiverAddress: string
    negotiatedPrice?: number
  }) => post<ApiResponse<Order>>('/orders', data),

  getById: (id: number) => get<ApiResponse<Order>>(`/orders/${id}`),

  getByNo: (orderNo: string) => get<ApiResponse<Order>>(`/orders/by-no/${orderNo}`),

  getList: (params: { page?: number; pageSize?: number; status?: number }) =>
    get<PaginationResponse<Order>>('/buyer/orders', params),

  getSellerOrders: (params: { page?: number; pageSize?: number; status?: number }) =>
    get<PaginationResponse<Order>>('/seller/orders', params),

  pay: (data: { orderNo: string; paymentMethod: string }) =>
    post<ApiResponse<null>>(`/orders/${data.orderNo}/pay`, data),

  ship: (data: { orderNo: string; trackingNo: string; trackingCompany: string }) =>
    post<ApiResponse<null>>(`/orders/${data.orderNo}/ship`, data),

  confirm: (orderNo: string) => post<ApiResponse<null>>(`/orders/${orderNo}/confirm`),

  cancel: (orderNo: string) => post<ApiResponse<null>>(`/orders/${orderNo}/cancel`),

  refund: (data: { orderNo: string; reason?: string }) =>
    post<ApiResponse<null>>('/orders/refund', data),

  handleRefund: (data: { orderId: number; approved: boolean }) =>
    post<ApiResponse<null>>('/orders/refund/handle', data),

  negotiate: (data: { orderNo: string; offeredPrice: number; message?: string }) =>
    post<ApiResponse<null>>('/orders/negotiate', data),

  handleNegotiation: (data: {
    orderNo: string
    accepted: boolean
    counterPrice?: number
    message?: string
  }) => post<ApiResponse<null>>('/orders/negotiation/handle', data),

  getStats: () => get<ApiResponse<OrderStats>>('/order/stats')
}

export const repairApi = {
  getServiceList: (params: {
    page?: number
    pageSize?: number
    serviceType?: string
    keyword?: string
    sortBy?: string
    technicianId?: number
  }) => get<PaginationResponse<RepairService>>('/services', params),

  getServiceById: (id: number) => get<ApiResponse<RepairService>>(`/services/${id}`),

  getServiceTypes: () => get<ApiResponse<{ serviceTypes: string[] }>>('/services/types'),

  createService: (data: {
    serviceType: string
    title: string
    description: string
    price: number
    minPrice?: number
    maxPrice?: number
    estimatedDays?: number
    images?: string
  }) => post<ApiResponse<RepairService>>('/services', data),

  updateService: (id: number, data: any) => put<ApiResponse<null>>(`/services/${id}`, data),

  deleteService: (id: number) => del<ApiResponse<null>>(`/services/${id}`),

  createOrder: (data: {
    technicianId: number
    serviceId: number
    deviceType: string
    deviceBrand: string
    deviceModel: string
    faultDescription: string
    contactName: string
    contactPhone: string
    address?: string
    servicePrice: number
  }) => post<ApiResponse<RepairOrder>>('/repair-orders', data),

  getOrderById: (id: number) => get<ApiResponse<RepairOrder>>(`/repair-orders/${id}`),

  getOrderList: (params: { page?: number; pageSize?: number; status?: number }) =>
    get<PaginationResponse<RepairOrder>>('/buyer/repair-orders', params),

  getTechnicianOrders: (params: { page?: number; pageSize?: number; status?: number }) =>
    get<PaginationResponse<RepairOrder>>('/technician/repair-orders', params),

  acceptOrder: (id: number) => post<ApiResponse<null>>(`/repair-orders/${id}/accept`),

  startRepair: (id: number) => post<ApiResponse<null>>(`/repair-orders/${id}/start`),

  completeRepair: (id: number, data?: { finalPrice?: number }) =>
    post<ApiResponse<null>>(`/repair-orders/${id}/complete`, data),

  pickUpDevice: (id: number) => post<ApiResponse<null>>(`/repair-orders/${id}/pickup`),

  cancelOrder: (id: number) => post<ApiResponse<null>>(`/repair-orders/${id}/cancel`)
}

export const reviewApi = {
  create: (data: {
    orderId?: number
    repairOrderId?: number
    revieweeId: number
    reviewType: string
    rating: number
    content?: string
    images?: string
    qualityScore?: number
    serviceScore?: number
  }) => post<ApiResponse<Review>>('/reviews', data),

  getById: (id: number) => get<ApiResponse<Review>>(`/reviews/${id}`),

  getList: (params: {
    page?: number
    pageSize?: number
    revieweeId?: number
    reviewType?: string
    minRating?: number
  }) => get<PaginationResponse<Review>>('/reviews', params),

  getAverageRating: (userId: number, reviewType?: string) =>
    get<ApiResponse<{ userId: number; avgRating: number; reviewType?: string }>>(`/users/${userId}/rating`, { reviewType }),

  delete: (id: number) => del<ApiResponse<null>>(`/reviews/${id}`)
}

export const reportApi = {
  create: (data: {
    targetType: string
    targetId: number
    reason: string
    description?: string
    images?: string
  }) => post<ApiResponse<Report>>('/reports', data),

  getById: (id: number) => get<ApiResponse<Report>>(`/reports/${id}`),

  getList: (params: { page?: number; pageSize?: number; status?: number; targetType?: string }) =>
    get<PaginationResponse<Report>>('/reports', params),

  handle: (id: number, data: { approved: boolean; handleResult?: string }) =>
    post<ApiResponse<null>>(`/reports/${id}/handle`, data)
}

export const warrantyApi = {
  create: (data: {
    orderId?: number
    repairOrderId?: number
    type: string
    description: string
    images?: string
  }) => post<ApiResponse<Warranty>>('/warranties', data),

  getById: (id: number) => get<ApiResponse<Warranty>>(`/warranties/${id}`),

  getList: (params: { page?: number; pageSize?: number; status?: number }) =>
    get<PaginationResponse<Warranty>>('/warranties', params),

  handle: (id: number, data: { status: number; handleResult?: string }) =>
    post<ApiResponse<null>>(`/warranties/${id}/handle`, data)
}

export const notificationApi = {
  getList: (params: { page?: number; pageSize?: number; isRead?: string }) =>
    get<PaginationResponse<Notification>>('/notifications', params),

  getById: (id: number) => get<ApiResponse<Notification>>(`/notifications/${id}`),

  markAsRead: (id: number) => post<ApiResponse<null>>(`/notifications/${id}/read`),

  markAllAsRead: () => post<ApiResponse<null>>('/notifications/read-all'),

  getUnreadCount: () => get<ApiResponse<{ unreadCount: number }>>('/notifications/unread-count'),

  delete: (id: number) => del<ApiResponse<null>>(`/notifications/${id}`)
}

export const messageApi = {
  send: (data: { receiverId: number; content: string; type?: string }) =>
    post<ApiResponse<Message>>('/messages', data),

  getList: (userId: number, params: { page?: number; pageSize?: number }) =>
    get<PaginationResponse<Message>>(`/messages/${userId}`, params),

  markAsRead: (userId: number) => post<ApiResponse<null>>(`/messages/${userId}/read`),

  getUnreadCount: () => get<ApiResponse<{ unreadCount: number }>>('/messages/unread-count'),

  getContacts: () => get<ApiResponse<any[]>>('/message-contacts')
}

export const adminApi = {
  getDashboardStats: () => get<ApiResponse<DashboardStats>>('/admin/dashboard/stats'),

  getTransactionStats: (params: { startDate?: string; endDate?: string }) =>
    get<ApiResponse<any[]>>('/admin/transaction-stats', params),

  getUserActivityStats: (days?: number) =>
    get<ApiResponse<any[]>>('/admin/user-activity-stats', { days }),

  getUsers: (params: { page?: number; pageSize?: number; role?: string; status?: number }) =>
    get<PaginationResponse<User>>('/admin/users', params),

  updateUserStatus: (id: number, data: { status: number }) =>
    put<ApiResponse<null>>(`/admin/users/${id}/status`, data),

  reviewRealNameAuth: (id: number, data: { approved: boolean; rejectReason?: string }) =>
    post<ApiResponse<null>>(`/admin/users/${id}/review-realname`, data),

  reviewTechnicianCert: (id: number, data: { approved: boolean; rejectReason?: string }) =>
    post<ApiResponse<null>>(`/admin/technician-certs/${id}/review`, data),

  getTransactions: (params: { page?: number; pageSize?: number; type?: string; status?: number }) =>
    get<PaginationResponse<Transaction>>('/admin/transactions', params),

  getAdminLogs: (params: { page?: number; pageSize?: number }) =>
    get<PaginationResponse<any>>('/admin/admin-logs', params)
}

export const uploadApi = {
  uploadImage: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return post<ApiResponse<{ url: string; filename: string }>>('/upload', formData)
  },

  uploadMultipleImages: (files: File[]) => {
    const formData = new FormData()
    files.forEach(file => formData.append('files', file))
    return post<ApiResponse<{ urls: string[] }>>('/upload/multiple', formData)
  }
}
