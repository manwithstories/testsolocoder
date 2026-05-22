import { useState, useEffect } from 'react'
import {
  Table, Button, Card, Typography, Tag, Space,
  Pagination, message, Modal, Form, Input
} from 'antd'
import { EyeOutlined, TruckOutlined } from '@ant-design/icons'
import { orderAPI } from '@/api'
import { Order } from '@/types'

const { Title } = Typography

const SellerOrders = () => {
  const [orders, setOrders] = useState<Order[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [status, setStatus] = useState<string>('')
  const [shipModalVisible, setShipModalVisible] = useState(false)
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null)
  const [shipForm] = Form.useForm()

  useEffect(() => {
    loadOrders()
  }, [page, status])

  const loadOrders = async () => {
    try {
      const params: any = { page, pageSize }
      if (status) params.status = status
      const res = await orderAPI.list(params)
      setOrders(res.data.data)
      setTotal(res.data.pagination.total)
    } catch (err) {
      console.error('加载订单失败', err)
    }
  }

  const handleShip = (order: Order) => {
    setCurrentOrder(order)
    setShipModalVisible(true)
  }

  const handleShipSubmit = async (values: any) => {
    if (!currentOrder) return
    try {
      await orderAPI.ship(currentOrder.id, {
        trackingCompany: values.trackingCompany,
        trackingNo: values.trackingNo
      })
      message.success('发货成功')
      setShipModalVisible(false)
      shipForm.resetFields()
      loadOrders()
    } catch (err: any) {
      message.error(err.message || '发货失败')
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
    }
    const info = statusMap[status] || { color: 'default', text: status }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const columns = [
    {
      title: '订单号',
      dataIndex: 'orderNo',
      key: 'orderNo',
      width: 180,
    },
    {
      title: '商品',
      key: 'items',
      width: 250,
      render: (_: any, record: Order) => (
        <div>
          {record.items?.map((item: any) => (
            <div key={item.id} style={{ display: 'flex', gap: 8, marginBottom: 4 }}>
              <img src={item.productImage} style={{ width: 40, height: 40, objectFit: 'cover' }} />
              <div>
                <div className="truncate" style={{ maxWidth: 160 }}>{item.productName}</div>
                <div style={{ color: '#999', fontSize: 12 }}>
                  {item.specs && Object.values(item.specs).join(' / ')} x{item.quantity}
                </div>
              </div>
            </div>
          ))}
        </div>
      ),
    },
    {
      title: '金额',
      dataIndex: 'totalAmount',
      key: 'totalAmount',
      width: 120,
      render: (amount: number) => <span className="price">¥{amount.toFixed(2)}</span>,
    },
    {
      title: '买家',
      dataIndex: 'userName',
      key: 'userName',
      width: 120,
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
      width: 150,
      render: (_: any, record: Order) => (
        <Space>
          {record.status === 'pending_ship' && (
            <Button
              type="primary"
              size="small"
              icon={<TruckOutlined />}
              onClick={() => handleShip(record)}
            >
              发货
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Title level={3}>订单管理</Title>

      <Card>
        <div style={{ marginBottom: 16 }}>
          <Button.Group>
            <Button type={status === '' ? 'primary' : 'default'} onClick={() => { setStatus(''); setPage(1) }}>
              全部
            </Button>
            <Button type={status === 'pending_ship' ? 'primary' : 'default'} onClick={() => { setStatus('pending_ship'); setPage(1) }}>
              待发货
            </Button>
            <Button type={status === 'shipped' ? 'primary' : 'default'} onClick={() => { setStatus('shipped'); setPage(1) }}>
              已发货
            </Button>
            <Button type={status === 'completed' ? 'primary' : 'default'} onClick={() => { setStatus('completed'); setPage(1) }}>
              已完成
            </Button>
          </Button.Group>
        </div>

        <Table
          columns={columns}
          dataSource={orders}
          rowKey="id"
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
      </Card>

      <Modal
        title="订单发货"
        open={shipModalVisible}
        onCancel={() => setShipModalVisible(false)}
        footer={null}
      >
        <Form form={shipForm} layout="vertical" onFinish={handleShipSubmit}>
          <Form.Item
            name="trackingCompany"
            label="物流公司"
            rules={[{ required: true, message: '请输入物流公司' }]}
          >
            <Input placeholder="请输入物流公司" />
          </Form.Item>
          <Form.Item
            name="trackingNo"
            label="物流单号"
            rules={[{ required: true, message: '请输入物流单号' }]}
          >
            <Input placeholder="请输入物流单号" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
              确认发货
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default SellerOrders
