export interface User {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  role: 'normal' | 'admin' | 'guide' | 'researcher'
  avatar: string
  museum_id?: number
  status: number
  created_at: string
}

export interface Museum {
  id: number
  name: string
  description: string
  address: string
  contact: string
  phone: string
  email: string
  logo: string
  open_time: string
  close_time: string
  created_at: string
}

export interface CollectionCategory {
  id: number
  name: string
  parent_id?: number
  sort_order: number
  created_at: string
  parent?: CollectionCategory
  children?: CollectionCategory[]
}

export interface Collection {
  id: number
  name: string
  category_id: number
  code: string
  era: string
  material: string
  size: string
  source: string
  condition: string
  description: string
  image_url: string
  status: 'active' | 'inactive' | 'repair'
  tags: string
  museum_id: number
  view_count: number
  created_at: string
  category?: CollectionCategory
  museum?: Museum
}

export interface Exhibition {
  id: number
  title: string
  description: string
  start_date: string
  end_date: string
  location: string
  hall_number: string
  ticket_price: number
  max_visitors: number
  image_url: string
  status: 'draft' | 'published' | 'closed'
  is_virtual: boolean
  virtual_url: string
  museum_id: number
  view_count: number
  created_at: string
  collections?: Collection[]
  museum?: Museum
}

export interface TimeSlot {
  id: number
  exhibition_id: number
  date: string
  start_time: string
  end_time: string
  max_capacity: number
  booked_count: number
  created_at: string
}

export interface Reservation {
  id: number
  user_id: number
  exhibition_id: number
  time_slot_id: number
  visitor_count: number
  guide_type: string
  total_price: number
  status: 'pending' | 'confirmed' | 'cancelled' | 'completed'
  qr_code: string
  remark: string
  cancelled_at?: string
  created_at: string
  user?: User
  exhibition?: Exhibition
  time_slot?: TimeSlot
}

export interface VisitRecord {
  id: number
  user_id: number
  reservation_id: number
  exhibition_id: number
  check_in_time?: string
  check_out_time?: string
  favorite: boolean
  rating: number
  comment: string
  created_at: string
}

export interface GuideSchedule {
  id: number
  guide_id: number
  date: string
  start_time: string
  end_time: string
  is_available: boolean
  reservation_id?: number
  created_at: string
  guide?: User
}

export interface GuideContent {
  id: number
  collection_id: number
  exhibition_id?: number
  language: string
  content: string
  audio_url?: string
  sort_order: number
  created_at: string
  collection?: Collection
}

export interface ResearchApplication {
  id: number
  user_id: number
  collection_id: number
  purpose: string
  institution: string
  status: 'pending' | 'approved' | 'rejected'
  reviewer_id?: number
  review_comment: string
  reviewed_at?: string
  approved_at?: string
  created_at: string
  user?: User
  collection?: Collection
  reviewer?: User
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  expires_at: string
  user: User
}

export interface RegisterRequest {
  username: string
  email: string
  phone?: string
  password: string
  nickname?: string
}
