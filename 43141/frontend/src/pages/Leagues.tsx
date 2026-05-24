import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, Input, Select, DatePicker, InputNumber, Space, Popconfirm, message, Tag, Card, Divider } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import api from '../api'
import type { League, Season } from '../types'

export default function Leagues({ isAdmin }: { isAdmin: boolean }) {
  const [leagues, setLeagues] = useState<League[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [seasonModalOpen, setSeasonModalOpen] = useState(false)
  const [editingLeague, setEditingLeague] = useState<League | null>(null)
  const [currentLeagueId, setCurrentLeagueId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [seasonForm] = Form.useForm()

  const fetchData = () => {
    api.get('/leagues').then((res) => setLeagues(res.data))
  }

  useEffect(() => { fetchData() }, [])

  const handleCreate = () => {
    setEditingLeague(null)
    form.resetFields()
    setModalOpen(true)
  }

  const handleEdit = (league: League) => {
    setEditingLeague(league)
    form.setFieldsValue(league)
    setModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/leagues/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '删除失败')
    }
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingLeague) {
        await api.put(`/leagues/${editingLeague.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/leagues', values)
        message.success('创建成功')
      }
      setModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const handleCreateSeason = (leagueId: number) => {
    setCurrentLeagueId(leagueId)
    seasonForm.resetFields()
    seasonForm.setFieldsValue({
      format: 'round_robin',
      group_count: 1,
      points_for_win: 3,
      points_for_draw: 1,
      points_for_loss: 0,
      max_teams: 16
    })
    setSeasonModalOpen(true)
  }

  const handleSeasonSubmit = async () => {
    try {
      const values = await seasonForm.validateFields()
      const payload: any = { ...values }
      if (values.start_date) payload.start_date = dayjs(values.start_date).format('YYYY-MM-DD')
      if (values.end_date) payload.end_date = dayjs(values.end_date).format('YYYY-MM-DD')
      await api.post(`/leagues/${currentLeagueId}/seasons`, payload)
      message.success('赛季创建成功')
      setSeasonModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '创建失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '名称', dataIndex: 'name' },
    { title: '运动', dataIndex: 'sport' },
    { title: '状态', dataIndex: 'status', render: (s: string) => <Tag color={s === 'active' ? 'green' : 'blue'}>{s}</Tag> },
    { title: '赛季数', render: (_: any, r: League) => r.seasons?.length || 0 },
    {
      title: '操作', width: 200, render: (_: any, r: League) => (
        <Space>
          {isAdmin && <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(r)}>编辑</Button>}
          {isAdmin && <Button type="primary" size="small" onClick={() => handleCreateSeason(r.id)}>添加赛季</Button>}
          {isAdmin && <Popconfirm title="确认删除?" onConfirm={() => handleDelete(r.id)}>
            <Button danger icon={<DeleteOutlined />} size="small">删除</Button>
          </Popconfirm>}
        </Space>
      )
    }
  ]

  const expandedRowRender = (league: League) => (
    <Card size="small" title="赛季列表">
      {league.seasons && league.seasons.length > 0 ? (
        <Table
          size="small"
          dataSource={league.seasons}
          rowKey="id"
          pagination={false}
          columns={[
            { title: 'ID', dataIndex: 'id', width: 60 },
            { title: '名称', dataIndex: 'name' },
            { title: '赛制', dataIndex: 'format' },
            { title: '开始日期', dataIndex: 'start_date', render: (d: string) => dayjs(d).format('YYYY-MM-DD') },
            { title: '结束日期', dataIndex: 'end_date', render: (d: string) => dayjs(d).format('YYYY-MM-DD') },
            { title: '状态', dataIndex: 'status', render: (s: string) => <Tag>{s}</Tag> },
            { title: '积分/胜', dataIndex: 'points_for_win' },
            { title: '报名费用', dataIndex: 'registration_fee' },
          ]}
        />
      ) : <p style={{ color: '#999' }}>暂无赛季</p>}
    </Card>
  )

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>联赛管理</h2>
        {isAdmin && <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>创建联赛</Button>}
      </div>
      <Table
        rowKey="id"
        dataSource={leagues}
        columns={columns}
        expandable={{ expandedRowRender }}
      />
      <Modal
        title={editingLeague ? '编辑联赛' : '创建联赛'}
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="sport" label="运动类型">
            <Select options={[
              { value: 'football', label: '足球' },
              { value: 'basketball', label: '篮球' },
              { value: 'volleyball', label: '排球' },
              { value: 'badminton', label: '羽毛球' },
              { value: 'tennis', label: '网球' },
            ]} />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="创建赛季"
        open={seasonModalOpen}
        onCancel={() => setSeasonModalOpen(false)}
        onOk={handleSeasonSubmit}
        width={600}
      >
        <Form form={seasonForm} layout="vertical">
          <Form.Item name="name" label="赛季名称" rules={[{ required: true }]}>
            <Input placeholder="如 2025-2026赛季" />
          </Form.Item>
          <Form.Item name="start_date" label="开始日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="end_date" label="结束日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Divider orientation="left">赛制规则</Divider>
          <Form.Item name="format" label="赛制">
            <Select options={[
              { value: 'round_robin', label: '单循环' },
              { value: 'double_round_robin', label: '双循环' },
              { value: 'group_knockout', label: '小组赛+淘汰赛' },
            ]} />
          </Form.Item>
          <Form.Item name="group_count" label="分组数">
            <InputNumber min={1} max={8} />
          </Form.Item>
          <Form.Item name="max_teams" label="最大参赛队数">
            <InputNumber min={2} max={64} />
          </Form.Item>
          <Divider orientation="left">积分制度</Divider>
          <Form.Item name="points_for_win" label="胜场积分">
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="points_for_draw" label="平局积分">
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="points_for_loss" label="负场积分">
            <InputNumber min={0} />
          </Form.Item>
          <Divider orientation="left">费用</Divider>
          <Form.Item name="registration_fee" label="报名费">
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="venue_fee" label="场地费">
            <InputNumber min={0} />
          </Form.Item>
          <Form.Item name="custom_rules" label="自定义规则">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="awards" label="奖项设置">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
