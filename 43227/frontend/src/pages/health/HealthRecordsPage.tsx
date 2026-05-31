import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Space, message, Card, Tag, InputNumber, DatePicker, Checkbox } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../../api'
import type { HealthRecord } from '../../types'
import dayjs from 'dayjs'

function HealthRecordsPage() {
  const [data, setData] = useState<HealthRecord[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [beehives, setBeehives] = useState<any[]>([])
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get('/health-records', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as HealthRecord[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchBeehives = async () => {
    try {
      const response = await api.get('/beehives', { params: { page: 1, page_size: 100 } })
      setBeehives(response.data || [])
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    fetchData()
    fetchBeehives()
  }, [page, pageSize])

  const handleAdd = () => {
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleSubmit = async (values: any) => {
    try {
      const submitData = {
        ...values,
        record_date: values.record_date?.format('YYYY-MM-DD'),
      }
      await api.post('/health-records', submitData)
      message.success('创建成功')
      setIsModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '创建失败')
    }
  }

  const columns = [
    {
      title: '蜂箱',
      dataIndex: ['beehive', 'name'],
      key: 'beehive_name',
      render: (_: any, record: HealthRecord) => record.beehive?.name || '-',
    },
    {
      title: '记录日期',
      dataIndex: 'record_date',
      key: 'record_date',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '蜂王状态',
      dataIndex: 'queen_status',
      key: 'queen_status',
      render: (status: string) => status ? <Tag color="green">{status}</Tag> : '-',
    },
    { title: '工蜂数量', dataIndex: 'worker_count', key: 'worker_count' },
    {
      title: '是否有病害',
      dataIndex: 'has_disease',
      key: 'has_disease',
      render: (has: boolean) => has ? <Tag color="red">是</Tag> : <Tag color="green">否</Tag>,
    },
    { title: '病害类型', dataIndex: 'disease_type', key: 'disease_type' },
    {
      title: '严重程度',
      dataIndex: 'disease_severity',
      key: 'disease_severity',
      render: (severity: string) => {
        if (!severity) return '-'
        const colorMap: Record<string, string> = { mild: 'yellow', moderate: 'orange', severe: 'red' }
        return <Tag color={colorMap[severity]}>{severity}</Tag>
      },
    },
    { title: '季节', dataIndex: 'season', key: 'season' },
  ]

  return (
    <Card title="健康记录" extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>添加记录</Button>}>
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
      <Modal title="添加健康记录" open={isModalVisible} onCancel={() => setIsModalVisible(false)} footer={null}>
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="beehive_id" label="选择蜂箱" rules={[{ required: true, message: '请选择蜂箱' }]}>
            <Select placeholder="请选择蜂箱">
              {beehives.map((b) => (
                <Select.Option key={b.id} value={b.id}>{b.name} ({b.code})</Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="record_date" label="记录日期" rules={[{ required: true, message: '请选择日期' }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="queen_status" label="蜂王状态">
            <Select>
              <Select.Option value="normal">正常</Select.Option>
              <Select.Option value="old">老化</Select.Option>
              <Select.Option value="missing">失踪</Select.Option>
              <Select.Option value="new">新王</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="worker_count" label="工蜂数量">
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="has_disease" label="是否有病害" valuePropName="checked">
            <Checkbox />
          </Form.Item>
          <Form.Item name="disease_type" label="病害类型">
            <Input placeholder="如：蜂螨、白垩病等" />
          </Form.Item>
          <Form.Item name="disease_severity" label="严重程度">
            <Select>
              <Select.Option value="mild">轻微</Select.Option>
              <Select.Option value="moderate">中等</Select.Option>
              <Select.Option value="severe">严重</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="treatment" label="处理措施">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>提交</Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default HealthRecordsPage
