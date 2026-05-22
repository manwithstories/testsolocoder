import React, { useState, useEffect } from 'react'
import { Table, Card, Button, Tag, Modal, Form, Input, Select, InputNumber, message, Switch } from 'antd'
import { serviceApi } from '@/services/service'
import { ServiceItem, ServiceCategory } from '@/types'
import { formatPrice, formatDuration } from '@/utils'

const MyServices: React.FC = () => {
  const [services, setServices] = useState<ServiceItem[]>([])
  const [categories, setCategories] = useState<ServiceCategory[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [editingService, setEditingService] = useState<ServiceItem | null>(null)
  const [form] = Form.useForm()

  useEffect(() => {
    loadServices()
    loadCategories()
  }, [])

  const loadServices = async () => {
    setLoading(true)
    try {
      const res = await serviceApi.getMyServices()
      setServices(res)
    } catch (error) {
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const loadCategories = async () => {
    try {
      const res = await serviceApi.getCategories()
      setCategories(res)
    } catch (error) {
      console.error(error)
    }
  }

  const handleAdd = () => {
    setEditingService(null)
    form.resetFields()
    form.setFieldsValue({ is_active: true, price_unit: 'hour' })
    setModalVisible(true)
  }

  const handleEdit = (service: ServiceItem) => {
    setEditingService(service)
    form.setFieldsValue(service)
    setModalVisible(true)
  }

  const handleToggleActive = async (service: ServiceItem) => {
    try {
      await serviceApi.update(service.id, { is_active: !service.is_active })
      message.success('状态更新成功')
      loadServices()
    } catch (error) {
      console.error(error)
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingService) {
        await serviceApi.update(editingService.id, values)
        message.success('更新成功')
      } else {
        await serviceApi.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      loadServices()
    } catch (error) {
      console.error(error)
    }
  }

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: '服务名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '分类',
      dataIndex: 'category_name',
      key: 'category',
      render: (_: any, record: ServiceItem) => record.category?.name || '-',
    },
    {
      title: '价格',
      dataIndex: 'base_price',
      key: 'base_price',
      render: (text: number, record: ServiceItem) => (
        <span style={{ color: '#ff4d4f', fontWeight: 600 }}>
          {formatPrice(text)}/{record.price_unit}
        </span>
      ),
    },
    {
      title: '时长',
      key: 'duration',
      render: (_: any, record: ServiceItem) => (
        <span>{formatDuration(record.min_duration)} - {formatDuration(record.max_duration)}</span>
      ),
    },
    {
      title: '评分',
      dataIndex: 'rating',
      key: 'rating',
      render: (text: number, record: ServiceItem) => <span>{text} ({record.review_count})</span>,
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (text: boolean, record: ServiceItem) => (
        <Switch
          checked={text}
          onChange={() => handleToggleActive(record)}
        />
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: ServiceItem) => (
        <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
      ),
    },
  ]

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">我的服务</h1>
        <Button type="primary" onClick={handleAdd}>
          添加服务
        </Button>
      </div>

      <Card>
        <Table
          rowKey="id"
          loading={loading}
          dataSource={services}
          columns={columns}
          pagination={false}
        />
      </Card>

      <Modal
        title={editingService ? '编辑服务' : '添加服务'}
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="category_id"
            label="服务分类"
            rules={[{ required: true, message: '请选择服务分类' }]}
          >
            <Select
              placeholder="请选择服务分类"
              options={categories.map((cat) => ({
                label: cat.name,
                value: cat.id,
              }))}
            />
          </Form.Item>
          <Form.Item
            name="name"
            label="服务名称"
            rules={[{ required: true, message: '请输入服务名称' }]}
          >
            <Input placeholder="请输入服务名称" />
          </Form.Item>
          <Form.Item
            name="description"
            label="服务描述"
            rules={[{ required: true, message: '请输入服务描述' }]}
          >
            <Input.TextArea rows={4} placeholder="请输入服务描述" />
          </Form.Item>
          <Form.Item
            name="base_price"
            label="基础价格"
            rules={[{ required: true, message: '请输入基础价格' }]}
          >
            <InputNumber
              min={0}
              precision={2}
              style={{ width: '100%' }}
              placeholder="请输入基础价格"
            />
          </Form.Item>
          <Form.Item
            name="price_unit"
            label="价格单位"
            rules={[{ required: true, message: '请选择价格单位' }]}
          >
            <Select
              placeholder="请选择价格单位"
              options={[
                { label: '小时', value: 'hour' },
                { label: '次', value: 'time' },
                { label: '天', value: 'day' },
                { label: '月', value: 'month' },
              ]}
            />
          </Form.Item>
          <Form.Item
            name="min_duration"
            label="最短时长(分钟)"
            rules={[{ required: true, message: '请输入最短时长' }]}
          >
            <InputNumber min={15} step={15} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item
            name="max_duration"
            label="最长时长(分钟)"
            rules={[{ required: true, message: '请输入最长时长' }]}
          >
            <InputNumber min={30} step={15} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="is_active" label="是否上架" valuePropName="checked">
            <Switch />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default MyServices
