import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, DatePicker, Space, message, Card, Tag } from 'antd'
import { PlusOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons'
import api from '../../api'
import type { Inspection, Inventory } from '../../types'
import { useAuthStore } from '../../store/authStore'
import dayjs from 'dayjs'

function InspectionsPage() {
  const { user } = useAuthStore()
  const [data, setData] = useState<Inspection[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [isResultModalVisible, setIsResultModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState<Inspection | null>(null)
  const [inventoryList, setInventoryList] = useState<Inventory[]>([])
  const [form] = Form.useForm()
  const [resultForm] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get('/inspections', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as Inspection[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchInventory = async () => {
    try {
      const response = await api.get('/inventory', { params: { page: 1, page_size: 100 } })
      setInventoryList(response.data || [])
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    fetchData()
    if (user?.role === 'beekeeper') {
      fetchInventory()
    }
  }, [page, pageSize, user?.role])

  const handleAdd = () => {
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleAssign = async (id: number) => {
    try {
      await api.put(`/inspections/${id}/assign`)
      message.success('预约已接受')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleSubmitResult = async (record: Inspection) => {
    setEditingRecord(record)
    resultForm.resetFields()
    setIsResultModalVisible(true)
  }

  const handleCancel = async (id: number) => {
    try {
      await api.put(`/inspections/${id}/cancel`)
      message.success('预约已取消')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '取消失败')
    }
  }

  const handleSubmitAppointment = async (values: any) => {
    try {
      const selectedInventory = inventoryList.find((inv) => inv.id === values.inventory_id)
      await api.post('/inspections', {
        ...values,
        appointment_date: values.appointment_date?.format('YYYY-MM-DD'),
        batch_code: selectedInventory?.batch_code || '',
      })
      message.success('预约成功')
      setIsModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '预约失败')
    }
  }

  const handleSubmitResultForm = async (values: any) => {
    try {
      await api.put(`/inspections/${editingRecord?.id}/result`, values)
      message.success('检测结果已提交')
      setIsResultModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '提交失败')
    }
  }

  const statusColorMap: Record<string, string> = {
    pending: 'default',
    confirmed: 'blue',
    in_progress: 'cyan',
    completed: 'green',
    cancelled: 'red',
  }

  const resultColorMap: Record<string, string> = {
    pass: 'green',
    fail: 'red',
    conditional: 'orange',
  }

  const columns = [
    { title: '批次号', dataIndex: 'batch_code', key: 'batch_code' },
    {
      title: user?.role === 'beekeeper' ? '预约人' : '预约人',
      key: 'user',
      render: (_: any, record: Inspection) => record.user?.username,
    },
    {
      title: '预约日期',
      dataIndex: 'appointment_date',
      key: 'appointment_date',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={statusColorMap[status]}>{status}</Tag>,
    },
    {
      title: '检测结果',
      dataIndex: 'result',
      key: 'result',
      render: (result: string) => result ? <Tag color={resultColorMap[result]}>{result}</Tag> : '-',
    },
    {
      title: '等级',
      dataIndex: 'grade',
      key: 'grade',
      render: (grade: string) => {
        if (!grade) return '-'
        const colorMap: Record<string, string> = { grade_a: 'gold', grade_b: 'blue', grade_c: 'green' }
        return <Tag color={colorMap[grade]}>{grade}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Inspection) => (
        <Space>
          {user?.role === 'inspector' && record.status === 'pending' && (
            <Button type="link" icon={<CheckOutlined />} onClick={() => handleAssign(record.id)}>
              接受预约
            </Button>
          )}
          {user?.role === 'inspector' && (record.status === 'confirmed' || record.status === 'in_progress') && (
            <Button type="link" onClick={() => handleSubmitResult(record)}>
              提交结果
            </Button>
          )}
          {user?.role === 'beekeeper' && (record.status === 'pending' || record.status === 'confirmed') && (
            <Button type="link" danger icon={<CloseOutlined />} onClick={() => handleCancel(record.id)}>
              取消
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <Card
      title={user?.role === 'beekeeper' ? '检测预约' : '检测任务'}
      extra={user?.role === 'beekeeper' && (
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          预约检测
        </Button>
      )}
    >
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize,
          total,
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps)
          },
        }}
      />
      <Modal
        title="预约检测"
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleSubmitAppointment} layout="vertical">
          <Form.Item name="inventory_id" label="选择库存" rules={[{ required: true, message: '请选择库存' }]}>
            <Select placeholder="请选择要检测的库存">
              {inventoryList.map((inv) => (
                <Select.Option key={inv.id} value={inv.id}>
                  {inv.honey_type} - 批次{inv.batch_code}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="appointment_date" label="预约日期" rules={[{ required: true, message: '请选择预约日期' }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交预约</Button>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="提交检测结果"
        open={isResultModalVisible}
        onCancel={() => setIsResultModalVisible(false)}
        footer={null}
      >
        <Form form={resultForm} onFinish={handleSubmitResultForm} layout="vertical">
          <Form.Item name="result" label="检测结果" rules={[{ required: true, message: '请选择检测结果' }]}>
            <Select>
              <Select.Option value="pass">合格</Select.Option>
              <Select.Option value="fail">不合格</Select.Option>
              <Select.Option value="conditional">有条件合格</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="grade" label="等级">
            <Select>
              <Select.Option value="grade_a">A级</Select.Option>
              <Select.Option value="grade_b">B级</Select.Option>
              <Select.Option value="grade_c">C级</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="report_url" label="检测报告链接">
            <Input />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交结果</Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default InspectionsPage
