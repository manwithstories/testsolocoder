import { useState, useEffect } from 'react'
import { Row, Col, Card, Typography, Carousel, Tag, Empty } from 'antd'
import { Link } from 'react-router-dom'
import { productAPI, categoryAPI, shopAPI } from '@/api'
import { Product, Category, Shop } from '@/types'

const { Title } = Typography

const Home = () => {
  const [products, setProducts] = useState<Product[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [hotShops, setHotShops] = useState<Shop[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const [productRes, categoryRes, shopRes] = await Promise.all([
        productAPI.list({ page: 1, pageSize: 8, sortBy: 'sales', sortOrder: 'desc' }),
        categoryAPI.list(),
        shopAPI.list({ page: 1, pageSize: 6, status: 'approved' }),
      ])
      setProducts(productRes.data.data)
      setCategories(categoryRes.data)
      setHotShops(shopRes.data.data)
    } catch (err) {
      console.error('加载首页数据失败', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <Carousel autoplay style={{ marginBottom: 32, borderRadius: 8, overflow: 'hidden' }}>
        <div>
          <div style={{
            height: 320,
            background: 'linear-gradient(90deg, #667eea 0%, #764ba2 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: 32,
            fontWeight: 'bold',
          }}>
            欢迎来到多商家商城
          </div>
        </div>
        <div>
          <div style={{
            height: 320,
            background: 'linear-gradient(90deg, #f093fb 0%, #f5576c 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: 32,
            fontWeight: 'bold',
          }}>
            新用户专享优惠
          </div>
        </div>
        <div>
          <div style={{
            height: 320,
            background: 'linear-gradient(90deg, #4facfe 0%, #00f2fe 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: 32,
            fontWeight: 'bold',
          }}>
            品质商家 优质商品
          </div>
        </div>
      </Carousel>

      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>热门分类</Title>
      </div>
      <Row gutter={[16, 16]} style={{ marginBottom: 32 }}>
        {categories.slice(0, 8).map((cat) => (
          <Col xs={6} sm={4} md={3} key={cat.id}>
            <Link to={`/products?categoryId=${cat.id}`}>
              <Card hoverable className="card-hover" style={{ textAlign: 'center' }}>
                <div style={{ fontSize: 32, marginBottom: 8 }}>📦</div>
                <div>{cat.name}</div>
              </Card>
            </Link>
          </Col>
        ))}
      </Row>

      <div className="page-header">
        <Title level={3} style={{ margin: 0 }}>热销商品</Title>
        <Link to="/products">查看全部</Link>
      </div>
      {products.length > 0 ? (
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
                        <div style={{ marginTop: 8, color: '#999', fontSize: 12 }}>
                          已售 {product.sales || 0}
                        </div>
                      </div>
                    }
                  />
                </Card>
              </Link>
            </Col>
          ))}
        </Row>
      ) : (
        <Empty description="暂无商品" />
      )}

      <div className="page-header" style={{ marginTop: 48 }}>
        <Title level={3} style={{ margin: 0 }}>热门店铺</Title>
      </div>
      <Row gutter={[16, 16]}>
        {hotShops.map((shop) => (
          <Col xs={24} sm={12} md={8} key={shop.id}>
            <Link to={`/shops/${shop.id}`}>
              <Card hoverable className="card-hover">
                <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                  <img
                    src={shop.logo || 'https://via.placeholder.com/80'}
                    alt={shop.name}
                    className="shop-logo"
                  />
                  <div style={{ flex: 1 }}>
                    <div style={{ fontSize: 16, fontWeight: 'bold', marginBottom: 4 }}>{shop.name}</div>
                    <div style={{ color: '#999', fontSize: 12 }}>{shop.description}</div>
                    <div style={{ marginTop: 8 }}>
                      <Tag color="blue">好评率 {shop.rating || 5}%</Tag>
                      <Tag color="orange">商品 {shop.productCount || 0}</Tag>
                    </div>
                  </div>
                </div>
              </Card>
            </Link>
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default Home
