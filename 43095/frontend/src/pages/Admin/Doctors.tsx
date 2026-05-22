import React, { useState, useEffect } from 'react'
import {
  Card,
  Table,
  Button,
  Space,
  Typography,
  Tag,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  message,
  Spin,
  Empty,
  Pagination,
  Popconfirm,
  Row,
  Col,
  Descriptions,
  List,
  Switch
} from 'antd'
import {
  UserOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  CloseOutlined,
  ScheduleOutlined,
  EyeOutlined
} from '@ant-design/icons'
import type { ColumnsType } from 'antd/es/table'
import type { Doctor, Schedule, Department, DoctorTitle, DayOfWeek } from '@/types'
import { doctorAPI, departmentAPI } from '@/services/api'

const { Title, Text } = Typography
const { TextArea } = Input
const { Option } = Select

const dayMap: Record<number, string> = {
  0: '周日',
  1: '周一',
  2: '周二',
  3: '周三',
  4: '周四',
  5: '周五',
  6: '周六'
}

const titleOptions = ['住院医师', '主治医师', '副主任医师', '主任医师', '教授']

interface DoctorFormData {
  user_id: number
  department_id: number
  title: DoctorTitle
  specialty: string
  introduction: string
  consultation_fee: number
  registration_fee: number
}

interface ScheduleFormData {
  day_of_week: DayOfWeek
  start_time: string
  end_time: string
  max_patients: number
  time_slot_minutes: number
  is_available: boolean
}

