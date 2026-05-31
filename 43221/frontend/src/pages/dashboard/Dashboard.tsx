import { Row, Col, Card, Statistic, List, Tag } from 'antd'
import { AppointmentOutlined, DollarOutlined, StarOutlined, ClockCircleOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { useAuthContext } from '@/contexts/AuthContext'
import { serviceApi } from '@/services/service'
import { appointmentApi } from '@/services/appointment'
import { statisticsApi } from '@/services/statistics'

export function Dashboard() {
  const { user } = useAuthContext()

  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ['dashboard-stats', user?.role],
    queryFn: async () => {
      if (user?.role === 'professional') {
        return statisticsApi.getProfessionalStats()
      }
      if (user?.role === 'admin') {
        return statisticsApi.getAdminStats()
      }
      return null
    },
    enabled: !!user && (user.role === 'professional' || user.role === 'admin'),
  })

  const { data: servicesData } = useQuery({
    queryKey: ['dashboard-services'],
    queryFn: () => serviceApi.getAll({ page: 1, page_size: 5 }),
  })

  const { data: appointmentsData } = useQuery({
    queryKey: ['dashboard-appointments'],
    queryFn: () => {
      if (user?.role === 'professional') {
        return appointmentApi.getProfessionalAppointments({ page: 1, page_size: 5 })
      }
      if (user?.role === 'client') {
        return appointmentApi.getClientAppointments({ page: 1, page_size: 5 })
      }
      return null
    },
    enabled: !!user && (user.role === 'professional' || user.role === 'client'),
  })

  return (
    <div className="page-container">
      <h2 style={{ marginBottom: 24 }}>欢迎回来，{user?.full_name}</h2>

      {(user?.role === 'professional' || user?.role === 'admin') && (
        <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
          <Col xs={24} sm={12} md={6}>
            <Card className="stat-card">
              <Statistic
                title="总预约数"
                value={stats?.appointments.total || 0}
                prefix={<AppointmentOutlined />}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Card className="stat-card" style={{ background: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)' }}>
              <Statistic
                title={user?.role === 'professional' ? '总收入' : '平台总收入'}
                value={stats?.revenue.total_revenue || 0}
                prefix={<DollarOutlined />}
                precision={2}
              />
            </Card>
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Card className="stat-card" style={{ background: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)' }}>
              <Statistic
                title="待处理"
                value={stats?.appointments.pending || 0}
                prefix={<ClockCircleOutlined />}
              />
            </Card>
          </Col>
          {user?.role === 'professional' && (
            <Col xs={24} sm={12} md={6}>
              <Card className="stat-card" style={{ background: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)' }}>
                <Statistic
                  title="平均评分"
                  value={stats?.reviews.average_rating || 0}
                  prefix={<StarOutlined />}
                  precision={1}
                />
              </Card>
            </Col>
          )}
        </Row>
      )}

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={14}>
          <Card title="热门服务">
            <List
              dataSource={servicesData?.items || []}
              renderItem={(service) => (
                <List.Item
                  actions={[
                    <Tag color="blue" key="type">{service.service_type}</Tag>,
                    <span key="price">¥{service.price}</span>,
                  ]}
                >
                  <List.Item.Meta
                    title={service.title}
                    description={service.description}
                  />
                </List.Item>
              )}
              loading={!servicesData}
            />
          </Card>
        </Col>
        <Col xs={24} lg={10}>
          <Card title="最近预约">
            <List
              dataSource={appointmentsData?.items || []}
              renderItem={(appointment) => (
                <List.Item>
                  <List.Item.Meta
                    title={`${appointment.service?.title || ''}`}
                    description={
                      <div>
                        <div>{appointment.schedule?.date} {appointment.schedule?.start_time}-{appointment.schedule?.end_time}</div>
                        <Tag color={
                          appointment.status === 'pending' ? 'orange' :
                          appointment.status === 'confirmed' ? 'blue' :
                          appointment.status === 'completed' ? 'green' :
                          appointment.status === 'cancelled' ? 'red' : 'purple'
                        }>
                          {appointment.status}
                        </Tag>
                      </div>
                    }
                  />
                </List.Item>
              )}
              loading={!appointmentsData}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}
