import { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Spin, message } from 'antd'
import { ApiOutlined, AppleOutlined, ShoppingOutlined, DollarOutlined } from '@ant-design/icons'
import api from '../../api'
import { useAuthStore } from '../../store/authStore'

function DashboardPage() {
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<any>({})

  useEffect(() => {
    if (user?.role === 'beekeeper') {
      fetchOverview()
    }
  }, [user])

  const fetchOverview = async () => {
    setLoading(true)
    try {
      const response = await api.get('/analytics/overview')
      setData(response.data)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  if (user?.role === 'beekeeper') {
    return (
      <Spin spinning={loading}>
        <Row gutter={16}>
          <Col span={6}>
            <Card>
              <Statistic
                title="蜂箱总数"
                value={data.total_beehives || 0}
                prefix={<ApiOutlined />}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="总采收量"
                value={data.total_harvest || 0}
                suffix="kg"
                prefix={<AppleOutlined />}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="商品数量"
                value={data.total_products || 0}
                prefix={<ShoppingOutlined />}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="总收入"
                value={data.total_revenue || 0}
                prefix={<DollarOutlined />}
                precision={2}
              />
            </Card>
          </Col>
        </Row>
      </Spin>
    )
  }

  return (
    <Card>
      <h2>欢迎使用养蜂管理与蜂蜜交易平台</h2>
      <p>您已成功登录，可以使用平台的各项功能。</p>
    </Card>
  )
}

export default DashboardPage
