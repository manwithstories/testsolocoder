import React, { useState, useEffect } from 'react'
import { Table, Tag, Button, Select, DatePicker, Card, Modal, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { orderApi } from '@/services/order'
import { Order } from '@/types'
import { formatPrice, formatDate, getOrderStatusText, getOrderStatusColor } from '@/utils'
import { usePagination } from '@/hooks'
import { useAppSelector } from '@/store/hooks'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

const OrderList: React.FC = () => {
  const navigate = useNavigate()
  const { userInfo } = useAppSelector((state) => state.auth)
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [orders, setOrders] = useState<Order[]>([])
  const [loading, setLoading] = useState(false)
  const [status, setStatus] = useState<string | undefined>()
  const [cancelModalVisible, setCancelModalVisible] = useState(false)
  const [cancelOrder, setCancelOrder] = useState<Order | null>(null)
  const [cancelReason, setCancelReason] = useState('')

  useEffect(() => {
    loadOrders()
  }, [page, pageSize, status])

  const loadOrders = async () => {
    setLoading(true)
    try {
      const res = await orderApi.getList({
        page,
        page_size: pageSize,
        status,
      })
      setOrders(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = async () => {
    if (!cancelOrder || !cancelReason) {
      message.warning('请填写取消原因')
      return
    }
    try {
      await orderApi.cancel(cancelOrder.id, { reason: cancelReason })
      message.success('订单已取消')
      setCancelModalVisible(false)
      setCancelOrder(null)
      setCancelReason('')
      loadOrders()
    } catch (error) {
      console.error(error)
    }
  }

  const handleStart = async (order: Order) => {
    try {
      await orderApi.start(order.id)
      message.success('服务已开始')
      loadOrders()
    } catch (error) {
      console.error(error)
    }
  }

  const handleComplete = async (order: Order) => {
    try {
      await orderApi.complete(order.id)
      message.success('服务已完成')
      loadOrders()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      key: 'order_no',
    },
    {
      title: '服务名称',
      dataIndex: 'service_name',
      key: 'service_name',
      render: (_: any, record: Order) => record.service_item?.name || '-',
    },
    {
      title: userInfo?.role === 'customer' ? '服务人员' : '客户',
      dataIndex: 'user_name',
      key: 'user',
      render: (_: any, record: Order) =>
        userInfo?.role === 'customer'
          ? record.provider?.nickname || '-'
          : record.customer?.nickname || '-',
    },
    {
      title: '预约时间',
      dataIndex: 'appointment_time',
      key: 'appointment_time',
      render: (text: string) => formatDate(text),
    },
    {
      title: '金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (text: number) => formatPrice(text),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (text: string) => (
        <Tag color={getOrderStatusColor(text)}>{getOrderStatusText(text)}</Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        <div style={{ display: 'flex', gap: 8 }}>
          <Button type="link" onClick={() => navigate(`/orders/${record.id}`)}>
            详情
          </Button>
          {userInfo?.role === 'service_provider' && record.status === 'confirmed' && (
            <Button type="link" onClick={() => handleStart(record)}>
              开始服务
            </Button>
          )}
          {userInfo?.role === 'service_provider' && record.status === 'in_service' && (
            <Button type="link" onClick={() => handleComplete(record)}>
              完成服务
            </Button>
          )}
          {(record.status === 'pending' || record.status === 'confirmed') && (
            <Button
              type="link"
              danger
              onClick={() => {
                setCancelOrder(record)
                setCancelModalVisible(true)
              }}
            >
              取消
            </Button>
          )}
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">订单管理</h1>
      </div>

      <div className="filter-bar">
        <Select
          placeholder="订单状态"
          style={{ width: 150 }}
          allowClear
          value={status}
          onChange={(value) => {
            setStatus(value)
            setPage(1)
          }}
          options={[
            { label: '待接单', value: 'pending' },
            { label: '已确认', value: 'confirmed' },
            { label: '服务中', value: 'in_service' },
            { label: '已完成', value: 'completed' },
            { label: '已取消', value: 'cancelled' },
          ]}
        />
      </div>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={orders}
          columns={columns}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            onChange: (p, ps) => {
              setPage(p)
              setPageSize(ps)
            },
          }}
        />
      </Card>

      <Modal
        title="取消订单"
        open={cancelModalVisible}
        onOk={handleCancel}
        onCancel={() => {
          setCancelModalVisible(false)
          setCancelOrder(null)
          setCancelReason('')
        }}
      >
        <p>确定要取消该订单吗？</p>
        <textarea
          style={{ width: '100%', minHeight: 80, padding: 8, border: '1px solid #d9d9d9', borderRadius: 4 }}
          placeholder="请输入取消原因"
          value={cancelReason}
          onChange={(e) => setCancelReason(e.target.value)}
        />
      </Modal>
    </div>
  )
}

export default OrderList
