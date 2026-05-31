import { useState } from 'react'
import { Row, Col, Card, Statistic, DatePicker, Button, Select, Tabs } from 'antd'
import {
  AppointmentOutlined,
  DollarOutlined,
  StarOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined,
  CloseCircleOutlined,
  DownloadOutlined,
} from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import ReactECharts from 'echarts-for-react'
import dayjs from 'dayjs'
import { useAuthContext } from '@/contexts/AuthContext'
import { statisticsApi } from '@/services/statistics'

const { RangePicker } = DatePicker

export function Statistics() {
  const { user } = useAuthContext()
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>([
    dayjs().subtract(1, 'month'),
    dayjs(),
  ])

  const { data: stats, isLoading } = useQuery({
    queryKey: ['statistics', user?.role, dateRange],
    queryFn: () => {
      const params = {
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      }
      if (user?.role === 'professional') {
        return statisticsApi.getProfessionalStats(params)
      }
      return statisticsApi.getAdminStats(params)
    },
    enabled: !!user && (user.role === 'professional' || user.role === 'admin'),
  })

  const handleExportAppointments = async () => {
    try {
      const blob = await statisticsApi.exportAppointments({
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `appointments_${dayjs().format('YYYYMMDD')}.xlsx`
      a.click()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('Export failed:', error)
    }
  }

  const handleExportRevenue = async () => {
    try {
      const blob = await statisticsApi.exportRevenue({
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `revenue_${dayjs().format('YYYYMMDD')}.xlsx`
      a.click()
      window.URL.revokeObjectURL(url)
    } catch (error) {
      console.error('Export failed:', error)
    }
  }

  const appointmentChartOption = {
    tooltip: { trigger: 'item' },
    legend: { bottom: 0 },
    series: [
      {
        name: '预约状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
        label: { show: false, position: 'center' },
        emphasis: {
          label: { show: true, fontSize: 20, fontWeight: 'bold' },
        },
        data: [
          { value: stats?.appointments.pending || 0, name: '待确认', itemStyle: { color: '#faad14' } },
          { value: stats?.appointments.confirmed || 0, name: '已确认', itemStyle: { color: '#1890ff' } },
          { value: stats?.appointments.completed || 0, name: '已完成', itemStyle: { color: '#52c41a' } },
          { value: stats?.appointments.cancelled || 0, name: '已取消', itemStyle: { color: '#ff4d4f' } },
          { value: stats?.appointments.refunded || 0, name: '已退款', itemStyle: { color: '#722ed1' } },
        ],
      },
    ],
  }

  return (
    <div className="page-container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
        <h2>数据统计</h2>
        <div style={{ display: 'flex', gap: 16 }}>
          <RangePicker
            value={dateRange}
            onChange={(dates) => {
              if (dates && dates[0] && dates[1]) {
                setDateRange([dates[0], dates[1]])
              }
            }}
          />
          <Button icon={<DownloadOutlined />} onClick={handleExportAppointments}>
            导出预约
          </Button>
          <Button icon={<DownloadOutlined />} onClick={handleExportRevenue}>
            导出收入
          </Button>
        </div>
      </div>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="总预约数"
              value={stats?.appointments.total || 0}
              prefix={<AppointmentOutlined />}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="已完成"
              value={stats?.appointments.completed || 0}
              prefix={<CheckCircleOutlined style={{ color: '#52c41a' }} />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="待处理"
              value={stats?.appointments.pending || 0}
              prefix={<ClockCircleOutlined style={{ color: '#faad14' }} />}
              valueStyle={{ color: '#faad14' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="已取消"
              value={stats?.appointments.cancelled || 0}
              prefix={<CloseCircleOutlined style={{ color: '#ff4d4f' }} />}
              valueStyle={{ color: '#ff4d4f' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} md={12}>
          <Card title="收入统计">
            <Row gutter={16}>
              <Col span={12}>
                <Statistic
                  title="总收入"
                  value={stats?.revenue.total_revenue || 0}
                  prefix={<DollarOutlined />}
                  precision={2}
                  valueStyle={{ color: '#ff4d4f' }}
                />
              </Col>
              <Col span={12}>
                <Statistic
                  title="支付笔数"
                  value={stats?.revenue.paid_count || 0}
                  prefix={<CheckCircleOutlined />}
                />
              </Col>
            </Row>
            {stats?.revenue.refunded_amount ? (
              <div style={{ marginTop: 16, color: '#ff4d4f' }}>
                已退款金额：¥{stats.revenue.refunded_amount.toFixed(2)}
              </div>
            ) : null}
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="预约状态分布">
            <ReactECharts option={appointmentChartOption} style={{ height: 280 }} />
          </Card>
        </Col>
      </Row>

      {user?.role === 'professional' && (
        <Card title="评价统计">
          <Row gutter={16}>
            <Col span={12}>
              <Statistic
                title="平均评分"
                value={stats?.reviews.average_rating || 0}
                prefix={<StarOutlined style={{ color: '#faad14' }} />}
                precision={1}
                valueStyle={{ color: '#faad14' }}
              />
            </Col>
            <Col span={12}>
              <Statistic
                title="总评价数"
                value={stats?.reviews.total_reviews || 0}
              />
            </Col>
          </Row>
        </Card>
      )}

      {user?.role === 'admin' && (
        <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
          <Col xs={12} md={6}>
            <Card>
              <Statistic title="总用户数" value={(stats as any)?.total_users || 0} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card>
              <Statistic title="客户数" value={(stats as any)?.total_clients || 0} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card>
              <Statistic title="专业人士数" value={(stats as any)?.total_professionals || 0} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card>
              <Statistic title="服务总数" value={(stats as any)?.total_services || 0} />
            </Card>
          </Col>
        </Row>
      )}
    </div>
  )
}
