export interface User {
  id: number
  username: string
  email: string
  full_name: string
  phone?: string
  avatar?: string
  role: string
  status: string
  created_at: string
}

export interface LoginRequest {
  account: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  full_name: string
  phone?: string
  role?: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PaginatedData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface Wedding {
  id: number
  user_id: number
  title: string
  groom_name: string
  bride_name: string
  wedding_date: string
  budget: number
  style?: string
  theme_color?: string
  guest_count: number
  venue?: string
  venue_address?: string
  venue_capacity: number
  description?: string
  status: string
  is_template: boolean
  created_at: string
  updated_at: string
}

export interface Vendor {
  id: number
  user_id: number
  wedding_id?: number
  name: string
  category: string
  contact_person?: string
  phone?: string
  email?: string
  address?: string
  website?: string
  service_area?: string
  price_range?: string
  rating: number
  review_count: number
  notes?: string
  status: string
  created_at: string
  updated_at: string
}

export interface Guest {
  id: number
  wedding_id: number
  first_name: string
  last_name: string
  full_name: string
  email?: string
  phone?: string
  group?: string
  relation?: string
  rsvp_status: string
  plus_one: boolean
  plus_one_name?: string
  table_id?: number
  seat_number: number
  notes?: string
  is_vip: boolean
  created_at: string
  updated_at: string
}

export interface GuestTable {
  id: number
  wedding_id: number
  table_name: string
  table_number: number
  capacity: number
  seats_json?: string
  notes?: string
  table_guests?: Guest[]
  created_at: string
  updated_at: string
}

export interface BudgetItem {
  id: number
  wedding_id: number
  vendor_id?: number
  category: string
  description?: string
  estimated_cost: number
  actual_cost: number
  paid_amount: number
  status: string
  due_date?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface Payment {
  id: number
  budget_item_id: number
  amount: number
  method?: string
  status: string
  paid_at?: string
  reference?: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface Task {
  id: number
  wedding_id: number
  title: string
  description?: string
  category?: string
  assignee?: string
  due_date?: string
  priority: string
  status: string
  completed_at?: string
  parent_id?: number
  template_id?: number
  is_template: boolean
  order: number
  created_at: string
  updated_at: string
}

export interface TaskTemplate {
  id: number
  name: string
  category?: string
  tasks_json: string
  is_default: boolean
  created_at: string
}

export interface Document {
  id: number
  wedding_id: number
  vendor_id?: number
  budget_item_id?: number
  file_name: string
  file_path: string
  file_size: number
  file_type: string
  category?: string
  version: number
  parent_id?: number
  uploaded_by: number
  notes?: string
  created_at: string
  updated_at: string
}

export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  is_read: boolean
  related_id: number
  created_at: string
}

export interface DashboardStats {
  total_weddings: number
  total_guests: number
  total_vendors: number
  budget: {
    total: number
    paid: number
    rate: number
  }
  rsvp: {
    Accepted: number
    Declined: number
    Pending: number
  }
  tasks: {
    Total: number
    Completed: number
    Pending: number
    Overdue: number
  }
}
