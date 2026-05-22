import { request } from './api'
import {
  User, Course, Chapter, Lesson, Order, Question, Answer,
  Review, Progress, Note, Coupon, InstructorApplication, FileInfo,
  LoginRequest, RegisterRequest, PaginatedData, ApiResponse,
} from '@/types'

export const authApi = {
  login: (data: LoginRequest) =>
    request<{ token: string; refresh_token: string; user: User }>({
      url: '/auth/login', method: 'post', data,
    }),
  register: (data: RegisterRequest) =>
    request<{ id: string; username: string; role: string }>({
      url: '/auth/register', method: 'post', data,
    }),
  refreshToken: () =>
    request<{ token: string }>({ url: '/auth/refresh', method: 'post' }),
  getProfile: () =>
    request<User>({ url: '/auth/profile', method: 'get' }),
  updateProfile: (data: Partial<User>) =>
    request<{ message: string }>({ url: '/auth/profile', method: 'put', data }),
  changePassword: (data: { old_password: string; new_password: string }) =>
    request<{ message: string }>({ url: '/auth/password', method: 'put', data }),
  applyInstructor: (data: any) =>
    request<any>({ url: '/auth/instructor-apply', method: 'post', data }),
}

export const courseApi = {
  list: (params?: any) =>
    request<PaginatedData<Course>>({ url: '/courses', method: 'get', params }),
  get: (id: string) =>
    request<Course>({ url: `/courses/${id}`, method: 'get' }),
  myCourses: (params?: any) =>
    request<PaginatedData<Course>>({ url: '/courses/my', method: 'get', params }),
  create: (data: any) =>
    request<Course>({ url: '/courses', method: 'post', data }),
  update: (id: string, data: any) =>
    request<{ message: string }>({ url: `/courses/${id}`, method: 'put', data }),
  updateStatus: (id: string, status: string) =>
    request<{ message: string }>({ url: `/courses/${id}/status`, method: 'put', data: { status } }),
  delete: (id: string) =>
    request<{ message: string }>({ url: `/courses/${id}`, method: 'delete' }),
  categories: () =>
    request<any[]>({ url: '/courses/categories', method: 'get' }),
}

export const chapterApi = {
  create: (courseId: string, data: any) =>
    request<Chapter>({ url: `/courses/${courseId}/chapters`, method: 'post', data }),
  update: (id: string, data: any) =>
    request<{ message: string }>({ url: `/chapters/${id}`, method: 'put', data }),
  delete: (id: string) =>
    request<{ message: string }>({ url: `/chapters/${id}`, method: 'delete' }),
}

export const lessonApi = {
  create: (chapterId: string, data: any) =>
    request<Lesson>({ url: `/chapters/${chapterId}/lessons`, method: 'post', data }),
  update: (id: string, data: any) =>
    request<{ message: string }>({ url: `/lessons/${id}`, method: 'put', data }),
  delete: (id: string) =>
    request<{ message: string }>({ url: `/lessons/${id}`, method: 'delete' }),
}

export const quizApi = {
  create: (lessonId: string, data: any) =>
    request<any>({ url: `/lessons/${lessonId}/quiz`, method: 'post', data }),
  submit: (quizId: string, data: any) =>
    request<{ score: number; total_score: number; is_passed: boolean }>({
      url: `/quizzes/${quizId}/submit`, method: 'post', data,
    }),
}

export const orderApi = {
  create: (data: any) =>
    request<Order>({ url: '/orders', method: 'post', data }),
  pay: (id: string) =>
    request<{ message: string; order: Order }>({ url: `/orders/${id}/pay`, method: 'post' }),
  get: (id: string) =>
    request<Order>({ url: `/orders/${id}`, method: 'get' }),
  myOrders: (params?: any) =>
    request<PaginatedData<Order>>({ url: '/orders/my', method: 'get', params }),
  listAll: (params?: any) =>
    request<PaginatedData<Order>>({ url: '/admin/orders', method: 'get', params }),
  refund: (id: string, reason: string) =>
    request<{ message: string }>({ url: `/orders/${id}/refund`, method: 'post', data: { reason } }),
  processRefund: (id: string, data: any) =>
    request<{ message: string }>({ url: `/admin/orders/${id}/refund/process`, method: 'put', data }),
  updateStatus: (id: string, status: string) =>
    request<{ message: string }>({ url: `/admin/orders/${id}/status`, method: 'put', data: { status } }),
}

export const couponApi = {
  validate: (code: string, courseId: string) =>
    request<{ coupon: Coupon; discount: number; final: number }>({
      url: '/coupons/validate', method: 'get', params: { code, course_id: courseId },
    }),
  list: () =>
    request<Coupon[]>({ url: '/admin/coupons', method: 'get' }),
  create: (data: any) =>
    request<Coupon>({ url: '/admin/coupons', method: 'post', data }),
  update: (id: string, data: any) =>
    request<{ message: string }>({ url: `/admin/coupons/${id}`, method: 'put', data }),
  delete: (id: string) =>
    request<{ message: string }>({ url: `/admin/coupons/${id}`, method: 'delete' }),
}

