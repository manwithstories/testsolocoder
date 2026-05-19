import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Button, DatePicker, List, Tag, Form, Input, Modal, message } from 'antd'
import { ArrowLeftOutlined, CalendarOutlined, BookOutlined } from '@ant-design/icons'
import { venueApi, orderApi } from '@/api'
import type { Venue } from '@/types'
import dayjs from 'dayjs'

const VenueDetail = () => {
  const { id } = useParams()
  const navigate = useNavigate()
  const [venue, setVenue] = useState<Venue | null>(null)
  const [availability, setAvailability] = useState<any>(null)
  const [selectedDate, setSelectedDate] = useState(dayjs().format('YYYY-MM-DD'))
  const [bookingModalVisible, setBookingModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (id) {
      loadVenue()
      loadAvailability(selectedDate)
    }
  }, [id, selectedDate])

  const loadVenue = async () => {
    try {
      const data = await venueApi.get(Number(id))
      setVenue(data)
    } catch (error: any) {
      message.error(error.message || '加载失败')
    }
  }

  const loadAvailability = async (date: string) => {
    setLoading(true)
    try {
      const data = await venueApi.getAvailability(Number(id), date)
      setAvailability(data)
    } catch (error: any) {
      message.error(error.message || '加载可用时段失败')
    } finally {
      setLoading(false)
    }
  }

  const handleBooking = (slot: any) => {
    form.setFieldsValue({
      start_time: `${selectedDate} ${slot.start}`,
      end_time: `${selectedDate} ${slot.end}`,
    })
    setBookingModalVisible(true)
  }

  const handleSubmitBooking = async (values: any) => {
    try {
      await orderApi.create({
        ...values,
        type: 'venue',
        item_id: Number(id),
      })
      message.success('预约提交成功')
      setBookingModalVisible(false)
      loadAvailability(selectedDate)
    } catch (error: any) {
      message.error(error.message || '预约失败')
    }
  }

  if (!venue) return null

  return (
    <div>
      <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/venues')} style={{ marginBottom: 16 }}>
        返回列表
      </Button>

      <Card title={venue.name}>
        {venue.cover_image && (
          <div style={{ marginBottom: 16 }}>
            <img src={venue.cover_image} alt={venue.name} style={{ width: '100%', maxHeight: 300, objectFit: 'cover' }} />
          </div>
        )}

        <Descriptions bordered column={2}>
          <Descriptions.Item label="位置">{venue.location}</Descriptions.Item>
          <Descriptions.Item label="容纳人数">{venue.capacity} 人</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color={venue.status === 'online' ? 'green' : 'red'}>
              {venue.status === 'online' ? '上架' : '下架'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="创建时间">{dayjs(venue.created_at).format('YYYY-MM-DD HH:mm')}</Descriptions.Item>
          <Descriptions.Item label="设施设备" span={2}>{venue.facilities || '-'}</Descriptions.Item>
          <Descriptions.Item label="描述" span={2}>{venue.description || '-'}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="可用时段查询" style={{ marginTop: 16 }}>
        <DatePicker
          value={dayjs(selectedDate)}
          onChange={(date) => date && setSelectedDate(date.format('YYYY-MM-DD'))}
          style={{ marginBottom: 16 }}
        />

        {availability?.available?.length > 0 && (
          <div style={{ marginBottom: 16 }}>
            <h4 style={{ marginBottom: 12 }}><CalendarOutlined /> 可预约时段</h4>
            <List
              grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4 }}
              dataSource={availability.available}
              renderItem={(slot: any) => (
                <List.Item>
                  <Card
                    hoverable
                    actions={[
                      <Button type="primary" icon={<BookOutlined />} onClick={() => handleBooking(slot)}>
                        预约
                      </Button>,
                    ]}
                  >
                    <Card.Meta
                      title={`${slot.start} - ${slot.end}`}
                      description={`¥${slot.price} / 小时`}
                    />
                  </Card>
                </List.Item>
              )}
            />
          </div>
        )}

        {availability?.booked?.length > 0 && (
          <div>
            <h4 style={{ marginBottom: 12 }}>已预约时段</h4>
            <List
              grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4 }}
              dataSource={availability.booked}
              renderItem={(slot: any) => (
                <List.Item>
                  <Card style={{ background: '#fff1f0' }}>
                    <Card.Meta
                      title={`${slot.start} - ${slot.end}`}
                      description={<Tag color="red">已被预约</Tag>}
                    />
                  </Card>
                </List.Item>
              )}
            />
          </div>
        )}
      </Card>

      <Modal
        title="提交预约"
        open={bookingModalVisible}
        onCancel={() => setBookingModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmitBooking}>
          <Form.Item
            name="start_time"
            label="开始时间"
            rules={[{ required: true, message: '请输入开始时间' }]}
          >
            <Input disabled />
          </Form.Item>
          <Form.Item
            name="end_time"
            label="结束时间"
            rules={[{ required: true, message: '请输入结束时间' }]}
          >
            <Input disabled />
          </Form.Item>
          <Form.Item
            name="purpose"
            label="用途说明"
            rules={[{ max: 500, message: '最多500个字符' }]}
          >
            <Input.TextArea rows={3} placeholder="请输入用途说明" />
          </Form.Item>
          <Form.Item
            name="contact_name"
            label="联系人姓名"
            rules={[{ required: true, message: '请输入联系人姓名' }]}
          >
            <Input placeholder="请输入联系人姓名" />
          </Form.Item>
          <Form.Item
            name="contact_phone"
            label="联系电话"
            rules={[{ required: true, message: '请输入联系电话' }]}
          >
            <Input placeholder="请输入联系电话" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              提交预约
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default VenueDetail
