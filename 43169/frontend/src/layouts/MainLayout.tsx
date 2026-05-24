import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, Button, Badge } from 'antd'
import {
  HomeOutlined,
  UserOutlined,
  HeartOutlined,
  CalendarOutlined,
  MessageOutlined,
  CrownOutlined,
  TeamOutlined,
  DashboardOutlined,
  LogoutOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '@/store/authStore'

const { Header, Sider, Content } = Layout

export default function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()

  const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: '首页' },
    { key: '/profile', icon: <UserOutlined />, label: '个人资料' },
    { key: '/match', icon: <HeartOutlined />, label: '智能匹配' },
    { key: '/dates', icon: <CalendarOutlined />, label: '约会管理' },
    { key: '/chat', icon: <MessageOutlined />, label: '消息' },
    { key: '/member', icon: <CrownOutlined />, label: '会员中心' },
    { key: '/matchmaker', icon: <TeamOutlined />, label: '红娘服务' },
  ]

  if (user?.role === 'admin') {
    menuItems.push({ key: '/admin', icon: <DashboardOutlined />, label: '管理后台' })
  }

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key)
  }

  const userMenu = {
    items: [
      { key: 'profile', icon: <UserOutlined />, label: '个人资料', onClick: () => navigate('/profile') },
      { key: 'verify', icon: <SafetyCertificateOutlined />, label: '实名认证', onClick: () => navigate('/verify') },
      { type: 'divider' as const },
      { key: 'logout', icon: <LogoutOutlined />, label: '退出登录', onClick: () => { logout(); navigate('/login') } },
    ],
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header className="app-header">
        <div className="app-logo">💕 相亲交友平台</div>
        <div className="app-header-user">
          <Badge dot>
            <Button type="text" icon={<MessageOutlined />} onClick={() => navigate('/chat')} />
          </Badge>
          <Dropdown menu={userMenu} placement="bottomRight">
            <span style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
              <Avatar src={user?.avatar} icon={<UserOutlined />} />
              <span>{user?.username}</span>
            </span>
          </Dropdown>
        </div>
      </Header>
      <Layout>
        <Sider width={200} theme="light">
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={handleMenuClick}
            style={{ height: '100%', borderRight: 0 }}
          />
        </Sider>
        <Layout>
          <Content className="page-container">
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  )
}
