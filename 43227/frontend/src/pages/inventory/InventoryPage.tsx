import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Space, message, Card, Tag, InputNumber, DatePicker, Alert } from 'antd'
import { EditOutlined, DeleteOutlined, WarningOutlined } from '@ant-design/icons'
import api from '../../api'
import type { Inventory } from '../../types'
import dayjs from 'dayjs'

function InventoryPage() {
  const [data, setData] = useState<Inventory[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState<Inventory | null>(null)
  const [alerts, setAlerts] = useState<any>({})
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get('/inventory', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as Inventory[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchAlerts = async () => {
    try {
      const response = await api.get('/inventory/alerts')
      setAlerts(response.data)
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    fetchData()
    fetchAlerts()
  }, [page, pageSize])

  const handleEdit = (record: Inventory) => {
    setEditingRecord(record)
    form.setFieldsValue({
      quantity: record.quantity,
      expiry_date: record.expiry_date ? dayjs(record.expiry_date) : null,
      inspection_report: record.inspection_report,
      grade: record.grade,
      threshold: record.threshold,
      price: record.price,
    })
    setIsModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/inventory/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  const handleSubmit = async (values: any) => {
    try {
      const submitData = {
        ...values,
        expiry_date: values.expiry_date?.format('YYYY-MM-DD'),
      }
      await api.put(`/inventory/${editingRecord?.id}`, submitData)
      message.success('更新成功')
      setIsModalVisible(false)
      fetchData()
      fetchAlerts()
    } catch (error: any) {
      message.error(error.message || '更新失败')
    }
  }

  const statusColorMap: Record<string, string> = {
    in_stock: 'green',
    low_stock: 'orange',
    expiring_soon: 'yellow',
    expired: 'red',
    sold_out: 'default',
  }

  const gradeColorMap: Record<string, string> = {
    ungraded: 'default',
    grade_a: 'gold',
    grade_b: 'blue',
    grade_c: 'green',
  }

  const columns = [
    { title: '批次号', dataIndex: 'batch_code', key: 'batch_code' },
    { title: '蜂蜜类型', dataIndex: 'honey_type', key: 'honey_type' },
    { title: '库存数量', dataIndex: 'quantity', key: 'quantity', render: (q: number) => `${q} kg` },
    {
      title: '等级',
      dataIndex: 'grade',
      key: 'grade',
      render: (g: string) => <Tag color={gradeColorMap[g]}>{g}</Tag>,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (s: string) => <Tag color={statusColorMap[s]}>{s}</Tag>,
    },
    {
      title: '保质期',
      dataIndex: 'expiry_date',
      key: 'expiry_date',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    { title: '价格', dataIndex: 'price', key: 'price', render: (p: number) => p ? `¥${p}` : '-' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Inventory) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  return (
    <div>
      {(alerts.low_stock?.length > 0 || alerts.expiring_soon?.length > 0) && (
        <div style={{ marginBottom: 16 }}>
          {alerts.low_stock?.length > 0 && (
            <Alert
              message="库存预警"
              description={`有 ${alerts.low_stock.length} 项库存低于阈值`}
              type="warning"
              showIcon
              icon={<WarningOutlined />}
              style={{ marginBottom: 8 }}
            />
          )}
          {alerts.expiring_soon?.length > 0 && (
            <Alert
              message="临期预警"
              description={`有 ${alerts.expiring_soon.length} 项库存即将到期`}
              type="error"
              showIcon
              icon={<WarningOutlined />}
            />
          )}
        </div>
      )}
      <Card title="库存管理">
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
        <Modal title="编辑库存" open={isModalVisible} onCancel={() => setIsModalVisible(false)} footer={null}>
          <Form form={form} onFinish={handleSubmit} layout="vertical">
            <Form.Item name="quantity" label="库存数量(kg)">
              <InputNumber min={0} step={0.1} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="expiry_date" label="保质期">
              <DatePicker style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="inspection_report" label="检测报告链接">
              <Input />
            </Form.Item>
            <Form.Item name="grade" label="等级">
              <Select>
                <Select.Option value="ungraded">未评级</Select.Option>
                <Select.Option value="grade_a">A级</Select.Option>
                <Select.Option value="grade_b">B级</Select.Option>
                <Select.Option value="grade_c">C级</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item name="threshold" label="预警阈值(kg)">
              <InputNumber min={0} step={1} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="price" label="价格(¥)">
              <InputNumber min={0} step={0.01} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" block>提交</Button>
            </Form.Item>
          </Form>
        </Modal>
      </Card>
    </div>
  )
}

export default InventoryPage
