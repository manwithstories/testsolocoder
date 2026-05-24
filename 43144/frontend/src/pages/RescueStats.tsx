import React, { useEffect, useState } from 'react'
import { Row, Col, Card, Statistic, Button, DatePicker, message } from 'antd'
import { SmileOutlined, HeartOutlined, FileTextOutlined, CheckCircleOutlined, ExportOutlined } from '@ant-design/icons'
import { getRescueStats } from '../api/rescue'
import { exportAdoptionReport } from '../api/export'
import { RescueStats } from '../types'
import dayjs, { Dayjs } from 'dayjs'

const { RangePicker } = DatePicker

const RescueStatsPage: React.FC = () => {
  const [stats, setStats] = useState<RescueStats | null>(null)
  const [loading, setLoading] = useState(false)
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs] | null>(null)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    setLoading(true)
    try {
      const response: any = await getRescueStats()
      if (response && response.code === 0 && response.data) {
        setStats(response.data as RescueStats)
      }
    } catch (error) {
      console.error('Failed to load stats:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleExport = async () => {
    try {
      const params: any = {}
      if (dateRange) {
        params.start_date = dateRange[0].format('YYYY-MM-DD')
        params.end_date = dateRange[1].format('YYYY-MM-DD')
      }
      const blob = await exportAdoptionReport(params)
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `adoption_report_${dayjs().format('YYYYMMDD_HHmmss')}.xlsx`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      message.success('导出成功')
    } catch (error: any) {
      message.error('导出失败')
    }
  }

  if (loading || !stats) {
    return <div>加载中...</div>
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2 style={{ margin: 0 }}>数据统计</h2>
        <div>
          <RangePicker
            value={dateRange as any}
            onChange={(dates) => {
              if (dates && dates.length === 2) {
                setDateRange([dates[0] as Dayjs, dates[1] as Dayjs])
              } else {
                setDateRange(null)
              }
            }}
            style={{ marginRight: 16 }}
          />
          <Button type="primary" icon={<ExportOutlined />} onClick={handleExport}>
            导出报表
          </Button>
        </div>
      </div>

      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="宠物总数"
              value={stats.total_pets}
              prefix={<SmileOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="待领养"
              value={stats.adoptable_pets}
              valueStyle={{ color: '#52c41a' }}
              prefix={<SmileOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="已领养"
              value={stats.adopted_pets}
              valueStyle={{ color: '#1890ff' }}
              prefix={<HeartOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="治疗中"
              value={stats.treatment_pets}
              valueStyle={{ color: '#faad14' }}
              prefix={<SmileOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="领养申请"
              value={stats.total_adoptions}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="待审核申请"
              value={stats.pending_applications}
              valueStyle={{ color: '#faad14' }}
              prefix={<FileTextOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="回访记录"
              value={stats.total_follow_ups}
              prefix={<CheckCircleOutlined />}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="领养率"
              value={stats.adoption_rate.toFixed(1)}
              suffix="%"
              valueStyle={{ color: '#52c41a' }}
            />
          </Card>
        </Col>
        <Col xs={24} sm={12} md={6}>
          <Card>
            <Statistic
              title="回访完成率"
              value={stats.follow_up_rate.toFixed(1)}
              suffix="%"
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default RescueStatsPage
