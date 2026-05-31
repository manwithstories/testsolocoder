import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, Badge, Button } from 'antd'
import {
  HomeOutlined,
  PawOutlined,
  CalendarOutlined,
  FileTextOutlined,
  StarOutlined,
  ShoppingCartOutlined,
  BellOutlined,
  UserOutlined,
  BarChartOutlined,
  LogoutOutlined,
  StoreOutlined,
  TeamOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '@/context/AuthContext'
import { useQuery } from '@tanstack/react-query'
import { alertApi } from '@/services/api'

const { Header, Sider, Content } = Layout

export default function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout, token } = useAuthStore()

  const { data: alertsData } = useQuery({
    queryKey: ['alerts', 'unread'],
    queryFn: () => alertApi.list({ is_read: false, page_size: 100 }),
    enabled: !!token,
  })

  const unreadCount = alertsData?.data?.total || 0

  const handleLogout = () => {
    logout()
    navigate('/login', { replace: true })
  }

  const getMenuItems = () => {
    const baseItems = [
      {
        key: '/',
        icon: <HomeOutlined />,
        label: '首页',
      },
      {
        key: '/pets',
        icon: <PawOutlined />,
        label: '宠物档案',
      },
      {
        key: '/reservations',
        icon: <CalendarOutlined />,
        label: '寄养预约',
      },
      {
        key: '/daily-records',
        icon: <FileTextOutlined />,
        label: '实时动态',
      },
      {
        key: '/reviews',
        icon: <StarOutlined />,
        label: '评价管理',
      },
      {
        key: '/orders',
        icon: <ShoppingCartOutlined />,
        label: '订单支付',
      },
      {
        key: '/alerts',
        icon: <Badge count={unreadCount} size="small"><BellOutlined /></Badge>,
        label: '健康提醒',
      },
    ]

    if (user?.role === 'store' || user?.role === 'admin') {
      baseItems.push({
        key: '/statistics',
        icon: <BarChartOutlined />,
        label: '数据统计',
      })
      baseItems.push({
        key: '/store-dashboard',
        icon: <StoreOutlined />,
        label: '门店管理',
      })
    }

    if (user?.role === 'keeper') {
      baseItems.push({
        key: '/keeper-dashboard',
        icon: <TeamOutlined />,
        label: '管家工作台',
      })
    }

    baseItems.push({
      key: '/profile',
      icon: <UserOutlined />,
      label: '个人中心',
    })

    return baseItems
  }

  const userMenu = {
    items: [
      {
        key: 'profile',
        icon: <UserOutlined />,
        label: '个人中心',
        onClick: () => navigate('/profile'),
      },
      {
        type: 'divider' as const,
      },
      {
        key: 'logout',
        icon: <LogoutOutlined />,
        label: '退出登录',
        onClick: handleLogout,
      },
    ],
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        theme="light"
        width={220}
        style={{ position: 'sticky', top: 0, height: '100vh', overflow: 'auto' }}
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontWeight: 'bold',
            fontSize: 18,
            color: '#0ea5e9',
            borderBottom: '1px solid #f0f0f0',
          }}
        >
          <PawOutlined style={{ marginRight: 8 }} />
          宠物寄养平台
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          items={getMenuItems()}
          onClick={({ key }) => navigate(key)}
          style={{ borderRight: 0 }}
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
            boxShadow: '0 1px 4px rgba(0,21,41,.08)',
          }}
        >
          <div style={{ fontSize: 16, fontWeight: 500 }}>
            {user?.role === 'owner' && '宠物主人工作台'}
            {user?.role === 'store' && '门店管理后台'}
            {user?.role === 'keeper' && '宠物管家工作台'}
            {user?.role === 'admin' && '系统管理后台'}
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <Button
              type="text"
              icon={<BellOutlined />}
              onClick={() => navigate('/alerts')}
            >
              <Badge count={unreadCount} size="small" offset={[-2, 2]} />
            </Button>
            <Dropdown menu={userMenu} placement="bottomRight">
              <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
                <Avatar icon={<UserOutlined />} src={user?.avatar_url} />
                <span>{user?.username}</span>
              </div>
            </Dropdown>
          </div>
        </Header>
        <Content style={{ margin: 0, background: '#f5f7fa' }}>
          <div className="page-container">
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  )
}
