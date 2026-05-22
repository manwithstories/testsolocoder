import { useState, useEffect } from 'react'
import { List, Card, Tag, Typography, Empty, Button, Pagination, message } from 'antd'
import { BellOutlined, ReadOutlined } from '@ant-design/icons'
import { notificationAPI } from '@/api'
import { Notification } from '@/types'

const { Title, Text } = Typography

const Notifications = () => {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(20)
  const [unreadCount, setUnreadCount] = useState(0)

  useEffect(() => {
    loadNotifications()
    loadUnreadCount()
  }, [page])

  const loadNotifications = async () => {
    try {
      const res = await notificationAPI.list({ page, pageSize })
      setNotifications(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载通知失败', err)
    }
  }

  const loadUnreadCount = async () => {
    try {
      const res = await notificationAPI.getUnreadCount() as any
      setUnreadCount(res.data.count)
    } catch (err) {
      console.error('加载未读数量失败', err)
    }
  }

  const handleMarkAllRead = async () => {
    try {
      await notificationAPI.markAllAsRead()
      message.success('已全部标记为已读')
      setNotifications(notifications.map(n => ({ ...n, isRead: true })))
      setUnreadCount(0)
    } catch (err: any) {
      message.error(err.message || '操作失败')
    }
  }

  const handleMarkRead = async (id: number) => {
    try {
      await notificationAPI.markAsRead(id)
      setNotifications(notifications.map(n => n.id === id ? { ...n, isRead: true } : n))
      setUnreadCount(Math.max(0, unreadCount - 1))
    } catch (err) {
      console.error('标记已读失败', err)
    }
  }

  const getTypeTag = (type: string) => {
    const typeMap: Record<string, { color: string; text: string }> = {
      new_order: { color: 'blue', text: '新订单' },
      order_shipped: { color: 'cyan', text: '已发货' },
      order_completed: { color: 'green', text: '订单完成' },
      refund_approved: { color: 'gold', text: '退款通过' },
      refund_rejected: { color: 'red', text: '退款拒绝' },
      shop_approved: { color: 'green', text: '店铺通过' },
      shop_rejected: { color: 'red', text: '店铺拒绝' },
    }
    const info = typeMap[type] || { color: 'default', text: type }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  return (
    <div>
      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>
          消息通知
          {unreadCount > 0 && <Tag color="red" style={{ marginLeft: 12 }}>{unreadCount} 条未读</Tag>}
        </Title>
        {unreadCount > 0 && (
          <Button icon={<ReadOutlined />} onClick={handleMarkAllRead}>
            全部已读
          </Button>
        )}
      </div>

      <Card>
        {notifications.length > 0 ? (
          <>
            <List
              dataSource={notifications}
              renderItem={(item) => (
                <List.Item
                  key={item.id}
                  onClick={() => !item.isRead && handleMarkRead(item.id)}
                  style={{
                    background: item.isRead ? '#fff' : '#f0f8ff',
                    padding: '16px',
                    marginBottom: '8px',
                    borderRadius: '4px',
                    cursor: 'pointer',
                  }}
                >
                  <List.Item.Meta
                    avatar={<BellOutlined style={{ fontSize: 24, color: item.isRead ? '#ccc' : '#1890ff' }} />}
                    title={
                      <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                        {getTypeTag(item.type)}
                        <Text strong>{item.title}</Text>
                        {!item.isRead && <Tag color="red">未读</Tag>}
                      </div>
                    }
                    description={
                      <div>
                        <p style={{ margin: '8px 0' }}>{item.content}</p>
                        <Text type="secondary" style={{ fontSize: 12 }}>{item.createdAt}</Text>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
            <div style={{ textAlign: 'center', marginTop: 24 }}>
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={setPage}
              />
            </div>
          </>
        ) : (
          <Empty description="暂无消息" />
        )}
      </Card>
    </div>
  )
}

export default Notifications
