import React, { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Row, Col, Card, Button, Tag, Rate, Tabs, List, Avatar,
  Typography, Modal, Form, Input, InputNumber, message, Descriptions,
  Empty, Progress as AntProgress, Space, Popconfirm,
} from 'antd'
import {
  PlayCircleOutlined, ShoppingCartOutlined, UserOutlined,
  StarOutlined, LikeOutlined, FileTextOutlined,
} from '@ant-design/icons'
import ReactPlayer from 'react-player'
import { courseApi, orderApi, couponApi, reviewApi, progressApi } from '@/services'
import { useAuthStore } from '@/store/auth'
import { Course, Review, Progress } from '@/types'

const { Title, Text, Paragraph } = Typography
const { TextArea } = Input

const CourseDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user, isAuthenticated } = useAuthStore()
  const [course, setCourse] = useState<Course | null>(null)
  const [loading, setLoading] = useState(true)
  const [reviews, setReviews] = useState<Review[]>([])
  const [myReview, setMyReview] = useState<Review | null>(null)
  const [reviewSummary, setReviewSummary] = useState<any>(null)
  const [courseProgress, setCourseProgress] = useState<any>(null)
  const [currentLesson, setCurrentLesson] = useState<any>(null)
  const [reviewModalVisible, setReviewModalVisible] = useState(false)
  const [orderModalVisible, setOrderModalVisible] = useState(false)
  const [couponCode, setCouponCode] = useState('')
  const [couponInfo, setCouponInfo] = useState<any>(null)
  const [isPurchased, setIsPurchased] = useState(false)
  const [reviewForm] = Form.useForm()

  const loadCourse = async () => {
    if (!id) return
    setLoading(true)
    try {
      const res = await courseApi.get(id)
      if (res.code === 0 && res.data) {
        setCourse(res.data)
      }
    } catch (error) {
      console.error('Failed to load course:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadReviews = async () => {
    if (!id) return
    try {
      const res = await reviewApi.list({ course_id: id, page: 1, page_size: 20 })
      if (res.code === 0 && res.data) {
        setReviews(res.data.items)
      }
    } catch (error) {
      console.error('Failed to load reviews:', error)
    }
  }

  const loadMyReview = async () => {
    if (!id || !isAuthenticated) return
    try {
      const res = await reviewApi.getMy(id)
      if (res.code === 0 && res.data) {
        setMyReview(res.data)
      }
    } catch (error) {
      console.error('Failed to load my review:', error)
    }
  }

  const loadReviewSummary = async () => {
    if (!id) return
    try {
      const res = await reviewApi.summary(id)
      if (res.code === 0 && res.data) {
        setReviewSummary(res.data)
      }
    } catch (error) {
      console.error('Failed to load review summary:', error)
    }
  }

  const loadProgress = async () => {
    if (!id || !isAuthenticated) return
    try {
      const res = await progressApi.getCourseProgress(id)
      if (res.code === 0 && res.data) {
        setCourseProgress(res.data)
      }
    } catch (error) {
      console.error('Failed to load progress:', error)
    }
  }

  const checkPurchased = async () => {
    if (!id || !isAuthenticated) return
    try {
      const res = await orderApi.myOrders({ course_id: id, status: 'paid' })
      if (res.code === 0 && res.data && res.data.total > 0) {
        setIsPurchased(true)
      }
    } catch (error) {
      console.error('Failed to check purchase:', error)
    }
  }

  useEffect(() => {
    loadCourse()
    loadReviews()
    loadReviewSummary()
    if (isAuthenticated) {
      loadMyReview()
      checkPurchased()
      loadProgress()
    }
  }, [id, isAuthenticated])

  const handleValidateCoupon = async () => {
    if (!id || !couponCode) return
    try {
      const res = await couponApi.validate(couponCode, id)
      if (res.code === 0 && res.data) {
        setCouponInfo(res.data)
        message.success(`优惠券可用，优惠：¥${res.data.discount.toFixed(2)}`)
      }
    } catch (error: any) {
      message.error(error.message || '优惠券无效')
    }
  }

  const handlePurchase = async () => {
    if (!id) return
    try {
      const data: any = { course_id: id }
      if (couponCode) data.coupon_code = couponCode
      const res = await orderApi.create(data)
      if (res.code === 0 && res.data) {
        const payRes = await orderApi.pay(res.data.id)
        if (payRes.code === 0) {
          message.success('购买成功')
          setOrderModalVisible(false)
          setIsPurchased(true)
          checkPurchased()
        }
      }
    } catch (error: any) {
      message.error(error.message || '购买失败')
    }
  }

  const handleSubmitReview = async () => {
    try {
      const values = await reviewForm.validateFields()
      const res = await reviewApi.create({ ...values, course_id: id })
      if (res.code === 0) {
        message.success('评价提交成功')
        setReviewModalVisible(false)
        loadReviews()
        loadMyReview()
        loadReviewSummary()
        loadCourse()
      }
    } catch (error: any) {
      message.error(error.message || '评价提交失败')
    }
  }

  if (loading) {
    return <div style={{ textAlign: 'center', padding: 48 }}>加载中...</div>
  }

  if (!course) {
    return <Empty description="课程不存在" />
  }

  const tabs = [
    {
      key: 'intro',
      label: '课程介绍',
      children: (
        <Card>
          <Descriptions column={2} bordered size="small">
            <Descriptions.Item label="分类">{course.category}</Descriptions.Item>
            <Descriptions.Item label="难度">
              {course.level === 'beginner' ? '入门' : course.level === 'intermediate' ? '中级' : '高级'}
            </Descriptions.Item>
            <Descriptions.Item label="时长">{course.total_hours} 小时</Descriptions.Item>
            <Descriptions.Item label="学员数">{course.student_count} 人</Descriptions.Item>
          </Descriptions>
          <div style={{ marginTop: 16 }}>
            <Title level={5}>课程简介</Title>
            <div className="rich-content" dangerouslySetInnerHTML={{ __html: course.description }} />
          </div>
          {course.tags && (
            <div style={{ marginTop: 16 }}>
              {course.tags.split(',').map((tag, i) => (
                <Tag key={i} color="blue">{tag}</Tag>
              ))}
            </div>
          )}
        </Card>
      ),
    },
    {
      key: 'chapters',
      label: '课程目录',
      children: (
        <Card>
          {course.chapters?.map((chapter) => (
            <div key={chapter.id} style={{ marginBottom: 24 }}>
              <Title level={5}>{chapter.position}. {chapter.title}</Title>
              <List
                size="small"
                dataSource={chapter.lessons}
                renderItem={(lesson) => (
                  <List.Item
                    actions={[
                      !isPurchased && !lesson.is_free && !course.is_free
                        ? <Tag color="gold">付费</Tag>
                        : null,
                    ]}
                  >
                    <List.Item.Meta
                      avatar={
                        lesson.type === 'video' ? <PlayCircleOutlined /> :
                        lesson.type === 'document' ? <FileTextOutlined /> :
                        <StarOutlined />
                      }
                      title={
                        <a onClick={() => {
                          if (isPurchased || lesson.is_free || course.is_free) {
                            setCurrentLesson(lesson)
                          } else {
                            message.info('请先购买课程')
                          }
                        }}>
                          {lesson.title}
                        </a>
                      }
                      description={
                        lesson.type === 'video'
                          ? `${Math.floor(lesson.video_length / 60)}分${lesson.video_length % 60}秒`
                          : lesson.type === 'quiz' ? '测验' : '文档'
                      }
                    />
                  </List.Item>
                )}
              />
            </div>
          ))}
          {(!course.chapters || course.chapters.length === 0) && (
            <Empty description="暂无章节" />
          )}
        </Card>
      ),
    },
    {
      key: 'reviews',
      label: `评价 (${course.review_count})`,
      children: (
        <Card>
          {reviewSummary && (
            <div style={{ textAlign: 'center', marginBottom: 24 }}>
              <Title level={2} style={{ color: '#faad14' }}>{reviewSummary.avg_rating?.toFixed(1) || '0.0'}</Title>
              <Rate disabled allowHalf defaultValue={reviewSummary.avg_rating || 0} />
              <div style={{ color: '#999' }}>共 {reviewSummary.review_count} 条评价</div>
            </div>
          )}
          {isAuthenticated && isPurchased && !myReview && (
            <div style={{ textAlign: 'center', marginBottom: 16 }}>
              <Button type="primary" onClick={() => setReviewModalVisible(true)}>
                写评价
              </Button>
            </div>
          )}
          <List
            dataSource={reviews}
            renderItem={(review) => (
              <List.Item key={review.id}>
                <List.Item.Meta
                  avatar={<Avatar icon={<UserOutlined />} src={!review.is_anonymous ? review.user?.avatar : undefined} />}
                  title={
                    <Space>
                      <span>{review.is_anonymous ? '匿名用户' : review.user?.nickname}</span>
                      <Rate disabled allowHalf defaultValue={review.rating} style={{ fontSize: 12 }} />
                    </Space>
                  }
                  description={
                    <div>
                      <Paragraph style={{ marginBottom: 0 }}>{review.content}</Paragraph>
                      <Text type="secondary" style={{ fontSize: 12 }}>
                        {new Date(review.created_at).toLocaleDateString()}
                      </Text>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
          {reviews.length === 0 && <Empty description="暂无评价" />}
        </Card>
      ),
    },
  ]

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          {currentLesson ? (
            <Card style={{ marginBottom: 16 }}>
              <div style={{ position: 'relative', paddingTop: '56.25%' }}>
                {currentLesson.type === 'video' && currentLesson.video_url ? (
                  <div style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%' }}>
                    <ReactPlayer
                      url={currentLesson.video_url}
                      controls
                      width="100%"
                      height="100%"
                      onProgress={(state) => {
                        progressApi.update(currentLesson.id, {
                          last_position: state.playedSeconds,
                          total_duration: state.loadedSeconds,
                          is_completed: state.played >= 0.95,
                        })
                      }}
                    />
                  </div>
                ) : currentLesson.type === 'document' && currentLesson.doc_url ? (
                  <div style={{ textAlign: 'center', padding: 48 }}>
                    <FileTextOutlined style={{ fontSize: 64, color: '#1890ff' }} />
                    <p><a href={currentLesson.doc_url} target="_blank" rel="noreferrer">
                      {currentLesson.doc_name || '下载文档'}
                    </a></p>
                  </div>
                ) : (
                  <div style={{ textAlign: 'center', padding: 48, color: '#999' }}>
                    暂无内容
                  </div>
                )}
              </div>
              <Title level={5} style={{ marginTop: 16 }}>{currentLesson.title}</Title>
              <Button onClick={() => setCurrentLesson(null)}>返回目录</Button>
            </Card>
          ) : (
            <Card cover={
              course.cover ? (
                <img alt={course.title} src={course.cover} style={{ height: 300, objectFit: 'cover' }} />
              ) : (
                <div style={{ height: 300, background: 'linear-gradient(135deg, #1890ff, #722ed1)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                  <PlayCircleOutlined style={{ fontSize: 64, color: '#fff' }} />
                </div>
              )
            }>
              <Title level={3}>{course.title}</Title>
              <Paragraph type="secondary">{course.subtitle}</Paragraph>
              <Space>
                <Tag color="blue">{course.category}</Tag>
                <Tag color={course.level === 'beginner' ? 'green' : course.level === 'intermediate' ? 'orange' : 'red'}>
                  {course.level === 'beginner' ? '入门' : course.level === 'intermediate' ? '中级' : '高级'}
                </Tag>
                <Rate disabled allowHalf defaultValue={course.avg_rating} />
                <Text>{course.avg_rating.toFixed(1)} ({course.review_count} 评价)</Text>
                <Text>{course.student_count} 人学习</Text>
              </Space>
            </Card>
          )}
          <Tabs defaultActiveKey="intro" items={tabs} style={{ marginTop: 16 }} />
        </Col>

        <Col xs={24} lg={8}>
          <Card style={{ position: 'sticky', top: 24 }}>
            <Title level={4} style={{ color: '#f5222d' }}>
              {course.is_free ? '免费课程' : `¥${course.price}`}
            </Title>
            {!course.is_free && course.original_price > course.price && (
              <Text delete type="secondary">¥{course.original_price}</Text>
            )}
            {isPurchased && (
              <Tag color="green" style={{ marginLeft: 8 }}>已购买</Tag>
            )}
            {courseProgress && (
              <div style={{ marginBottom: 16 }}>
                <Text>学习进度：{courseProgress.completion_rate?.toFixed(1) || 0}%</Text>
                <AntProgress percent={courseProgress.completion_rate || 0} size="small" />
              </div>
            )}
            {!course.is_free && !isPurchased && (
              <>
                <Space.Compact style={{ width: '100%', marginBottom: 16 }}>
                  <Input
                    placeholder="优惠券码"
                    value={couponCode}
                    onChange={(e) => setCouponCode(e.target.value)}
                  />
                  <Button onClick={handleValidateCoupon}>验证</Button>
                </Space.Compact>
                {couponInfo && (
                  <Text type="success" style={{ display: 'block', marginBottom: 16 }}>
                    优惠：¥{couponInfo.discount.toFixed(2)}，实付：¥{couponInfo.final.toFixed(2)}
                  </Text>
                )}
                <Button
                  type="primary"
                  size="large"
                  icon={<ShoppingCartOutlined />}
                  block
                  onClick={() => {
                    if (!isAuthenticated) {
                      navigate('/login')
                    } else {
                      setOrderModalVisible(true)
                    }
                  }}
                >
                  立即购买
                </Button>
              </>
            )}
            {course.is_free && !isPurchased && (
              <Button type="primary" size="large" block onClick={() => setIsPurchased(true)}>
                免费学习
              </Button>
            )}
            {isPurchased && (
              <Button type="primary" size="large" block onClick={() => setCurrentLesson(course.chapters?.[0]?.lessons?.[0])}>
                继续学习
              </Button>
            )}
            {course.instructor && (
              <Card size="small" style={{ marginTop: 16 }}>
                <Space>
                  <Avatar src={course.instructor.avatar} icon={<UserOutlined />} />
                  <div>
                    <div><Text strong>{course.instructor.nickname}</Text></div>
                    <Text type="secondary" style={{ fontSize: 12 }}>{course.instructor.bio || '讲师'}</Text>
                  </div>
                </Space>
              </Card>
            )}
          </Card>
        </Col>
      </Row>

      <Modal
        title="购买课程"
        open={orderModalVisible}
        onCancel={() => setOrderModalVisible(false)}
        onOk={handlePurchase}
        okText="确认支付"
      >
        <Descriptions column={1} size="small" bordered>
          <Descriptions.Item label="课程名称">{course.title}</Descriptions.Item>
          <Descriptions.Item label="原价">¥{course.price}</Descriptions.Item>
          <Descriptions.Item label="优惠">¥{couponInfo?.discount?.toFixed(2) || '0.00'}</Descriptions.Item>
          <Descriptions.Item label="实付">
            <Text strong style={{ color: '#f5222d', fontSize: 18 }}>
              ¥{couponInfo?.final?.toFixed(2) || course.price.toFixed(2)}
            </Text>
          </Descriptions.Item>
        </Descriptions>
      </Modal>

      <Modal
        title="写评价"
        open={reviewModalVisible}
        onCancel={() => setReviewModalVisible(false)}
        footer={null}
      >
        <Form form={reviewForm} onFinish={handleSubmitReview} layout="vertical">
          <Form.Item name="rating" label="评分" rules={[{ required: true }]} initialValue={5}>
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容" rules={[{ required: true, message: '请输入评价内容' }]}>
            <TextArea rows={4} placeholder="分享你的学习体验..." />
          </Form.Item>
          <Form.Item name="is_anonymous" valuePropName="checked" initialValue={false}>
            匿名评价
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交评价</Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CourseDetailPage
