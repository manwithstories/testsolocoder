import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, DatePicker, message, Tag, Row, Col, Card, Tooltip, Checkbox } from 'antd'
import { PlusOutlined, EditOutlined, CalendarOutlined, PlayCircleOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import api from '../api'
import type { Match, Season, Venue } from '../types'

export default function Matches({ isAdmin }: { isAdmin: boolean }) {
  const [matches, setMatches] = useState<Match[]>([])
  const [seasons, setSeasons] = useState<Season[]>([])
  const [venues, setVenues] = useState<Venue[]>([])
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [scoreModalOpen, setScoreModalOpen] = useState(false)
  const [scheduleModalOpen, setScheduleModalOpen] = useState(false)
  const [venuesModalOpen, setVenuesModalOpen] = useState(false)
  const [editingMatch, setEditingMatch] = useState<Match | null>(null)
  const [currentSeasonId, setCurrentSeasonId] = useState<number | null>(null)
  const [selectedSeasonId, setSelectedSeasonId] = useState<number | null>(null)
  const [form] = Form.useForm()
  const [scoreForm] = Form.useForm()
  const [scheduleForm] = Form.useForm()
  const [venueForm] = Form.useForm()

  const fetchData = () => {
    api.get('/leagues').then((res) => {
      const all: Season[] = []
      res.data.forEach((l: any) => l.seasons?.forEach((s: Season) => all.push(s)))
      setSeasons(all)
    })
    api.get('/venues').then((res) => setVenues(res.data))
    if (selectedSeasonId) {
      api.get(`/matches?season_id=${selectedSeasonId}`).then((res) => setMatches(res.data))
    }
  }

  useEffect(() => { fetchData() }, [selectedSeasonId])

  const handleEdit = (match: Match) => {
    setEditingMatch(match)
    form.setFieldsValue({
      ...match,
      match_time: match.match_time ? dayjs(match.match_time) : null,
      venue_id: match.venue_id
    })
    setEditModalOpen(true)
  }

  const handleEditSubmit = async () => {
    try {
      const values = await form.validateFields()
      const payload: any = { ...values }
      if (values.match_time) payload.match_time = dayjs(values.match_time).format()
      await api.put(`/matches/${editingMatch?.id}`, payload)
      message.success('更新成功')
      setEditModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '更新失败')
    }
  }

  const handleReportScore = (match: Match) => {
    setEditingMatch(match)
    scoreForm.resetFields()
    setScoreModalOpen(true)
  }

  const handleScoreSubmit = async () => {
    try {
      const values = await scoreForm.validateFields()
      await api.post(`/matches/${editingMatch?.id}/report-score`, values)
      message.success('比分上报成功')
      setScoreModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '上报失败')
    }
  }

  const handleGenerateSchedule = () => {
    scheduleForm.resetFields()
    setScheduleModalOpen(true)
  }

  const handleScheduleSubmit = async () => {
    try {
      const values = await scheduleForm.validateFields()
      const payload: any = { ...values }
      if (values.start_date) payload.start_date = dayjs(values.start_date).format('YYYY-MM-DD')
      await api.post(`/seasons/${currentSeasonId}/generate-schedule`, payload)
      message.success('赛程生成成功')
      setScheduleModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '生成失败')
    }
  }

  const handleCreateVenue = () => {
    venueForm.resetFields()
    setVenuesModalOpen(true)
  }

  const handleVenueSubmit = async () => {
    try {
      const values = await venueForm.validateFields()
      await api.post('/venues', values)
      message.success('场地创建成功')
      setVenuesModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '创建失败')
    }
  }

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 60 },
    { title: '轮次', dataIndex: 'round', width: 60 },
    { title: '淘汰赛阶段', dataIndex: 'knockout_stage', width: 120, render: (v: string) => v || '-' },
    { title: '分组', dataIndex: 'group_name', width: 80 },
    {
      title: '主队', dataIndex: ['home_team', 'name'], render: (v: string) => <strong>{v}</strong>
    },
    {
      title: '比分', width: 100, render: (_: any, r: Match) => {
        if (r.home_score != null && r.away_score != null) {
          let score = `${r.home_score} - ${r.away_score}`
          if (r.has_ot) score += ' (OT)'
          if (r.has_penalty) score += ' (PEN)'
          return <Tag color="blue">{score}</Tag>
        }
        return <Tag>VS</Tag>
      }
    },
    { title: '客队', dataIndex: ['away_team', 'name'] },
    { title: '场地', dataIndex: ['venue', 'name'], render: (v: string) => v || '-' },
    { title: '裁判', dataIndex: ['referee', 'full_name'], render: (v: string) => v || '-' },
    {
      title: '比赛时间', dataIndex: 'match_time', render: (v: string) => v ? dayjs(v).format('YYYY-MM-DD HH:mm') : '-'
    },
    {
      title: '状态', dataIndex: 'status', render: (s: string) => {
        const colors: Record<string, string> = { scheduled: 'blue', 'finished': 'green', 'cancelled': 'red', in_progress: 'orange' }
        return <Tag color={colors[s] || 'default'}>{s}</Tag>
      }
    },
    {
      title: '操作', width: 200, render: (_: any, r: Match) => (
        <Space>
          {isAdmin && <Button icon={<EditOutlined />} size="small" onClick={() => handleEdit(r)}>编辑</Button>}
          {isAdmin && r.status === 'scheduled' && (
            <Button type="primary" icon={<PlayCircleOutlined />} size="small" onClick={() => handleReportScore(r)}>上报比分</Button>
          )}
        </Space>
      )
    }
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>比赛赛程</h2>
        {isAdmin && (
          <Space>
            <Button icon={<PlusOutlined />} onClick={handleCreateVenue}>添加场地</Button>
            <Button icon={<CalendarOutlined />} onClick={handleGenerateSchedule} disabled={!currentSeasonId}>生成赛程</Button>
          </Space>
        )}
      </div>
      <Card size="small" style={{ marginBottom: 16 }}>
        <Row gutter={16}>
          <Col span={8}>
            <label style={{ marginRight: 8 }}>选择赛季:</label>
            <Select
              style={{ width: 200 }}
              placeholder="选择赛季"
              value={selectedSeasonId}
              onChange={(v) => setSelectedSeasonId(v)}
              options={seasons.map((s) => ({ value: s.id, label: s.name }))}
            />
          </Col>
          <Col span={8}>
            <label style={{ marginRight: 8 }}>生成赛程赛季:</label>
            <Select
              style={{ width: 200 }}
              placeholder="选择赛季"
              value={currentSeasonId}
              onChange={(v) => setCurrentSeasonId(v)}
              options={seasons.map((s) => ({ value: s.id, label: s.name }))}
            />
          </Col>
        </Row>
      </Card>
      <Table
        rowKey="id"
        dataSource={matches}
        columns={columns}
      />
      <Modal
        title="编辑比赛"
        open={editModalOpen}
        onCancel={() => setEditModalOpen(false)}
        onOk={handleEditSubmit}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="match_time" label="比赛时间">
            <DatePicker showTime style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="venue_id" label="场地">
            <Select options={venues.map((v) => ({ value: v.id, label: v.name }))} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select options={[
              { value: 'scheduled', label: '已安排' },
              { value: 'in_progress', label: '进行中' },
              { value: 'finished', label: '已结束' },
              { value: 'cancelled', label: '已取消' },
            ]} />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="上报比分"
        open={scoreModalOpen}
        onCancel={() => setScoreModalOpen(false)}
        onOk={handleScoreSubmit}
      >
        <Form form={scoreForm} layout="vertical">
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="home_score" label="主队得分" rules={[{ required: true }]}>
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="away_score" label="客队得分" rules={[{ required: true }]}>
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="has_ot" label="包含加时赛" valuePropName="checked">
            <Checkbox />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="ot_home_score" label="加时主队">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="ot_away_score" label="加时客队">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="has_penalty" label="包含点球大战" valuePropName="checked">
            <Checkbox />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="pen_home_score" label="点球主队">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="pen_away_score" label="点球客队">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>
      <Modal
        title="生成赛程"
        open={scheduleModalOpen}
        onCancel={() => setScheduleModalOpen(false)}
        onOk={handleScheduleSubmit}
      >
        <Form form={scheduleForm} layout="vertical">
          <Form.Item name="start_date" label="开始日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="interval" label="比赛间隔">
            <Select options={[
              { value: 'daily', label: '每天' },
              { value: '3days', label: '每3天' },
              { value: 'weekly', label: '每周' },
            ]} />
          </Form.Item>
          <Form.Item name="venue_ids" label="场地 (可选)">
            <Select mode="multiple" options={venues.map((v) => ({ value: v.id, label: v.name }))} />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        title="添加场地"
        open={venuesModalOpen}
        onCancel={() => setVenuesModalOpen(false)}
        onOk={handleVenueSubmit}
      >
        <Form form={venueForm} layout="vertical">
          <Form.Item name="name" label="场地名称" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="address" label="地址">
            <Input />
          </Form.Item>
          <Form.Item name="capacity" label="容量">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
