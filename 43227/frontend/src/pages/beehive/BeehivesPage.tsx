import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, Select, Space, message, Card, Tag, InputNumber } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import api from '../../api'
import type { Beehive } from '../../types'

function BeehivesPage() {
  const navigate = useNavigate()
  const [data, setData] = useState<Beehive[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [isModalVisible, setIsModalVisible] = useState(false)
  const [editingRecord, setEditingRecord] = useState<Beehive | null>(null)
  const [form] = Form.useForm()

  const fetchData = async () => {
    setLoading(true)
    try {
      const response = await api.get('/beehives', {
        params: { page, page_size: pageSize },
      })
      setData(response.data as Beehive[])
      setTotal(response.total || 0)
    } catch (error: any) {
      message.error(error.message || '获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [page, pageSize])

  const handleAdd = () => {
    setEditingRecord(null)
    form.resetFields()
    setIsModalVisible(true)
  }

  const handleEdit = (record: Beehive) => {
    setEditingRecord(record)
    form.setFieldsValue({
      name: record.name,
      latitude: record.latitude,
      longitude: record.longitude,
      region: record.region,
      bee_species: record.bee_species,
      group_name: record.group_name,
      queen_status: record.queen_status,
      worker_count: record.worker_count,
      notes: record.notes,
    })
    setIsModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/beehives/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  const handleSubmit = async (values: any) => {
    try {
      if (editingRecord) {
        await api.put(`/beehives/${editingRecord.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/beehives', values)
        message.success('创建成功')
      }
      setIsModalVisible(false)
      fetchData()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const columns = [
    { title: '名称', dataIndex: 'name', key: 'name' },
    { title: '编号', dataIndex: 'code', key: 'code' },
    { title: '区域', dataIndex: 'region', key: 'region' },
    { title: '蜜蜂品种', dataIndex: 'bee_species', key: 'bee_species' },
    { title: '分组', dataIndex: 'group_name', key: 'group_name' },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          active: 'green',
          inactive: 'default',
          harvesting: 'orange',
        }
        return <Tag color={colorMap[status]}>{status}</Tag>
      },
    },
    {
      title: '健康状态',
      dataIndex: 'health_status',
      key: 'health_status',
      render: (status: string) => {
        const colorMap: Record<string, string> = {
          healthy: 'green',
          warning: 'orange',
          critical: 'red',
        }
        return <Tag color={colorMap[status]}>{status}</Tag>
      },
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Beehive) => (
        <Space>
          <Button type="link" icon={<EyeOutlined />} onClick={() => navigate(`/beehives/${record.id}`)}>
            查看
          </Button>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Button type="link" danger icon={<DeleteOutlined />} onClick={() => handleDelete(record.id)}>
            删除
          </Button>
        </Space>
      ),
    },
  ]

  return (
    <Card title="蜂箱管理" extra={<Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>添加蜂箱</Button>}>
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
        title={editingRecord ? '编辑蜂箱' : '添加蜂箱'}
        open={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleSubmit} layout="vertical">
          {!editingRecord && (
            <Form.Item name="code" label="编号" rules={[{ required: true, message: '请输入编号' }]}>
              <Input placeholder="蜂箱编号" />
            </Form.Item>
          )}
          <Form.Item name="name" label="名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input placeholder="蜂箱名称" />
          </Form.Item>
          <Form.Item name="latitude" label="纬度" rules={[{ required: true, message: '请输入纬度' }]}>
            <InputNumber style={{ width: '100%' }} min={-90} max={90} step={0.0000001} />
          </Form.Item>
          <Form.Item name="longitude" label="经度" rules={[{ required: true, message: '请输入经度' }]}>
            <InputNumber style={{ width: '100%' }} min={-180} max={180} step={0.0000001} />
          </Form.Item>
          <Form.Item name="region" label="区域">
            <Input placeholder="区域名称" />
          </Form.Item>
          <Form.Item name="bee_species" label="蜜蜂品种">
            <Input placeholder="蜜蜂品种" />
          </Form.Item>
          <Form.Item name="group_name" label="分组">
            <Input placeholder="分组名称" />
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
          <Form.Item name="notes" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              提交
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default BeehivesPage
