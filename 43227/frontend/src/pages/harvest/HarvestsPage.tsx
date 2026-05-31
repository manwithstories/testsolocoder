import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Space, message, Card, Tag, InputNumber, DatePicker } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../../api'
import type { Harvest } from '../../types'
import dayjs from 'dayjs'

function HarvestsPage() {
  const [data, setData] = useState<Harvest[]>([])
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
      const response = await api.get('/harvests', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as Harvest[])
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

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/harvests/${id}`)
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
        harvest_date: values.harvest_date?.format('YYYY-MM-DD'),
      }
      await api.post('/harvests', submitData)
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
      render: (_: any, record: Harvest) => record.beehive?.name || '-',
    },
    {
      title: '采收日期',
      dataIndex: 'harvest_date',
      key: 'harvest_date',
      render: (date: string) => dayjs(date).format('YYYY-MM-DD'),
    },
    { title: '蜂蜜类型', dataIndex: 'honey_type', key: 'honey_type' },
    { title: '产量', dataIndex: 'quantity', key: 'quantity', render: (q: number) => `${q} kg` },
    {
      title: '质量等级',
      dataIndex: 'quality',
      key: 'quality',
      render: (q: string) => {
        const colorMap: Record<string, string> = { normal: 'blue', good: 'green', premium: 'gold' }
        return <Tag color={colorMap[q]}>{q}</Tag>
      },
    },
    { title: '批次号', dataIndex: 'batch_code', key: 'batch_code' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Harvest) => (
        <Space>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <Card title="采收记录" extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>添加采收</Button>}>
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
      <Modal title="添加采收记录" open={isModalVisible} onCancel={() => setIsModalVisible(false)} footer={null}>
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          <Form.Item name="beehive_id" label="选择蜂箱" rules={[{ required: true, message: '请选择蜂箱' }]}>
            <Select placeholder="请选择蜂箱">
              {beehives.map((b) => (
                <Select.Option key={b.id} value={b.id}>{b.name} ({b.code})</Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="harvest_date" label="采收日期" rules={[{ required: true, message: '请选择日期' }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="honey_type" label="蜂蜜类型" rules={[{ required: true, message: '请输入蜂蜜类型' }]}>
            <Input placeholder="如：洋槐蜜、枣花蜜等" />
          </Form.Item>
          <Form.Item name="quantity" label="产量(kg)" rules={[{ required: true, message: '请输入产量' }]}>
            <InputNumber min={0} step={0.1} />
          </Form.Item>
          <Form.Item name="quality" label="质量等级">
            <Select>
              <Select.Option value="normal">普通</Select.Option>
              <Select.Option value="good">良好</Select.Option>
              <Select.Option value="premium">优质</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="batch_code" label="批次号">
            <Input placeholder="留空自动生成" />
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

export default HarvestsPage
