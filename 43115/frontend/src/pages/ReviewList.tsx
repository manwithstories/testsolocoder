import React, { useState, useEffect } from 'react'
import { Table, Card, Rate, Avatar, Button, Modal, Input, message } from 'antd'
import { reviewApi } from '@/services/review'
import { Review } from '@/types'
import { formatDate } from '@/utils'
import { usePagination } from '@/hooks'
import { useAppSelector } from '@/store/hooks'

const { TextArea } = Input

const ReviewList: React.FC = () => {
  const { userInfo } = useAppSelector((state) => state.auth)
  const { page, pageSize, total, setPage, setPageSize, setTotal } = usePagination()
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(false)
  const [replyModalVisible, setReplyModalVisible] = useState(false)
  const [replyReview, setReplyReview] = useState<Review | null>(null)
  const [replyContent, setReplyContent] = useState('')

  useEffect(() => {
    loadReviews()
  }, [page, pageSize])

  const loadReviews = async () => {
    setLoading(true)
    try {
      const res = await reviewApi.getList({
        page,
        page_size: pageSize,
      })
      setReviews(res.list)
      setTotal(res.total)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const handleReply = async () => {
    if (!replyReview || !replyContent) {
      message.warning('请填写回复内容')
      return
    }
    try {
      await reviewApi.reply(replyReview.id, { content: replyContent })
      message.success('回复成功')
      setReplyModalVisible(false)
      setReplyReview(null)
      setReplyContent('')
      loadReviews()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: '评价ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '评分',
      dataIndex: 'overall_rating',
      key: 'overall_rating',
      render: (text: number) => <Rate disabled value={text} allowHalf />,
    },
    {
      title: '评价内容',
      dataIndex: 'content',
      key: 'content',
      ellipsis: true,
    },
    {
      title: '评价人',
      key: 'reviewer',
      render: (_: any, record: Review) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <Avatar size="small" src={record.reviewer?.avatar}>
            {record.reviewer?.nickname?.charAt(0)}
          </Avatar>
          <span>{record.reviewer?.nickname}</span>
        </div>
      ),
    },
    {
      title: '被评价人',
      key: 'reviewee',
      render: (_: any, record: Review) => (
        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
          <Avatar size="small" src={record.reviewee?.avatar}>
            {record.reviewee?.nickname?.charAt(0)}
          </Avatar>
          <span>{record.reviewee?.nickname}</span>
        </div>
      ),
    },
    {
      title: '评价时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => formatDate(text),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Review) => (
        <div>
          {userInfo?.role === 'service_provider' &&
            record.reviewee_id === userInfo.id &&
            !record.reply_content && (
              <Button
                type="link"
                onClick={() => {
                  setReplyReview(record)
                  setReplyModalVisible(true)
                }}
              >
                回复
              </Button>
            )}
          {record.reply_content && (
            <div style={{ color: '#999', fontSize: 12 }}>已回复</div>
          )}
        </div>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">评价管理</h1>
      </div>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={reviews}
          columns={columns}
          pagination={{
            current: page,
            pageSize,
            total,
            showSizeChanger: true,
            onChange: (p, ps) => {
              setPage(p)
              setPageSize(ps)
            },
          }}
        />
      </Card>

      <Modal
        title="回复评价"
        open={replyModalVisible}
        onOk={handleReply}
        onCancel={() => {
          setReplyModalVisible(false)
          setReplyReview(null)
          setReplyContent('')
        }}
      >
        {replyReview && (
          <div style={{ marginBottom: 16 }}>
            <div>原评价：{replyReview.content}</div>
            <div style={{ marginTop: 8 }}>
              <Rate disabled value={replyReview.overall_rating} allowHalf />
            </div>
          </div>
        )}
        <TextArea
          rows={3}
          placeholder="请输入回复内容"
          value={replyContent}
          onChange={(e) => setReplyContent(e.target.value)}
        />
      </Modal>
    </div>
  )
}

export default ReviewList
