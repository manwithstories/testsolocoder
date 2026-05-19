import request from '@/utils/request'
import type {
  User,
  Venue,
  Device,
  DeviceCategory,
  Order,
  Payment,
  Review,
  CalendarEvent,
  StatsOverview,
  BookingStats,
  RevenueStats,
  PopularVenue,
  PaginationResponse,
  PaginationParams,
} from '@/types'

export const authApi = {
  register: (data: { username: string; email: string; password: string; real_name?: string; phone?: string }) =>
    request.post('/auth/register', data),
  login: (data: { username: string; password: string }) =>
    request.post('/auth/login', data),
  logout: () => request.post('/auth/logout'),
  sendVerifyEmail: (email: string) => request.post('/auth/send-verify-email', { email }),
  verifyEmail: (data: { email: string; code: string }) => request.post('/auth/verify-email', data),
  forgotPassword: (email: string) => request.post('/auth/forgot-password', { email }),
  resetPassword: (data: { email: string; code: string; new_password: string }) =>
    request.post('/auth/reset-password', data),
}

export const userApi = {
  getMe: () => request.get<User>('/users/me'),
  updateMe: (data: { real_name?: string; phone?: string; avatar?: string }) =>
    request.put<User>('/users/me', data),
  list: (params?: PaginationParams) =>
    request.get<PaginationResponse<User>>('/users', { params }),
  updateRole: (id: number, role: string) => request.put(`/users/${id}/role`, { role }),
}

export const venueApi = {
  list: (params?: PaginationParams & { status?: string }) =>
    request.get<PaginationResponse<Venue>>('/venues', { params }),
  get: (id: number) => request.get<Venue>(`/venues/${id}`),
  create: (data: Partial<Venue>) => request.post<Venue>('/venues', data),
  update: (id: number, data: Partial<Venue>) => request.put<Venue>(`/venues/${id}`, data),
  delete: (id: number) => request.delete(`/venues/${id}`),
  updateStatus: (id: number, status: string) =>
    request.patch(`/venues/${id}/status`, { status }),
  setPrice: (id: number, data: { day_of_week: number; time_slots: { start: string; end: string; price: number }[] }) =>
    request.post(`/venues/${id}/prices`, data),
  getAvailability: (id: number, date: string) =>
    request.get(`/venues/${id}/availability?date=${date}`),
}

export const deviceApi = {
  listCategories: () => request.get<DeviceCategory[]>('/devices/categories'),
  createCategory: (data: { name: string; description?: string; sort_order?: number }) =>
    request.post<DeviceCategory>('/devices/categories', data),
  list: (params?: PaginationParams & { category_id?: number; status?: string }) =>
    request.get<PaginationResponse<Device>>('/devices', { params }),
  get: (id: number) => request.get<Device>(`/devices/${id}`),
  create: (data: Partial<Device>) => request.post<Device>('/devices', data),
  update: (id: number, data: Partial<Device>) => request.put<Device>(`/devices/${id}`, data),
  updateStatus: (id: number, status: string) =>
    request.patch(`/devices/${id}/status`, { status }),
  getAvailability: (id: number, date: string) =>
    request.get(`/devices/${id}/availability?date=${date}`),
  batchImport: (data: any[]) => request.post('/devices/batch-import', data),
}

export const orderApi = {
  create: (data: {
    type: string
    item_id: number
    start_time: string
    end_time: string
    quantity?: number
    purpose?: string
    contact_name: string
    contact_phone: string
  }) => request.post<Order>('/bookings', data),
  getCalendar: (params: { start_date: string; end_date: string; type?: string; item_id?: number }) =>
    request.get<CalendarEvent[]>('/bookings/calendar', { params }),
  list: (params?: PaginationParams & { status?: string; type?: string }) =>
    request.get<PaginationResponse<Order>>('/orders', { params }),
  get: (id: number) => request.get<Order>(`/orders/${id}`),
  cancel: (id: number, reason?: string) =>
    request.put(`/orders/${id}/cancel`, { reason }),
  confirm: (id: number, note?: string) =>
    request.put(`/orders/${id}/confirm`, { note }),
  complete: (id: number) => request.put(`/orders/${id}/complete`),
}

export const paymentApi = {
  list: (params?: PaginationParams & { status?: string }) =>
    request.get<PaginationResponse<Payment>>('/payments', { params }),
  confirm: (id: number, data: { transaction_no: string; payment_method: string; amount: number }) =>
    request.post<Payment>(`/payments/${id}/confirm`, data),
  export: (params: { start_date: string; end_date: string }) =>
    request.get('/payments/export', { params, responseType: 'blob' }),
}

export const reviewApi = {
  create: (data: { order_id: number; rating: number; content?: string }) =>
    request.post<Review>('/reviews', data),
  list: (params?: PaginationParams & { status?: string }) =>
    request.get<PaginationResponse<Review>>('/reviews', { params }),
  get: (id: number) => request.get<Review>(`/reviews/${id}`),
  approve: (id: number, note?: string) =>
    request.put(`/reviews/${id}/approve`, { note }),
  reject: (id: number, note?: string) =>
    request.put(`/reviews/${id}/reject`, { note }),
}

export const statsApi = {
  getOverview: () => request.get<StatsOverview>('/stats/overview'),
  getBookingStats: (params: { start_date: string; end_date: string }) =>
    request.get<BookingStats[]>('/stats/bookings', { params }),
  getRevenueStats: (params: { start_date: string; end_date: string }) =>
    request.get<RevenueStats[]>('/stats/revenue', { params }),
  getPopularVenues: (params: { start_date: string; end_date: string }) =>
    request.get<PopularVenue[]>('/stats/popular-venues', { params }),
}
