export interface User {
  id: number
  username: string
  phone: string
  email: string
  real_name: string
  avatar: string
  role: string
  status: string
  balance: number
  address: string
  longitude: number
  latitude: number
}

export interface TechnicianProfile {
  id: number
  user_id: number
  certificate_image: string
  certificate_no: string
  specialty: string
  experience_years: number
  rating: number
  total_orders: number
  completed_orders: number
  active_orders: number
  max_active_orders: number
  is_verified: boolean
  verify_status: string
  verify_remark: string
  service_radius: number
  bank_account: string
  bank_name: string
}

export interface Category {
  id: number
  name: string
  code: string
  icon: string
  sort: number
  status: boolean
}

export interface ServiceItem {
  id: number
  category_id: number
  category?: Category
  name: string
  description: string
  min_price: number
  max_price: number
  estimated_time: number
  image: string
  sort: number
  status: boolean
}

export interface Order {
  id: number
  order_no: string
  customer_id: number
  customer?: User
  technician_id?: number
  technician?: User
  service_item_id: number
  service_item?: ServiceItem
  title: string
  description: string
  images: string
  address: string
  longitude: number
  latitude: number
  contact_name: string
  contact_phone: string
  appointment_time?: string
  quoted_price: number
  final_price: number
  status: string
  urgent_level: number
  cancel_reason: string
  refund_reason: string
  refund_amount: number
  assigned_at?: string
  accepted_at?: string
  arrived_at?: string
  started_at?: string
  completed_at?: string
  cancelled_at?: string
  refunded_at?: string
  created_at: string
  updated_at: string
}

export interface OrderLog {
  id: number
  order_id: number
  user_id: number
  action: string
  content: string
  longitude: number
  latitude: number
  created_at: string
}

export interface Review {
  id: number
  order_id: number
  customer_id: number
  customer?: User
  technician_id: number
  technician?: User
  rating: number
  content: string
  images: string
  reply: string
  replied_at?: string
  is_intervened: boolean
  intervene_note: string
  created_at: string
}

export interface Part {
  id: number
  name: string
  code: string
  category: string
  description: string
  price: number
  stock: number
  min_stock: number
  image: string
  status: boolean
}

export interface PartRequest {
  id: number
  request_no: string
  technician_id: number
  technician?: User
  status: string
  total_amount: number
  remark: string
  approved_at?: string
  shipped_at?: string
  received_at?: string
  created_at: string
}

export interface PartRequestItem {
  id: number
  part_request_id: number
  part_id: number
  part?: Part
  quantity: number
  price: number
}

export interface PartUsage {
  id: number
  order_id: number
  part_id: number
  part?: Part
  technician_id: number
  quantity: number
  price: number
  note: string
  created_at: string
}

export interface WithdrawRequest {
  id: number
  request_no: string
  technician_id: number
  technician?: User
  amount: number
  bank_account: string
  bank_name: string
  status: string
  remark: string
  approved_at?: string
  transferred_at?: string
  created_at: string
}

export interface Transaction {
  id: number
  transaction_no: string
  user_id: number
  type: string
  amount: number
  balance_after: number
  order_id?: number
  description: string
  created_at: string
}

export interface MonthlyReport {
  id: number
  month: string
  total_orders: number
  total_revenue: number
  platform_income: number
  technician_pay: number
  total_withdraw: number
  new_technicians: number
  new_customers: number
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface LoginData {
  token: string
  user_id: number
  username: string
  role: string
  status: string
}

export interface PaginationData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}
