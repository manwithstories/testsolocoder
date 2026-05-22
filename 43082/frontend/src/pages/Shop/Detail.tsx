import { useState, useEffect } from 'react'
import { Row, Col, Card, Typography, Button, Tag, Avatar, Pagination, Empty } from 'antd'
import { HeartOutlined, HeartFilled, ShopOutlined } from '@ant-design/icons'
import { Link, useParams } from 'react-router-dom'
import { shopAPI, productAPI, favoriteAPI } from '@/api'
import { Shop, Product } from '@/types'
import { useAppStore } from '@/store'

const { Title, Text } = Typography

const ShopDetail = () => {
  const { id } = useParams()
  const { user } = useAppStore()
  const [shop, setShop] = useState<Shop | null>(null)
  const [products, setProducts] = useState<Product[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(12)
  const [isFavorite, setIsFavorite] = useState(false)

  useEffect(() => {
    if (id) {
      loadShop()
      loadProducts()
    }
  }, [id, page])

  const loadShop = async () => {
    try {
      const res = await shopAPI.getById(Number(id)) as any
      setShop(res.data)
    } catch (err) {
      console.error('加载店铺失败', err)
    }
  }

  const loadProducts = async () => {
    try {
      const res = await productAPI.list({ shopId: Number(id), page, pageSize, status: 'on_sale' })
      setProducts(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载商品失败', err)
    }
  }

  const toggleFavorite = async () => {
    if (!user) {
      return
    }
    try {
      await favoriteAPI.toggleShop(shop!.id)
      setIsFavorite(!isFavorite)
    } catch (err) {
      console.error('收藏操作失败', err)
    }
  }

  if (!shop) {
    return <div style={{ textAlign: 'center', padding: 48 }}>加载中...</div>
  }

  return (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 24 }}>
          <img
            src={shop.logo || 'https://via.placeholder.com/120'}
            alt={shop.name}
            style={{ width: 120, height: 120, borderRadius: '50%', objectFit: 'cover' }}
          />
          <div style={{ flex: 1 }}>
            <Title level={3} style={{ marginBottom: 8 }}>{shop.name}</Title>
            <div style={{ color: '#666', marginBottom: 12 }}>{shop.description}</div>
            <div style={{ display: 'flex', gap: 24, marginBottom: 16 }}>
              <div>
                <Text type="secondary">好评率</Text>
                <div style={{ fontSize: 20, fontWeight: 'bold', color: '#52c41a' }}>
                  {shop.rating || 5}%
                </div>
              </div>
              <div>
                <Text type="secondary">商品数量</Text>
                <div style={{ fontSize: 20, fontWeight: 'bold' }}>{shop.productCount || 0}</div>
              </div>
              <div>
                <Text type="secondary">销量</Text>
                <div style={{ fontSize: 20, fontWeight: 'bold' }}>{shop.soldCount || 0}</div>
              </div>
            </div>
            <div style={{ display: 'flex', gap: 12 }}>
              <Button type="primary" icon={<ShopOutlined />}>
                联系商家
              </Button>
              <Button
                icon={isFavorite ? <HeartFilled style={{ color: '#ff4d4f' }} /> : <HeartOutlined />}
                onClick={toggleFavorite}
              >
                {isFavorite ? '已收藏' : '收藏店铺'}
              </Button>
            </div>
          </div>
        </div>
      </Card>

      <Title level={4} style={{ marginBottom: 16 }}>店铺商品</Title>
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
          <div style={{ textAlign: 'center', marginTop: 32 }}>
            <Pagination
              current={page}
              pageSize={pageSize}
              total={total}
              onChange={setPage}
            />
          </div>
        </>
      ) : (
        <Empty description="该店铺暂无商品" />
      )}
    </div>
  )
}

export default ShopDetail
