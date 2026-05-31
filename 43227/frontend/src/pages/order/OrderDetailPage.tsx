import { useEffect, useState } from 'react'
import { Card, Descriptions, Button, Tag, Space, Modal, Form, Input, Rate, InputNumber, message } from 'antd'
import { ArrowLeftOutlined, PayCircleOutlined, TruckOutlined, CheckCircleOutlined, StarOutlined } from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../../api'
import type { Order } from '../../types'
import { useAuthStore } from '../../store/authStore'
import dayjs from 'dayjs'

function OrderDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<Order | null>(null)
  const [isRateModalVisible, setIsRateModalVisible] = useState(false)
  const [isShipModalVisible, setIsShipModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [shipForm] = Form.useForm()

  useEffect(() => {
    if (id) {
      fetchData()
    }
  }, [id])

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/orders/${id}`)
      setData(response.data as Order)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const handlePay = async () => {
    try {
      await api.put(`/orders/${id}/pay`)
      message.success('支付成功')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '支付失败')
    }
  }

  const handleShip = async () => {
    try {
      const values = await shipForm.validateFields()
      await api.put(`/orders/${id}/ship`, values)
      message.success('发货成功')
      setIsShipModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '发货失败')
    }
  }

  const handleDeliver = async () => {
    try {
      await api.put(`/orders/${id}/deliver`)
      message.success('确认收货成功')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleComplete = async () => {
    try {
      await api.put(`/orders/${id}/complete`)
      message.success('订单已完成')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleRate = async () => {
    try {
      const values = await form.validateFields()
      await api.put(`/orders/${id}/rate`, values)
      message.success('评价成功')
      setIsRateModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '评价失败')
    }
  }

  const handleCancel = async () => {
    try {
      await api.put(`/orders/${id}/cancel`)
      message.success('订单已取消')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '取消失败')
    }
  }

  const statusColorMap: Record<string, string> = {
    pending: 'default',
    paid: 'blue',
    shipped: 'cyan',
    delivered: 'geekblue',
    completed: 'green',
    cancelled: 'red',
    refunded: 'orange',
  }

  if (loading) {
    return <Card loading />
  }

  if (!data) {
    return <Card>订单不存在</Card>
  }

  const canCancel = data.status === 'pending' || data.status === 'paid'
  const canPay = data.status === 'pending'
  const canShip = data.status === 'paid' && user?.role === 'beekeeper'
  const canDeliver = data.status === 'shipped' && user?.role === 'buyer'
  const canComplete = data.status === 'delivered' && user?.role === 'buyer'
  const canRate = data.status === 'completed' && user?.role === 'buyer' && data.buyer_rating == null

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(-1)}>
          返回
        </Button>
      </Space>
      <Card title={`订单详情 - ${data.order_no}`}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="订单状态">
            <Tag color={statusColorMap[data.status]}>{data.status}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="支付状态">{data.payment_status}</Descriptions.Item>
          <Descriptions.Item label="买家">{data.buyer?.username}</Descriptions.Item>
          <Descriptions.Item label="卖家">{data.seller?.username}</Descriptions.Item>
          <Descriptions.Item label="商品">{data.product?.title}</Descriptions.Item>
          <Descriptions.Item label="蜂蜜类型">{data.product?.honey_type}</Descriptions.Item>
          <Descriptions.Item label="数量">{data.quantity} {data.product?.unit}</Descriptions.Item>
          <Descriptions.Item label="单价">¥{data.unit_price}/{data.product?.unit}</Descriptions.Item>
          <Descriptions.Item label="总金额" span={2}>
            <span style={{ color: '#f5222d', fontWeight: 'bold', fontSize: 20 }}>¥{data.total_amount}</span>
          </Descriptions.Item>
          <Descriptions.Item label="收货地址" span={2}>{data.shipping_address}</Descriptions.Item>
          {data.tracking_number && (
            <Descriptions.Item label="物流单号">{data.tracking_number}</Descriptions.Item>
          )}
          {data.tracking_status && (
            <Descriptions.Item label="物流状态">{data.tracking_status}</Descriptions.Item>
          )}
          <Descriptions.Item label="创建时间">{dayjs(data.created_at).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
          {data.payment_time && (
            <Descriptions.Item label="支付时间">{dayjs(data.payment_time).format('YYYY-MM-DD HH:mm:ss')}</Descriptions.Item>
          )}
          {data.buyer_rating != null && (
            <Descriptions.Item label="买家评价">
              <Space>
                <Rate disabled allowHalf value={data.buyer_rating} />
                {data.buyer_comment && <span>{data.buyer_comment}</span>}
              </Space>
            </Descriptions.Item>
          )}
          {data.seller_rating != null && (
            <Descriptions.Item label="卖家评价">
              <Space>
                <Rate disabled allowHalf value={data.seller_rating} />
                {data.seller_comment && <span>{data.seller_comment}</span>}
              </Space>
            </Descriptions.Item>
          )}
        </Descriptions>

        <Space style={{ marginTop: 24 }}>
          {canPay && (
            <Button type="primary" icon={<PayCircleOutlined />} onClick={handlePay}>
              立即支付
            </Button>
          )}
          {canShip && (
            <Button type="primary" icon={<TruckOutlined />} onClick={() => setIsShipModalVisible(true)}>
              发货
            </Button>
          )}
          {canDeliver && (
            <Button type="primary" icon={<CheckCircleOutlined />} onClick={handleDeliver}>
              确认收货
            </Button>
          )}
          {canComplete && (
            <Button type="primary" icon={<CheckCircleOutlined />} onClick={handleComplete}>
              完成订单
            </Button>
          )}
          {canRate && (
            <Button type="primary" icon={<StarOutlined />} onClick={() => setIsRateModalVisible(true)}>
              评价
            </Button>
          )}
          {canCancel && (
            <Button danger onClick={handleCancel}>
              取消订单
            </Button>
          )}
        </Space>
      </Card>

      <Modal
        title="评价订单"
        open={isRateModalVisible}
        onCancel={() => setIsRateModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleRate} layout="vertical">
          <Form.Item name="rating" label="评分" rules={[{ required: true, message: '请评分' }]}>
            <Rate allowHalf />
          </Form.Item>
          <Form.Item name="comment" label="评价内容">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交评价</Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="发货"
        open={isShipModalVisible}
        onCancel={() => setIsShipModalVisible(false)}
        footer={null}
      >
        <Form form={shipForm} onFinish={handleShip} layout="vertical">
          <Form.Item name="tracking_number" label="物流单号" rules={[{ required: true, message: '请输入物流单号' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="tracking_status" label="物流状态">
            <Input placeholder="如：已揽收" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>确认发货</Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default OrderDetailPage
