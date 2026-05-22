import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import router from '@/router'

const request: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 15000
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
    const res = response.data
    if (res.code !== 0 && res.code !== undefined) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        router.push('/login')
        ElMessage.error('登录已过期，请重新登录')
      } else if (status === 403) {
        ElMessage.error('没有权限访问')
      } else {
        const msg = error.response.data?.message || '请求失败'
        ElMessage.error(msg)
      }
    } else {
      ElMessage.error('网络连接失败')
    }
    return Promise.reject(error)
  }
)

export const api = {
  login: (data: { email: string; password: string }) =>
    request.post('/auth/login', data),

  register: (data: { username: string; email: string; password: string; real_name?: string; phone?: string; department?: string }) =>
    request.post('/auth/register', data),

  getProfile: () =>
    request.get('/user/profile'),

  updateProfile: (data: any) =>
    request.put('/user/profile', data),

  listUsers: (params?: { page?: number; page_size?: number; role?: string }) =>
    request.get('/admin/users', { params }),

  updateUserRole: (id: number, role: string) =>
    request.put(`/admin/users/${id}/role`, { role }),

  deleteUser: (id: number) =>
    request.delete(`/admin/users/${id}`),

  getRooms: (params?: { page?: number; page_size?: number; floor?: string; equipment?: string }) =>
    request.get('/rooms', { params }),

  listAllRooms: () =>
    request.get('/rooms/list'),

  getRoom: (id: number) =>
    request.get(`/rooms/${id}`),

  createRoom: (data: any) =>
    request.post('/admin/rooms', data),

  updateRoom: (id: number, data: any) =>
    request.put(`/admin/rooms/${id}`, data),

  deleteRoom: (id: number) =>
    request.delete(`/admin/rooms/${id}`),

  uploadRoomPhoto: (id: number, file: File) => {
    const formData = new FormData()
    formData.append('photo', file)
    return request.post(`/admin/rooms/${id}/photos`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  deleteRoomPhoto: (id: number) =>
    request.delete(`/admin/rooms/photos/${id}`),

  getFloors: () =>
    request.get('/rooms/floors'),

  createBooking: (data: any) =>
    request.post('/bookings', data),

  getBookings: (params?: { page?: number; page_size?: number; room_id?: number; status?: number }) =>
    request.get('/bookings', { params }),

  getBooking: (id: number) =>
    request.get(`/bookings/${id}`),

  cancelBooking: (id: number, reason: string) =>
    request.delete(`/bookings/${id}`, { data: { reason } }),

  rescheduleBooking: (id: number, data: { start_time: string; end_time: string }) =>
    request.put(`/bookings/${id}/reschedule`, data),

  approveBooking: (id: number) =>
    request.put(`/admin/bookings/${id}/approve`),

  completeBooking: (id: number) =>
    request.put(`/admin/bookings/${id}/complete`),

  getWeekCalendar: (params: { date?: string; room_id?: number; floor?: string }) =>
    request.get('/calendar/week', { params }),

  getMonthCalendar: (params: { date?: string; room_id?: number; floor?: string }) =>
    request.get('/calendar/month', { params }),

  getRoomAvailability: (id: number, date: string) =>
    request.get(`/calendar/rooms/${id}/availability`, { params: { date } }),

  uploadMaterial: (bookingId: number, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return request.post(`/bookings/${bookingId}/materials`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  getMaterials: (bookingId: number) =>
    request.get(`/bookings/${bookingId}/materials`),

  deleteMaterial: (id: number) =>
    request.delete(`/materials/${id}`),

  getStats: (params: { start_date: string; end_date: string; department?: string }) =>
    request.get('/stats', { params }),

  exportStats: (params: { start_date: string; end_date: string; department?: string }) =>
    request.get('/stats/export', { params, responseType: 'blob' as const })
}

export default request
