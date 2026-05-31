import axios, { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig, AxiosResponse } from 'axios'

const API_BASE_URL = '/api/v1'

const api: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
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

export default api

export const authApi = {
  register: (data: { username: string; email: string; password: string; nickname?: string; user_type?: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
}

export const userApi = {
  getProfile: () => api.get('/users/profile'),
  updateProfile: (data: FormData | object) => api.put('/users/profile', data),
  changePassword: (data: { old_password: string; new_password: string }) =>
    api.put('/users/password', data),
  getUserById: (id: string) => api.get(`/users/${id}`),
}

export const plotApi = {
  create: (data: object) => api.post('/plots', data),
  getAll: (params?: object) => api.get('/plots', { params }),
  getById: (id: string) => api.get(`/plots/${id}`),
  update: (id: string, data: object) => api.put(`/plots/${id}`, data),
  delete: (id: string) => api.delete(`/plots/${id}`),
}

export const plantApi = {
  create: (data: object) => api.post('/plants', data),
  getAll: (params?: object) => api.get('/plants', { params }),
  getById: (id: string) => api.get(`/plants/${id}`),
}

export const plantingRecordApi = {
  create: (data: object) => api.post('/planting-records', data),
  getAll: (params?: object) => api.get('/planting-records', { params }),
  getById: (id: string) => api.get(`/planting-records/${id}`),
  update: (id: string, data: object) => api.put(`/planting-records/${id}`, data),
  delete: (id: string) => api.delete(`/planting-records/${id}`),
  getReport: (id: string) => api.get(`/planting-records/${id}/report`),
}

export const growthLogApi = {
  create: (data: object) => api.post('/growth-logs', data),
  getAll: (params?: object) => api.get('/growth-logs', { params }),
  getById: (id: string) => api.get(`/growth-logs/${id}`),
  update: (id: string, data: object) => api.put(`/growth-logs/${id}`, data),
  delete: (id: string) => api.delete(`/growth-logs/${id}`),
}

export const uploadApi = {
  upload: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post('/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },
}

export const calendarApi = {
  createEvent: (data: object) => api.post('/calendar/events', data),
  getEvents: (params?: object) => api.get('/calendar/events', { params }),
  getEvent: (id: string) => api.get(`/calendar/events/${id}`),
  updateEvent: (id: string, data: object) => api.put(`/calendar/events/${id}`, data),
  deleteEvent: (id: string) => api.delete(`/calendar/events/${id}`),
  getRecommendations: () => api.get('/calendar/recommendations'),
}

export const diseaseApi = {
  create: (data: object) => api.post('/disease-diagnosis', data),
  getAll: (params?: object) => api.get('/disease-diagnosis', { params }),
  getById: (id: string) => api.get(`/disease-diagnosis/${id}`),
  update: (id: string, data: object) => api.put(`/disease-diagnosis/${id}`, data),
  delete: (id: string) => api.delete(`/disease-diagnosis/${id}`),
}

export const postApi = {
  create: (data: object) => api.post('/posts', data),
  getAll: (params?: object) => api.get('/posts', { params }),
  getById: (id: string) => api.get(`/posts/${id}`),
  update: (id: string, data: object) => api.put(`/posts/${id}`, data),
  delete: (id: string) => api.delete(`/posts/${id}`),
  like: (id: string) => api.post(`/posts/${id}/like`),
  createComment: (postId: string, data: object) => api.post(`/posts/${postId}/comments`, data),
  getComments: (postId: string) => api.get(`/posts/${postId}/comments`),
  deleteComment: (commentId: string) => api.delete(`/posts/comments/${commentId}`),
}

export const followApi = {
  follow: (userId: string) => api.post(`/follows/${userId}`),
  getFollowers: (userId: string) => api.get(`/follows/${userId}/followers`),
  getFollowing: (userId: string) => api.get(`/follows/${userId}/following`),
}

export const exchangeApi = {
  create: (data: object) => api.post('/seed-exchanges', data),
  getAll: (params?: object) => api.get('/seed-exchanges', { params }),
  getById: (id: string) => api.get(`/seed-exchanges/${id}`),
  update: (id: string, data: object) => api.put(`/seed-exchanges/${id}`, data),
  delete: (id: string) => api.delete(`/seed-exchanges/${id}`),
  createOffer: (exchangeId: string, data: object) => api.post(`/seed-exchanges/${exchangeId}/offers`, data),
  getOffers: (exchangeId: string) => api.get(`/seed-exchanges/${exchangeId}/offers`),
  updateOffer: (offerId: string, data: object) => api.put(`/seed-exchanges/offers/${offerId}`, data),
  deleteOffer: (offerId: string) => api.delete(`/seed-exchanges/offers/${offerId}`),
}

export const productApi = {
  create: (data: object) => api.post('/products', data),
  getAll: (params?: object) => api.get('/products', { params }),
  getById: (id: string) => api.get(`/products/${id}`),
  update: (id: string, data: object) => api.put(`/products/${id}`, data),
  delete: (id: string) => api.delete(`/products/${id}`),
}

export const cartApi = {
  add: (data: object) => api.post('/cart', data),
  getAll: () => api.get('/cart'),
  update: (id: string, data: object) => api.put(`/cart/${id}`, data),
  remove: (id: string) => api.delete(`/cart/${id}`),
  clear: () => api.delete('/cart'),
}

export const orderApi = {
  create: (data: object) => api.post('/orders', data),
  getAll: (params?: object) => api.get('/orders', { params }),
  getById: (id: string) => api.get(`/orders/${id}`),
  update: (id: string, data: object) => api.put(`/orders/${id}`, data),
  cancel: (id: string) => api.post(`/orders/${id}/cancel`),
}
