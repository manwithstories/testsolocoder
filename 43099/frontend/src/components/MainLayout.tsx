import { Layout, Menu, Avatar, Dropdown, Button } from 'antd'
import {
  DashboardOutlined,
  ShopOutlined,
  ToolOutlined,
  CalendarOutlined,
  FileTextOutlined,
  PayCircleOutlined,
  StarOutlined,
  BarChartOutlined,
  LogoutOutlined,
  UserOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons'
import { useNavigate, useLocation, Outlet } from 'react-router-dom'
import { useAuthStore } from '@/store/authStore'
import { useAppStore } from '@/store/appStore'
import type { MenuProps } from 'antd'

const { Header, Sider, Content } = Layout

const MainLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const { collapsed, toggleCollapsed } = useAppStore()

  const isAdmin = user?.role === 'admin' || user?.role === 'super_admin'
  const isSuperAdmin = user?.role === 'super_admin'

  const menuItems: MenuProps['items'] = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: '首页',
    },
    {
      key: '/venues',
      icon: <ShopOutlined />,
      label: '场地管理',
    },
    {
      key: '/devices',
      icon: <ToolOutlined />,
      label: '设备管理',
    },
    {
      key: '/calendar',
      icon: <CalendarOutlined />,
      label: '预约日历',
    },
    {
      key: '/orders',
      icon: <FileTextOutlined />,
      label: '订单管理',
    },
    isAdmin && {
      key: '/payments',
      icon: <PayCircleOutlined />,
      label: '支付管理',
    },
    {
      key: '/reviews',
      icon: <StarOutlined />,
      label: '评价管理',
    },
    isAdmin && {
      key: '/stats',
      icon: <BarChartOutlined />,
      label: '数据统计',
    },
    isSuperAdmin && {
      key: '/users',
      icon: <UserOutlined />,
      label: '用户管理',
    },
  ].filter(Boolean) as MenuProps['items']

  const userMenu: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: () => {
        logout()
        navigate('/login')
      },
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed} theme="dark">
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white', fontSize: collapsed ? 14 : 18, fontWeight: 'bold' }}>
          {collapsed ? 'VBS' : '场地预约系统'}
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
        <Header style={{ padding: '0 16px', background: '#fff', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={toggleCollapsed}
          />
          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <span>欢迎，{user?.real_name || user?.username}</span>
            <Dropdown menu={{ items: userMenu }} placement="bottomRight">
              <Avatar src={user?.avatar} icon={<UserOutlined />} style={{ cursor: 'pointer' }} />
            </Dropdown>
          </div>
        </Header>
        <Content style={{ margin: '16px', padding: 24, background: '#fff', borderRadius: 8, minHeight: 'calc(100vh - 112px)' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
