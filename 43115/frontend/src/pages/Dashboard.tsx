import React, { useState, useEffect } from 'react'
import { Row, Col, Card, Statistic, DatePicker, Select, Button, Table } from 'antd'
import { statisticsApi } from '@/services/statistics'
import { formatPrice, formatDate } from '@/utils'
import dayjs from 'dayjs'
import * as echarts from 'echarts'

const { RangePicker } = DatePicker

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<any>({
    order_stats: { total: 0, today: 0, yesterday: 0 },
    revenue_stats: { total: 0, today: 0, yesterday: 0 },
    user_stats: { total: 0, today: 0, yesterday: 0 },
    provider_stats: { total: 0, pending: 0, approved: 0 },
    service_stats: { total: 0, active: 0 },
    review_stats: { total: 0, average: 0 },
  })
  const [loading, setLoading] = useState(false)
  const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>([
    dayjs().startOf('month'),
    dayjs(),
  ])
  const [chartData, setChartData] = useState<any>({
    order_trend: [],
    service_distribution: [],
    user_activity: [],
    top_providers: [],
  })

  useEffect(() => {
    loadStats()
  }, [dateRange])

  useEffect(() => {
    initCharts()
  }, [chartData])

  const loadStats = async () => {
    setLoading(true)
    try {
      const res = await statisticsApi.getDashboard({
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      })
      setStats({
        order_stats: { total: res.total_orders || 0, today: 0, yesterday: 0 },
        revenue_stats: { total: res.total_revenue || 0, today: 0, yesterday: 0 },
        user_stats: { total: res.new_customers || 0, today: 0, yesterday: 0 },
        provider_stats: { total: res.active_providers || 0, pending: 0, approved: 0 },
        service_stats: { total: 0, active: 0 },
        review_stats: { total: 0, average: 0 },
      })
      setChartData({
        order_trend: res.order_trend || [],
        service_distribution: res.service_type_distribution || [],
        user_activity: [],
        top_providers: res.top_providers || [],
      })
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const initCharts = () => {
    const orderTrendChart = echarts.init(document.getElementById('order-trend-chart'))
    orderTrendChart.setOption({
      title: { text: '订单量趋势' },
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: chartData.order_trend.map((item: any) => item.date),
      },
      yAxis: { type: 'value' },
      series: [
        {
          data: chartData.order_trend.map((item: any) => item.order_count),
          type: 'line',
          smooth: true,
        },
      ],
    })

    const serviceDistChart = echarts.init(document.getElementById('service-dist-chart'))
    serviceDistChart.setOption({
      title: { text: '服务类型分布' },
      tooltip: { trigger: 'item' },
      series: [
        {
          type: 'pie',
          data: chartData.service_distribution.map((item: any) => ({
            value: item.order_count,
            name: item.category_name,
          })),
        },
      ],
    })
  }

  const handleExport = async () => {
    try {
      await statisticsApi.export({
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      })
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">数据看板</h1>
        <div>
          <RangePicker
            value={dateRange}
            onChange={(dates) => {
              if (dates && dates[0] && dates[1]) {
                setDateRange([dates[0], dates[1]])
              }
            }}
          />
          <Button type="primary" onClick={handleExport} style={{ marginLeft: 16 }}>
            导出报表
          </Button>
        </div>
      </div>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="总订单数"
              value={stats.order_stats?.total || 0}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="总收入"
              value={stats.revenue_stats?.total || 0}
              precision={2}
              prefix="¥"
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="总用户数"
              value={stats.user_stats?.total || 0}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="服务人员"
              value={stats.provider_stats?.total || 0}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={8}>
          <Card loading={loading}>
            <Statistic title="今日订单" value={stats.order_stats?.today || 0} />
          </Card>
        </Col>
        <Col span={8}>
          <Card loading={loading}>
            <Statistic title="今日收入" value={stats.revenue_stats?.today || 0} precision={2} prefix="¥" />
          </Card>
        </Col>
        <Col span={8}>
          <Card loading={loading}>
            <Statistic title="待审核服务人员" value={stats.provider_stats?.pending || 0} valueStyle={{ color: '#ff4d4f' }} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Card>
            <div id="order-trend-chart" style={{ height: 300 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card>
            <div id="service-dist-chart" style={{ height: 300 }} />
          </Card>
        </Col>
      </Row>

      <Card title="服务人员评分排行" style={{ marginTop: 24 }}>
        <Table
          rowKey="id"
          dataSource={chartData.top_providers}
          pagination={false}
          columns={[
            { title: '排名', key: 'rank', render: (_: any, __: any, index: number) => index + 1 },
            { title: '姓名', dataIndex: 'nickname', key: 'name', render: (_: any, record: any) => record.user?.nickname || '-' },
            { title: '评分', dataIndex: 'rating', key: 'rating' },
            { title: '订单数', dataIndex: 'order_count', key: 'order_count' },
            { title: '收入', dataIndex: 'total_income', key: 'total_income', render: (text: number) => formatPrice(text) },
          ]}
        />
      </Card>
    </div>
  )
}

export default Dashboard
