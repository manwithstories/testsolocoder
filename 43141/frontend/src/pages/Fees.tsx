import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, message, Tag, Card, Statistic, Row, Col, Popconfirm } from 'antd'
import { PlusOutlined, CheckOutlined, DeleteOutlined, FileTextOutlined } from '@ant-design/icons'
import api from '../api'
import type { Fee, Season, Team } from '../types'

export default function Fees({ isAdmin }: { isAdmin: boolean }) {
  const [fees, setFees] = useState<Fee[]>([])
  const [seasons, setSeasons] = useState<Season[]>([])
  const [teams, setTeams] = useState<Team[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchData = () => {
    api.get('/fees').then((res) => setFees(res.data))
    api.get('/leagues').then((res) => {
      const all: Season[] = []
      res.data.forEach((l: any) => l.seasons?.forEach((s: Season) => all.push(s)))
      setSeasons(all)
    })
    api.get('/teams').then((res) => setTeams(res.data))
  }

  useEffect(() => { fetchData() }, [])

  const handleCreate = () => {
    form.resetFields()
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      await api.post('/fees', values)
      message.success('创建成功')
      setModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '创建失败')
    }
  }

  const handleMarkPaid = async (id: number) => {
    try {
      await api.put(`/fees/${id}/paid`)
      message.success('已标记为已付款')
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/fees/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '删除失败')
    }
  }

  const totalUnpaid = fees.filter((f) => f.status === 'unpaid').reduce((sum, f) => sum + f.amount, 0)
  const totalPaid = fees.filter((f) => f.status === 'paid').reduce((sum, f) => sum + f.amount, 0)

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '发票号', dataIndex: 'invoice_no' },
    { title: '赛季', dataIndex: ['season', 'name'] },
    { title: '球队', dataIndex: ['team', 'name'] },
    { title: '类型', dataIndex: 'type', render: (v: string) => <Tag>{v}</Tag> },
    { title: '金额', dataIndex: 'amount', render: (v: number) => `¥${v.toFixed(2)}` },
    {
      title: '状态', dataIndex: 'status', render: (s: string) => (
        <Tag color={s === 'paid' ? 'green' : 'orange'}>{s === 'paid' ? '已付款' : '未付款'}</Tag>
      )
    },
    { title: '付款时间', dataIndex: 'paid_at', render: (v: string) => v ? new Date(v).toLocaleString() : '-' },
    {
      title: '操作', width: 200, render: (_: any, r: Fee) => (
        <Space>
          {isAdmin && r.status === 'unpaid' && (
            <Button type="primary" icon={<CheckOutlined />} size="small" onClick={() => handleMarkPaid(r.id)}>
              标记已付
            </Button>
          )}
          {isAdmin && <Popconfirm title="确认删除?" onConfirm={() => handleDelete(r.id)}>
            <Button danger icon={<DeleteOutlined />} size="small">删除</Button>
          </Popconfirm>}
        </Space>
      )
    }
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>费用管理</h2>
        {isAdmin && <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>创建费用</Button>}
      </div>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic title="已收费用" value={totalPaid} prefix="¥" precision={2} valueStyle={{ color: '#52c41a' }} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="待收费用" value={totalUnpaid} prefix="¥" precision={2} valueStyle={{ color: '#faad14' }} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="总发票数" value={fees.length} />
          </Card>
        </Col>
      </Row>

      <Card title="费用列表">
        <Table rowKey="id" dataSource={fees} columns={columns} />
      </Card>

      <Modal
        title="创建费用"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="season_id" label="赛季" rules={[{ required: true }]}>
            <Select options={seasons.map((s) => ({ value: s.id, label: s.name }))} />
          </Form.Item>
          <Form.Item name="team_id" label="球队" rules={[{ required: true }]}>
            <Select options={teams.map((t) => ({ value: t.id, label: t.name }))} showSearch optionFilterProp="label" />
          </Form.Item>
          <Form.Item name="type" label="费用类型" rules={[{ required: true }]}>
            <Select options={[
              { value: 'registration', label: '报名费' },
              { value: 'venue', label: '场地费' },
              { value: 'other', label: '其他' },
            ]} />
          </Form.Item>
          <Form.Item name="amount" label="金额" rules={[{ required: true }]}>
            <InputNumber min={0} step={0.01} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="note" label="备注">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
