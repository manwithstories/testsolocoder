import { Layout, Menu, Button, Dropdown, Avatar } from 'antd'
import {
  ShopOutlined,
  ProductOutlined,
  ShoppingOutlined,
  BarChartOutlined,
  LogoutOutlined,
  UserOutlined,
} from '@ant-design/icons'
import { Link, Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAppStore } from '@/store'

const { Header, Sider, Content } = Layout

const SellerLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAppStore()

  const menuItems = [
    {
      key: '/seller/shop',
      icon: <ShopOutlined />,
      label: '店铺管理',
      onClick: () => navigate('/seller/shop'),
    },
    {
      key: '/seller/products',
      icon: <ProductOutlined />,
      label: '商品管理',
      onClick: () => navigate('/seller/products'),
    },
    {
      key: '/seller/orders',
      icon: <ShoppingOutlined />,
      label: '订单管理',
      onClick: () => navigate('/seller/orders'),
    },
    {
      key: '/seller/statistics',
      icon: <BarChartOutlined />,
      label: '数据统计',
      onClick: () => navigate('/seller/statistics'),
    },
  ]

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#001529', padding: '0 24px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Link to="/" style={{ color: '#fff', fontSize: 18, fontWeight: 'bold' }}>
          <ShopOutlined /> 商家中心
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

export default SellerLayout
