import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Space, message, Card, Tag, List } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, LikeOutlined, MessageOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import api from '../../api'
import type { Post } from '../../types'
import dayjs from 'dayjs'

function CommunityPage() {
  const navigate = useNavigate()
  const [data, setData] = useState<Post[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState<Post | null>(null)
  const [categories, setCategories] = useState<any[]>([])
  const [form] = Form.useForm()

  const fetchCategories = async () => {
    try {
      const response = await api.get('/posts/categories')
      setCategories(response.data || [])
    } catch (error) {
      console.error(error)
    }
  }

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get('/posts', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as Post[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
    fetchCategories()
  }, [page, pageSize])

  const handleAdd = () => {
    setEditingRecord(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (record: Post) => {
    setEditingRecord(record)
    form.setFieldsValue({
      title: record.title,
      content: record.content,
      category: record.category,
      tags: record.tags,
    })
    setIsModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/posts/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  const handleLike = async (id: number) => {
    try {
      await api.post(`/posts/${id}/like`)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingRecord) {
        await api.put(`/posts/${editingRecord.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/posts', values)
        message.success('发布成功')
      }
      setIsModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  return (
    <Card
      title="蜂农社区"
      extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>发布帖子</Button>}
    >
      <List
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
        dataSource={data}
        renderItem={(item) => (
          <List.Item
            key={item.id}
            actions={[
              <Button type="link" icon={<LikeOutlined />} onClick={() => handleLike(item.id)}>
                {item.like_count}
              </Button>,
              <Button type="link" icon={<MessageOutlined />}>
                {item.comment_count}
              </Button>,
              <Button type="link" onClick={() => navigate(`/community/${item.id}`)}>查看</Button>,
            ]}
          >
            <List.Item.Meta
              title={
                <Space>
                  <a onClick={() => navigate(`/community/${item.id}`)}>{item.title}</a>
                  <Tag color="blue">{item.category}</Tag>
                  {item.tags?.map((tag, index) => (
                    <Tag key={index}>{tag}</Tag>
                  ))}
                </Space>
              }
              description={
                <Space>
                  <span>作者: {item.user?.username}</span>
                  <span>浏览: {item.view_count}</span>
                  <span>发布于: {dayjs(item.created_at).format('YYYY-MM-DD HH:mm')}</span>
                </Space>
              }
            />
          </List.Item>
        )}
      />
      <Modal
        title={editingRecord ? '编辑帖子' : '发布帖子'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
        width={700}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input placeholder="帖子标题" />
          </Form.Item>
          <Form.Item name="category" label="分类" rules={[{ required: true, message: '请选择分类' }]}>
            <Select placeholder="请选择分类">
              {categories.map((cat) => (
                <Select.Option key={cat.key} value={cat.key}>{cat.name}</Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="tags" label="标签">
            <Select mode="tags" placeholder="输入标签后回车" />
          </Form.Item>
          <Form.Item name="content" label="内容" rules={[{ required: true, message: '请输入内容' }]}>
            <Input.TextArea rows={8} placeholder="帖子内容" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交</Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default CommunityPage
