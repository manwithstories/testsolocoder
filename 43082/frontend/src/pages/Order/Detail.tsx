import { useState, useEffect } from 'react'
import {
  Card, Typography, Descriptions, Table, Tag, Button, message,
  Steps, Space, Divider
} from 'antd'
import { ArrowLeftOutlined, CheckOutlined, ClockCircleOutlined } from '@ant-design/icons'
import { Link, useParams, useNavigate } from 'react-router-dom'
import { orderAPI } from '@/api'
import { Order } from '@/types'

const { Title, Text } = Typography

const OrderDetail = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (id) {
      loadOrder()
    }
  }, [id])

  const loadOrder = async () => {
    try {
      const res = await orderAPI.getById(Number(id)) as any
      setOrder(res.data)
    } catch (err) {
      console.error('加载订单失败', err)
    }
  }

  const handlePay = async () => {
    if (!order) return
    setLoading(true)
    try {
      await orderAPI.pay(order.id)
      message.success('支付成功')
      loadOrder()
    } catch (err: any) {
      message.error(err.message || '支付失败')
    } finally {
      setLoading(false)
    }
  }

  const handleConfirm = async () => {
    if (!order) return
    setLoading(true)
    try {
      await orderAPI.confirm(order.id)
      message.success('已确认收货')
      loadOrder()
    } catch (err: any) {
      message.error(err.message || '操作失败')
    } finally {
      setLoading(false)
    }
  }

  const getStatusSteps = (status: string) => {
    const steps = [
      { title: '提交订单', status: 'finish' as const },
      { title: '支付', status: status === 'pending_payment' ? 'process' as const : 'finish' as const },
      { title: '发货', status: ['shipped', 'completed', 'refunded'].includes(status) ? 'finish' as const : (status === 'pending_ship' ? 'process' as const : 'wait' as const) },
      { title: '收货', status: ['completed', 'refunded'].includes(status) ? 'finish' as const : (status === 'shipped' ? 'process' as const : 'wait' as const) },
      { title: '完成', status: status === 'completed' ? 'finish' as const : 'wait' as const },
    ]
    return steps
  }

  const itemColumns = [
    {
      title: '商品',
      dataIndex: 'productName',
      key: 'productName',
      render: (text: string, record: any) => (
        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <img
            src={record.productImage}
            style={{ width: 60, height: 60, objectFit: 'cover', borderRadius: 4 }}
          />
          <div>
            <div>{text}</div>
            {record.specs && Object.values(record.specs).length > 0 && (
              <Tag color="blue" style={{ marginTop: 4 }}>
                {Object.values(record.specs).join(' / ')}
              </Tag>
            )}
          </div>
        </div>
      ),
    },
    {
      title: '单价',
      dataIndex: 'price',
      key: 'price',
      width: 120,
      render: (price: number) => <span>¥{price.toFixed(2)}</span>,
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      width: 100,
    },
    {
      title: '小计',
      key: 'subtotal',
      width: 120,
      render: (_: any, record: any) => <span className="price">¥{(record.price * record.quantity).toFixed(2)}</span>,
    },
  ]

  if (!order) {
    return <div style={{ textAlign: 'center', padding: 48 }}>加载中...</div>
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate('/orders')}
        >
          返回订单列表
        </Button>
      </div>

      <Card>
        <div className="page-header">
          <Title level={3} style={{ margin: 0 }}>
            订单详情
            <Tag style={{ marginLeft: 12 }}>{order.orderNo}</Tag>
          </Title>
        </div>

        <Steps
          current={getStatusSteps(order.status).findIndex(s => s.status === 'process')}
          items={getStatusSteps(order.status)}
          style={{ marginBottom: 32 }}
        />

        <Descriptions column={2} style={{ marginBottom: 24 }}>
          <Descriptions.Item label="订单状态">
            <Tag color={order.status === 'completed' ? 'green' : order.status === 'pending_payment' ? 'orange' : 'blue'}>
              {order.statusText}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="下单时间">{order.createdAt}</Descriptions.Item>
          <Descriptions.Item label="店铺">{order.shopName}</Descriptions.Item>
          <Descriptions.Item label="支付方式">{order.paymentMethod || '在线支付'}</Descriptions.Item>
        </Descriptions>

        <Divider />

        <Title level={5} style={{ marginBottom: 16 }}>商品信息</Title>
        <Table
          columns={itemColumns}
          dataSource={order.items}
          rowKey="id"
          pagination={false}
          style={{ marginBottom: 24 }}
        />

        <Divider />

        <div style={{ textAlign: 'right' }}>
          <div style={{ marginBottom: 8 }}>
            <Text type="secondary">商品金额: </Text>
            <Text>¥{order.totalAmount.toFixed(2)}</Text>
          </div>
          <div style={{ marginBottom: 8 }}>
            <Text type="secondary">运费: </Text>
            <Text>¥{order.shippingFee?.toFixed(2) || '0.00'}</Text>
          </div>
          <div style={{ fontSize: 18, fontWeight: 'bold' }}>
            <Text type="secondary">实付金额: </Text>
            <span className="price">¥{order.totalAmount.toFixed(2)}</span>
          </div>
        </div>

        <Divider />

        <Descriptions column={2} title="收货信息">
          <Descriptions.Item label="收货人">{order.receiverName}</Descriptions.Item>
          <Descriptions.Item label="联系电话">{order.receiverPhone}</Descriptions.Item>
          <Descriptions.Item label="收货地址" span={2}>{order.receiverAddress}</Descriptions.Item>
          {order.trackingCompany && (
            <>
              <Descriptions.Item label="物流公司">{order.trackingCompany}</Descriptions.Item>
              <Descriptions.Item label="物流单号">{order.trackingNo}</Descriptions.Item>
            </>
          )}
        </Descriptions>

        {order.status === 'pending_payment' && (
          <div style={{ textAlign: 'right', marginTop: 24 }}>
            <Space>
              <Button onClick={() => navigate('/orders')}>取消</Button>
              <Button type="primary" size="large" onClick={handlePay} loading={loading}>
                立即支付
              </Button>
            </Space>
          </div>
        )}

        {order.status === 'shipped' && (
          <div style={{ textAlign: 'right', marginTop: 24 }}>
            <Button type="primary" size="large" onClick={handleConfirm} loading={loading}>
              确认收货
            </Button>
          </div>
        )}
      </Card>
    </div>
  )
}

export default OrderDetail
