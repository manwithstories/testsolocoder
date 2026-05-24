import { api } from './client'
import { ApiResponse, PaginatedData } from './auth'

export interface Job {
  id: number
  company_id: number
  department_id?: number
  title: string
  location: string
  salary_min: number
  salary_max: number
  salary_type: string
  description: string
  requirements?: string
  skills?: string
  job_type: 'full-time' | 'part-time' | 'contract' | 'internship' | 'remote'
  status: 'open' | 'paused' | 'closed'
  views: number
  company?: Company
  department?: Department
  applications?: Application[]
  created_at: string
  updated_at: string
}

export interface Department {
  id: number
  company_id: number
  name: string
  jobs?: Job[]
}

export interface Application {
  id: number
  job_id: number
  jobseeker_id: number
  resume_id: number
  status: 'pending' | 'reviewed' | 'interview' | 'accepted' | 'rejected' | 'hold'
  cover_letter?: string
  job?: Job
  jobseeker?: JobSeeker
  resume?: Resume
  interviews?: Interview[]
  created_at: string
  updated_at: string
}

export interface Resume {
  id: number
  jobseeker_id: number
  title: string
  full_name: string
  email?: string
  phone?: string
  location?: string
  summary?: string
  file_path?: string
  file_name?: string
  is_default: boolean
  education_list?: Education[]
  work_experiences?: WorkExperience[]
  skills?: Skill[]
  created_at: string
  updated_at: string
}

export interface Education {
  id: number
  resume_id: number
  school: string
  degree?: string
  major?: string
  start_date?: string
  end_date?: string
  description?: string
}

export interface WorkExperience {
  id: number
  resume_id: number
  company: string
  position: string
  start_date?: string
  end_date?: string
  description?: string
}

export interface Skill {
  id: number
  resume_id: number
  name: string
  level?: string
}

export interface Interview {
  id: number
  application_id: number
  scheduled_at: string
  duration: number
  location: string
  interviewer: string
  interviewer_email?: string
  notes?: string
  status: 'scheduled' | 'confirmed' | 'declined' | 'completed' | 'cancelled'
  application?: Application
  review?: Review
  created_at: string
  updated_at: string
}

export interface Review {
  id: number
  interview_id: number
  rating: number
  feedback: string
  strengths?: string
  weaknesses?: string
  status: 'offer' | 'pass' | 'reject' | 'pending'
  interview?: Interview
  created_at: string
  updated_at: string
}

export interface Company {
  id: number
  user_id: number
  company_name: string
  industry?: string
  company_size?: string
  company_type?: string
  address?: string
  website?: string
  logo?: string
  description?: string
}

export interface JobSeeker {
  id: number
  user_id: number
  user?: {
    id: number
    name: string
    email: string
    phone?: string
  }
}

export const jobApi = {
  search: (params?: Record<string, string>) =>
    api.get<ApiResponse<PaginatedData<Job>>>('/jobs/search', { params }),

  getById: (id: number) =>
    api.get<ApiResponse<Job>>(`/jobs/${id}`),

  getSimilar: (id: number) =>
    api.get<ApiResponse<Job[]>>(`/jobs/${id}/similar`),

  create: (data: Partial<Job>) =>
    api.post<ApiResponse<Job>>('/jobs', data),

  update: (id: number, data: Partial<Job>) =>
    api.put<ApiResponse<Job>>(`/jobs/${id}`, data),

  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/jobs/${id}`),

  listCompanyJobs: (params?: Record<string, string>) =>
    api.get<ApiResponse<PaginatedData<Job>>>('/company/jobs', { params }),

  updateStatus: (id: number, status: string) =>
    api.put<ApiResponse<null>>(`/jobs/${id}/status`, { status }),

  getStats: () =>
    api.get<ApiResponse<{ total_jobs: number; open_jobs: number; total_views: number; total_applications: number }>>('/company/jobs/stats'),

  getJobStatistics: (id: number) =>
    api.get<ApiResponse<{ job_id: number; job_title: string; views: number; applications: number; interviews: number; offers: number; conversion_rate: number; interview_rate: number }>>(`/company/jobs/${id}/statistics`),

  getJobTypes: () =>
    api.get<ApiResponse<{ value: string; label: string }[]>>('/job-types'),

  getLocations: () =>
    api.get<ApiResponse<string[]>>('/locations'),

  getSalaryRanges: () =>
    api.get<ApiResponse<{ value: string; label: string; min: number; max: number }[]>>('/salary-ranges'),
}

export const departmentApi = {
  list: () =>
    api.get<ApiResponse<Department[]>>('/departments'),

  create: (data: { name: string }) =>
    api.post<ApiResponse<Department>>('/departments', data),

  update: (id: number, data: { name: string }) =>
    api.put<ApiResponse<Department>>(`/departments/${id}`, data),

  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/departments/${id}`),
}

