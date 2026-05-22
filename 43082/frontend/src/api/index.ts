import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'
import { message } from 'antd'
import {
  User,
  LoginResponse,
  Shop,
  Category,
  Product,
  CartItem,
  Order,
  Review,
  Notification,
  PaginatedResponse,
  ProductQueryParams,
  ShopStatistics,
  AdminStatistics,
} from '@/types'

const camelToSnake = (str: string): string => {
  return str.replace(/[A-Z]/g, (letter) => `_${letter.toLowerCase()}`)
}

const snakeToCamel = (str: string): string => {
  return str.replace(/_([a-z])/g, (_, letter) => letter.toUpperCase())
}

const convertKeysToSnake = (obj: any): any => {
  if (Array.isArray(obj)) {
    return obj.map(convertKeysToSnake)
  } else if (obj !== null && typeof obj === 'object' && !(obj instanceof FormData) && !(obj instanceof Date)) {
    return Object.keys(obj).reduce((acc: any, key) => {
      acc[camelToSnake(key)] = convertKeysToSnake(obj[key])
      return acc
    }, {})
  }
  return obj
}

const convertKeysToCamel = (obj: any): any => {
  if (Array.isArray(obj)) {
    return obj.map(convertKeysToCamel)
  } else if (obj !== null && typeof obj === 'object' && !(obj instanceof Date)) {
    return Object.keys(obj).reduce((acc: any, key) => {
      acc[snakeToCamel(key)] = convertKeysToCamel(obj[key])
      return acc
    }, {})
  }
  return obj
}

const api: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    if (config.data && typeof config.data === 'object') {
      config.data = convertKeysToSnake(config.data)
    }
    if (config.params) {
      config.params = convertKeysToSnake(config.params)
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => {
    const raw = response.data
    if (raw && typeof raw === 'object') {
      return convertKeysToCamel(raw)
    }
    return raw
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    if (error.response?.data?.message) {
      message.error(error.response.data.message)
    } else {
      message.error('请求失败')
    }
    return Promise.reject(error)
  }
)

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export const authAPI = {
  register: (data: {
    username: string
    email: string
    phone: string
    password: string
    nickname?: string
    role: string
  }) => api.post<any, ApiResponse<LoginResponse>>('/auth/register', data),

  login: (data: { account: string; password: string }) =>
    api.post<any, ApiResponse<LoginResponse>>('/auth/login', data),

  getProfile: () => api.get<any, ApiResponse<User>>('/auth/profile'),

  updateProfile: (data: { nickname?: string; avatar?: string; phone?: string }) =>
    api.put('/auth/profile', data),

  changePassword: (data: { oldPassword: string; newPassword: string }) =>
    api.put('/auth/password', data),
}

export const categoryAPI = {
  list: () => api.get<any, ApiResponse<Category[]>>('/categories'),
  getAll: () => api.get<any, ApiResponse<Category[]>>('/categories/all'),
  create: (data: { name: string; icon?: string; parentId?: number; sort?: number }) =>
    api.post('/categories', data),
  update: (id: number, data: { name?: string; icon?: string; sort?: number; status?: string }) =>
    api.put(`/categories/${id}`, data),
  delete: (id: number) => api.delete(`/categories/${id}`),
}

export const shopAPI = {
  list: (params?: { page?: number; pageSize?: number; status?: string; keyword?: string }) =>
    api.get<any, ApiResponse<PaginatedResponse<Shop[]>>>('/shops', { params }),
  getById: (id: number) => api.get<any, ApiResponse<Shop>>(`/shops/${id}`),
  apply: (data: {
    name: string
    description?: string
    logo?: string
    contactName: string
    contactPhone: string
    address?: string
    idCardFront?: string
    idCardBack?: string
    businessLicense?: string
  }) => api.post('/shops/apply', data),
  getMyShop: () => api.get<any, ApiResponse<Shop>>('/shops/my'),
  update: (data: {
    name?: string
    description?: string
    logo?: string
    contactName?: string
    contactPhone?: string
    address?: string
  }) => api.put('/shops/my', data),
  review: (id: number, data: { status: string; rejectReason?: string }) =>
    api.put(`/shops/${id}/review`, data),
}

export const productAPI = {
  list: (params?: ProductQueryParams) =>
    api.get<any, ApiResponse<PaginatedResponse<Product[]>>>('/products', { params }),
  getById: (id: number) => api.get<any, ApiResponse<Product>>(`/products/${id}`),
  create: (data: any) => api.post('/products', data),
  update: (id: number, data: any) => api.put(`/products/${id}`, data),
  delete: (id: number) => api.delete(`/products/${id}`),
  myProducts: (params?: { page?: number; pageSize?: number; status?: string }) =>
    api.get<any, ApiResponse<PaginatedResponse<Product[]>>>('/products/my/list', { params }),
}

