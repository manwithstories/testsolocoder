import { request } from './request'
import type {
  User,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  Department,
  Doctor,
  Schedule,
  Patient,
  HealthRecord,
  Appointment,
  Consultation,
  Prescription,
  ExaminationReport,
  Notification,
  Payment,
  Review,
  PaginationParams,
  PaginationResult,
  TimeSlot
} from '@/types'

export const authAPI = {
  login: (data: LoginRequest) => request<LoginResponse>({ url: '/auth/login', method: 'post', data }),
  register: (data: RegisterRequest) => request<User>({ url: '/auth/register', method: 'post', data }),
  getCurrentUser: () => request<User>({ url: '/auth/me', method: 'get' }),
  updatePassword: (data: { old_password: string; new_password: string }) =>
    request<void>({ url: '/auth/password', method: 'put', data }),
  logout: () => request<void>({ url: '/auth/logout', method: 'post' })
}

export const departmentAPI = {
  getList: (params?: PaginationParams & { keyword?: string }) =>
    request<PaginationResult<Department>>({ url: '/departments', method: 'get', params }),
  getDetail: (id: number) => request<Department>({ url: `/departments/${id}`, method: 'get' }),
  create: (data: Partial<Department>) => request<Department>({ url: '/departments', method: 'post', data }),
  update: (id: number, data: Partial<Department>) =>
    request<Department>({ url: `/departments/${id}`, method: 'put', data }),
  delete: (id: number) => request<void>({ url: `/departments/${id}`, method: 'delete' })
}

export const doctorAPI = {
  getList: (params?: PaginationParams & { department_id?: number; keyword?: string; sort_by_rating?: boolean }) =>
    request<PaginationResult<Doctor>>({ url: '/doctors', method: 'get', params }),
  getDetail: (id: number) => request<Doctor>({ url: `/doctors/${id}`, method: 'get' }),
  create: (data: Partial<Doctor> & { user_id: number; department_id: number }) =>
    request<Doctor>({ url: '/doctors', method: 'post', data }),
  update: (id: number, data: Partial<Doctor>) => request<Doctor>({ url: `/doctors/${id}`, method: 'put', data }),
  delete: (id: number) => request<void>({ url: `/doctors/${id}`, method: 'delete' }),
  getSchedules: (doctorId: number) => request<Schedule[]>({ url: `/doctors/${doctorId}/schedules`, method: 'get' }),
  createSchedule: (doctorId: number, data: Partial<Schedule>) =>
    request<Schedule>({ url: `/doctors/${doctorId}/schedules`, method: 'post', data }),
  updateSchedule: (scheduleId: number, data: Partial<Schedule>) =>
    request<Schedule>({ url: `/doctors/schedules/${scheduleId}`, method: 'put', data }),
  deleteSchedule: (scheduleId: number) =>
    request<void>({ url: `/doctors/schedules/${scheduleId}`, method: 'delete' }),
  getAvailableSlots: (doctorId: number, date: string) =>
    request<TimeSlot[]>({ url: `/doctors/${doctorId}/available-slots`, method: 'get', params: { date } })
}

export const patientAPI = {
  getList: (params?: PaginationParams & { keyword?: string }) =>
    request<PaginationResult<Patient>>({ url: '/patients', method: 'get', params }),
  getDetail: (id: number) => request<Patient>({ url: `/patients/${id}`, method: 'get' }),
  update: (id: number, data: Partial<Patient>) =>
    request<Patient>({ url: `/patients/${id}`, method: 'put', data }),
  delete: (id: number) => request<void>({ url: `/patients/${id}`, method: 'delete' })
}

export const healthRecordAPI = {
  get: (patientUserId: number) =>
    request<HealthRecord>({ url: `/health-records/${patientUserId}`, method: 'get' }),
  update: (patientUserId: number, data: Partial<HealthRecord>) =>
    request<HealthRecord>({ url: `/health-records/${patientUserId}`, method: 'put', data }),
  getVisitHistory: (patientUserId: number, params?: PaginationParams) =>
    request<PaginationResult<Appointment>>({
      url: `/health-records/${patientUserId}/visit-history`,
      method: 'get',
      params
    }),
  export: (patientUserId: number) =>
    request<Blob>({ url: `/health-records/${patientUserId}/export`, method: 'get', responseType: 'blob' })
}

