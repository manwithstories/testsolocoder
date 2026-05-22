import React from 'react'
import { Layout, Menu, Avatar, Dropdown, Button, Space } from 'antd'
import {
  HomeOutlined,
  BookOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  SettingOutlined,
  DashboardOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'

const { Header } = Layout

const AppHeader: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, isAuthenticated, logout } = useAuthStore()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenu = {
    items: [
      ...(user?.role === 'instructor' || user?.role === 'admin'
        ? [{ key: 'instructor', label: <Link to="/instructor/dashboard">讲师中心</Link>, icon: <DashboardOutlined /> }]
        : []),
      ...(user?.role === 'admin'
        ? [{ key: 'admin', label: <Link to="/admin/dashboard">管理后台</Link>, icon: <SettingOutlined /> }]
        : []),
      { key: 'profile', label: <Link to="/profile">个人中心</Link>, icon: <UserOutlined /> },
      { key: 'my-courses', label: <Link to="/my-courses">我的课程</Link>, icon: <BookOutlined /> },
      { key: 'orders', label: <Link to="/orders">我的订单</Link>, icon: <ShoppingCartOutlined /> },
      { type: 'divider' as const },
      { key: 'logout', label: '退出登录', icon: <LogoutOutlined />, onClick: handleLogout },
    ],
  }

  const navItems = [
    { key: '/', label: <Link to="/">首页</Link>, icon: <HomeOutlined /> },
    { key: '/courses', label: <Link to="/courses">课程</Link>, icon: <BookOutlined /> },
  ]

  return (
    <Header
      style={{
        display: 'flex',
        alignItems: 'center',
        background: '#fff',
        boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
        padding: '0 24px',
      }}
    >
      <Link to="/" style={{ fontSize: 20, fontWeight: 600, marginRight: 48, color: '#1890ff' }}>
        在线学习平台
      </Link>
      <Menu
        theme="light"
        mode="horizontal"
        selectedKeys={[location.pathname]}
        items={navItems}
        style={{ flex: 1, borderRight: 'none' }}
      />
      <Space>
        {isAuthenticated ? (
          <Dropdown menu={userMenu} placement="bottomRight">
            <Space style={{ cursor: 'pointer' }}>
              <Avatar size="small" src={user?.avatar} icon={<UserOutlined />} />
              <span>{user?.nickname || user?.username}</span>
            </Space>
          </Dropdown>
        ) : (
          <Space>
            <Button type="link" onClick={() => navigate('/login')}>登录</Button>
            <Button type="primary" onClick={() => navigate('/register')}>注册</Button>
          </Space>
        )}
      </Space>
    </Header>
  )
}

export default AppHeader
