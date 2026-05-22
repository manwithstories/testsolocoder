import request from '@/utils/request'
import type { ApiResponse, Notification, PageResult, PageParams, AuditLog } from '@/types'

export const getNotifications = (params: PageParams) => {
  return request.get<ApiResponse<PageResult<Notification>>>('/notifications', { params })
}

export const getUnreadCount = () => {
  return request.get<ApiResponse<{ unread_count: number }>>('/notifications/unread-count')
}

export const markAsRead = (id: number) => {
  return request.put<ApiResponse<null>>(`/notifications/${id}/read`)
}

export const markAllAsRead = () => {
  return request.put<ApiResponse<null>>('/notifications/read-all')
}

export const getAuditLogs = (params: PageParams & { user_id?: number; module?: string; action?: string }) => {
  return request.get<ApiResponse<PageResult<AuditLog>>>('/audits', { params })
}
