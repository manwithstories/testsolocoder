import axios from 'axios';
import type {
  User,
  Client,
  Project,
  TimeEntry,
  Invoice,
  DashboardStats,
  AuthResponse,
  TokenPair,
  ApiResponse,
  PaginatedResponse,
} from '../types';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      try {
        const refreshToken = localStorage.getItem('refresh_token');
        if (refreshToken) {
          const response = await api.post<ApiResponse<TokenPair>>('/auth/refresh', {
            refresh_token: refreshToken,
          });
          if (response.data.success && response.data.data) {
            localStorage.setItem('access_token', response.data.data.access_token);
            localStorage.setItem('refresh_token', response.data.data.refresh_token);
            api.defaults.headers.common['Authorization'] = `Bearer ${response.data.data.access_token}`;
            return api(originalRequest);
          }
        }
      } catch (refreshError) {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export const authAPI = {
  register: (data: { email: string; password: string; first_name: string; last_name: string }) =>
    api.post<ApiResponse<AuthResponse>>('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post<ApiResponse<AuthResponse>>('/auth/login', data),
  refresh: (refresh_token: string) =>
    api.post<ApiResponse<TokenPair>>('/auth/refresh', { refresh_token }),
  me: () => api.get<ApiResponse<User>>('/auth/me'),
};

export const clientAPI = {
  create: (data: Partial<Client>) => api.post<ApiResponse<Client>>('/clients', data),
  list: (params?: { page?: number; per_page?: number }) =>
    api.get<PaginatedResponse<Client[]>>('/clients', { params }),
  get: (id: number) => api.get<ApiResponse<Client>>(`/clients/${id}`),
  update: (id: number, data: Partial<Client>) =>
    api.put<ApiResponse<Client>>(`/clients/${id}`, data),
  delete: (id: number) => api.delete(`/clients/${id}`),
};

export const projectAPI = {
  create: (data: any) => api.post<ApiResponse<Project>>('/projects', data),
  list: (params?: { page?: number; per_page?: number; status?: string; client_id?: number }) =>
    api.get<PaginatedResponse<Project[]>>('/projects', { params }),
  get: (id: number) => api.get<ApiResponse<Project>>(`/projects/${id}`),
  update: (id: number, data: any) => api.put<ApiResponse<Project>>(`/projects/${id}`, data),
  delete: (id: number) => api.delete(`/projects/${id}`),
  addMilestone: (projectId: number, data: any) =>
    api.post<ApiResponse<any>>(`/projects/${projectId}/milestones`, data),
  updateMilestone: (milestoneId: number, data: any) =>
    api.put<ApiResponse<any>>(`/projects/milestones/${milestoneId}`, data),
  deleteMilestone: (milestoneId: number) =>
    api.delete(`/projects/milestones/${milestoneId}`),
};

export const timeEntryAPI = {
  create: (data: any) => api.post<ApiResponse<TimeEntry>>('/time-entries', data),
  list: (params?: { page?: number; per_page?: number; project_id?: number; start_date?: string; end_date?: string }) =>
    api.get<PaginatedResponse<TimeEntry[]>>('/time-entries', { params }),
  getActiveTimer: () => api.get<ApiResponse<TimeEntry>>('/time-entries/active-timer'),
  get: (id: number) => api.get<ApiResponse<TimeEntry>>(`/time-entries/${id}`),
  update: (id: number, data: any) =>
    api.put<ApiResponse<TimeEntry>>(`/time-entries/${id}`, data),
  delete: (id: number) => api.delete(`/time-entries/${id}`),
  startTimer: (data: { project_id: number; description?: string }) =>
    api.post<ApiResponse<TimeEntry>>('/time-entries/timer/start', data),
  stopTimer: (id: number, data?: { description?: string }) =>
    api.post<ApiResponse<TimeEntry>>(`/time-entries/timer/${id}/stop`, data),
};

export const invoiceAPI = {
  create: (data: any) => api.post<ApiResponse<Invoice>>('/invoices', data),
  list: (params?: { page?: number; per_page?: number; status?: string; client_id?: number; year?: number }) =>
    api.get<PaginatedResponse<Invoice[]>>('/invoices', { params }),
  get: (id: number) => api.get<ApiResponse<Invoice>>(`/invoices/${id}`),
  updateStatus: (id: number, status: string) =>
    api.put<ApiResponse<Invoice>>(`/invoices/${id}/status`, { status }),
  downloadPDF: (id: number) =>
    api.get(`/invoices/${id}/download`, { responseType: 'blob' }),
  delete: (id: number) => api.delete(`/invoices/${id}`),
};

export const dashboardAPI = {
  getStats: () => api.get<ApiResponse<DashboardStats>>('/dashboard'),
};

export default api;
