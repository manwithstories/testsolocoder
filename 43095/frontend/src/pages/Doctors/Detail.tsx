import React, { useState, useEffect } from 'react'
import {
  Card,
  Avatar,
  Rate,
  Button,
  Spin,
  Alert,
  Tag,
  Table,
  Space,
  Descriptions,
  List,
  Empty,
  Row,
  Col,
  Statistic,
  Modal,
  Pagination
} from 'antd'
import {
  UserOutlined,
  ArrowLeftOutlined,
  CalendarOutlined,
  ClockCircleOutlined,
  StarOutlined
} from '@ant-design/icons'
import { useParams, useNavigate } from 'react-router-dom'
import dayjs from 'dayjs'
import { doctorAPI, reviewAPI } from '@/services/api'
import type { Doctor, Schedule, Review, DayOfWeek } from '@/types'

const weekDays: Record<DayOfWeek, string> = {
  0: '周日',
  1: '周一',
  2: '周二',
  3: '周三',
  4: '周四',
  5: '周五',
  6: '周六'
}

const DoctorDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [doctor, setDoctor] = useState<Doctor | null>(null)
  const [schedules, setSchedules] = useState<Schedule[]>([])
  const [reviews, setReviews] = useState<Review[]>([])
  const [loading, setLoading] = useState(false)
  const [reviewLoading, setReviewLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [reviewPage, setReviewPage] = useState(1)
  const [reviewPageSize] = useState(5)
  const [reviewTotal, setReviewTotal] = useState(0)
  const [bookingModalVisible, setBookingModalVisible] = useState(false)

  useEffect(() => {
    if (id) {
      fetchDoctorDetail(parseInt(id))
      fetchSchedules(parseInt(id))
    }
  }, [id])

  useEffect(() => {
    if (id) {
      fetchReviews(parseInt(id))
    }
  }, [id, reviewPage])

  const fetchDoctorDetail = async (doctorId: number) => {
    setLoading(true)
    setError(null)
    try {
      const response = await doctorAPI.getDetail(doctorId)
      setDoctor(response)
    } catch (err: any) {
      setError(err.message || '获取医生详情失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchSchedules = async (doctorId: number) => {
    try {
      const response = await doctorAPI.getSchedules(doctorId)
      setSchedules(response)
    } catch (err) {
      console.error('获取排班信息失败:', err)
    }
  }

  const fetchReviews = async (doctorId: number) => {
    setReviewLoading(true)
    try {
      const response = await reviewAPI.getList({
        page: reviewPage,
        pageSize: reviewPageSize,
        doctor_id: doctorId
      })
      setReviews(response.list)
      setReviewTotal(response.total)
    } catch (err) {
      console.error('获取评价列表失败:', err)
    } finally {
      setReviewLoading(false)
    }
  }

  const handleBack = () => {
    navigate('/doctors')
  }

  const handleBook = () => {
    navigate(`/appointments/create?doctor_id=${id}`)
  }

  const scheduleColumns = [
    {
      title: '星期',
      dataIndex: 'day_of_week',
      key: 'day_of_week',
      render: (day: DayOfWeek) => <span className="font-medium">{weekDays[day]}</span>
    },
    {
      title: '开始时间',
      dataIndex: 'start_time',
      key: 'start_time'
    },
    {
      title: '结束时间',
      dataIndex: 'end_time',
      key: 'end_time'
    },
    {
      title: '时段时长',
      dataIndex: 'time_slot_minutes',
      key: 'time_slot_minutes',
      render: (minutes: number) => `${minutes}分钟`
    },
    {
      title: '最大接诊数',
      dataIndex: 'max_patients',
      key: 'max_patients'
    },
    {
      title: '状态',
      dataIndex: 'is_available',
      key: 'is_available',
      render: (available: boolean) =>
        available ? (
          <Tag color="green">可预约</Tag>
        ) : (
          <Tag color="red">暂停接诊</Tag>
        )
    }
  ]

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={handleBack}
        className="mb-4"
        type="text"
      >
        返回医生列表
      </Button>

      {error && (
        <Alert
          message="错误"
          description={error}
          type="error"
          showIcon
          className="mb-4"
          closable
          onClose={() => setError(null)}
        />
      )}

      <Spin spinning={loading} tip="加载中...">
        {doctor && (
          <>
            <Card className="mb-6">
              <Row gutter={24} align="middle">
                <Col flex="120px">
                  <Avatar
                    size={120}
                    src={doctor.user?.avatar_url}
                    icon={!doctor.user?.avatar_url && <UserOutlined style={{ fontSize: 60 }} />}
                    className="block mx-auto"
                  />
                </Col>
                <Col flex="auto">
                  <div className="mb-4">
                    <Space size="large" align="center">
                      <h1 className="text-3xl font-bold mb-0">
                        {doctor.user?.full_name}
                      </h1>
                      <Tag color="blue" className="text-base py-1 px-3">
                        {doctor.title}
                      </Tag>
                    </Space>
                    <p className="text-gray-500 mt-2 text-lg">
                      {doctor.department?.name}
                    </p>
                  </div>

                  <Row gutter={32} className="mb-4">
                    <Col>
                      <Statistic
                        title="平均评分"
                        value={doctor.average_rating}
                        precision={1}
                        prefix={<StarOutlined style={{ color: '#faad14' }} />}
                        suffix="/ 5.0"
                      />
                    </Col>
                    <Col>
                      <Statistic
                        title="评价数量"
                        value={doctor.review_count}
                        suffix="条"
                      />
                    </Col>
                    <Col>
                      <Statistic
                        title="挂号费"
                        value={doctor.registration_fee}
                        prefix="¥"
                      />
                    </Col>
                    <Col>
                      <Statistic
                        title="问诊费"
                        value={doctor.consultation_fee}
                        prefix="¥"
                      />
                    </Col>
                  </Row>

                  <Space size="middle">
                    <Button
                      type="primary"
                      size="large"
                      icon={<CalendarOutlined />}
                      onClick={handleBook}
                    >
                      立即预约
                    </Button>
                  </Space>
                </Col>
              </Row>
            </Card>

            <Card title="个人简介" className="mb-6">
              <Descriptions column={1} bordered size="middle">
                <Descriptions.Item label="擅长领域">
                  {doctor.specialty || '暂无信息'}
                </Descriptions.Item>
                <Descriptions.Item label="个人介绍">
                  {doctor.introduction || '暂无介绍'}
                </Descriptions.Item>
              </Descriptions>
            </Card>

            <Card
              title={
                <Space>
                  <ClockCircleOutlined />
                  排班信息
                </Space>
              }
              className="mb-6"
            >
              {schedules.length > 0 ? (
                <Table
                  columns={scheduleColumns}
                  dataSource={schedules}
                  rowKey="id"
                  pagination={false}
                />
              ) : (
                <Empty description="暂无排班信息" />
              )}
            </Card>

            <Card
              title={
                <Space>
                  <StarOutlined />
                  患者评价
                  <Tag color="blue">{reviewTotal}条</Tag>
                </Space>
              }
            >
              <Spin spinning={reviewLoading}>
                {reviews.length > 0 ? (
                  <>
                    <List
                      dataSource={reviews}
                      renderItem={(review) => (
                        <List.Item key={review.id}>
                          <List.Item.Meta
                            avatar={
                              <Avatar
                                src={review.is_anonymous ? undefined : ''}
                                icon={<UserOutlined />}
                              />
                            }
                            title={
                              <Space>
                                <span>
                                  {review.is_anonymous ? '匿名用户' : '患者'}
                                </span>
                                <Rate
                                  disabled
                                  value={review.rating}
                                  allowHalf
                                  className="text-xs"
                                />
                              </Space>
                            }
                            description={
                              <div>
                                <p className="text-gray-700 mb-1">
                                  {review.content || '暂无评价内容'}
                                </p>
                                <span className="text-xs text-gray-400">
                                  {dayjs(review.created_at).format('YYYY-MM-DD HH:mm')}
                                </span>
                              </div>
                            }
                          />
                        </List.Item>
                      )}
                    />
                    {reviewTotal > reviewPageSize && (
                      <div className="mt-4 flex justify-end">
                        <Pagination
                          current={reviewPage}
                          pageSize={reviewPageSize}
                          total={reviewTotal}
                          onChange={setReviewPage}
                          size="small"
                        />
                      </div>
                    )}
                  </>
                ) : (
                  <Empty description="暂无评价" />
                )}
              </Spin>
            </Card>
          </>
        )}
      </Spin>

      <Modal
        title="预约挂号"
        open={bookingModalVisible}
        onCancel={() => setBookingModalVisible(false)}
        footer={null}
      >
        <p className="text-center py-8">即将跳转到预约页面...</p>
      </Modal>
    </div>
  )
}

export default DoctorDetail
