import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types'

const request: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

request.interceptors.request.use(
  (config) => {
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

request.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse
    if (res.code !== 200) {
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res.data as any
  },
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
  register: (data: { username: string; password: string; phone: string; email?: string }) =>
    request.post('/auth/register', data),
  login: (data: { username: string; password: string }) =>
    request.post('/auth/login', data)
}

export const userApi = {
  getInfo: () => request.get('/user/info'),
  updateInfo: (data: any) => request.put('/user/info', data),
  changePassword: (data: { old_password: string; new_password: string }) =>
    request.post('/user/change-password', data),
  getMemberLevels: () => request.get('/user/member-levels'),
  exchangeCoupon: (data: { points: number }) => request.post('/user/exchange-coupon', data),
  getCoupons: () => request.get('/user/coupons')
}

export const showApi = {
  list: (params?: any) => request.get('/shows', { params }),
  get: (id: number) => request.get(`/shows/${id}`),
  create: (data: any) => request.post('/shows', data),
  update: (id: number, data: any) => request.put(`/shows/${id}`, data),
  delete: (id: number) => request.delete(`/shows/${id}`),
  getSessions: (showId: number) => request.get(`/shows/${showId}/sessions`),
  createSession: (data: any) => request.post('/shows/sessions', data)
}

export const seatApi = {
  getSeats: (sessionId: number) => request.get(`/seats/${sessionId}`),
  create: (data: any) => request.post('/seats', data),
  batchCreate: (data: any) => request.post('/seats/batch', data),
  lock: (data: { session_id: number; seat_ids: number[] }) => request.post('/seats/lock', data),
  unlock: (data: { session_id: number; seat_ids: number[] }) => request.post('/seats/unlock', data),
  getAreas: (sessionId: number) => request.get(`/seat-areas/${sessionId}`),
  createArea: (data: any) => request.post('/seat-areas', data),
  getChart: (sessionId: number) => request.get(`/seat-charts/${sessionId}`),
  updateChart: (data: any) => request.post('/seat-charts', data)
}

export const orderApi = {
  list: (params?: any) => request.get('/orders', { params }),
  get: (orderNo: string) => request.get(`/orders/${orderNo}`),
  create: (data: any) => request.post('/orders', data),
  pay: (data: { order_no: string; pay_type: number }) => request.post('/orders/pay', data),
  cancel: (orderNo: string) => request.post(`/orders/${orderNo}/cancel`),
  refund: (data: { order_no: string; reason: string }) => request.post('/orders/refund', data),
  auditRefund: (data: { refund_no: string; status: number; audit_remark?: string }) =>
    request.post('/orders/refund/audit', data),
  export: (params?: any) =>
    request.get('/orders/export/excel', { params, responseType: 'blob' })
}

export const checkinApi = {
  checkin: (data: { ticket_no: string }) => request.post('/checkin', data)
}

export const statisticsApi = {
  getSales: (params?: any) => request.get('/statistics/sales', { params }),
  getAreaSales: (sessionId: number) => request.get(`/statistics/area-sales/${sessionId}`),
  getSeatHeatmap: (sessionId: number) => request.get(`/statistics/seat-heatmap/${sessionId}`),
  getAudienceProfile: (params?: any) => request.get('/statistics/audience-profile', { params }),
  getDailySales: (params?: { start_date: string; end_date: string }) =>
    request.get('/statistics/daily-sales', { params }),
  exportPDF: (params?: any) =>
    request.get('/statistics/export/pdf', { params, responseType: 'blob' })
}

export default request
