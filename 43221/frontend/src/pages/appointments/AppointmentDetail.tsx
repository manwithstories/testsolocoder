import { useParams } from 'react-router-dom'
import { Card, Descriptions, Button, Tag, Statistic, CountDown, message, Modal } from 'antd'
import { DollarOutlined, ClockCircleOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useAuthContext } from '@/contexts/AuthContext'
import { appointmentApi } from '@/services/appointment'
import { Appointment } from '@/types'

export function AppointmentDetail() {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuthContext()
  const queryClient = useQueryClient()

  const { data: appointment, isLoading } = useQuery({
    queryKey: ['appointment', id],
    queryFn: () => appointmentApi.getById(id!),
    enabled: !!id,
  })

  const payMutation = useMutation({
    mutationFn: () => appointmentApi.pay(id!, { transaction_id: `mock_${Date.now()}` }),
    onSuccess: () => {
      message.success('支付成功')
      queryClient.invalidateQueries({ queryKey: ['appointment', id] })
    },
    onError: (error: any) => {
      message.error(error.message || '支付失败')
    },
  })

  const refundMutation = useMutation({
    mutationFn: (reason: string) => appointmentApi.refund(id!, { reason }),
    onSuccess: () => {
      message.success('退款申请已提交')
      queryClient.invalidateQueries({ queryKey: ['appointment', id] })
    },
    onError: (error: any) => {
      message.error(error.message || '退款失败')
    },
  })

  if (isLoading) {
    return <div className="loading-spin">加载中...</div>
  }

  if (!appointment) {
    return <div>预约不存在</div>
  }

  const handlePay = () => {
    Modal.confirm({
      title: '确认支付',
      content: `您将支付 ¥${appointment.payment?.amount?.toFixed(2)}`,
      onOk: () => payMutation.mutate(),
    })
  }

  const handleRefund = () => {
    Modal.confirm({
      title: '申请退款',
      content: '请输入退款原因',
      input: { type: 'textarea' as const },
      onOk: (value) => refundMutation.mutate(value || '用户申请退款'),
    })
  }

  const getStatusTag = (status: string) => {
    const colorMap: Record<string, string> = {
      pending: 'orange',
      confirmed: 'blue',
      completed: 'green',
      cancelled: 'red',
      refunded: 'purple',
    }
    const textMap: Record<string, string> = {
      pending: '待确认',
      confirmed: '已确认',
      completed: '已完成',
      cancelled: '已取消',
      refunded: '已退款',
    }
    return <Tag color={colorMap[status]}>{textMap[status]}</Tag>
  }

  return (
    <div className="page-container">
      <Card
        title="预约详情"
        extra={getStatusTag(appointment.status)}
      >
        <Descriptions column={2} bordered>
          <Descriptions.Item label="服务名称" span={2}>
            {appointment.service?.title}
          </Descriptions.Item>
          <Descriptions.Item label="客户">
            {appointment.client?.full_name}
          </Descriptions.Item>
          <Descriptions.Item label="专业人士">
            {appointment.professional?.full_name}
          </Descriptions.Item>
          <Descriptions.Item label="预约时间">
            {appointment.schedule?.date} {appointment.schedule?.start_time}-{appointment.schedule?.end_time}
          </Descriptions.Item>
          <Descriptions.Item label="服务时长">
            {appointment.service?.duration_minutes}分钟
          </Descriptions.Item>
          <Descriptions.Item label="备注" span={2}>
            {appointment.notes || '无'}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="支付信息" style={{ marginTop: 24 }}>
        <Descriptions column={2} bordered>
          <Descriptions.Item label="支付金额">
            <span style={{ color: '#ff4d4f', fontSize: 20, fontWeight: 600 }}>
              ¥{appointment.payment?.amount?.toFixed(2)}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="支付状态">
            <Tag color={
              appointment.payment?.status === 'paid' ? 'green' :
              appointment.payment?.status === 'pending' ? 'orange' :
              appointment.payment?.status === 'refunded' ? 'purple' : 'red'
            }>
              {appointment.payment?.status === 'paid' ? '已支付' :
               appointment.payment?.status === 'pending' ? '待支付' :
               appointment.payment?.status === 'refunded' ? '已退款' : '已取消'}
            </Tag>
          </Descriptions.Item>
          {appointment.payment?.expires_at && appointment.payment.status === 'pending' && (
            <Descriptions.Item label="支付剩余时间" span={2}>
              <CountDown
                value={new Date(appointment.payment.expires_at).getTime()}
                format="HH:mm:ss"
                onFinish={() => message.warning('支付已超时')}
              />
            </Descriptions.Item>
          )}
        </Descriptions>

        <div style={{ marginTop: 24, textAlign: 'right' }}>
          {appointment.payment?.status === 'pending' && user?.role === 'client' && (
            <Button
              type="primary"
              size="large"
              icon={<DollarOutlined />}
              onClick={handlePay}
              loading={payMutation.isPending}
            >
              立即支付
            </Button>
          )}
          {appointment.payment?.status === 'paid' && user?.role === 'client' && appointment.status === 'confirmed' && (
            <Button
              danger
              size="large"
              onClick={handleRefund}
              loading={refundMutation.isPending}
            >
              申请退款
            </Button>
          )}
        </div>
      </Card>

      {appointment.consult_record && (
        <Card title="咨询记录" style={{ marginTop: 24 }}>
          <Descriptions column={1} bordered>
            <Descriptions.Item label="咨询总结">
              {appointment.consult_record.summary || '无'}
            </Descriptions.Item>
            <Descriptions.Item label="后续建议">
              {appointment.consult_record.advice || '无'}
            </Descriptions.Item>
            {appointment.consult_record.follow_up_date && (
              <Descriptions.Item label="下次回访日期">
                {appointment.consult_record.follow_up_date}
              </Descriptions.Item>
            )}
          </Descriptions>
        </Card>
      )}
    </div>
  )
}
