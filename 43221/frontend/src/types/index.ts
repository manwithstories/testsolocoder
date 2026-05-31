export type UserRole = 'client' | 'professional' | 'admin'

export type VerificationStatus = 'pending' | 'approved' | 'rejected'

export interface User {
  id: string
  username: string
  email: string
  role: UserRole
  full_name: string
  avatar?: string
  phone?: string
  verification_status: VerificationStatus
  verification_note?: string
  is_active: boolean
  last_login_at?: string
  created_at: string
  updated_at: string
}

export type ServiceType = 'legal' | 'counseling' | 'financial' | 'other'

export type ServiceStatus = 'active' | 'inactive'

export interface Service {
  id: string
  professional_id: string
  title: string
  description?: string
  service_type: ServiceType
  price: number
  duration_minutes: number
  status: ServiceStatus
  average_rating: number
  review_count: number
  tags?: string
  created_at: string
  updated_at: string
  professional?: User
  schedules?: Schedule[]
}

export interface Schedule {
  id: string
  service_id: string
  date: string
  start_time: string
  end_time: string
  is_booked: boolean
  is_available: boolean
  created_at: string
  updated_at: string
  appointment?: Appointment
}

export type AppointmentStatus = 'pending' | 'confirmed' | 'completed' | 'cancelled' | 'refunded'

export interface Appointment {
  id: string
  client_id: string
  professional_id: string
  service_id: string
  schedule_id: string
  status: AppointmentStatus
  notes?: string
  cancel_reason?: string
  cancelled_at?: string
  completed_at?: string
  created_at: string
  updated_at: string
  client?: User
  professional?: User
  service?: Service
  schedule?: Schedule
  payment?: Payment
  consult_record?: ConsultRecord
  review?: Review
}

export type PaymentMethod = 'online' | 'offline'

export type PaymentStatus = 'pending' | 'paid' | 'refunded' | 'failed' | 'cancelled'

export interface Payment {
  id: string
  appointment_id: string
  client_id: string
  professional_id: string
  amount: number
  payment_method: PaymentMethod
  status: PaymentStatus
  transaction_id?: string
  refund_reason?: string
  refund_status?: string
  refunded_at?: string
  paid_at?: string
  expires_at: string
  created_at: string
  updated_at: string
}

export interface ConsultRecord {
  id: string
  appointment_id: string
  client_id: string
  professional_id: string
  summary?: string
  advice?: string
  follow_up_date?: string
  is_confidential: boolean
  created_at: string
  updated_at: string
  appointment?: Appointment
  client?: User
  professional?: User
}

export type ReviewStatus = 'pending' | 'approved' | 'rejected'

export interface Review {
  id: string
  appointment_id: string
  client_id: string
  professional_id: string
  service_id: string
  rating: number
  content?: string
  status: ReviewStatus
  reject_reason?: string
  created_at: string
  updated_at: string
  appointment?: Appointment
  client?: User
  professional?: User
  service?: Service
}

export type NotificationType =
  | 'appointment_success'
  | 'appointment_cancel'
  | 'appointment_remind'
  | 'payment_success'
  | 'payment_refund'
  | 'review_reply'
  | 'system'

export interface Notification {
  id: string
  user_id: string
  type: NotificationType
  title: string
  content: string
  data?: string
  is_read: boolean
  read_at?: string
  created_at: string
  user?: User
}

export interface NotificationTemplate {
  id: string
  type: string
  title: string
  content: string
  variables?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginatedResponse<T = any> {
  total: number
  page: number
  page_size: number
  items: T[]
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_at: string
  user: User
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  role: UserRole
  full_name: string
  phone?: string
  verification_docs?: string
}

export interface CreateServiceRequest {
  title: string
  description?: string
  service_type: ServiceType
  price: number
  duration_minutes: number
  tags?: string
}

export interface CreateScheduleRequest {
  service_id: string
  date: string
  start_time: string
  end_time: string
}

export interface BatchCreateScheduleRequest {
  service_id: string
  start_date: string
  end_date: string
  time_slots: { start_time: string; end_time: string }[]
}

export interface CreateAppointmentRequest {
  service_id: string
  schedule_id: string
  notes?: string
}

export interface CreateConsultRecordRequest {
  appointment_id: string
  summary?: string
  advice?: string
  follow_up_date?: string
  is_confidential?: boolean
}

export interface CreateReviewRequest {
  appointment_id: string
  rating: number
  content?: string
}

export interface ProfessionalStats {
  appointments: {
    total: number
    pending: number
    confirmed: number
    completed: number
    cancelled: number
    refunded: number
  }
  revenue: {
    total_revenue: number
    paid_count: number
    refunded_amount: number
  }
  reviews: {
    average_rating: number
    total_reviews: number
  }
}

export interface AdminStats {
  total_users: number
  total_clients: number
  total_professionals: number
  total_services: number
  appointments: {
    total: number
    pending: number
    confirmed: number
    completed: number
    cancelled: number
    refunded: number
  }
  revenue: {
    total_revenue: number
    paid_count: number
    refunded_amount: number
  }
}
