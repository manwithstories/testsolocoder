export interface User {
  id: number
  email: string
  role: UserRole
  status: UserStatus
  company_id?: number
  company?: Company
  profile?: ApplicantProfile
  last_login_at?: string
  created_at: string
  updated_at: string
}

export type UserRole = 'admin' | 'company' | 'applicant'
export type UserStatus = 'active' | 'inactive' | 'banned'

export interface ApplicantProfile {
  id: number
  user_id: number
  full_name: string
  phone?: string
  avatar?: string
  gender?: string
  birth_date?: string
  location?: string
  education?: string
  experience?: string
  skills?: string
  summary?: string
  resume_file?: string
  created_at: string
  updated_at: string
}

export interface Company {
  id: number
  name: string
  logo?: string
  industry?: string
  size?: string
  address?: string
  website?: string
  description?: string
  verified: boolean
  created_at: string
  updated_at: string
}

export type JobStatus = 'open' | 'closed' | 'paused' | 'draft'

export interface Job {
  id: number
  company_id: number
  company?: Company
  title: string
  description: string
  salary_min?: number
  salary_max?: number
  salary_type?: string
  location: string
  job_type?: string
  experience?: string
  education?: string
  skills?: string
  requirements?: string
  benefits?: string
  deadline?: string
  status: JobStatus
  view_count: number
  apply_count: number
  created_at: string
  updated_at: string
}

export interface Resume {
  id: number
  user_id: number
  user?: User
  title: string
  full_name: string
  email?: string
  phone?: string
  location?: string
  education?: string
  experience?: string
  skills?: string
  summary?: string
  projects?: string
  file_url?: string
  file_type?: string
  file_size?: number
  is_default: boolean
  created_at: string
  updated_at: string
}

export type ApplicationStatus = 'pending' | 'viewed' | 'interested' | 'interview' | 'accepted' | 'rejected' | 'withdrawn'

export interface Application {
  id: number
  job_id: number
  job?: Job
  applicant_id: number
  applicant?: User
  resume_id: number
  resume?: Resume
  status: ApplicationStatus
  cover_letter?: string
  hr_note?: string
  applied_at: string
  last_update_at: string
  created_at: string
  updated_at: string
}

export interface ApplicationHistory {
  id: number
  application_id: number
  old_status?: ApplicationStatus
  new_status: ApplicationStatus
  changed_by: number
  change_reason?: string
  created_at: string
}

export type InterviewStatus = 'pending' | 'accepted' | 'rejected' | 'completed' | 'cancelled'

export interface Interview {
  id: number
  application_id: number
  application?: Application
  job_id: number
  applicant_id: number
  interviewer?: string
  interview_type?: string
  location?: string
  meeting_link?: string
  scheduled_at: string
  duration: number
  status: InterviewStatus
  notes?: string
  feedback?: string
  rating?: number
  created_at: string
  updated_at: string
}

export interface DailyStatistics {
  date: string
  new_jobs: number
  new_applications: number
  new_users: number
  total_views: number
  interviews: number
  hires: number
}

export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
  errors?: any
}

export interface PaginatedResponse<T> {
  items: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  role: 'company' | 'applicant'
  full_name: string
  phone?: string
  company?: {
    name: string
    industry?: string
    size?: string
    address?: string
  }
}

export interface CreateJobRequest {
  title: string
  description: string
  salary_min?: number
  salary_max?: number
  salary_type?: string
  location: string
  job_type?: string
  experience?: string
  education?: string
  skills?: string
  requirements?: string
  benefits?: string
  deadline?: string
}

export interface CreateResumeRequest {
  title: string
  full_name: string
  email?: string
  phone?: string
  location?: string
  education?: string
  experience?: string
  skills?: string
  summary?: string
  projects?: string
  is_default?: boolean
}

export interface ApplyRequest {
  job_id: number
  resume_id: number
  cover_letter?: string
}

export interface ScheduleInterviewRequest {
  application_id: number
  interviewer?: string
  interview_type?: string
  location?: string
  meeting_link?: string
  scheduled_at: string
  duration?: number
  notes?: string
}
