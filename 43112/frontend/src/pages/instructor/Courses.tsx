import React, { useEffect, useState } from 'react'
import { Table, Button, Tag, Space, Popconfirm, Modal, Form, Input, Select, Upload, message, InputNumber, Row, Col } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, UploadOutlined } from '@ant-design/icons'
import { Link, useNavigate } from 'react-router-dom'
import { courseApi, uploadApi } from '@/services'
import { Course, PaginatedData } from '@/types'

const statusMap: Record<string, { color: string; text: string }> = {
  draft: { color: 'default', text: '草稿' },
  published: { color: 'green', text: '已上架' },
  offline: { color: 'orange', text: '已下架' },
  rejected: { color: 'red', text: '已拒绝' },
}

const InstructorCourses: React.FC = () => {
  const navigate = useNavigate()
  const [data, setData] = useState<PaginatedData<Course>>({
    items: [], total: 0, page: 1, page_size: 10, total_pages: 0,
  })
  const [loading, setLoading] = useState(false)

  const loadCourses = async (page = 1) => {
    setLoading(true)
    try {
      const res = await courseApi.myCourses({ page, page_size: 10 })
      if (res.code === 0 && res.data) {
        setData(res.data)
      }
    } catch (error) {
      console.error('Failed to load courses:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadCourses()
  }, [])

  const handleUpdateStatus = async (id: string, status: string) => {
    try {
      const res = await courseApi.updateStatus(id, status)
      if (res.code === 0) {
        message.success('状态更新成功')
        loadCourses()
      }
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleDelete = async (id: string) => {
    try {
      const res = await courseApi.delete(id)
      if (res.code === 0) {
        message.success('删除成功')
        loadCourses()
      }
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  const columns = [
    {
      title: '课程封面',
      dataIndex: 'cover',
      key: 'cover',
      render: (val: string) => val ? (
        <img src={val} alt="" style={{ width: 60, height: 40, objectFit: 'cover', borderRadius: 4 }} />
      ) : <div style={{ width: 60, height: 40, background: '#e8e8e8', borderRadius: 4 }} />,
    },
    {
      title: '课程名称',
      dataIndex: 'title',
      key: 'title',
    },
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (val: number, record: Course) => (
        <span style={{ color: '#f5222d' }}>
          {record.is_free ? '免费' : `¥${val}`}
        </span>
      ),
    },
    {
      title: '学员数',
      dataIndex: 'student_count',
      key: 'student_count',
    },
    {
      title: '评分',
      key: 'rating',
      render: (_: any, record: Course) => `${record.avg_rating.toFixed(1)} (${record.review_count})`,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const info = statusMap[status] || { color: 'default', text: status }
        return <Tag color={info.color}>{info.text}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Course) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => navigate(`/instructor/courses/${record.id}/edit`)}>
            编辑
          </Button>
          {record.status === 'draft' && (
            <Button type="link" onClick={() => handleUpdateStatus(record.id, 'published')}>
              上架
            </Button>
          )}
          {record.status === 'published' && (
            <Button type="link" onClick={() => handleUpdateStatus(record.id, 'offline')}>
              下架
            </Button>
          )}
          <Popconfirm title="确定删除？" onConfirm={() => handleDelete(record.id)}>
            <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <h2>我的课程</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/instructor/courses/new')}>
          创建课程
        </Button>
      </div>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={data.items}
        loading={loading}
        pagination={{
          current: data.page,
          total: data.total,
          pageSize: data.page_size,
          onChange: (page) => loadCourses(page),
        }}
      />
    </div>
  )
}

export default InstructorCourses
