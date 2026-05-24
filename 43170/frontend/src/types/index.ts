export interface User {
  id: number
  username: string
  email: string
  realName: string
  idCard: string
  role: 'renter' | 'owner' | 'admin'
  verified: boolean
  phone: string
  avatar: string
  depositBalance: number
  createdAt: string
  updatedAt: string
}

export interface Equipment {
  id: number
  ownerId: number
  owner?: User
  name: string
  category: string
  brand: string
  model: string
  purchaseDate: string
  status: 'available' | 'rented' | 'maintenance'
  deposit: number
  dailyRent: number
  description: string
  rating: number
  reviewCount: number
  images: EquipmentImage[]
  createdAt: string
  updatedAt: string
}

export interface EquipmentImage {
  id: number
  equipmentId: number
  imageUrl: string
  sortOrder: number
  createdAt: string
}

export interface Order {
  id: number
  equipmentId: number
  equipment?: Equipment
  renterId: number
  renter?: User
  ownerId: number
  owner?: User
  startDate: string
  endDate: string
  totalRent: number
  deposit: number
  status: 'pending' | 'confirmed' | 'rented' | 'completed' | 'cancelled'
  deliveryMethod: string
  deliveryAddress: string
  rejectReason?: string
  createdAt: string
  updatedAt: string
}

export interface Review {
  id: number
  orderId: number
  fromUserId: number
  toUserId: number
  equipmentId: number
  rating: number
  content: string
  createdAt: string
}

export interface Settlement {
  id: number
  orderId: number
  order?: Order
  totalRent: number
  deposit: number
  refundDeposit: number
  finalAmount: number
  damageFee: number
  status: string
  remark: string
  settledAt: string
  createdAt: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T> {
  code: number
  message: string
  data: T[]
  total: number
  page: number
  pageSize: number
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
  expiresAt: number
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  role: 'renter' | 'owner' | 'admin'
  realName?: string
  idCard?: string
  phone?: string
}

export interface CreateEquipmentRequest {
  name: string
  category: string
  brand: string
  model: string
  purchaseDate?: string
  deposit: number
  dailyRent: number
  description?: string
  images?: string[]
}

export interface CreateOrderRequest {
  equipmentId: number
  startDate: string
  endDate: string
  deliveryMethod: string
  deliveryAddress?: string
}

export interface CreateReviewRequest {
  orderId: number
  toUserId: number
  equipmentId: number
  rating: number
  content: string
}

export interface SearchRequest {
  category?: string
  brand?: string
  minPrice?: number
  maxPrice?: number
  status?: string
  startDate?: string
  endDate?: string
  page?: number
  pageSize?: number
  sortBy?: string
  sortOrder?: string
}
