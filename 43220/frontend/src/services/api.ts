import axios from 'axios'
import { useAuthStore } from '@/context/AuthContext'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => {
    const res = response.data
    if (res.code !== 0 && res.code !== undefined) {
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  register: (data: any) => api.post('/auth/register', data),
  login: (data: any) => api.post('/auth/login', data),
  getProfile: () => api.get('/auth/profile'),
  updateProfile: (data: any) => api.put('/auth/profile', data),
  changePassword: (data: any) => api.put('/auth/password', data),
  updateStoreInfo: (data: any) => api.put('/store/info', data),
  updateKeeperInfo: (data: any) => api.put('/keeper/info', data),
  listUsers: (params?: any) => api.get('/admin/users', { params }),
}

export const petApi = {
  create: (data: any) => api.post('/pets', data),
  list: (params?: any) => api.get('/pets', { params }),
  get: (id: string) => api.get(`/pets/${id}`),
  update: (id: string, data: any) => api.put(`/pets/${id}`, data),
  delete: (id: string) => api.delete(`/pets/${id}`),
  addVaccine: (data: any) => api.post('/pets/vaccines', data),
  getVaccines: (id: string) => api.get(`/pets/${id}/vaccines`),
  addDeworm: (data: any) => api.post('/pets/deworms', data),
  getDeworms: (id: string) => api.get(`/pets/${id}/deworms`),
}

export const packageApi = {
  create: (data: any) => api.post('/packages', data),
  list: (params?: any) => api.get('/packages', { params }),
  get: (id: string) => api.get(`/packages/${id}`),
  update: (id: string, data: any) => api.put(`/packages/${id}`, data),
  delete: (id: string) => api.delete(`/packages/${id}`),
}

export const reservationApi = {
  create: (data: any) => api.post('/reservations', data),
  list: (params?: any) => api.get('/reservations', { params }),
  get: (id: string) => api.get(`/reservations/${id}`),
  confirm: (id: string, data: any) => api.put(`/reservations/${id}/confirm`, data),
  checkIn: (id: string, data?: any) => api.put(`/reservations/${id}/check-in`, data || {}),
  checkOut: (id: string, data?: any) => api.put(`/reservations/${id}/check-out`, data || {}),
  cancel: (id: string, data: any) => api.put(`/reservations/${id}/cancel`, data),
}

export const dailyRecordApi = {
  create: (data: any) => api.post('/daily-records', data),
  listByReservation: (params?: any) => api.get('/daily-records', { params }),
  listByPet: (petId: string, params?: any) => api.get(`/daily-records/pet/${petId}`, { params }),
  get: (id: string) => api.get(`/daily-records/${id}`),
  update: (id: string, data: any) => api.put(`/daily-records/${id}`, data),
}

export const reviewApi = {
  create: (data: any) => api.post('/reviews', data),
  get: (id: string) => api.get(`/reviews/${id}`),
  listByStore: (params?: any) => api.get('/stores/reviews', { params }),
  listByKeeper: (params?: any) => api.get('/keepers/reviews', { params }),
  reply: (id: string, data: any) => api.put(`/reviews/${id}/reply`, data),
}

export const orderApi = {
  pay: (data: any) => api.post('/orders/pay', data),
  settle: (data: any) => api.post('/orders/settle', data),
  refund: (id: string, data: any) => api.post(`/orders/${id}/refund`, data),
  list: (params?: any) => api.get('/orders', { params }),
  get: (id: string) => api.get(`/orders/${id}`),
  getByReservation: (reservationId: string) => api.get(`/reservations/${reservationId}/orders`),
}

export const alertApi = {
  list: (params?: any) => api.get('/alerts', { params }),
  markAsRead: (id: string) => api.put(`/alerts/${id}/read`),
  markAllAsRead: () => api.put('/alerts/read-all'),
}

export const statisticsApi = {
  revenueTrend: (params?: any) => api.get('/statistics/revenue', { params }),
  occupancyRate: (params?: any) => api.get('/statistics/occupancy', { params }),
  petTypeDistribution: (params?: any) => api.get('/statistics/pet-types', { params }),
  orderStatistics: (params?: any) => api.get('/statistics/orders', { params }),
}

export const exportApi = {
  excel: (params?: any) => api.get('/export/excel', { params, responseType: 'blob' }),
  pdf: (params?: any) => api.get('/export/pdf', { params, responseType: 'blob' }),
}

export const uploadApi = {
  upload: (formData: FormData) => api.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }),
  uploadMultiple: (formData: FormData) => api.post('/upload/multiple', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }),
}

export default api
