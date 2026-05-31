import { Card, Table, Tag, Button, Modal, Form, Select, InputNumber, message, Space } from 'antd'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { orderApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { Order } from '@/types'
import dayjs from 'dayjs'
import { useState } from 'react'

export default function Orders() {
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [refundModal, setRefundModal] = useState(false)
  const [refundForm] = Form.useForm()
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['orders', 'list'],
    queryFn: () => orderApi.list({ page_size: 100 }),
  })

  const orders: Order[] = data?.data?.items || []
  const total = data?.data?.total || 0

  const refundMutation = useMutation({
    mutationFn: ({ id, amount, reason }: { id: string; amount: number; reason: string }) =>
      orderApi.refund(id, { amount, reason }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      message.success('退款成功')
      setRefundModal(false)
      refundForm.resetFields()
    },
    onError: (err: any) => message.error(err.message || '退款失败'),
  })

  const payMutation = useMutation({
    mutationFn: (values: any) => orderApi.pay(values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      message.success('支付成功')
    },
    onError: (err: any) => message.error(err.message || '支付失败'),
  })

  const handleRefund = (order: Order) => {
    setCurrentOrder(order)
    refundForm.resetFields()
    setRefundModal(true)
  }

  const columns = [
    { title: '订单号', dataIndex: 'order_no', key: 'order_no', width: 140 },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (t: string) => {
        const map: Record<string, { color: string; text: string }> = {
          prepay: { color: 'blue', text: '预付款' },
          settlement: { color: 'green', text: '结算款' },
          refund: { color: 'red', text: '退款' },
        }
        const s = map[t] || { color: 'default', text: t }
        return <Tag color={s.color}>{s.text}</Tag>
      },
    },
    { title: '金额', dataIndex: 'amount', key: 'amount', render: (a: number) => `¥${a.toFixed(2)}` },
    {
      title: '支付状态',
      dataIndex: 'pay_status',
      key: 'pay_status',
      render: (s: string) => (
        <Tag color={s === 'paid' ? 'green' : s === 'unpaid' ? 'orange' : 'red'}>
          {s === 'paid' ? '已支付' : s === 'unpaid' ? '未支付' : '已退款'}
        </Tag>
      ),
    },
    { title: '支付方式', dataIndex: 'pay_method', key: 'pay_method', render: (m: string) => m || '-' },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, order: Order) => (
        <Space>
          {user?.role === 'owner' && order.pay_status === 'unpaid' && (
            <Button
              size="small"
              type="primary"
              onClick={() => {
                const amountHash = require('crypto').createHash('md5').update(`${order.amount.toFixed(2)}-${order.reservation_id}`).digest('hex')
                payMutation.mutate({
                  reservation_id: order.reservation_id,
                  amount: order.amount,
                  amount_hash: amountHash,
                  pay_method: 'wechat',
                })
              }}
            >
              支付
            </Button>
          )}
          {user?.role === 'store' && order.pay_status === 'paid' && (
            <Button size="small" danger onClick={() => handleRefund(order)}>
              退款
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div className="space-y-4">
      <Card title="订单管理">
        <Table columns={columns} dataSource={orders} rowKey="id" loading={isLoading} pagination={{ pageSize: 10, total }} />
      </Card>

      <Modal
        title="订单退款"
        open={refundModal}
        onCancel={() => setRefundModal(false)}
        onOk={() => refundForm.submit()}
      >
        <Form form={refundForm} layout="vertical" onFinish={(v) => refundMutation.mutate({ id: currentOrder!.id, ...v })}>
          <Form.Item name="amount" label="退款金额" rules={[{ required: true }]}>
            <InputNumber className="w-full" min={0} max={currentOrder?.amount} step={0.01} />
          </Form.Item>
          <Form.Item name="reason" label="退款原因">
            <Input.TextArea rows={3} placeholder="请输入退款原因" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
