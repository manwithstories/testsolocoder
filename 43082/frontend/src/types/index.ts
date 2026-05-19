export interface Member {
  id: number
  name: string
  phone: string
  email: string
  gender: string
  birthday: string | null
  address: string
  profile_photo: string
  status: number
  membership_id: number | null
  membership: Membership | null
  created_at: string
  updated_at: string
}

export interface Membership {
  id: number
  member_id: number
  type: 'monthly' | 'quarter' | 'yearly'
  start_date: string
  end_date: string
  price: number
  status: number
  auto_renew: boolean
  created_at: string
  updated_at: string
}

export interface Coach {
  id: number
  name: string
  phone: string
  specialty: string
  description: string
  photo: string
  status: number
  created_at: string
  updated_at: string
}

export interface Course {
  id: number
  name: string
  description: string
  coach_id: number
  coach: Coach | null
  capacity: number
  duration: number
  type: 'single' | 'weekly' | 'monthly'
  weekdays: string
  start_date: string
  end_date: string | null
  start_time: string
  location: string
  status: number
  schedules: CourseSchedule[]
  created_at: string
  updated_at: string
}

export interface CourseSchedule {
  id: number
  course_id: number
  course: Course | null
  start_time: string
  end_time: string
  capacity: number
  booked_count: number
  status: number
  created_at: string
  updated_at: string
}

export interface Booking {
  id: number
  member_id: number
  member: Member | null
  schedule_id: number
  schedule: CourseSchedule | null
  status: number
  booking_time: string
  cancel_time: string | null
  check_in_id: number | null
  check_in: CheckIn | null
  created_at: string
  updated_at: string
}

export interface Waitlist {
  id: number
  member_id: number
  member: Member | null
  schedule_id: number
  schedule: CourseSchedule | null
  position: number
  notified: boolean
  status: number
  created_at: string
  updated_at: string
}

export interface CheckIn {
  id: number
  member_id: number
  member: Member | null
  schedule_id: number | null
  schedule: CourseSchedule | null
  check_in_time: string
  check_type: number
  remark: string
  created_at: string
  updated_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  errors?: any
}

export interface PaginatedResponse<T> extends ApiResponse<T> {
  pagination: {
    page: number
    page_size: number
    total: number
    total_page: number
  }
}

export interface LoginRequest {
  phone: string
  password: string
}

export interface LoginResponse {
  token: string
  member: Member
}
