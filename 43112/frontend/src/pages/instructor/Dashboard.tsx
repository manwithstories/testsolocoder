import React, { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Typography, List } from 'antd'
import { BookOutlined, UserOutlined, DollarOutlined, ShoppingOutlined } from '@ant-design/icons'
import { analyticsApi, courseApi, orderApi } from '@/services'
import { useAuthStore } from '@/store/auth'
import dayjs from 'dayjs'

const { Title } = Typography

const InstructorDashboard: React.FC = () => {
  const { user } = useAuthStore()
  const [stats, setStats] = useState<any>({})
  const [recentCourses, setRecentCourses] = useState<any[]>([])
  const [recentOrders, setRecentOrders] = useState<any[]>([])

  const loadData = async () => {
    try {
      const [dashRes, courseRes] = await Promise.all([
        analyticsApi.instructorDashboard(),
        courseApi.myCourses({ page: 1, page_size: 5 }),
      ])
      if (dashRes.code === 0 && dashRes.data) {
        setStats(dashRes.data)
        setRecentOrders(dashRes.data.recent_orders || [])
      }
      if (courseRes.code === 0 && courseRes.data) {
        setRecentCourses(courseRes.data.items)
      }
    } catch (error) {
      console.error('Failed to load dashboard:', error)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  return (
    <div>
      <Title level={3}>欢迎回来，{user?.nickname || user?.username}</Title>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="课程数量"
              value={stats.course_count || 0}
              prefix={<BookOutlined style={{ color: '#1890ff' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="学员总数"
              value={stats.total_students || 0}
              prefix={<UserOutlined style={{ color: '#52c41a' }} />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="总收入"
              value={stats.total_revenue || 0}
              precision={2}
              prefix={<DollarOutlined style={{ color: '#faad14' }} />}
              prefixCls=""
              suffix="元"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="订单数量"
              value={stats.order_count || 0}
              prefix={<ShoppingOutlined style={{ color: '#722ed1' }} />}
            />
          </Card>
        </Col>
      </Row>
      <Row gutter={[16, 16]}>
        <Col xs={24} md={12}>
          <Card title="最近课程">
            <List
              dataSource={recentCourses}
              renderItem={(item: any) => (
                <List.Item key={item.id}>
                  <List.Item.Meta
                    avatar={<BookOutlined />}
                    title={item.title}
                    description={`${item.student_count} 人学习 · ¥${item.price}`}
                  />
                  <div style={{ color: item.status === 'published' ? '#52c41a' : '#faad14' }}>
                    {item.status === 'published' ? '已上架' : '草稿'}
                  </div>
                </List.Item>
              )}
            />
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="最近订单">
            <List
              dataSource={recentOrders}
              renderItem={(item: any) => (
                <List.Item key={item.id}>
                  <List.Item.Meta
                    avatar={<ShoppingOutlined />}
                    title={item.course?.title || '-'}
                    description={`${item.user?.username || ''} · ${dayjs(item.created_at).format('MM-DD HH:mm')}`}
                  />
                  <div style={{ color: '#f5222d', fontWeight: 600 }}>¥{item.amount.toFixed(2)}</div>
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default InstructorDashboard
