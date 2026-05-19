import { useState, useEffect } from 'react'
import { Table, Button, Tag, Space, Modal, Form, Input, Rate, Select, message, Card } from 'antd'
import { CheckOutlined, CloseOutlined, StarOutlined } from '@ant-design/icons'
import { reviewApi } from '@/api'
import { useAuthStore } from '@/store/authStore'
import type { Review } from '@/types'
import dayjs from 'dayjs'

const ReviewList = () => {
  const { user } = useAuthStore()
  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin'
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [statusFilter, setStatusFilter] = useState<string>()
  const [reviewModalVisible, setReviewModalVisible] = useState(false)
  const [orderId, setOrderId] = useState<number | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    loadReviews()
  }, [page, pageSize, statusFilter])

  const loadReviews = async () => {
    setLoading(true)
    try {
      const params: any = { page, page_size: pageSize }
      if (statusFilter) params.status = statusFilter
      const data: any = await reviewApi.list(params)
      setReviews(data.list)
      setTotal(data.total)
    } catch (error: any) {
      message.error(error.message || '加载失败')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmitReview = async (values: any) => {
    if (!orderId) return
    try {
      await reviewApi.create({
        order_id: orderId,
        ...values,
      })
      message.success('评价提交成功，等待审核')
      setReviewModalVisible(false)
      loadReviews()
    } catch (error: any) {
      message.error(error.message || '提交失败')
    }
  }

  const handleApprove = (id: number) => {
    Modal.confirm({
      title: '审核通过',
      content: '确认要通过这个评价吗？',
      onOk: async () => {
        try {
          await reviewApi.approve(id)
          message.success('审核通过')
          loadReviews()
        } catch (error: any) {
          message.error(error.message || '操作失败')
        }
      },
    })
  }

  const handleReject = (id: number) => {
    Modal.confirm({
      title: '审核拒绝',
      content: '确认要拒绝这个评价吗？',
      onOk: async () => {
        try {
          await reviewApi.reject(id)
          message.success('已拒绝')
          loadReviews()
        } catch (error: any) {
          message.error(error.message || '操作失败')
        }
      },
    })
  }

  const getStatusTag = (status: string) => {
    const map: Record<string, { text: string; color: string }> = {
      pending: { text: '待审核', color: 'orange' },
      approved: { text: '已通过', color: 'green' },
      rejected: { text: '已拒绝', color: 'red' },
    }
    const info = map[status] || { text: status, color: 'default' }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 60,
    },
    {
      title: '用户',
      dataIndex: ['user', 'username'],
      key: 'user',
    },
    {
      title: '订单号',
      dataIndex: ['order', 'order_no'],
      key: 'order_no',
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      render: (rating: number) => (
        <span>
          {[...Array(rating)].map((_, i) => (
            <StarOutlined key={i} style={{ color: '#fadb14' }} />
          ))}
        </span>
      ),
    },
    {
      title: '内容',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
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
    isAdmin && {
      title: '操作',
      key: 'action',
      render: (_: any, record: Review) => (
        <Space>
          {record.status === 'pending' && (
            <>
              <Button type="link" icon={<CheckOutlined />} onClick={() => handleApprove(record.id)}>
                通过
              </Button>
              <Button type="link" danger icon={<CloseOutlined />} onClick={() => handleReject(record.id)}>
                拒绝
              </Button>
            </>
          )}
        </Space>
      ),
    },
  ].filter(Boolean)

  return (
    <Card
      title="评价管理"
      extra={
        <Space>
          <Select
            placeholder="选择状态"
            style={{ width: 150 }}
            allowClear
            value={statusFilter}
            onChange={setStatusFilter}
          >
            <Select.Option value="pending">待审核</Select.Option>
            <Select.Option value="approved">已通过</Select.Option>
            <Select.Option value="rejected">已拒绝</Select.Option>
          </Select>
        </Space>
      }
    >
      <Table
        columns={columns}
        dataSource={reviews}
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

      <Modal
        title="提交评价"
        open={reviewModalVisible}
        onCancel={() => setReviewModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmitReview}>
          <Form.Item
            name="rating"
            label="评分"
            rules={[{ required: true, message: '请选择评分' }]}
          >
            <Rate />
          </Form.Item>
          <Form.Item
            name="content"
            label="评价内容"
            rules={[{ max: 1000, message: '最多1000个字符' }]}
          >
            <Input.TextArea rows={4} placeholder="请输入评价内容" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              提交评价
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default ReviewList
