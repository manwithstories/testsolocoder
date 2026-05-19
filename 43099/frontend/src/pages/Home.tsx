import { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Table, Button } from 'antd'
import {
  CalendarOutlined,
  DollarOutlined,
  UserOutlined,
  ShopOutlined,
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { statsApi, orderApi } from '@/api'
import type { StatsOverview, Order } from '@/types'
import dayjs from 'dayjs'

const Home = () => {
  const navigate = useNavigate()
  const [stats, setStats] = useState<StatsOverview | null>(null)
  const [recentOrders, setRecentOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const [overview, ordersRes] = await Promise.all([
        statsApi.getOverview(),
        orderApi.list({ page: 1, page_size: 5 }),
      ])
      setStats(overview)
      setRecentOrders((ordersRes as any).list || [])
    } catch (error) {
      console.error('Load data error:', error)
    } finally {
      setLoading(false)
    }
  }

  const orderColumns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
    },
    {
      title: '项目',
      dataIndex: 'item_name',
      key: 'item_name',
    },
    {
      title: '金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (val: number) => `¥${val.toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusMap: Record<string, { text: string; color: string }> = {
          pending: { text: '待确认', color: 'orange' },
          confirmed: { text: '已确认', color: 'blue' },
          paid: { text: '已支付', color: 'green' },
          completed: { text: '已完成', color: 'gray' },
          cancelled: { text: '已取消', color: 'red' },
        }
        const info = statusMap[status] || { text: status, color: 'default' }
        return <span style={{ color: info.color }}>{info.text}</span>
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => dayjs(val).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        <Button type="link" onClick={() => navigate(`/orders`)}>
          查看
        </Button>
      ),
    },
  ]

  return (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="总预约量"
              value={stats?.total_bookings || 0}
              prefix={<CalendarOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="总收入"
              value={stats?.total_revenue || 0}
              precision={2}
              prefix={<DollarOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="用户总数"
              value={stats?.total_users || 0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={6}>
          <Card>
            <Statistic
              title="待处理订单"
              value={stats?.pending_orders || 0}
              prefix={<ShopOutlined />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      <Card title="最近订单" extra={<Button type="link" onClick={() => navigate('/orders')}>查看全部</Button>}>
        <Table
          columns={orderColumns}
          dataSource={recentOrders}
          rowKey="id"
          loading={loading}
          pagination={false}
        />
      </Card>
    </div>
  )
}

export default Home
