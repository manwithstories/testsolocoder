import { useState, useEffect } from 'react'
import {
  Table, Button, Card, Typography, Tabs, Tag, Empty, message,
  Popconfirm, Space, Pagination
} from 'antd'
import { Link } from 'react-router-dom'
import { orderAPI } from '@/api'
import { Order } from '@/types'

const { Title } = Typography
const { TabPane } = Tabs

const OrderList = () => {
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [status, setStatus] = useState<string>('')

  useEffect(() => {
    loadOrders()
  }, [page, status])

  const loadOrders = async () => {
    setLoading(true)
    try {
      const params: any = { page, pageSize }
      if (status) params.status = status
      const res = await orderAPI.list(params)
      setOrders(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载订单失败', err)
    } finally {
      setLoading(false)
    }
  }

  const handlePay = async (id: number) => {
    try {
      await orderAPI.pay(id)
      message.success('支付成功')
      loadOrders()
    } catch (err: any) {
      message.error(err.message || '支付失败')
    }
  }

  const handleCancel = async (id: number) => {
    try {
      await orderAPI.cancel(id)
      message.success('已取消订单')
      loadOrders()
    } catch (err: any) {
      message.error(err.message || '取消失败')
    }
  }

  const handleConfirm = async (id: number) => {
    try {
      await orderAPI.confirm(id)
      message.success('已确认收货')
      loadOrders()
    } catch (err: any) {
      message.error(err.message || '操作失败')
    }
  }

  const handleRefund = async (orderId: number, orderItemId?: number) => {
    try {
      await orderAPI.applyRefund({
        orderId,
        orderItemId,
        reason: '申请退款',
        type: 'refund',
      })
      message.success('退款申请已提交')
      loadOrders()
    } catch (err: any) {
      message.error(err.message || '申请失败')
    }
  }

  const getStatusTag = (status: string) => {
    const statusMap: Record<string, { color: string; text: string }> = {
      pending_payment: { color: 'orange', text: '待付款' },
      pending_ship: { color: 'blue', text: '待发货' },
      shipped: { color: 'cyan', text: '已发货' },
      completed: { color: 'green', text: '已完成' },
      cancelled: { color: 'default', text: '已取消' },
      refunded: { color: 'red', text: '已退款' },
      refund_pending: { color: 'gold', text: '退款中' },
    }
    const info = statusMap[status] || { color: 'default', text: status }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const columns = [
    {
      title: '订单信息',
      key: 'orderInfo',
      width: 300,
      render: (_: any, record: Order) => (
        <Link to={`/orders/${record.id}`}>
          <div style={{ display: 'flex', gap: 12 }}>
            <img
              src={record.items?.[0]?.productImage || 'https://via.placeholder.com/80'}
              style={{ width: 60, height: 60, objectFit: 'cover', borderRadius: 4 }}
            />
            <div>
              <div className="truncate" style={{ maxWidth: 180 }}>
                {record.items?.[0]?.productName}
              </div>
              {record.items && record.items.length > 1 && (
                <div style={{ color: '#999', fontSize: 12 }}>等{record.items.length}件商品</div>
              )}
              <div style={{ color: '#999', fontSize: 12, marginTop: 4 }}>
                订单号: {record.orderNo}
              </div>
            </div>
          </div>
        </Link>
      ),
    },
    {
      title: '店铺',
      dataIndex: 'shopName',
      key: 'shopName',
      width: 150,
    },
    {
      title: '金额',
      key: 'amount',
      width: 120,
      render: (_: any, record: Order) => (
        <span className="price">¥{record.totalAmount.toFixed(2)}</span>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status: string) => getStatusTag(status),
    },
    {
      title: '下单时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      width: 180,
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_: any, record: Order) => (
        <Space>
          <Link to={`/orders/${record.id}`}>
            <Button size="small">查看详情</Button>
          </Link>
          {record.status === 'pending_payment' && (
            <>
              <Button type="primary" size="small" onClick={() => handlePay(record.id)}>
                去支付
              </Button>
              <Popconfirm title="确定取消订单？" onConfirm={() => handleCancel(record.id)}>
                <Button size="small" danger>取消</Button>
              </Popconfirm>
            </>
          )}
          {record.status === 'shipped' && (
            <Button type="primary" size="small" onClick={() => handleConfirm(record.id)}>
              确认收货
            </Button>
          )}
          {record.status === 'completed' && (
            <Popconfirm title="确定申请退款？" onConfirm={() => handleRefund(record.id, record.items?.[0]?.id)}>
              <Button size="small" danger>申请退款</Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Title level={3}>我的订单</Title>

      <Card>
        <Tabs activeKey={status} onChange={(key) => { setStatus(key); setPage(1) }}>
          <TabPane tab="全部" key="" />
          <TabPane tab="待付款" key="pending_payment" />
          <TabPane tab="待发货" key="pending_ship" />
          <TabPane tab="待收货" key="shipped" />
          <TabPane tab="已完成" key="completed" />
          <TabPane tab="退款" key="refunded" />
        </Tabs>

        {orders.length > 0 ? (
          <>
            <Table
              columns={columns}
              dataSource={orders}
              rowKey="id"
              loading={loading}
              pagination={false}
            />
            <div style={{ textAlign: 'center', marginTop: 24 }}>
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={setPage}
              />
            </div>
          </>
        ) : (
          <Empty description="暂无订单" />
        )}
      </Card>
    </div>
  )
}

export default OrderList
