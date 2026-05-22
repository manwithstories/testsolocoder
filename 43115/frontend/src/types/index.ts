export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface User {
  id: number
  phone: string
  nickname: string
  avatar: string
  role: UserRole
  real_name: string
  id_card: string
  gender: string
  age: number
  provider_status: ProviderStatus
  certification_desc: string
  reject_reason: string
  rating: number
  order_count: number
  balance: number
  total_income: number
  is_active: boolean
  last_login_at: string | null
  created_at: string
  updated_at: string
  addresses?: Address[]
  certifications?: ServiceProviderCert[]
}

export type UserRole = 'customer' | 'service_provider' | 'admin'

export type ProviderStatus = 'pending' | 'approved' | 'rejected' | 'suspended'

export interface Address {
  id: number
  user_id: number
  contact_name: string
  contact_phone: string
  province: string
  city: string
  district: string
  address: string
  longitude: number
  latitude: number
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface ServiceProviderCert {
  id: number
  provider_id: number
  cert_type: string
  cert_name: string
  cert_number: string
  cert_image: string
  issued_at: string | null
  expired_at: string | null
  is_verified: boolean
  created_at: string
  updated_at: string
}

export interface ServiceCategory {
  id: number
  name: string
  description: string
  icon: string
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
  service_items?: ServiceItem[]
}

export interface ServiceItem {
  id: number
  category_id: number
  provider_id: number
  name: string
  description: string
  images: string
  base_price: number
  price_unit: string
  min_duration: number
  max_duration: number
  rating: number
  review_count: number
  order_count: number
  is_active: boolean
  created_at: string
  updated_at: string
  category?: ServiceCategory
  provider?: User
  service_areas?: ServiceArea[]
}

export interface ServiceArea {
  id: number
  province: string
  city: string
  district: string
  created_at: string
  updated_at: string
}

export type OrderStatus = 'pending' | 'confirmed' | 'in_service' | 'completed' | 'cancelled' | 'disputed'

export type InvitationStatus = 'pending' | 'accepted' | 'rejected' | 'expired'

export interface Order {
  id: number
  order_no: string
  customer_id: number
  provider_id: number | null
  service_item_id: number
  status: OrderStatus
  address_id: number
  service_address: string
  contact_name: string
  contact_phone: string
  longitude: number
  latitude: number
  appointment_time: string
  actual_start_time: string | null
  actual_end_time: string | null
  start_location: string
  end_location: string
  duration: number
  base_price: number
  total_amount: number
  platform_fee: number
  provider_income: number
  penalty_amount: number
  remark: string
  cancel_reason: string
  cancelled_by: string
  cancelled_at: string | null
  completed_at: string | null
  customer_rated: boolean
  provider_rated: boolean
  created_at: string
  updated_at: string
  customer?: User
  provider?: User
  service_item?: ServiceItem
  invitations?: OrderInvitation[]
  review?: Review
}

export interface OrderInvitation {
  id: number
  order_id: number
  provider_id: number
  status: InvitationStatus
  responded_at: string | null
  reject_reason: string
  expired_at: string
  created_at: string
  updated_at: string
  provider?: User
  order?: Order
}

export type ReviewType = 'customer_to_provider' | 'provider_to_customer'

export interface Review {
  id: number
  order_id: number
  reviewer_id: number
  reviewee_id: number
  review_type: ReviewType
  attitude_rating: number
  punctual_rating: number
  professional_rating: number
  overall_rating: number
  content: string
  images: string
  is_anonymous: boolean
  reply_content: string
  replied_at: string | null
  created_at: string
  updated_at: string
  reviewer?: User
  reviewee?: User
}

export type ComplaintStatus = 'pending' | 'processing' | 'resolved' | 'rejected'

export interface Complaint {
  id: number
  order_id: number
  complainant_id: number
  respondent_id: number
  status: ComplaintStatus
  title: string
  content: string
  images: string
  handler_id: number | null
  handle_result: string
  handled_at: string | null
  created_at: string
  updated_at: string
  complainant?: User
  respondent?: User
  order?: Order
}

export type BillType = 'income' | 'withdraw' | 'penalty' | 'refund' | 'commission'

export type BillStatus = 'pending' | 'completed' | 'failed'

export interface Bill {
  id: number
  provider_id: number
  order_id: number | null
  bill_type: BillType
  amount: number
  balance: number
  status: BillStatus
  description: string
  transaction_no: string
  settled_at: string | null
  created_at: string
  updated_at: string
  provider?: User
  order?: Order
}

export type WithdrawStatus = 'pending' | 'approved' | 'rejected' | 'completed' | 'failed'

export interface WithdrawRequest {
  id: number
  provider_id: number
  amount: number
  bank_name: string
  bank_account: string
  account_holder: string
  status: WithdrawStatus
  handler_id: number | null
  handle_remark: string
  remark: string
  transfer_no: string
  handled_at: string | null
  created_at: string
  updated_at: string
  provider?: User
}

export type MessageType = 'system' | 'order' | 'invitation' | 'review' | 'complaint' | 'withdraw'

export interface Message {
  id: number
  user_id: number
  type: MessageType
  title: string
  content: string
  related_id: number | null
  is_read: boolean
  read_at: string | null
  created_at: string
  updated_at: string
}

export interface PaginatedResponse<T> {
  total: number
  page: number
  page_size: number
  list: T[]
}

export interface LoginParams {
  phone: string
  password: string
}

export interface RegisterParams {
  phone: string
  password: string
  nickname: string
  role: UserRole
}

export interface CreateOrderParams {
  service_item_id: number
  address_id: number
  appointment_time: string
  duration: number
  remark?: string
}

export interface CreateReviewParams {
  order_id: number
  attitude_rating: number
  punctual_rating: number
  professional_rating: number
  content?: string
  images?: string
  is_anonymous?: boolean
}

export interface CreateComplaintParams {
  order_id: number
  title: string
  content: string
  images?: string
}

export interface CreateWithdrawParams {
  amount: number
  bank_name: string
  bank_account: string
  account_holder: string
}
