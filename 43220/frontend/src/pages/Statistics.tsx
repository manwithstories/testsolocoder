import { useState } from 'react'
import { Card, Row, Col, DatePicker, Statistic, Button, Space, message } from 'antd'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  Legend,
} from 'recharts'
import { DownloadOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { statisticsApi, exportApi } from '@/services/api'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

export default function Statistics() {
  const [dateRange, setDateRange] = useState<any>([
    dayjs().subtract(30, 'day'),
    dayjs(),
  ])

  const params = {
    store_id: '00000000-0000-0000-0000-000000000000',
    start_date: dateRange[0]?.toISOString(),
    end_date: dateRange[1]?.toISOString(),
  }

  const { data: revenueData } = useQuery({
    queryKey: ['stats', 'revenue', params],
    queryFn: () => statisticsApi.revenueTrend(params),
  })

  const { data: occupancyData } = useQuery({
    queryKey: ['stats', 'occupancy', params],
    queryFn: () => statisticsApi.occupancyRate(params),
  })

  const { data: petTypeData } = useQuery({
    queryKey: ['stats', 'petTypes', params],
    queryFn: () => statisticsApi.petTypeDistribution(params),
  })

  const { data: orderStatsData } = useQuery({
    queryKey: ['stats', 'orders', params],
    queryFn: () => statisticsApi.orderStatistics(params),
  })

  const revenue = revenueData?.data || []
  const occupancy = occupancyData?.data?.occupancy_rate || 0
  const petTypes = petTypeData?.data || []
  const orderStats = orderStatsData?.data || {}

  const COLORS = ['#0ea5e9', '#22c55e', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4']

  const handleExport = async (type: 'excel' | 'pdf') => {
    try {
      const res = type === 'excel'
        ? await exportApi.excel(params)
        : await exportApi.pdf(params)
      const blob = new Blob([res])
      const url = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = type === 'excel' ? 'statistics.xlsx' : 'statistics.pdf'
      link.click()
      URL.revokeObjectURL(url)
      message.success('导出成功')
    } catch (err: any) {
      message.error(err.message || '导出失败')
    }
  }

  return (
    <div className="space-y-4">
      <Card
        title="数据统计"
        extra={
          <Space>
            <RangePicker value={dateRange} onChange={setDateRange} />
            <Button icon={<DownloadOutlined />} onClick={() => handleExport('excel')}>
              导出Excel
            </Button>
            <Button icon={<DownloadOutlined />} onClick={() => handleExport('pdf')}>
              导出PDF
            </Button>
          </Space>
        }
      >
        <Row gutter={16} className="mb-6">
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="总收入" value={orderStats.total_amount || 0} prefix="¥" precision={2} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="订单数" value={orderStats.order_count || 0} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="入住率" value={occupancy * 100} suffix="%" precision={1} />
            </Card>
          </Col>
          <Col xs={12} md={6}>
            <Card size="small">
              <Statistic title="退款金额" value={orderStats.refund_amount || 0} prefix="¥" precision={2} />
            </Card>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} lg={16}>
            <Card title="营收趋势" size="small">
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={revenue}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip formatter={(value: number) => `¥${value}`} />
                  <Line type="monotone" dataKey="total" stroke="#0ea5e9" strokeWidth={2} dot={{ fill: '#0ea5e9' }} />
                </LineChart>
              </ResponsiveContainer>
            </Card>
          </Col>
          <Col xs={24} lg={8}>
            <Card title="宠物品类分布" size="small">
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={petTypes}
                    dataKey="count"
                    nameKey="species"
                    cx="50%"
                    cy="50%"
                    outerRadius={100}
                    label={(entry: any) => entry.species}
                  >
                    {petTypes.map((_: any, index: number) => (
                      <Cell key={index} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>
      </Card>
    </div>
  )
}
