export interface User {
  id: number
  username: string
  real_name: string
  id_card: string
  phone: string
  email: string
  verified: boolean
  role: string
  status: number
  last_login_at?: string
  created_at: string
}

export interface EventItem {
  id: number
  event_id: number
  name: string
  category: string
  gender: string
  min_age: number
  max_age: number
  quota: number
  waitlist_quota: number
  fee: number
  requirements: string
  status: string
}

export interface Event {
  id: number
  name: string
  description: string
  location: string
  start_date: string
  end_date: string
  registration_deadline: string
  status: string
  organizer: string
  cover_image: string
  is_published: boolean
  items?: EventItem[]
}

export interface Registration {
  id: number
  user_id: number
  event_id: number
  event_item_id: number
  team_name?: string
  team_members?: string
  reg_type: string
  status: string
  queue_position: number
  payment_status: string
  amount: number
  remark?: string
  created_at: string
  user?: User
  event?: Event
  event_item?: EventItem
}

export interface Score {
  id: number
  event_id: number
  event_item_id: number
  user_id: number
  registration_id: number
  score: number
  rank: number
  points: number
  time_used?: string
  remarks?: string
  is_valid: boolean
}

export interface Certificate {
  id: number
  user_id: number
  event_id: number
  event_item_id: number
  score_id?: number
  certificate_no: string
  certificate_name: string
  rank: number
  score: number
  file_path?: string
  status: string
  retry_count: number
  generated_at?: string
  created_at: string
}

export interface Message {
  id: number
  user_id: number
  type: string
  title: string
  content: string
  is_read: boolean
  extra?: string
  created_at: string
}

export interface Paged<T> {
  total: number
  page: number
  page_size: number
  list: T[]
}

export interface ApiResp<T> {
  code: number
  message: string
  data: T
}

export interface LoginResp {
  token: string
  user_id: number
  username: string
  role: string
  real_name: string
  verified: boolean
}
