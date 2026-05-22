import React, { useState, useEffect } from 'react'
import { Row, Col, Card, Button, Input, Select, Rate, Tag } from 'antd'
import { useNavigate } from 'react-router-dom'
import { serviceApi } from '@/services/service'
import { ServiceItem, ServiceCategory } from '@/types'
import { formatPrice } from '@/utils'

const { Search } = Input

const Home: React.FC = () => {
  const navigate = useNavigate()
  const [categories, setCategories] = useState<ServiceCategory[]>([])
  const [services, setServices] = useState<ServiceItem[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadCategories()
    loadServices()
  }, [])

  const loadCategories = async () => {
    try {
      const res = await serviceApi.getCategories()
      setCategories(res)
    } catch (error) {
      console.error(error)
    }
  }

  const loadServices = async () => {
    setLoading(true)
    try {
      const res = await serviceApi.getList({ page_size: 8 })
      setServices(res.list)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <div
        style={{
          background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
          padding: '60px 24px',
          borderRadius: 8,
          marginBottom: 32,
          textAlign: 'center',
          color: '#fff',
        }}
      >
        <h1 style={{ color: '#fff', fontSize: 36, marginBottom: 16 }}>
          专业家政服务平台
        </h1>
        <p style={{ fontSize: 18, marginBottom: 24, opacity: 0.9 }}>
          保洁、月嫂、护工、育婴等专业服务，一键预约，放心到家
        </p>
        <Search
          placeholder="搜索服务..."
          size="large"
          style={{ width: 400 }}
          onSearch={(value) => navigate(`/services?keyword=${value}`)}
        />
      </div>

      <div className="page-header">
        <h2 className="page-title">服务分类</h2>
      </div>
      <Row gutter={[16, 16]} style={{ marginBottom: 32 }}>
        {categories.map((cat) => (
          <Col key={cat.id} xs={12} sm={8} md={6} lg={4}>
            <Card
              hoverable
              style={{ textAlign: 'center', cursor: 'pointer' }}
              onClick={() => navigate(`/services?category_id=${cat.id}`)}
            >
              <div style={{ fontSize: 32, marginBottom: 8 }}>{cat.icon || '🏠'}</div>
              <div>{cat.name}</div>
            </Card>
          </Col>
        ))}
      </Row>

      <div className="page-header">
        <h2 className="page-title">热门服务</h2>
        <Button type="link" onClick={() => navigate('/services')}>
          查看全部
        </Button>
      </div>
      <Row gutter={[16, 16]}>
        {services.map((service) => (
          <Col key={service.id} xs={24} sm={12} md={8} lg={6}>
            <Card
              hoverable
              className="service-card"
              onClick={() => navigate(`/services/${service.id}`)}
            >
              <Card.Meta
                title={service.name}
                description={
                  <div>
                    <div style={{ marginBottom: 8 }}>
                      <Rate disabled value={service.rating} allowHalf />
                      <span style={{ marginLeft: 8, color: '#999' }}>
                        {service.rating} ({service.review_count}条评价)
                      </span>
                    </div>
                    <div style={{ color: '#ff4d4f', fontSize: 18, fontWeight: 600 }}>
                      {formatPrice(service.base_price)}
                      <span style={{ fontSize: 14, color: '#999' }}>/{service.price_unit}</span>
                    </div>
                    {service.category && (
                      <Tag color="blue" style={{ marginTop: 8 }}>
                        {service.category.name}
                      </Tag>
                    )}
                  </div>
                }
              />
            </Card>
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default Home
