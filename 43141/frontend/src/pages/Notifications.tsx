import { useState, useEffect } from 'react'
import { List, Button, Tag, Card, Space, message, Badge } from 'antd'
import { CheckOutlined, BellOutlined } from '@ant-design/icons'
import api from '../api'
import type { Notification } from '../types'

export default function Notifications({ onRead }: { onRead: () => void }) {
  const [notifications, setNotifications] = useState<Notification[]>([])

  const fetchData = () => {
    api.get('/notifications').then((res) => setNotifications(res.data))
  }

  useEffect(() => { fetchData() }, [])

  const handleMarkRead = async (id: number) => {
    try {
      await api.put(`/notifications/${id}/read`)
      message.success('已标记为已读')
      fetchData()
      onRead()
    } catch (e: any) {
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const handleMarkAllRead = async () => {
    try {
      await api.put('/notifications/read-all')
      message.success('全部标记为已读')
      fetchData()
      onRead()
    } catch (e: any) {
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const typeColors: Record<string, string> = {
    match_result: 'blue',
    referee_assignment: 'purple',
    match_schedule: 'green',
    promotion: 'gold',
    system: 'default',
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>消息通知</h2>
        <Button icon={<CheckOutlined />} onClick={handleMarkAllRead}>全部标记已读</Button>
      </div>
      <Card>
        <List
          itemLayout="horizontal"
          dataSource={notifications}
          renderItem={(item) => (
            <List.Item
              actions={[
                !item.is_read && (
                  <Button type="link" size="small" onClick={() => handleMarkRead(item.id)}>
                    标记已读
                  </Button>
                )
              ]}
            >
              <List.Item.Meta
                avatar={<Badge dot={!item.is_read} offset={[2, 2]}><BellOutlined style={{ fontSize: 24 }} /></Badge>}
                title={
                  <Space>
                    <span>{item.title}</span>
                    <Tag color={typeColors[item.type] || 'default'}>{item.type}</Tag>
                    {!item.is_read && <Tag color="red">未读</Tag>}
                  </Space>
                }
                description={
                  <div>
                    <div>{item.content}</div>
                    <small style={{ color: '#999' }}>{new Date(item.created_at).toLocaleString()}</small>
                  </div>
                }
              />
            </List.Item>
          )}
        />
        {notifications.length === 0 && <p style={{ textAlign: 'center', color: '#999' }}>暂无通知</p>}
      </Card>
    </div>
  )
}
