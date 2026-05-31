import { useState } from 'react'
import { Row, Col, Card, Select, Input, Tag, Rate, Button, Pagination, Empty } from 'antd'
import { SearchOutlined, StarOutlined, ClockCircleOutlined, DollarOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { serviceApi } from '@/services/service'
import { Service, ServiceType } from '@/types'

export function ServiceList() {
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const [pageSize] = useState(12)
  const [searchText, setSearchText] = useState('')
  const [serviceType, setServiceType] = useState<string>('')

  const { data, isLoading } = useQuery({
    queryKey: ['services', page, pageSize, serviceType, searchText],
    queryFn: () =>
      serviceApi.getAll({
        page,
        page_size: pageSize,
        service_type: serviceType || undefined,
      }),
  })

  const filteredServices = (data?.items || []).filter(
    (service: Service) =>
      !searchText ||
      service.title.toLowerCase().includes(searchText.toLowerCase()) ||
      service.description?.toLowerCase().includes(searchText.toLowerCase())
  )

  const serviceTypeOptions = [
    { value: '', label: '全部类型' },
    { value: 'legal', label: '法律咨询' },
    { value: 'counseling', label: '心理咨询' },
    { value: 'financial', label: '财务咨询' },
    { value: 'other', label: '其他服务' },
  ]

  return (
    <div className="page-container">
      <div style={{ marginBottom: 24, display: 'flex', gap: 16, flexWrap: 'wrap' }}>
        <Input
          placeholder="搜索服务..."
          prefix={<SearchOutlined />}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          style={{ width: 240 }}
          allowClear
        />
        <Select
          value={serviceType}
          onChange={setServiceType}
          options={serviceTypeOptions}
          style={{ width: 160 }}
        />
      </div>

      {filteredServices.length === 0 && !isLoading ? (
        <Empty description="暂无服务" />
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {filteredServices.map((service: Service) => (
              <Col xs={24} sm={12} md={8} lg={6} key={service.id}>
                <Card
                  hoverable
                  className="service-card"
                  onClick={() => navigate(`/services/${service.id}`)}
                  cover={
                    <div
                      style={{
                        height: 120,
                        background: `linear-gradient(135deg, ${
                          service.service_type === 'legal'
                            ? '#667eea 0%, #764ba2 100%'
                            : service.service_type === 'counseling'
                            ? '#f093fb 0%, #f5576c 100%'
                            : service.service_type === 'financial'
                            ? '#4facfe 0%, #00f2fe 100%'
                            : '#43e97b 0%, #38f9d7 100%'
                        })`,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: '#fff',
                        fontSize: 48,
                      }}
                    >
                      {service.service_type === 'legal' ? '⚖️' :
                       service.service_type === 'counseling' ? '💬' :
                       service.service_type === 'financial' ? '💰' : '📋'}
                    </div>
                  }
                >
                  <Card.Meta
                    title={
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <span style={{ fontSize: 16, fontWeight: 500 }}>{service.title}</span>
                        <Tag color="blue">{service.service_type}</Tag>
                      </div>
                    }
                    description={
                      <div>
                        <div style={{ marginBottom: 8, color: '#666', fontSize: 13 }}>
                          {service.description?.substring(0, 50)}...
                        </div>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                          <div>
                            <Rate disabled allowHalf value={service.average_rating} style={{ fontSize: 12 }} />
                            <span style={{ marginLeft: 8, fontSize: 12, color: '#999' }}>
                              ({service.review_count})
                            </span>
                          </div>
                          <div style={{ color: '#ff4d4f', fontWeight: 600 }}>¥{service.price}</div>
                        </div>
                        <div style={{ marginTop: 8, fontSize: 12, color: '#999' }}>
                          <ClockCircleOutlined /> {service.duration_minutes}分钟
                        </div>
                      </div>
                    }
                  />
                </Card>
              </Col>
            ))}
          </Row>

          {data && data.total > pageSize && (
            <div style={{ marginTop: 24, textAlign: 'center' }}>
              <Pagination
                current={page}
                pageSize={pageSize}
                total={data.total}
                onChange={setPage}
              />
            </div>
          )}
        </>
      )}
    </div>
  )
}
