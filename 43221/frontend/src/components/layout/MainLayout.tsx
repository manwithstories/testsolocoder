import { Outlet, useNavigate } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, Badge, Button } from 'antd'
import {
  HomeOutlined,
  ScheduleOutlined,
  UserOutlined,
  MessageOutlined,
  BarChartOutlined,
  SettingOutlined,
  LogoutOutlined,
  TeamOutlined,
  FileTextOutlined,
  StarOutlined,
  AppstoreOutlined,
} from '@ant-design/icons'
import { useAuthContext } from '@/contexts/AuthContext'
import { useUnreadCount } from '@/hooks/useNotifications'

const { Header, Sider, Content } = Layout

export function MainLayout() {
  const navigate = useNavigate()
  const { user, logout } = useAuthContext()
  const { data: unreadData } = useUnreadCount()

  const getMenuItems = () => {
    const baseItems = [
      {
        key: '/',
        icon: <HomeOutlined />,
        label: '首页',
      },
      {
        key: '/services',
        icon: <AppstoreOutlined />,
        label: '服务市场',
      },
    ]

    if (user?.role === 'professional' || user?.role === 'client') {
      baseItems.push({
        key: '/appointments',
        icon: <ScheduleOutlined />,
        label: '我的预约',
      })
    }

    if (user?.role === 'professional') {
      baseItems.push({
        key: '/services/mine',
        icon: <FileTextOutlined />,
        label: '我的服务',
      })
      baseItems.push({
        key: '/reviews',
        icon: <StarOutlined />,
        label: '我的评价',
      })
      baseItems.push({
        key: '/statistics',
        icon: <BarChartOutlined />,
        label: '数据统计',
      })
    }

    if (user?.role === 'client') {
      baseItems.push({
        key: '/reviews',
        icon: <StarOutlined />,
        label: '我的评价',
      })
    }

    if (user?.role === 'admin') {
      baseItems.push(
        {
          key: '/admin',
          icon: <TeamOutlined />,
          label: '用户管理',
        },
        {
          key: '/admin/verifications',
          icon: <UserOutlined />,
          label: '资质审核',
        },
        {
          key: '/reviews/management',
          icon: <StarOutlined />,
          label: '评价管理',
        },
        {
          key: '/statistics',
          icon: <BarChartOutlined />,
          label: '数据统计',
        },
        {
          key: '/admin/templates',
          icon: <SettingOutlined />,
          label: '模板管理',
        }
      )
    }

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
        onClick: () => {
          logout()
          navigate('/login')
        },
      },
    ],
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={220}>
        <div style={{ padding: '16px', textAlign: 'center', fontSize: '18px', fontWeight: 600, borderBottom: '1px solid #f0f0f0' }}>
          咨询服务平台
        </div>
        <Menu
          mode="inline"
          selectedKeys={[window.location.pathname]}
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
            boxShadow: '0 1px 4px rgba(0,0,0,0.1)',
          }}
        >
          <div style={{ fontSize: '16px', fontWeight: 500 }}>
            {user?.role === 'professional' && '专业人士工作台'}
            {user?.role === 'client' && '客户工作台'}
            {user?.role === 'admin' && '管理后台'}
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
            <Badge count={unreadData?.unread_count || 0} size="small">
              <Button
                type="text"
                icon={<MessageOutlined />}
                onClick={() => navigate('/notifications')}
              />
            </Badge>
            <Dropdown menu={userMenu} placement="bottomRight">
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                <Avatar src={user?.avatar}>
                  {user?.full_name?.charAt(0) || 'U'}
                </Avatar>
                <span>{user?.full_name}</span>
              </div>
            </Dropdown>
          </div>
        </Header>
        <Content style={{ margin: '24px' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
