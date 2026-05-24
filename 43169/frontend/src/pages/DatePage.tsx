import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Card, Table, Button, Modal, Form, Input, DatePicker, InputNumber, Tag, Space, message, Rate, Tabs } from 'antd'
import { PlusOutlined, CalendarOutlined, CheckOutlined, CloseOutlined, DeleteOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { dateApi, DateRecord, DateReview } from '@/api/endpoints'

export default function DatePage() {
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [inviteModal, setInviteModal] = useState(false)
  const [reviewModal, setReviewModal] = useState(false)
  const [currentDate, setCurrentDate] = useState<DateRecord | null>(null)
  const [inviteForm] = Form.useForm()
  const [reviewForm] = Form.useForm()
  const queryClient = useQueryClient()

  const { data: datesData, isLoading } = useQuery({
    queryKey: ['dates', page],
    queryFn: () => dateApi.list({ page, page_size: pageSize }),
  })

  const { data: reviewsData } = useQuery({
    queryKey: ['reviews'],
    queryFn: () => dateApi.getReviews({ page: 1, page_size: 100 }),
  })

  const inviteMutation = useMutation({
    mutationFn: dateApi.createInvite,
    onSuccess: () => {
      message.success('邀请已发送')
      setInviteModal(false)
      inviteForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['dates'] })
    },
  })

  const handleAction = (id: number, action: 'accept' | 'reject' | 'cancel' | 'complete') => {
    const apiMap: Record<string, any> = {
      accept: dateApi.accept,
      reject: dateApi.reject,
      cancel: dateApi.cancel,
      complete: dateApi.complete,
    }
    apiMap[action](id).then(() => {
      message.success('操作成功')
      queryClient.invalidateQueries({ queryKey: ['dates'] })
    })
  }

  const reviewMutation = useMutation({
    mutationFn: dateApi.createReview,
    onSuccess: () => {
      message.success('评价成功')
      setReviewModal(false)
      reviewForm.resetFields()
      queryClient.invalidateQueries({ queryKey: ['reviews'] })
    },
  })

  const statusColors: Record<string, string> = {
    pending: 'orange',
    accepted: 'blue',
    rejected: 'red',
    canceled: 'default',
    completed: 'green',
  }

  const statusText: Record<string, string> = {
    pending: '待确认',
    accepted: '已接受',
    rejected: '已拒绝',
    canceled: '已取消',
    completed: '已完成',
  }

  const columns = [
    { title: '标题', dataIndex: 'title', key: 'title' },
    { title: '地点', dataIndex: 'location', key: 'location' },
    {
      title: '时间',
      dataIndex: 'date_at',
      key: 'date_at',
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
    },
    { title: '时长(分钟)', dataIndex: 'duration', key: 'duration' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={statusColors[status]}>{statusText[status]}</Tag>,
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: DateRecord) => (
        <Space>
          {record.status === 'pending' && (
            <>
              <Button size="small" type="primary" icon={<CheckOutlined />} onClick={() => handleAction(record.id, 'accept')}>
                接受
              </Button>
              <Button size="small" danger icon={<CloseOutlined />} onClick={() => handleAction(record.id, 'reject')}>
                拒绝
              </Button>
            </>
          )}
          {record.status === 'accepted' && (
            <>
              <Button size="small" type="primary" onClick={() => handleAction(record.id, 'complete')}>
                完成
              </Button>
              <Button size="small" danger icon={<DeleteOutlined />} onClick={() => handleAction(record.id, 'cancel')}>
                取消
              </Button>
            </>
          )}
          {record.status === 'completed' && (
            <Button size="small" onClick={() => { setCurrentDate(record); setReviewModal(true) }}>
              评价
            </Button>
          )}
        </Space>
      ),
    },
  ]

  const reviewColumns = [
    {
      title: '评价对象',
      dataIndex: 'target_id',
      key: 'target_id',
      render: (id: number) => `用户#${id}`,
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      render: (rating: number) => <Rate disabled value={rating} />,
    },
    { title: '内容', dataIndex: 'content', key: 'content' },
    {
      title: '时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (v: string) => dayjs(v).format('YYYY-MM-DD HH:mm'),
    },
  ]

  return (
    <div>
      <Card
        title="约会管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setInviteModal(true)}>
            发起约会
          </Button>
        }
      >
        <Tabs
          items={[
            {
              key: 'dates',
              label: '约会列表',
              children: (
                <Table
                  columns={columns}
                  dataSource={datesData?.list}
                  loading={isLoading}
                  rowKey="id"
                  pagination={{
                    current: page,
                    pageSize,
                    total: datesData?.total || 0,
                    onChange: setPage,
                  }}
                />
              ),
            },
            {
              key: 'reviews',
              label: '我的评价',
              children: (
                <Table
                  columns={reviewColumns}
                  dataSource={reviewsData?.list}
                  rowKey="id"
                  pagination={false}
                />
              ),
            },
          ]}
        />
      </Card>

      <Modal
        title="发起约会邀请"
        open={inviteModal}
        onCancel={() => setInviteModal(false)}
        footer={null}
      >
        <Form form={inviteForm} onFinish={(v) => inviteMutation.mutate({ ...v, date_at: dayjs(v.date_at).format('YYYY-MM-DD HH:mm:ss') })} layout="vertical">
          <Form.Item name="receiver_id" label="对方ID" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="title" label="约会标题" rules={[{ required: true }]}>
            <Input placeholder="例如：一起喝咖啡" />
          </Form.Item>
          <Form.Item name="location" label="地点">
            <Input placeholder="约会地点" />
          </Form.Item>
          <Form.Item name="date_at" label="约会时间" rules={[{ required: true }]}>
            <DatePicker showTime style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="duration" label="时长(分钟)">
            <InputNumber min={30} max={480} defaultValue={60} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="note" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={inviteMutation.isPending} block>
              发送邀请
            </Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="评价约会"
        open={reviewModal}
        onCancel={() => setReviewModal(false)}
        footer={null}
      >
        <Form form={reviewForm} onFinish={(v) => reviewMutation.mutate({ ...v, date_id: currentDate?.id, target_id: currentDate?.initiator_id })} layout="vertical">
          <Form.Item name="rating" label="评分" rules={[{ required: true }]}>
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容">
            <Input.TextArea rows={4} placeholder="分享你的约会体验..." />
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
