import { Card, List, Avatar, Rate, Tag, Button, Space, Modal, Input, message } from 'antd'
import { CheckOutlined, CloseOutlined, UserOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { recordApi } from '@/services/record'
import { Review } from '@/types'

export function ReviewManagement() {
  const queryClient = useQueryClient()

  const { data, isLoading } = useQuery({
    queryKey: ['pending-reviews'],
    queryFn: () => recordApi.getPendingReviews({ page: 1, page_size: 50 }),
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, status, rejectReason }: { id: string; status: string; rejectReason?: string }) =>
      recordApi.updateReviewStatus({ review_id: id, status, reject_reason: rejectReason }),
    onSuccess: () => {
      message.success('操作成功')
      queryClient.invalidateQueries({ queryKey: ['pending-reviews'] })
    },
    onError: (error: any) => {
      message.error(error.message || '操作失败')
    },
  })

  const handleReject = (review: Review) => {
    Modal.confirm({
      title: '拒绝评价',
      content: (
        <Input.TextArea
          rows={4}
          placeholder="请输入拒绝原因"
          onChange={(e) => {
            Modal.confirm({
              title: '确认拒绝',
              content: `拒绝原因：${e.target.value}`,
              onOk: () => updateMutation.mutate({ id: review.id, status: 'rejected', rejectReason: e.target.value }),
            })
          }}
        />
      ),
    })
  }

  return (
    <div className="page-container">
      <h2 style={{ marginBottom: 24 }}>评价审核</h2>

      <Card>
        <List
          dataSource={data?.items || []}
          loading={isLoading}
          renderItem={(review: Review) => (
            <List.Item
              actions={[
                <Space key="actions">
                  <Button
                    type="primary"
                    icon={<CheckOutlined />}
                    onClick={() => updateMutation.mutate({ id: review.id, status: 'approved' })}
                    loading={updateMutation.isPending}
                  >
                    通过
                  </Button>
                  <Button
                    danger
                    icon={<CloseOutlined />}
                    onClick={() => handleReject(review)}
                  >
                    拒绝
                  </Button>
                </Space>,
              ]}
            >
              <List.Item.Meta
                avatar={<Avatar icon={<UserOutlined />} src={review.client?.avatar} />}
                title={
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <span>{review.client?.full_name}</span>
                    <Rate disabled value={review.rating} style={{ fontSize: 14 }} />
                    <Tag color="orange">待审核</Tag>
                  </div>
                }
                description={
                  <div>
                    <div style={{ color: '#666', marginBottom: 8 }}>
                      服务：{review.service?.title} | 专业人士：{review.professional?.full_name}
                    </div>
                    {review.content && <div style={{ background: '#f5f5f5', padding: 12, borderRadius: 4 }}>{review.content}</div>}
                    <div style={{ color: '#999', fontSize: 12, marginTop: 8 }}>
                      {new Date(review.created_at).toLocaleString()}
                    </div>
                  </div>
                }
              />
            </List.Item>
          )}
        />
      </Card>
    </div>
  )
}
