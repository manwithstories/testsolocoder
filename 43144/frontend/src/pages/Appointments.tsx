import React, { useEffect, useState } from 'react'
import { Table, Tag, Card, Button, Modal, Form, Input, Select, DatePicker, TimePicker, message, Space, Empty, Spin, Upload } from 'antd'
import { PlusOutlined, UploadOutlined, CalendarOutlined, DeleteOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { listHealthRecords, createHealthRecord, deleteHealthRecord, uploadHealthReport } from '../api/health'
import { listAppointments, createAppointment, cancelAppointment, confirmAppointment, completeAppointment } from '../api/appointment'
import { HealthRecord, Appointment, CreateHealthRecordRequest, CreateAppointmentRequest } from '../types'
import dayjs from 'dayjs'

const { Option } = Select
const { TextArea } = Input

const Appointments: React.FC = () => {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    loadAppointments()
  }, [])

  const loadAppointments = async () => {
    setLoading(true)
    try {
      const response = await listAppointments({ page_size: 100 })
      if (response.code === 0 && response.data) {
        setAppointments((response.data as any).items || [])
      }
    } catch (error) {
      console.error('Failed to load appointments:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      setSubmitting(true)

      const data: CreateAppointmentRequest = {
        pet_id: values.pet_id,
        appointment_type: values.appointment_type,
        appointment_date: values.appointment_date.format('YYYY-MM-DD'),
        start_time: values.start_time.format('HH:mm'),
        end_time: values.end_time.format('HH:mm'),
        notes: values.notes,
      }

      await createAppointment(data)
      message.success('预约成功')
      setModalVisible(false)
      form.resetFields()
      loadAppointments()
    } catch (error: any) {
      message.error(error.message || '预约失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleCancel = async (id: number) => {
    Modal.confirm({
      title: '确认取消预约',
      content: '确定要取消这个预约吗？',
      onOk: async () => {
        try {
          await cancelAppointment(id)
          message.success('已取消')
          loadAppointments()
        } catch (error: any) {
          message.error('取消失败')
        }
      },
    })
  }

  const handleConfirm = async (id: number) => {
    try {
      await confirmAppointment(id)
      message.success('已确认')
      loadAppointments()
    } catch (error: any) {
      message.error('确认失败')
    }
  }

  const handleComplete = async (id: number) => {
    try {
      await completeAppointment(id)
      message.success('已完成')
      loadAppointments()
    } catch (error: any) {
      message.error('操作失败')
    }
  }

  const columns = [
    { title: '预约编号', dataIndex: 'id', key: 'id' },
    {
      title: '宠物',
      key: 'pet',
      render: (_: any, r: Appointment) => r.pet?.name || '-',
    },
    {
      title: '预约人',
      key: 'user',
      render: (_: any, r: Appointment) => r.user?.name || '-',
    },
    { title: '类型', dataIndex: 'appointment_type', key: 'type', render: (t: string) => t === 'visit' ? '上门看望' : '体检' },
    {
      title: '预约日期',
      dataIndex: 'appointment_date',
      key: 'date',
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    { title: '时间段', key: 'time', render: (_: any, r: Appointment) => `${r.start_time} - ${r.end_time}` },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          pending: 'orange',
          confirmed: 'blue',
          cancelled: 'default',
          completed: 'green',
          rescheduled: 'cyan',
        }
        return <Tag color={colorMap[status]}>{status}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Appointment) => (
        <Space>
          {record.status === 'pending' && user?.role === 'rescue' && (
            <Button type="link" size="small" onClick={() => handleConfirm(record.id)}>
              确认
            </Button>
          )}
          {record.status === 'confirmed' && user?.role === 'rescue' && (
            <Button type="link" size="small" onClick={() => handleComplete(record.id)}>
              完成
            </Button>
          )}
          {(record.status === 'pending' || record.status === 'confirmed') && (
            <Button type="link" size="small" danger onClick={() => handleCancel(record.id)}>
              取消
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>预约管理</h2>
        {user?.role === 'adopter' && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
            新建预约
          </Button>
        )}
      </div>

      <Spin spinning={loading}>
        {appointments.length > 0 ? (
          <Card>
            <Table
              dataSource={appointments}
              columns={columns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
          </Card>
        ) : (
          !loading && (
            <Card>
              <Empty description="暂无预约记录" />
            </Card>
          )
        )}
      </Spin>

      <Modal
        title="新建预约"
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false)
          form.resetFields()
        }}
        footer={null}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="pet_id" label="选择宠物" rules={[{ required: true }]}>
            <Select
              placeholder="请先在宠物列表选择"
              showSearch
              optionFilterProp="children"
            >
              <Option value="">请先前往宠物详情页预约</Option>
            </Select>
          </Form.Item>
          <Form.Item name="appointment_type" label="预约类型" rules={[{ required: true }]}>
            <Select>
              <Option value="visit">上门看望</Option>
              <Option value="checkup">体检</Option>
            </Select>
          </Form.Item>
          <Form.Item name="appointment_date" label="预约日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Space style={{ width: '100%' }}>
            <Form.Item name="start_time" label="开始时间" rules={[{ required: true }]} style={{ flex: 1 }}>
              <TimePicker format="HH:mm" style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="end_time" label="结束时间" rules={[{ required: true }]} style={{ flex: 1 }}>
              <TimePicker format="HH:mm" style={{ width: '100%' }} />
            </Form.Item>
          </Space>
          <Form.Item name="notes" label="备注">
            <TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting} block>
              提交预约
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default Appointments
