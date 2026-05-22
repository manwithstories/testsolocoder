import { useState, useEffect } from 'react'
import { Row, Col, Card, Input, Select, Slider, Button, Pagination, Empty, Typography, Tag } from 'antd'
import { SearchOutlined, FilterOutlined } from '@ant-design/icons'
import { Link, useSearchParams } from 'react-router-dom'
import { productAPI, categoryAPI } from '@/api'
import { Product, Category } from '@/types'

const { Title } = Typography
const { Option } = Select

const ProductList = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(16)
  const [keyword, setKeyword] = useState(searchParams.get('keyword') || '')
  const [categoryId, setCategoryId] = useState<number | null>(searchParams.get('categoryId') ? Number(searchParams.get('categoryId')) : null)
  const [priceRange, setPriceRange] = useState<[number, number]>([0, 10000])
  const [sortBy, setSortBy] = useState('created_at')
  const [sortOrder, setSortOrder] = useState('desc')

  useEffect(() => {
    loadCategories()
  }, [])

  useEffect(() => {
    loadProducts()
  }, [page, categoryId, sortBy, sortOrder])

  const loadCategories = async () => {
    try {
      const res = await categoryAPI.list() as any
      setCategories(res.data)
    } catch (err) {
      console.error('加载分类失败', err)
    }
  }

  const loadProducts = async () => {
    setLoading(true)
    try {
      const params: any = {
        page,
        pageSize,
        sortBy,
        sortOrder,
        minPrice: priceRange[0],
        maxPrice: priceRange[1],
      }
      if (keyword) params.keyword = keyword
      if (categoryId) params.categoryId = categoryId

      const res = await productAPI.list(params)
      setProducts(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载商品失败', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = () => {
    setPage(1)
    loadProducts()
  }

  const handleCategoryChange = (val: number | null) => {
    setCategoryId(val)
    setPage(1)
    searchParams.set('categoryId', val?.toString() || '')
    setSearchParams(searchParams)
  }

  return (
    <div>
      <Title level={3}>商品列表</Title>

      <Card style={{ marginBottom: 24 }}>
        <Row gutter={[16, 16]} align="middle">
          <Col xs={24} sm={12} md={8}>
            <Input.Search
              placeholder="搜索商品名称..."
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              onSearch={handleSearch}
              enterButton={<SearchOutlined />}
              allowClear
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Select
              placeholder="选择分类"
              allowClear
              style={{ width: '100%' }}
              value={categoryId}
              onChange={handleCategoryChange}
            >
              {categories.map((cat) => (
                <Option key={cat.id} value={cat.id}>{cat.name}</Option>
              ))}
            </Select>
          </Col>
          <Col xs={24} md={6}>
            <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
              <span style={{ whiteSpace: 'nowrap' }}>价格:</span>
              <Slider
                range
                min={0}
                max={10000}
                step={10}
                value={priceRange}
                onChange={(val) => setPriceRange(val as [number, number])}
                style={{ flex: 1 }}
              />
              <span style={{ whiteSpace: 'nowrap' }}>¥{priceRange[0]}-¥{priceRange[1]}</span>
            </div>
          </Col>
          <Col xs={24} md={4}>
            <Button type="primary" icon={<FilterOutlined />} onClick={handleSearch}>
              筛选
            </Button>
          </Col>
        </Row>

        <Row style={{ marginTop: 16 }} gutter={[8, 8]} align="middle">
          <Col>排序:</Col>
          <Col>
            <Button.Group>
              <Button
                type={sortBy === 'created_at' ? 'primary' : 'default'}
                onClick={() => { setSortBy('created_at'); setSortOrder('desc') }}
              >
                最新
              </Button>
              <Button
                type={sortBy === 'sales' ? 'primary' : 'default'}
                onClick={() => { setSortBy('sales'); setSortOrder('desc') }}
              >
                销量
              </Button>
              <Button
                type={sortBy === 'price' && sortOrder === 'asc' ? 'primary' : 'default'}
                onClick={() => { setSortBy('price'); setSortOrder('asc') }}
              >
                价格↑
              </Button>
              <Button
                type={sortBy === 'price' && sortOrder === 'desc' ? 'primary' : 'default'}
                onClick={() => { setSortBy('price'); setSortOrder('desc') }}
              >
                价格↓
              </Button>
            </Button.Group>
          </Col>
        </Row>
      </Card>

      {products.length > 0 ? (
        <>
          <Row gutter={[16, 16]}>
            {products.map((product) => (
              <Col xs={12} sm={8} md={6} lg={6} key={product.id}>
                <Link to={`/products/${product.id}`}>
                  <Card hoverable className="card-hover" cover={
                    <img
                      alt={product.name}
                      src={product.images?.[0] || 'https://via.placeholder.com/400x300'}
                      className="product-image"
                    />
                  }>
                    <Card.Meta
                      title={<div className="product-name">{product.name}</div>}
                      description={
                        <div>
                          <span className="price">¥{product.price}</span>
                          {product.originalPrice && product.originalPrice > product.price && (
                            <span className="original-price">¥{product.originalPrice}</span>
                          )}
                          <div style={{ marginTop: 8, display: 'flex', justifyContent: 'space-between', fontSize: 12, color: '#999' }}>
                            <span>已售 {product.sales || 0}</span>
                            {product.status === 'on_sale' ? (
                              <Tag color="green">在售</Tag>
                            ) : (
                              <Tag color="red">已下架</Tag>
                            )}
                          </div>
                        </div>
                      }
                    />
                  </Card>
                </Link>
              </Col>
            ))}
          </Row>

          <div style={{ textAlign: 'center', marginTop: 32 }}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              onChange={setPage}
              showSizeChanger={false}
            />
          </div>
        </>
      ) : (
        <Empty description="暂无商品" />
      )}
    </div>
  )
}

export default ProductList