export const applicationApi = {
  create: (data: { job_id: number; resume_id: number; cover_letter?: string }) =>
    api.post<ApiResponse<Application>>('/applications', data),

  getById: (id: number) =>
    api.get<ApiResponse<Application>>(`/company/applications/${id}`),

  listJobSeeker: (params?: Record<string, string>) =>
    api.get<ApiResponse<Application[]>>('/jobseeker/applications', { params }),

  listCompany: (params?: Record<string, string>) =>
    api.get<ApiResponse<Application[]>>('/company/applications', { params }),

  updateStatus: (id: number, status: string) =>
    api.put<ApiResponse<null>>(`/company/applications/${id}/status`, { status }),

  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/applications/${id}`),
}

export const resumeApi = {
  list: () =>
    api.get<ApiResponse<Resume[]>>('/resumes'),

  getById: (id: number) =>
    api.get<ApiResponse<Resume>>(`/resumes/${id}`),

  create: (data: Partial<Resume> & {
    education_list?: Omit<Education, 'id' | 'resume_id'>[]
    work_experiences?: Omit<WorkExperience, 'id' | 'resume_id'>[]
    skills?: Omit<Skill, 'id' | 'resume_id'>[]
  }) => api.post<ApiResponse<Resume>>('/resumes', data),

  update: (id: number, data: Partial<Resume> & {
    education_list?: Omit<Education, 'id' | 'resume_id'>[]
    work_experiences?: Omit<WorkExperience, 'id' | 'resume_id'>[]
    skills?: Omit<Skill, 'id' | 'resume_id'>[]
  }) => api.put<ApiResponse<Resume>>(`/resumes/${id}`, data),

  delete: (id: number) =>
    api.delete<ApiResponse<null>>(`/resumes/${id}`),

  uploadFile: (id: number, formData: FormData) =>
    api.upload<ApiResponse<{ file_path: string; file_name: string }>>(`/resumes/${id}/upload`, formData),

  setDefault: (id: number) =>
    api.put<ApiResponse<null>>(`/resumes/${id}/default`),
}

export const interviewApi = {
  create: (data: {
    application_id: number
    scheduled_at: string
    duration?: number
    location: string
    interviewer: string
    interviewer_email?: string
    notes?: string
  }) => api.post<ApiResponse<Interview>>('/interviews', data),

  getById: (id: number) =>
    api.get<ApiResponse<Interview>>(`/interviews/${id}`),

  listCompany: (params?: Record<string, string>) =>
    api.get<ApiResponse<Interview[]>>('/company/interviews', { params }),

  listJobSeeker: (params?: Record<string, string>) =>
    api.get<ApiResponse<Interview[]>>('/jobseeker/interviews', { params }),

  update: (id: number, data: Partial<Interview>) =>
    api.put<ApiResponse<null>>(`/interviews/${id}`, data),

  confirm: (id: number) =>
    api.put<ApiResponse<null>>(`/interviews/${id}/confirm`),

  decline: (id: number, data?: { reason: string }) =>
    api.put<ApiResponse<null>>(`/interviews/${id}/decline`, data),

  cancel: (id: number) =>
    api.put<ApiResponse<null>>(`/interviews/${id}/cancel`),
}

export const reviewApi = {
  create: (data: {
    interview_id: number
    rating: number
    feedback: string
    strengths?: string
    weaknesses?: string
    status: 'offer' | 'pass' | 'reject' | 'pending'
  }) => api.post<ApiResponse<Review>>('/reviews', data),

  getByInterview: (interviewId: number) =>
    api.get<ApiResponse<Review>>(`/interviews/${interviewId}/review`),

  update: (id: number, data: Partial<Review>) =>
    api.put<ApiResponse<null>>(`/reviews/${id}`, data),

  listCompany: (params?: Record<string, string>) =>
    api.get<ApiResponse<Review[]>>('/company/reviews', { params }),

  listJobSeeker: () =>
    api.get<ApiResponse<Review[]>>('/jobseeker/reviews'),
}

export const recommendationApi = {
  getRecommended: () =>
    api.get<ApiResponse<(Job & { match_score: number; matched_skills: string[] })[]>>('/recommendations'),
}

export const statisticsApi = {
  getCompanyStats: () =>
    api.get<ApiResponse<{
      total_jobs: number
      open_jobs: number
      total_applications: number
      total_interviews: number
      total_offers: number
      interview_rate: number
      offer_rate: number
      total_views: number
      conversion_rate: number
    }>>('/company/statistics'),

  getJobSeekerStats: () =>
    api.get<ApiResponse<{
      total_applications: number
      interview_count: number
      offer_count: number
      rejection_count: number
      pending_count: number
      success_rate: number
      resume_count: number
    }>>('/jobseeker/statistics'),

  getApplicationTrend: (params?: Record<string, string>) =>
    api.get<ApiResponse<{ date: string; applications: number; interviews: number }[]>>('/company/trends', { params }),

  getNotifications: () =>
    api.get<ApiResponse<Notification[]>>('/notifications'),

  markNotificationRead: (id: number) =>
    api.put<ApiResponse<null>>(`/notifications/${id}/read`),

  markAllNotificationsRead: () =>
    api.put<ApiResponse<null>>('/notifications/read-all'),
}

export interface Notification {
  id: number
  user_id: number
  type: string
  title: string
  message: string
  read: boolean
  related_id?: number
  created_at: string
}

export const exportApi = {
  applications: (params?: Record<string, string>) =>
    api.get<Blob>('/export/applications', { params, responseType: 'blob' }),

  interviews: (params?: Record<string, string>) =>
    api.get<Blob>('/export/interviews', { params, responseType: 'blob' }),

  jobs: () =>
    api.get<Blob>('/export/jobs', { responseType: 'blob' }),
}
