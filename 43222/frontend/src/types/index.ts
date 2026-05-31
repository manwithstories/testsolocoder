export interface User {
  id: string
  username: string
  email: string
  nickname: string
  avatar: string
  phone: string
  region: string
  climate_zone: string
  user_type: string
  credit_score: number
  is_verified: boolean
  is_active: boolean
  last_login_at: string | null
  created_at: string
}

export interface Plot {
  id: string
  user_id: string
  name: string
  description: string
  soil_type: string
  sunlight: string
  area: number
  location: string
  grid_config: string
  irrigation_device: string
  sensor_data: string
  created_at: string
  updated_at: string
  planting_records?: PlantingRecord[]
}

export interface Plant {
  id: string
  name: string
  latin_name: string
  category: string
  description: string
  growth_cycle: number
  water_frequency: string
  fertilizer_need: string
  sunlight_need: string
  soil_ph: string
  planting_season: string
  harvest_season: string
  disease_info: string
  pest_info: string
  image_url: string
  sowing_depth: string
  spacing: string
  difficulty: string
  climate_zone: string
  created_at: string
  updated_at: string
}

export interface PlantingRecord {
  id: string
  plot_id: string
  plant_id: string
  user_id: string
  quantity: number
  planting_date: string
  expected_harvest_date: string | null
  actual_harvest_date: string | null
  status: string
  notes: string
  created_at: string
  updated_at: string
  plant?: Plant
  plot?: Plot
  growth_logs?: GrowthLog[]
}

export interface GrowthLog {
  id: string
  planting_record_id: string
  title: string
  description: string
  image_url: string
  log_type: string
  log_date: string
  created_at: string
  updated_at: string
}

export interface Post {
  id: string
  user_id: string
  title: string
  content: string
  image_urls: string
  category: string
  tags: string
  view_count: number
  like_count: number
  is_pinned: boolean
  created_at: string
  updated_at: string
  user?: User
  comments?: Comment[]
}

export interface Comment {
  id: string
  post_id: string
  user_id: string
  content: string
  parent_id: string | null
  like_count: number
  created_at: string
  updated_at: string
  user?: User
}

export interface SeedExchange {
  id: string
  owner_id: string
  title: string
  seed_name: string
  description: string
  image_urls: string
  quantity: number
  exchange_type: string
  want_seeds: string
  location: string
  status: string
  created_at: string
  updated_at: string
  owner?: User
  exchange_offers?: ExchangeOffer[]
}

export interface ExchangeOffer {
  id: string
  seed_exchange_id: string
  offerer_id: string
  offer_seeds: string
  message: string
  status: string
  created_at: string
  updated_at: string
  offerer?: User
}

export interface Product {
  id: string
  seller_id: string
  name: string
  description: string
  category: string
  price: number
  stock: number
  image_urls: string
  specifications: string
  is_active: boolean
  created_at: string
  updated_at: string
  seller?: User
}

export interface Cart {
  id: string
  user_id: string
  product_id: string
  quantity: number
  created_at: string
  updated_at: string
  product?: Product
}

export interface Order {
  id: string
  user_id: string
  order_no: string
  total_amount: number
  status: string
  payment_method: string
  payment_status: string
  shipping_address: string
  shipping_phone: string
  shipping_name: string
  tracking_number: string
  shipping_status: string
  remark: string
  created_at: string
  updated_at: string
  order_items?: OrderItem[]
}

export interface OrderItem {
  id: string
  order_id: string
  product_id: string
  product_name: string
  price: number
  quantity: number
  subtotal: number
  created_at: string
}

export interface CalendarEvent {
  id: string
  user_id: string
  title: string
  event_type: string
  event_date: string
  description: string
  plant_id: string | null
  plot_id: string | null
  is_completed: boolean
  completed_at: string | null
  created_at: string
  updated_at: string
}

export interface DiseaseDiagnosis {
  id: string
  user_id: string
  plant_name: string
  image_url: string
  description: string
  symptoms: string
  diagnosis: string
  severity: string
  treatment: string
  confidence: number
  status: string
  created_at: string
  updated_at: string
}

export interface ApiResponse<T> {
  data?: T
  message?: string
  error?: string
  total?: number
  page?: number
  page_size?: number
}

export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}
