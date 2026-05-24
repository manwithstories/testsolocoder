import axios, { AxiosInstance, InternalAxiosRequestConfig } from 'axios'

const api: AxiosInstance = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  register: (data: { email: string; password: string; firstName: string; lastName: string; role: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
}

export const userApi = {
  getProfile: () => api.get('/users/me'),
  updateProfile: (data: any) => api.put('/users/me', data),
  getUserById: (id: string) => api.get(`/users/${id}`),
}

export const teacherApi = {
  getProfile: () => api.get('/teacher/profile'),
  updateProfile: (data: any) => api.put('/teacher/profile', data),
  addSubject: (data: any) => api.post('/teacher/subjects', data),
  removeSubject: (subjectId: string) => api.delete(`/teacher/subjects/${subjectId}`),
  addAvailability: (data: any) => api.post('/teacher/availability', data),
  removeAvailability: (slotId: string) => api.delete(`/teacher/availability/${slotId}`),
  getEarnings: () => api.get('/teacher/earnings'),
  list: (params?: any) => api.get('/teachers', { params }),
  getById: (id: string) => api.get(`/teachers/${id}`),
  getReviews: (id: string) => api.get(`/teachers/${id}/reviews`),
}

export const studentApi = {
  getProfile: () => api.get('/student/profile'),
  updateProfile: (data: any) => api.put('/student/profile', data),
  addGoal: (data: any) => api.post('/student/goals', data),
  updateGoal: (goalId: string, data: any) => api.put(`/student/goals/${goalId}`, data),
  deleteGoal: (goalId: string) => api.delete(`/student/goals/${goalId}`),
  getAssessmentQuestions: (params?: any) => api.get('/student/assessment/questions', { params }),
  submitAssessment: (data: any) => api.post('/student/assessment/submit', data),
  getMyAssessment: () => api.get('/student/assessment'),
  matchTeachers: () => api.get('/student/match-teachers'),
  getMilestones: () => api.get('/student/milestones'),
  createMilestone: (data: any) => api.post('/student/milestones', data),
  updateMilestone: (id: string, data: any) => api.put(`/student/milestones/${id}`, data),
  getById: (id: string) => api.get(`/students/${id}`),
}

export const bookingApi = {
  getAll: (params?: any) => api.get('/bookings', { params }),
  getById: (id: string) => api.get(`/bookings/${id}`),
  create: (data: any) => api.post('/bookings', data),
  confirm: (id: string) => api.post(`/bookings/${id}/confirm`),
  reschedule: (data: any) => api.post('/bookings/reschedule', data),
  cancel: (data: any) => api.post('/bookings/cancel', data),
  complete: (id: string) => api.post(`/bookings/${id}/complete`),
}

export const videoApi = {
  createSession: (data: any) => api.post('/video/sessions', data),
  getSession: (id: string) => api.get(`/video/sessions/${id}`),
  getSessionByBooking: (bookingId: string) => api.get(`/video/sessions/booking/${bookingId}`),
  startSession: (data: any) => api.post('/video/sessions/start', data),
  endSession: (data: any) => api.post('/video/sessions/end', data),
  handleEvent: (sessionId: string, event: string) => api.post(`/video/sessions/${sessionId}/events?event=${event}`),
  getQuality: (id: string) => api.get(`/video/sessions/${id}/quality`),
}

export const walletApi = {
  getWallet: () => api.get('/wallet'),
  getTransactions: (params?: any) => api.get('/wallet/transactions', { params }),
  getTransactionById: (id: string) => api.get(`/wallet/transactions/${id}`),
  deposit: (data: any) => api.post('/wallet/deposit', data),
  withdraw: (data: any) => api.post('/wallet/withdraw', data),
  getWithdrawRequests: () => api.get('/wallet/withdraw-requests'),
}

export const notesApi = {
  getAll: (params?: any) => api.get('/notes', { params }),
  getById: (id: string) => api.get(`/notes/${id}`),
  create: (data: any) => api.post('/notes', data),
  update: (id: string, data: any) => api.put(`/notes/${id}`, data),
  delete: (id: string) => api.delete(`/notes/${id}`),
}

export const homeworkApi = {
  getAll: (params?: any) => api.get('/homework', { params }),
  create: (data: any) => api.post('/homework', data),
  submit: (data: any) => api.post('/homework/submit', data),
  grade: (id: string, data: any) => api.post(`/homework/${id}/grade`, data),
}

export const feedbackApi = {
  getAll: (params?: any) => api.get('/feedback', { params }),
  create: (data: any) => api.post('/feedback', data),
}

export const reviewApi = {
  getAll: (params?: any) => api.get('/reviews', { params }),
  getById: (id: string) => api.get(`/reviews/${id}`),
  create: (data: any) => api.post('/reviews', data),
  reply: (id: string, data: any) => api.post(`/reviews/${id}/reply`, data),
}

export const messageApi = {
  getMessages: (params?: any) => api.get('/messages', { params }),
  getConversations: () => api.get('/messages/conversations'),
  getUnreadCount: () => api.get('/messages/unread-count'),
  send: (data: any) => api.post('/messages', data),
  markRead: (params: any) => api.post('/messages/mark-read', null, { params }),
  delete: (id: string) => api.delete(`/messages/${id}`),
}

export const notificationApi = {
  getAll: (params?: any) => api.get('/notifications', { params }),
  getUnreadCount: () => api.get('/notifications/unread-count'),
  markRead: (id: string) => api.put(`/notifications/${id}/read`),
  markAllRead: () => api.put('/notifications/read-all'),
  delete: (id: string) => api.delete(`/notifications/${id}`),
}

export const subjectApi = {
  getAll: () => api.get('/subjects'),
}

export const adminApi = {
  getStats: () => api.get('/admin/stats'),
  getLogs: (params?: any) => api.get('/admin/logs', { params }),
  getAdminActions: () => api.get('/admin/admin-actions'),
  getPendingApprovals: () => api.get('/admin/pending-approvals'),
  approveTeacher: (id: string) => api.post(`/admin/teachers/${id}/approve`),
  rejectTeacher: (id: string, data: any) => api.post(`/admin/teachers/${id}/reject`, data),
  getWithdrawRequests: () => api.get('/admin/withdraw-requests'),
  processWithdraw: (id: string, data: any) => api.post(`/admin/withdraw-requests/${id}/process`, data),
  createSubject: (data: any) => api.post('/admin/subjects', data),
  updateSubject: (id: string, data: any) => api.put(`/admin/subjects/${id}`, data),
  getPaymentConfigs: () => api.get('/admin/payment-configs'),
  updatePaymentConfig: (id: string, data: any) => api.put(`/admin/payment-configs/${id}`, data),
  hideReview: (id: string) => api.post(`/admin/reviews/${id}/hide`),
}

export default api
