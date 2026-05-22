import React, { useEffect, useState } from 'react'
import { Table, Tag, Button, Space, Modal, Form, Input, InputNumber, Select, DatePicker, message, Popconfirm, Switch } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { couponApi } from '@/services'
import { Coupon } from '@/types'
import dayjs from 'dayjs'

const AdminCoupons: React.FC = () => {
  const [coupons, setCoupons] = useState<Coupon[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  const loadCoupons = async () => {
    setLoading(true)
    try {
      const res = await couponApi.list()
      if (res.code === 0 && res.data) {
        setCoupons(res.data)
      }
    } catch (error) {
      console.error('Failed to load coupons:', error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadCoupons()
  }, [])

  const handleCreate = async () => {
    try {
      const values = await form.validateFields()
      if (values.expires_at) {
        values.expires_at = values.expires_at.toISOString()
      }
      const res = await couponApi.create(values)
      if (res.code === 0) {
        message.success('优惠券创建成功')
        setModalVisible(false)
        form.resetFields()
        loadCoupons()
      }
    } catch (error: any) {
      if (error.errorFields) return
      message.error(error.message || '创建失败')
    }
  }

  const handleToggleActive = async (id: string, active: boolean) => {
    try {
      await couponApi.update(id, { is_active: active })
      message.success('状态更新成功')
      loadCoupons()
    } catch (error: any) {
      message.error(error.message || '操作失败')
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await couponApi.delete(id)
      message.success('删除成功')
      loadCoupons()
    } catch (error: any) {
      message.error(error.message || '删除失败')
    }
  }

  const columns = [
    { title: '优惠券码', dataIndex: 'code', key: 'code' },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => <Tag color={type === 'fixed' ? 'blue' : 'green'}>{type === 'fixed' ? '固定金额' : '折扣百分比'}</Tag>,
    },
    {
      title: '面值',
      dataIndex: 'value',
      key: 'value',
      render: (val: number, record: Coupon) => record.type === 'fixed' ? `¥${val}` : `${val}%`,
    },
    { title: '最低消费', dataIndex: 'min_amount', key: 'min_amount', render: (val: number) => `¥${val}` },
    {
      title: '使用情况',
      key: 'usage',
      render: (_: any, record: Coupon) => `${record.used_count}/${record.total_count || '∞'}`,
    },
    {
      title: '有效期',
      dataIndex: 'expires_at',
      key: 'expires_at',
      render: (val: string) => val ? dayjs(val).format('YYYY-MM-DD') : '永久有效',
    },
    {
      title: '状态',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (active: boolean, record: Coupon) => (
        <Switch
          checked={active}
          onChange={(checked) => handleToggleActive(record.id, checked)}
        />
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Coupon) => (
        <Popconfirm title="确定删除？" onConfirm={() => handleDelete(record.id)}>
          <Button type="link" danger>删除</Button>
        </Popconfirm>
      ),
    },
  ]

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <h2>优惠券管理</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
          创建优惠券
        </Button>
      </div>
      <Table
        rowKey="id"
        columns={columns}
        dataSource={coupons}
        loading={loading}
        pagination={false}
      />
      <Modal
        title="创建优惠券"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={handleCreate}
        okText="创建"
      >
        <Form form={form} layout="vertical">
          <Form.Item name="code" label="优惠券码" rules={[{ required: true }]}>
            <Input placeholder="如：NEWYEAR2024" />
          </Form.Item>
          <Form.Item name="type" label="类型" rules={[{ required: true }]} initialValue="fixed">
            <Select>
              <Select.Option value="fixed">固定金额</Select.Option>
              <Select.Option value="percent">折扣百分比</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="value" label="面值" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} min={0} placeholder="金额或百分比" />
          </Form.Item>
          <Form.Item name="min_amount" label="最低消费" initialValue={0}>
            <InputNumber style={{ width: '100%' }} min={0} />
          </Form.Item>
          <Form.Item name="max_discount" label="最大折扣金额" initialValue={0}>
            <InputNumber style={{ width: '100%' }} min={0} placeholder="0表示不限制" />
          </Form.Item>
          <Form.Item name="total_count" label="可使用次数" initialValue={0}>
            <InputNumber style={{ width: '100%' }} min={0} placeholder="0表示不限制" />
          </Form.Item>
          <Form.Item name="expires_at" label="有效期">
            <DatePicker style={{ width: '100%' }} placeholder="不填则永久有效" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default AdminCoupons
