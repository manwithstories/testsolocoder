import { useState } from 'react'
import { Card, Table, Tag, Button, Select, DatePicker, Space, Modal, message, Popconfirm } from 'antd'
import { PlusOutlined, EyeOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { reservationApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { Reservation } from '@/types'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker

export default function Reservations() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [status, setStatus] = useState<string>('')
  const [dateRange, setDateRange] = useState<any>(null)

  const params: any = { page_size: 100 }
  if (status) params.status = status
  if (dateRange && dateRange.length === 2) {
    params.start_date = dateRange[0].toISOString()
    params.end_date = dateRange[1].toISOString()
  }

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['reservations', 'list', status, dateRange],
    queryFn: () => reservationApi.list(params),
  })

  const reservations: Reservation[] = data?.data?.items || []
  const total = data?.data?.total || 0

  const confirmMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      reservationApi.confirm(id, { status, reason: '' }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reservations'] })
      message.success('操作成功')
    },
    onError: (err: any) => message.error(err.message || '操作失败'),
  })

  const cancelMutation = useMutation({
    mutationFn: (id: string) => reservationApi.cancel(id, { reason: '用户取消' }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reservations'] })
      message.success('取消成功')
    },
    onError: (err: any) => message.error(err.message || '操作失败'),
  })

  const handleCancel = (id: string) => {
    Modal.confirm({
      title: '确定取消预约？',
      content: '取消后将无法恢复',
      onOk: () => cancelMutation.mutate(id),
    })
  }

  const columns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no', width: 140 },
    {
      title: '日期范围',
      key: 'dates',
      render: (_: any, r: Reservation) => (
        <div>
          <div>{dayjs(r.check_in_date).format('YYYY-MM-DD')}</div>
          <div className="text-gray-400 text-xs">至 {dayjs(r.check_out_date).format('YYYY-MM-DD')}</div>
        </div>
      ),
    },
    { title: '天数', dataIndex: 'total_days', key: 'total_days', width: 60 },
    {
      title: '金额',
      key: 'amount',
      render: (_: any, r: Reservation) => `¥${r.total_amount.toFixed(2)}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (s: string) => <StatusTag status={s} />,
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, r: Reservation) => (
        <Space>
          <Button icon={<EyeOutlined />} size="small" onClick={() => navigate(`/reservations/${r.id}`)}>
            详情
          </Button>
          {user?.role === 'store' && r.status === 'pending' && (
            <>
              <Button size="small" type="primary" onClick={() => confirmMutation.mutate({ id: r.id, status: 'confirmed' })}>
                确认
              </Button>
              <Button size="small" danger onClick={() => confirmMutation.mutate({ id: r.id, status: 'cancelled' })}>
                拒绝
              </Button>
            </>
          )}
          {user?.role === 'owner' && (r.status === 'pending' || r.status === 'confirmed') && (
            <Popconfirm title="确定取消？" onConfirm={() => handleCancel(r.id)}>
              <Button size="small" danger>取消</Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div className="space-y-4">
      <Card
        title="寄养预约"
        extra={
          <Space>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/reservations/new')}>
              新建预约
            </Button>
          </Space>
        }
      >
        <div className="mb-4 flex gap-4 flex-wrap">
          <Select
            placeholder="选择状态"
            style={{ width: 150 }}
            allowClear
            value={status || undefined}
            onChange={setStatus}
            options={[
              { value: 'pending', label: '待确认' },
              { value: 'confirmed', label: '已确认' },
              { value: 'checked_in', label: '已入住' },
              { value: 'completed', label: '已完成' },
              { value: 'cancelled', label: '已取消' },
            ]}
          />
          <RangePicker value={dateRange} onChange={setDateRange} />
          <Button onClick={() => { setStatus(''); setDateRange(null); refetch() }}>
            重置
          </Button>
        </div>
        <Table
          columns={columns}
          dataSource={reservations}
          rowKey="id"
          loading={isLoading}
          pagination={{ pageSize: 10, total }}
        />
      </Card>
    </div>
  )
}

function StatusTag({ status }: { status: string }) {
  const map: Record<string, { color: string; text: string }> = {
    pending: { color: 'orange', text: '待确认' },
    confirmed: { color: 'blue', text: '已确认' },
    checked_in: { color: 'green', text: '已入住' },
    completed: { color: 'purple', text: '已完成' },
    cancelled: { color: 'red', text: '已取消' },
  }
  const s = map[status] || { color: 'default', text: status }
  return <Tag color={s.color}>{s.text}</Tag>
}
