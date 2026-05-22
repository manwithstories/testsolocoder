import React, { useState, useEffect } from 'react'
import {
  Card,
  Form,
  Select,
  DatePicker,
  Button,
  Input,
  Spin,
  Alert,
  Space,
  Avatar,
  Tag,
  message,
  Row,
  Col,
  Descriptions,
  Divider
} from 'antd'
import {
  ArrowLeftOutlined,
  UserOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  ClockCircleOutlined
} from '@ant-design/icons'
import { useNavigate, useSearchParams } from 'react-router-dom'
import dayjs, { Dayjs } from 'dayjs'
import { departmentAPI, doctorAPI, appointmentAPI } from '@/services/api'
import type { Department, Doctor, TimeSlot } from '@/types'

const { Option } = Select
const { TextArea } = Input

const CreateAppointment: React.FC = () => {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [form] = Form.useForm()

  const doctorIdParam = searchParams.get('doctor_id')
  const appointmentIdParam = searchParams.get('appointment_id')

  const [departments, setDepartments] = useState<Department[]>([])
  const [doctors, setDoctors] = useState<Doctor[]>([])
  const [timeSlots, setTimeSlots] = useState<TimeSlot[]>([])
  const [selectedDoctor, setSelectedDoctor] = useState<Doctor | null>(null)

  const [loading, setLoading] = useState(false)
  const [doctorLoading, setDoctorLoading] = useState(false)
  const [slotLoading, setSlotLoading] = useState(false)
  const [submitLoading, setSubmitLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const [selectedDepartmentId, setSelectedDepartmentId] = useState<number | undefined>()
  const [selectedDoctorId, setSelectedDoctorId] = useState<number | undefined>()
  const [selectedDate, setSelectedDate] = useState<Dayjs | null>(null)
  const [selectedSlot, setSelectedSlot] = useState<TimeSlot | null>(null)

  useEffect(() => {
    fetchDepartments()
  }, [])

  useEffect(() => {
    if (selectedDepartmentId) {
      fetchDoctors(selectedDepartmentId)
    } else {
      setDoctors([])
      setSelectedDoctor(null)
      setSelectedDoctorId(undefined)
    }
  }, [selectedDepartmentId])

  useEffect(() => {
    if (selectedDoctorId && selectedDate) {
      fetchAvailableSlots(selectedDoctorId, selectedDate.format('YYYY-MM-DD'))
    } else {
      setTimeSlots([])
      setSelectedSlot(null)
    }
  }, [selectedDoctorId, selectedDate])

  useEffect(() => {
    if (doctorIdParam) {
      const id = parseInt(doctorIdParam)
      fetchDoctorDetail(id)
    }
  }, [doctorIdParam])

  const fetchDepartments = async () => {
    setLoading(true)
    try {
      const response = await departmentAPI.getList({ page: 1, pageSize: 100 })
      setDepartments(response.list)
    } catch (err: any) {
      setError(err.message || '获取科室列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchDoctors = async (departmentId: number) => {
    setDoctorLoading(true)
    try {
      const response = await doctorAPI.getList({
        page: 1,
        pageSize: 100,
        department_id: departmentId
      })
      setDoctors(response.list)
    } catch (err: any) {
      message.error(err.message || '获取医生列表失败')
    } finally {
      setDoctorLoading(false)
    }
  }

  const fetchDoctorDetail = async (doctorId: number) => {
    try {
      const response = await doctorAPI.getDetail(doctorId)
      setSelectedDoctor(response)
      setSelectedDoctorId(doctorId)
      if (response.department_id) {
        setSelectedDepartmentId(response.department_id)
      }
    } catch (err: any) {
      message.error(err.message || '获取医生详情失败')
    }
  }

  const fetchAvailableSlots = async (doctorId: number, date: string) => {
    setSlotLoading(true)
    try {
      const response = await doctorAPI.getAvailableSlots(doctorId, date)
      setTimeSlots(response)
    } catch (err: any) {
      message.error(err.message || '获取可用时段失败')
    } finally {
      setSlotLoading(false)
    }
  }

  const handleDepartmentChange = (value: number) => {
    setSelectedDepartmentId(value)
    setSelectedDoctorId(undefined)
    setSelectedDoctor(null)
    setSelectedDate(null)
    setSelectedSlot(null)
  }

  const handleDoctorChange = async (value: number) => {
    setSelectedDoctorId(value)
    setSelectedDate(null)
    setSelectedSlot(null)
    try {
      const response = await doctorAPI.getDetail(value)
      setSelectedDoctor(response)
    } catch (err) {
      console.error('获取医生详情失败:', err)
    }
  }

  const handleDateChange = (date: Dayjs | null) => {
    setSelectedDate(date)
    setSelectedSlot(null)
  }

  const handleSlotSelect = (slot: TimeSlot) => {
    if (slot.available) {
      setSelectedSlot(slot)
    }
  }

  const handleBack = () => {
    navigate('/appointments')
  }

  const handleSubmit = async () => {
    if (!selectedDoctorId || !selectedDate || !selectedSlot) {
      message.warning('请完成所有必填项')
      return
    }

    try {
      const values = await form.validateFields()
      setSubmitLoading(true)

      const appointmentData = {
        doctor_id: selectedDoctorId,
        appointment_date: selectedDate.format('YYYY-MM-DD'),
        start_time: selectedSlot.start,
        end_time: selectedSlot.end,
        symptoms: values.symptoms
      }

      if (appointmentIdParam) {
        await appointmentAPI.reschedule(parseInt(appointmentIdParam), appointmentData)
        message.success('预约改签成功')
      } else {
        await appointmentAPI.create(appointmentData)
        message.success('预约创建成功')
      }

      navigate('/appointments')
    } catch (err: any) {
      message.error(err.message || '提交预约失败')
    } finally {
      setSubmitLoading(false)
    }
  }

  const disabledDate = (current: Dayjs) => {
    return current && current < dayjs().startOf('day')
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={handleBack}
        className="mb-4"
        type="text"
      >
        返回预约列表
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

      <Card
        title={
          <div className="flex items-center gap-2">
            <CalendarOutlined />
            {appointmentIdParam ? '改签预约' : '新建预约'}
          </div>
        }
      >
        <Spin spinning={loading}>
          <Form form={form} layout="vertical">
            <Row gutter={24}>
              <Col xs={24} md={12}>
                <Form.Item
                  label="选择科室"
                  name="department_id"
                  rules={[{ required: true, message: '请选择科室' }]}
                >
                  <Select
                    placeholder="请选择科室"
                    size="large"
                    onChange={handleDepartmentChange}
                    value={selectedDepartmentId}
                    loading={loading}
                  >
                    {departments.map((dept) => (
                      <Option key={dept.id} value={dept.id}>
                        {dept.name}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
              </Col>
              <Col xs={24} md={12}>
                <Form.Item
                  label="选择医生"
                  name="doctor_id"
                  rules={[{ required: true, message: '请选择医生' }]}
                >
                  <Select
                    placeholder="请选择医生"
                    size="large"
                    onChange={handleDoctorChange}
                    value={selectedDoctorId}
                    loading={doctorLoading}
                    disabled={!selectedDepartmentId}
                    showSearch
                    optionFilterProp="children"
                  >
                    {doctors.map((doctor) => (
                      <Option key={doctor.id} value={doctor.id}>
                        {doctor.user?.full_name} - {doctor.title}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
              </Col>
            </Row>

            {selectedDoctor && (
              <Card size="small" className="mb-4 bg-gray-50">
                <Row gutter={16} align="middle">
                  <Col flex="80px">
                    <Avatar
                      size={64}
                      src={selectedDoctor.user?.avatar_url}
                      icon={!selectedDoctor.user?.avatar_url && <UserOutlined />}
                    />
                  </Col>
                  <Col flex="auto">
                    <Space size="middle" align="center">
                      <span className="text-lg font-semibold">
                        {selectedDoctor.user?.full_name}
                      </span>
                      <Tag color="blue">{selectedDoctor.title}</Tag>
                      <Tag color="green">
                        {selectedDoctor.department?.name}
                      </Tag>
                    </Space>
                    <p className="text-gray-500 mt-1 text-sm">
                      擅长：{selectedDoctor.specialty || '暂无信息'}
                    </p>
                    <p className="text-orange-500 font-medium mt-1">
                      挂号费：¥{selectedDoctor.registration_fee}
                    </p>
                  </Col>
                </Row>
              </Card>
            )}

            <Form.Item
              label="选择日期"
              name="appointment_date"
              rules={[{ required: true, message: '请选择预约日期' }]}
            >
              <DatePicker
                placeholder="请选择预约日期"
                size="large"
                style={{ width: '100%' }}
                onChange={handleDateChange}
                value={selectedDate}
                disabledDate={disabledDate}
                disabled={!selectedDoctorId}
              />
            </Form.Item>

            <Form.Item
              label="选择时间段"
              name="time_slot"
              rules={[{ required: true, message: '请选择时间段' }]}
            >
              <Spin spinning={slotLoading}>
                {selectedDate && selectedDoctorId ? (
                  timeSlots.length > 0 ? (
                    <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-2">
                      {timeSlots.map((slot, index) => (
                        <div
                          key={index}
                          onClick={() => handleSlotSelect(slot)}
                          className={`
                            p-3 rounded-lg border-2 text-center cursor-pointer transition-all
                            ${!slot.available
                              ? 'bg-gray-100 border-gray-200 text-gray-400 cursor-not-allowed'
                              : selectedSlot?.start === slot.start && selectedSlot?.end === slot.end
                                ? 'bg-blue-50 border-blue-500 text-blue-600'
                                : 'hover:border-blue-300 hover:bg-blue-50 border-gray-200'
                            }
                          `}
                        >
                          <div className="text-sm font-medium">
                            {slot.start}
                          </div>
                          <div className="text-xs mt-1">
                            {slot.available ? (
                              <CheckCircleOutlined className="text-green-500" />
                            ) : (
                              '已满'
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <Alert
                      message="该日期暂无可用时段"
                      type="warning"
                      showIcon
                    />
                  )
                ) : (
                  <div className="text-gray-400 text-center py-8 border-2 border-dashed border-gray-200 rounded-lg">
                    <ClockCircleOutlined className="text-3xl mb-2" />
                    <p>请先选择医生和日期</p>
                  </div>
                )}
              </Spin>
            </Form.Item>

            <Divider />

            <Form.Item label="症状描述" name="symptoms">
              <TextArea
                rows={4}
                placeholder="请简要描述您的症状（选填）"
                maxLength={500}
                showCount
              />
            </Form.Item>

            {selectedSlot && (
              <Card size="small" className="mb-4 bg-blue-50 border-blue-200">
                <Descriptions column={1} size="small">
                  <Descriptions.Item label="预约信息">
                    <Space>
                      <Tag color="blue">
                        {selectedDate?.format('YYYY-MM-DD')}
                      </Tag>
                      <Tag color="green">
                        {selectedSlot.start} - {selectedSlot.end}
                      </Tag>
                    </Space>
                  </Descriptions.Item>
                </Descriptions>
              </Card>
            )}

            <Form.Item>
              <Space size="large">
                <Button
                  type="primary"
                  size="large"
                  onClick={handleSubmit}
                  loading={submitLoading}
                  disabled={!selectedSlot}
                >
                  {appointmentIdParam ? '确认改签' : '提交预约'}
                </Button>
                <Button size="large" onClick={handleBack}>
                  取消
                </Button>
              </Space>
            </Form.Item>
          </Form>
        </Spin>
      </Card>
    </div>
  )
}

export default CreateAppointment
