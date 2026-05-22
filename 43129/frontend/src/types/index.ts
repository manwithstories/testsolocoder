export interface User {
  id: number
  phone: string
  nickname?: string
  avatar?: string
  role: string
  status: number
  last_login_at?: string
  created_at: string
  updated_at: string
}

export interface Customer {
  id: number
  user_id: number
  user?: User
  name?: string
  gender?: string
  age?: number
  skin_type?: string
  hair_preference?: string
  allergy_history?: string
  notes?: string
  member_level: number
  points: number
  total_spent: number
  visit_count: number
  last_visit_at?: string
  created_at: string
  updated_at: string
}

export interface Technician {
  id: number
  user_id: number
  user?: User
  name: string
  title?: string
  avatar?: string
  specialties?: string
  description?: string
  rating: number
  review_count: number
  work_start_time: string
  work_end_time: string
  work_days: string
  status: number
  created_at: string
  updated_at: string
}

export interface Service {
  id: number
  name: string
  category: string
  description?: string
  price: number
  duration: number
  required_skill?: string
  products?: string
  is_package: boolean
  package_count: number
  dynamic_pricing: boolean
  weekend_price: number
  holiday_price: number
  status: number
  created_at: string
  updated_at: string
}

export interface Appointment {
  id: number
  customer_id: number
  customer?: Customer
  technician_id: number
  technician?: Technician
  service_id: number
  service?: Service
  package_id?: number
  appointment_date: string
  start_time: string
  end_time: string
  status: string
  remark?: string
  cancel_reason?: string
  points_deducted: number
  created_at: string
  updated_at: string
}

export interface Payment {
  id: number
  appointment_id: number
  appointment?: Appointment
  customer_id: number
  amount: number
  pay_method: string
  points_used: number
  card_id?: number
  status: string
  transaction_no: string
  created_at: string
  updated_at: string
}

export interface MemberCard {
  id: number
  customer_id: number
  customer?: Customer
  card_no: string
  card_type?: string
  balance: number
  discount: number
  status: number
  created_at: string
  updated_at: string
}

export interface CustomerPackage {
  id: number
  customer_id: number
  customer?: Customer
  service_id: number
  service?: Service
  total_count: number
  used_count: number
  purchase_date: string
  expire_date: string
  status: number
  created_at: string
  updated_at: string
}

export interface Product {
  id: number
  name: string
  category?: string
  description?: string
  unit: string
  stock: number
  threshold: number
  price: number
  retail_price: number
  supplier?: string
  status: number
  created_at: string
  updated_at: string
}

export interface ProductRecord {
  id: number
  product_id: number
  product?: Product
  change_type: string
  quantity: number
  before_stock: number
  after_stock: number
  appointment_id?: number
  operator_id?: number
  remark?: string
  created_at: string
}

export interface ProductSale {
  id: number
  customer_id: number
  customer?: Customer
  product_id: number
  product?: Product
  quantity: number
  unit_price: number
  total_price: number
  pay_method: string
  operator_id?: number
  created_at: string
}

export interface Notification {
  id: number
  user_id: number
  user?: User
  type: string
  title: string
  content: string
  is_read: boolean
  created_at: string
}

export interface AuditLog {
  id: number
  user_id?: number
  user?: User
  action: string
  module: string
  detail: string
  ip: string
  created_at: string
}

export interface Review {
  id: number
  appointment_id: number
  appointment?: Appointment
  customer_id: number
  customer?: Customer
  technician_id: number
  technician?: Technician
  service_id: number
  service?: Service
  rating: number
  content?: string
  created_at: string
  updated_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageParams {
  page?: number
  page_size?: number
  [key: string]: any
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface AvailableSlot {
  date: string
  slots: string[]
}

export interface RevenueReport {
  total_revenue: number
  date_range: string
  daily_data: DailyRevenue[]
}

export interface DailyRevenue {
  date: string
  revenue: number
  count: number
}

export interface TechnicianPerformance {
  technician_id: number
  technician_name: string
  service_count: number
  revenue: number
}

export interface ServiceRank {
  service_id: number
  service_name: string
  count: number
  revenue: number
}

export interface ReportData {
  revenue_report: RevenueReport
  technician_performances: TechnicianPerformance[]
  service_ranks: ServiceRank[]
}
