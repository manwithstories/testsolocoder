import { Card, Row, Col, Statistic, Table, Tag, Button, Empty } from 'antd'
import { PawOutlined, CalendarOutlined, BellOutlined, ShoppingCartOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { petApi, reservationApi, alertApi, orderApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import dayjs from 'dayjs'

export default function Dashboard() {
  const navigate = useNavigate()
  const { user } = useAuthStore()

  const { data: petsData } = useQuery({
    queryKey: ['pets', 'list'],
    queryFn: () => petApi.list({ page_size: 5 }),
  })

  const { data: reservationsData } = useQuery({
    queryKey: ['reservations', 'list'],
    queryFn: () => reservationApi.list({ page_size: 5 }),
  })

  const { data: alertsData } = useQuery({
    queryKey: ['alerts', 'list'],
    queryFn: () => alertApi.list({ page_size: 5 }),
  })

  const { data: ordersData } = useQuery({
    queryKey: ['orders', 'list'],
    queryFn: () => orderApi.list({ page_size: 5 }),
  })

  const pets = petsData?.data?.items || []
  const reservations = reservationsData?.data?.items || []
  const alerts = alertsData?.data?.items || []
  const orders = ordersData?.data?.items || []

  const petColumns = [
    { title: '名字', dataIndex: 'name', key: 'name' },
    { title: '品种', dataIndex: 'breed', key: 'breed' },
    { title: '物种', dataIndex: 'species', key: 'species' },
  ]

  const reservationColumns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    {
      title: '日期',
      key: 'dates',
      render: (_: any, r: any) =>
        `${dayjs(r.check_in_date).format('MM-DD')} ~ ${dayjs(r.check_out_date).format('MM-DD')}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <StatusTag status={status} />,
    },
  ]

  const alertColumns = [
    { title: '标题', dataIndex: 'title', key: 'title' },
    { title: '类型', dataIndex: 'alert_type', key: 'alert_type' },
    {
      title: '到期时间',
      dataIndex: 'expire_at',
      key: 'expire_at',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
  ]

  return (
    <div className="space-y-6">
      <Row gutter={16}>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="宠物数量"
              value={petsData?.data?.total || 0}
              prefix={<PawOutlined className="text-sky-500" />}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="预约数量"
              value={reservationsData?.data?.total || 0}
              prefix={<CalendarOutlined className="text-sky-500" />}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="未读提醒"
              value={alertsData?.data?.total || 0}
              prefix={<BellOutlined className="text-orange-500" />}
              valueStyle={{ color: alertsData?.data?.total > 0 ? '#f97316' : undefined }}
            />
          </Card>
        </Col>
        <Col xs={12} md={6}>
          <Card>
            <Statistic
              title="订单数量"
              value={ordersData?.data?.total || 0}
              prefix={<ShoppingCartOutlined className="text-sky-500" />}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col xs={24} lg={12}>
          <Card
            title="我的宠物"
            extra={
              <Button type="link" onClick={() => navigate('/pets')}>
                查看全部
              </Button>
            }
          >
            {pets.length > 0 ? (
              <Table
                columns={petColumns}
                dataSource={pets}
                rowKey="id"
                pagination={false}
                size="small"
              />
            ) : (
              <Empty description="暂无宠物" />
            )}
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card
            title="最近预约"
            extra={
              <Button type="link" onClick={() => navigate('/reservations')}>
                查看全部
              </Button>
            }
          >
            {reservations.length > 0 ? (
              <Table
                columns={reservationColumns}
                dataSource={reservations}
                rowKey="id"
                pagination={false}
                size="small"
              />
            ) : (
              <Empty description="暂无预约" />
            )}
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col xs={24} lg={12}>
          <Card
            title="健康提醒"
            extra={
              <Button type="link" onClick={() => navigate('/alerts')}>
                查看全部
              </Button>
            }
          >
            {alerts.length > 0 ? (
              <Table
                columns={alertColumns}
                dataSource={alerts}
                rowKey="id"
                pagination={false}
                size="small"
              />
            ) : (
              <Empty description="暂无提醒" />
            )}
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card
            title="快捷操作"
          >
            <div className="flex flex-wrap gap-3">
              <Button type="primary" onClick={() => navigate('/pets')}>
                管理宠物
              </Button>
              <Button type="primary" onClick={() => navigate('/reservations/new')}>
                新建预约
              </Button>
              <Button onClick={() => navigate('/daily-records')}>
                查看动态
              </Button>
              <Button onClick={() => navigate('/orders')}>
                订单支付
              </Button>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

function StatusTag({ status }: { status: string }) {
  const statusMap: Record<string, { color: string; text: string }> = {
    pending: { color: 'orange', text: '待确认' },
    confirmed: { color: 'blue', text: '已确认' },
    checked_in: { color: 'green', text: '已入住' },
    completed: { color: 'purple', text: '已完成' },
    cancelled: { color: 'red', text: '已取消' },
  }
  const s = statusMap[status] || { color: 'default', text: status }
  return <Tag color={s.color}>{s.text}</Tag>
}
