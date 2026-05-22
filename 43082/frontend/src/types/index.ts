export interface User {
  id: number
  username: string
  email: string
  phone: string
  nickname: string
  avatar: string
  role: 'buyer' | 'seller' | 'admin'
}

export interface LoginRequest {
  account: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  phone: string
  password: string
  nickname?: string
  role: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface Shop {
  id: number
  userId: number
  name: string
  description: string
  logo: string
  status: 'pending' | 'approved' | 'rejected' | 'suspended'
  rating: number
  createdAt: string
  contactName?: string
  contactPhone?: string
  address?: string
  rejectReason?: string
  productCount?: number
  soldCount?: number
}

export interface Category {
  id: number
  name: string
  icon: string
  parentId?: number
  level: number
  sort: number
  status: string
  children?: Category[]
}

export interface Product {
  id: number
  shopId: number
  categoryId: number
  name: string
  mainImage: string
  price: number
  stock: number
  sales: number
  status: string
  isHot: boolean
  isRecommend: boolean
  createdAt: string
  shopName: string
  categoryName: string
  description?: string
  images?: string[]
  specs?: ProductSpec[]
  skus?: SKU[]
  shop?: Shop
  originalPrice?: number
}

export interface ProductSpec {
  id: number
  name: string
  values: string[]
}

export interface SKU {
  id: number
  specs: Record<string, string>
  price: number
  stock: number
  skuCode: string
}

export interface CartItem {
  id: number
  productId: number
  skuId?: number
  quantity: number
  product: Product
  sku?: SKU
  productName: string
  productImage: string
  specValues?: string[]
  price: number
  stock: number
}

export interface Order {
  id: number
  orderNo: string
  shopId: number
  shopName: string
  totalAmount: number
  shippingFee: number
  status: string
  statusText: string
  receiverName: string
  receiverPhone: string
  receiverAddress: string
  remark: string
  createdAt: string
  paidAt?: string
  shippedAt?: string
  completedAt?: string
  trackingNo?: string
  trackingCompany?: string
  items?: OrderItem[]
  paymentMethod?: string
}

export interface OrderItem {
  id: number
  productId: number
  skuId?: number
  productName: string
  productImage: string
  specs?: Record<string, string>
  price: number
  quantity: number
  subtotal: number
  reviewed: boolean
}

export interface Review {
  id: number
  userId: number
  username: string
  avatar: string
  productId: number
  rating: number
  content: string
  images: string[]
  reply?: string
  replyAt?: string
  createdAt: string
}

export interface Notification {
  id: number
  type: string
  title: string
  content: string
  isRead: boolean
  createdAt: string
}

export interface PaginationParams {
  page?: number
  pageSize?: number
}

export interface ProductQueryParams extends PaginationParams {
  keyword?: string
  categoryId?: number
  shopId?: number
  minPrice?: number
  maxPrice?: number
  sortBy?: string
  sortOrder?: string
  isHot?: boolean
  status?: string
}

export interface PaginatedResponse<T> {
  data: T
  pagination: {
    total: number
    page: number
    pageSize: number
    pages: number
  }
}

export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

export interface DailySales {
  date: string
  amount: number
  orders: number
}

export interface ProductSales {
  id: number
  name: string
  sales: number
  amount: number
}

export interface ShopStatistics {
  totalOrders: number
  totalSales: number
  totalProducts: number
  newOrders: number
  dailySales: DailySales[]
  topProducts: ProductSales[]
}

export interface AdminStatistics {
  totalUsers: number
  totalShops: number
  totalProducts: number
  totalOrders: number
  totalSales: number
  pendingShops: number
  openDisputes: number
  dailySales: DailySales[]
}

export interface Dispute {
  id: number
  orderId: number
  orderNo: string
  userId: number
  userName: string
  type: string
  description: string
  status: string
  result?: string
  createdAt: string
}

export interface Favorite {
  id: number
  targetId: number
  targetType: 'product' | 'shop'
  product?: Product
  shop?: Shop
  createdAt: string
}
