import { http } from './request'
import {
  Notification,
  NotificationTemplate,
  PaginatedResponse,
} from '@/types'

export const notificationApi = {
  getNotifications: (params?: { page?: number; page_size?: number; only_unread?: boolean }) =>
    http.get<PaginatedResponse<Notification>>('/notifications', { params }),
  markAsRead: (id: string) =>
    http.put(`/notifications/${id}/read`),
  markAllAsRead: () =>
    http.put('/notifications/read-all'),
  getUnreadCount: () =>
    http.get<{ unread_count: number }>('/notifications/unread-count'),
  getTemplates: (params?: { page?: number; page_size?: number }) =>
    http.get<PaginatedResponse<NotificationTemplate>>('/notifications/templates', { params }),
  updateTemplate: (id: string, data: { title: string; content: string }) =>
    http.put(`/notifications/templates/${id}`, data),
}

export { Notification, NotificationTemplate }
