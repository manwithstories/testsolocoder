import React, { useState, useEffect } from 'react'
import {
  Card,
  Table,
  Button,
  Space,
  Tag,
  Select,
  DatePicker,
  Spin,
  Alert,
  Modal,
  Input,
  Form,
  message,
  Empty,
  Popconfirm
} from 'antd'
import {
  PlusOutlined,
  EyeOutlined,
  EditOutlined,
  CloseOutlined,
  CalendarOutlined
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import dayjs, { Dayjs } from 'dayjs'
import { appointmentAPI } from '@/services/api'
import type { Appointment, AppointmentStatus } from '@/types'

const { RangePicker } = DatePicker
const { Option } = Select
const { TextArea } = Input

const statusMap: Record<AppointmentStatus, { text: string; color: string }> = {
  pending: { text: '待确认', color: 'orange' },
  confirmed: { text: '已确认', color: 'blue' },
  completed: { text: '已完成', color: 'green' },
  cancelled: { text: '已取消', color: 'red' },
  no_show: { text: '未就诊', color: 'default' }
}

const AppointmentList: React.FC = () => {
  const navigate = useNavigate()
  const [appointments, setAppointments] = useState<Appointment[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [status, setStatus] = useState<AppointmentStatus | undefined>()
  const [dateRange, setDateRange] = useState<[Dayjs | null, Dayjs | null] | null>(null)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [total, setTotal] = useState(0)

  const [cancelModalVisible, setCancelModalVisible] = useState(false)
  const [cancelAppointmentId, setCancelAppointmentId] = useState<number | null>(null)
  const [cancelForm] = Form.useForm()
  const [cancelLoading, setCancelLoading] = useState(false)

  useEffect(() => {
    fetchAppointments()
  }, [status, dateRange, page, pageSize])

  const fetchAppointments = async () => {
    setLoading(true)
    setError(null)
    try {
      const params: any = { page, pageSize }
      if (status) params.status = status
      if (dateRange && dateRange[0] && dateRange[1]) {
        params.start_date = dateRange[0].format('YYYY-MM-DD')
        params.end_date = dateRange[1].format('YYYY-MM-DD')
      }
      const response = await appointmentAPI.getList(params)
      setAppointments(response.list)
      setTotal(response.total)
    } catch (err: any) {
      setError(err.message || '获取预约列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleStatusChange = (value: AppointmentStatus | undefined) => {
    setStatus(value)
    setPage(1)
  }

  const handleDateRangeChange = (dates: [Dayjs | null, Dayjs | null] | null) => {
    setDateRange(dates)
    setPage(1)
  }

  const handleViewDetail = (id: number) => {
    navigate(`/appointments/${id}`)
  }

  const handleCreate = () => {
    navigate('/appointments/create')
  }

  const handleCancel = (id: number) => {
    setCancelAppointmentId(id)
    setCancelModalVisible(true)
    cancelForm.resetFields()
  }

  const handleCancelSubmit = async () => {
    if (!cancelAppointmentId) return
    setCancelLoading(true)
    try {
      const values = await cancelForm.validateFields()
      await appointmentAPI.cancel(cancelAppointmentId, {
        cancel_reason: values.cancel_reason
      })
      message.success('预约已取消')
      setCancelModalVisible(false)
      fetchAppointments()
    } catch (err: any) {
      message.error(err.message || '取消预约失败')
    } finally {
      setCancelLoading(false)
    }
  }

  const handleReschedule = (id: number) => {
    navigate(`/appointments/create?appointment_id=${id}`)
  }

  const columns = [
    {
      title: '预约编号',
      dataIndex: 'id',
      key: 'id',
      width: 100
    },
    {
      title: '医生',
      dataIndex: 'doctor',
      key: 'doctor',
      render: (doctor: Appointment['doctor']) =>
        doctor ? doctor.user?.full_name || '-' : '-'
    },
    {
      title: '科室',
      dataIndex: 'doctor',
      key: 'department',
      render: (doctor: Appointment['doctor']) =>
        doctor ? doctor.department?.name || '-' : '-'
    },
    {
      title: '预约日期',
      dataIndex: 'appointment_date',
      key: 'appointment_date',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD')
    },
    {
      title: '时间段',
      key: 'time',
      render: (_: unknown, record: Appointment) =>
        `${record.start_time} - ${record.end_time}`
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: AppointmentStatus) => {
        const info = statusMap[status]
        return <Tag color={info.color}>{info.text}</Tag>
      }
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')
    },
    {
      title: '操作',
      key: 'actions',
      width: 200,
      render: (_: unknown, record: Appointment) => (
        <Space size="small">
          <Button
            type="link"
            size="small"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record.id)}
          >
            详情
          </Button>
          {(record.status === 'pending' || record.status === 'confirmed') && (
            <Button
              type="link"
              size="small"
              icon={<EditOutlined />}
              onClick={() => handleReschedule(record.id)}
            >
              改签
            </Button>
          )}
          {(record.status === 'pending' || record.status === 'confirmed') && (
            <Popconfirm
              title="确定要取消这个预约吗？"
              description="取消后将无法恢复"
              onConfirm={() => handleCancel(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="link" size="small" danger icon={<CloseOutlined />}>
                取消
              </Button>
            </Popconfirm>
          )}
        </Space>
      )
    }
  ]

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-800 mb-2">预约列表</h1>
          <p className="text-gray-500">查看和管理您的预约记录</p>
        </div>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          size="large"
          onClick={handleCreate}
        >
          新建预约
        </Button>
      </div>

      <Card className="mb-6">
        <Space size="large" wrap>
          <div className="flex items-center gap-2">
            <span className="text-gray-600">状态：</span>
            <Select
              placeholder="全部状态"
              allowClear
              style={{ width: 150 }}
              value={status}
              onChange={handleStatusChange}
            >
              <Option value="pending">待确认</Option>
              <Option value="confirmed">已确认</Option>
              <Option value="completed">已完成</Option>
              <Option value="cancelled">已取消</Option>
              <Option value="no_show">未就诊</Option>
            </Select>
          </div>
          <div className="flex items-center gap-2">
            <span className="text-gray-600">日期范围：</span>
            <RangePicker
              value={dateRange}
              onChange={handleDateRangeChange}
              placeholder={['开始日期', '结束日期']}
            />
          </div>
          <Button onClick={fetchAppointments} icon={<CalendarOutlined />}>
            刷新
          </Button>
        </Space>
      </Card>

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
        {appointments.length > 0 ? (
          <Card>
            <Table
              columns={columns}
              dataSource={appointments}
              rowKey="id"
              pagination={{
                current: page,
                pageSize,
                total,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total) => `共 ${total} 条记录`,
                onChange: (page, pageSize) => {
                  setPage(page)
                  setPageSize(pageSize)
                }
              }}
            />
          </Card>
        ) : (
          !loading && (
            <Empty
              description="暂无预约记录"
              image={Empty.PRESENTED_IMAGE_SIMPLE}
              className="py-16"
            />
          )
        )}
      </Spin>

      <Modal
        title="取消预约"
        open={cancelModalVisible}
        onCancel={() => setCancelModalVisible(false)}
        footer={[
          <Button key="back" onClick={() => setCancelModalVisible(false)}>
            返回
          </Button>,
          <Button
            key="submit"
            type="primary"
            danger
            loading={cancelLoading}
            onClick={handleCancelSubmit}
          >
            确认取消
          </Button>
        ]}
      >
        <Form form={cancelForm} layout="vertical">
          <Form.Item
            name="cancel_reason"
            label="取消原因"
            rules={[{ required: true, message: '请输入取消原因' }]}
          >
            <TextArea rows={4} placeholder="请输入取消预约的原因" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AppointmentList
