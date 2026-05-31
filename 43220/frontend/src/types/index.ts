export interface User {
  id: string
  username: string
  email: string
  phone?: string
  role: 'owner' | 'store' | 'keeper' | 'admin'
  avatar_url?: string
  real_name?: string
  status: string
  created_at: string
  store_info?: StoreInfo
  keeper_info?: KeeperInfo
}

export interface StoreInfo {
  id: string
  store_name: string
  address?: string
  license_no?: string
  business_hours?: string
  description?: string
  rating: number
  review_count: number
}

export interface KeeperInfo {
  id: string
  nick_name: string
  experience: number
  specialty?: string
  rating: number
  review_count: number
  certification?: string
  store_id?: string
}

export interface Pet {
  id: string
  owner_id: string
  name: string
  species: string
  breed?: string
  gender: string
  birth_date?: string
  weight: number
  color?: string
  avatar_url?: string
  allergies?: string
  diet_habit?: string
  temperament?: string
  vaccine_records?: VaccineRecord[]
  created_at: string
  updated_at: string
}

export interface VaccineRecord {
  id: string
  pet_id: string
  vaccine_name: string
  vaccinated_at: string
  expire_at: string
  hospital?: string
  proof_url?: string
  is_valid: boolean
}

export interface DewormRecord {
  id: string
  pet_id: string
  deworm_type: string
  dewormed_at: string
  expire_at: string
  medicine?: string
  is_valid: boolean
}

export interface BoardingPackage {
  id: string
  store_id: string
  name: string
  type: 'daycare' | 'boarding'
  description?: string
  price_per_day: number
  capacity: number
  features?: string
  is_available: boolean
  sort_order: number
  created_at: string
}

export interface Reservation {
  id: string
  order_no: string
  owner_id: string
  pet_id: string
  store_id: string
  package_id: string
  package_type: string
  check_in_date: string
  check_out_date: string
  total_days: number
  total_amount: number
  deposit_amount: number
  status: 'pending' | 'confirmed' | 'checked_in' | 'completed' | 'cancelled'
  keeper_id?: string
  remark?: string
  cancel_reason?: string
  cancelled_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
}

export interface DailyRecord {
  id: string
  reservation_id: string
  pet_id: string
  keeper_id: string
  record_date: string
  feed_status?: string
  activity?: string
  health_status?: string
  mood?: string
  photos?: string
  remark?: string
  created_at: string
  updated_at: string
}

export interface Review {
  id: string
  reservation_id: string
  owner_id: string
  store_id: string
  keeper_id?: string
  store_rating: number
  keeper_rating?: number
  content?: string
  images?: string
  reply?: string
  reply_at?: string
  created_at: string
}

export interface Order {
  id: string
  order_no: string
  reservation_id: string
  owner_id: string
  store_id: string
  type: 'prepay' | 'settlement' | 'refund'
  amount: number
  pay_status: 'unpaid' | 'paid' | 'refunded'
  pay_method?: string
  transaction_id?: string
  paid_at?: string
  refund_amount: number
  refund_at?: string
  remark?: string
  amount_hash: string
  created_at: string
}

export interface HealthAlert {
  id: string
  user_id: string
  pet_id: string
  alert_type: string
  title: string
  content: string
  record_id: string
  expire_at: string
  is_read: boolean
  created_at: string
}

export interface PaginationParams {
  page?: number
  page_size?: number
}

export interface PagedResult<T> {
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

export interface LoginData {
  token: string
  user_id: string
  username: string
  role: string
  expires_at: string
}

export interface RegisterData {
  username: string
  email: string
  phone?: string
  password: string
  role: 'owner' | 'store' | 'keeper'
  real_name?: string
}
