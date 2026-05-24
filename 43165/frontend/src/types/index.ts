export interface User {
  id: string;
  username: string;
  email: string;
  phone: string;
  real_name: string;
  role: 'employer' | 'agent' | 'temporary';
  avatar: string;
  credit_score: number;
  status: 'active' | 'inactive' | 'banned';
  company: string;
  address: string;
  id_card: string;
  last_login_at: string | null;
  last_login_ip: string;
  created_at: string;
  updated_at: string;
}

export interface JobPosting {
  id: string;
  employer_id: string;
  employer?: User;
  title: string;
  description: string;
  activity_type: string;
  position: string;
  location: string;
  latitude: number;
  longitude: number;
  start_date: string;
  end_date: string;
  salary_per_hour: number;
  salary_type: string;
  work_hours: string;
  headcount: number;
  applicants: number;
  hired_count: number;
  requirements: string;
  benefits: string;
  contact_person: string;
  contact_phone: string;
  status: string;
  tags: string;
  is_urgent: boolean;
  created_at: string;
  updated_at: string;
}

export interface JobApplication {
  id: string;
  job_id: string;
  job_posting?: JobPosting;
  temporary_id: string;
  temporary?: User;
  agent_id?: string | null;
  agent?: User | null;
  message: string;
  status: 'pending' | 'approved' | 'rejected';
  applied_at: string;
  reviewed_at: string | null;
  review_note: string;
  created_at: string;
  updated_at: string;
}

export interface Schedule {
  id: string;
  job_id: string;
  job_posting?: JobPosting;
  temporary_id: string;
  temporary?: User;
  shift_date: string;
  start_time: string;
  end_time: string;
  location: string;
  notes: string;
  status: 'scheduled' | 'in_progress' | 'completed';
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface CheckIn {
  id: string;
  schedule_id: string;
  schedule?: Schedule;
  temporary_id: string;
  temporary?: User;
  check_in_type: 'qr' | 'location' | 'face';
  check_in_time: string;
  check_out_time: string | null;
  latitude: number;
  longitude: number;
  location: string;
  face_verified: boolean;
  qr_code: string;
  status: 'checked_in' | 'checked_out';
  remarks: string;
  work_hours: number;
  created_at: string;
  updated_at: string;
}

export interface SalaryRecord {
  id: string;
  temporary_id: string;
  temporary?: User;
  employer_id: string;
  employer?: User;
  job_id: string;
  job_posting?: JobPosting;
  period_start: string;
  period_end: string;
  total_hours: number;
  base_salary: number;
  overtime_hours: number;
  overtime_pay: number;
  deductions: number;
  total_salary: number;
  status: 'pending' | 'paid';
  payment_method: string;
  payment_at: string | null;
  transaction_id: string;
  remark: string;
  created_at: string;
  updated_at: string;
}

export interface SalaryDetail {
  id: string;
  salary_id: string;
  check_in_id: string;
  date: string;
  work_hours: number;
  hourly_rate: number;
  amount: number;
  type: string;
  description: string;
}

export interface Evaluation {
  id: string;
  job_id: string;
  job_posting?: JobPosting;
  from_user_id: string;
  from_user?: User;
  to_user_id: string;
  to_user?: User;
  rating: number;
  content: string;
  tags: string;
  type: 'employer_to_temp' | 'temp_to_employer';
  is_anonymous: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data?: T;
}

export interface PaginatedData<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface LoginRequest {
  login: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  phone?: string;
  password: string;
  real_name: string;
  role: 'employer' | 'agent' | 'temporary';
  company?: string;
  address?: string;
  id_card?: string;
}

export interface LoginResponse {
  token: string;
  user: User;
  expire: string;
}

export interface CreateJobRequest {
  title: string;
  description: string;
  activity_type?: string;
  position: string;
  location: string;
  latitude?: number;
  longitude?: number;
  start_date: string;
  end_date: string;
  salary_per_hour: number;
  salary_type?: string;
  work_hours?: string;
  headcount: number;
  requirements?: string;
  benefits?: string;
  contact_person?: string;
  contact_phone?: string;
  tags?: string;
  is_urgent?: boolean;
}

export interface CreateScheduleRequest {
  job_id: string;
  temporary_id: string;
  shift_date: string;
  start_time: string;
  end_time: string;
  location?: string;
  notes?: string;
}

export interface BatchCreateScheduleRequest {
  job_id: string;
  temporary_ids: string[];
  shift_date: string;
  start_time: string;
  end_time: string;
  location?: string;
  notes?: string;
}

export interface CalculateSalaryRequest {
  temporary_id: string;
  employer_id: string;
  job_id: string;
  period_start: string;
  period_end: string;
  overtime_rate?: number;
  deductions?: number;
  deduction_note?: string;
}

export interface CreateEvaluationRequest {
  job_id: string;
  to_user_id: string;
  rating: number;
  content?: string;
  tags?: string;
  type: 'employer_to_temp' | 'temp_to_employer';
  is_anonymous?: boolean;
}
