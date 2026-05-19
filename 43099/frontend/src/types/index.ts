export interface User {
  id: number
  username: string
  email: string
  real_name: string
  phone: string
  avatar: string
  role: 'user' | 'admin' | 'super_admin'
  email_verified: boolean
  status: 'active' | 'inactive'
  created_at: string
  updated_at: string
}

export interface Venue {
  id: number
  name: string
  location: string
  capacity: number
  facilities: string
  description: string
  cover_image: string
  status: 'online' | 'offline'
  created_by: number
  created_at: string
  updated_at: string
}

export interface DeviceCategory {
  id: number
  name: string
  description: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Device {
  id: number
  category_id: number
  name: string
  description: string
  specification: string
  stock_quantity: number
  available_quantity: number
  rental_price: number
  deposit_amount: number
  status: 'online' | 'offline'
  created_at: string
  updated_at: string
  category?: DeviceCategory
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  type: 'venue' | 'device'
  item_id: number
  item_name: string
  start_time: string
  end_time: string
  total_hours: number
  quantity: number
  total_amount: number
  deposit_amount: number
  status: 'pending' | 'confirmed' | 'paid' | 'completed' | 'cancelled'
  purpose: string
  contact_name: string
  contact_phone: string
  cancel_reason?: string
  cancelled_at?: string
  reviewed_by?: number
  review_note?: string
  reviewed_at?: string
  created_at: string
  updated_at: string
  user?: User
}

export interface Payment {
  id: number
  order_id: number
  transaction_no: string
  amount: number
  payment_method: 'wechat' | 'alipay' | 'cash'
  status: 'pending' | 'success' | 'failed' | 'refunded'
  paid_at?: string
  created_at: string
  order?: Order
}

export interface Review {
  id: number
  order_id: number
  user_id: number
  rating: number
  content: string
  status: 'pending' | 'approved' | 'rejected'
  review_note?: string
  reviewed_by?: number
  reviewed_at?: string
  created_at: string
  order?: Order
  user?: User
}

export interface CalendarEvent {
  id: number
  title: string
  start: string
  end: string
  status: string
  type: string
  item_id: number
  item_name: string
  color: string
}

export interface StatsOverview {
  total_bookings: number
  total_revenue: number
  total_users: number
  total_venues: number
  total_devices: number
  pending_orders: number
  today_bookings: number
  today_revenue: number
}

export interface BookingStats {
  date: string
  count: number
}

export interface RevenueStats {
  date: string
  amount: number
}

export interface PopularVenue {
  id: number
  name: string
  bookings: number
  revenue: number
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginationResponse<T> {
  total: number
  page: number
  page_size: number
  list: T[]
}

export interface PaginationParams {
  page?: number
  page_size?: number
}
