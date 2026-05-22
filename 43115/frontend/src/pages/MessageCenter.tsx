import React, { useState, useEffect } from 'react'
import { List, Card, Button, Tag, Badge, Empty, message, Modal, Input } from 'antd'
import { messageApi } from '@/services/message'
import { Message } from '@/types'
import { formatDate, getMessageTypeText } from '@/utils'
import { usePagination } from '@/hooks'

const MessageCenter: React.FC = () => {
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination(20)
  const [messages, setMessages] = useState<Message[]>([])
  const [loading, setLoading] = useState(false)
  const [unreadCount, setUnreadCount] = useState(0)
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [selectedMessage, setSelectedMessage] = useState<Message | null>(null)

  useEffect(() => {
    loadMessages()
    loadUnreadCount()
  }, [page, pageSize])

  const loadMessages = async () => {
    setLoading(true)
    try {
      const res = await messageApi.getList({
        page,
        page_size: pageSize,
      })
      setMessages(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const loadUnreadCount = async () => {
    try {
      const res = await messageApi.getUnreadCount()
      setUnreadCount(res.unread_count)
    } catch (error) {
      console.error(error)
    }
  }

  const handleMarkAsRead = async (id: number) => {
    try {
      await messageApi.read(id)
      loadMessages()
      loadUnreadCount()
    } catch (error) {
      console.error(error)
    }
  }

  const handleMarkAllAsRead = async () => {
    try {
      await messageApi.readAll()
      message.success('全部标记为已读')
      loadMessages()
      loadUnreadCount()
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">消息中心</h1>
        <Badge count={unreadCount} offset={[-10, 0]}>
          <Button onClick={handleMarkAllAsRead}>全部已读</Button>
        </Badge>
      </div>

      <Card>
        {messages.length === 0 ? (
          <Empty description="暂无消息" />
        ) : (
          <List
            loading={loading}
            dataSource={messages}
            renderItem={(item) => (
              <List.Item
                key={item.id}
                onClick={() => {
                  setSelectedMessage(item)
                  setDetailModalVisible(true)
                  if (!item.is_read) {
                    handleMarkAsRead(item.id)
                  }
                }}
                style={{
                  cursor: 'pointer',
                  background: item.is_read ? '#fff' : '#f0f5ff',
                }}
              >
                <List.Item.Meta
                  title={
                    <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                      {!item.is_read && <Badge status="processing" />}
                      <span style={{ fontWeight: item.is_read ? 400 : 600 }}>
                        {item.title}
                      </span>
                      <Tag color="blue">{getMessageTypeText(item.type)}</Tag>
                    </div>
                  }
                  description={
                    <div>
                      <div style={{ color: '#333' }}>{item.content.slice(0, 100)}...</div>
                      <div style={{ color: '#999', marginTop: 4 }}>{formatDate(item.created_at)}</div>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        )}
      </Card>

      <Modal
        title={selectedMessage?.title}
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={null}
        width={600}
      >
        {selectedMessage && (
          <div>
            <div style={{ marginBottom: 16, color: '#999' }}>
              <Tag color="blue">{getMessageTypeText(selectedMessage.type)}</Tag>
              <span>{formatDate(selectedMessage.created_at)}</span>
            </div>
            <div style={{ whiteSpace: 'pre-wrap', lineHeight: 1.8 }}>
              {selectedMessage.content}
            </div>
          </div>
        )}
      </Modal>
    </div>
  )
}

export default MessageCenter
