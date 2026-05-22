import React, { useState, useEffect } from 'react'
import { Row, Col, Card, Input, Select, Slider, Empty, Rate, Tag } from 'antd'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { serviceApi } from '@/services/service'
import { ServiceItem, ServiceCategory } from '@/types'
import { formatPrice } from '@/utils'
import { usePagination } from '@/hooks'

const { Search } = Input

const ServiceList: React.FC = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination(12)
  const [categories, setCategories] = useState<ServiceCategory[]>([])
  const [services, setServices] = useState<ServiceItem[]>([])
  const [loading, setLoading] = useState(false)
  const [categoryId, setCategoryId] = useState<number | undefined>(
    searchParams.get('category_id') ? Number(searchParams.get('category_id')) : undefined
  )
  const [keyword, setKeyword] = useState(searchParams.get('keyword') || '')
  const [priceRange, setPriceRange] = useState<[number, number]>([0, 1000])
  const [minRating, setMinRating] = useState<number | undefined>()

  useEffect(() => {
    loadCategories()
  }, [])

  useEffect(() => {
    loadServices()
  }, [page, pageSize, categoryId, minRating, priceRange])

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
      const res = await serviceApi.getList({
        page,
        page_size: pageSize,
        category_id: categoryId,
        min_price: priceRange[0],
        max_price: priceRange[1],
        min_rating: minRating,
        keyword: keyword || undefined,
      })
      setServices(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">服务市场</h1>
      </div>

      <div className="filter-bar">
        <Search
          placeholder="搜索服务名称"
          size="large"
          style={{ width: 300 }}
          allowClear
          defaultValue={keyword}
          onSearch={(value) => {
            setKeyword(value)
            setPage(1)
            loadServices()
          }}
        />
        <Select
          placeholder="服务分类"
          size="large"
          style={{ width: 200 }}
          allowClear
          value={categoryId}
          onChange={(value) => {
            setCategoryId(value)
            setPage(1)
          }}
          options={categories.map((cat) => ({
            label: cat.name,
            value: cat.id,
          }))}
        />
        <Select
          placeholder="最低评分"
          size="large"
          style={{ width: 150 }}
          allowClear
          value={minRating}
          onChange={(value) => {
            setMinRating(value)
            setPage(1)
          }}
          options={[
            { label: '3星以上', value: 3 },
            { label: '4星以上', value: 4 },
            { label: '4.5星以上', value: 4.5 },
          ]}
        />
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <span>价格区间:</span>
          <Slider
            range
            style={{ width: 200 }}
            min={0}
            max={1000}
            step={10}
            value={priceRange}
            onChange={(value) => setPriceRange(value as [number, number])}
            onChangeComplete={() => setPage(1)}
          />
          <span>{priceRange[0]} - {priceRange[1]}</span>
        </div>
      </div>

      {services.length === 0 ? (
        <Empty description="暂无服务" />
      ) : (
        <>
          <Row gutter={[16, 16]}>
            {services.map((service) => (
              <Col key={service.id} xs={24} sm={12} md={8} lg={6}>
                <Card
                  hoverable
                  className="service-card"
                  loading={loading}
                  onClick={() => navigate(`/services/${service.id}`)}
                >
                  <Card.Meta
                    title={service.name}
                    description={
                      <div>
                        <div style={{ marginBottom: 8 }}>
                          <Rate disabled value={service.rating} allowHalf />
                          <span style={{ marginLeft: 8, color: '#999' }}>
                            {service.rating} ({service.review_count})
                          </span>
                        </div>
                        <div style={{ color: '#ff4d4f', fontSize: 18, fontWeight: 600 }}>
                          {formatPrice(service.base_price)}
                          <span style={{ fontSize: 14, color: '#999' }}>/{service.price_unit}</span>
                        </div>
                        <div style={{ marginTop: 8 }}>
                          {service.category && (
                            <Tag color="blue">{service.category.name}</Tag>
                          )}
                          <Tag color="green">已售{service.order_count}</Tag>
                        </div>
                      </div>
                    }
                  />
                </Card>
              </Col>
            ))}
          </Row>
          <div style={{ textAlign: 'center', marginTop: 24 }}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              showSizeChanger
              showQuickJumper
              onChange={(p, ps) => {
                setPage(p)
                setPageSize(ps)
              }}
            />
          </div>
        </>
      )}
    </div>
  )
}

export default ServiceList
