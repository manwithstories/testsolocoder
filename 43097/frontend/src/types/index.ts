export enum RoomStatus {
  AVAILABLE = 'available',
  OCCUPIED = 'occupied',
  MAINTENANCE = 'maintenance',
  CLEANING = 'cleaning',
  RESERVED = 'reserved'
}

export enum BookingStatus {
  PENDING = 'pending',
  CONFIRMED = 'confirmed',
  CHECKED_IN = 'checked_in',
  CHECKED_OUT = 'checked_out',
  CANCELLED = 'cancelled',
  EXPIRED = 'expired'
}

export enum CheckInStatus {
  CHECKED_IN = 'checked_in',
  CHECKED_OUT = 'checked_out',
  NO_SHOW = 'no_show'
}

export enum PaymentMethod {
  CASH = 'cash',
  WECHAT = 'wechat',
  ALIPAY = 'alipay',
  CREDIT_CARD = 'credit_card',
  DEBIT_CARD = 'debit_card',
  TRANSFER = 'transfer'
}

export enum PaymentStatus {
  PENDING = 'pending',
  PAID = 'paid',
  REFUNDED = 'refunded',
  FAILED = 'failed'
}

export enum MemberLevel {
  NORMAL = 'normal',
  SILVER = 'silver',
  GOLD = 'gold',
  PLATINUM = 'platinum',
  DIAMOND = 'diamond'
}

export enum UserRole {
  ADMIN = 'admin',
  MANAGER = 'manager',
  RECEPTIONIST = 'receptionist',
  STAFF = 'staff'
}

export interface User {
  id: number
  username: string
  name: string
  phone: string
  email?: string
  role: UserRole
  avatar?: string
  status: boolean
  createdAt: string
  updatedAt: string
}

export interface RoomType {
  id: number
  name: string
  description?: string
  basePrice: number
  bedCount: number
  maxGuests: number
  facilities: string[]
  createdAt: string
  updatedAt: string
}

export interface Room {
  id: number
  roomNo: string
  floor: number
  roomTypeId: number
  roomType?: RoomType
  status: RoomStatus
  price: number
  facilities: string[]
  description?: string
  lastCleanedAt?: string
  createdAt: string
  updatedAt: string
}

export interface Booking {
  id: number
  bookingNo: string
  guestName: string
  guestPhone: string
  guestIdCard?: string
  roomTypeId?: number
  roomType?: RoomType
  roomId: number
  room?: Room
  checkInDate: string
  checkOutDate: string
  days: number
  adults?: number
  children?: number
  totalPrice: number
  deposit?: number
  status: BookingStatus
  source?: string
  remarks?: string
  memberId?: number
  member?: Member
  paidAmount?: number
  cancelDeadline?: string
  createdBy?: number
  createdAt: string
  updatedAt: string
}

export interface CheckIn {
  id: number
  checkInNo: string
  bookingId?: number
  booking?: Booking
  guestName: string
  guestPhone: string
  guestIdCard?: string
  roomId: number
  room?: Room
  checkInTime: string
  expectedCheckOut: string
  actualCheckOut?: string
  adults?: number
  children?: number
  deposit: number
  extraCharges?: number
  totalAmount: number
  status: CheckInStatus
  remarks?: string
  memberId?: number
  member?: Member
  createdBy?: number
  createdAt: string
  updatedAt: string
}

export interface Payment {
  id: number
  paymentNo: string
  bookingId?: number
  booking?: Booking
  checkInId?: number
  checkIn?: CheckIn
  amount: number
  method: PaymentMethod
  status: PaymentStatus
  transactionId?: string
  remark?: string
  paidAt?: string
  createdBy: number
  createdByUser?: User
  createdAt: string
  updatedAt: string
}

export interface Member {
  id: number
  memberNo: string
  name: string
  phone: string
  idCard: string
  level: MemberLevel
  points: number
  balance: number
  birthday?: string
  email?: string
  address?: string
  status: boolean
  totalSpent: number
  totalStays: number
  lastVisitAt?: string
  createdAt: string
  updatedAt: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PageParams {
  page: number
  pageSize: number
  keyword?: string
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface DashboardStats {
  todayCheckIns: number
  todayCheckOuts: number
  occupiedRooms: number
  availableRooms: number
  todayRevenue: number
  monthRevenue: number
  todayBookings: number
  pendingBookings: number
}

export interface ReportData {
  date: string
  revenue: number
  bookings: number
  checkIns: number
  occupancyRate: number
}
