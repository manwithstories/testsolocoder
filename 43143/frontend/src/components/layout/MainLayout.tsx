import React, { useState } from 'react';
import { Layout, Menu, Avatar, Dropdown, Badge } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import {
  HomeOutlined,
  AppstoreOutlined,
  CalendarOutlined,
  MessageOutlined,
  UserOutlined,
  WalletOutlined,
  BarChartOutlined,
  LogoutOutlined,
  SettingOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { RootState, logout } from '@/store';
import { messageApi } from '@/api/message';

const { Header, Sider, Content } = Layout;

const MainLayout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch();
  const user = useSelector((state: RootState) => state.auth.user);
  const [collapsed, setCollapsed] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);

  React.useEffect(() => {
    const fetchUnreadCount = async () => {
      try {
        const data = await messageApi.getUnreadCount();
        setUnreadCount(data.unread_count);
      } catch (error) {
        console.error('获取未读消息数失败:', error);
      }
    };

    fetchUnreadCount();
    const interval = setInterval(fetchUnreadCount, 30000);
    return () => clearInterval(interval);
  }, []);

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: '首页',
    },
    {
      key: '/skills',
      icon: <AppstoreOutlined />,
      label: '技能库',
    },
    {
      key: '/postings/create',
      icon: <PlusOutlined />,
      label: '发布技能',
    },
    {
      key: '/bookings',
      icon: <CalendarOutlined />,
      label: '我的预约',
    },
    {
      key: '/messages',
      icon: <MessageOutlined />,
      label: '消息中心',
    },
    {
      key: '/schedule',
      icon: <CalendarOutlined />,
      label: '日程管理',
    },
    {
      key: '/wallet',
      icon: <WalletOutlined />,
      label: '我的钱包',
    },
    {
      key: '/statistics',
      icon: <BarChartOutlined />,
      label: '数据统计',
    },
    {
      key: '/profile',
      icon: <UserOutlined />,
      label: '个人中心',
    },
  ];

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人资料',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'edit',
      icon: <SettingOutlined />,
      label: '编辑资料',
      onClick: () => navigate('/profile/edit'),
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
  ];

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        theme="dark"
      >
        <div
          style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: '#fff',
            fontSize: collapsed ? 14 : 18,
            fontWeight: 'bold',
          }}
        >
          {collapsed ? 'Skill' : '技能共享平台'}
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
        <Header
          style={{
            background: '#fff',
            padding: '0 24px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
          }}
        >
          <div style={{ fontSize: 18, fontWeight: 500 }}>
            {menuItems.find((item) => item.key === location.pathname)?.label || '首页'}
          </div>
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div
              style={{
                display: 'flex',
                alignItems: 'center',
                cursor: 'pointer',
                gap: 12,
              }}
            >
              <Badge count={unreadCount} size="small" offset={[-4, 4]}>
                <MessageOutlined
                  style={{ fontSize: 20, cursor: 'pointer' }}
                  onClick={() => navigate('/messages')}
                />
              </Badge>
              <Avatar src={user?.avatar} icon={<UserOutlined />} />
              <span>{user?.nickname || '用户'}</span>
            </div>
          </Dropdown>
        </Header>
        <Content
          style={{
            margin: 16,
            padding: 24,
            background: '#fff',
            minHeight: 280,
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;
