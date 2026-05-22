import axios, { AxiosInstance, AxiosResponse } from 'axios'
import type {
  User,
  AuthResponse,
  TravelPlan,
  PlanDetail,
  Activity,
  File,
  Checklist,
  ChecklistItem,
  Reminder,
  BudgetSummary,
  MapLocation,
  ApiResponse,
  PaginatedResponse,
  PaginationParams,
} from '@/types'

const api: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
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

api.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    if (response.data.code !== 0) {
      return Promise.reject(new Error(response.data.message || 'Request failed'))
    }
    return response.data.data as any
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

export const authAPI = {
  register: (data: { username: string; email: string; password: string; first_name?: string; last_name?: string }) =>
    api.post('/auth/register', data),
  login: (data: { email: string; password: string }) =>
    api.post<AuthResponse>('/auth/login', data) as Promise<AuthResponse>,
  getCurrentUser: () => api.get<User>('/auth/me') as Promise<User>,
  updateProfile: (data: Partial<User>) => api.put<User>('/auth/profile', data) as Promise<User>,
}

export const planAPI = {
  getPlans: (params?: PaginationParams & { status?: string }) =>
    api.get<PaginatedResponse<TravelPlan>>('/plans', { params }) as Promise<PaginatedResponse<TravelPlan>>,
  getPlan: (id: string) => api.get<PlanDetail>(`/plans/${id}`) as Promise<PlanDetail>,
  createPlan: (data: {
    title: string
    description?: string
    destination: string
    start_date: string
    end_date: string
    budget?: number
    currency?: string
    cover_image?: string
    is_public?: boolean
  }) => api.post<TravelPlan>('/plans', data) as Promise<TravelPlan>,
  updatePlan: (id: string, data: Partial<TravelPlan>) =>
    api.put<TravelPlan>(`/plans/${id}`, data) as Promise<TravelPlan>,
  deletePlan: (id: string) => api.delete(`/plans/${id}`),
  addParticipant: (planId: string, data: { user_id: string; role: string; can_edit: boolean; can_delete: boolean }) =>
    api.post(`/plans/${planId}/participants`, data),
  removeParticipant: (planId: string, participantId: string) =>
    api.delete(`/plans/${planId}/participants/${participantId}`),
  exportJSON: (id: string) => api.get(`/plans/${id}/export/json`, { responseType: 'blob' }),
  exportPDF: (id: string) => api.get(`/plans/${id}/export/pdf`, { responseType: 'blob' }),
}

export const activityAPI = {
  getActivities: (planId: string, params?: { date?: string; type?: string }) =>
    api.get(`/plans/${planId}/activities`, { params }),
  getActivity: (planId: string, id: string) =>
    api.get<Activity>(`/plans/${planId}/activities/${id}`) as Promise<Activity>,
  createActivity: (planId: string, data: {
    title: string
    description?: string
    type: string
    date: string
    start_time?: string
    end_time?: string
    location?: string
    latitude?: number
    longitude?: number
    cost?: number
    currency?: string
    notes?: string
    booked?: boolean
    confirmation?: string
    contact_info?: string
    order_index?: number
  }) => api.post<Activity>(`/plans/${planId}/activities`, data) as Promise<Activity>,
  updateActivity: (planId: string, id: string, data: Partial<Activity>) =>
    api.put<Activity>(`/plans/${planId}/activities/${id}`, data) as Promise<Activity>,
  deleteActivity: (planId: string, id: string) =>
    api.delete(`/plans/${planId}/activities/${id}`),
  getBudgetSummary: (planId: string) =>
    api.get<BudgetSummary>(`/plans/${planId}/budget`) as Promise<BudgetSummary>,
}

export const fileAPI = {
  getFiles: (planId: string, params?: { category?: string }) =>
    api.get<File[]>(`/plans/${planId}/files`, { params }) as Promise<File[]>,
  uploadFile: (planId: string, formData: FormData, onProgress?: (progress: number) => void) =>
    api.post(`/plans/${planId}/files`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(percentCompleted)
        }
      },
    }),
  deleteFile: (planId: string, fileId: string) =>
    api.delete(`/plans/${planId}/files/${fileId}`),
  downloadFile: (fileId: string) =>
    api.get(`/files/download/${fileId}`, { responseType: 'blob' }),
}

export const checklistAPI = {
  getChecklists: (planId: string) =>
    api.get<Checklist[]>(`/plans/${planId}/checklists`) as Promise<Checklist[]>,
  getChecklist: (planId: string, id: string) =>
    api.get<Checklist>(`/plans/${planId}/checklists/${id}`) as Promise<Checklist>,
  createChecklist: (planId: string, data: { title: string; type: string }) =>
    api.post<Checklist>(`/plans/${planId}/checklists`, data) as Promise<Checklist>,
  updateChecklist: (planId: string, id: string, data: Partial<Checklist>) =>
    api.put<Checklist>(`/plans/${planId}/checklists/${id}`, data) as Promise<Checklist>,
  deleteChecklist: (planId: string, id: string) =>
    api.delete(`/plans/${planId}/checklists/${id}`),
  addItem: (planId: string, checklistId: string, data: {
    title: string
    description?: string
    category?: string
    quantity?: number
    order_index?: number
  }) => api.post<ChecklistItem>(`/plans/${planId}/checklists/${checklistId}/items`, data) as Promise<ChecklistItem>,
  updateItem: (planId: string, checklistId: string, itemId: string, data: Partial<ChecklistItem>) =>
    api.put<ChecklistItem>(`/plans/${planId}/checklists/${checklistId}/items/${itemId}`, data) as Promise<ChecklistItem>,
  deleteItem: (planId: string, checklistId: string, itemId: string) =>
    api.delete(`/plans/${planId}/checklists/${checklistId}/items/${itemId}`),
}

export const reminderAPI = {
  getReminders: () => api.get<Reminder[]>('/reminders') as Promise<Reminder[]>,
  createReminder: (planId: string, data: {
    title: string
    description?: string
    reminder_time: string
    channel: string
    activity_id?: string
  }) => api.post<Reminder>(`/plans/${planId}/reminders`, data) as Promise<Reminder>,
  updateReminder: (id: string, data: Partial<Reminder>) =>
    api.put<Reminder>(`/reminders/${id}`, data) as Promise<Reminder>,
  deleteReminder: (id: string) => api.delete(`/reminders/${id}`),
}

export const mapAPI = {
  getMapData: (planId: string) =>
    api.get<{ locations: MapLocation[]; total: number }>(`/plans/${planId}/map`) as Promise<{ locations: MapLocation[]; total: number }>,
}

export default api
