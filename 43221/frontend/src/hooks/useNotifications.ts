import { useQuery } from '@tanstack/react-query'
import { notificationApi } from '@/services/notification'

export function useNotifications(params?: { page?: number; page_size?: number; only_unread?: boolean }) {
  return useQuery({
    queryKey: ['notifications', params],
    queryFn: () => notificationApi.getNotifications(params),
  })
}

export function useUnreadCount() {
  return useQuery({
    queryKey: ['unread-count'],
    queryFn: () => notificationApi.getUnreadCount(),
    refetchInterval: 30000,
  })
}
