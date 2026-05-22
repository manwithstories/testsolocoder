import { useState, useEffect } from 'react'
import {
  Row, Col, Card, Image, Typography, Tag, Button, InputNumber,
  Descriptions, Rate, List, Avatar, message, Divider
} from 'antd'
import { ShoppingCartOutlined, HeartOutlined, HeartFilled, StarOutlined } from '@ant-design/icons'
import { Link, useParams, useNavigate } from 'react-router-dom'
import { productAPI, cartAPI, favoriteAPI, reviewAPI } from '@/api'
import { useAppStore } from '@/store'
import { Product, Review, SKU } from '@/types'

const { Title, Text } = Typography

const ProductDetail = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user, addToCart } = useAppStore()
  const [product, setProduct] = useState<Product | null>(null)
  const [reviews, setReviews] = useState<Review[]>([])
  const [skus, setSkus] = useState<SKU[]>([])
  const [selectedSku, setSelectedSku] = useState<SKU | null>(null)
  const [quantity, setQuantity] = useState(1)
  const [isFavorite, setIsFavorite] = useState(false)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (id) {
      loadProduct()
      loadReviews()
    }
  }, [id])

  const loadProduct = async () => {
    try {
      const res = await productAPI.getById(Number(id)) as any
      setProduct(res.data)
      // SKU信息在product详情中
      if (res.data.skus && res.data.skus.length > 0) {
        setSkus(res.data.skus)
        setSelectedSku(res.data.skus[0])
      }
    } catch (err) {
      console.error('加载商品失败', err)
    }
  }

  // SKU信息已包含在商品详情中，无需单独加载
  // const loadSkus = async () => {
  //   try {
  //     const res = await productAPI.getProductSkus(Number(id))
  //     setSkus(res.data)
  //     if (res.data.length > 0) {
  //       setSelectedSku(res.data[0])
  //     }
  //   } catch (err) {
  //     console.error('加载SKU失败', err)
  //   }
  // }

  const loadReviews = async () => {
    try {
      const res = await reviewAPI.getByProduct(Number(id), { page: 1, pageSize: 10 })
      setReviews(res.data.data)
    } catch (err) {
      console.error('加载评价失败', err)
    }
  }

  const handleAddToCart = async () => {
    if (!user) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    if (!selectedSku) {
      message.warning('请选择规格')
      return
    }
    setLoading(true)
    try {
      await cartAPI.add({ productId: product!.id, skuId: selectedSku.id, quantity })
      addToCart({ productId: product!.id, skuId: selectedSku.id, quantity })
      message.success('已加入购物车')
    } catch (err: any) {
      message.error(err.message || '加入购物车失败')
    } finally {
      setLoading(false)
    }
  }

  const handleBuyNow = async () => {
    if (!user) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    if (!selectedSku) {
      message.warning('请选择规格')
      return
    }
    setLoading(true)
    try {
      await cartAPI.add({ productId: product!.id, skuId: selectedSku.id, quantity })
      addToCart({ productId: product!.id, skuId: selectedSku.id, quantity })
      navigate('/cart')
    } catch (err: any) {
      message.error(err.message || '操作失败')
    } finally {
      setLoading(false)
    }
  }

  const toggleFavorite = async () => {
    if (!user) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    try {
      await favoriteAPI.toggleProduct(product!.id)
      setIsFavorite(!isFavorite)
      message.success(isFavorite ? '已取消收藏' : '已收藏')
    } catch (err: any) {
      message.error(err.message || '操作失败')
    }
  }

  if (!product) {
    return <div style={{ textAlign: 'center', padding: 48 }}>加载中...</div>
  }

  return (
    <div>
      <Row gutter={[24, 24]}>
        <Col xs={24} md={10}>
          <Card>
            <Image
              width="100%"
              src={product.images?.[0] || 'https://via.placeholder.com/500x500'}
              alt={product.name}
            />
            {product.images && product.images.length > 1 && (
              <div style={{ marginTop: 12, display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                {product.images.slice(1).map((img, idx) => (
                  <Image
                    key={idx}
                    width={80}
                    height={80}
                    src={img}
                    style={{ objectFit: 'cover', borderRadius: 4 }}
                  />
                ))}
              </div>
            )}
          </Card>
        </Col>

        <Col xs={24} md={14}>
          <Card>
            <Title level={3} style={{ marginBottom: 16 }}>{product.name}</Title>

            <div style={{ marginBottom: 16 }}>
              <span className="price" style={{ fontSize: 28 }}>
                ¥{selectedSku?.price || product.price}
              </span>
              {product.originalPrice && product.originalPrice > product.price && (
                <span className="original-price" style={{ fontSize: 18 }}>
                  ¥{product.originalPrice}
                </span>
              )}
              <Tag color="red" style={{ marginLeft: 12 }}>已售 {product.sales || 0}</Tag>
            </div>

            <Descriptions column={1} style={{ marginBottom: 16 }}>
              <Descriptions.Item label="店铺">
                <Link to={`/shops/${product.shopId}`} style={{ color: '#1890ff' }}>
                  {product.shopName}
                </Link>
              </Descriptions.Item>
              <Descriptions.Item label="库存">{selectedSku?.stock || product.stock || 0} 件</Descriptions.Item>
            </Descriptions>

            {skus.length > 0 && (
              <div style={{ marginBottom: 16 }}>
                <Text strong style={{ marginBottom: 8, display: 'block' }}>选择规格:</Text>
                <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap' }}>
                  {skus.map((sku) => (
                    <Button
                      key={sku.id}
                      type={selectedSku?.id === sku.id ? 'primary' : 'default'}
                      onClick={() => setSelectedSku(sku)}
                    >
                      {Object.values(sku.specs || {}).join(' / ') || `SKU-${sku.id}`}
                    </Button>
                  ))}
                </div>
              </div>
            )}

            <div style={{ marginBottom: 16 }}>
              <Text strong style={{ marginRight: 12 }}>数量:</Text>
              <InputNumber
                min={1}
                max={selectedSku?.stock || product.stock || 99}
                value={quantity}
                onChange={(val) => setQuantity(val || 1)}
              />
            </div>

            <div style={{ display: 'flex', gap: 12 }}>
              <Button
                type="primary"
                size="large"
                icon={<ShoppingCartOutlined />}
                onClick={handleAddToCart}
                loading={loading}
              >
                加入购物车
              </Button>
              <Button
                size="large"
                onClick={handleBuyNow}
                loading={loading}
                style={{ background: '#ff4d4f', borderColor: '#ff4d4f', color: '#fff' }}
              >
                立即购买
              </Button>
              <Button
                size="large"
                icon={isFavorite ? <HeartFilled style={{ color: '#ff4d4f' }} /> : <HeartOutlined />}
                onClick={toggleFavorite}
              >
                {isFavorite ? '已收藏' : '收藏'}
              </Button>
            </div>

            <Divider />

            <div style={{ color: '#666' }}>
              <Text strong>商品描述:</Text>
              <p style={{ marginTop: 8 }}>{product.description}</p>
            </div>
          </Card>
        </Col>
      </Row>

      <Card style={{ marginTop: 24 }} title="商品评价">
        {reviews.length > 0 ? (
          <List
            dataSource={reviews}
            renderItem={(review) => (
              <List.Item key={review.id}>
                <List.Item.Meta
                  avatar={<Avatar>{review.username?.[0]}</Avatar>}
                  title={
                    <div>
                      <span style={{ marginRight: 12 }}>{review.username}</span>
                      <Rate disabled value={review.rating} style={{ fontSize: 14 }} />
                    </div>
                  }
                  description={
                    <div>
                      <p>{review.content}</p>
                      {review.reply && (
                        <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, marginTop: 8 }}>
                          <Text type="secondary">商家回复:</Text>
                          <p style={{ marginTop: 4, marginBottom: 0 }}>{review.reply}</p>
                        </div>
                      )}
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        ) : (
          <div className="empty-state">暂无评价</div>
        )}
      </Card>
    </div>
  )
}

export default ProductDetail
