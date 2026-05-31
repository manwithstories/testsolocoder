import { Card, List, Tag, Button, Badge, Empty, Checkbox } from 'antd'
import {
  CheckOutlined,
  BellOutlined,
  CalendarOutlined,
  DollarOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuthContext } from '@/contexts/AuthContext'
import { notificationApi } from '@/services/notification'
import { Notification, NotificationType } from '@/types'

export function Notifications() {
  const { user } = useAuthContext()
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['notifications'],
    queryFn: () => notificationApi.getNotifications({ page: 1, page_size: 50 }),
    enabled: !!user,
  })

  const { data: unreadData } = useQuery({
    queryKey: ['unread-count'],
    queryFn: () => notificationApi.getUnreadCount(),
    enabled: !!user,
  })

  const markAsReadMutation = useMutation({
    mutationFn: (id: string) => notificationApi.markAsRead(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
      queryClient.invalidateQueries({ queryKey: ['unread-count'] })
    },
  })

  const markAllAsReadMutation = useMutation({
    mutationFn: () => notificationApi.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] })
      queryClient.invalidateQueries({ queryKey: ['unread-count'] })
    },
  })

  const getNotificationIcon = (type: NotificationType) => {
    switch (type) {
      case 'appointment_success':
      case 'appointment_remind':
        return <CalendarOutlined style={{ color: '#1890ff' }} />
      case 'appointment_cancel':
        return <CalendarOutlined style={{ color: '#ff4d4f' }} />
      case 'payment_success':
        return <DollarOutlined style={{ color: '#52c41a' }} />
      case 'payment_refund':
        return <DollarOutlined style={{ color: '#ff4d4f' }} />
      default:
        return <BellOutlined style={{ color: '#1890ff' }} />
    }
  }

  const getNotificationTypeText = (type: NotificationType) => {
    const map: Record<NotificationType, string> = {
      appointment_success: '预约成功',
      appointment_cancel: '预约取消',
      appointment_remind: '预约提醒',
      payment_success: '支付成功',
      payment_refund: '退款通知',
      review_reply: '评价回复',
      system: '系统通知',
    }
    return map[type] || '通知'
  }

  return (
    <div className="page-container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
        <h2>
          消息通知
          {unreadData?.unread_count ? (
            <Badge
              count={unreadData.unread_count}
              style={{ marginLeft: 8 }}
              className="notification-badge"
            />
          ) : null}
        </h2>
        {unreadData?.unread_count ? (
          <Button
            icon={<CheckOutlined />}
            onClick={() => markAllAsReadMutation.mutate()}
            loading={markAllAsReadMutation.isPending}
          >
            全部已读
          </Button>
        ) : null}
      </div>

      <Card>
        {data?.items && data.items.length > 0 ? (
          <List
            dataSource={data.items}
            loading={isLoading}
            renderItem={(notification: Notification) => (
              <List.Item
                style={{
                  background: notification.is_read ? 'transparent' : '#f6ffed',
                  padding: '16px',
                  borderRadius: 4,
                  marginBottom: 8,
                }}
                actions={[
                  !notification.is_read && (
                    <Button
                      key="read"
                      type="link"
                      size="small"
                      onClick={() => markAsReadMutation.mutate(notification.id)}
                    >
                      标记已读
                    </Button>
                  ),
                ]}
              >
                <List.Item.Meta
                  avatar={
                    <div
                      style={{
                        width: 40,
                        height: 40,
                        borderRadius: '50%',
                        background: '#e6f7ff',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        fontSize: 20,
                      }}
                    >
                      {getNotificationIcon(notification.type)}
                    </div>
                  }
                  title={
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      <span style={{ fontWeight: notification.is_read ? 400 : 600 }}>
                        {notification.title}
                      </span>
                      <Tag color="blue">{getNotificationTypeText(notification.type)}</Tag>
                      {!notification.is_read && <Badge status="processing" />}
                    </div>
                  }
                  description={
                    <div>
                      <div style={{ color: '#666', marginBottom: 4 }}>
                        {notification.content}
                      </div>
                      <div style={{ color: '#999', fontSize: 12 }}>
                        {new Date(notification.created_at).toLocaleString()}
                      </div>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        ) : (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <span>
                <InfoCircleOutlined /> 暂无消息通知
              </span>
            }
          />
        )}
      </Card>
    </div>
  )
}
