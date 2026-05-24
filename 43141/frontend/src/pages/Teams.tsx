import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, Popconfirm, message, DatePicker, Tag, Card } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, UserAddOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import api from '../api'
import type { Team, Player, Registration, Season } from '../types'

export default function Teams({ isAdmin }: { isAdmin: boolean }) {
  const [teams, setTeams] = useState<Team[]>([])
  const [seasons, setSeasons] = useState<Season[]>([])
  const [registrations, setRegistrations] = useState<Registration[]>([])
  const [teamModalOpen, setTeamModalOpen] = useState(false)
  const [playerModalOpen, setPlayerModalOpen] = useState(false)
  const [regModalOpen, setRegModalOpen] = useState(false)
  const [editingTeam, setEditingTeam] = useState<Team | null>(null)
  const [currentTeamId, setCurrentTeamId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [playerForm] = Form.useForm()
  const [regForm] = Form.useForm()

  const fetchData = () => {
    api.get('/teams').then((res) => setTeams(res.data))
    api.get('/leagues').then((res) => {
      const all: Season[] = []
      res.data.forEach((l: any) => l.seasons?.forEach((s: Season) => all.push(s)))
      setSeasons(all)
    })
    if (isAdmin) {
      api.get('/registrations').then((res) => setRegistrations(res.data))
    }
  }

  useEffect(() => { fetchData() }, [])

  const handleCreateTeam = () => {
    setEditingTeam(null)
    form.resetFields()
    setTeamModalOpen(true)
  }

  const handleEditTeam = (team: Team) => {
    setEditingTeam(team)
    form.setFieldsValue(team)
    setTeamModalOpen(true)
  }

  const handleDeleteTeam = async (id: number) => {
    try {
      await api.delete(`/teams/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '删除失败')
    }
  }

  const handleTeamSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (editingTeam) {
        await api.put(`/teams/${editingTeam.id}`, values)
        message.success('更新成功')
      } else {
        await api.post('/teams', values)
        message.success('创建成功')
      }
      setTeamModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const handleAddPlayer = (teamId: number) => {
    setCurrentTeamId(teamId)
    playerForm.resetFields()
    setPlayerModalOpen(true)
  }

  const handlePlayerSubmit = async () => {
    try {
      const values = await playerForm.validateFields()
      const payload: any = { ...values }
      if (values.birth_date) payload.birth_date = dayjs(values.birth_date).format('YYYY-MM-DD')
      await api.post(`/teams/${currentTeamId}/players`, payload)
      message.success('球员添加成功')
      setPlayerModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '添加失败')
    }
  }

  const handleRegisterTeam = (teamId: number) => {
    setCurrentTeamId(teamId)
    regForm.resetFields()
    setRegModalOpen(true)
  }

  const handleRegSubmit = async () => {
    try {
      const values = await regForm.validateFields()
      await api.post(`/teams/${currentTeamId}/register`, values)
      message.success('报名成功')
      setRegModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '报名失败')
    }
  }

  const handleApproveReg = async (id: number) => {
    try {
      await api.put(`/registrations/${id}/approve`)
      message.success('已批准')
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '球队名称', dataIndex: 'name' },
    { title: '队长', dataIndex: ['captain', 'full_name'], render: (v: string) => v || '-' },
    { title: '联系邮箱', dataIndex: 'contact_email' },
    { title: '联系电话', dataIndex: 'contact_phone' },
    { title: '球员数', render: (_: any, r: Team) => r.players?.length || 0 },
    {
      title: '操作', width: 280, render: (_: any, r: Team) => (
        <Space>
          {isAdmin && <Button icon={<EditOutlined />} size="small" onClick={() => handleEditTeam(r)}>编辑</Button>}
          {isAdmin && <Button icon={<UserAddOutlined />} size="small" onClick={() => handleAddPlayer(r.id)}>添加球员</Button>}
          <Button size="small" onClick={() => handleRegisterTeam(r.id)}>报名参赛</Button>
          {isAdmin && <Popconfirm title="确认删除?" onConfirm={() => handleDeleteTeam(r.id)}>
            <Button danger icon={<DeleteOutlined />} size="small">删除</Button>
          </Popconfirm>}
        </Space>
      )
    }
  ]

  const expandedRowRender = (team: Team) => (
    <Card size="small" title={`球员名单 (${team.players?.length || 0})`}>
      {team.players && team.players.length > 0 ? (
        <Table
          size="small"
          dataSource={team.players}
          rowKey="id"
          pagination={false}
          columns={[
            { title: '号码', dataIndex: 'number', width: 80 },
            { title: '姓名', dataIndex: 'name' },
            { title: '位置', dataIndex: 'position' },
            { title: '出生日期', dataIndex: 'birth_date', render: (d: string) => d ? dayjs(d).format('YYYY-MM-DD') : '-' },
            { title: '状态', dataIndex: 'is_active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? '活跃' : '停赛'}</Tag> },
          ]}
        />
      ) : <p style={{ color: '#999' }}>暂无球员</p>}
    </Card>
  )

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>球队管理</h2>
        {isAdmin && <Button type="primary" icon={<PlusOutlined />} onClick={handleCreateTeam}>创建球队</Button>}
      </div>
      <Table
        rowKey="id"
        dataSource={teams}
        columns={columns}
        expandable={{ expandedRowRender }}
      />
      {isAdmin && registrations.length > 0 && (
        <Card title="报名审批" style={{ marginTop: 24 }}>
          <Table
            size="small"
            dataSource={registrations}
            rowKey="id"
            pagination={false}
            columns={[
              { title: 'ID', dataIndex: 'id', width: 60 },
              { title: '球队', dataIndex: ['team', 'name'] },
              { title: '赛季', dataIndex: ['season', 'name'] },
              { title: '分组', dataIndex: 'group_name' },
              { title: '状态', dataIndex: 'status', render: (s: string) => <Tag color={s === 'approved' ? 'green' : s === 'pending' ? 'orange' : 'red'}>{s}</Tag> },
              { title: '已付费', dataIndex: 'paid', render: (v: boolean) => v ? '是' : '否' },
              {
                title: '操作', render: (_: any, r: Registration) => r.status === 'pending' ? (
                  <Button type="primary" size="small" onClick={() => handleApproveReg(r.id)}>批准</Button>
                ) : null
              }
            ]}
          />
        </Card>
      )}
      <Modal
        title={editingTeam ? '编辑球队' : '创建球队'}
        open={teamModalOpen}
        onCancel={() => setTeamModalOpen(false)}
        onOk={handleTeamSubmit}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="球队名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="logo" label="Logo URL">
            <Input />
          </Form.Item>
          <Form.Item name="description" label="简介">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Form.Item name="contact_email" label="联系邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="contact_phone" label="联系电话">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="添加球员"
        open={playerModalOpen}
        onCancel={() => setPlayerModalOpen(false)}
        onOk={handlePlayerSubmit}
      >
        <Form form={playerForm} layout="vertical">
          <Form.Item name="name" label="姓名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="number" label="号码" rules={[{ required: true }]}>
            <InputNumber min={1} max={99} />
          </Form.Item>
          <Form.Item name="position" label="位置">
            <Select options={[
              { value: '前锋', label: '前锋' },
              { value: '中场', label: '中场' },
              { value: '后卫', label: '后卫' },
              { value: '门将', label: '门将' },
              { value: '替补', label: '替补' },
            ]} />
          </Form.Item>
          <Form.Item name="birth_date" label="出生日期">
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="球队报名"
        open={regModalOpen}
        onCancel={() => setRegModalOpen(false)}
        onOk={handleRegSubmit}
      >
        <Form form={regForm} layout="vertical">
          <Form.Item name="season_id" label="选择赛季" rules={[{ required: true }]}>
            <Select options={seasons.map((s) => ({ value: s.id, label: s.name }))} />
          </Form.Item>
          <Form.Item name="group_name" label="分组">
            <Input placeholder="如 A, B, C..." />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
