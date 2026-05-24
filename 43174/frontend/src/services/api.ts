import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

const API_BASE_URL = '/api/v1';

const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error.response?.data || error.message);
  }
);

export const api = {
  get: <T>(url: string, config?: AxiosRequestConfig) =>
    apiClient.get<unknown, T>(url, config),

  post: <T>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
    apiClient.post<unknown, T>(url, data, config),

  put: <T>(url: string, data?: unknown, config?: AxiosRequestConfig) =>
    apiClient.put<unknown, T>(url, data, config),

  delete: <T>(url: string, config?: AxiosRequestConfig) =>
    apiClient.delete<unknown, T>(url, config),

  upload: <T>(url: string, formData: FormData) =>
    apiClient.post<unknown, T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
};

export const authApi = {
  login: (data: { username: string; password: string }) =>
    api.post<{ token: string; user: any }>('/auth/login', data),

  register: (data: any) =>
    api.post('/auth/register', data),
};

export const userApi = {
  getProfile: () => api.get('/users/profile'),

  updateProfile: (data: any) => api.put('/users/profile', data),

  changePassword: (data: { old_password: string; new_password: string }) =>
    api.put('/users/password', data),

  uploadAvatar: (formData: FormData) => api.upload('/users/avatar', formData),

  getUsers: (params?: any) => api.get('/users', { params }),

  getUserById: (id: string) => api.get(`/users/${id}`),

  updateUserStatus: (id: string, status: string) =>
    api.put(`/users/${id}/status`, { status }),

  deleteUser: (id: string) => api.delete(`/users/${id}`),

  getTopRated: (limit?: number) => api.get('/users/top-rated', { params: { limit } }),
};

export const categoryApi = {
  getAll: () => api.get('/categories'),

  getById: (id: string) => api.get(`/categories/${id}`),

  create: (data: any) => api.post('/categories', data),

  update: (id: string, data: any) => api.put(`/categories/${id}`, data),

  delete: (id: string) => api.delete(`/categories/${id}`),
};

export const textbookApi = {
  getAll: (params?: any) => api.get('/textbooks', { params }),

  getById: (id: string) => api.get(`/textbooks/${id}`),

  searchByISBN: (isbn: string) => api.get('/textbooks/search/isbn', { params: { isbn } }),

  getPopular: (limit?: number) => api.get('/textbooks/popular', { params: { limit } }),

  getBySeller: (sellerId: string, params?: any) =>
    api.get(`/textbooks/seller/${sellerId}`, { params }),

  create: (data: any) => api.post('/textbooks', data),

  uploadCoverImage: (formData: FormData) => api.upload('/textbooks/cover-image', formData),

  update: (id: string, data: any) => api.put(`/textbooks/${id}`, data),

  updateStatus: (id: string, status: string) =>
    api.put(`/textbooks/${id}/status`, { status }),

  delete: (id: string) => api.delete(`/textbooks/${id}`),
};

export const noteApi = {
  getAll: (params?: any) => api.get('/notes', { params }),

  getById: (id: string) => api.get(`/notes/${id}`),

  getFeatured: (limit?: number) => api.get('/notes/featured', { params: { limit } }),

  getByUploader: (uploaderId: string, params?: any) =>
    api.get(`/notes/uploader/${uploaderId}`, { params }),

  create: (data: any) => api.post('/notes', data),

  uploadFile: (formData: FormData) => api.upload('/notes/upload', formData),

  update: (id: string, data: any) => api.put(`/notes/${id}`, data),

  delete: (id: string) => api.delete(`/notes/${id}`),

  incrementDownload: (id: string) => api.post(`/notes/${id}/download`),

  setFeatured: (id: string, isFeatured: boolean) =>
    api.put(`/notes/${id}/featured`, { is_featured: isFeatured }),
};

export const transactionApi = {
  create: (data: any) => api.post('/transactions', data),

  getAll: (params?: any) => api.get('/transactions', { params }),

  getById: (id: string) => api.get(`/transactions/${id}`),

  getByBuyer: (buyerId: string, params?: any) =>
    api.get(`/transactions/buyer/${buyerId}`, { params }),

  getBySeller: (sellerId: string, params?: any) =>
    api.get(`/transactions/seller/${sellerId}`, { params }),

  confirm: (id: string) => api.put(`/transactions/${id}/confirm`),

  complete: (id: string) => api.put(`/transactions/${id}/complete`),

  cancel: (id: string, reason?: string) =>
    api.put(`/transactions/${id}/cancel`, { reason }),

  negotiate: (id: string, price: number) =>
    api.put(`/transactions/${id}/negotiate`, { price }),
};

export const orderApi = {
  create: (data: any) => api.post('/orders', data),

  getAll: (params?: any) => api.get('/orders', { params }),

  getMyOrders: (params?: any) => api.get('/orders/my', { params }),

  getById: (id: string) => api.get(`/orders/${id}`),

  getByOrderNo: (orderNo: string) => api.get(`/orders/order-no/${orderNo}`),

  pay: (id: string) => api.put(`/orders/${id}/pay`),

  ship: (id: string, trackingNumber: string) =>
    api.put(`/orders/${id}/ship`, { tracking_number: trackingNumber }),

  deliver: (id: string) => api.put(`/orders/${id}/deliver`),

  complete: (id: string) => api.put(`/orders/${id}/complete`),

  cancel: (id: string, reason?: string) =>
    api.put(`/orders/${id}/cancel`, { reason }),

  updateStatus: (id: string, status: string, remark?: string) =>
    api.put(`/orders/${id}/status`, { status, remark }),
};

export const messageApi = {
  create: (data: any) => api.post('/messages', data),

  getConversation: (userId1: string, userId2: string, params?: any) =>
    api.get('/messages/conversation', { params: { user_id_1: userId1, user_id_2: userId2, ...params } }),

  getUnreadCount: () => api.get('/messages/unread-count'),

  markAsRead: () => api.put('/messages/mark-read'),
};

export const reviewApi = {
  create: (data: any) => api.post('/reviews', data),

  getByTextbook: (textbookId: string, params?: any) =>
    api.get(`/reviews/textbook/${textbookId}`, { params }),

  getByNote: (noteId: string, params?: any) =>
    api.get(`/reviews/note/${noteId}`, { params }),

  getAll: (params?: any) => api.get('/reviews', { params }),

  hide: (id: string) => api.put(`/reviews/${id}/hide`),

  markMalicious: (id: string, isMalicious: boolean) =>
    api.put(`/reviews/${id}/malicious`, { is_malicious: isMalicious }),
};

export const notificationApi = {
  getAll: (params?: any) => api.get('/notifications', { params }),

  markAsRead: (id: string) => api.put(`/notifications/${id}/read`),

  markAllAsRead: () => api.put('/notifications/read-all'),
};

export const statisticsApi = {
  getTextbookStats: () => api.get('/statistics/textbooks'),

  getUserStats: () => api.get('/statistics/users'),

  getOrderStats: () => api.get('/statistics/orders'),

  getPopularTextbooks: (limit?: number) =>
    api.get('/statistics/popular-textbooks', { params: { limit } }),

  getTopUsers: (limit?: number) => api.get('/statistics/top-users', { params: { limit } }),

  getMonthlyStats: (months?: number) =>
    api.get('/statistics/monthly', { params: { months } }),

  exportReport: (month: string) =>
    api.get('/statistics/export', { params: { month }, responseType: 'blob' }),
};
