import { useState, useEffect } from 'react'
import { Row, Col, Card, Statistic, DatePicker, Table } from 'antd'
import {
  CalendarOutlined,
  DollarOutlined,
  UserOutlined,
  ShopOutlined,
  ToolOutlined,
  ClockCircleOutlined,
  RiseOutlined,
} from '@ant-design/icons'
import ReactECharts from 'echarts-for-react'
import { statsApi } from '@/api'
import type { StatsOverview, BookingStats, RevenueStats, PopularVenue } from '@/types'
import dayjs from 'dayjs'

const StatsPage = () => {
  const [overview, setOverview] = useState<StatsOverview | null>(null)
  const [bookingStats, setBookingStats] = useState<BookingStats[]>([])
  const [revenueStats, setRevenueStats] = useState<RevenueStats[]>([])
  const [popularVenues, setPopularVenues] = useState<PopularVenue[]>([])
  const [dateRange, setDateRange] = useState<[string, string]>([
    dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
    dayjs().format('YYYY-MM-DD'),
  ])

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [overviewData, bookingData, revenueData, venuesData] = await Promise.all([
        statsApi.getOverview(),
        statsApi.getBookingStats({ start_date: dateRange[0], end_date: dateRange[1] }),
        statsApi.getRevenueStats({ start_date: dateRange[0], end_date: dateRange[1] }),
        statsApi.getPopularVenues({ start_date: dateRange[0], end_date: dateRange[1] }),
      ])
      setOverview(overviewData)
      setBookingStats(bookingData)
      setRevenueStats(revenueData)
      setPopularVenues(venuesData)
    } catch (error: any) {
      console.error('Load stats error:', error)
    }
  }

  const bookingChartOption = {
    title: { text: '预约量趋势' },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: bookingStats.map((s) => s.date),
    },
    yAxis: { type: 'value' },
    series: [
      {
        data: bookingStats.map((s) => s.count),
        type: 'line',
        smooth: true,
        areaStyle: { color: 'rgba(24, 144, 255, 0.2)' },
        lineStyle: { color: '#1890ff' },
      },
    ],
  }

  const revenueChartOption = {
    title: { text: '收入趋势' },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: revenueStats.map((s) => s.date),
    },
    yAxis: { type: 'value' },
    series: [
      {
        data: revenueStats.map((s) => s.amount),
        type: 'bar',
        itemStyle: { color: '#52c41a' },
      },
    ],
  }

  const popularVenuesColumns = [
    {
      title: '排名',
      key: 'rank',
      render: (_: any, __: any, index: number) => index + 1,
      width: 60,
    },
    {
      title: '场地名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '预约次数',
      dataIndex: 'bookings',
      key: 'bookings',
    },
    {
      title: '收入',
      dataIndex: 'revenue',
      key: 'revenue',
      render: (val: number) => `¥${val.toFixed(2)}`,
    },
  ]

  return (
    <div>
      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="总预约量"
              value={overview?.total_bookings || 0}
              prefix={<CalendarOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="总收入"
              value={overview?.total_revenue || 0}
              precision={2}
              prefix={<DollarOutlined />}
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="用户总数"
              value={overview?.total_users || 0}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="待处理订单"
              value={overview?.pending_orders || 0}
              prefix={<ClockCircleOutlined />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="今日预约"
              value={overview?.today_bookings || 0}
              prefix={<RiseOutlined />}
              valueStyle={{ color: '#13c2c2' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="今日收入"
              value={overview?.today_revenue || 0}
              precision={2}
              prefix={<DollarOutlined />}
              valueStyle={{ color: '#eb2f96' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="场地总数"
              value={overview?.total_venues || 0}
              prefix={<ShopOutlined />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
        <Col xs={12} sm={8} md={6}>
          <Card>
            <Statistic
              title="设备总数"
              value={overview?.total_devices || 0}
              prefix={<ToolOutlined />}
              valueStyle={{ color: '#a0d911' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginBottom: 16 }}>
        <Col xs={24} lg={12}>
          <Card>
            <ReactECharts option={bookingChartOption} style={{ height: 350 }} />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card>
            <ReactECharts option={revenueChartOption} style={{ height: 350 }} />
          </Card>
        </Col>
      </Row>

      <Card title="热门场地排行">
        <Table
          columns={popularVenuesColumns}
          dataSource={popularVenues}
          rowKey="id"
          pagination={false}
        />
      </Card>
    </div>
  )
}

export default StatsPage
