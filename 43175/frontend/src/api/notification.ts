import request from './request'

export interface Notification {
  id: number
  userId: number
  type: string
  title: string
  content: string
  isRead: boolean
  createdAt: string
}

export function listNotifications(params?: { isRead?: string; type?: string }): Promise<any> {
  return request.get('/notifications', { params })
}

export function markNotificationRead(id: number): Promise<void> {
  return request.put(`/notifications/${id}/read`)
}

export function markAllRead(): Promise<void> {
  return request.put('/notifications/read-all')
}

export function deleteNotification(id: number): Promise<void> {
  return request.delete(`/notifications/${id}`)
}

export const notificationTypeMap: Record<string, { label: string; color: string; icon: string }> = {
  alert: { label: '告警', color: 'danger', icon: 'Warning' },
  info: { label: '信息', color: 'info', icon: 'InfoFilled' },
  success: { label: '成功', color: 'success', icon: 'CircleCheck' },
  invitation: { label: '邀请', color: 'warning', icon: 'Bell' }
}
