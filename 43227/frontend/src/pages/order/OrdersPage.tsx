import { useEffect, useState } from 'react'
import { Table, Button, Tag, Space, message, Card } from 'antd'
import { EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import api from '../../api'
import type { Order } from '../../types'
import { useAuthStore } from '../../store/authStore'
import dayjs from 'dayjs'

function OrdersPage() {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [data, setData] = useState<Order[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get('/orders', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as Order[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, pageSize])

  const statusColorMap: Record<string, string> = {
    pending: 'default',
    paid: 'blue',
    shipped: 'cyan',
    delivered: 'geekblue',
    completed: 'green',
    cancelled: 'red',
    refunded: 'orange',
  }

  const paymentStatusColorMap: Record<string, string> = {
    unpaid: 'red',
    paid: 'green',
    refunded: 'orange',
  }

  const columns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    {
      title: user?.role === 'buyer' ? '卖家' : '买家',
      key: 'user',
      render: (_: any, record: Order) => (
        user?.role === 'buyer' ? record.seller?.username : record.buyer?.username
      ),
    },
    {
      title: '商品',
      key: 'product',
      render: (_: any, record: Order) => record.product?.title,
    },
    { title: '数量', dataIndex: 'quantity', key: 'quantity', render: (q: number, r: Order) => `${q} ${r.product?.unit || 'kg'}` },
    {
      title: '总金额',
      dataIndex: 'total_amount',
      key: 'total_amount',
      render: (amount: number) => <span style={{ color: '#f5222d', fontWeight: 'bold' }}>¥{amount}</span>,
    },
    {
      title: '订单状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={statusColorMap[status]}>{status}</Tag>,
    },
    {
      title: '支付状态',
      dataIndex: 'payment_status',
      key: 'payment_status',
      render: (status: string) => <Tag color={paymentStatusColorMap[status]}>{status}</Tag>,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Order) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />} onClick={() => navigate(`/orders/${record.id}`)}>
            查看
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <Card title={user?.role === 'buyer' ? '我的订单' : '订单管理'}>
      <Table
        columns={columns}
        dataSource={data}
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

export default OrdersPage
