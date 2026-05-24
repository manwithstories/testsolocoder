import { apiRequest } from './index';
import { Message, Payment, Wallet, Schedule, TeacherStats, MonthlyReport, PaginatedData } from '@/types';

export const messageApi = {
  send: (data: { receiver_id: string; type: string; content?: string; file_url?: string; file_name?: string; file_size?: number }) =>
    apiRequest.post<Message>('/messages', data),

  getConversation: (userId: string, params?: { page?: number; page_size?: number }) =>
    apiRequest.get<PaginatedData<Message>>(`/messages/${userId}`, { params }),

  getConversations: () => apiRequest.get<Message[]>('/messages/conversations'),

  getUnreadCount: () => apiRequest.get<{ unread_count: number }>('/messages/unread'),

  markAsRead: (userId: string) => apiRequest.put(`/messages/${userId}/read`),
};

export const paymentApi = {
  create: (data: { booking_id: string; method: string }) =>
    apiRequest.post<Payment>('/payments', data),

  list: (params?: { page?: number; page_size?: number }) =>
    apiRequest.get<PaginatedData<Payment>>('/payments', { params }),

  get: (id: string) => apiRequest.get<Payment>(`/payments/${id}`),

  getWallet: () => apiRequest.get<Wallet>('/payments/wallet'),

  withdraw: (amount: number) =>
    apiRequest.post('/payments/withdraw', { amount }),
};

export const scheduleApi = {
  create: (data: any) => apiRequest.post<Schedule>('/schedules', data),

  list: () => apiRequest.get<Schedule[]>('/schedules'),

  update: (id: string, data: any) =>
    apiRequest.put<Schedule>(`/schedules/${id}`, data),

  delete: (id: string) => apiRequest.delete(`/schedules/${id}`),

  getAvailability: (userId: string, params?: { day_of_week?: string }) =>
    apiRequest.get<Schedule[]>(`/schedules/availability/${userId}`, { params }),
};

export const statsApi = {
  getTeacherStats: (params?: { start_date?: string; end_date?: string }) =>
    apiRequest.get<TeacherStats>('/stats/teacher', { params }),

  getMonthlyReport: (params?: { year?: number; month?: number }) =>
    apiRequest.get<MonthlyReport>('/stats/monthly', { params }),

  exportReport: () => apiRequest.get('/stats/export'),
};
