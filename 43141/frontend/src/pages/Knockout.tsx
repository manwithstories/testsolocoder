import { useState, useEffect } from 'react'
import { Table, Button, Modal, Form, Select, Space, message, Card, Tag, Tree } from 'antd'
import { PlusOutlined, CrownOutlined } from '@ant-design/icons'
import api from '../api'
import type { Season, Match, Team } from '../types'

export default function Knockout({ isAdmin }: { isAdmin: boolean }) {
  const [seasons, setSeasons] = useState<Season[]>([])
  const [matches, setMatches] = useState<Match[]>([])
  const [teams, setTeams] = useState<Team[]>([])
  const [modalOpen, setModalOpen] = useState(false)
  const [currentSeasonId, setCurrentSeasonId] = useState<number | null>(null)
  const [form] = Form.useForm()

  const fetchData = () => {
    api.get('/leagues').then((res) => {
      const all: Season[] = []
      res.data.forEach((l: any) => l.seasons?.forEach((s: Season) => all.push(s)))
      setSeasons(all)
    })
    api.get('/teams').then((res) => setTeams(res.data))
  }

  const fetchMatches = (seasonId: number) => {
    api.get(`/matches?season_id=${seasonId}`).then((res) => {
      const knockout = res.data.filter((m: Match) => m.knockout_stage)
      setMatches(knockout)
    })
  }

  useEffect(() => { fetchData() }, [])

  const handleGenerate = () => {
    form.resetFields()
    setModalOpen(true)
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      await api.post(`/seasons/${currentSeasonId}/generate-knockout`, values)
      message.success('淘汰赛生成成功')
      setModalOpen(false)
      if (currentSeasonId) fetchMatches(currentSeasonId)
    } catch (e: any) {
      message.error(e.response?.data?.error || '生成失败')
    }
  }

  const stages = ['Round of 16', 'Quarter-final', 'Semi-final', 'Final']

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2 style={{ margin: 0 }}>淘汰赛对阵</h2>
        {isAdmin && (
          <Space>
            <Select
              style={{ width: 200 }}
              placeholder="选择赛季"
              value={currentSeasonId}
              onChange={(v) => { setCurrentSeasonId(v); fetchMatches(v) }}
              options={seasons.map((s) => ({ value: s.id, label: s.name }))}
            />
            <Button type="primary" icon={<CrownOutlined />} onClick={handleGenerate} disabled={!currentSeasonId}>
              生成淘汰赛
            </Button>
          </Space>
        )}
      </div>

      {currentSeasonId && (
        <Card title="淘汰赛对阵图">
          {stages.map((stage) => {
            const stageMatches = matches.filter((m) => m.knockout_stage === stage)
            if (stageMatches.length === 0) return null
            return (
              <div key={stage} style={{ marginBottom: 24 }}>
                <h3 style={{ borderBottom: '2px solid #1890ff', paddingBottom: 8 }}>{stage}</h3>
                <Table
                  size="small"
                  rowKey="id"
                  dataSource={stageMatches}
                  pagination={false}
                  columns={[
                    { title: '主队', dataIndex: ['home_team', 'name'], render: (v: string) => <strong>{v}</strong> },
                    {
                      title: '比分', width: 100, render: (_: any, r: Match) => {
                        if (r.home_score != null && r.away_score != null) {
                          return <Tag color="blue">{r.home_score} - {r.away_score}</Tag>
                        }
                        return <Tag>VS</Tag>
                      }
                    },
                    { title: '客队', dataIndex: ['away_team', 'name'] },
                    { title: '胜者', dataIndex: 'winner_team_id', render: (v: number) => v ? teams.find((t) => t.id === v)?.name : '-' },
                    { title: '场地', dataIndex: ['venue', 'name'], render: (v: string) => v || '-' },
                  ]}
                />
              </div>
            )
          })}
        </Card>
      )}

      <Modal
        title="生成淘汰赛"
        open={modalOpen}
        onCancel={() => setModalOpen(false)}
        onOk={handleSubmit}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="format" label="赛制" rules={[{ required: true }]}>
            <Select options={[
              { value: 'single_elimination', label: '单败淘汰' },
              { value: 'double_elimination', label: '双败淘汰' },
              { value: 'group_knockout', label: '小组赛+淘汰赛' },
            ]} />
          </Form.Item>
          <Form.Item name="team_ids" label="参赛球队" rules={[{ required: true }]}>
            <Select mode="multiple" options={teams.map((t) => ({ value: t.id, label: t.name }))} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
