export interface User {
  id: number
  username: string
  role: 'buyer' | 'seller' | 'appraiser'
  email?: string
  phone?: string
  avatar?: string
  credit_score: number
  review_count: number
  real_name?: string
  certified?: boolean
}

export interface Watch {
  id: number
  seller_id: number
  brand: string
  model: string
  reference_no?: string
  year: number
  movement?: string
  case_size_mm?: number
  case_material?: string
  dial_color?: string
  bracelet?: string
  condition?: string
  description?: string
  price: number
  status: string
  authed: boolean
  photos?: WatchPhoto[]
}

export interface WatchPhoto {
  id: number
  watch_id: number
  url: string
}

export interface AuthOrder {
  id: number
  user_id: number
  watch_id: number
  appraiser_id?: number
  status: 'pending' | 'assigned' | 'reported' | 'rejected'
  note?: string
  photos?: AuthPhoto[]
  report?: AuthReport
}

export interface AuthPhoto {
  id: number
  auth_order_id: number
  url: string
}

export interface AuthReport {
  id: number
  auth_order_id: number
  appraiser_id: number
  conclusion: string
  authentic: boolean
  details: string
  estimated_value: number
  pdf_path?: string
}

export interface Trade {
  id: number
  seller_id: number
  buyer_id?: number
  watch_id: number
  start_price: number
  final_price: number
  status: string
  remark?: string
}

export interface TradeBid {
  id: number
  trade_id: number
  buyer_id: number
  price: number
  message?: string
  accepted: boolean
}

export interface FavoriteGroup {
  id: number
  user_id: number
  brand: string
  name: string
}

export interface Favorite {
  id: number
  user_id: number
  watch_id: number
  group_id?: number
}

export interface Review {
  id: number
  trade_id: number
  from_user_id: number
  to_user_id: number
  role: 'buyer' | 'seller'
  rating: number
  content?: string
}

export interface Message {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  read: boolean
}
