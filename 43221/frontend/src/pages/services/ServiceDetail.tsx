import { useState } from 'react'
import { Row, Col, Card, Rate, Button, DatePicker, TimePicker, Tag, Modal, Form, Input, message, List, Avatar, Divider } from 'antd'
import { DollarOutlined, ClockCircleOutlined, UserOutlined, CalendarOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useParams, useNavigate } from 'react-router-dom'
import dayjs from 'dayjs'
import { useAuthContext } from '@/contexts/AuthContext'
import { serviceApi } from '@/services/service'
import { appointmentApi } from '@/services/appointment'
import { recordApi } from '@/services/record'
import { Schedule } from '@/types'

export function ServiceDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuthContext()
  const queryClient = useQueryClient()
  const [selectedDate, setSelectedDate] = useState<dayjs.Dayjs | null>(null)
  const [selectedSchedule, setSelectedSchedule] = useState<Schedule | null>(null)
  const [bookingModalVisible, setBookingModalVisible] = useState(false)
  const [notes, setNotes] = useState('')

  const { data: service, isLoading } = useQuery({
    queryKey: ['service', id],
    queryFn: () => serviceApi.getById(id!),
    enabled: !!id,
  })

  const { data: schedules, isLoading: schedulesLoading } = useQuery({
    queryKey: ['schedules', id, selectedDate?.format('YYYY-MM-DD')],
    queryFn: () =>
      serviceApi.getSchedules(id!, {
        date: selectedDate?.format('YYYY-MM-DD'),
        only_available: true,
      }),
    enabled: !!id && !!selectedDate,
  })

  const { data: reviews } = useQuery({
    queryKey: ['service-reviews', id],
    queryFn: () => recordApi.getServiceReviews(id!, { page: 1, page_size: 10 }),
    enabled: !!id,
  })

  const createAppointmentMutation = useMutation({
    mutationFn: (data: { service_id: string; schedule_id: string; notes?: string }) =>
      appointmentApi.create(data),
    onSuccess: () => {
      message.success('预约成功！请及时完成支付')
      setBookingModalVisible(false)
      queryClient.invalidateQueries({ queryKey: ['schedules'] })
      navigate('/appointments')
    },
    onError: (error: any) => {
      message.error(error.message || '预约失败')
    },
  })

  const handleBook = () => {
    if (!selectedSchedule || !service) return

    createAppointmentMutation.mutate({
      service_id: service.id,
      schedule_id: selectedSchedule.id,
      notes,
    })
  }

  const disabledDate = (current: dayjs.Dayjs) => {
    return current && current < dayjs().startOf('day')
  }

  if (isLoading) {
    return <div className="loading-spin">加载中...</div>
  }

  if (!service) {
    return <div>服务不存在</div>
  }

  return (
    <div className="page-container">
      <Row gutter={[24, 24]}>
        <Col xs={24} lg={16}>
          <Card>
            <div style={{ display: 'flex', gap: 24, marginBottom: 24 }}>
              <div
                style={{
                  width: 120,
                  height: 120,
                  background: `linear-gradient(135deg, ${
                    service.service_type === 'legal'
                      ? '#667eea 0%, #764ba2 100%'
                      : service.service_type === 'counseling'
                      ? '#f093fb 0%, #f5576c 100%'
                      : service.service_type === 'financial'
                      ? '#4facfe 0%, #00f2fe 100%'
                      : '#43e97b 0%, #38f9d7 100%'
                  })`,
                  borderRadius: 8,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: 48,
                  color: '#fff',
                }}
              >
                {service.service_type === 'legal' ? '⚖️' :
                 service.service_type === 'counseling' ? '💬' :
                 service.service_type === 'financial' ? '💰' : '📋'}
              </div>
              <div style={{ flex: 1 }}>
                <h1 style={{ marginBottom: 8 }}>{service.title}</h1>
                <div style={{ marginBottom: 8 }}>
                  <Tag color="blue">{service.service_type}</Tag>
                  <Rate disabled allowHalf value={service.average_rating} style={{ fontSize: 14 }} />
                  <span style={{ marginLeft: 8, color: '#999' }}>({service.review_count}条评价)</span>
                </div>
                <div style={{ color: '#666', marginBottom: 16 }}>
                  <UserOutlined /> {service.professional?.full_name}
                </div>
                <div style={{ display: 'flex', gap: 24 }}>
                  <div>
                    <DollarOutlined style={{ color: '#ff4d4f' }} />
                    <span style={{ marginLeft: 8, fontSize: 20, fontWeight: 600, color: '#ff4d4f' }}>
                      ¥{service.price}
                    </span>
                  </div>
                  <div>
                    <ClockCircleOutlined style={{ color: '#1890ff' }} />
                    <span style={{ marginLeft: 8 }}>{service.duration_minutes}分钟</span>
                  </div>
                </div>
              </div>
            </div>

            <Divider />

            <h3>服务描述</h3>
            <p style={{ whiteSpace: 'pre-wrap', color: '#666' }}>{service.description || '暂无描述'}</p>

            {service.tags && (
              <>
                <Divider />
                <h3>标签</h3>
                <div>
                  {service.tags.split(',').map((tag, index) => (
                    <Tag key={index} color="blue">{tag.trim()}</Tag>
                  ))}
                </div>
              </>
            )}
          </Card>

          <Card title="用户评价" style={{ marginTop: 24 }}>
            {reviews?.items && reviews.items.length > 0 ? (
              <List
                dataSource={reviews.items}
                renderItem={(review) => (
                  <List.Item key={review.id}>
                    <List.Item.Meta
                      avatar={<Avatar>{review.client?.full_name?.charAt(0)}</Avatar>}
                      title={
                        <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                          <span>{review.client?.full_name}</span>
                          <Rate disabled value={review.rating} style={{ fontSize: 12 }} />
                        </div>
                      }
                      description={review.content}
                    />
                  </List.Item>
                )}
              />
            ) : (
              <div className="empty-state">暂无评价</div>
            )}
          </Card>
        </Col>

        <Col xs={24} lg={8}>
          <Card title="预约服务" style={{ position: 'sticky', top: 24 }}>
            <Form layout="vertical">
              <Form.Item label="选择日期">
                <DatePicker
                  style={{ width: '100%' }}
                  disabledDate={disabledDate}
                  value={selectedDate}
                  onChange={(date) => {
                    setSelectedDate(date)
                    setSelectedSchedule(null)
                  }}
                  placeholder="请选择预约日期"
                />
              </Form.Item>

              {selectedDate && (
                <Form.Item label="选择时间">
                  {schedulesLoading ? (
                    <div>加载中...</div>
                  ) : schedules && schedules.length > 0 ? (
                    <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                      {schedules.map((schedule) => (
                        <Button
                          key={schedule.id}
                          type={selectedSchedule?.id === schedule.id ? 'primary' : 'default'}
                          onClick={() => setSelectedSchedule(schedule)}
                          disabled={schedule.is_booked}
                        >
                          {schedule.start_time} - {schedule.end_time}
                        </Button>
                      ))}
                    </div>
                  ) : (
                    <div className="empty-state">该日期暂无可用时段</div>
                  )}
                </Form.Item>
              )}

              <Form.Item label="备注">
                <Input.TextArea
                  rows={3}
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  placeholder="请输入您的需求或问题（选填）"
                />
              </Form.Item>

              <Button
                type="primary"
                block
                size="large"
                disabled={!selectedSchedule}
                onClick={() => setBookingModalVisible(true)}
                loading={createAppointmentMutation.isPending}
              >
                <CalendarOutlined /> 立即预约
              </Button>
            </Form>
          </Card>
        </Col>
      </Row>

      <Modal
        title="确认预约"
        open={bookingModalVisible}
        onOk={handleBook}
        onCancel={() => setBookingModalVisible(false)}
        confirmLoading={createAppointmentMutation.isPending}
      >
        <p>服务：{service.title}</p>
        <p>日期：{selectedDate?.format('YYYY-MM-DD')}</p>
        <p>时间：{selectedSchedule?.start_time} - {selectedSchedule?.end_time}</p>
        <p>费用：¥{service.price}</p>
        <p style={{ color: '#ff4d4d' }}>请在30分钟内完成支付，否则预约将自动取消</p>
      </Modal>
    </div>
  )
}
