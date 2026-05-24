import { Routes, Route, Navigate, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Dropdown, Avatar, Badge, Button } from 'antd'
import {
  TrophyOutlined,
  TeamOutlined,
  CalendarOutlined,
  BarChartOutlined,
  CrownOutlined,
  UserOutlined,
  DollarOutlined,
  BellOutlined,
  ExportOutlined,
  LogoutOutlined,
  FileTextOutlined
} from '@ant-design/icons'
import { useState, useEffect } from 'react'
import api from './api'
import type { User } from './types'
import Login from './pages/Login'
import Leagues from './pages/Leagues'
import Teams from './pages/Teams'
import Matches from './pages/Matches'
import Standings from './pages/Standings'
import Knockout from './pages/Knockout'
import Stats from './pages/Stats'
import Referees from './pages/Referees'
import Fees from './pages/Fees'
import Notifications from './pages/Notifications'
import Export from './pages/Export'

const { Header, Sider, Content } = Layout

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem('token')
  if (!token) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

function App() {
  const navigate = useNavigate()
  const location = useLocation()
  const [user, setUser] = useState<User | null>(null)
  const [unreadCount, setUnreadCount] = useState(0)
  const [collapsed, setCollapsed] = useState(false)

  useEffect(() => {
    const saved = localStorage.getItem('user')
    if (saved) {
      setUser(JSON.parse(saved))
    }
    fetchUnread()
  }, [])

  const fetchUnread = () => {
    api.get('/notifications/unread-count').then((res) => {
      setUnreadCount(res.data.count)
    }).catch(() => {})
  }

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    setUser(null)
    navigate('/login')
  }

  const menuItems = [
    { key: '/leagues', icon: <TrophyOutlined />, label: '联赛管理' },
    { key: '/teams', icon: <TeamOutlined />, label: '球队管理' },
    { key: '/matches', icon: <CalendarOutlined />, label: '比赛赛程' },
    { key: '/standings', icon: <BarChartOutlined />, label: '积分榜' },
    { key: '/knockout', icon: <CrownOutlined />, label: '淘汰赛' },
    { key: '/stats', icon: <FileTextOutlined />, label: '球员统计' },
    { key: '/referees', icon: <UserOutlined />, label: '裁判管理' },
    { key: '/fees', icon: <DollarOutlined />, label: '费用管理' },
    { key: '/notifications', icon: <BellOutlined />, label: '消息通知' },
    { key: '/export', icon: <ExportOutlined />, label: '数据导出' },
  ]

  if (!user && location.pathname !== '/login') {
    return <Navigate to="/login" replace />
  }

  if (location.pathname === '/login') {
    return <Login onLogin={(u) => { setUser(u); navigate('/leagues') }} />
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed} theme="dark">
        <div style={{ color: '#fff', textAlign: 'center', padding: '16px', fontSize: collapsed ? 14 : 18, fontWeight: 'bold' }}>
          {collapsed ? 'SLM' : '体育联赛管理'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
          <Badge count={unreadCount} offset={[-4, 4]}>
            <Button type="text" icon={<BellOutlined />} onClick={() => { navigate('/notifications'); fetchUnread() }} />
          </Badge>
          <Dropdown
            menu={{
              items: [
                { key: 'profile', label: user?.full_name },
                { key: 'role', label: `角色: ${user?.role}` },
                { type: 'divider' },
                { key: 'logout', label: '退出登录', icon: <LogoutOutlined />, onClick: handleLogout }
              ]
            }}
            trigger={['click']}
          >
            <Avatar style={{ marginLeft: 16, cursor: 'pointer' }} icon={<UserOutlined />} />
          </Dropdown>
        </Header>
        <Content style={{ margin: 16, padding: 24, background: '#fff', borderRadius: 8 }}>
          <Routes>
            <Route path="/" element={<Navigate to="/leagues" replace />} />
            <Route path="/login" element={<Login onLogin={(u) => { setUser(u); navigate('/leagues') }} />} />
            <Route path="/leagues" element={<PrivateRoute><Leagues isAdmin={user?.role === 'admin'} /></PrivateRoute>} />
            <Route path="/teams" element={<PrivateRoute><Teams isAdmin={user?.role === 'admin'} /></PrivateRoute>} />
            <Route path="/matches" element={<PrivateRoute><Matches isAdmin={user?.role === 'admin'} /></PrivateRoute>} />
            <Route path="/standings" element={<PrivateRoute><Standings /></PrivateRoute>} />
            <Route path="/knockout" element={<PrivateRoute><Knockout isAdmin={user?.role === 'admin'} /></PrivateRoute>} />
            <Route path="/stats" element={<PrivateRoute><Stats isAdmin={user?.role === 'admin'} /></PrivateRoute>} />
            <Route path="/referees" element={<PrivateRoute><Referees isAdmin={user?.role === 'admin'} /></PrivateRoute>} />
            <Route path="/fees" element={<PrivateRoute><Fees isAdmin={user?.role === 'admin'} /></PrivateRoute>} />
            <Route path="/notifications" element={<PrivateRoute><Notifications onRead={fetchUnread} /></PrivateRoute>} />
            <Route path="/export" element={<PrivateRoute><Export /></PrivateRoute>} />
          </Routes>
        </Content>
      </Layout>
    </Layout>
  )
}

export default App
