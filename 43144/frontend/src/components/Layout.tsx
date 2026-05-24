import React from 'react'
import { Layout, Menu, Avatar, Dropdown, Button } from 'antd'
import {
  HomeOutlined,
  SmileOutlined,
  CalendarOutlined,
  FileTextOutlined,
  DashboardOutlined,
  UserOutlined,
  LogoutOutlined,
  SettingOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const { Header, Sider, Content } = Layout

const MainLayout: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()

  const getMenuItems = () => {
    const items: any[] = [
      {
        key: '/',
        icon: <HomeOutlined />,
        label: '首页',
        onClick: () => navigate('/'),
      },
      {
        key: '/pets',
        icon: <SmileOutlined />,
        label: '宠物列表',
        onClick: () => navigate('/pets'),
      },
    ]

    if (user) {
      if (user.role === 'adopter') {
        items.push(
          {
            key: '/adopted',
            icon: <SmileOutlined />,
            label: '我的领养',
            onClick: () => navigate('/adopted'),
          },
          {
            key: '/appointments',
            icon: <CalendarOutlined />,
            label: '我的预约',
            onClick: () => navigate('/appointments'),
          },
        )
      }

      if (user.role === 'rescue') {
        items.push(
          {
            key: '/my-pets',
            icon: <SmileOutlined />,
            label: '宠物管理',
            onClick: () => navigate('/my-pets'),
          },
          {
            key: '/adoption-applications',
            icon: <FileTextOutlined />,
            label: '领养申请',
            onClick: () => navigate('/adoption-applications'),
          },
          {
            key: '/rescue-stats',
            icon: <DashboardOutlined />,
            label: '数据统计',
            onClick: () => navigate('/rescue-stats'),
          },
        )
      }

      if (user.role === 'admin') {
        items.push(
          {
            key: '/admin/rescues',
            icon: <DashboardOutlined />,
            label: '救助站管理',
            onClick: () => navigate('/admin/rescues'),
          },
          {
            key: '/admin/users',
            icon: <UserOutlined />,
            label: '用户管理',
            onClick: () => navigate('/admin/users'),
          },
        )
      }

      if (user.role === 'adopter' || user.role === 'admin') {
        items.push(
          {
            key: '/health-records',
            icon: <FileTextOutlined />,
            label: '健康档案',
            onClick: () => navigate('/health-records'),
          },
        )
      }
    }

    return items
  }

  const userMenu = {
    items: [
      {
        key: 'profile',
        icon: <SettingOutlined />,
        label: '个人设置',
        onClick: () => navigate('/profile'),
      },
      { type: 'divider' as const },
      {
        key: 'logout',
        icon: <LogoutOutlined />,
        label: '退出登录',
        onClick: () => {
          logout()
          navigate('/login')
        },
      },
    ],
  }

  return (
    <Layout className="app-container">
      <Sider
        breakpoint="lg"
        collapsedWidth="0"
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'sticky',
          top: 0,
          left: 0,
        }}
      >
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <SmileOutlined style={{ fontSize: 24, color: '#fff' }} />
          <span style={{ color: '#fff', fontSize: 16, marginLeft: 8 }}>宠物领养平台</span>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={getMenuItems()}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <div style={{ fontSize: 18, fontWeight: 600 }}>
            宠物领养与健康档案管理平台
          </div>
          {user && (
            <Dropdown menu={userMenu} placement="bottomRight">
              <Button type="text" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Avatar size="small" icon={<UserOutlined />} />
                <span>{user.name}</span>
              </Button>
            </Dropdown>
          )}
        </Header>
        <Content style={{ background: '#f0f2f5', minHeight: 'calc(100vh - 64px)' }}>
          <div className="content-container">{children}</div>
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
