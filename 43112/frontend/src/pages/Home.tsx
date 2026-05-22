import React, { useEffect, useState } from 'react'
import { Row, Col, Card, Typography, Tag, Rate, Space, Button, Input, Select } from 'antd'
import { PlayCircleOutlined, StarOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { courseApi } from '@/services'
import { Course } from '@/types'

const { Title, Paragraph } = Typography
const { Search } = Input
const { Option } = Select

const HomePage: React.FC = () => {
  const [courses, setCourses] = useState<Course[]>([])
  const [categories, setCategories] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState<string>('')
  const [level, setLevel] = useState<string>('')

  const loadCourses = async () => {
    setLoading(true)
    try {
      const params: any = { page: 1, page_size: 12 }
      if (search) params.search = search
      if (category) params.category = category
      if (level) params.level = level
      const res = await courseApi.list(params)
      if (res.code === 0 && res.data) {
        setCourses(res.data.items)
      }
    } catch (error) {
      console.error('Failed to load courses:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadCategories = async () => {
    try {
      const res = await courseApi.categories()
      if (res.code === 0 && res.data) {
        setCategories(res.data)
      }
    } catch (error) {
      console.error('Failed to load categories:', error)
    }
  }

  useEffect(() => {
    loadCourses()
    loadCategories()
  }, [])

  return (
    <div>
      <div
        style={{
          background: 'linear-gradient(135deg, #1890ff 0%, #722ed1 100%)',
          padding: '60px 24px',
          marginBottom: 24,
          borderRadius: 8,
          textAlign: 'center',
          color: '#fff',
        }}
      >
        <Title level={2} style={{ color: '#fff', marginBottom: 16 }}>
          在线学习与课程管理平台
        </Title>
        <Paragraph style={{ color: 'rgba(255,255,255,0.85)', fontSize: 16 }}>
          海量优质课程，随时随地学习，助你成就更好的自己
        </Paragraph>
        <Search
          placeholder="搜索课程..."
          size="large"
          style={{ maxWidth: 500, marginTop: 16 }}
          onSearch={(value) => {
            setSearch(value)
            loadCourses()
          }}
        />
      </div>

      <Space style={{ marginBottom: 16 }} wrap>
        <Select
          placeholder="全部分类"
          style={{ width: 160 }}
          allowClear
          value={category || undefined}
          onChange={(value) => {
            setCategory(value || '')
            loadCourses()
          }}
        >
          {categories.map((cat) => (
            <Option key={cat.category} value={cat.category}>
              {cat.category} ({cat.count})
            </Option>
          ))}
        </Select>
        <Select
          placeholder="全部难度"
          style={{ width: 140 }}
          allowClear
          value={level || undefined}
          onChange={(value) => {
            setLevel(value || '')
            loadCourses()
          }}
        >
          <Option value="beginner">入门</Option>
          <Option value="intermediate">中级</Option>
          <Option value="advanced">高级</Option>
        </Select>
      </Space>

      <Row gutter={[16, 16]}>
        {courses.map((course) => (
          <Col xs={24} sm={12} md={8} lg={6} key={course.id}>
            <Link to={`/courses/${course.id}`}>
              <Card
                hoverable
                loading={loading}
                className="card-hover"
                cover={
                  course.cover ? (
                    <img
                      alt={course.title}
                      src={course.cover}
                      style={{ height: 160, objectFit: 'cover' }}
                    />
                  ) : (
                    <div
                      style={{
                        height: 160,
                        background: '#1890ff',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                      }}
                    >
                      <PlayCircleOutlined style={{ fontSize: 48, color: '#fff' }} />
                    </div>
                  )
                }
              >
                <Card.Meta
                  title={
                    <div
                      style={{
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        fontSize: 16,
                      }}
                    >
                      {course.title}
                    </div>
                  }
                  description={
                    <div>
                      <Tag color="blue" style={{ marginBottom: 8 }}>
                        {course.category}
                      </Tag>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <span style={{ color: '#f5222d', fontWeight: 600, fontSize: 16 }}>
                          {course.is_free ? '免费' : `¥${course.price}`}
                        </span>
                        <Space size={4}>
                          <Rate
                            disabled
                            allowHalf
                            defaultValue={course.avg_rating}
                            style={{ fontSize: 12 }}
                          />
                          <span style={{ fontSize: 12, color: '#999' }}>
                            ({course.review_count})
                          </span>
                        </Space>
                      </div>
                      <div style={{ fontSize: 12, color: '#999', marginTop: 4 }}>
                        {course.student_count} 人学习
                      </div>
                    </div>
                  }
                />
              </Card>
            </Link>
          </Col>
        ))}
      </Row>
    </div>
  )
}

export default HomePage
