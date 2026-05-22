import { useState, useEffect } from 'react'
import { Row, Col, Card, Typography, DatePicker, Button, Table, Tag, message } from 'antd'
import {
  DollarOutlined, ShoppingOutlined, ShopOutlined,
  ArrowUpOutlined, ArrowDownOutlined, DownloadOutlined
} from '@ant-design/icons'
import { statisticsAPI } from '@/api'
import type { Dayjs } from 'dayjs'

const { Title, Text } = Typography
const { RangePicker } = DatePicker

const SellerStatistics = () => {
  const [statistics, setStatistics] = useState<any>(null)
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs] | null>(null)
  const [topProducts, setTopProducts] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadStatistics()
  }, [dateRange])

  const loadStatistics = async () => {
    setLoading(true)
    try {
      const params: any = {}
      if (dateRange) {
        params.start_date = dateRange[0].format('YYYY-MM-DD')
        params.end_date = dateRange[1].format('YYYY-MM-DD')
      }
      const res = await statisticsAPI.getShopStatistics(params) as any
      setStatistics(res.data)
      setTopProducts(res.data.topProducts || [])
    } catch (err) {
      console.error('加载统计数据失败', err)
    } finally {
      setLoading(false)
    }
  }

  const handleExport = async () => {
    try {
      const params: any = {}
      if (dateRange) {
        params.start_date = dateRange[0].format('YYYY-MM-DD')
        params.end_date = dateRange[1].format('YYYY-MM-DD')
      }
      const res: any = await statisticsAPI.exportShopOrders(params)
      const url = window.URL.createObjectURL(new Blob([res.data]))
      const link = document.createElement('a')
      link.href = url
      link.download = `销售报表_${new Date().toLocaleDateString()}.xlsx`
      link.click()
      message.success('导出成功')
    } catch (err) {
      console.error('导出失败', err)
      message.error('导出失败')
    }
  }

  const productColumns = [
    {
      title: '商品',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: any) => (
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <img src={record.image} style={{ width: 50, height: 50, objectFit: 'cover', borderRadius: 4 }} />
          <span className="truncate" style={{ maxWidth: 200 }}>{text}</span>
        </div>
      ),
    },
    {
      title: '销量',
      dataIndex: 'sales',
      key: 'sales',
      width: 100,
    },
    {
      title: '销售额',
      dataIndex: 'revenue',
      key: 'revenue',
      width: 120,
      render: (val: number) => <span className="price">¥{val.toFixed(2)}</span>,
    },
  ]

  return (
    <div>
      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>数据统计</Title>
        <div style={{ display: 'flex', gap: 12 }}>
          <RangePicker value={dateRange} onChange={(val) => setDateRange(val as [Dayjs, Dayjs])} />
          <Button icon={<DownloadOutlined />} onClick={handleExport}>
            导出报表
          </Button>
        </div>
      </div>

      <Row gutter={[16, 16]} style={{ marginBottom: 24 }}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <div>
                <Text type="secondary">总销售额</Text>
                <div style={{ fontSize: 28, fontWeight: 'bold', color: '#1890ff', marginTop: 8 }}>
                  ¥{statistics?.totalRevenue?.toFixed(2) || 0}
                </div>
              </div>
              <div style={{
                width: 48,
                height: 48,
                borderRadius: '50%',
                background: '#e6f7ff',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}>
                <DollarOutlined style={{ fontSize: 24, color: '#1890ff' }} />
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <div>
                <Text type="secondary">订单总数</Text>
                <div style={{ fontSize: 28, fontWeight: 'bold', color: '#52c41a', marginTop: 8 }}>
                  {statistics?.totalOrders || 0}
                </div>
              </div>
              <div style={{
                width: 48,
                height: 48,
                borderRadius: '50%',
                background: '#f6ffed',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}>
                <ShoppingOutlined style={{ fontSize: 24, color: '#52c41a' }} />
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <div>
                <Text type="secondary">商品数量</Text>
                <div style={{ fontSize: 28, fontWeight: 'bold', color: '#faad14', marginTop: 8 }}>
                  {statistics?.totalProducts || 0}
                </div>
              </div>
              <div style={{
                width: 48,
                height: 48,
                borderRadius: '50%',
                background: '#fffbe6',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}>
                <ShopOutlined style={{ fontSize: 24, color: '#faad14' }} />
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
              <div>
                <Text type="secondary">平均客单价</Text>
                <div style={{ fontSize: 28, fontWeight: 'bold', color: '#722ed1', marginTop: 8 }}>
                  ¥{statistics?.avgOrderAmount?.toFixed(2) || 0}
                </div>
              </div>
              <div style={{
                width: 48,
                height: 48,
                borderRadius: '50%',
                background: '#f9f0ff',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}>
                <DollarOutlined style={{ fontSize: 24, color: '#722ed1' }} />
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <Card title="商品销售排行">
        <Table
          columns={productColumns}
          dataSource={topProducts}
          rowKey="id"
          pagination={false}
        />
      </Card>
    </div>
  )
}

export default SellerStatistics
