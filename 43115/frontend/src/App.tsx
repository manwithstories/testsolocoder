import React from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Menu, Dropdown, Avatar, Badge, Button } from 'antd'
import {
  HomeOutlined,
  AppstoreOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  MessageOutlined,
  SettingOutlined,
  WalletOutlined,
  FileTextOutlined,
  StarOutlined,
  DashboardOutlined,
  TeamOutlined,
  FundOutlined,
  SafetyCertificateOutlined,
} from '@ant-design/icons'
import { useAppSelector, useAppDispatch } from '@/store/hooks'
import { logout } from '@/store'
import { messageApi } from '@/services/message'
import { useState, useEffect } from 'react'

const { Header, Sider, Content } = Layout

const App: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const { userInfo, token } = useAppSelector((state) => state.auth)
  const [unreadCount, setUnreadCount] = useState(0)

  useEffect(() => {
    if (token) {
      messageApi.getUnreadCount().then((res) => {
        setUnreadCount(res.unread_count)
      })
    }
  }, [token])

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  const getUserMenuItems = () => {
    if (!userInfo) return []

    const items: any[] = [
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

    if (userInfo.role === 'customer') {
      items.push(
        {
          key: '/orders',
          icon: <ShoppingCartOutlined />,
          label: '我的订单',
        },
        {
          key: '/reviews',
          icon: <StarOutlined />,
          label: '评价管理',
        },
        {
          key: '/complaints',
          icon: <FileTextOutlined />,
          label: '投诉管理',
        },
        {
          key: '/addresses',
          icon: <SettingOutlined />,
          label: '地址管理',
        }
      )
    }

    if (userInfo.role === 'service_provider') {
      items.push(
        {
          key: '/invitations',
          icon: <MessageOutlined />,
          label: '预约邀请',
        },
        {
          key: '/orders',
          icon: <ShoppingCartOutlined />,
          label: '订单管理',
        },
        {
          key: '/my-services',
          icon: <AppstoreOutlined />,
          label: '我的服务',
        },
        {
          key: '/bills',
          icon: <WalletOutlined />,
          label: '收入账单',
        },
        {
          key: '/withdraws',
          icon: <FundOutlined />,
          label: '提现管理',
        },
        {
          key: '/certification',
          icon: <SafetyCertificateOutlined />,
          label: '认证管理',
        }
      )
    }

    if (userInfo.role === 'admin') {
      items.push(
        {
          key: '/admin/dashboard',
          icon: <DashboardOutlined />,
          label: '数据看板',
        },
        {
          key: '/admin/users',
          icon: <TeamOutlined />,
          label: '用户管理',
        },
        {
          key: '/admin/categories',
          icon: <AppstoreOutlined />,
          label: '分类管理',
        },
        {
          key: '/admin/withdraws',
          icon: <FundOutlined />,
          label: '提现审核',
        },
        {
          key: '/admin/complaints',
          icon: <FileTextOutlined />,
          label: '投诉处理',
        }
      )
    }

    return items
  }

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'messages',
      icon: <MessageOutlined />,
      label: '消息中心',
      onClick: () => navigate('/messages'),
    },
    { type: 'divider' as const },
    {
      key: 'logout',
      icon: <SettingOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header
        style={{
          background: '#fff',
          padding: '0 24px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
        }}
      >
        <div
          style={{ fontSize: 20, fontWeight: 'bold', cursor: 'pointer' }}
          onClick={() => navigate('/')}
        >
          家政服务平台
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          <Badge count={unreadCount} size="small">
            <Button
              type="text"
              icon={<MessageOutlined />}
              onClick={() => navigate('/messages')}
            />
          </Badge>
          {token ? (
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
                <Avatar icon={<UserOutlined />} src={userInfo?.avatar} />
                <span>{userInfo?.nickname || userInfo?.phone}</span>
              </div>
            </Dropdown>
          ) : (
            <Button type="primary" onClick={() => navigate('/login')}>
              登录
            </Button>
          )}
        </div>
      </Header>
      <Layout>
        {token && (
          <Sider width={200} theme="light">
            <Menu
              mode="inline"
              selectedKeys={[location.pathname]}
              items={getUserMenuItems()}
              onClick={({ key }) => navigate(key)}
              style={{ height: '100%', borderRight: 0 }}
            />
          </Sider>
        )}
        <Layout style={{ padding: 24 }}>
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

export default App