export const cartAPI = {
  list: () => api.get<any, ApiResponse<CartItem[]>>('/cart'),
  add: (data: { productId: number; skuId?: number; quantity: number }) =>
    api.post('/cart', data),
  update: (id: number, data: { quantity: number }) => api.put(`/cart/${id}`, data),
  delete: (id: number) => api.delete(`/cart/${id}`),
}

export const orderAPI = {
  list: (params?: { page?: number; pageSize?: number; status?: string }) =>
    api.get<any, ApiResponse<PaginatedResponse<Order[]>>>('/orders', { params }),
  getById: (id: number) => api.get<any, ApiResponse<Order>>(`/orders/${id}`),
  create: (data: {
    cartIds?: number[]
    receiverName: string
    receiverPhone: string
    receiverAddress: string
    remark?: string
  }) => api.post('/orders', data),
  pay: (id: number) => api.post(`/orders/${id}/pay`),
  ship: (id: number, data: { trackingNo: string; trackingCompany: string }) =>
    api.post(`/orders/${id}/ship`, data),
  confirm: (id: number) => api.post(`/orders/${id}/confirm`),
  cancel: (id: number) => api.post(`/orders/${id}/cancel`),
  applyRefund: (data: { orderId: number; orderItemId?: number; reason: string; type: string }) =>
    api.post('/orders/refund', data),
  reviewRefund: (id: number, data: { status: string; rejectReason?: string }) =>
    api.put(`/orders/refund/${id}/review`, data),
}

export const reviewAPI = {
  getByProduct: (productId: number, params?: { page?: number; pageSize?: number }) =>
    api.get<any, ApiResponse<PaginatedResponse<Review[]>>>(`/reviews/product/${productId}`, { params }),
  create: (data: {
    orderId: number
    orderItemId: number
    productId: number
    rating: number
    content: string
    images?: string[]
  }) => api.post('/reviews', data),
  reply: (id: number, data: { reply: string }) => api.put(`/reviews/${id}/reply`, data),
  getShopReviews: (params?: { page?: number; pageSize?: number }) =>
    api.get<any, ApiResponse<PaginatedResponse<Review[]>>>('/reviews/my/shop', { params }),
}

export const favoriteAPI = {
  toggleShop: (shopId: number) => api.post(`/favorites/shop/${shopId}`),
  toggleProduct: (productId: number) => api.post(`/favorites/product/${productId}`),
  getShops: (params?: { page?: number; pageSize?: number }) =>
    api.get<any, ApiResponse<PaginatedResponse<Shop[]>>>('/favorites/shops', { params }),
  getProducts: (params?: { page?: number; pageSize?: number }) =>
    api.get<any, ApiResponse<PaginatedResponse<Product[]>>>('/favorites/products', { params }),
}

export const notificationAPI = {
  list: (params?: { page?: number; pageSize?: number; isRead?: boolean }) =>
    api.get<any, ApiResponse<PaginatedResponse<Notification[]>>>('/notifications', { params }),
  getUnreadCount: () => api.get<any, ApiResponse<{ count: number }>>('/notifications/unread-count'),
  markAsRead: (id: number) => api.put(`/notifications/${id}/read`),
  markAllAsRead: () => api.put('/notifications/read-all'),
}

export const adminAPI = {
  getUsers: (params?: { page?: number; pageSize?: number; role?: string; keyword?: string }) =>
    api.get<any, ApiResponse<PaginatedResponse<User[]>>>('/admin/users', { params }),
  getDisputes: (params?: { page?: number; pageSize?: number; status?: string }) =>
    api.get<any, ApiResponse<PaginatedResponse<any[]>>>('/admin/disputes', { params }),
  resolveDispute: (id: number, data: { result: string }) =>
    api.put(`/admin/disputes/${id}/resolve`, data),
  getStatistics: (params?: { startDate?: string; endDate?: string }) =>
    api.get<any, ApiResponse<AdminStatistics>>('/admin/statistics', { params }),
  exportOrders: (params?: { startDate?: string; endDate?: string }) =>
    api.get('/admin/export/orders', { params, responseType: 'blob' }),
}

export const statisticsAPI = {
  getShopStatistics: (params?: { startDate?: string; endDate?: string }) =>
    api.get<any, ApiResponse<ShopStatistics>>('/statistics/shop', { params }),
  exportShopOrders: (params?: { startDate?: string; endDate?: string }) =>
    api.get('/statistics/shop/export', { params, responseType: 'blob' }),
}

export const uploadAPI = {
  uploadImage: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return api.post<any, { url: string }>('/upload/image', formData)
  },
  uploadMultiple: (files: File[]) => {
    const formData = new FormData()
    files.forEach((file) => formData.append('files', file))
    return api.post<any, { urls: string[] }>('/upload/multiple', formData)
  },
}

export default api
