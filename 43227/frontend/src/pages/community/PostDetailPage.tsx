import { useEffect, useState } from 'react'
import { Card, Button, Space, Form, Input, List, Avatar, message, Rate } from 'antd'
import { ArrowLeftOutlined, UserOutlined, LikeOutlined } from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../../api'
import type { Post, Comment } from '../../types'
import { useAuthStore } from '../../store/authStore'
import dayjs from 'dayjs'

function PostDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [data, setData] = useState<Post | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get(`/posts/${id}`)
      setData(response.data as Post)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchComments = async () => {
    try {
      const response = await api.get(`/posts/${id}/comments`)
      setComments(response.data || [])
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    if (id) {
      fetchData()
      fetchComments()
    }
  }, [id])

  const handleLike = async () => {
    try {
      await api.post(`/posts/${id}/like`)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleComment = async (values: any) => {
    try {
      await api.post('/comments', {
        post_id: id,
        content: values.content,
      })
      message.success('评论成功')
      form.resetFields()
      fetchComments()
    } catch (error: any) {
      message.error(error.message || '评论失败')
    }
  }

  if (loading) {
    return <Card loading />
  }

  if (!data) {
    return <Card>帖子不存在</Card>
  }

  return (
    <div>
      <Space style={{ marginBottom: 16 }}>
        <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(-1)}>
          返回
        </Button>
      </Space>
      <Card>
        <h1>{data.title}</h1>
        <Space style={{ marginBottom: 16 }}>
          <Avatar icon={<UserOutlined />} src={data.user?.avatar} />
          <span>{data.user?.username}</span>
          <span>发布于: {dayjs(data.created_at).format('YYYY-MM-DD HH:mm')}</span>
          <span>浏览: {data.view_count}</span>
        </Space>
        <div style={{ whiteSpace: 'pre-wrap', fontSize: 16, lineHeight: 1.8 }}>
          {data.content}
        </div>
        <Space style={{ marginTop: 16 }}>
          <Button icon={<LikeOutlined />} onClick={handleLike}>
            {data.like_count}
          </Button>
        </Space>
      </Card>

      <Card title="评论" style={{ marginTop: 16 }}>
        {user && (
          <Form form={form} onFinish={handleComment} style={{ marginBottom: 24 }}>
            <Form.Item name="content" rules={[{ required: true, message: '请输入评论内容' }]}>
              <Input.TextArea rows={3} placeholder="写下你的评论..." />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit">发表评论</Button>
            </Form.Item>
          </Form>
        )}
        <List
          dataSource={comments}
          renderItem={(item) => (
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar icon={<UserOutlined />} src={item.user?.avatar} />}
                title={
                  <Space>
                    <span>{item.user?.username}</span>
                    <span style={{ color: '#999' }}>{dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}</span>
                  </Space>
                }
                description={item.content}
              />
            </List.Item>
          )}
        />
      </Card>
    </div>
  )
}

export default PostDetailPage
