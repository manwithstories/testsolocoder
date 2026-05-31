export interface User {
  id: number
  username: string
  email: string
  phone?: string
  role: string
  avatar?: string
  reputation: number
  address?: string
}

export interface Beehive {
  id: number
  user_id: number
  name: string
  code: string
  latitude: number
  longitude: number
  region?: string
  bee_species?: string
  group_name?: string
  status: string
  health_status: string
  queen_status: string
  worker_count: number
  last_inspection?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface HealthRecord {
  id: number
  beehive_id: number
  record_date: string
  queen_status?: string
  worker_count?: number
  has_disease: boolean
  disease_type?: string
  disease_severity?: string
  treatment?: string
  season: string
  recommendations?: string
  notes?: string
  created_at: string
  beehive?: Beehive
}

export interface Harvest {
  id: number
  user_id: number
  beehive_id: number
  harvest_date: string
  honey_type: string
  quantity: number
  unit: string
  quality: string
  batch_code: string
  notes?: string
  created_at: string
  beehive?: Beehive
}

export interface Inventory {
  id: number
  user_id: number
  harvest_id: number
  honey_type: string
  batch_code: string
  quantity: number
  unit: string
  expiry_date: string
  inspection_report?: string
  grade: string
  status: string
  threshold: number
  price?: number
  created_at: string
  updated_at: string
  harvest?: Harvest
}

export interface Product {
  id: number
  user_id: number
  inventory_id: number
  title: string
  description?: string
  honey_type: string
  batch_code?: string
  price: number
  stock: number
  unit: string
  images?: string[]
  grade?: string
  status: string
  view_count: number
  created_at: string
  updated_at: string
  user?: User
  inventory?: Inventory
}

export interface Order {
  id: number
  order_no: string
  buyer_id: number
  seller_id: number
  product_id: number
  quantity: number
  unit_price: number
  total_amount: number
  status: string
  payment_status: string
  payment_time?: string
  shipping_address: string
  tracking_number?: string
  tracking_status?: string
  buyer_rating?: number
  seller_rating?: number
  buyer_comment?: string
  seller_comment?: string
  created_at: string
  updated_at: string
  buyer?: User
  seller?: User
  product?: Product
}

export interface Inspection {
  id: number
  user_id: number
  inspector_id?: number
  inventory_id: number
  batch_code: string
  appointment_date: string
  inspection_date?: string
  status: string
  report_url?: string
  result?: string
  grade?: string
  notes?: string
  created_at: string
  updated_at: string
  user?: User
  inspector?: User
  inventory?: Inventory
}

export interface Post {
  id: number
  user_id: number
  title: string
  content: string
  category: string
  tags?: string[]
  images?: string[]
  view_count: number
  like_count: number
  comment_count: number
  created_at: string
  updated_at: string
  user?: User
}

export interface Comment {
  id: number
  post_id: number
  user_id: number
  content: string
  parent_id?: number
  like_count: number
  created_at: string
  user?: User
}

export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content?: string
  related_id?: number
  is_read: boolean
  created_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  total?: number
}

export interface PageParams {
  page?: number
  page_size?: number
  [key: string]: any
}