export const appointmentAPI = {
  create: (data: { doctor_id: number; appointment_date: string; start_time: string; end_time: string; symptoms?: string }) =>
    request<Appointment>({ url: '/appointments', method: 'post', data }),
  getList: (params?: PaginationParams & { doctor_id?: number; patient_id?: number; status?: string; date?: string }) =>
    request<PaginationResult<Appointment>>({ url: '/appointments', method: 'get', params }),
  getDetail: (id: number) => request<Appointment>({ url: `/appointments/${id}`, method: 'get' }),
  cancel: (id: number, data: { cancel_reason?: string }) =>
    request<Appointment>({ url: `/appointments/${id}/cancel`, method: 'put', data }),
  reschedule: (id: number, data: { appointment_date: string; start_time: string; end_time: string }) =>
    request<Appointment>({ url: `/appointments/${id}/reschedule`, method: 'put', data }),
  confirm: (id: number) => request<Appointment>({ url: `/appointments/${id}/confirm`, method: 'put' }),
  complete: (id: number) => request<Appointment>({ url: `/appointments/${id}/complete`, method: 'put' }),
  checkAvailability: (params: { doctor_id: number; date: string; start_time: string; end_time: string }) =>
    request<{ available: boolean; reason?: string }>({ url: '/appointments/check-availability', method: 'get', params })
}

export const consultationAPI = {
  create: (data: { appointment_id: number; diagnosis: string; treatment_plan?: string; doctor_notes?: string }) =>
    request<Consultation>({ url: '/consultations', method: 'post', data }),
  getDetail: (id: number) => request<Consultation>({ url: `/consultations/${id}`, method: 'get' }),
  update: (id: number, data: Partial<Consultation>) =>
    request<Consultation>({ url: `/consultations/${id}`, method: 'put', data }),
  createPrescription: (data: {
    consultation_id: number
    notes?: string
    items: Array<{
      drug_name: string
      specification?: string
      dosage: string
      frequency?: string
      duration?: string
      quantity?: number
      unit_price?: number
      notes?: string
    }>
  }) => request<Prescription>({ url: '/consultations/prescriptions', method: 'post', data }),
  uploadReport: (data: FormData) =>
    request<ExaminationReport>({ url: '/consultations/reports', method: 'post', data }),
  getHistory: (params?: PaginationParams) =>
    request<PaginationResult<Consultation>>({ url: '/consultations/patient/history', method: 'get', params })
}

export const notificationAPI = {
  getList: (params?: PaginationParams & { is_read?: boolean }) =>
    request<PaginationResult<Notification>>({ url: '/notifications', method: 'get', params }),
  getUnreadCount: () => request<{ count: number }>({ url: '/notifications/unread-count', method: 'get' }),
  markAsRead: (id: number) => request<Notification>({ url: `/notifications/${id}/read`, method: 'put' }),
  markAllAsRead: () => request<void>({ url: '/notifications/read-all', method: 'put' }),
  delete: (id: number) => request<void>({ url: `/notifications/${id}`, method: 'delete' })
}

export const paymentAPI = {
  create: (data: { appointment_id: number }) => request<Payment>({ url: '/payments', method: 'post', data }),
  getDetail: (id: number) => request<Payment>({ url: `/payments/${id}`, method: 'get' }),
  getList: (params?: PaginationParams & { status?: string; start_date?: string; end_date?: string }) =>
    request<PaginationResult<Payment>>({ url: '/payments', method: 'get', params }),
  updateStatus: (id: number, data: { status: string; method?: string; transaction_no?: string }) =>
    request<Payment>({ url: `/payments/${id}/status`, method: 'put', data }),
  getReport: (params: { start_date: string; end_date: string; group_by?: 'date' | 'department' | 'doctor' }) =>
    request<Array<Record<string, unknown>>>({ url: '/payments/report', method: 'get', params }),
  exportReport: (params: { start_date: string; end_date: string; group_by?: 'date' | 'department' | 'doctor' }) =>
    request<Blob>({ url: '/payments/export', method: 'get', params, responseType: 'blob' })
}

export const reviewAPI = {
  create: (data: { appointment_id: number; rating: number; content?: string; is_anonymous?: boolean }) =>
    request<Review>({ url: '/reviews', method: 'post', data }),
  getList: (params?: PaginationParams & { doctor_id?: number }) =>
    request<PaginationResult<Review>>({ url: '/reviews', method: 'get', params }),
  getDetail: (id: number) => request<Review>({ url: `/reviews/${id}`, method: 'get' }),
  delete: (id: number) => request<void>({ url: `/reviews/${id}`, method: 'delete' })
}

export const uploadAPI = {
  upload: (file: File, onProgress?: (percent: number) => void) => {
    const formData = new FormData()
    formData.append('file', file)
    return request<{ file_url: string; file_name: string }>({
      url: '/upload',
      method: 'post',
      data: formData,
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          onProgress(Math.round((progressEvent.loaded * 100) / progressEvent.total))
        }
      }
    })
  },
  download: (filename: string) =>
    request<Blob>({ url: `/upload/${filename}`, method: 'get', responseType: 'blob' }),
  delete: (filename: string) => request<void>({ url: `/upload/${filename}`, method: 'delete' })
}
