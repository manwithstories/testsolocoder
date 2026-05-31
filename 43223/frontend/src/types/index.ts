export type UserRole = 'admin' | 'roaster' | 'user'
export type UserStatus = 'active' | 'disabled' | 'pending'

export interface User {
  id: number
  username: string
  email: string
  phone?: string
  nickname?: string
  avatar?: string
  role: UserRole
  status: UserStatus
  bio?: string
  address?: string
  is_certified: boolean
  certification_id?: number
  created_at: string
  updated_at: string
}

export type RoastLevel = 'light' | 'medium' | 'medium_dark' | 'dark'
export type ProcessMethod = 'washed' | 'natural' | 'honey' | 'anaerobic' | 'wet_hulled'
export type ProductStatus = 'on_sale' | 'offline' | 'draft'

export interface ProductImage {
  id: number
  product_id: number
  url: string
  sort_order: number
  is_cover: boolean
  created_at: string
}

export interface Product {
  id: number
  name: string
  origin: string
  farm?: string
  variety?: string
  altitude?: string
  process_method: ProcessMethod
  roast_level: RoastLevel
  flavor_notes?: string
  cupping_score: number
  description?: string
  price: number
  weight: number
  stock: number
  status: ProductStatus
  roaster_id: number
  roaster?: User
  images?: ProductImage[]
  created_at: string
  updated_at: string
}

export interface RoastingDataPoint {
  id: number
  roasting_record_id: number
  time_elapsed: number
  bean_temp: number
  env_temp: number
  rate_of_rise: number
  created_at: string
}

export interface RoastingRecord {
  id: number
  product_id: number
  product?: Product
  roaster_id: number
  roaster?: User
  batch_number: string
  green_bean_weight: number
  roasted_weight: number
  input_temp: number
  turning_point: number
  turning_time: number
  first_crack_temp: number
  first_crack_time: number
  second_crack_temp: number
  second_crack_time: number
  drop_temp: number
  total_roast_time: number
  notes?: string
  data_points?: RoastingDataPoint[]
  roasted_at: string
  created_at: string
  updated_at: string
}

export type OrderStatus = 'pending' | 'paid' | 'processing' | 'shipped' | 'delivered' | 'cancelled' | 'refunded'
export type PaymentStatus = 'pending' | 'success' | 'failed' | 'refunded'

export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  product?: Product
  product_name: string
  price: number
  quantity: number
  subtotal: number
  created_at: string
}

export interface Order {
  id: number
  order_no: string
  user_id: number
  user?: User
  total_amount: number
  status: OrderStatus
  payment_status: PaymentStatus
  receiver_name: string
  receiver_phone: string
  address: string
  remark?: string
  items?: OrderItem[]
  paid_at?: string
  shipped_at?: string
  delivered_at?: string
  cancelled_at?: string
  created_at: string
  updated_at: string
}

export interface CartItem {
  id: number
  user_id: number
  product_id: number
  product?: Product
  quantity: number
  created_at: string
  updated_at: string
}

export interface CuppingScore {
  id: number
  product_id: number
  product?: Product
  user_id: number
  user?: User
  dry_fragrance: number
  wet_aroma: number
  body: number
  acidity: number
  sweetness: number
  aftertaste: number
  balance: number
  overall_score: number
  notes?: string
  created_at: string
  updated_at: string
}

export type CertStatus = 'pending' | 'approved' | 'rejected'

export interface RoasterCertification {
  id: number
  user_id: number
  user?: User
  cert_name: string
  cert_number: string
  org_name: string
  cert_file?: string
  experience: string
  specialty?: string
  status: CertStatus
  reviewer_id?: number
  reviewer?: User
  review_comment?: string
  reviewed_at?: string
  created_at: string
  updated_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
  errors?: any
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
}

export interface LoginRequest {
  account: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  phone?: string
  nickname?: string
}

export interface CreateProductRequest {
  name: string
  origin: string
  farm?: string
  variety?: string
  altitude?: string
  process_method: ProcessMethod
  roast_level: RoastLevel
  flavor_notes?: string
  cupping_score?: number
  description?: string
  price: number
  weight: number
  stock?: number
  status?: ProductStatus
}

export interface CreateRoastingRecordRequest {
  product_id: number
  batch_number: string
  green_bean_weight?: number
  roasted_weight?: number
  input_temp?: number
  turning_point?: number
  turning_time?: number
  first_crack_temp?: number
  first_crack_time?: number
  second_crack_temp?: number
  second_crack_time?: number
  drop_temp?: number
  total_roast_time?: number
  notes?: string
  data_points?: {
    time_elapsed: number
    bean_temp: number
    env_temp: number
    rate_of_rise: number
  }[]
  roasted_at: string
}

export interface CreateOrderRequest {
  receiver_name: string
  receiver_phone: string
  address: string
  remark?: string
  items: { product_id: number; quantity: number }[]
}

export interface CreateCuppingScoreRequest {
  product_id: number
  dry_fragrance: number
  wet_aroma: number
  body: number
  acidity: number
  sweetness: number
  aftertaste: number
  balance: number
  notes?: string
}
