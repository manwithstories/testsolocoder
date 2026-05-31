import { useEffect, useState } from 'react'
import { Card, Row, Col, Statistic, Button, Select, Space, DatePicker, message } from 'antd'
import { DownloadOutlined } from '@ant-design/icons'
import ReactECharts from 'echarts-for-react'
import api from '../../api'
import { useAuthStore } from '../../store/authStore'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

function AnalyticsPage() {
  const { user } = useAuthStore()
  const [overview, setOverview] = useState<any>({})
  const [productionData, setProductionData] = useState<any>({})
  const [diseaseData, setDiseaseData] = useState<any>({})
  const [salesData, setSalesData] = useState<any>({})
  const [dateRange, setDateRange] = useState<any>(null)

  const fetchOverview = async () => {
    if (user?.role !== 'beekeeper') return
    try {
      const response = await api.get('/analytics/overview')
      setOverview(response.data)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    }
  }

  const fetchProductionStats = async () => {
    if (user?.role !== 'beekeeper') return
    try {
      const params: any = {}
      if (dateRange && dateRange.length === 2) {
        params.start_date = dateRange[0].format('YYYY-MM-DD')
        params.end_date = dateRange[1].format('YYYY-MM-DD')
      }
      const response = await api.get('/analytics/production', { params })
      setProductionData(response.data)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    }
  }

  const fetchDiseaseStats = async () => {
    if (user?.role !== 'beekeeper') return
    try {
      const response = await api.get('/analytics/disease')
      setDiseaseData(response.data)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    }
  }

  const fetchSalesStats = async () => {
    if (user?.role !== 'beekeeper') return
    try {
      const response = await api.get('/analytics/sales')
      setSalesData(response.data)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    }
  }

  useEffect(() => {
    fetchOverview()
    fetchProductionStats()
    fetchDiseaseStats()
    fetchSalesStats()
  }, [user?.role])

  const handleExport = (type: string) => {
    window.open(`/api/v1/analytics/export?type=${type}`, '_blank')
  }

  const productionChartOption = {
    title: { text: '月度采收趋势' },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: productionData.monthly_stats?.map((item: any) => item.month) || [],
    },
    yAxis: { type: 'value', name: '产量(kg)' },
    series: [
      {
        name: '产量',
        type: 'line',
        smooth: true,
        data: productionData.monthly_stats?.map((item: any) => item.total_qty) || [],
        areaStyle: {},
      },
    ],
  }

  const honeyTypeChartOption = {
    title: { text: '蜂蜜类型产量分布' },
    tooltip: { trigger: 'item' },
    series: [
      {
        name: '产量',
        type: 'pie',
        radius: '60%',
        data: productionData.honey_type_stats?.map((item: any) => ({
          name: item.honey_type,
          value: item.total_qty,
        })) || [],
      },
    ],
  }

  const diseaseTypeChartOption = {
    title: { text: '病害类型分布' },
    tooltip: { trigger: 'item' },
    series: [
      {
        name: '次数',
        type: 'pie',
        radius: '60%',
        data: diseaseData.disease_type_stats?.map((item: any) => ({
          name: item.disease_type || '未知',
          value: item.count,
        })) || [],
      },
    ],
  }

  const salesChartOption = {
    title: { text: '月度销售趋势' },
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: salesData.monthly_sales_stats?.map((item: any) => item.month) || [],
    },
    yAxis: [
      { type: 'value', name: '订单数' },
      { type: 'value', name: '销售额(¥)' },
    ],
    series: [
      {
        name: '订单数',
        type: 'bar',
        data: salesData.monthly_sales_stats?.map((item: any) => item.orders_count) || [],
      },
      {
        name: '销售额',
        type: 'line',
        yAxisIndex: 1,
        data: salesData.monthly_sales_stats?.map((item: any) => item.total_revenue) || [],
      },
    ],
  }

  if (user?.role !== 'beekeeper') {
    return (
      <Card>
        <h2>数据分析</h2>
        <p>该功能仅对蜂农角色开放</p>
      </Card>
    )
  }

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Space>
          <span>日期范围：</span>
          <RangePicker
            value={dateRange}
            onChange={(dates) => setDateRange(dates)}
            onCalendarChange={() => {}}
          />
          <Button type="primary" onClick={fetchProductionStats}>查询</Button>
        </Space>
        <Space style={{ marginLeft: 24 }}>
          <Button icon={<DownloadOutlined />} onClick={() => handleExport('harvest')}>
            导出采收报表
          </Button>
          <Button icon={<DownloadOutlined />} onClick={() => handleExport('inventory')}>
            导出现存报表
          </Button>
          <Button icon={<DownloadOutlined />} onClick={() => handleExport('orders')}>
            导出订单报表
          </Button>
        </Space>
      </Card>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic title="蜂箱总数" value={overview.total_beehives || 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="总采收量" value={overview.total_harvest || 0} suffix="kg" />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="商品数量" value={overview.total_products || 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="总收入" value={overview.total_revenue || 0} prefix="¥" precision={2} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={12}>
          <Card>
            <ReactECharts option={productionChartOption} style={{ height: 350 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card>
            <ReactECharts option={honeyTypeChartOption} style={{ height: 350 }} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={12}>
          <Card title="病害统计">
            <Statistic
              title="病害发生率"
              value={diseaseData.disease_rate || '0%'}
              valueStyle={{ color: '#f5222d' }}
            />
            <ReactECharts option={diseaseTypeChartOption} style={{ height: 250 }} />
          </Card>
        </Col>
        <Col span={12}>
          <Card>
            <ReactECharts option={salesChartOption} style={{ height: 350 }} />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default AnalyticsPage
