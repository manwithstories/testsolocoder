export interface User {
  id: number
  phone: string
  nickname: string
  avatar: string
  role: UserRole
  status: UserStatus
  balance: number
  rating: number
  order_count: number
  created_at: string
}

export type UserRole = 'publisher' | 'courier' | 'admin'
export type UserStatus = 'active' | 'frozen' | 'verified'

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_at: string
  user: UserInfo
}

export interface UserInfo {
  id: number
  phone: string
  nickname: string
  avatar: string
  role: UserRole
  status: UserStatus
}

export interface Task {
  id: number
  publisher_id: number
  publisher?: User
  courier_id?: number
  courier?: User
  type: TaskType
  title: string
  description: string
  start_addr: string
  start_lat: number
  start_lng: number
  end_addr: string
  end_lat: number
  end_lng: number
  publish_time: string
  deadline: string
  reward: number
  status: TaskStatus
  images?: TaskImage[]
  order?: Order
  distance?: number
  created_at: string
  updated_at: string
}

export type TaskType = 'buy' | 'pickup' | 'deliver' | 'queue' | 'errand'
export type TaskStatus = 'pending' | 'accepted' | 'in_progress' | 'completed' | 'cancelled' | 'timeout'

export interface TaskImage {
  id: number
  task_id: number
  image_url: string
  sort_order: number
  created_at: string
}

export interface Order {
  id: number
  task_id: number
  task?: Task
  publisher_id: number
  publisher?: User
  courier_id: number
  courier?: User
  status: OrderStatus
  reward: number
  service_fee: number
  actual_payment: number
  start_time?: string
  end_time?: string
  tracks?: OrderTrack[]
  proof_images?: OrderProof[]
  created_at: string
  updated_at: string
}

export type OrderStatus = 'pending' | 'accepted' | 'in_progress' | 'delivered' | 'completed' | 'cancelled'

export interface OrderTrack {
  id: number
  order_id: number
  latitude: number
  longitude: number
  address: string
  message: string
  event_type: string
  created_at: string
}

export interface OrderProof {
  id: number
  order_id: number
  image_url: string
  proof_type: string
  created_at: string
}

export interface Transaction {
  id: number
  order_id?: number
  user_id: number
  type: TransactionType
  amount: number
  balance_before: number
  balance_after: number
  status: TransactionStatus
  description: string
  payment_method: string
  transaction_no: string
  failure_reason?: string
  completed_at?: string
  created_at: string
}

export type TransactionType = 'deposit' | 'withdraw' | 'payment' | 'refund' | 'settlement' | 'service_fee'
export type TransactionStatus = 'pending' | 'completed' | 'failed' | 'cancelled'

export interface Review {
  id: number
  order_id: number
  order?: Order
  reviewer_id: number
  reviewer?: User
  reviewee_id: number
  reviewee?: User
  review_type: ReviewType
  rating: number
  content: string
  tags: string
  created_at: string
}

export type ReviewType = 'courier' | 'publisher'

export interface CourierProfile {
  id: number
  user_id: number
  user?: User
  status: string
  level: number
  total_orders: number
  completed_orders: number
  cancelled_orders: number
  rating: number
  current_task_id?: number
  created_at: string
  updated_at: string
}

export interface ChatMessage {
  id: number
  order_id: number
  sender_id: number
  sender?: User
  content: string
  msg_type: string
  is_read: boolean
  created_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PaginationParams {
  page?: number
  page_size?: number
}

export interface PaginatedData<T> {
  total: number
  page: number
  page_size: number
  items: T[]
}
