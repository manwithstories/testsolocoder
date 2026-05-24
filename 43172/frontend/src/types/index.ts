export interface User {
  id: number
  username: string
  email: string
  phone?: string
  real_name?: string
  avatar?: string
  role: UserRole
  status: UserStatus
  credit_score: number
  address?: string
  created_at: string
  updated_at: string
  authenticator_profile?: AuthenticatorProfile
  seller_profile?: SellerProfile
  buyer_profile?: BuyerProfile
}

export type UserRole = 'buyer' | 'seller' | 'authenticator' | 'admin'
export type UserStatus = 'active' | 'inactive' | 'pending' | 'banned'

export interface AuthenticatorProfile {
  id: number
  user_id: number
  license_number: string
  license_file: string
  certifications?: string
  specialties?: string
  status: AuthenticatorStatus
  rating: number
  completed_count: number
  rejection_reason?: string
  verified_at?: string
  created_at?: string
  user?: {
    id: number
    username: string
    avatar?: string
    email?: string
  }
}

export type AuthenticatorStatus = 'pending' | 'approved' | 'rejected'

export interface SellerProfile {
  id: number
  user_id: number
  store_name?: string
  store_logo?: string
  description?: string
  business_license?: string
  total_sales: number
}

export interface BuyerProfile {
  id: number
  user_id: number
  preferred_brands?: string
  total_purchases: number
  total_spent: number
}

export interface Brand {
  id: number
  name: string
  name_cn?: string
  logo?: string
  country?: string
  description?: string
  category: ProductCategory
  popularity: number
  created_at: string
}

export interface Product {
  id: number
  seller_id: number
  title: string
  description: string
  category: ProductCategory
  brand_id?: number
  brand_name?: string
  original_price?: number
  price: number
  condition?: string
  color?: string
  size?: string
  material?: string
  stock: number
  status: ProductStatus
  views: number
  favorites: number
  is_authenticated: boolean
  created_at: string
  updated_at: string
  seller?: User
  brand?: Brand
  images?: ProductImage[]
}

export type ProductCategory = 'bag' | 'jewelry' | 'watch' | 'clothing' | 'shoes' | 'other'
export type ProductStatus = 'draft' | 'on_sale' | 'sold' | 'removed'

export interface ProductImage {
  id: number
  product_id: number
  image_url: string
  image_type: string
  is_primary: boolean
  sort_order: number
}

export interface Order {
  id: number
  order_number: string
  buyer_id: number
  seller_id: number
  product_id: number
  price: number
  status: OrderStatus
  payment_status: PaymentStatus
  payment_method?: string
  payment_time?: string
  shipping_address?: string
  tracking_number?: string
  shipped_at?: string
  delivered_at?: string
  completed_at?: string
  cancelled_at?: string
  cancel_reason?: string
  need_auth: boolean
  remark?: string
  created_at: string
  updated_at: string
  buyer?: User
  seller?: User
  product?: Product
  authentication?: Authentication
  review?: Review
}

export type OrderStatus = 'pending' | 'paid' | 'shipped' | 'delivered' | 'completed' | 'cancelled' | 'refunded'
export type PaymentStatus = 'pending' | 'success' | 'failed' | 'refunded'

export interface Authentication {
  id: number
  order_id: number
  product_id: number
  buyer_id: number
  authenticator_id?: number
  status: AuthenticationStatus
  result?: AuthenticationResult
  report_file?: string
  report_content?: string
  authenticator_notes?: string
  rejection_reason?: string
  accepted_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
  order?: Order
  product?: Product
  buyer?: User
  authenticator?: User
}

export type AuthenticationStatus = 'pending' | 'accepted' | 'completed' | 'rejected' | 'cancelled'
export type AuthenticationResult = 'genuine' | 'counterfeit' | 'inconclusive'

export interface Review {
  id: number
  order_id: number
  reviewer_id: number
  reviewee_id: number
  rating: ReviewRating
  content?: string
  images?: string
  is_anonymous: boolean
  created_at: string
  order?: Order
  reviewer?: User
  reviewee?: User
}

export type ReviewRating = 1 | 2 | 3 | 4 | 5

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
  total?: number
  page?: number
  page_size?: number
}

export interface LoginResponse {
  token: string
  user: User
  expire: string
}

export interface PaginationParams {
  page?: number
  page_size?: number
}

export interface DashboardStats {
  total_orders: number
  total_amount: number
  total_users: number
  total_products: number
  transaction_trend: TransactionTrend[]
  brand_rankings: BrandRanking[]
  auth_stats: AuthStats
  recent_orders: Order[]
}

export interface TransactionTrend {
  date: string
  order_count: number
  total_amount: number
}

export interface BrandRanking {
  brand_name: string
  product_count: number
  order_count: number
}

export interface AuthStats {
  total: number
  completed: number
  passed: number
  pass_rate: number
}

export const CATEGORY_OPTIONS = [
  { label: '包包', value: 'bag' },
  { label: '首饰', value: 'jewelry' },
  { label: '手表', value: 'watch' },
  { label: '服装', value: 'clothing' },
  { label: '鞋履', value: 'shoes' },
  { label: '其他', value: 'other' }
]

export const ORDER_STATUS_OPTIONS = [
  { label: '待支付', value: 'pending', type: 'warning' },
  { label: '已支付', value: 'paid', type: 'primary' },
  { label: '已发货', value: 'shipped', type: 'info' },
  { label: '已送达', value: 'delivered', type: 'success' },
  { label: '已完成', value: 'completed', type: 'success' },
  { label: '已取消', value: 'cancelled', type: 'info' },
  { label: '已退款', value: 'refunded', type: 'danger' }
]

export const AUTH_STATUS_OPTIONS = [
  { label: '待接单', value: 'pending', type: 'warning' },
  { label: '鉴定中', value: 'accepted', type: 'primary' },
  { label: '已完成', value: 'completed', type: 'success' },
  { label: '已拒绝', value: 'rejected', type: 'danger' },
  { label: '已取消', value: 'cancelled', type: 'info' }
]

export const AUTH_RESULT_OPTIONS = [
  { label: '正品', value: 'genuine', type: 'success' },
  { label: '赝品', value: 'counterfeit', type: 'danger' },
  { label: '无法鉴定', value: 'inconclusive', type: 'warning' }
]
