import React, { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Button, DatePicker, Table, Space } from 'antd'
import {
  DollarOutlined, UserOutlined, BookOutlined,
  ShoppingOutlined, DownloadOutlined,
} from '@ant-design/icons'
import { analyticsApi } from '@/services'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

const InstructorAnalytics: React.FC = () => {
  const [stats, setStats] = useState<any>({})
  const [revenueData, setRevenueData] = useState<any>([])
  const [loading, setLoading] = useState(false)

  const loadData = async (params?: any) => {
    setLoading(true)
    try {
      const [dashRes, revRes] = await Promise.all([
        analyticsApi.instructorDashboard(),
        analyticsApi.instructorRevenue(params),
      ])
      if (dashRes.code === 0 && dashRes.data) {
        setStats(dashRes.data)
      }
      if (revRes.code === 0 && revRes.data) {
        setRevenueData(revRes.data.by_course || [])
      }
    } catch (error) {
      console.error('Failed to load analytics:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleExport = () => {
    window.open('/api/v1/analytics/instructor/export', '_blank')
  }

  useEffect(() => {
    loadData()
  }, [])

  const columns = [
    { title: '课程名称', dataIndex: 'course_title', key: 'course_title' },
    { title: '订单数', dataIndex: 'order_count', key: 'order_count' },
    {
      title: '收入',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (val: number) => <span style={{ color: '#52c41a', fontWeight: 600 }}>¥{val.toFixed(2)}</span>,
    },
  ]

  return (
    <div>
      <h2>数据分析</h2>
      <div style={{ marginBottom: 16 }}>
        <Space>
          <RangePicker
            onChange={(dates) => {
              if (dates && dates[0] && dates[1]) {
                loadData({
                  start_date: dates[0].format('YYYY-MM-DD'),
                  end_date: dates[1].format('YYYY-MM-DD'),
                })
              } else {
                loadData()
              }
            }}
          />
          <Button icon={<DownloadOutlined />} onClick={handleExport}>
            导出Excel
          </Button>
        </Space>
      </div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic title="课程数" value={stats.course_count || 0} prefix={<BookOutlined style={{ color: '#1890ff' }} />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic title="学员总数" value={stats.total_students || 0} prefix={<UserOutlined style={{ color: '#52c41a' }} />} />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="总收入"
              value={stats.total_revenue || 0}
              precision={2}
              prefix={<DollarOutlined style={{ color: '#faad14' }} />}
              suffix="元"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic title="订单数" value={stats.order_count || 0} prefix={<ShoppingOutlined style={{ color: '#722ed1' }} />} />
          </Card>
        </Col>
      </Row>
      <Card title="各课程收入统计">
        <Table
          rowKey="course_id"
          columns={columns}
          dataSource={revenueData}
          loading={loading}
          pagination={false}
        />
      </Card>
    </div>
  )
}

export default InstructorAnalytics
