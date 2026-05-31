import { useState } from 'react'
import { Card, Table, Tag, Button, Select, Space, Modal, message, Descriptions, Input, Rate, Form, Popconfirm } from 'antd'
import { EyeOutlined, CheckOutlined, CloseOutlined, PayCircleOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { useAuthContext } from '@/contexts/AuthContext'
import { appointmentApi } from '@/services/appointment'
import { recordApi } from '@/services/record'
import { Appointment, AppointmentStatus } from '@/types'

export function MyAppointments() {
  const navigate = useNavigate()
  const { user } = useAuthContext()
  const queryClient = useQueryClient()
  const [status, setStatus] = useState<string>('')
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [detailVisible, setDetailVisible] = useState(false)
  const [selectedAppointment, setSelectedAppointment] = useState<Appointment | null>(null)
  const [cancelVisible, setCancelVisible] = useState(false)
  const [cancelReason, setCancelReason] = useState('')
  const [reviewVisible, setReviewVisible] = useState(false)
  const [reviewForm] = Form.useForm()

  const { data, isLoading } = useQuery({
    queryKey: ['appointments', page, pageSize, status, user?.role],
    queryFn: () => {
      if (user?.role === 'professional') {
        return appointmentApi.getProfessionalAppointments({ page, page_size: pageSize, status: status || undefined })
      }
      return appointmentApi.getClientAppointments({ page, page_size: pageSize, status: status || undefined })
    },
    enabled: !!user,
  })

  const confirmMutation = useMutation({
    mutationFn: (id: string) => appointmentApi.confirm(id, { appointment_id: id }),
    onSuccess: () => {
      message.success('预约已确认')
      queryClient.invalidateQueries({ queryKey: ['appointments'] })
    },
    onError: (error: any) => {
      message.error(error.message || '确认失败')
    },
  })

  const cancelMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) =>
      appointmentApi.cancel(id, { appointment_id: id, reason }),
    onSuccess: () => {
      message.success('预约已取消')
      setCancelVisible(false)
      queryClient.invalidateQueries({ queryKey: ['appointments'] })
    },
    onError: (error: any) => {
      message.error(error.message || '取消失败')
    },
  })

  const completeMutation = useMutation({
    mutationFn: (id: string) => appointmentApi.complete(id),
    onSuccess: () => {
      message.success('预约已完成')
      queryClient.invalidateQueries({ queryKey: ['appointments'] })
    },
    onError: (error: any) => {
      message.error(error.message || '操作失败')
    },
  })

  const reviewMutation = useMutation({
    mutationFn: (data: { appointment_id: string; rating: number; content?: string }) =>
      recordApi.createReview(data),
    onSuccess: () => {
      message.success('评价已提交，等待审核')
      setReviewVisible(false)
      reviewForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['appointments'] })
    },
    onError: (error: any) => {
      message.error(error.message || '评价提交失败')
    },
  })

  const columns = [
    {
      title: '服务名称',
      dataIndex: ['service', 'title'],
      key: 'service',
    },
    {
      title: user?.role === 'professional' ? '客户' : '专业人士',
      dataIndex: user?.role === 'professional' ? ['client', 'full_name'] : ['professional', 'full_name'],
      key: 'person',
    },
    {
      title: '预约时间',
      key: 'time',
      render: (_: any, record: Appointment) =>
        `${record.schedule?.date} ${record.schedule?.start_time}-${record.schedule?.end_time}`,
    },
    {
      title: '金额',
      dataIndex: ['payment', 'amount'],
      key: 'amount',
      render: (amount: number) => `¥${amount?.toFixed(2) || '0.00'}`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: AppointmentStatus) => {
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
      },
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, record: Appointment) => (
        <Space>
          <Button
            size="small"
            icon={<EyeOutlined />}
            onClick={() => {
              setSelectedAppointment(record)
              setDetailVisible(true)
            }}
          >
            详情
          </Button>
          {user?.role === 'professional' && record.status === 'pending' && (
            <Button
              size="small"
              type="primary"
              icon={<CheckOutlined />}
              onClick={() => confirmMutation.mutate(record.id)}
            >
              确认
            </Button>
          )}
          {record.status === 'pending' && (
            <Button
              size="small"
              danger
              icon={<CloseOutlined />}
              onClick={() => {
                setSelectedAppointment(record)
                setCancelVisible(true)
              }}
            >
              取消
            </Button>
          )}
          {user?.role === 'professional' && record.status === 'confirmed' && (
            <Button
              size="small"
              type="primary"
              onClick={() => completeMutation.mutate(record.id)}
            >
              完成
            </Button>
          )}
          {user?.role === 'client' && record.status === 'completed' && !record.review && (
            <Button
              size="small"
              onClick={() => {
                setSelectedAppointment(record)
                setReviewVisible(true)
              }}
            >
              评价
            </Button>
          )}
          {record.payment?.status === 'pending' && (
            <Button
              size="small"
              type="primary"
              icon={<PayCircleOutlined />}
              onClick={() => navigate(`/appointments/${record.id}`)}
            >
              支付
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div className="page-container">
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 24 }}>
        <h2>我的预约</h2>
        <Select
          style={{ width: 160 }}
          value={status}
          onChange={setStatus}
          allowClear
          placeholder="筛选状态"
          options={[
            { value: '', label: '全部' },
            { value: 'pending', label: '待确认' },
            { value: 'confirmed', label: '已确认' },
            { value: 'completed', label: '已完成' },
            { value: 'cancelled', label: '已取消' },
            { value: 'refunded', label: '已退款' },
          ]}
        />
      </div>

      <Card>
        <Table
          columns={columns}
          dataSource={data?.items || []}
          loading={isLoading}
          rowKey="id"
          pagination={{
            current: page,
            pageSize,
            total: data?.total || 0,
            onChange: setPage,
          }}
        />
      </Card>

      <Modal
        title="预约详情"
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={null}
        width={600}
      >
        {selectedAppointment && (
          <Descriptions column={1} bordered>
            <Descriptions.Item label="服务名称">
              {selectedAppointment.service?.title}
            </Descriptions.Item>
            <Descriptions.Item label={user?.role === 'professional' ? '客户' : '专业人士'}>
              {user?.role === 'professional'
                ? selectedAppointment.client?.full_name
                : selectedAppointment.professional?.full_name}
            </Descriptions.Item>
            <Descriptions.Item label="预约时间">
              {selectedAppointment.schedule?.date} {selectedAppointment.schedule?.start_time}-
              {selectedAppointment.schedule?.end_time}
            </Descriptions.Item>
            <Descriptions.Item label="金额">
              ¥{selectedAppointment.payment?.amount?.toFixed(2) || '0.00'}
            </Descriptions.Item>
            <Descriptions.Item label="支付状态">
              <Tag color={
                selectedAppointment.payment?.status === 'paid' ? 'green' :
                selectedAppointment.payment?.status === 'pending' ? 'orange' :
                selectedAppointment.payment?.status === 'refunded' ? 'purple' : 'red'
              }>
                {selectedAppointment.payment?.status === 'paid' ? '已支付' :
                 selectedAppointment.payment?.status === 'pending' ? '待支付' :
                 selectedAppointment.payment?.status === 'refunded' ? '已退款' : '已取消'}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="备注">
              {selectedAppointment.notes || '无'}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>

      <Modal
        title="取消预约"
        open={cancelVisible}
        onOk={() => {
          if (selectedAppointment && cancelReason) {
            cancelMutation.mutate({ id: selectedAppointment.id, reason: cancelReason })
          }
        }}
        onCancel={() => {
          setCancelVisible(false)
          setCancelReason('')
        }}
        confirmLoading={cancelMutation.isPending}
      >
        <Input.TextArea
          rows={4}
          value={cancelReason}
          onChange={(e) => setCancelReason(e.target.value)}
          placeholder="请输入取消原因"
        />
      </Modal>

      <Modal
        title="评价服务"
        open={reviewVisible}
        onCancel={() => {
          setReviewVisible(false)
          reviewForm.resetFields()
        }}
        footer={null}
      >
        <Form form={reviewForm} onFinish={(values) => {
          if (selectedAppointment) {
            reviewMutation.mutate({
              appointment_id: selectedAppointment.id,
              rating: values.rating,
              content: values.content,
            })
          }
        }}>
          <Form.Item
            name="rating"
            label="评分"
            rules={[{ required: true, message: '请评分' }]}
          >
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容">
            <Input.TextArea rows={4} placeholder="请输入您的评价" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={reviewMutation.isPending} block>
              提交评价
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
