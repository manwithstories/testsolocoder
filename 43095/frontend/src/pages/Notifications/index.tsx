import React, { useState, useEffect } from 'react'
import {
  Card,
  List,
  Tabs,
  Button,
  Space,
  Typography,
  Tag,
  Modal,
  message,
  Spin,
  Empty,
  Popconfirm,
  Pagination,
  Badge,
  Divider
} from 'antd'
import {
  BellOutlined,
  CheckCircleOutlined,
  CheckOutlined,
  DeleteOutlined,
  EyeOutlined,
  CloseOutlined
} from '@ant-design/icons'
import type { TabsProps } from 'antd/es/tabs'
import type { Notification } from '@/types'
import { notificationAPI } from '@/services/api'

const { Title, Text, Paragraph } = Typography

const typeMap: Record<string, { text: string; color: string; icon: React.ReactNode }> = {
  appointment_confirmation: { text: '预约确认', color: 'blue', icon: <BellOutlined /> },
  appointment_reminder: { text: '预约提醒', color: 'orange', icon: <BellOutlined /> },
  appointment_cancelled: { text: '预约取消', color: 'red', icon: <BellOutlined /> },
  consultation_completed: { text: '问诊完成', color: 'green', icon: <CheckCircleOutlined /> },
  payment_success: { text: '支付成功', color: 'green', icon: <CheckCircleOutlined /> },
  system: { text: '系统通知', color: 'default', icon: <BellOutlined /> }
}

