import request from './request'
import { Message, PaginatedResponse } from '@/types'

export const messageApi = {
  getList: (params?: {
    page?: number
    page_size?: number
    type?: string
    is_read?: string
  }) => {
    return request.get<any, PaginatedResponse<Message> & { unread_count: number }>('/messages', { params })
  },

  getUnreadCount: () => {
    return request.get<any, { unread_count: number }>('/messages/unread-count')
  },

  read: (id: number) => {
    return request.put<any, any>(`/messages/${id}/read`)
  },

  readAll: () => {
    return request.put<any, any>('/messages/read-all')
  },

  delete: (id: number) => {
    return request.delete<any, any>(`/messages/${id}`)
  },
}
