export interface User {
  id: number
  username: string
  email: string
  phone: string
  nickname: string
  avatar: string
  balance: number
  role: string
  status: number
  created_at: string
}

export interface Category {
  id: number
  name: string
  parent_id: number
  sort: number
  status: number
}

export interface AuctionImage {
  id: number
  auction_item_id: number
  url: string
  sort: number
  is_main: number
  created_at: string
}

export interface AuctionItem {
  id: number
  title: string
  description: string
  category_id: number
  seller_id: number
  start_price: number
  reserve_price: number
  current_price: number
  view_count: number
  bid_count: number
  status: number
  location: string
  condition: string
  created_at: string
  category?: Category
  seller?: User
  images?: AuctionImage[]
}

export interface AuctionSession {
  id: number
  name: string
  description: string
  start_time: string
  end_time: string
  min_increment: number
  extend_time: number
  status: number
  created_at: string
  auction_items?: AuctionItemSession[]
}

export interface AuctionItemSession {
  id: number
  session_id: number
  auction_item_id: number
  sort: number
  start_time: string
  end_time: string
  status: number
  auction_item?: AuctionItem
}

export interface Bid {
  id: number
  auction_item_id: number
  session_id: number
  user_id: number
  amount: number
  max_auto_bid: number
  is_auto_bid: number
  is_winning: number
  created_at: string
  auction_item?: AuctionItem
  user?: User
}

export interface AutoBid {
  id: number
  auction_item_id: number
  user_id: number
  max_price: number
  current_bid: number
  status: number
  created_at: string
  auction_item?: AuctionItem
}

export interface Order {
  id: number
  order_no: string
  auction_item_id: number
  buyer_id: number
  seller_id: number
  price: number
  status: number
  payment_time?: string
  shipping_info: string
  tracking_no: string
  created_at: string
  auction_item?: AuctionItem
  buyer?: User
  seller?: User
  payments?: Payment[]
}

export interface Payment {
  id: number
  order_id: number
  payment_no: string
  amount: number
  method: string
  status: number
  transaction_id: string
  created_at: string
  paid_at?: string
}

export interface Review {
  id: number
  order_id: number
  reviewer_id: number
  reviewee_id: number
  rating: number
  content: string
  created_at: string
  reviewer?: User
  reviewee?: User
}

export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  related_id: number
  is_read: number
  created_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user_info: User
}

export interface RegisterRequest {
  username: string
  email: string
  phone: string
  password: string
  nickname?: string
}

export interface Statistics {
  total_auctions: number
  total_bids: number
  total_orders: number
  total_amount: number
  success_rate: number
  active_users: number
  new_users: number
  average_bid_amount: number
}
