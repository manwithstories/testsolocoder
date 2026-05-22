import React, { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Typography, List } from 'antd'
import {
  UserOutlined, BookOutlined, ShoppingOutlined,
  DollarOutlined, FormOutlined,
} from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { analyticsApi } from '@/services'
import dayjs from 'dayjs'

const { Title } = Typography

const AdminDashboard: React.FC = () => {
  const [stats, setStats] = useState<any>({})
  const [loading, setLoading] = useState(false)

  const loadData = async () => {
    setLoading(true)
    try {
      const res = await analyticsApi.adminDashboard()
      if (res.code === 0 && res.data) {
        setStats(res.data)
      }
    } catch (error) {
      console.error('Failed to load dashboard:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  return (
    <div>
      <Title level={3}>管理后台</Title>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <Link to="/admin/users">
            <Card hoverable>
              <Statistic title="学员数" value={stats.user_count || 0} prefix={<UserOutlined style={{ color: '#1890ff' }} />} />
            </Card>
          </Link>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card hoverable>
            <Statistic title="讲师数" value={stats.instructor_count || 0} prefix={<UserOutlined style={{ color: '#52c41a' }} />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card hoverable>
            <Statistic title="课程数" value={stats.course_count || 0} prefix={<BookOutlined style={{ color: '#faad14' }} />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Link to="/admin/orders">
            <Card hoverable>
              <Statistic title="订单数" value={stats.order_count || 0} prefix={<ShoppingOutlined style={{ color: '#722ed1' }} />} />
            </Card>
          </Link>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card hoverable>
            <Statistic
              title="总收入"
              value={stats.total_revenue || 0}
              precision={2}
              prefix={<DollarOutlined style={{ color: '#f5222d' }} />}
              suffix="元"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Link to="/admin/applications">
            <Card hoverable>
              <Statistic
                title="待审核讲师"
                value={stats.pending_applications || 0}
                prefix={<FormOutlined style={{ color: '#fa8c16' }} />}
                valueStyle={{ color: stats.pending_applications > 0 ? '#fa8c16' : undefined }}
              />
            </Card>
          </Link>
        </Col>
      </Row>
      <Row gutter={[16, 16]}>
        <Col xs={24} md={12}>
          <Card title="最近订单">
            <List
              dataSource={stats.recent_orders || []}
              renderItem={(item: any) => (
                <List.Item key={item.id}>
                  <List.Item.Meta
                    avatar={<ShoppingOutlined />}
                    title={item.course?.title || '-'}
                    description={`${item.user?.username || ''} · ${dayjs(item.created_at).format('MM-DD HH:mm')}`}
                  />
                  <div style={{ color: '#f5222d' }}>¥{item.amount.toFixed(2)}</div>
                </List.Item>
              )}
            />
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="最近注册">
            <List
              dataSource={stats.recent_users || []}
              renderItem={(item: any) => (
                <List.Item key={item.id}>
                  <List.Item.Meta
                    avatar={<UserOutlined />}
                    title={item.username}
                    description={dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default AdminDashboard
