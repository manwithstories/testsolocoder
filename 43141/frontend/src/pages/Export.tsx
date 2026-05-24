import { useState, useEffect } from 'react'
import { Button, Select, Card, message, Row, Col } from 'antd'
import { DownloadOutlined, FileExcelOutlined, FilePdfOutlined } from '@ant-design/icons'
import api from '../api'
import type { Season, Match, Standings, PlayerRanking, League } from '../types'

async function downloadFile(url: string, filename: string) {
  const token = localStorage.getItem('token')
  const res = await fetch(`/api${url}`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || `HTTP ${res.status}`)
  }
  const blob = await res.blob()
  const downloadUrl = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = downloadUrl
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  window.URL.revokeObjectURL(downloadUrl)
}

function buildPdfHtml(
  season: Season,
  league: League | undefined,
  matches: Match[],
  standings: Standings[],
  rankings: PlayerRanking[]
): string {
  const finishedMatches = matches.filter((m) => m.status === 'finished')
  const scheduledMatches = matches.filter((m) => m.status === 'scheduled')

  const matchesRows = matches.map((m) => {
    const home = m.home_team?.name || '-'
    const away = m.away_team?.name || '-'
    const venue = m.venue?.name || '-'
    const time = m.match_time ? new Date(m.match_time).toLocaleString() : '-'
    const score = m.home_score != null ? `${m.home_score} - ${m.away_score}` : 'VS'
    return `<tr>
      <td>${m.round}</td>
      <td>${m.group_name || m.knockout_stage || '-'}</td>
      <td>${home}</td>
      <td>${score}</td>
      <td>${away}</td>
      <td>${venue}</td>
      <td>${time}</td>
      <td>${m.status}</td>
    </tr>`
  }).join('')

  const standingsRows = standings.map((s, i) => `
    <tr>
      <td>${i + 1}</td>
      <td>${s.team_name}</td>
      <td>${s.group_name}</td>
      <td>${s.played}</td>
      <td>${s.wins}</td>
      <td>${s.draws}</td>
      <td>${s.losses}</td>
      <td>${s.goals_for}</td>
      <td>${s.goals_against}</td>
      <td>${s.goal_diff > 0 ? '+' + s.goal_diff : s.goal_diff}</td>
      <td><strong>${s.points}</strong></td>
    </tr>
  `).join('')

  const rankingRows = rankings.slice(0, 20).map((r, i) => `
    <tr>
      <td>${i + 1}</td>
      <td>${r.player_name}</td>
      <td>${r.team_name}</td>
      <td>${r.matches}</td>
      <td>${r.goals}</td>
      <td>${r.assists}</td>
      <td>${r.fouls}</td>
      <td>${r.yellow_card}</td>
      <td>${r.red_card}</td>
      <td>${r.minutes}</td>
    </tr>
  `).join('')

  return `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>赛季报告</title>
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Microsoft YaHei', sans-serif; margin: 40px; color: #333; }
  h1 { color: #1890ff; border-bottom: 3px solid #1890ff; padding-bottom: 10px; }
  h2 { color: #333; margin-top: 30px; border-left: 4px solid #1890ff; padding-left: 10px; }
  table { width: 100%; border-collapse: collapse; margin-top: 10px; font-size: 13px; }
  th { background: #f5f7fa; color: #333; padding: 8px; text-align: left; border-bottom: 2px solid #d9d9d9; }
  td { padding: 6px 8px; border-bottom: 1px solid #f0f0f0; }
  tr:hover td { background: #fafafa; }
  .meta { background: #f5f7fa; padding: 16px; border-radius: 6px; margin: 16px 0; }
  .meta p { margin: 4px 0; }
  .stats-grid { display: flex; gap: 16px; margin: 16px 0; }
  .stat-card { flex: 1; background: #e6f7ff; border-radius: 6px; padding: 16px; text-align: center; }
  .stat-card .num { font-size: 28px; font-weight: bold; color: #1890ff; }
  .stat-card .label { font-size: 13px; color: #666; margin-top: 4px; }
</style></head><body>
<h1>🏆 赛季综合报告</h1>
<div class="meta">
  <p><strong>联赛:</strong> ${league?.name || '-'}</p>
  <p><strong>赛季:</strong> ${season.name}</p>
  <p><strong>赛制:</strong> ${season.format}</p>
  <p><strong>时间:</strong> ${new Date(season.start_date).toLocaleDateString()} ~ ${new Date(season.end_date).toLocaleDateString()}</p>
  <p><strong>生成时间:</strong> ${new Date().toLocaleString()}</p>
</div>

<div class="stats-grid">
  <div class="stat-card"><div class="num">${matches.length}</div><div class="label">总比赛数</div></div>
  <div class="stat-card"><div class="num">${finishedMatches.length}</div><div class="label">已结束</div></div>
  <div class="stat-card"><div class="num">${scheduledMatches.length}</div><div class="label">待进行</div></div>
  <div class="stat-card"><div class="num">${standings.length}</div><div class="label">参赛队伍</div></div>
</div>

<h2>📊 积分榜</h2>
<table>
<thead><tr><th>#</th><th>球队</th><th>分组</th><th>场次</th><th>胜</th><th>平</th><th>负</th><th>进</th><th>失</th><th>净</th><th>积分</th></tr></thead>
<tbody>${standingsRows || '<tr><td colspan="11" style="text-align:center;color:#999">暂无数据</td></tr>'}</tbody>
</table>

<h2>📅 赛程表</h2>
<table>
<thead><tr><th>轮次</th><th>阶段</th><th>主队</th><th>比分</th><th>客队</th><th>场地</th><th>时间</th><th>状态</th></tr></thead>
<tbody>${matchesRows || '<tr><td colspan="8" style="text-align:center;color:#999">暂无数据</td></tr>'}</tbody>
</table>

<h2>⚽ 球员排行榜 TOP 20</h2>
<table>
<thead><tr><th>#</th><th>球员</th><th>球队</th><th>场次</th><th>进球</th><th>助攻</th><th>犯规</th><th>黄牌</th><th>红牌</th><th>分钟</th></tr></thead>
<tbody>${rankingRows || '<tr><td colspan="10" style="text-align:center;color:#999">暂无数据</td></tr>'}</tbody>
</table>

<p style="margin-top:40px;color:#999;text-align:center;font-size:12px">本报告由体育联赛管理平台自动生成</p>
</body></html>`
}

