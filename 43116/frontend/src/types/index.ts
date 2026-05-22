export interface User {
  id: number
  username: string
  email: string
  phone: string
  real_name: string
  id_card: string
  license_number: string
  license_image: string
  id_card_front: string
  id_card_back: string
  auth_status: string
  status: string
  role_id: number
  role?: Role
  avatar: string
  last_login_at: string
  created_at: string
}

export interface Role {
  id: number
  name: string
  description: string
}

export interface LoginRequest {
  account: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  email: string
  phone: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserInfo
}

export interface UserInfo {
  id: number
  username: string
  email: string
  phone: string
  real_name: string
  auth_status: string
  status: string
  role_id: number
  role_name: string
  avatar: string
}

export interface Car {
  id: number
  brand: string
  model: string
  year: number
  seats: number
  transmission: string
  fuel_type: string
  daily_rent: number
  deposit: number
  status: string
  license_plate: string
  color: string
  mileage: number
  rating: number
  review_count: number
  store_id: number
  store?: Store
  images?: CarImage[]
  features: string
  description: string
  created_at: string
}

export interface CarImage {
  id: number
  car_id: number
  url: string
  is_cover: boolean
  sort_order: number
}

export interface City {
  id: number
  name: string
  code: string
  province: string
  stores?: Store[]
}

export interface Store {
  id: number
  name: string
  city_id: number
  city?: City
  address: string
  phone: string
  business_hours: string
  latitude: number
  longitude: number
  is_active: boolean
}

export interface Booking {
  id: number
  booking_no: string
  user_id: number
  user?: User
  car_id: number
  car?: Car
  pickup_store_id: number
  pickup_store?: Store
  return_store_id: number
  return_store?: Store
  pickup_time: string
  return_time: string
  actual_pickup_time?: string
  actual_return_time?: string
  total_days: number
  base_price: number
  total_price: number
  discount: number
  final_price: number
  deposit: number
  promo_code_id?: number
  promo_code?: PromoCode
  status: string
  remark: string
  cancel_reason: string
  cancelled_at?: string
  created_at: string
}

export interface CreateBookingRequest {
  car_id: number
  pickup_store_id: number
  return_store_id: number
  pickup_time: string
  return_time: string
  promo_code?: string
}

export interface PriceCalculation {
  base_price: number
  total_price: number
  discount: number
  final_price: number
  total_days: number
  price_details: DayPrice[]
}

export interface DayPrice {
  date: string
  base_price: number
  multiplier: number
  day_price: number
}

export interface PromoCode {
  id: number
  code: string
  name: string
  type: string
  value: number
  min_amount: number
  max_discount: number
  usage_limit: number
  used_count: number
  start_date: string
  end_date: string
  is_active: boolean
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  user?: User
  car_id: number
  car?: Car
  booking_id: number
  booking?: Booking
  total_amount: number
  discount: number
  final_amount: number
  paid_amount: number
  payment_status: string
  payment_method: string
  paid_at?: string
  status: string
  refund_reason: string
  refunded_at?: string
  created_at: string
}

export interface Review {
  id: number
  user_id: number
  user?: User
  car_id: number
  car?: Car
  booking_id: number
  booking?: Booking
  rating: number
  content: string
  images: string
  is_anonymous: boolean
  likes: number
  is_hidden: boolean
  created_at: string
}

export interface CreateReviewRequest {
  booking_id: number
  rating: number
  content?: string
  images?: string
  is_anonymous?: boolean
}

export interface MaintenancePlan {
  id: number
  car_id: number
  car?: Car
  title: string
  description: string
  start_date: string
  end_date: string
  actual_start?: string
  actual_end?: string
  cost: number
  status: string
  notes: string
  created_by: number
  created_at: string
}

export interface Message {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  related_id?: number
  is_read: boolean
  created_at: string
}

export interface DashboardStats {
  monthly_revenue: number
  monthly_orders: number
  car_utilization: number
  total_cars: number
  active_bookings: number
  popular_cars: CarStat[]
  revenue_trend: RevenuePoint[]
  status_breakdown: StatusCount[]
}

export interface CarStat {
  id: number
  brand: string
  model: string
  rating: number
  review_count: number
  booking_count: number
}

export interface RevenuePoint {
  date: string
  revenue: number
}

export interface StatusCount {
  status: string
  count: number
}

export interface PageResult<T> {
  total: number
  page: number
  page_size: number
  items: T[]
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PricingRule {
  id: number
  name: string
  rule_type: string
  start_date?: string
  end_date?: string
  weekdays: string
  multiplier: number
  car_model: string
  min_days: number
  max_days: number
  is_active: boolean
  priority: number
}
