import React, { useEffect, useState } from 'react'
import { Row, Col, Card, Pagination, Input, Select, Space, Tag, Rate } from 'antd'
import { PlayCircleOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { courseApi } from '@/services'
import { Course, PaginatedData } from '@/types'

const { Search } = Input
const { Option } = Select

const CourseListPage: React.FC = () => {
  const [data, setData] = useState<PaginatedData<Course>>({
    items: [], total: 0, page: 1, page_size: 12, total_pages: 0,
  })
  const [loading, setLoading] = useState(false)
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')
  const [level, setLevel] = useState('')
  const [categories, setCategories] = useState<any[]>([])

  const loadCourses = async (page = 1) => {
    setLoading(true)
    try {
      const params: any = { page, page_size: 12 }
      if (search) params.search = search
      if (category) params.category = category
      if (level) params.level = level
      const res = await courseApi.list(params)
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
    courseApi.categories().then((res) => {
      if (res.code === 0 && res.data) setCategories(res.data)
    })
  }, [])

  return (
    <div>
      <h2>课程列表</h2>
      <Space style={{ marginBottom: 16 }} wrap>
        <Search
          placeholder="搜索课程..."
          style={{ width: 300 }}
          onSearch={(value) => { setSearch(value); loadCourses(1) }}
          allowClear
        />
        <Select
          placeholder="全部分类"
          style={{ width: 160 }}
          allowClear
          value={category || undefined}
          onChange={(value) => { setCategory(value || ''); loadCourses(1) }}
        >
          {categories.map((cat) => (
            <Option key={cat.category} value={cat.category}>{cat.category}</Option>
          ))}
        </Select>
        <Select
          placeholder="全部难度"
          style={{ width: 140 }}
          allowClear
          value={level || undefined}
          onChange={(value) => { setLevel(value || ''); loadCourses(1) }}
        >
          <Option value="beginner">入门</Option>
          <Option value="intermediate">中级</Option>
          <Option value="advanced">高级</Option>
        </Select>
      </Space>
      <Row gutter={[16, 16]}>
        {data.items.map((course) => (
          <Col xs={24} sm={12} md={8} lg={6} key={course.id}>
            <Link to={`/courses/${course.id}`}>
              <Card
                hoverable
                loading={loading}
                className="card-hover"
                cover={
                  course.cover ? (
                    <img alt={course.title} src={course.cover} style={{ height: 160, objectFit: 'cover' }} />
                  ) : (
                    <div style={{ height: 160, background: '#1890ff', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <PlayCircleOutlined style={{ fontSize: 48, color: '#fff' }} />
                    </div>
                  )
                }
              >
                <Card.Meta
                  title={<div style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', fontSize: 16 }}>{course.title}</div>}
                  description={
                    <div>
                      <Tag color="blue" style={{ marginBottom: 8 }}>{course.category}</Tag>
                      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <span style={{ color: '#f5222d', fontWeight: 600 }}>
                          {course.is_free ? '免费' : `¥${course.price}`}
                        </span>
                        <Space size={4}>
                          <Rate disabled allowHalf defaultValue={course.avg_rating} style={{ fontSize: 12 }} />
                          <span style={{ fontSize: 12, color: '#999' }}>({course.review_count})</span>
                        </Space>
                      </div>
                    </div>
                  }
                />
              </Card>
            </Link>
          </Col>
        ))}
      </Row>
      {data.total > 0 && (
        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Pagination
            current={data.page}
            total={data.total}
            pageSize={data.page_size}
            onChange={(page) => loadCourses(page)}
            showSizeChanger={false}
          />
        </div>
      )}
    </div>
  )
}

export default CourseListPage
