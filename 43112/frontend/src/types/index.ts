export interface User {
  id: string
  username: string
  email: string
  nickname: string
  avatar: string
  role: 'student' | 'instructor' | 'admin'
  status: 'active' | 'disabled'
  phone?: string
  bio?: string
  instructor_status: 'pending' | 'approved' | 'rejected'
  email_verified: boolean
  created_at: string
}

export interface Course {
  id: string
  instructor_id: string
  title: string
  subtitle: string
  description: string
  cover: string
  category: string
  level: 'beginner' | 'intermediate' | 'advanced'
  price: number
  original_price: number
  status: 'draft' | 'published' | 'offline' | 'rejected'
  is_free: boolean
  tags: string
  total_hours: number
  student_count: number
  avg_rating: number
  review_count: number
  published_at?: string
  created_at: string
  instructor?: User
  chapters?: Chapter[]
}

export interface Chapter {
  id: string
  course_id: string
  title: string
  position: number
  is_free: boolean
  created_at: string
  lessons?: Lesson[]
}

export interface Lesson {
  id: string
  chapter_id: string
  title: string
  type: 'video' | 'document' | 'quiz'
  content: string
  video_url: string
  video_length: number
  doc_url: string
  doc_name: string
  position: number
  is_free: boolean
  is_published: boolean
  created_at: string
  quiz?: Quiz
}

export interface Quiz {
  id: string
  lesson_id: string
  title: string
  pass_score: number
  time_limit: number
  questions?: QuizQuestion[]
}

export interface QuizQuestion {
  id: string
  quiz_id: string
  content: string
  type: 'single' | 'multiple'
  score: number
  position: number
  options?: QuizOption[]
}

export interface QuizOption {
  id: string
  question_id: string
  content: string
  is_correct?: boolean
}

export interface Order {
  id: string
  user_id: string
  course_id: string
  order_no: string
  amount: number
  original_price: number
  discount: number
  coupon_id?: string
  status: 'pending' | 'paid' | 'refunding' | 'refunded' | 'cancelled' | 'failed'
  pay_method: 'alipay' | 'wechat' | 'balance'
  transaction_no: string
  paid_at?: string
  refund_reason?: string
  refunded_at?: string
  created_at: string
  course?: Course
  user?: User
}

export interface Coupon {
  id: string
  code: string
  type: 'fixed' | 'percent'
  value: number
  min_amount: number
  max_discount: number
  total_count: number
  used_count: number
  is_active: boolean
  expires_at?: string
  created_at: string
}

export interface Question {
  id: string
  course_id: string
  user_id: string
  lesson_id?: string
  title: string
  content: string
  reply_count: number
  view_count: number
  like_count: number
  is_resolved: boolean
  created_at: string
  user?: User
  answers?: Answer[]
}

export interface Answer {
  id: string
  question_id: string
  user_id: string
  content: string
  like_count: number
  is_best: boolean
  created_at: string
  user?: User
}

export interface Review {
  id: string
  course_id: string
  user_id: string
  rating: number
  content: string
  like_count: number
  is_anonymous: boolean
  created_at: string
  user?: User
}

export interface Progress {
  id: string
  user_id: string
  course_id: string
  lesson_id: string
  last_position: number
  total_duration: number
  is_completed: boolean
  completed_at?: string
  created_at: string
}

export interface Note {
  id: string
  user_id: string
  lesson_id: string
  content: string
  timestamp: number
  created_at: string
  user?: User
}

export interface InstructorApplication {
  id: string
  user_id: string
  real_name: string
  qualification: string
  experience: string
  certificates: string
  status: 'pending' | 'approved' | 'rejected'
  review_remark?: string
  reviewed_at?: string
  created_at: string
  user?: User
}

export interface FileInfo {
  id: string
  user_id: string
  file_type: 'image' | 'video' | 'document' | 'other'
  file_name: string
  file_url: string
  url: string
  file_size: number
  created_at: string
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

export interface PaginatedData<T = any> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface LoginRequest {
  account: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  nickname?: string
  role: 'student' | 'instructor'
}
