import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response;
  },
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const response = await axios.post('/api/v1/auth/refresh', {
            refresh_token: refreshToken,
          });
          const { access_token, refresh_token } = response.data;
          localStorage.setItem('access_token', access_token);
          localStorage.setItem('refresh_token', refresh_token);
          if (error.config) {
            error.config.headers.Authorization = `Bearer ${access_token}`;
            return axios(error.config);
          }
        } catch (refreshError) {
          localStorage.removeItem('access_token');
          localStorage.removeItem('refresh_token');
          localStorage.removeItem('user');
          window.location.href = '/login';
        }
      } else {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const authApi = {
  register: (data: any) => api.post('/auth/register', data),
  login: (data: any) => api.post('/auth/login', data),
  refreshToken: (refreshToken: string) => api.post('/auth/refresh', { refresh_token: refreshToken }),
  getProfile: () => api.get('/user/profile'),
  updateProfile: (data: any) => api.put('/user/profile', data),
  updateDesignerProfile: (data: any) => api.put('/user/profile/designer', data),
  updatePrinterProfile: (data: any) => api.put('/user/profile/printer', data),
  getUserStats: () => api.get('/user/stats'),
  getNotifications: (params?: any) => api.get('/user/notifications', { params }),
  markNotificationRead: (id: string) => api.put(`/user/notifications/${id}/read`),
  getTransactions: (params?: any) => api.get('/user/transactions', { params }),
};

export const modelApi = {
  list: (params?: any) => api.get('/models', { params }),
  getHot: (limit: number = 10) => api.get('/models/hot', { params: { limit } }),
  get: (id: string) => api.get(`/models/${id}`),
  create: (data: any) => api.post('/models', data),
  update: (id: string, data: any) => api.put(`/models/${id}`, data),
  delete: (id: string) => api.delete(`/models/${id}`),
  uploadFile: (id: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post(`/models/${id}/file`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
  uploadThumbnail: (id: string, file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post(`/models/${id}/thumbnail`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
  purchase: (id: string, data: any) => api.post(`/models/${id}/purchase', data),
  download: (id: string) => api.get(`/models/${id}/download`),
  addFavorite: (id: string) => api.post(`/models/${id}/favorite'),
  removeFavorite: (id: string) => api.delete(`/models/${id}/favorite'),
  getVersions: (id: string) => api.get(`/models/${id}/versions`),
  createVersion: (id: string, formData: FormData) => api.post(`/models/${id}/version', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
  getDesignerModels: (designerId: string, params?: any) =>
    api.get(`/designers/${designerId}/models`, { params }),
  getMyModels: (params?: any) => api.get('/models/designer/my', { params }),
  getFavorites: (params?: any) => api.get('/favorites', { params }),
  getPurchases: (params?: any) => api.get('/purchases', { params }),
  validateFile: (file: File) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post('/models/validate', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};

export const orderApi = {
  estimate: (data: any) => api.post('/orders/estimate', data),
  create: (data: any) => api.post('/orders', data),
  get: (id: string) => api.get(`/orders/${id}`),
  getByNo: (orderNo: string) => api.get(`/orders/no/${orderNo}`),
  getHistory: (id: string) => api.get(`/orders/${id}/history`),
  listCustomerOrders: (params?: any) => api.get('/orders/my/customer', { params }),
  listPrinterOrders: (params?: any) => api.get('/orders/my/printer', { params }),
  getPending: () => api.get('/orders/pending'),
  assignPrinter: (id: string) => api.put(`/orders/${id}/assign`),
  startPrinting: (id: string) => api.put(`/orders/${id}/print/start'),
  completePrinting: (id: string) => api.put(`/orders/${id}/print/complete'),
  approveQuality: (id: string) => api.put(`/orders/${id}/quality/approve'),
  shipOrder: (id: string, data: any) => api.put(`/orders/${id}/ship`, data),
  deliverOrder: (id: string) => api.put(`/orders/${id}/deliver`),
  completeOrder: (id: string) => api.put(`/orders/${id}/complete`),
  cancelOrder: (id: string, data: any) => api.put(`/orders/${id}/cancel`, data),
};

export const printerApi = {
  getMaterials: () => api.get('/materials'),
  getMaterial: (id: string) => api.get(`/materials/${id}`),
  getDevices: () => api.get('/printer/devices'),
  createDevice: (data: any) => api.post('/printer/devices', data),
  updateDevice: (id: string, data: any) => api.put(`/printer/devices/${id}`, data),
  deleteDevice: (id: string) => api.delete(`/printer/devices/${id}`),
  getIdleDevices: () => api.get('/printer/devices/idle'),
  getInventory: () => api.get('/printer/inventory'),
  createInventory: (data: any) => api.post('/printer/inventory', data),
  updateInventory: (id: string, data: any) => api.put(`/printer/inventory/${id}`, data),
  deleteInventory: (id: string) => api.delete(`/printer/inventory/${id}`),
  getSchedules: (params?: any) => api.get('/printer/schedules', { params }),
  createSchedule: (data: any) => api.post('/printer/schedules', data),
  createReview: (data: any) => api.post('/reviews', data),
  getReview: (id: string) => api.get(`/reviews/${id}`),
  getModelReviews: (modelId: string, params?: any) => api.get(`/reviews/model/${modelId}`, { params }),
  getPrinterReviews: (printerId: string, params?: any) => api.get(`/reviews/printer/${printerId}`, { params }),
  listDesigners: (params?: any) => api.get('/designers', { params }),
  listPrinters: (params?: any) => api.get('/printers', { params }),
};

export const fileApi = {
  initiateUpload: (data: any) => api.post('/files/initiate', data),
  uploadChunk: (uploadId: string, chunkNumber: number, chunk: Blob) => {
    const formData = new FormData();
    formData.append('chunk', chunk);
    return api.post(`/files/${uploadId}/chunk?chunk=${chunkNumber}`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
  completeUpload: (uploadId: string) => api.post(`/files/${uploadId}/complete`),
  getUpload: (id: string) => api.get(`/files/${id}`),
  getMyUploads: (params?: any) => api.get('/files/my', { params }),
  deleteUpload: (id: string) => api.delete(`/files/${id}`),
  getAccessLogs: (id: string, params?: any) => api.get(`/files/${id}/logs`, { params }),
};

export const statsApi = {
  getPlatformStats: (params?: any) => api.get('/stats/platform', { params }),
  getRevenueStats: (params?: any) => api.get('/stats/revenue', { params }),
  getMaterialStats: (params?: any) => api.get('/stats/materials', { params }),
  exportStats: (params?: any) => api.get('/stats/export', { params, responseType: 'blob' }),
};

export default api;
