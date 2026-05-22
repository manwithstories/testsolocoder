import { Layout, Menu, Button, Dropdown, Avatar } from 'antd'
import {
  DashboardOutlined,
  ShopOutlined,
  AppstoreOutlined,
  UserOutlined,
  WarningOutlined,
  LogoutOutlined,
  SafetyOutlined,
} from '@ant-design/icons'
import { Link, Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAppStore } from '@/store'

const { Header, Sider, Content } = Layout

const AdminLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAppStore()

  const menuItems = [
    {
      key: '/admin/statistics',
      icon: <DashboardOutlined />,
      label: '数据概览',
      onClick: () => navigate('/admin/statistics'),
    },
    {
      key: '/admin/shops',
      icon: <ShopOutlined />,
      label: '商家管理',
      onClick: () => navigate('/admin/shops'),
    },
    {
      key: '/admin/categories',
      icon: <AppstoreOutlined />,
      label: '分类管理',
      onClick: () => navigate('/admin/categories'),
    },
    {
      key: '/admin/users',
      icon: <UserOutlined />,
      label: '用户管理',
      onClick: () => navigate('/admin/users'),
    },
    {
      key: '/admin/disputes',
      icon: <WarningOutlined />,
      label: '纠纷处理',
      onClick: () => navigate('/admin/disputes'),
    },
  ]

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#faad14', padding: '0 24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Link to="/" style={{ color: '#fff', fontSize: 18, fontWeight: 'bold' }}>
          <SafetyOutlined /> 管理后台
        </Link>
        <Dropdown
          menu={{
            items: [
              {
                key: '1',
                icon: <UserOutlined />,
                label: user?.nickname || user?.username,
                disabled: true,
              },
              {
                key: '2',
                icon: <LogoutOutlined />,
                label: '退出登录',
                onClick: handleLogout,
              },
            ],
          }}
          placement="bottomRight"
        >
          <Button type="text" style={{ color: '#fff' }}>
            <Avatar icon={<UserOutlined />} size="small" style={{ marginRight: 8 }} />
            {user?.nickname || user?.username}
          </Button>
        </Dropdown>
      </Header>
      <Layout>
        <Sider width={200} style={{ background: '#fff' }}>
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            style={{ height: '100%', borderRight: 0 }}
          />
        </Sider>
        <Layout style={{ padding: '24px' }}>
          <Content
            style={{
              background: '#fff',
              padding: 24,
              margin: 0,
              minHeight: 280,
              borderRadius: 8,
            }}
          >
            <Outlet />
          </Content>
        </Layout>
      </Layout>
    </Layout>
  )
}

export default AdminLayout
