import { useState, useEffect } from 'react'
import { notificationApi } from '@/services/api'
import { Notification } from '@/types'
import { Bell, X } from 'lucide-react'

interface NotificationDropdownProps {
  onClose: () => void
}

export default function NotificationDropdown({ onClose }: NotificationDropdownProps) {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadNotifications()
  }, [])

  const loadNotifications = async () => {
    try {
      setLoading(true)
      const res = await notificationApi.getAll()
      setNotifications(res.data)
    } catch (error) {
      console.error('Failed to load notifications:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleMarkAllRead = async () => {
    try {
      await notificationApi.markAllRead()
      setNotifications(notifications.map(n => ({ ...n, isRead: true })))
    } catch (error) {
      console.error('Failed to mark notifications as read:', error)
    }
  }

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'booking_created':
      case 'booking_updated':
        return '📅'
      case 'booking_cancelled':
        return '❌'
      case 'lesson_start':
      case 'lesson_end':
        return '📹'
      case 'new_review':
        return '⭐'
      case 'payment_received':
        return '💰'
      case 'new_message':
        return '💬'
      case 'homework_assigned':
      case 'homework_graded':
        return '📝'
      case 'milestone':
        return '🏆'
      default:
        return '🔔'
    }
  }

  return (
    <div className="absolute right-0 mt-2 w-96 bg-white rounded-xl shadow-lg border border-gray-200 z-50">
      <div className="p-4 border-b border-gray-200 flex items-center justify-between">
        <h3 className="font-semibold text-gray-900">通知</h3>
        <button
          onClick={onClose}
          className="p-1 hover:bg-gray-100 rounded"
        >
          <X className="h-4 w-4" />
        </button>
      </div>

      <div className="max-h-96 overflow-y-auto">
        {loading ? (
          <div className="p-8 text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
          </div>
        ) : notifications.length > 0 ? (
          <div className="divide-y divide-gray-100">
            {notifications.map((notification) => (
              <div
                key={notification.id}
                className={`p-4 hover:bg-gray-50 cursor-pointer ${
                  !notification.isRead ? 'bg-primary-50' : ''
                }`}
              >
                <div className="flex items-start gap-3">
                  <span className="text-xl">{getNotificationIcon(notification.type)}</span>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-gray-900 text-sm">
                      {notification.title}
                    </div>
                    <p className="text-gray-500 text-sm truncate">
                      {notification.content}
                    </p>
                    <span className="text-xs text-gray-400">
                      {new Date(notification.createdAt).toLocaleDateString('zh-CN', {
                        month: 'short',
                        day: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </span>
                  </div>
                  {!notification.isRead && (
                    <div className="w-2 h-2 bg-primary-600 rounded-full"></div>
                  )}
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="p-8 text-center">
            <Bell className="h-12 w-12 text-gray-300 mx-auto mb-2" />
            <p className="text-gray-500">暂无通知</p>
          </div>
        )}
      </div>

      {notifications.some(n => !n.isRead) && (
        <div className="p-4 border-t border-gray-200">
          <button
            onClick={handleMarkAllRead}
            className="w-full text-primary-600 hover:text-primary-700 text-sm font-medium"
          >
            全部标为已读
          </button>
        </div>
      )}
    </div>
  )
}
