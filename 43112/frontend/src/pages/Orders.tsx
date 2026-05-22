import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Modal, Input, Form, message, Popconfirm } from 'antd'
import { orderApi } from '@/services'
import { Order } from '@/types'

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'orange', text: '待支付' },
  paid: { color: 'green', text: '已支付' },
  refunding: { color: 'gold', text: '退款中' },
  refunded: { color: 'default', text: '已退款' },
  cancelled: { color: 'default', text: '已取消' },
  failed: { color: 'red', text: '支付失败' },
}

const OrdersPage: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [loading, setLoading] = useState(false)
  const [refundModalVisible, setRefundModalVisible] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null)
  const [refundForm] = Form.useForm()

  const loadOrders = async (p = page) => {
    setLoading(true)
    try {
      const res = await orderApi.myOrders({ page: p, page_size: pageSize })
      if (res.code === 0 && res.data) {
        setOrders(res.data.items)
        setTotal(res.data.total)
      }
    } catch (error) {
      console.error('Failed to load orders:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadOrders()
  }, [])

  const handleRefund = async () => {
    if (!selectedOrder) return
    try {
      const values = await refundForm.validateFields()
      const res = await orderApi.refund(selectedOrder.id, values.reason)
      if (res.code === 0) {
        message.success('退款申请已提交')
        setRefundModalVisible(false)
        loadOrders()
      }
    } catch (error: any) {
      message.error(error.message || '退款申请失败')
    }
  }

  const columns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
    },
    {
      title: '课程',
      key: 'course',
      render: (_: any, record: Order) => record.course?.title || '-',
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (val: number) => <span style={{ color: '#f5222d', fontWeight: 600 }}>¥{val.toFixed(2)}</span>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const info = statusMap[status] || { color: 'default', text: status }
        return <Tag color={info.color}>{info.text}</Tag>
      },
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => new Date(val).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        record.status === 'paid' ? (
          <Popconfirm
            title="确定申请退款？"
            onConfirm={() => {
              setSelectedOrder(record)
              setRefundModalVisible(true)
            }}
          >
            <Button type="link" danger>申请退款</Button>
          </Popconfirm>
        ) : null
      ),
    },
  ]

  return (
    <div>
      <h2>我的订单</h2>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={orders}
        loading={loading}
        pagination={{
          current: page,
          total,
          pageSize,
          onChange: (p) => {
            setPage(p)
            loadOrders(p)
          },
        }}
      />
      <Modal
        title="申请退款"
        open={refundModalVisible}
        onCancel={() => setRefundModalVisible(false)}
        onOk={handleRefund}
        okText="提交申请"
      >
        <Form form={refundForm} layout="vertical">
          <Form.Item name="reason" label="退款原因" rules={[{ required: true, message: '请输入退款原因' }]}>
            <Input.TextArea rows={4} placeholder="请说明退款原因..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default OrdersPage
