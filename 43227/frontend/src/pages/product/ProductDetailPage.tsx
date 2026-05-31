import { useEffect, useState } from 'react'
import { Card, Descriptions, Button, Image, InputNumber, Form, Modal, message, Space, Rate, Tag } from 'antd'
import { ArrowLeftOutlined, ShoppingCartOutlined } from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../../api'
import type { Product } from '../../types'
import { useAuthStore } from '../../store/authStore'

function ProductDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<Product | null>(null)
  const [isBuyModalVisible, setIsBuyModalVisible] = useState(false)
  const [form] = Form.useForm()

  useEffect(() => {
    if (id) {
      fetchData()
    }
  }, [id])

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/products/${id}`)
      setData(response.data as Product)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const handleBuy = () => {
    if (!user || user.role !== 'buyer') {
      message.warning('只有买家角色可以购买商品')
      return
    }
    form.resetFields()
    setIsBuyModalVisible(true)
  }

  const handleSubmitOrder = async (values: any) => {
    try {
      await api.post('/orders', {
        product_id: data?.id,
        quantity: values.quantity,
        shipping_address: values.shipping_address,
      })
      message.success('订单创建成功')
      setIsBuyModalVisible(false)
      navigate('/orders')
    } catch (error: any) {
      message.error(error.message || '创建订单失败')
    }
  }

  if (loading) {
    return <Card loading />
  }

  if (!data) {
    return <Card>商品不存在</Card>
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(-1)}>
          返回
        </Button>
      </Space>
      <Card>
        <div style={{ display: 'flex', gap: 24 }}>
          <div style={{ width: 400 }}>
            {data.images && data.images.length > 0 ? (
              <Image src={data.images[0]} style={{ width: '100%', maxHeight: 400, objectFit: 'cover' }} />
            ) : (
              <div style={{ width: 400, height: 400, background: '#f5f5f5', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 18 }}>
                无图片
              </div>
            )}
          </div>
          <div style={{ flex: 1 }}>
            <h1 style={{ marginBottom: 16 }}>{data.title}</h1>
            <Descriptions column={1} size="small" style={{ marginBottom: 16 }}>
              <Descriptions.Item label="蜂蜜类型">{data.honey_type}</Descriptions.Item>
              <Descriptions.Item label="批次号">{data.batch_code || '-'}</Descriptions.Item>
              <Descriptions.Item label="等级">
                {data.grade ? <Tag color={data.grade === 'grade_a' ? 'gold' : data.grade === 'grade_b' ? 'blue' : 'green'}>{data.grade}</Tag> : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="库存">{data.stock} {data.unit}</Descriptions.Item>
              <Descriptions.Item label="浏览次数">{data.view_count}</Descriptions.Item>
            </Descriptions>
            <div style={{ fontSize: 32, color: '#f5222d', fontWeight: 'bold', marginBottom: 16 }}>
              ¥{data.price}/{data.unit}
            </div>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="卖家">
                <Space>
                  <span>{data.user?.username}</span>
                  {data.user?.reputation != null && <Rate disabled allowHalf value={data.user.reputation} style={{ fontSize: 14 }} />}
                </Space>
              </Descriptions.Item>
            </Descriptions>
            {data.description && (
              <div style={{ marginTop: 16 }}>
                <h3>商品描述</h3>
                <p style={{ whiteSpace: 'pre-wrap' }}>{data.description}</p>
              </div>
            )}
            <Button
              type="primary"
              size="large"
              icon={<ShoppingCartOutlined />}
              onClick={handleBuy}
              style={{ marginTop: 16 }}
              disabled={data.stock <= 0}
            >
              立即购买
            </Button>
          </div>
        </div>
      </Card>
      <Modal
        title="确认购买"
        open={isBuyModalVisible}
        onCancel={() => setIsBuyModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleSubmitOrder} layout="vertical">
          <Form.Item name="quantity" label="购买数量" rules={[{ required: true, message: '请输入购买数量' }]}>
            <InputNumber min={0.1} max={data?.stock} step={0.1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="shipping_address" label="收货地址" rules={[{ required: true, message: '请输入收货地址' }]}>
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              提交订单
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default ProductDetailPage
