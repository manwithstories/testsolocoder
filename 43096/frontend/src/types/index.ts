export interface User {
  id: number
  username: string
  email: string
  phone: string
  real_name: string
  id_card: string
  member_level: number
  points: number
  status: number
  avatar: string
  role: string
  created_at: string
  updated_at: string
}

export interface MemberLevel {
  id: number
  name: string
  level: number
  discount: number
  min_points: number
  priority: number
  description: string
}

export interface Coupon {
  id: number
  code: string
  name: string
  type: number
  value: number
  min_amount: number
  user_id: number
  status: number
  expire_at: string
}

export interface Show {
  id: number
  name: string
  description: string
  poster: string
  artist: string
  duration: number
  status: number
  organizer: string
  venue: string
  address: string
  sessions?: Session[]
  created_at: string
  updated_at: string
}

export interface Session {
  id: number
  show_id: number
  start_time: string
  end_time: string
  status: number
  total_seats: number
  sold_seats: number
  seat_areas?: SeatArea[]
}

export interface SeatArea {
  id: number
  session_id: number
  name: string
  color: string
  price: number
  total_seats: number
  sold_seats: number
  sort_order: number
}

export interface Seat {
  id: number
  session_id: number
  area_id: number
  row: string
  col: number
  seat_no: string
  status: number
  x: number
  y: number
  width: number
  height: number
}

export interface SeatChart {
  id: number
  session_id: number
  config: string
  background: string
  width: number
  height: number
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  show_id: number
  session_id: number
  total_amount: number
  discount: number
  pay_amount: number
  status: number
  pay_type: number
  pay_time: string | null
  real_name: string
  id_card: string
  phone: string
  email: string
  remark: string
  tickets?: Ticket[]
  refund?: Refund
  created_at: string
}

export interface Ticket {
  id: number
  order_id: number
  user_id: number
  show_id: number
  session_id: number
  seat_id: number
  area_id: number
  ticket_no: string
  qr_code: string
  price: number
  seat_info: string
  status: number
  checked_in: number
  checkin_time: string | null
  real_name: string
  id_card: string
}

export interface Refund {
  id: number
  order_id: number
  user_id: number
  refund_no: string
  refund_amount: number
  reason: string
  status: number
  audit_time: string | null
  audit_remark: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginatedResponse<T> {
  list: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    pages: number
  }
}

export const OrderStatus = {
  PENDING: 0,
  PAID: 1,
  CANCELED: 2,
  REFUNDING: 4,
  REFUNDED: 3
} as const

export const SeatStatus = {
  AVAILABLE: 0,
  LOCKED: 1,
  SOLD: 2,
  DISABLED: 3
} as const

export const TicketStatus = {
  VALID: 0,
  USED: 1,
  REFUNDED: 2
} as const

export const RefundStatus = {
  PENDING: 0,
  APPROVED: 1,
  REJECTED: 2
} as const

export const OrderStatusText: Record<number, string> = {
  0: '待支付',
  1: '已支付',
  2: '已取消',
  3: '已退款',
  4: '退款中'
}

export const PayTypeText: Record<number, string> = {
  0: '未支付',
  1: '支付宝',
  2: '微信支付'
}

export const SeatStatusText: Record<number, string> = {
  0: '可选',
  1: '已锁定',
  2: '已售出',
  3: '不可用'
}
