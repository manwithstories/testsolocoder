import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Layout, Menu, Avatar, Dropdown, Button, Badge } from 'antd';
import {
  BookOutlined,
  HomeOutlined,
  UserOutlined,
  ShoppingCartOutlined,
  MessageOutlined,
  FileTextOutlined,
  LogoutOutlined,
  SettingOutlined,
  DashboardOutlined,
} from '@ant-design/icons';
import { useAuthStore } from '../context/authStore';

const { Header, Content, Sider } = Layout;

interface MainLayoutProps {
  children: React.ReactNode;
}

export const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const navigate = useNavigate();
  const { user, isAuthenticated, logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };

  const userMenu = (
    <Menu>
      <Menu.Item key="profile" icon={<UserOutlined />}>
        <Link to="/profile">个人中心</Link>
      </Menu.Item>
      <Menu.Item key="settings" icon={<SettingOutlined />}>
        <Link to="/settings">设置</Link>
      </Menu.Item>
      <Menu.Divider />
      <Menu.Item key="logout" icon={<LogoutOutlined />} onClick={handleLogout}>
        退出登录
      </Menu.Item>
    </Menu>
  );

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: <Link to="/">首页</Link>,
    },
    {
      key: '/textbooks',
      icon: <BookOutlined />,
      label: <Link to="/textbooks">教材市场</Link>,
    },
    {
      key: '/notes',
      icon: <FileTextOutlined />,
      label: <Link to="/notes">笔记共享</Link>,
    },
    {
      key: '/orders',
      icon: <ShoppingCartOutlined />,
      label: <Link to="/orders">我的订单</Link>,
    },
    {
      key: '/messages',
      icon: <MessageOutlined />,
      label: <Link to="/messages">消息中心</Link>,
    },
  ];

  return (
    <Layout className="min-h-screen">
      <Header className="bg-white shadow-md px-6 flex items-center justify-between fixed w-full z-10">
        <div className="flex items-center gap-8">
          <Link to="/" className="text-xl font-bold text-blue-600">
            校园教材交易平台
          </Link>
          <Menu
            mode="horizontal"
            selectedKeys={[window.location.pathname]}
            items={menuItems}
            className="border-0 flex-1"
          />
        </div>
        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <>
              <Badge count={0} size="small">
                <Button type="text" icon={<MessageOutlined />} onClick={() => navigate('/messages')} />
              </Badge>
              <Dropdown overlay={userMenu} placement="bottomRight">
                <div className="flex items-center gap-2 cursor-pointer">
                  <Avatar src={user?.avatar} icon={<UserOutlined />} />
                  <span className="hidden md:inline">{user?.username}</span>
                </div>
              </Dropdown>
            </>
          ) : (
            <div className="flex gap-2">
              <Button onClick={() => navigate('/login')}>登录</Button>
              <Button type="primary" onClick={() => navigate('/register')}>
                注册
              </Button>
            </div>
          )}
        </div>
      </Header>
      <Content className="pt-16">{children}</Content>
    </Layout>
  );
};

interface AdminLayoutProps {
  children: React.ReactNode;
}

export const AdminLayout: React.FC<AdminLayoutProps> = ({ children }) => {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };

  const adminMenuItems = [
    {
      key: '/admin/dashboard',
      icon: <DashboardOutlined />,
      label: <Link to="/admin/dashboard">数据概览</Link>,
    },
    {
      key: '/admin/users',
      icon: <UserOutlined />,
      label: <Link to="/admin/users">用户管理</Link>,
    },
    {
      key: '/admin/textbooks',
      icon: <BookOutlined />,
      label: <Link to="/admin/textbooks">教材管理</Link>,
    },
    {
      key: '/admin/notes',
      icon: <FileTextOutlined />,
      label: <Link to="/admin/notes">笔记管理</Link>,
    },
    {
      key: '/admin/orders',
      icon: <ShoppingCartOutlined />,
      label: <Link to="/admin/orders">订单管理</Link>,
    },
    {
      key: '/admin/reviews',
      icon: <MessageOutlined />,
      label: <Link to="/admin/reviews">评价管理</Link>,
    },
  ];

  return (
    <Layout className="min-h-screen">
      <Sider theme="light" width={220} className="fixed h-full shadow-md">
        <div className="p-4 text-center border-b">
          <h2 className="text-lg font-bold text-blue-600">管理后台</h2>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[window.location.pathname]}
          items={adminMenuItems}
          className="mt-4"
        />
      </Sider>
      <Layout className="ml-[220px]">
        <Header className="bg-white shadow-md flex justify-end items-center px-6">
          <Dropdown
            overlay={
              <Menu>
                <Menu.Item key="logout" icon={<LogoutOutlined />} onClick={handleLogout}>
                  退出登录
                </Menu.Item>
              </Menu>
            }
          >
            <div className="flex items-center gap-2 cursor-pointer">
              <Avatar src={user?.avatar} icon={<UserOutlined />} />
              <span>{user?.username}</span>
            </div>
          </Dropdown>
        </Header>
        <Content className="p-6">{children}</Content>
      </Layout>
    </Layout>
  );
};
