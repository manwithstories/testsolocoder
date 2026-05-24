import { useState, useEffect } from 'react'
import { Table, Select, Card } from 'antd'
import api from '../api'
import type { Season, Standings as StandingsType } from '../types'

export default function Standings() {
  const [seasons, setSeasons] = useState<Season[]>([])
  const [selectedSeasonId, setSelectedSeasonId] = useState<number | null>(null)
  const [standings, setStandings] = useState<StandingsType[]>([])

  const fetchSeasons = () => {
    api.get('/leagues').then((res) => {
      const all: Season[] = []
      res.data.forEach((l: any) => l.seasons?.forEach((s: Season) => all.push(s)))
      setSeasons(all)
    })
  }

  const fetchStandings = (seasonId: number) => {
    api.get(`/seasons/${seasonId}/standings`).then((res) => setStandings(res.data))
  }

  useEffect(() => { fetchSeasons() }, [])

  const columns = [
    { title: '排名', key: 'rank', render: (_: any, __: any, i: number) => i + 1, width: 60 },
    { title: '球队', dataIndex: 'team_name' },
    { title: '分组', dataIndex: 'group_name', width: 80 },
    { title: '场次', dataIndex: 'played', width: 60 },
    { title: '胜', dataIndex: 'wins', width: 60 },
    { title: '平', dataIndex: 'draws', width: 60 },
    { title: '负', dataIndex: 'losses', width: 60 },
    { title: '进球', dataIndex: 'goals_for', width: 60 },
    { title: '失球', dataIndex: 'goals_against', width: 60 },
    { title: '净胜', dataIndex: 'goal_diff', width: 60, render: (v: number) => v > 0 ? `+${v}` : v },
    { title: '积分', dataIndex: 'points', width: 60, render: (v: number) => <strong>{v}</strong> },
  ]

  return (
    <div>
      <h2 style={{ marginTop: 0 }}>积分榜</h2>
      <Card size="small" style={{ marginBottom: 16 }}>
        <label style={{ marginRight: 8 }}>选择赛季:</label>
        <Select
          style={{ width: 250 }}
          placeholder="选择赛季"
          value={selectedSeasonId}
          onChange={(v) => { setSelectedSeasonId(v); fetchStandings(v) }}
          options={seasons.map((s) => ({ value: s.id, label: s.name }))}
        />
      </Card>
      {selectedSeasonId && (
        <Table
          rowKey="team_id"
          dataSource={standings}
          columns={columns}
          pagination={false}
        />
      )}
    </div>
  )
}
