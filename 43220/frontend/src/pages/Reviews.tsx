import { useState } from 'react'
import { Card, Table, Rate, Tag, Button, Modal, Form, Input, message, Avatar } from 'antd'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { reviewApi } from '@/services/api'
import { useAuthStore } from '@/context/AuthContext'
import type { Review } from '@/types'
import dayjs from 'dayjs'

export default function Reviews() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [replyModal, setReplyModal] = useState(false)
  const [replyForm] = Form.useForm()
  const [currentReview, setCurrentReview] = useState<Review | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['reviews', 'list'],
    queryFn: () => reviewApi.listByStore({ page_size: 100 }),
  })

  const reviews: Review[] = data?.data?.items || []
  const total = data?.data?.total || 0

  const replyMutation = useMutation({
    mutationFn: ({ id, reply }: { id: string; reply: string }) => reviewApi.reply(id, { reply }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] })
      message.success('回复成功')
      setReplyModal(false)
      replyForm.resetFields()
    },
    onError: (err: any) => message.error(err.message || '回复失败'),
  })

  const handleReply = (review: Review) => {
    setCurrentReview(review)
    replyForm.resetFields()
    setReplyModal(true)
  }

  const columns = [
    {
      title: '门店评分',
      dataIndex: 'store_rating',
      key: 'store_rating',
      render: (rating: number) => <Rate disabled allowHalf value={rating} />,
    },
    {
      title: '管家评分',
      dataIndex: 'keeper_rating',
      key: 'keeper_rating',
      render: (rating: number) => rating ? <Rate disabled allowHalf value={rating} /> : '-',
    },
    {
      title: '评价内容',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
    },
    {
      title: '商家回复',
      dataIndex: 'reply',
      key: 'reply',
      render: (reply: string) => reply || '-',
    },
    {
      title: '评价时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (d: string) => dayjs(d).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '操作',
      key: 'actions',
      render: (_: any, review: Review) =>
      user?.role === 'store' && !review.reply ? (
        <Button size="small" onClick={() => handleReply(review)}>
          回复
        </Button>
      ) : null,
    },
  ]

  return (
    <div className="space-y-4">
      <Card title="评价管理">
        <Table columns={columns} dataSource={reviews} rowKey="id" loading={isLoading} pagination={{ pageSize: 10, total }} />
      </Card>

      <Modal
        title="回复评价"
        open={replyModal}
        onCancel={() => setReplyModal(false)}
        onOk={() => replyForm.submit()}
      >
        <Form form={replyForm} layout="vertical" onFinish={(v) => replyMutation.mutate({ id: currentReview!.id, reply: v.reply })}>
          <Form.Item name="reply" label="回复内容" rules={[{ required: true }]}>
            <Input.TextArea rows={4} placeholder="请输入回复内容..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
