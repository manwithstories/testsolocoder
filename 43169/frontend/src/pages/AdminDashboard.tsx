import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Card, Row, Col, Statistic, Table, Button, DatePicker, Tabs, Tag, Select, Modal, message } from 'antd'
import {
  UserOutlined,
  TeamOutlined,
  CalendarOutlined,
  MessageOutlined,
  SafetyCertificateOutlined,
  DownloadOutlined,
  CheckOutlined,
  CloseOutlined,
  ExportOutlined,
} from '@ant-design/icons'
import dayjs, { Dayjs } from 'dayjs'
import { RangePickerProps } from 'antd/es/date-picker'
import ReactECharts from 'echarts-for-react'
import { statsApi, userApi, matchmakerApi, PlatformStats, DailyStats, UserInfo } from '@/api/endpoints'

const { RangePicker } = DatePicker

export default function AdminDashboard() {
  const [dateRange, setDateRange] = useState<[Dayjs, Dayjs]>([
    dayjs().subtract(14, 'day'),
    dayjs(),
  ])
  const [verifyModal, setVerifyModal] = useState(false)
  const [selectedUser, setSelectedUser] = useState<UserInfo | null>(null)

  const { data: statsData } = useQuery({
    queryKey: ['platformStats'],
    queryFn: statsApi.getPlatformStats,
  })

  const { data: dailyData } = useQuery({
    queryKey: ['dailyStats', dateRange],
    queryFn: () =>
      statsApi.getDailyStats({
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      }),
  })

  const { data: usersData, refetch: refetchUsers } = useQuery({
    queryKey: ['adminUsers'],
    queryFn: () => userApi.listUsers({ page: 1, page_size: 50 }),
  })

  const { data: matchmakerStatsData } = useQuery({
    queryKey: ['matchmakerStatsAdmin'],
    queryFn: () => statsApi.getMatchmakerStats({}),
  })

  const stats = statsData as PlatformStats | undefined
  const daily = dailyData as DailyStats[] || []

  const handleUserAction = (user: UserInfo, action: 'approve' | 'reject' | 'disable' | 'enable') => {
    const actionMap: Record<string, any> = {
      approve: userApi.approveVerify,
      reject: userApi.rejectVerify,
      disable: userApi.disableUser,
      enable: userApi.enableUser,
    }
    actionMap[action](user.id).then(() => {
      message.success('操作成功')
      refetchUsers()
    })
  }

  const handleExport = (type: 'excel' | 'pdf') => {
    const params = {
      start_date: dateRange[0].format('YYYY-MM-DD'),
      end_date: dateRange[1].format('YYYY-MM-DD'),
    }
    const apiCall = type === 'excel' ? statsApi.exportExcel : statsApi.exportPDF
    apiCall(params).then((data: any) => {
      const blob = data?.data || data
      const url = window.URL.createObjectURL(new Blob([blob]))
      const link = document.createElement('a')
      link.href = url
      link.download = type === 'excel' ? 'report.xlsx' : 'report.pdf'
      link.click()
      window.URL.revokeObjectURL(url)
      message.success('导出成功')
    })
  }

  const userColumns = [
    { title: 'ID', dataIndex: 'id', key: 'id' },
    { title: '用户名', dataIndex: 'username', key: 'username' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    { title: '角色', dataIndex: 'role', key: 'role', render: (role: string) => <Tag color={role === 'admin' ? 'red' : role === 'matchmaker' ? 'purple' : 'blue'}>{role}</Tag> },
    {
      title: '认证状态',
      dataIndex: 'verify_status',
      key: 'verify_status',
      render: (status: string) => (
        <Tag color={status === 'verified' ? 'green' : status === 'pending' ? 'orange' : 'red'}>
          {status === 'verified' ? '已认证' : status === 'pending' ? '待审核' : '已拒绝'}
        </Tag>
      ),
    },
    {
      title: '账号状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>
          {status === 'active' ? '正常' : '已禁用'}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: UserInfo) => (
        <div style={{ display: 'flex', gap: 4 }}>
          {record.verify_status === 'pending' && (
            <>
              <Button size="small" type="primary" icon={<CheckOutlined />} onClick={() => handleUserAction(record, 'approve')}>
                通过
              </Button>
              <Button size="small" danger icon={<CloseOutlined />} onClick={() => handleUserAction(record, 'reject')}>
                拒绝
              </Button>
            </>
          )}
          {record.status === 'active' ? (
            <Button size="small" danger onClick={() => handleUserAction(record, 'disable')}>
              禁用
            </Button>
          ) : (
            <Button size="small" type="primary" onClick={() => handleUserAction(record, 'enable')}>
              启用
            </Button>
          )}
        </div>
      ),
    },
  ]

  const chartOption = {
    tooltip: { trigger: 'axis' },
    legend: { data: ['新增用户', '活跃用户', '约会数', '消息数'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', boundaryGap: false, data: daily.map((d) => d.date) },
    yAxis: { type: 'value' },
    series: [
      { name: '新增用户', type: 'line', data: daily.map((d) => d.new_users), smooth: true, areaStyle: {} },
      { name: '活跃用户', type: 'line', data: daily.map((d) => d.active_users), smooth: true },
      { name: '约会数', type: 'line', data: daily.map((d) => d.dates_created), smooth: true },
      { name: '消息数', type: 'line', data: daily.map((d) => d.messages_count), smooth: true },
    ],
  }

  return (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
          <h2 style={{ margin: 0 }}>数据统计</h2>
          <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
            <RangePicker
              value={dateRange}
              onChange={(dates) => dates && setDateRange(dates as [Dayjs, Dayjs])}
            />
            <Button icon={<ExportOutlined />} onClick={() => handleExport('excel')}>
              导出Excel
            </Button>
            <Button icon={<ExportOutlined />} onClick={() => handleExport('pdf')}>
              导出PDF
            </Button>
          </div>
        </div>
        <Row gutter={16}>
          <Col xs={12} sm={6}>
            <Card><Statistic title="总用户数" value={stats?.total_users || 0} prefix={<UserOutlined />} /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card><Statistic title="今日活跃" value={stats?.active_today || 0} /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card><Statistic title="已认证用户" value={stats?.verified_users || 0} prefix={<SafetyCertificateOutlined />} /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card><Statistic title="红娘总数" value={stats?.total_matchmakers || 0} prefix={<TeamOutlined />} /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card><Statistic title="约会总数" value={stats?.total_dates || 0} prefix={<CalendarOutlined />} /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card><Statistic title="完成约会" value={stats?.completed_dates || 0} /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card><Statistic title="匹配成功率" value={stats?.match_success_rate || 0} precision={1} suffix="%" /></Card>
          </Col>
          <Col xs={12} sm={6}>
            <Card><Statistic title="消息总数" value={stats?.total_messages || 0} prefix={<MessageOutlined />} /></Card>
          </Col>
        </Row>
      </Card>

      <Card title="趋势图表" style={{ marginBottom: 24 }}>
        <ReactECharts option={chartOption} style={{ height: 350 }} />
      </Card>

      <Card>
        <Tabs
          items={[
            {
              key: 'users',
              label: '用户管理',
              children: (
                <Table
                  columns={userColumns}
                  dataSource={usersData?.list || []}
                  rowKey="id"
                  pagination={false}
                />
              ),
            },
            {
              key: 'matchmakers',
              label: '红娘业绩',
              children: (
                <Table
                  dataSource={matchmakerStatsData || []}
                  rowKey="id"
                  pagination={false}
                  columns={[
                    { title: '红娘', dataIndex: 'matchmaker_name', key: 'matchmaker_name' },
                    { title: '会员数', dataIndex: 'total_members', key: 'total_members' },
                    { title: '服务数', dataIndex: 'total_services', key: 'total_services' },
                    { title: '约会数', dataIndex: 'total_dates', key: 'total_dates' },
                    { title: '成功约会', dataIndex: 'success_dates', key: 'success_dates' },
                    { title: '平均分', dataIndex: 'avg_rating', key: 'avg_rating' },
                  ]}
                />
              ),
            },
          ]}
        />
      </Card>
    </div>
  )
}
