import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Avatar, Dropdown, Badge, Button } from 'antd'
import {
  DashboardOutlined,
  ApiOutlined,
  MedicineBoxOutlined,
  AppleOutlined,
  ShoppingOutlined,
  ShoppingCartOutlined,
  FileSearchOutlined,
  TeamOutlined,
  BarChartOutlined,
  UserOutlined,
  LogoutOutlined,
  BellOutlined,
} from '@ant-design/icons'
import { useAuthStore } from '../../store/authStore'
import { useState, useEffect } from 'react'
import api from '../../api'

const { Header, Sider, Content } = Layout

function MainLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout, isAuthenticated } = useAuthStore()
  const [unreadCount, setUnreadCount] = useState(0)

  useEffect(() => {
    if (isAuthenticated) {
      fetchUnreadCount()
    }
  }, [isAuthenticated])

  const fetchUnreadCount = async () => {
    try {
      const response = await api.get('/notifications/unread-count')
      setUnreadCount(response.data.count)
    } catch (error) {
      console.error('Failed to fetch unread count:', error)
    }
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const getMenuItems = () => {
    const items = [
      {
        key: '/',
        icon: <DashboardOutlined />,
        label: '仪表盘',
      },
    ]

    if (user?.role === 'beekeeper') {
      items.push(
        {
          key: '/beehives',
          icon: <ApiOutlined />,
          label: '蜂场管理',
        },
        {
          key: '/health-records',
          icon: <MedicineBoxOutlined />,
          label: '健康监控',
        },
        {
          key: '/harvests',
          icon: <AppleOutlined />,
          label: '采收记录',
        },
        {
          key: '/inventory',
          icon: <ShoppingOutlined />,
          label: '库存管理',
        },
        {
          key: '/my-products',
          icon: <ShoppingCartOutlined />,
          label: '我的商品',
        },
        {
          key: '/inspections',
          icon: <FileSearchOutlined />,
          label: '检测预约',
        },
        {
          key: '/orders',
          icon: <ShoppingCartOutlined />,
          label: '订单管理',
        }
      )
    }

    if (user?.role === 'buyer') {
      items.push(
        {
          key: '/products',
          icon: <ShoppingOutlined />,
          label: '交易市场',
        },
        {
          key: '/orders',
          icon: <ShoppingCartOutlined />,
          label: '我的订单',
        }
      )
    }

    if (user?.role === 'inspector') {
      items.push({
        key: '/inspections',
        icon: <FileSearchOutlined />,
        label: '检测任务',
      })
    }

    items.push(
      {
        key: '/community',
        icon: <TeamOutlined />,
        label: '蜂农社区',
      },
      {
        key: '/analytics',
        icon: <BarChartOutlined />,
        label: '数据分析',
      }
    )

    return items
  }

  const userMenu = {
    items: [
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
      <Sider theme="dark" width={200}>
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#fff', fontSize: 18, fontWeight: 'bold' }}>
          养蜂平台
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={getMenuItems()}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
          <div style={{ marginRight: 24 }}>
            <Badge count={unreadCount} size="small">
              <Button
                type="text"
                icon={<BellOutlined />}
                onClick={() => navigate('/notifications')}
              />
            </Badge>
          </div>
          <Dropdown menu={userMenu}>
            <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
              <Avatar icon={<UserOutlined />} src={user?.avatar} />
              <span>{user?.username}</span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
