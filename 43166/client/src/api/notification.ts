import request from '@/utils/request'
import {
  Notification,
  NotificationTemplate,
  SendNotificationRequest,
  PaginatedData
} from '@/types'

export const notificationApi = {
  getList: (params?: {
    page?: number
    pageSize?: number
    isRead?: boolean
  }) => {
    return request.get<any, PaginatedData<Notification>>('/notifications', { params })
  },

  markAsRead: (id: number) => {
    return request.put<any, null>(`/notifications/${id}/read`)
  },

  markAllAsRead: () => {
    return request.put<any, null>('/notifications/read-all')
  },

  getUnreadCount: () => {
    return request.get<any, { count: number }>('/notifications/unread-count')
  },

  send: (data: SendNotificationRequest) => {
    return request.post<any, null>('/notifications/send', data)
  },

  getTemplates: () => {
    return request.get<any, NotificationTemplate[]>('/admin/notification-templates')
  },

  createTemplate: (data: Partial<NotificationTemplate>) => {
    return request.post<any, NotificationTemplate>('/admin/notification-templates', data)
  },

  updateTemplate: (id: number, data: Partial<NotificationTemplate>) => {
    return request.put<any, null>(`/admin/notification-templates/${id}`, data)
  },

  deleteTemplate: (id: number) => {
    return request.delete<any, null>(`/admin/notification-templates/${id}`)
  }
}
