import { useState, useEffect } from 'react'
import { Table, Button, Tag, Space, Modal, message, Select, Card } from 'antd'
import { CheckOutlined, CloseOutlined, CompleteOutlined } from '@ant-design/icons'
import { orderApi } from '@/api'
import { useAuthStore } from '@/store/authStore'
import type { Order } from '@/types'
import dayjs from 'dayjs'

const OrderList = () => {
  const { user } = useAuthStore()
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin'
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [statusFilter, setStatusFilter] = useState<string>()

  useEffect(() => {
    loadOrders()
  }, [page, pageSize, statusFilter])

  const loadOrders = async () => {
    setLoading(true)
    try {
      const params: any = { page, page_size: pageSize }
      if (statusFilter) params.status = statusFilter
      const data: any = await orderApi.list(params)
      setOrders(data.list)
      setTotal(data.total)
    } catch (error: any) {
      message.error(error.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = (order: Order) => {
    Modal.confirm({
      title: '取消预约',
      content: '确认要取消这个预约吗？',
      onOk: async () => {
        try {
          await orderApi.cancel(order.id)
          message.success('取消成功')
          loadOrders()
        } catch (error: any) {
          message.error(error.message || '取消失败')
        }
      },
    })
  }

  const handleConfirm = (order: Order) => {
    Modal.confirm({
      title: '确认订单',
      content: '确认要通过这个预约吗？',
      onOk: async () => {
        try {
          await orderApi.confirm(order.id)
          message.success('确认成功')
          loadOrders()
        } catch (error: any) {
          message.error(error.message || '确认失败')
        }
      },
    })
  }

  const handleComplete = (order: Order) => {
    Modal.confirm({
      title: '标记完成',
      content: '确认要标记这个订单为已完成吗？',
      onOk: async () => {
        try {
          await orderApi.complete(order.id)
          message.success('操作成功')
          loadOrders()
        } catch (error: any) {
          message.error(error.message || '操作失败')
        }
      },
    })
  }

  const getStatusTag = (status: string) => {
    const map: Record<string, { text: string; color: string }> = {
      pending: { text: '待确认', color: 'orange' },
      confirmed: { text: '已确认', color: 'blue' },
      paid: { text: '已支付', color: 'green' },
      completed: { text: '已完成', color: 'default' },
      cancelled: { text: '已取消', color: 'red' },
    }
    const info = map[status] || { text: status, color: 'default' }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const canCancel = (order: Order) => {
    if (order.status === 'cancelled' || order.status === 'completed') return false
    const cancelDeadline = dayjs(order.start_time).subtract(24, 'hour')
    return dayjs().isBefore(cancelDeadline)
  }

  const columns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (type === 'venue' ? '场地' : '设备'),
    },
    {
      title: '项目',
      dataIndex: 'item_name',
      key: 'item_name',
    },
    {
      title: '时间',
      key: 'time',
      render: (_: any, record: Order) => (
        <div>
          <div>{dayjs(record.start_time).format('YYYY-MM-DD HH:mm')}</div>
          <div style={{ color: '#999' }}>
            至 {dayjs(record.end_time).format('YYYY-MM-DD HH:mm')}
          </div>
        </div>
      ),
    },
    {
      title: '时长',
      dataIndex: 'total_hours',
      key: 'total_hours',
      render: (val: number) => `${val} 小时`,
    },
    {
      title: '金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (val: number) => `¥${val.toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: getStatusTag,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => dayjs(val).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        <Space>
          {record.status === 'pending' && !isAdmin && canCancel(record) && (
            <Button type="link" danger icon={<CloseOutlined />} onClick={() => handleCancel(record)}>
              取消
            </Button>
          )}
          {isAdmin && record.status === 'pending' && (
            <Button type="link" icon={<CheckOutlined />} onClick={() => handleConfirm(record)}>
              确认
            </Button>
          )}
          {isAdmin && (record.status === 'confirmed' || record.status === 'paid') && (
            <Button type="link" icon={<CompleteOutlined />} onClick={() => handleComplete(record)}>
              完成
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <Card
      title="订单管理"
      extra={
        <Select
          placeholder="选择状态"
          style={{ width: 150 }}
          allowClear
          value={statusFilter}
          onChange={setStatusFilter}
        >
          <Select.Option value="pending">待确认</Select.Option>
          <Select.Option value="confirmed">已确认</Select.Option>
          <Select.Option value="paid">已支付</Select.Option>
          <Select.Option value="completed">已完成</Select.Option>
          <Select.Option value="cancelled">已取消</Select.Option>
        </Select>
      }
    >
      <Table
        columns={columns}
        dataSource={orders}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps)
          },
        }}
      />
    </Card>
  )
}

export default OrderList
