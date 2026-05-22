export interface User {
  id: string
  username: string
  email: string
  first_name: string
  last_name: string
  avatar: string
}

export interface AuthResponse {
  token: string
  expires_at: string
  user: User
}

export interface TravelPlan {
  id: string
  title: string
  description: string
  destination: string
  start_date: string
  end_date: string
  budget: number
  currency: string
  status: string
  owner_id: string
  version: number
  cover_image: string
  is_public: boolean
  created_at: string
  updated_at: string
  owner?: User
  participants?: PlanParticipant[]
  activities?: Activity[]
  files?: File[]
}

export interface PlanDetail extends TravelPlan {
  total_spent: number
  activity_count: number
  participant_count: number
}

export interface PlanParticipant {
  id: string
  plan_id: string
  user_id: string
  role: string
  can_edit: boolean
  can_delete: boolean
  joined_at: string
  user?: User
}

export interface Activity {
  id: string
  plan_id: string
  title: string
  description: string
  type: 'sightseeing' | 'transport' | 'accommodation' | 'food' | 'other'
  date: string
  start_time: string
  end_time: string
  location: string
  latitude: number
  longitude: number
  cost: number
  currency: string
  notes: string
  status: string
  booked: boolean
  confirmation: string
  contact_info: string
  order_index: number
  created_by: string
  created_at: string
  updated_at: string
  files?: File[]
}

export interface File {
  id: string
  plan_id: string
  activity_id: string
  uploaded_by: string
  file_name: string
  original_name: string
  file_type: string
  file_size: number
  file_url: string
  category: string
  description: string
  tags: string[]
  is_public: boolean
  created_at: string
  uploader?: User
}

export interface Checklist {
  id: string
  plan_id: string
  title: string
  type: 'packing' | 'preparation' | 'other'
  created_by: string
  created_at: string
  items: ChecklistItem[]
}

export interface ChecklistItem {
  id: string
  checklist_id: string
  title: string
  description: string
  category: string
  quantity: number
  is_completed: boolean
  completed_by: string
  completed_at: string
  order_index: number
  created_at: string
}

export interface Reminder {
  id: string
  plan_id: string
  activity_id: string
  user_id: string
  title: string
  description: string
  reminder_time: string
  is_sent: boolean
  sent_at: string
  channel: 'email' | 'app'
  created_at: string
  plan?: TravelPlan
  activity?: Activity
}

export interface BudgetSummary {
  plan_budget: number
  plan_currency: string
  total_spent: number
  budget_remaining: number
  budget_usage: number
  by_category: {
    type: string
    total: number
    count: number
  }[]
}

export interface MapLocation {
  id: string
  title: string
  type: string
  location: string
  latitude: number
  longitude: number
  date: string
  start_time: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  errors?: string[]
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface PaginationParams {
  page?: number
  page_size?: number
  search?: string
  sort_by?: string
  sort_dir?: 'asc' | 'desc'
}
