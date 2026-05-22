import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Space, Modal, Input, Form, message, Popconfirm, Select, DatePicker } from 'antd'
import { orderApi } from '@/services'
import { Order } from '@/types'
import dayjs from 'dayjs'

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'orange', text: '待支付' },
  paid: { color: 'green', text: '已支付' },
  refunding: { color: 'gold', text: '退款中' },
  refunded: { color: 'default', text: '已退款' },
  cancelled: { color: 'default', text: '已取消' },
  failed: { color: 'red', text: '支付失败' },
}

const AdminOrders: React.FC = () => {
  const [data, setData] = useState<Order[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [loading, setLoading] = useState(false)
  const [statusFilter, setStatusFilter] = useState('')
  const [search, setSearch] = useState('')
  const [refundModalVisible, setRefundModalVisible] = useState(false)
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null)
  const [refundForm] = Form.useForm()

  const loadOrders = async (p = page) => {
    setLoading(true)
    try {
      const params: any = { page: p, page_size: pageSize }
      if (statusFilter) params.status = statusFilter
      if (search) params.search = search
      const res = await orderApi.listAll(params)
      if (res.code === 0 && res.data) {
        setData(res.data.items)
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

  const handleProcessRefund = async (approved: boolean) => {
    if (!selectedOrder) return
    try {
      const values = await refundForm.validateFields()
      const res = await orderApi.processRefund(selectedOrder.id, { ...values, approved })
      if (res.code === 0) {
        message.success('退款处理成功')
        setRefundModalVisible(false)
        loadOrders()
      }
    } catch (error: any) {
      message.error(error.message || '处理失败')
    }
  }

  const columns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    {
      title: '课程',
      key: 'course',
      render: (_: any, record: Order) => record.course?.title || '-',
    },
    {
      title: '学员',
      key: 'user',
      render: (_: any, record: Order) => record.user?.username || '-',
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
      title: '支付时间',
      dataIndex: 'paid_at',
      key: 'paid_at',
      render: (val: string) => val ? dayjs(val).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        <Space>
          {record.status === 'refunding' && (
            <Button
              type="link"
              onClick={() => {
                setSelectedOrder(record)
                setRefundModalVisible(true)
              }}
            >
              处理退款
            </Button>
          )}
          {record.status === 'pending' && (
            <Popconfirm
              title="确定取消订单？"
              onConfirm={async () => {
                try {
                  await orderApi.updateStatus(record.id, 'cancelled')
                  message.success('订单已取消')
                  loadOrders()
                } catch (error: any) {
                  message.error(error.message || '操作失败')
                }
              }}
            >
              <Button type="link" danger>取消</Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <h2>订单管理</h2>
      <Space style={{ marginBottom: 16 }}>
        <Select
          placeholder="全部状态"
          style={{ width: 140 }}
          allowClear
          value={statusFilter || undefined}
          onChange={(val) => { setStatusFilter(val || ''); loadOrders(1) }}
        >
          {Object.entries(statusMap).map(([key, val]) => (
            <Select.Option key={key} value={key}>{val.text}</Select.Option>
          ))}
        </Select>
        <Input.Search
          placeholder="搜索订单号/学员"
          style={{ width: 200 }}
          onSearch={(val) => { setSearch(val); loadOrders(1) }}
          allowClear
        />
      </Space>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={data}
        loading={loading}
        pagination={{
          current: page,
          total,
          pageSize,
          onChange: (p) => { setPage(p); loadOrders(p) },
        }}
      />
      <Modal
        title="处理退款"
        open={refundModalVisible}
        onCancel={() => setRefundModalVisible(false)}
        footer={[
          <Button key="reject" danger onClick={() => handleProcessRefund(false)}>
            拒绝退款
          </Button>,
          <Button key="approve" type="primary" onClick={() => handleProcessRefund(true)}>
            同意退款
          </Button>,
        ]}
      >
        {selectedOrder && (
          <div>
            <p>订单号：{selectedOrder.order_no}</p>
            <p>课程：{selectedOrder.course?.title}</p>
            <p>金额：¥{selectedOrder.amount.toFixed(2)}</p>
            <p>退款原因：{selectedOrder.refund_reason}</p>
          </div>
        )}
        <Form form={refundForm} layout="vertical">
          <Form.Item name="remark" label="处理备注">
            <Input.TextArea rows={3} placeholder="请输入处理备注..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AdminOrders
