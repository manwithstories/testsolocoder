import { useState, useEffect } from 'react'
import { Row, Col, Card, Tabs, Empty, Button, message } from 'antd'
import { HeartOutlined, DeleteOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { favoriteAPI } from '@/api'
import { Favorite } from '@/types'

const { TabPane } = Tabs

const Favorites = () => {
  const [products, setProducts] = useState<Favorite[]>([])
  const [shops, setShops] = useState<Favorite[]>([])
  const [activeTab, setActiveTab] = useState('product')

  useEffect(() => {
    loadFavorites()
  }, [])

  const loadFavorites = async () => {
    try {
      const [productRes, shopRes] = await Promise.all([
        favoriteAPI.getProducts() as any,
        favoriteAPI.getShops() as any,
      ])
      setProducts(productRes.data.data)
      setShops(shopRes.data.data)
    } catch (err) {
      console.error('加载收藏失败', err)
    }
  }

  const handleRemove = async (id: number) => {
    try {
      if (activeTab === 'product') {
        await favoriteAPI.toggleProduct(id)
      } else {
        await favoriteAPI.toggleShop(id)
      }
      message.success('已取消收藏')
      loadFavorites()
    } catch (err: any) {
      message.error(err.message || '操作失败')
    }
  }

  return (
    <div>
      <Card>
        <Tabs activeKey={activeTab} onChange={setActiveTab}>
          <TabPane tab="商品收藏" key="product">
            {products.length > 0 ? (
              <Row gutter={[16, 16]}>
                {products.map((item) => (
                  <Col xs={12} sm={8} md={6} key={item.targetId}>
                    <Card hoverable cover={
                      <Link to={`/products/${item.targetId}`}>
                        <img
                          alt=""
                          src={item.product?.images?.[0] || 'https://via.placeholder.com/400x300'}
                          className="product-image"
                        />
                      </Link>
                    }>
                      <Card.Meta
                        title={<div className="product-name">{item.product?.name}</div>}
                        description={
                          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                            <span className="price">¥{item.product?.price}</span>
                            <Button
                              type="text"
                              danger
                              icon={<DeleteOutlined />}
                              onClick={() => handleRemove(item.targetId)}
                            />
                          </div>
                        }
                      />
                    </Card>
                  </Col>
                ))}
              </Row>
            ) : (
              <Empty description="暂无收藏商品" />
            )}
          </TabPane>

          <TabPane tab="店铺收藏" key="shop">
            {shops.length > 0 ? (
              <Row gutter={[16, 16]}>
                {shops.map((item) => (
                  <Col xs={24} sm={12} md={8} key={item.targetId}>
                    <Link to={`/shops/${item.targetId}`}>
                      <Card hoverable>
                        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                          <img
                            src={item.shop?.logo || 'https://via.placeholder.com/80'}
                            alt={item.shop?.name}
                            className="shop-logo"
                          />
                          <div style={{ flex: 1 }}>
                            <div style={{ fontSize: 16, fontWeight: 'bold', marginBottom: 4 }}>
                              {item.shop?.name}
                            </div>
                            <div style={{ color: '#999', fontSize: 12 }}>
                              {item.shop?.description}
                            </div>
                            <Button
                              type="text"
                              danger
                              icon={<DeleteOutlined />}
                              style={{ padding: 0, marginTop: 8 }}
                              onClick={(e) => {
                                e.preventDefault()
                                handleRemove(item.targetId)
                              }}
                            >
                              取消收藏
                            </Button>
                          </div>
                        </div>
                      </Card>
                    </Link>
                  </Col>
                ))}
              </Row>
            ) : (
              <Empty description="暂无收藏店铺" />
            )}
          </TabPane>
        </Tabs>
      </Card>
    </div>
  )
}

export default Favorites