const NotificationsPage: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [markingRead, setMarkingRead] = useState(false)
  const [deleting, setDeleting] = useState<number | null>(null)
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [unreadCount, setUnreadCount] = useState(0)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [activeTab, setActiveTab] = useState<'all' | 'unread'>('all')
  const [detailModalVisible, setDetailModalVisible] = useState(false)
  const [selectedNotification, setSelectedNotification] = useState<Notification | null>(null)

  useEffect(() => {
    fetchNotifications()
    fetchUnreadCount()
  }, [page, pageSize, activeTab])

  const fetchNotifications = async () => {
    setLoading(true)
    try {
      const params = {
        page,
        pageSize,
        ...(activeTab === 'unread' ? { is_read: false } : {})
      }
      const result = await notificationAPI.getList(params)
      setNotifications(result.list)
      setTotal(result.total)
    } catch (error) {
      console.error('获取通知列表失败:', error)
      message.error('获取通知列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchUnreadCount = async () => {
    try {
      const result = await notificationAPI.getUnreadCount()
      setUnreadCount(result.count)
    } catch (error) {
      console.error('获取未读数量失败:', error)
    }
  }

  const handleMarkAsRead = async (id: number) => {
    setMarkingRead(true)
    try {
      await notificationAPI.markAsRead(id)
      message.success('标记已读成功')
      setNotifications((prev) =>
        prev.map((item) =>
          item.id === id ? { ...item, is_read: true, read_at: new Date().toISOString() } : item
        )
      )
      setUnreadCount((prev) => Math.max(0, prev - 1))
    } catch (error) {
      console.error('标记已读失败:', error)
      message.error('标记已读失败')
    } finally {
      setMarkingRead(false)
    }
  }

  const handleMarkAllAsRead = async () => {
    setMarkingRead(true)
    try {
      await notificationAPI.markAllAsRead()
      message.success('全部标记已读成功')
      setNotifications((prev) =>
        prev.map((item) => ({ ...item, is_read: true, read_at: new Date().toISOString() }))
      )
      setUnreadCount(0)
    } catch (error) {
      console.error('全部标记已读失败:', error)
      message.error('全部标记已读失败')
    } finally {
      setMarkingRead(false)
    }
  }

  const handleDelete = async (id: number) => {
    setDeleting(id)
    try {
      await notificationAPI.delete(id)
      message.success('删除成功')
      setNotifications((prev) => prev.filter((item) => item.id !== id))
      setTotal((prev) => prev - 1)
    } catch (error) {
      console.error('删除失败:', error)
      message.error('删除失败')
    } finally {
      setDeleting(null)
    }
  }

  const handleViewDetail = (notification: Notification) => {
    setSelectedNotification(notification)
    setDetailModalVisible(true)
    if (!notification.is_read) {
      handleMarkAsRead(notification.id)
    }
  }

  const handlePageChange = (newPage: number) => {
    setPage(newPage)
  }

  const handleTabChange = (key: string) => {
    setActiveTab(key as 'all' | 'unread')
    setPage(1)
  }

  const tabItems: TabsProps['items'] = [
    {
      key: 'all',
      label: '全部',
      children: null
    },
    {
      key: 'unread',
      label: (
        <span>
          未读
          {unreadCount > 0 && (
            <Badge count={unreadCount} size="small" className="ml-2" />
          )}
        </span>
      ),
      children: null
    }
  ]

  return (
    <div className="space-y-6">
      <Card>
        <div className="flex items-center justify-between mb-4">
          <Title level={3} style={{ margin: 0 }}>
            <BellOutlined className="mr-2" />
            通知中心
          </Title>
          <Space>
            <Button
              icon={<CheckOutlined />}
              loading={markingRead}
              onClick={handleMarkAllAsRead}
              disabled={unreadCount === 0}
            >
              全部已读
            </Button>
          </Space>
        </div>

        <Tabs
          activeKey={activeTab}
          onChange={handleTabChange}
          items={tabItems}
        />

        {loading && notifications.length === 0 ? (
          <div className="flex justify-center items-center min-h-[400px]">
            <Spin size="large" />
          </div>
        ) : notifications.length === 0 ? (
          <Empty description="暂无通知" />
        ) : (
          <>
            <List
              itemLayout="vertical"
              dataSource={notifications}
              renderItem={(item) => (
                <List.Item
                  key={item.id}
                  className={`relative p-4 border-b last:border-b-0 ${!item.is_read ? 'bg-blue-50' : ''}`}
                  style={{ paddingLeft: item.is_read ? '16px' : '16px' }}
                >
                  <div className="absolute left-0 top-4 bottom-0 w-1">
                    {!item.is_read && (
                      <div className="w-2 h-2 bg-blue-500 rounded-full ml-2" />
                    )}
                  </div>
                  <List.Item.Meta
                    title={
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          {typeMap[item.type]?.icon}
                          <Text strong>{item.title}</Text>
                          {typeMap[item.type] && (
                            <Tag color={typeMap[item.type].color}>
                              {typeMap[item.type].text}
                            </Tag>
                          )}
                        </div>
                        <Text type="secondary" className="text-xs">
                          {new Date(item.created_at).toLocaleString()}
                        </Text>
                      </div>
                    }
                    description={
                      <div className="mt-2">
                        <Paragraph ellipsis={{ rows: 2 }} className="mb-0">
                          {item.content}
                        </Paragraph>
                        <div className="flex justify-end mt-2">
                          <Space>
                            <Button
                              type="link"
                              size="small"
                              icon={<EyeOutlined />}
                              onClick={() => handleViewDetail(item)}
                            >
                              查看详情
                            </Button>
                            {!item.is_read && (
                              <Button
                                type="link"
                                size="small"
                                icon={<CheckCircleOutlined />}
                                onClick={() => handleMarkAsRead(item.id)}
                              >
                                标记已读
                              </Button>
                            )}
                            <Popconfirm
                              title="确定要删除这条通知吗？"
                              onConfirm={() => handleDelete(item.id)}
                              okText="确定"
                              cancelText="取消"
                            >
                              <Button
                                type="link"
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                                loading={deleting === item.id}
                              >
                                删除
                              </Button>
                            </Popconfirm>
                          </Space>
                        </div>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
            <div className="flex justify-end mt-4">
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={handlePageChange}
                showSizeChanger={false}
                showQuickJumper
                showTotal={(total) => `共 ${total} 条通知`}
              />
            </div>
          </>
        )}
      </Card>

      <Modal
        title="通知详情"
        open={detailModalVisible}
        onCancel={() => setDetailModalVisible(false)}
        footer={[
          <Button
            key="close"
            icon={<CloseOutlined />}
            onClick={() => setDetailModalVisible(false)}
          >
            关闭
          </Button>
        ]}
      >
        {selectedNotification && (
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              {typeMap[selectedNotification.type]?.icon}
              <Text strong className="text-lg">{selectedNotification.title}</Text>
              {typeMap[selectedNotification.type] && (
                <Tag color={typeMap[selectedNotification.type].color}>
                  {typeMap[selectedNotification.type].text}
                </Tag>
              )}
            </div>
            <Text type="secondary" className="block">
              {new Date(selectedNotification.created_at).toLocaleString()}
            </Text>
            <Divider />
            <Paragraph className="text-base">{selectedNotification.content}</Paragraph>
          </div>
        )}
      </Modal>
    </div>
  )
}

export default NotificationsPage
