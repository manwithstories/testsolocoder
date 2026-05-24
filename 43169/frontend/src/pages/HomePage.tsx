import { useQuery } from '@tanstack/react-query'
import { Row, Col, Card, Statistic, Button } from 'antd'
import { UserOutlined, HeartOutlined, CalendarOutlined, MessageOutlined, StarOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { statsApi, matchApi } from '@/api/endpoints'

export default function HomePage() {
  const navigate = useNavigate()
  const { user } = useAuthStore()

  const { data: stats } = useQuery({
    queryKey: ['platformStats'],
    queryFn: statsApi.getPlatformStats,
    enabled: user?.role === 'admin',
  })

  const { data: matchData } = useQuery({
    queryKey: ['smartMatch'],
    queryFn: () => matchApi.smartMatch({ page: 1, page_size: 6 }),
  })

  return (
    <div>
      <Card style={{ marginBottom: 24 }}>
        <h2 style={{ marginBottom: 16 }}>欢迎回来，{user?.username} 👋</h2>
        <Row gutter={16}>
          <Col xs={12} sm={8} md={6}>
            <Card size="small" hoverable onClick={() => navigate('/match')}>
              <Statistic title="智能匹配" value={0} prefix={<HeartOutlined style={{ color: '#ff6b81' }} />} />
              <p style={{ color: '#888', marginTop: 8 }}>发现心仪的TA</p>
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card size="small" hoverable onClick={() => navigate('/dates')}>
              <Statistic title="约会管理" value={0} prefix={<CalendarOutlined style={{ color: '#52c41a' }} />} />
              <p style={{ color: '#888', marginTop: 8 }}>管理你的约会</p>
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card size="small" hoverable onClick={() => navigate('/chat')}>
              <Statistic title="消息中心" value={0} prefix={<MessageOutlined style={{ color: '#1890ff' }} />} />
              <p style={{ color: '#888', marginTop: 8 }}>与TA畅聊</p>
            </Card>
          </Col>
          <Col xs={12} sm={8} md={6}>
            <Card size="small" hoverable onClick={() => navigate('/member')}>
              <Statistic title="会员中心" value={0} prefix={<StarOutlined style={{ color: '#faad14' }} />} />
              <p style={{ color: '#888', marginTop: 8 }}>升级会员享特权</p>
            </Card>
          </Col>
        </Row>
      </Card>

      {user?.role === 'admin' && stats && (
        <Card title="平台数据概览" style={{ marginBottom: 24 }}>
          <Row gutter={16}>
            <Col span={6}><Statistic title="总用户数" value={stats.total_users} /></Col>
            <Col span={6}><Statistic title="今日活跃" value={stats.active_today} /></Col>
            <Col span={6}><Statistic title="已认证用户" value={stats.verified_users} /></Col>
            <Col span={6}><Statistic title="匹配成功率" value={stats.match_success_rate} suffix="%" /></Col>
          </Row>
        </Card>
      )}

      <Card
        title="为你推荐"
        extra={<Button type="link" onClick={() => navigate('/match')}>查看更多</Button>}
      >
        <Row gutter={[16, 16]}>
          {matchData?.list.slice(0, 6).map((item) => (
            <Col xs={12} sm={8} md={6} lg={4} key={item.user_id}>
              <Card
                hoverable
                cover={
                  <div style={{ height: 160, background: '#f0f0f0', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                    <UserOutlined style={{ fontSize: 48, color: '#ccc' }} />
                  </div>
                }
                style={{ cursor: 'pointer' }}
                onClick={() => navigate(`/user/${item.user_id}`)}
              >
                <Card.Meta
                  title={item.profile?.nickname || `用户${item.user_id}`}
                  description={`匹配度 ${item.match_score}%`}
                />
              </Card>
            </Col>
          ))}
        </Row>
      </Card>
    </div>
  )
}
