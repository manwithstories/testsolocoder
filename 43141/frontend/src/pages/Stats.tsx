import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, InputNumber, Select, Space, Popconfirm, message, Card } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import api from '../api'
import type { PlayerRanking, Season, Team, Player, Match } from '../types'

export default function Stats({ isAdmin }: { isAdmin: boolean }) {
  const [rankings, setRankings] = useState<PlayerRanking[]>([])
  const [seasons, setSeasons] = useState<Season[]>([])
  const [teams, setTeams] = useState<Team[]>([])
  const [players, setPlayers] = useState<Player[]>([])
  const [matches, setMatches] = useState<Match[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [form] = Form.useForm()

  const fetchData = () => {
    api.get('/stats/rankings').then((res) => setRankings(res.data))
    api.get('/leagues').then((res) => {
      const all: Season[] = []
      res.data.forEach((l: any) => l.seasons?.forEach((s: Season) => all.push(s)))
      setSeasons(all)
    })
    api.get('/teams').then((res) => setTeams(res.data))
  }

  useEffect(() => { fetchData() }, [])

  const handleAddStat = () => {
    form.resetFields()
    setModalOpen(true)
  }

  const fetchPlayers = (teamId: number) => {
    api.get(`/teams/${teamId}`).then((res) => setPlayers(res.data.players || []))
  }

  const fetchMatches = (seasonId: number) => {
    api.get(`/matches?season_id=${seasonId}`).then((res) => setMatches(res.data))
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      await api.post('/stats', values)
      message.success('数据录入成功')
      setModalOpen(false)
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '录入失败')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/stats/${id}`)
      message.success('删除成功')
      fetchData()
    } catch (e: any) {
      message.error(e.response?.data?.error || '删除失败')
    }
  }

  const columns = [
    { title: '排名', key: 'rank', render: (_: any, __: any, i: number) => i + 1, width: 60 },
    { title: '球员', dataIndex: 'player_name' },
    { title: '球队', dataIndex: 'team_name' },
    { title: '场次', dataIndex: 'matches', width: 60 },
    { title: '进球', dataIndex: 'goals', width: 60, render: (v: number) => <strong style={{ color: v > 0 ? '#52c41a' : undefined }}>{v}</strong> },
    { title: '助攻', dataIndex: 'assists', width: 60 },
    { title: '犯规', dataIndex: 'fouls', width: 60 },
    { title: '黄牌', dataIndex: 'yellow_card', width: 60, render: (v: number) => v > 0 ? <span style={{ color: '#faad14' }}>{v}</span> : v },
    { title: '红牌', dataIndex: 'red_card', width: 60, render: (v: number) => v > 0 ? <span style={{ color: '#ff4d4f' }}>{v}</span> : v },
    { title: '出场分钟', dataIndex: 'minutes', width: 80 },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>球员数据统计</h2>
        {isAdmin && <Button type="primary" icon={<PlusOutlined />} onClick={handleAddStat}>录入数据</Button>}
      </div>
      <Card title="球员排行榜">
        <Table
          rowKey="player_id"
          dataSource={rankings}
          columns={columns}
        />
      </Card>

      <Modal
        title="录入球员数据"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="season_id" label="赛季" rules={[{ required: true }]}>
            <Select
              options={seasons.map((s) => ({ value: s.id, label: s.name }))}
              onChange={(v) => fetchMatches(v)}
            />
          </Form.Item>
          <Form.Item name="team_id" label="球队" rules={[{ required: true }]}>
            <Select
              options={teams.map((t) => ({ value: t.id, label: t.name }))}
              onChange={(v) => fetchPlayers(v)}
            />
          </Form.Item>
          <Form.Item name="match_id" label="比赛" rules={[{ required: true }]}>
            <Select options={matches.map((m) => ({ value: m.id, label: `${m.home_team?.name || ''} vs ${m.away_team?.name || ''}` }))} />
          </Form.Item>
          <Form.Item name="player_id" label="球员" rules={[{ required: true }]}>
            <Select options={players.map((p) => ({ value: p.id, label: `${p.number} - ${p.name}` }))} />
          </Form.Item>
          <Form.Item name="goals" label="进球" initialValue={0}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="assists" label="助攻" initialValue={0}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="fouls" label="犯规" initialValue={0}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="yellow_card" label="黄牌" initialValue={0}>
            <InputNumber min={0} max={2} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="red_card" label="红牌" initialValue={0}>
            <InputNumber min={0} max={1} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="minutes" label="出场分钟" initialValue={0}>
            <InputNumber min={0} max={120} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
