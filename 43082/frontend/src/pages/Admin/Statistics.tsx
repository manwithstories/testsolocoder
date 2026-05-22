import { useState, useEffect } from 'react'
import { Row, Col, Card, Typography, DatePicker, Button, Table, Tag, message } from 'antd'
import {
  DollarOutlined, ShoppingOutlined, ShopOutlined, UserOutlined,
  DownloadOutlined
} from '@ant-design/icons'
import { adminAPI } from '@/api'
import type { Dayjs } from 'dayjs'

const { Title, Text } = Typography
const { RangePicker } = DatePicker

const AdminStatistics = () => {
  const [statistics, setStatistics] = useState<any>(null)
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs] | null>(null)
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
      const res = await adminAPI.getStatistics(params) as any
      setStatistics(res.data)
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
      const res = await adminAPI.exportOrders(params)
      const url = window.URL.createObjectURL(new Blob([res as unknown as BlobPart]))
      const link = document.createElement('a')
      link.href = url
      link.download = `平台报表_${new Date().toLocaleDateString()}.xlsx`
      link.click()
      message.success('导出成功')
    } catch (err) {
      console.error('导出失败', err)
      message.error('导出失败')
    }
  }

  return (
    <div>
      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>数据概览</Title>
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
                <Text type="secondary">平台总销售额</Text>
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
                <Text type="secondary">入驻商家</Text>
                <div style={{ fontSize: 28, fontWeight: 'bold', color: '#faad14', marginTop: 8 }}>
                  {statistics?.totalShops || 0}
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
                <Text type="secondary">注册用户</Text>
                <div style={{ fontSize: 28, fontWeight: 'bold', color: '#722ed1', marginTop: 8 }}>
                  {statistics?.totalUsers || 0}
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
                <UserOutlined style={{ fontSize: 24, color: '#722ed1' }} />
              </div>
            </div>
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]}>
        <Col xs={24} md={12}>
          <Card title="待处理事项">
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <div className="flex-between">
                <span>待审核商家</span>
                <Tag color="orange">{statistics?.pendingShops || 0}</Tag>
              </div>
              <div className="flex-between">
                <span>待处理纠纷</span>
                <Tag color="red">{statistics?.pendingDisputes || 0}</Tag>
              </div>
              <div className="flex-between">
                <span>待处理退款</span>
                <Tag color="gold">{statistics?.pendingRefunds || 0}</Tag>
              </div>
            </div>
          </Card>
        </Col>
        <Col xs={24} md={12}>
          <Card title="平台商品统计">
            <div style={{ display: 'flex', flexDirection: 'column', gap: 16 }}>
              <div className="flex-between">
                <span>在售商品总数</span>
                <span style={{ fontWeight: 'bold' }}>{statistics?.totalProducts || 0}</span>
              </div>
              <div className="flex-between">
                <span>今日新增商品</span>
                <span style={{ fontWeight: 'bold', color: '#52c41a' }}>{statistics?.todayProducts || 0}</span>
              </div>
              <div className="flex-between">
                <span>今日新增订单</span>
                <span style={{ fontWeight: 'bold', color: '#1890ff' }}>{statistics?.todayOrders || 0}</span>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default AdminStatistics