const AdminDoctorsPage: React.FC = () => {
  const [form] = Form.useForm<DoctorFormData>()
  const [scheduleForm] = Form.useForm<ScheduleFormData>()
  const [loading, setLoading] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [doctors, setDoctors] = useState<Doctor[]>([])
  const [departments, setDepartments] = useState<Department[]>([])
  const [schedules, setSchedules] = useState<Schedule[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize] = useState(10)
  const [modalVisible, setModalVisible] = useState(false)
  const [scheduleModalVisible, setScheduleModalVisible] = useState(false)
  const [editingDoctor, setEditingDoctor] = useState<Doctor | null>(null)
  const [selectedDoctor, setSelectedDoctor] = useState<Doctor | null>(null)
  const [editingSchedule, setEditingSchedule] = useState<Schedule | null>(null)

  useEffect(() => {
    fetchDoctors()
    fetchDepartments()
  }, [page, pageSize])

  const fetchDoctors = async () => {
    setLoading(true)
    try {
      const result = await doctorAPI.getList({ page, pageSize })
      setDoctors(result.list)
      setTotal(result.total)
    } catch (error) {
      console.error('获取医生列表失败:', error)
      message.error('获取医生列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchDepartments = async () => {
    try {
      const result = await departmentAPI.getList({ page: 1, pageSize: 100 })
      setDepartments(result.list)
    } catch (error) {
      console.error('获取科室列表失败:', error)
    }
  }

  const fetchSchedules = async (doctorId: number) => {
    try {
      const result = await doctorAPI.getSchedules(doctorId)
      setSchedules(result)
    } catch (error) {
      console.error('获取排班信息失败:', error)
      message.error('获取排班信息失败')
    }
  }

  const handleAdd = () => {
    setEditingDoctor(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (doctor: Doctor) => {
    setEditingDoctor(doctor)
    form.setFieldsValue({
      user_id: doctor.user_id,
      department_id: doctor.department_id,
      title: doctor.title,
      specialty: doctor.specialty,
      introduction: doctor.introduction,
      consultation_fee: doctor.consultation_fee,
      registration_fee: doctor.registration_fee
    })
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await doctorAPI.delete(id)
      message.success('删除成功')
      fetchDoctors()
    } catch (error) {
      console.error('删除失败:', error)
      message.error('删除失败')
    }
  }

  const handleSubmit = async (values: DoctorFormData) => {
    setSubmitting(true)
    try {
      if (editingDoctor) {
        await doctorAPI.update(editingDoctor.id, values)
        message.success('更新成功')
      } else {
        await doctorAPI.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchDoctors()
    } catch (error) {
      console.error('提交失败:', error)
      message.error('提交失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handleManageSchedule = (doctor: Doctor) => {
    setSelectedDoctor(doctor)
    fetchSchedules(doctor.id)
    setScheduleModalVisible(true)
  }

  const handleEditSchedule = (schedule: Schedule) => {
    setEditingSchedule(schedule)
    scheduleForm.setFieldsValue({
      day_of_week: schedule.day_of_week,
      start_time: schedule.start_time,
      end_time: schedule.end_time,
      max_patients: schedule.max_patients,
      time_slot_minutes: schedule.time_slot_minutes,
      is_available: schedule.is_available
    })
  }

  const handleDeleteSchedule = async (scheduleId: number) => {
    try {
      await doctorAPI.deleteSchedule(scheduleId)
      message.success('删除成功')
      if (selectedDoctor) {
        fetchSchedules(selectedDoctor.id)
      }
    } catch (error) {
      console.error('删除排班失败:', error)
      message.error('删除排班失败')
    }
  }

  const handleScheduleSubmit = async (values: ScheduleFormData) => {
    if (!selectedDoctor) return

    setSubmitting(true)
    try {
      if (editingSchedule) {
        await doctorAPI.updateSchedule(editingSchedule.id, values)
        message.success('更新成功')
      } else {
        await doctorAPI.createSchedule(selectedDoctor.id, values)
        message.success('创建成功')
      }
      setEditingSchedule(null)
      scheduleForm.resetFields()
      fetchSchedules(selectedDoctor.id)
    } catch (error) {
      console.error('提交失败:', error)
      message.error('提交失败')
    } finally {
      setSubmitting(false)
    }
  }

  const handlePageChange = (newPage: number) => {
    setPage(newPage)
  }

  const columns: ColumnsType<Doctor> = [
    {
      title: '姓名',
      key: 'name',
      width: 120,
      render: (_, record) => record.user?.full_name || '-'
    },
    {
      title: '科室',
      key: 'department',
      width: 120,
      render: (_, record) => record.department?.name || '-'
    },
    {
      title: '职称',
      dataIndex: 'title',
      key: 'title',
      width: 120
    },
    {
      title: '专长',
      dataIndex: 'specialty',
      key: 'specialty',
      ellipsis: true
    },
    {
      title: '挂号费',
      dataIndex: 'registration_fee',
      key: 'registration_fee',
      width: 100,
      render: (fee: number) => `¥${fee.toFixed(2)}`
    },
    {
      title: '诊疗费',
      dataIndex: 'consultation_fee',
      key: 'consultation_fee',
      width: 100,
      render: (fee: number) => `¥${fee.toFixed(2)}`
    },
    {
      title: '评分',
      key: 'rating',
      width: 100,
      render: (_, record) => (
        <Space>
          <Text type="warning">★</Text>
          <Text>{record.average_rating?.toFixed(1) || '0.0'}</Text>
          <Text type="secondary">({record.review_count})</Text>
        </Space>
      )
    },
    {
      title: '操作',
      key: 'action',
      width: 240,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            size="small"
            icon={<ScheduleOutlined />}
            onClick={() => handleManageSchedule(record)}
          >
            排班
          </Button>
          <Popconfirm
            title="确定要删除这名医生吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ]

  return (
    <div className="space-y-6">
      <Card>
        <div className="flex items-center justify-between mb-4">
          <Title level={3} style={{ margin: 0 }}>
            <UserOutlined className="mr-2" />
            医生管理
          </Title>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            添加医生
          </Button>
        </div>

        {loading && doctors.length === 0 ? (
          <div className="flex justify-center items-center min-h-[400px]">
            <Spin size="large" />
          </div>
        ) : doctors.length === 0 ? (
          <Empty description="暂无医生数据" />
        ) : (
          <>
            <Table
              columns={columns}
              dataSource={doctors}
              rowKey="id"
              pagination={false}
              loading={loading}
              scroll={{ x: 1200 }}
            />
            <div className="flex justify-end mt-4">
              <Pagination
                current={page}
                pageSize={pageSize}
                total={total}
                onChange={handlePageChange}
                showSizeChanger={false}
                showQuickJumper
                showTotal={(total) => `共 ${total} 条记录`}
              />
            </div>
          </>
        )}
      </Card>

      <Modal
        title={editingDoctor ? '编辑医生' : '添加医生'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={[
          <Button key="cancel" onClick={() => setModalVisible(false)}>
            取消
          </Button>,
          <Button
            key="submit"
            type="primary"
            loading={submitting}
            onClick={() => form.submit()}
          >
            确定
          </Button>
        ]}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <Form.Item
                name="user_id"
                label="关联用户ID"
                rules={[{ required: true, message: '请输入关联用户ID' }]}
              >
                <InputNumber min={1} style={{ width: '100%' }} placeholder="请输入用户ID" />
              </Form.Item>
            </Col>
            <Col xs={24} sm={12}>
              <Form.Item
                name="department_id"
                label="所属科室"
                rules={[{ required: true, message: '请选择科室' }]}
              >
                <Select placeholder="请选择科室">
                  {departments.map((dept) => (
                    <Option key={dept.id} value={dept.id}>
                      {dept.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <Form.Item
                name="title"
                label="职称"
                rules={[{ required: true, message: '请选择职称' }]}
              >
                <Select placeholder="请选择职称">
                  {titleOptions.map((title) => (
                    <Option key={title} value={title}>
                      {title}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col xs={24} sm={12}>
              <Form.Item
                name="specialty"
                label="专长"
                rules={[{ required: true, message: '请输入专长' }]}
              >
                <Input placeholder="请输入专长" />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <Form.Item
                name="registration_fee"
                label="挂号费"
                rules={[{ required: true, message: '请输入挂号费' }]}
              >
                <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="请输入挂号费" />
              </Form.Item>
            </Col>
            <Col xs={24} sm={12}>
              <Form.Item
                name="consultation_fee"
                label="诊疗费"
                rules={[{ required: true, message: '请输入诊疗费' }]}
              >
                <InputNumber min={0} precision={2} style={{ width: '100%' }} placeholder="请输入诊疗费" />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item
            name="introduction"
            label="个人简介"
            rules={[{ required: true, message: '请输入个人简介' }]}
          >
            <TextArea rows={4} placeholder="请输入个人简介" />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={`排班管理 - ${selectedDoctor?.user?.full_name || ''}`}
        open={scheduleModalVisible}
        onCancel={() => {
          setScheduleModalVisible(false)
          setEditingSchedule(null)
          scheduleForm.resetFields()
        }}
        footer={[
          <Button
            key="close"
            icon={<CloseOutlined />}
            onClick={() => {
              setScheduleModalVisible(false)
              setEditingSchedule(null)
              scheduleForm.resetFields()
            }}
          >
            关闭
          </Button>
        ]}
        width={800}
      >
        {selectedDoctor && (
          <div className="space-y-4">
            <Descriptions bordered column={2} size="small">
              <Descriptions.Item label="姓名">
                {selectedDoctor.user?.full_name || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="科室">
                {selectedDoctor.department?.name || '-'}
              </Descriptions.Item>
              <Descriptions.Item label="职称">
                {selectedDoctor.title}
              </Descriptions.Item>
              <Descriptions.Item label="专长">
                {selectedDoctor.specialty}
              </Descriptions.Item>
            </Descriptions>

            <Card size="small" title="添加/编辑排班">
              <Form
                form={scheduleForm}
                layout="vertical"
                onFinish={handleScheduleSubmit}
              >
                <Row gutter={16}>
                  <Col xs={24} sm={8}>
                    <Form.Item
                      name="day_of_week"
                      label="星期"
                      rules={[{ required: true, message: '请选择星期' }]}
                    >
                      <Select placeholder="请选择星期">
                        {Object.entries(dayMap).map(([key, value]) => (
                          <Option key={key} value={Number(key)}>
                            {value}
                          </Option>
                        ))}
                      </Select>
                    </Form.Item>
                  </Col>
                  <Col xs={24} sm={8}>
                    <Form.Item
                      name="start_time"
                      label="开始时间"
                      rules={[{ required: true, message: '请选择开始时间' }]}
                    >
                      <Input type="time" placeholder="HH:mm" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} sm={8}>
                    <Form.Item
                      name="end_time"
                      label="结束时间"
                      rules={[{ required: true, message: '请选择结束时间' }]}
                    >
                      <Input type="time" placeholder="HH:mm" />
                    </Form.Item>
                  </Col>
                </Row>
                <Row gutter={16}>
                  <Col xs={24} sm={8}>
                    <Form.Item
                      name="max_patients"
                      label="最大接诊数"
                      rules={[{ required: true, message: '请输入最大接诊数' }]}
                    >
                      <InputNumber min={1} style={{ width: '100%' }} placeholder="请输入最大接诊数" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} sm={8}>
                    <Form.Item
                      name="time_slot_minutes"
                      label="每号时长(分钟)"
                      rules={[{ required: true, message: '请输入每号时长' }]}
                    >
                      <InputNumber min={5} step={5} style={{ width: '100%' }} placeholder="请输入每号时长" />
                    </Form.Item>
                  </Col>
                  <Col xs={24} sm={8}>
                    <Form.Item
                      name="is_available"
                      label="是否可用"
                      valuePropName="checked"
                    >
                      <Switch />
                    </Form.Item>
                  </Col>
                </Row>
                <div className="flex justify-end">
                  <Space>
                    {editingSchedule && (
                      <Button onClick={() => {
                        setEditingSchedule(null)
                        scheduleForm.resetFields()
                      }}>
                        取消
                      </Button>
                    )}
                    <Button type="primary" htmlType="submit" loading={submitting}>
                      {editingSchedule ? '更新排班' : '添加排班'}
                    </Button>
                  </Space>
                </div>
              </Form>
            </Card>

            <Card size="small" title="排班列表">
              {schedules.length === 0 ? (
                <Empty description="暂无排班" />
              ) : (
                <List
                  dataSource={schedules}
                  renderItem={(schedule) => (
                    <List.Item
                      key={schedule.id}
                      actions={[
                        <Button
                          type="link"
                          size="small"
                          icon={<EditOutlined />}
                          onClick={() => handleEditSchedule(schedule)}
                        >
                          编辑
                        </Button>,
                        <Popconfirm
                          title="确定要删除这个排班吗？"
                          onConfirm={() => handleDeleteSchedule(schedule.id)}
                          okText="确定"
                          cancelText="取消"
                        >
                          <Button type="link" size="small" danger icon={<DeleteOutlined />}>
                            删除
                          </Button>
                        </Popconfirm>
                      ]}
                    >
                      <List.Item.Meta
                        title={
                          <Space>
                            <Text strong>{dayMap[schedule.day_of_week]}</Text>
                            <Tag color={schedule.is_available ? 'green' : 'default'}>
                              {schedule.is_available ? '可用' : '不可用'}
                            </Tag>
                          </Space>
                        }
                        description={
                          <Space direction="vertical" size="small" className="mt-2">
                            <Text>时间: {schedule.start_time} - {schedule.end_time}</Text>
                            <Text>最大接诊: {schedule.max_patients}人</Text>
                            <Text>每号时长: {schedule.time_slot_minutes}分钟</Text>
                          </Space>
                        }
                      />
                    </List.Item>
                  )}
                />
              )}
            </Card>
          </div>
        )}
      </Modal>
    </div>
  )
}

export default AdminDoctorsPage
