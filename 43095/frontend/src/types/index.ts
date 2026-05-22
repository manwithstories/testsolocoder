export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PaginationParams {
  page: number
  pageSize: number
}

export interface PaginationResult<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export type UserRole = 'admin' | 'doctor' | 'patient'

export interface User {
  id: number
  username: string
  email: string
  phone: string
  role: UserRole
  full_name: string
  gender: string
  birth_date: string | null
  avatar_url: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Department {
  id: number
  name: string
  description: string
  location: string
  created_at: string
  updated_at: string
}

export type DoctorTitle = '住院医师' | '主治医师' | '副主任医师' | '主任医师' | '教授'

export interface Doctor {
  id: number
  user_id: number
  user?: User
  department_id: number
  department?: Department
  title: DoctorTitle
  specialty: string
  introduction: string
  consultation_fee: number
  registration_fee: number
  average_rating: number
  review_count: number
  created_at: string
  updated_at: string
}

export type DayOfWeek = 0 | 1 | 2 | 3 | 4 | 5 | 6

export interface Schedule {
  id: number
  doctor_id: number
  doctor?: Doctor
  day_of_week: DayOfWeek
  start_time: string
  end_time: string
  max_patients: number
  time_slot_minutes: number
  is_available: boolean
  created_at: string
  updated_at: string
}

export interface Patient {
  id: number
  user_id: number
  user?: User
  id_card_no: string
  address: string
  emergency_contact_name: string
  emergency_contact_phone: string
  created_at: string
  updated_at: string
}

export interface HealthRecord {
  id: number
  patient_id: number
  medical_history: Record<string, unknown>
  allergies: Record<string, unknown>
  medications: Record<string, unknown>
  vaccinations: Record<string, unknown>
  family_history: string
  life_habits: string
  remarks: string
  created_at: string
  updated_at: string
}

export type AppointmentStatus = 'pending' | 'confirmed' | 'completed' | 'cancelled' | 'no_show'

export interface Appointment {
  id: number
  patient_id: number
  patient?: Patient
  doctor_id: number
  doctor?: Doctor
  appointment_date: string
  start_time: string
  end_time: string
  status: AppointmentStatus
  symptoms: string
  notes: string
  cancel_reason?: string
  consultation?: Consultation
  payment?: Payment
  review?: Review
  created_at: string
  updated_at: string
}

export interface Consultation {
  id: number
  appointment_id: number
  diagnosis: string
  treatment_plan: string
  doctor_notes: string
  follow_up_date: string | null
  prescription?: Prescription
  reports?: ExaminationReport[]
  created_at: string
  updated_at: string
}

export interface Prescription {
  id: number
  consultation_id: number
  prescription_no: string
  items?: PrescriptionItem[]
  notes: string
  is_fulfilled: boolean
  fulfilled_at: string | null
  created_at: string
  updated_at: string
}

export interface PrescriptionItem {
  id: number
  prescription_id: number
  drug_name: string
  specification: string
  dosage: string
  frequency: string
  duration: string
  quantity: number
  unit_price: number
  subtotal: number
  notes: string
  created_at: string
  updated_at: string
}

export interface ExaminationReport {
  id: number
  consultation_id: number
  report_type: string
  report_name: string
  file_url: string
  file_size: number
  content_type: string
  uploaded_by: number
  findings: string
  conclusion: string
  created_at: string
  updated_at: string
}

export type NotificationType = 'appointment_confirmation' | 'appointment_reminder' | 'appointment_cancelled' | 'consultation_completed' | 'payment_success' | 'system'

export type NotificationChannel = 'email' | 'in_app' | 'sms'

export interface Notification {
  id: number
  user_id: number
  type: NotificationType
  title: string
  content: string
  channel: NotificationChannel
  is_read: boolean
  read_at: string | null
  related_id: number
  related_type: string
  created_at: string
  updated_at: string
}

export type PaymentStatus = 'pending' | 'paid' | 'failed' | 'refunded'

export type PaymentMethod = 'wechat' | 'alipay' | 'card' | 'cash'

export interface Payment {
  id: number
  appointment_id: number
  appointment?: Appointment
  transaction_no: string
  registration_fee: number
  consultation_fee: number
  drug_fee: number
  examination_fee: number
  other_fee: number
  total_amount: number
  status: PaymentStatus
  method: PaymentMethod
  paid_at: string | null
  refunded_at: string | null
  notes: string
  created_at: string
  updated_at: string
}

export interface Review {
  id: number
  appointment_id: number
  patient_id: number
  doctor_id: number
  rating: number
  content: string
  is_anonymous: boolean
  is_verified: boolean
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface RegisterRequest {
  username: string
  email: string
  phone: string
  password: string
  full_name: string
  role: UserRole
}

export interface TimeSlot {
  start: string
  end: string
  available: boolean
}
