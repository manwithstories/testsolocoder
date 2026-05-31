import { useState } from 'react'
import { Card, Descriptions, Button, Tag, Timeline, Table, Modal, Form, Rate, Input, message, Divider, Space } from 'antd'
import { ArrowLeftOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { reservationApi, dailyRecordApi, reviewApi, orderApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { Reservation, DailyRecord, Review, Order } from '@/types'
import dayjs from 'dayjs'

export default function ReservationDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [reviewModal, setReviewModal] = useState(false)
  const [reviewForm] = Form.useForm()

  const { data: resData, isLoading } = useQuery({
    queryKey: ['reservation', id],
    queryFn: () => reservationApi.get(id!),
    enabled: !!id,
  })

  const { data: recordsData } = useQuery({
    queryKey: ['reservation', id, 'records'],
    queryFn: () => dailyRecordApi.listByReservation({ reservation_id: id, page_size: 100 }),
    enabled: !!id,
  })

  const { data: reviewData } = useQuery({
    queryKey: ['reservation', id, 'review'],
    queryFn: async () => {
      try {
        const res = await reservationApi.get(id!)
        return res.data
      } catch {
        return null
      }
    },
    enabled: !!id,
  })

  const { data: ordersData } = useQuery({
    queryKey: ['reservation', id, 'orders'],
    queryFn: () => orderApi.getByReservation(id!),
    enabled: !!id,
  })

  const reservation: Reservation | undefined = resData?.data
  const records: DailyRecord[] = recordsData?.data?.items || []
  const orders: Order[] = ordersData?.data || []

  const reviewMutation = useMutation({
    mutationFn: (values: any) => reviewApi.create({ ...values, reservation_id: id }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reservation', id] })
      message.success('评价成功')
      setReviewModal(false)
      reviewForm.resetFields()
    },
    onError: (err: any) => message.error(err.message || '评价失败'),
  })

  const checkInMutation = useMutation({
    mutationFn: () => reservationApi.checkIn(id!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reservation', id] })
      message.success('入住成功')
    },
    onError: (err: any) => message.error(err.message || '操作失败'),
  })

  const checkOutMutation = useMutation({
    mutationFn: () => reservationApi.checkOut(id!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reservation', id] })
      message.success('退房成功')
    },
    onError: (err: any) => message.error(err.message || '操作失败'),
  })

  const statusText: Record<string, { color: string; text: string }> = {
    pending: { color: 'orange', text: '待确认' },
    confirmed: { color: 'blue', text: '已确认' },
    checked_in: { color: 'green', text: '已入住' },
    completed: { color: 'purple', text: '已完成' },
    cancelled: { color: 'red', text: '已取消' },
  }

  const recordColumns = [
    {
      title: '日期',
      dataIndex: 'record_date',
      key: 'record_date',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD'),
    },
    { title: '饮食', dataIndex: 'feed_status', key: 'feed_status' },
    { title: '活动', dataIndex: 'activity', key: 'activity' },
    { title: '健康', dataIndex: 'health_status', key: 'health_status' },
    { title: '心情', dataIndex: 'mood', key: 'mood' },
  ]

  const orderColumns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no' },
    { title: '类型', dataIndex: 'type', key: 'type', render: (t: string) => t === 'prepay' ? '预付' : t === 'settlement' ? '结算' : '退款' },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (a: number) => `¥${a.toFixed(2)}` },
    {
      title: '状态',
      dataIndex: 'pay_status',
      key: 'pay_status',
      render: (s: string) => (
        <Tag color={s === 'paid' ? 'green' : s === 'unpaid' ? 'orange' : 'default'}>
          {s === 'paid' ? '已支付' : s === 'unpaid' ? '未支付' : '已退款'}
        </Tag>
      ),
    },
  ]

  if (isLoading || !reservation) {
    return <div className="text-center py-10">加载中...</div>
  }

  const s = statusText[reservation.status] || { color: 'default', text: reservation.status }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/reservations')}>
            返回
          </Button>
          <h2 className="text-xl font-semibold m-0">预约详情</h2>
          <Tag color={s.color}>{s.text}</Tag>
        </div>
        <Space>
          {user?.role === 'store' && reservation.status === 'confirmed' && (
            <Button type="primary" icon={<CheckOutlined />} onClick={() => checkInMutation.mutate()}>
              办理入住
            </Button>
          )}
          {user?.role === 'store' && reservation.status === 'checked_in' && (
            <Button type="primary" icon={<CloseOutlined />} onClick={() => checkOutMutation.mutate()}>
              办理退房
            </Button>
          )}
          {user?.role === 'owner' && reservation.status === 'completed' && (
            <Button type="primary" onClick={() => setReviewModal(true)}>
              去评价
            </Button>
          )}
        </Space>
      </div>

      <Card title="基本信息">
        <Descriptions column={2} size="small">
          <Descriptions.Item label="订单号">{reservation.order_no}</Descriptions.Item>
          <Descriptions.Item label="套餐类型">{reservation.package_type === 'daycare' ? '日托' : '寄养'}</Descriptions.Item>
          <Descriptions.Item label="入住日期">{dayjs(reservation.check_in_date).format('YYYY-MM-DD')}</Descriptions.Item>
          <Descriptions.Item label="退房日期">{dayjs(reservation.check_out_date).format('YYYY-MM-DD')}</Descriptions.Item>
          <Descriptions.Item label="寄养天数">{reservation.total_days} 天</Descriptions.Item>
          <Descriptions.Item label="总金额">¥{reservation.total_amount.toFixed(2)}</Descriptions.Item>
          <Descriptions.Item label="预付定金">¥{reservation.deposit_amount.toFixed(2)}</Descriptions.Item>
          <Descriptions.Item label="备注">{reservation.remark || '-'}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="每日动态">
        {records.length > 0 ? (
          <Table columns={recordColumns} dataSource={records} rowKey="id" size="small" pagination={false} />
        ) : (
          <div className="text-center text-gray-400 py-8">暂无动态记录</div>
        )}
      </Card>

      <Card title="订单信息">
        {orders.length > 0 ? (
          <Table columns={orderColumns} dataSource={orders} rowKey="id" size="small" pagination={false} />
        ) : (
          <div className="text-center text-gray-400 py-8">暂无订单</div>
        )}
      </Card>

      <Modal
        title="提交评价"
        open={reviewModal}
        onCancel={() => setReviewModal(false)}
        onOk={() => reviewForm.submit()}
      >
        <Form form={reviewForm} layout="vertical" onFinish={(v) => reviewMutation.mutate(v)}>
          <Form.Item name="store_rating" label="门店评分" rules={[{ required: true }]}>
            <Rate />
          </Form.Item>
          <Form.Item name="keeper_rating" label="管家评分">
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容">
            <Input.TextArea rows={4} placeholder="分享您的寄养体验..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
