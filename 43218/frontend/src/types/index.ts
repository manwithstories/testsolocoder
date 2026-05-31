export interface User {
  id: number
  username: string
  email?: string
  phone?: string
  nickname?: string
  avatar?: string
  role: 'admin' | 'seller' | 'buyer' | 'technician'
  realName?: string
  idCard?: string
  isAuthenticated: boolean
  creditScore: number
  walletBalance: number
  status: number
  createdAt: string
}

export interface LoginResponse {
  user: User
  accessToken: string
  refreshToken: string
  tokenType: string
  expiresIn: number
}

export interface Product {
  id: number
  sellerId: number
  title: string
  category: string
  brand: string
  model: string
  condition: string
  price: number
  originalPrice?: number
  description: string
  warrantyDays: number
  images: string
  status: number
  viewCount: number
  favoriteCount: number
  soldCount: number
  rejectReason?: string
  createdAt: string
  updatedAt: string
  seller?: User
}

export interface ProductListParams {
  page?: number
  pageSize?: number
  category?: string
  condition?: string
  minPrice?: number
  maxPrice?: number
  keyword?: string
  sortBy?: string
  sellerId?: number
}

export interface PaginationResponse<T> {
  code: number
  message: string
  data: T[]
  pagination: {
    page: number
    pageSize: number
    total: number
  }
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface Order {
  id: number
  orderNo: string
  buyerId: number
  sellerId: number
  productId: number
  productTitle: string
  productImage?: string
  originalPrice: number
  finalPrice: number
  negotiated: boolean
  status: number
  paymentMethod?: string
  paidAt?: string
  shippedAt?: string
  deliveredAt?: string
  completedAt?: string
  cancelledAt?: string
  refundedAt?: string
  trackingNo?: string
  trackingCompany?: string
  receiverName?: string
  receiverPhone?: string
  receiverAddress?: string
  warrantyDays: number
  warrantyUntil?: string
  commission: number
  remark?: string
  createdAt: string
  updatedAt: string
  buyer?: User
  seller?: User
  product?: Product
  negotiations?: Negotiation[]
}

export interface Negotiation {
  id: number
  orderId: number
  buyerId: number
  sellerId: number
  offeredPrice: number
  counterPrice?: number
  status: number
  buyerMessage?: string
  sellerMessage?: string
  expireAt?: string
  createdAt: string
  updatedAt: string
}

export interface RepairService {
  id: number
  technicianId: number
  serviceType: string
  title: string
  description: string
  price: number
  minPrice?: number
  maxPrice?: number
  estimatedDays: number
  images?: string
  status: number
  orderCount: number
  rating: number
  createdAt: string
  updatedAt: string
  technician?: User
}

export interface RepairOrder {
  id: number
  orderNo: string
  buyerId: number
  technicianId: number
  serviceId: number
  deviceType: string
  deviceBrand: string
  deviceModel: string
  faultDescription: string
  contactName: string
  contactPhone: string
  address?: string
  servicePrice: number
  finalPrice: number
  status: number
  paymentMethod?: string
  paidAt?: string
  acceptedAt?: string
  completedAt?: string
  pickedUpAt?: string
  warrantyDays: number
  warrantyUntil?: string
  remark?: string
  createdAt: string
  updatedAt: string
  buyer?: User
  technician?: User
  service?: RepairService
}

export interface Review {
  id: number
  orderId?: number
  repairOrderId?: number
  reviewerId: number
  revieweeId: number
  reviewType: string
  rating: number
  content?: string
  images?: string
  qualityScore?: number
  serviceScore?: number
  anonymous: boolean
  status: number
  createdAt: string
  updatedAt: string
  reviewer?: User
  reviewee?: User
}

export interface Report {
  id: number
  reporterId: number
  targetType: string
  targetId: number
  reason: string
  description?: string
  images?: string
  status: number
  handleResult?: string
  handledBy?: number
  handledAt?: string
  createdAt: string
  updatedAt: string
  reporter?: User
}

export interface Warranty {
  id: number
  orderId?: number
  repairOrderId?: number
  userId: number
  type: string
  description: string
  images?: string
  status: number
  handleResult?: string
  handledBy?: number
  handledAt?: string
  createdAt: string
  updatedAt: string
}

export interface Notification {
  id: number
  userId: number
  type: string
  title: string
  content?: string
  orderNo?: string
  isRead: boolean
  extraData?: string
  createdAt: string
}

export interface Message {
  id: number
  senderId: number
  receiverId: number
  content: string
  type: string
  isRead: boolean
  createdAt: string
  sender?: User
  receiver?: User
}

export interface WalletLog {
  id: number
  userId: number
  type: string
  amount: number
  balance: number
  orderNo?: string
  description?: string
  createdAt: string
}

export interface Transaction {
  id: number
  userId: number
  type: string
  amount: number
  orderNo?: string
  paymentMethod?: string
  transactionNo?: string
  status: number
  remark?: string
  createdAt: string
  updatedAt: string
}

export interface Favorite {
  id: number
  userId: number
  productId: number
  createdAt: string
  product?: Product
}

export interface DashboardStats {
  totalUsers: number
  sellers: number
  buyers: number
  technicians: number
  totalProducts: number
  pendingProducts: number
  totalOrders: number
  totalRepairOrders: number
  totalAmount: number
  totalCommission: number
}

export interface OrderStats {
  pending: number
  paid: number
  shipped: number
  completed: number
}

export interface UserStats {
  sellOrders: number
  buyOrders: number
  products: number
  reviews: number
}

export const OrderStatus = {
  PENDING: 1,
  PAID: 2,
  SHIPPED: 3,
  DELIVERED: 4,
  COMPLETED: 5,
  CANCELLED: 6,
  REFUNDING: 7,
  REFUNDED: 8,
  NEGOTIATING: 9
}

export const OrderStatusText: Record<number, string> = {
  1: '待支付',
  2: '已支付',
  3: '已发货',
  4: '已送达',
  5: '已完成',
  6: '已取消',
  7: '退款中',
  8: '已退款',
  9: '议价中'
}

export const RepairStatus = {
  PENDING: 1,
  ACCEPTED: 2,
  REPAIRING: 3,
  COMPLETED: 4,
  PICKED_UP: 5,
  CANCELLED: 6,
  REFUNDING: 7,
  REFUNDED: 8
}

export const RepairStatusText: Record<number, string> = {
  1: '待接单',
  2: '已接单',
  3: '维修中',
  4: '待取件',
  5: '已取件',
  6: '已取消',
  7: '退款中',
  8: '已退款'
}

export const ProductStatus = {
  PENDING: 0,
  APPROVED: 1,
  REJECTED: 2,
  ON_SALE: 3,
  SOLD_OUT: 4,
  OFF_SHELF: 5
}

export const ProductStatusText: Record<number, string> = {
  0: '待审核',
  1: '审核通过',
  2: '审核拒绝',
  3: '在售',
  4: '已售出',
  5: '已下架'
}

export const ProductCategories = [
  '手机', '电脑', '相机', '耳机', '平板', '智能手表', '游戏机', '其他数码'
]

export const ProductConditions = [
  '全新', '95新', '9成新', '8成新', '7成新及以下'
]

export const ServiceTypes = [
  '屏幕更换', '电池更换', '主板维修', '外壳更换', '摄像头维修',
  '充电接口维修', '扬声器维修', '系统维护', '其他'
]
