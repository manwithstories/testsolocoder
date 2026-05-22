import { Layout, Menu, Badge, Input, Dropdown, Avatar } from 'antd'
import {
  HomeOutlined,
  ShoppingCartOutlined,
  HeartOutlined,
  BellOutlined,
  UserOutlined,
  ShopOutlined,
  LogoutOutlined,
  SettingOutlined,
} from '@ant-design/icons'
import { Link, Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAppStore } from '@/store'

const { Header, Content, Footer } = Layout
const { Search } = Input

const MainLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, cartCount, logout } = useAppStore()

  const handleSearch = (value: string) => {
    if (value.trim()) {
      navigate(`/products?keyword=${encodeURIComponent(value)}`)
    }
  }

  const userMenuItems = [
    {
      key: '1',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      key: '2',
      icon: <ShoppingCartOutlined />,
      label: '我的订单',
      onClick: () => navigate('/orders'),
    },
    {
      key: '3',
      icon: <HeartOutlined />,
      label: '我的收藏',
      onClick: () => navigate('/favorites'),
    },
    {
      key: '4',
      icon: <SettingOutlined />,
      label: '账号设置',
      onClick: () => navigate('/settings'),
    },
    {
      type: 'divider' as const,
    },
    {
      key: '5',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: () => {
        logout()
        navigate('/login')
      },
    },
  ]

  const navItems = [
    { key: '/', icon: <HomeOutlined />, label: <Link to="/">首页</Link> },
    { key: '/products', label: <Link to="/products">商品</Link> },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#fff', padding: '0 24px', boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', height: '100%' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 48 }}>
            <Link to="/" style={{ fontSize: 20, fontWeight: 'bold', color: '#1890ff' }}>
              多商家商城
            </Link>
            <Menu
              mode="horizontal"
              selectedKeys={[location.pathname]}
              items={navItems}
              style={{ borderBottom: 'none' }}
            />
          </div>

          <Search
            placeholder="搜索商品..."
            allowClear
            style={{ width: 400 }}
            onSearch={handleSearch}
          />

          <div style={{ display: 'flex', alignItems: 'center', gap: 24 }}>
            <Link to="/cart" style={{ fontSize: 16 }}>
              <Badge count={cartCount} size="small">
                <ShoppingCartOutlined style={{ fontSize: 20 }} />
              </Badge>
            </Link>
            <Link to="/favorites" style={{ fontSize: 16 }}>
              <HeartOutlined style={{ fontSize: 20 }} />
            </Link>
            <Link to="/notifications" style={{ fontSize: 16 }}>
              <Badge dot>
                <BellOutlined style={{ fontSize: 20 }} />
              </Badge>
            </Link>

            {user ? (
              <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
                <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
                  <Avatar icon={<UserOutlined />} src={user.avatar} />
                  <span>{user.nickname || user.username}</span>
                </div>
              </Dropdown>
            ) : (
              <div>
                <Link to="/login" style={{ marginRight: 16 }}>登录</Link>
                <Link to="/register">注册</Link>
              </div>
            )}

            {user?.role === 'seller' && (
              <Link to="/seller" target="_blank" style={{ color: '#52c41a' }}>
                <ShopOutlined /> 商家中心
              </Link>
            )}
            {user?.role === 'admin' && (
              <Link to="/admin" target="_blank" style={{ color: '#faad14' }}>
                管理后台
              </Link>
            )}
          </div>
        </div>
      </Header>

      <Content style={{ padding: '24px 0' }}>
        <div className="container">
          <Outlet />
        </div>
      </Content>

      <Footer style={{ textAlign: 'center', background: '#fff', marginTop: 24 }}>
        多商家商城平台 ©{new Date().getFullYear()}
      </Footer>
    </Layout>
  )
}

export default MainLayout
