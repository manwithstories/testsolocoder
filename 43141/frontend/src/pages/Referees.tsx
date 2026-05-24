import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, Select, Space, message, Tag, Card } from 'antd'
import { CheckOutlined, CloseOutlined } from '@ant-design/icons'
import api from '../api'
import type { RefereeAssignment, Match, User } from '../types'

export default function Referees({ isAdmin }: { isAdmin: boolean }) {
  const [assignments, setAssignments] = useState<RefereeAssignment[]>([])
  const [referees, setReferees] = useState<User[]>([])
  const [matches, setMatches] = useState<Match[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchData = () => {
    api.get('/referees/assignments').then((res) => setAssignments(res.data))
    api.get('/referees').then((res) => setReferees(res.data))
  }

  useEffect(() => { fetchData() }, [])

  const handleAssign = () => {
    form.resetFields()
    api.get('/matches').then((res) => setMatches(res.data))
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      await api.post('/referees/assign', values)
      message.success('指派成功')
      setModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '指派失败')
    }
  }

  const handleRespond = async (id: number, action: string) => {
    try {
      await api.put(`/referees/assignments/${id}?action=${action}`)
      message.success(action === 'accept' ? '已接受' : '已拒绝')
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '操作失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '比赛', dataIndex: ['match', 'id'], render: (v: number, r: any) => {
      const m = r.match
      if (!m) return '-'
      return `${m.home_team?.name || ''} vs ${m.away_team?.name || ''}`
    }},
    { title: '裁判', dataIndex: ['referee', 'full_name'] },
    {
      title: '状态', dataIndex: 'status', render: (s: string) => {
        const colors: Record<string, string> = { assigned: 'blue', accepted: 'green', rejected: 'red' }
        return <Tag color={colors[s] || 'default'}>{s}</Tag>
      }
    },
    { title: '指派时间', dataIndex: 'assigned_at', render: (v: string) => v ? new Date(v).toLocaleString() : '-' },
    {
      title: '操作', render: (_: any, r: RefereeAssignment) => {
        if (r.status === 'assigned') {
          return (
            <Space>
              <Button icon={<CheckOutlined />} type="primary" size="small" onClick={() => handleRespond(r.id, 'accept')}>接受</Button>
              <Button icon={<CloseOutlined />} danger size="small" onClick={() => handleRespond(r.id, 'reject')}>拒绝</Button>
            </Space>
          )
        }
        return null
      }
    }
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>裁判管理</h2>
        {isAdmin && <Button type="primary" onClick={handleAssign}>指派裁判</Button>}
      </div>
      <Card title="执法任务">
        <Table rowKey="id" dataSource={assignments} columns={columns} />
      </Card>
      <Card title="裁判列表" style={{ marginTop: 16 }}>
        <Table
          size="small"
          rowKey="id"
          dataSource={referees}
          pagination={false}
          columns={[
            { title: 'ID', dataIndex: 'id', width: 60 },
            { title: '姓名', dataIndex: 'full_name' },
            { title: '邮箱', dataIndex: 'email' },
            { title: '状态', dataIndex: 'is_active', render: (v: boolean) => <Tag color={v ? 'green' : 'red'}>{v ? '活跃' : '停用'}</Tag> },
          ]}
        />
      </Card>

      <Modal
        title="指派裁判"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="match_id" label="选择比赛" rules={[{ required: true }]}>
            <Select options={matches.map((m) => ({
              value: m.id,
              label: `${m.home_team?.name || ''} vs ${m.away_team?.name || ''}`
            }))} showSearch optionFilterProp="label" />
          </Form.Item>
          <Form.Item name="referee_id" label="选择裁判" rules={[{ required: true }]}>
            <Select options={referees.map((r) => ({ value: r.id, label: r.full_name }))} showSearch optionFilterProp="label" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
