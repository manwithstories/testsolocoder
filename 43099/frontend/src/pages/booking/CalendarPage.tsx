import { useState, useEffect, useRef } from 'react'
import { Card, Select, Button, Modal, Form, Input, message, DatePicker, Space, Tabs } from 'antd'
import FullCalendar from '@fullcalendar/react'
import dayGridPlugin from '@fullcalendar/daygrid'
import timeGridPlugin from '@fullcalendar/timegrid'
import interactionPlugin from '@fullcalendar/interaction'
import { BookOutlined } from '@ant-design/icons'
import { orderApi, venueApi, deviceApi } from '@/api'
import type { CalendarEvent } from '@/types'
import dayjs from 'dayjs'

const CalendarPage = () => {
  const calendarRef = useRef<FullCalendar>(null)
  const [type, setType] = useState<'venue' | 'device'>('venue')
  const [itemId, setItemId] = useState<number | undefined>()
  const [venues, setVenues] = useState<any[]>([])
  const [devices, setDevices] = useState<any[]>([])
  const [events, setEvents] = useState<CalendarEvent[]>([])
  const [bookingModalVisible, setBookingModalVisible] = useState(false)
  const [selectedSlot, setSelectedSlot] = useState<any>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    loadItems()
  }, [type])

  useEffect(() => {
    if (itemId) {
      loadCalendarEvents()
    }
  }, [type, itemId])

  const loadItems = async () => {
    try {
      if (type === 'venue') {
        const data: any = await venueApi.list({ page_size: 100 })
        setVenues(data.list || [])
        if (data.list?.length > 0) {
          setItemId(data.list[0].id)
        }
      } else {
          const data: any = await deviceApi.list({ page_size: 100 })
          setDevices(data.list || [])
          if (data.list?.length > 0) {
            setItemId(data.list[0].id)
          }
        }
    } catch (error: any) {
      message.error(error.message || '加载失败')
    }
  }

  const loadCalendarEvents = async () => {
    if (!itemId) return
    try {
      const start = dayjs().startOf('month').format('YYYY-MM-DD')
      const end = dayjs().endOf('month').format('YYYY-MM-DD')
      const data = await orderApi.getCalendar({
        start_date: start,
        end_date: end,
        type,
        item_id: itemId,
      })
      setEvents(data)
    } catch (error: any) {
      message.error(error.message || '加载日历数据失败')
    }
  }

  const handleDateClick = (info: any) => {
    const start = dayjs(info.dateStr)
    const end = start.add(1, 'hour')
    setSelectedSlot({
      start: start.format('YYYY-MM-DD HH:mm'),
      end: end.format('YYYY-MM-DD HH:mm'),
    })
    form.setFieldsValue({
      start_time: start,
      end_time: end,
    })
    setBookingModalVisible(true)
  }

  const handleEventClick = (info: any) => {
    message.info(`订单详情: ${info.event.title}`)
  }

  const handleSubmitBooking = async (values: any) => {
    if (!itemId) return
    try {
      const submitData = {
        ...values,
        start_time: values.start_time?.format ? values.start_time.format('YYYY-MM-DD HH:mm') : values.start_time,
        end_time: values.end_time?.format ? values.end_time.format('YYYY-MM-DD HH:mm') : values.end_time,
        type,
        item_id: itemId,
      }
      await orderApi.create(submitData)
      message.success('预约提交成功')
      setBookingModalVisible(false)
      loadCalendarEvents()
    } catch (error: any) {
      message.error(error.message || '预约失败')
    }
  }

  const calendarEvents = events.map((event) => ({
    id: String(event.id),
    title: event.title,
    start: event.start,
    end: event.end,
    backgroundColor: event.color,
    borderColor: event.color,
  }))

  return (
    <Card
      title="预约日历"
      extra={
        <Space>
          <Tabs
            activeKey={type}
            onChange={(key) => setType(key as any)}
            items={[
              { key: 'venue', label: '场地' },
              { key: 'device', label: '设备' },
            ]}
          />
          <Select
            style={{ width: 200 }}
            value={itemId}
            onChange={setItemId}
            placeholder="选择场地/设备"
          >
            {(type === 'venue' ? venues : devices).map((item) => (
              <Select.Option key={item.id} value={item.id}>
              {item.name}
            </Select.Option>
          ))}
          </Select>
        </Space>
      }
    >
      <FullCalendar
        ref={calendarRef}
        plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin]}
        initialView="timeGridWeek"
        headerToolbar={{
          left: 'prev,next today',
          center: 'title',
          right: 'dayGridMonth,timeGridWeek,timeGridDay',
        }}
        events={calendarEvents}
        dateClick={handleDateClick}
        eventClick={handleEventClick}
        locale="zh-cn"
        height="600px"
        slotMinTime="08:00:00"
        slotMaxTime="22:00:00"
        allDaySlot={false}
      />

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
            rules={[{ required: true, message: '请选择开始时间' }]
          >
            <DatePicker showTime style={{ width: '100%' }} format="YYYY-MM-DD HH:mm" />
          </Form.Item>
          <Form.Item
            name="end_time"
            label="结束时间"
            rules={[{ required: true, message: '请选择结束时间' }]
          >
            <DatePicker showTime style={{ width: '100%' }} format="YYYY-MM-DD HH:mm" />
          </Form.Item>
          {type === 'device' && (
            <Form.Item
              name="quantity"
              label="租赁数量"
              rules={[{ required: true, message: '请输入数量' }]}
            >
              <Input type="number" min={1} defaultValue={1} />
            </Form.Item>
          )}
          <Form.Item
            name="purpose"
            label="用途说明"
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
            <Button type="primary" htmlType="submit" block icon={<BookOutlined />}>
              提交预约
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default CalendarPage