export default function Export() {
  const [seasons, setSeasons] = useState<Season[]>([])
  const [leagues, setLeagues] = useState<League[]>([])
  const [selectedSeasonId, setSelectedSeasonId] = useState<number | null>(null)
  const [loading, setLoading] = useState<string | null>(null)

  useEffect(() => {
    api.get('/leagues').then((res) => {
      setLeagues(res.data)
      const all: Season[] = []
      res.data.forEach((l: League) => l.seasons?.forEach((s: Season) => all.push({ ...s, league_id: l.id })))
      setSeasons(all)
    })
  }, [])

  const selectedSeason = seasons.find((s) => s.id === selectedSeasonId)
  const selectedLeague = leagues.find((l) => l.id === selectedSeason?.league_id)

  const handleExportExcel = async (type: string) => {
    if (!selectedSeasonId) {
      message.warning('请先选择赛季')
      return
    }
    setLoading(type)
    try {
      let url = ''
      let filename = ''
      if (type === 'schedule') {
        url = `/export/schedule/${selectedSeasonId}`
        filename = `赛程表_赛季${selectedSeasonId}.xlsx`
      } else if (type === 'standings') {
        url = `/export/standings/${selectedSeasonId}`
        filename = `积分榜_赛季${selectedSeasonId}.xlsx`
      } else if (type === 'stats') {
        url = `/export/stats?season_id=${selectedSeasonId}`
        filename = `球员数据_赛季${selectedSeasonId}.xlsx`
      }
      await downloadFile(url, filename)
      message.success('导出成功')
    } catch (e: any) {
      message.error(`导出失败: ${e.message}`)
    } finally {
      setLoading(null)
    }
  }

  const handleExportPDF = async () => {
    if (!selectedSeasonId || !selectedSeason) {
      message.warning('请先选择赛季')
      return
    }
    setLoading('pdf')
    try {
      const [matchesRes, standingsRes, statsRes] = await Promise.all([
        api.get(`/matches?season_id=${selectedSeasonId}`),
        api.get(`/seasons/${selectedSeasonId}/standings`),
        api.get(`/stats/rankings?season_id=${selectedSeasonId}`),
      ])

      const html = buildPdfHtml(
        selectedSeason,
        selectedLeague,
        matchesRes.data || [],
        standingsRes.data || [],
        statsRes.data || []
      )

      const printWindow = window.open('', '_blank')
      if (!printWindow) {
        message.error('请允许弹出窗口以查看PDF报告')
        setLoading(null)
        return
      }
      printWindow.document.write(html)
      printWindow.document.close()
      printWindow.focus()
      setTimeout(() => {
        printWindow.print()
      }, 500)

      message.success('报告已生成，请在新窗口中打印或保存为PDF')
    } catch (e: any) {
      message.error(`生成失败: ${e.response?.data?.error || e.message}`)
    } finally {
      setLoading(null)
    }
  }

  return (
    <div>
      <h2 style={{ marginTop: 0 }}>数据导出</h2>
      <Card>
        <div style={{ marginBottom: 24 }}>
          <label style={{ marginRight: 8 }}>选择赛季:</label>
          <Select
            style={{ width: 250 }}
            placeholder="选择赛季"
            value={selectedSeasonId}
            onChange={setSelectedSeasonId}
            options={seasons.map((s) => ({ value: s.id, label: s.name }))}
          />
        </div>

        <Row gutter={16}>
          <Col span={6}>
            <Card
              hoverable
              style={{ textAlign: 'center' }}
              actions={[
                <Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  loading={loading === 'schedule'}
                  onClick={() => handleExportExcel('schedule')}
                  disabled={!selectedSeasonId}
                >
                  导出赛程表
                </Button>
              ]}
            >
              <FileExcelOutlined style={{ fontSize: 48, color: '#1890ff' }} />
              <p style={{ marginTop: 12 }}>赛程表 (Excel)</p>
            </Card>
          </Col>
          <Col span={6}>
            <Card
              hoverable
              style={{ textAlign: 'center' }}
              actions={[
                <Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  loading={loading === 'standings'}
                  onClick={() => handleExportExcel('standings')}
                  disabled={!selectedSeasonId}
                >
                  导出积分榜
                </Button>
              ]}
            >
              <FileExcelOutlined style={{ fontSize: 48, color: '#52c41a' }} />
              <p style={{ marginTop: 12 }}>积分榜 (Excel)</p>
            </Card>
          </Col>
          <Col span={6}>
            <Card
              hoverable
              style={{ textAlign: 'center' }}
              actions={[
                <Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  loading={loading === 'stats'}
                  onClick={() => handleExportExcel('stats')}
                  disabled={!selectedSeasonId}
                >
                  导出球员数据
                </Button>
              ]}
            >
              <FileExcelOutlined style={{ fontSize: 48, color: '#faad14' }} />
              <p style={{ marginTop: 12 }}>球员数据 (Excel)</p>
            </Card>
          </Col>
          <Col span={6}>
            <Card
              hoverable
              style={{ textAlign: 'center' }}
              actions={[
                <Button
                  type="primary"
                  icon={<DownloadOutlined />}
                  loading={loading === 'pdf'}
                  onClick={handleExportPDF}
                  disabled={!selectedSeasonId}
                >
                  生成PDF报告
                </Button>
              ]}
            >
              <FilePdfOutlined style={{ fontSize: 48, color: '#ff4d4f' }} />
              <p style={{ marginTop: 12 }}>综合报告 (PDF)</p>
            </Card>
          </Col>
        </Row>

        <Card type="inner" title="导出说明" style={{ marginTop: 24 }}>
          <ul>
            <li>赛程表导出: 包含所有比赛的轮次、分组、对阵双方、场地、时间和状态 (Excel格式)</li>
            <li>积分榜导出: 包含各球队的场次、胜负平、进失球、净胜球、积分 (Excel格式)</li>
            <li>球员数据导出: 包含球员进球、助攻、犯规、黄牌、红牌等统计 (Excel格式)</li>
            <li>PDF报告: 生成赛季综合报告，在新窗口中预览并可打印/保存为PDF</li>
          </ul>
        </Card>
      </Card>
    </div>
  )
}
