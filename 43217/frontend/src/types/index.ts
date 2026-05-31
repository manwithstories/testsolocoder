export interface User {
  id: number
  username: string
  real_name: string
  phone: string
  email: string
  role: UserRole
  status: number
  company_id?: number
  agency_id?: number
  avatar?: string
  last_login_at?: string
  created_at: string
  updated_at: string
}

export type UserRole = 'hr' | 'employee' | 'agency' | 'admin'

export interface LoginResponse {
  token: string
  user: User
  expires_at: number
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  real_name: string
  phone: string
  email?: string
  role: UserRole
  company_id?: number
  agency_id?: number
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PaginationParams {
  page: number
  page_size: number
}

export interface PaginationResponse<T> {
  total: number
  page: number
  page_size: number
  items: T[]
}

export interface Company {
  id: number
  name: string
  unified_code: string
  legal_person: string
  contact_phone: string
  contact_email: string
  address: string
  business_license?: string
  status: number
  annual_budget: number
  used_budget: number
  payment_type: number
  balance: number
  credit_limit: number
  hr_user_id?: number
  created_at: string
  updated_at: string
  departments?: Department[]
}

export interface Department {
  id: number
  company_id: number
  parent_id?: number
  name: string
  manager_name?: string
  manager_phone?: string
  sort_order: number
  status: number
  created_at: string
  updated_at: string
}

export interface Employee {
  id: number
  company_id: number
  department_id: number
  user_id?: number
  employee_no?: string
  real_name: string
  gender: number
  birthday?: string
  id_card?: string
  phone?: string
  email?: string
  position?: string
  entry_date?: string
  status: number
  quota: number
  used_quota: number
  created_at: string
  updated_at: string
  department?: Department
}

export interface Agency {
  id: number
  name: string
  unified_code: string
  legal_person: string
  contact_phone: string
  contact_email: string
  address: string
  description?: string
  rating: number
  status: number
  created_at: string
  updated_at: string
  packages?: Package[]
}

export interface Package {
  id: number
  agency_id: number
  name: string
  description?: string
  original_price: number
  price: number
  suitable_for?: string
  gender_limit: number
  min_age: number
  max_age: number
  notes?: string
  status: number
  view_count: number
  sale_count: number
  created_at: string
  updated_at: string
  agency?: Agency
  items?: PackageItem[]
  time_slots?: TimeSlot[]
}

export interface PackageItem {
  id: number
  package_id: number
  item_name: string
  item_code?: string
  description?: string
  department?: string
  normal_range?: string
  unit?: string
  sort_order: number
  is_required: boolean
}

export interface TimeSlot {
  id: number
  package_id: number
  date: string
  start_time: string
  end_time: string
  total: number
  booked: number
  status: number
}

export interface Appointment {
  id: number
  employee_id: number
  company_id: number
  agency_id: number
  package_id: number
  time_slot_id: number
  appointment_no: string
  appointment_date: string
  start_time: string
  end_time: string
  status: number
  remark?: string
  is_reminded: boolean
  reminded_at?: string
  cancelled_at?: string
  cancel_reason?: string
  completed_at?: string
  created_at: string
  updated_at: string
  employee?: Employee
  company?: Company
  agency?: Agency
  package?: Package
  time_slot?: TimeSlot
  report?: Report
}

export interface Report {
  id: number
  appointment_id: number
  employee_id: number
  agency_id: number
  package_id: number
  report_no: string
  report_date: string
  doctor_name?: string
  summary?: string
  suggestion?: string
  has_abnormal: boolean
  status: number
  pdf_file?: string
  viewed_at?: string
  created_at: string
  updated_at: string
  items?: ReportItem[]
}

export interface ReportItem {
  id: number
  report_id: number
  package_item_id?: number
  item_name: string
  item_code?: string
  result: string
  unit?: string
  normal_range?: string
  is_abnormal: boolean
  abnormal_type?: string
  description?: string
  suggestion?: string
  department?: string
}

export interface HealthRecord {
  id: number
  employee_id: number
  company_id: number
  record_year: number
  report_id?: number
  height: number
  weight: number
  bmi: number
  blood_pressure?: string
  heart_rate: number
  has_abnormal: boolean
  abnormal_count: number
  tags?: string
  summary?: string
  created_at: string
  updated_at: string
}

export interface AbnormalItem {
  id: number
  employee_id: number
  health_record_id: number
  item_name: string
  item_code?: string
  result: string
  normal_range?: string
  abnormal_type?: string
  level: number
  suggestion?: string
  recheck_date?: string
  recheck_status: number
  created_at: string
  updated_at: string
}

export interface RecheckReminder {
  id: number
  employee_id: number
  abnormal_id: number
  remind_date: string
  remind_type: string
  content: string
  is_read: boolean
  is_sent: boolean
  sent_at?: string
  created_at: string
  updated_at: string
}

export interface Billing {
  id: number
  billing_no: string
  company_id: number
  agency_id: number
  period: string
  total_amount: number
  paid_amount: number
  status: number
  due_date?: string
  paid_at?: string
  invoice_no?: string
  invoice_file?: string
  invoice_status: number
  remark?: string
  created_at: string
  updated_at: string
  items?: BillingItem[]
  company?: Company
  agency?: Agency
}

export interface BillingItem {
  id: number
  billing_id: number
  appointment_id: number
  employee_id: number
  package_id: number
  package_name: string
  unit_price: number
  quantity: number
  amount: number
  appointment_date: string
}

export interface Transaction {
  id: number
  transaction_no: string
  company_id: number
  type: string
  amount: number
  balance: number
  payment_method: string
  status: number
  remark?: string
  created_at: string
  updated_at: string
}

export interface CompanyBudget {
  id: number
  company_id: number
  year: number
  total_budget: number
  used_budget: number
  frequency: number
  period_start?: string
  period_end?: string
}

export interface DepartmentAppointment {
  id: number
  company_id: number
  department_id: number
  agency_id: number
  package_id: number
  year: number
  total_quota: number
  used_quota: number
  start_date?: string
  end_date?: string
  status: number
  created_at: string
  updated_at: string
  agency?: Agency
  package?: Package
}
