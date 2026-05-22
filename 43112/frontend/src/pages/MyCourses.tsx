import React, { useEffect, useState } from 'react'
import { List, Card, Rate, Tag, Empty, Progress, Row, Col } from 'antd'
import { BookOutlined, PlayCircleOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { orderApi, progressApi } from '@/services'
import { Order } from '@/types'

const MyCoursesPage: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([])
  const [progressMap, setProgressMap] = useState<Record<string, any>>({})
  const [loading, setLoading] = useState(false)

  const loadData = async () => {
    setLoading(true)
    try {
      const res = await orderApi.myOrders({ status: 'paid', page_size: 50 })
      if (res.code === 0 && res.data) {
        setOrders(res.data.items)
        for (const order of res.data.items) {
          try {
            const progressRes = await progressApi.getCourseProgress(order.course_id)
            if (progressRes.code === 0 && progressRes.data) {
              setProgressMap((prev) => ({ ...prev, [order.course_id]: progressRes.data }))
            }
          } catch (e) { /* ignore */ }
        }
      }
    } catch (error) {
      console.error('Failed to load courses:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [])

  if (!loading && orders.length === 0) {
    return <Empty description="还没有购买课程" />
  }

  return (
    <div>
      <h2>我的课程</h2>
      <Row gutter={[16, 16]}>
        {orders.map((order) => (
          <Col xs={24} sm={12} md={8} lg={6} key={order.id}>
            <Link to={`/courses/${order.course_id}`}>
              <Card
                hoverable
                className="card-hover"
                cover={
                  order.course?.cover ? (
                    <img
                      alt={order.course?.title}
                      src={order.course.cover}
                      style={{ height: 140, objectFit: 'cover' }}
                    />
                  ) : (
                    <div style={{ height: 140, background: '#1890ff', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <PlayCircleOutlined style={{ fontSize: 48, color: '#fff' }} />
                    </div>
                  )
                }
              >
                <Card.Meta
                  title={
                    <div style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                      {order.course?.title}
                    </div>
                  }
                  description={
                    <div>
                      <Tag color="blue">{order.course?.category}</Tag>
                      <div style={{ marginTop: 8 }}>
                        <Progress
                          percent={progressMap[order.course_id]?.completion_rate || 0}
                          size="small"
                          format={(p) => `${p}%`}
                        />
                      </div>
                      <div style={{ fontSize: 12, color: '#999', marginTop: 4 }}>
                        购买于 {new Date(order.created_at).toLocaleDateString()}
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

export default MyCoursesPage
