import { Card, Row, Col, Statistic, Table, Tag, Button } from 'antd'
import {
  CalendarOutlined,
  ClockCircleOutlined,
  DollarOutlined,
  UserOutlined,
  RightOutlined
} from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { ColumnsType } from 'antd/es/table'
import type { Appointment } from '@/types'
import dayjs from 'dayjs'

const mockAppointments: Appointment[] = [
  {
    id: 1,
    patient_id: 1,
    doctor_id: 1,
    appointment_date: '2024-01-15',
    start_time: '09:00:00',
    end_time: '09:30:00',
    status: 'confirmed',
    symptoms: '头痛、发热',
    notes: '',
    created_at: '2024-01-10T10:00:00Z',
    updated_at: '2024-01-10T10:00:00Z',
    doctor: {
      id: 1,
      user_id: 2,
      department_id: 1,
      title: '主任医师',
      specialty: '神经内科',
      introduction: '',
      consultation_fee: 100,
      registration_fee: 50,
      average_rating: 4.8,
      review_count: 120,
      created_at: '',
      updated_at: '',
      user: {
        id: 2,
        username: 'doctor1',
        email: 'doctor1@example.com',
        phone: '13800138001',
        role: 'doctor',
        full_name: '张医生',
        gender: '男',
        birth_date: null,
        avatar_url: '',
        is_active: true,
        created_at: '',
        updated_at: ''
      }
    }
  },
  {
    id: 2,
    patient_id: 1,
    doctor_id: 2,
    appointment_date: '2024-01-15',
    start_time: '14:00:00',
    end_time: '14:30:00',
    status: 'pending',
    symptoms: '胃部不适',
    notes: '',
    created_at: '2024-01-11T14:00:00Z',
    updated_at: '2024-01-11T14:00:00Z',
    doctor: {
      id: 2,
      user_id: 3,
      department_id: 2,
      title: '副主任医师',
      specialty: '消化内科',
      introduction: '',
      consultation_fee: 80,
      registration_fee: 50,
      average_rating: 4.6,
      review_count: 89,
      created_at: '',
      updated_at: '',
      user: {
        id: 3,
        username: 'doctor2',
        email: 'doctor2@example.com',
        phone: '13800138002',
        role: 'doctor',
        full_name: '李医生',
        gender: '女',
        birth_date: null,
        avatar_url: '',
        is_active: true,
        created_at: '',
        updated_at: ''
      }
    }
  },
  {
    id: 3,
    patient_id: 1,
    doctor_id: 3,
    appointment_date: '2024-01-12',
    start_time: '10:00:00',
    end_time: '10:30:00',
    status: 'completed',
    symptoms: '咳嗽、胸闷',
    notes: '',
    created_at: '2024-01-08T09:00:00Z',
    updated_at: '2024-01-12T11:00:00Z',
    doctor: {
      id: 3,
      user_id: 4,
      department_id: 3,
      title: '主治医师',
      specialty: '呼吸内科',
      introduction: '',
      consultation_fee: 60,
      registration_fee: 50,
      average_rating: 4.5,
      review_count: 67,
      created_at: '',
      updated_at: '',
      user: {
        id: 4,
        username: 'doctor3',
        email: 'doctor3@example.com',
        phone: '13800138003',
        role: 'doctor',
        full_name: '王医生',
        gender: '男',
        birth_date: null,
        avatar_url: '',
        is_active: true,
        created_at: '',
        updated_at: ''
      }
    }
  },
  {
    id: 4,
    patient_id: 1,
    doctor_id: 1,
    appointment_date: '2024-01-10',
    start_time: '15:00:00',
    end_time: '15:30:00',
    status: 'cancelled',
    symptoms: '头晕',
    notes: '',
    cancel_reason: '患者有事取消',
    created_at: '2024-01-05T16:00:00Z',
    updated_at: '2024-01-09T08:00:00Z',
    doctor: {
      id: 1,
      user_id: 2,
      department_id: 1,
      title: '主任医师',
      specialty: '神经内科',
      introduction: '',
      consultation_fee: 100,
      registration_fee: 50,
      average_rating: 4.8,
      review_count: 120,
      created_at: '',
      updated_at: '',
      user: {
        id: 2,
        username: 'doctor1',
        email: 'doctor1@example.com',
        phone: '13800138001',
        role: 'doctor',
        full_name: '张医生',
        gender: '男',
        birth_date: null,
        avatar_url: '',
        is_active: true,
        created_at: '',
        updated_at: ''
      }
    }
  }
]

const statusMap: Record<string, { color: string; text: string }> = {
  pending: { color: 'orange', text: '待确认' },
  confirmed: { color: 'blue', text: '已确认' },
  completed: { color: 'green', text: '已完成' },
  cancelled: { color: 'red', text: '已取消' },
  no_show: { color: 'default', text: '未就诊' }
}

const Dashboard = () => {
  const navigate = useNavigate()

  const columns: ColumnsType<Appointment> = [
    {
      title: '日期',
      dataIndex: 'appointment_date',
      key: 'date',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD')
    },
    {
      title: '时间',
      key: 'time',
      render: (_, record) => `${record.start_time.slice(0, 5)} - ${record.end_time.slice(0, 5)}`
    },
    {
      title: '医生',
      key: 'doctor',
      render: (_, record) => record.doctor?.user?.full_name || '-'
    },
    {
      title: '科室',
      key: 'department',
      render: (_, record) => record.doctor?.specialty || '-'
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const info = statusMap[status] || { color: 'default', text: status }
        return <Tag color={info.color}>{info.text}</Tag>
      }
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Button type="link" onClick={() => navigate(`/appointments/${record.id}`)}>
          详情
        </Button>
      )
    }
  ]

  return (
    <div className="space-y-6">
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} lg={6}>
          <Card className="shadow-sm">
            <Statistic
              title="今日预约"
              value={3}
              prefix={<CalendarOutlined className="text-blue-500" />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card className="shadow-sm">
            <Statistic
              title="待处理"
              value={1}
              prefix={<ClockCircleOutlined className="text-orange-500" />}
              valueStyle={{ color: '#fa8c16' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card className="shadow-sm">
            <Statistic
              title="总收入"
              value={1280}
              precision={2}
              prefix={<DollarOutlined className="text-green-500" />}
              valueStyle={{ color: '#52c41a' }}
              suffix="元"
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} lg={6}>
          <Card className="shadow-sm">
            <Statistic
              title="就诊医生"
              value={5}
              prefix={<UserOutlined className="text-purple-500" />}
              valueStyle={{ color: '#722ed1' }}
            />
          </Card>
        </Col>
      </Row>

      <Card
        className="shadow-sm"
        title="最近预约"
        extra={
          <Button type="link" onClick={() => navigate('/appointments')}>
            查看全部 <RightOutlined />
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={mockAppointments}
          rowKey="id"
          pagination={false}
          scroll={{ x: 600 }}
        />
      </Card>
    </div>
  )
}

export default Dashboard