export const progressApi = {
  update: (lessonId: string, data: any) =>
    request<{ message: string }>({ url: `/progress/lessons/${lessonId}`, method: 'put', data }),
  getCourseProgress: (courseId: string) =>
    request<{ progresses: Progress[]; total_lessons: number; completed_lessons: number; completion_rate: number }>({
      url: `/progress/courses/${courseId}`, method: 'get',
    }),
  getLessonProgress: (lessonId: string) =>
    request<Progress>({ url: `/progress/lessons/${lessonId}`, method: 'get' }),
}

export const noteApi = {
  create: (data: any) =>
    request<Note>({ url: '/notes', method: 'post', data }),
  list: (params?: any) =>
    request<PaginatedData<Note>>({ url: '/notes', method: 'get', params }),
  update: (id: string, data: any) =>
    request<{ message: string }>({ url: `/notes/${id}`, method: 'put', data }),
  delete: (id: string) =>
    request<{ message: string }>({ url: `/notes/${id}`, method: 'delete' }),
}

export const qaApi = {
  listQuestions: (params?: any) =>
    request<PaginatedData<Question>>({ url: '/questions', method: 'get', params }),
  getQuestion: (id: string) =>
    request<Question>({ url: `/questions/${id}`, method: 'get' }),
  createQuestion: (data: any) =>
    request<Question>({ url: '/questions', method: 'post', data }),
  deleteQuestion: (id: string) =>
    request<{ message: string }>({ url: `/questions/${id}`, method: 'delete' }),
  likeQuestion: (id: string) =>
    request<{ liked: boolean }>({ url: `/questions/${id}/like`, method: 'post' }),
  createAnswer: (data: any) =>
    request<Answer>({ url: '/answers', method: 'post', data }),
  markBest: (answerId: string) =>
    request<{ message: string }>({ url: `/answers/${answerId}/best`, method: 'post' }),
  likeAnswer: (answerId: string) =>
    request<{ liked: boolean }>({ url: `/answers/${answerId}/like`, method: 'post' }),
  deleteAnswer: (answerId: string) =>
    request<{ message: string }>({ url: `/answers/${answerId}`, method: 'delete' }),
}

export const reviewApi = {
  list: (params?: any) =>
    request<PaginatedData<Review>>({ url: '/reviews', method: 'get', params }),
  getMy: (courseId: string) =>
    request<Review>({ url: '/reviews/my', method: 'get', params: { course_id: courseId } }),
  create: (data: any) =>
    request<Review>({ url: '/reviews', method: 'post', data }),
  delete: (id: string) =>
    request<{ message: string }>({ url: `/reviews/${id}`, method: 'delete' }),
  summary: (courseId: string) =>
    request<{ avg_rating: number; review_count: number; rating_distribution: any[] }>({
      url: `/reviews/summary/${courseId}`, method: 'get',
    }),
}

export const uploadApi = {
  upload: (file: File, fileType: string, onProgress?: (progress: number) => void) => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('file_type', fileType)
    return request<FileInfo>({
      url: '/upload',
      method: 'post',
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          onProgress(Math.round((progressEvent.loaded * 100) / progressEvent.total))
        }
      },
    })
  },
  list: (params?: any) =>
    request<PaginatedData<FileInfo>>({ url: '/files', method: 'get', params }),
  delete: (id: string) =>
    request<{ message: string }>({ url: `/files/${id}`, method: 'delete' }),
}

export const analyticsApi = {
  instructorDashboard: () =>
    request<any>({ url: '/analytics/instructor/dashboard', method: 'get' }),
  instructorRevenue: (params?: any) =>
    request<any>({ url: '/analytics/instructor/revenue', method: 'get', params }),
  instructorStats: () =>
    request<any>({ url: '/analytics/instructor/stats', method: 'get' }),
  adminDashboard: () =>
    request<any>({ url: '/admin/analytics/dashboard', method: 'get' }),
  adminRevenue: (params?: any) =>
    request<any>({ url: '/admin/analytics/revenue', method: 'get', params }),
  adminStats: () =>
    request<any>({ url: '/admin/analytics/stats', method: 'get' }),
}

export const userApi = {
  list: (params?: any) =>
    request<PaginatedData<User>>({ url: '/admin/users', method: 'get', params }),
  updateStatus: (id: string, status: string) =>
    request<{ message: string }>({ url: `/admin/users/${id}/status`, method: 'put', data: { status } }),
  listApplications: (params?: any) =>
    request<PaginatedData<InstructorApplication>>({
      url: '/admin/instructor-applications', method: 'get', params,
    }),
  reviewApplication: (id: string, data: any) =>
    request<{ message: string }>({
      url: `/admin/instructor-applications/${id}/review`, method: 'put', data,
    }),
}
