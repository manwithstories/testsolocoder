import React from 'react';
import { Layout, Menu, Dropdown, Avatar, Badge, Button } from 'antd';
import {
  HomeOutlined,
  AppstoreOutlined,
  ShoppingCartOutlined,
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  NotificationOutlined,
  BarChartOutlined,
  PrinterOutlined,
  UploadOutlined,
  HeartOutlined,
  WalletOutlined,
  InboxOutlined,
  DashboardOutlined,
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { logout } from '@/store/authSlice';
import { useAuth } from '@/hooks/useAuth';
import { useNotifications } from '@/hooks/useNotifications';
import { getRoleText } from '@/utils/format';
import type { MenuProps } from 'antd';

const { Header, Sider, Content } = Layout;

interface MainLayoutProps {
  children: React.ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  const { user } = useAppSelector((state) => state.auth);
  const { userRole, userId } = useAuth();
  const { unreadCount, notifications, markAsRead } = useNotifications();

  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  const getMenuItems = (): MenuProps['items'] => {
    const items: MenuProps['items'] = [
      {
        key: '/',
        icon: <HomeOutlined />,
        label: '首页',
        onClick: () => navigate('/'),
      },
      {
        key: '/models',
        icon: <AppstoreOutlined />,
        label: '模型市场',
        onClick: () => navigate('/models'),
      },
    ];

    if (userRole === 'customer' || !userRole) {
      items.push(
        {
          key: '/my-orders',
          icon: <ShoppingCartOutlined />,
          label: '我的订单',
          onClick: () => navigate('/my-orders'),
        },
        {
          key: '/favorites',
          icon: <HeartOutlined />,
          label: '我的收藏',
          onClick: () => navigate('/favorites'),
        },
        {
          key: '/purchases',
          icon: <InboxOutlined />,
          label: '已购模型',
          onClick: () => navigate('/purchases'),
        }
      );
    }

    if (userRole === 'designer') {
      items.push(
        {
          key: '/my-models',
          icon: <UploadOutlined />,
          label: '我的模型',
          onClick: () => navigate('/my-models'),
        },
        {
          key: '/model-upload',
          icon: <UploadOutlined />,
          label: '上传模型',
          onClick: () => navigate('/model-upload'),
        }
      );
    }

    if (userRole === 'printer') {
      items.push(
        {
          key: '/printer-orders',
          icon: <PrinterOutlined />,
          label: '打印订单',
          onClick: () => navigate('/printer-orders'),
        },
        {
          key: '/printer-devices',
          icon: <DashboardOutlined />,
          label: '设备管理',
          onClick: () => navigate('/printer-devices'),
        },
        {
          key: '/printer-inventory',
          icon: <InboxOutlined />,
          label: '材料库存',
          onClick: () => navigate('/printer-inventory'),
        },
        {
          key: '/printer-schedules',
          icon: <BarChartOutlined />,
          label: '排产调度',
          onClick: () => navigate('/printer-schedules'),
        }
      );
    }

    if (userRole === 'designer' || userRole === 'printer' || userRole === 'admin') {
      items.push({
        key: '/stats',
        icon: <BarChartOutlined />,
        label: '数据统计',
        onClick: () => navigate('/stats'),
      });
    }

    items.push(
      {
        key: '/wallet',
        icon: <WalletOutlined />,
        label: '我的钱包',
        onClick: () => navigate('/wallet'),
      },
      {
        key: '/profile',
        icon: <UserOutlined />,
        label: '个人中心',
        onClick: () => navigate('/profile'),
      }
    );

    return items;
  };

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心',
      onClick: () => navigate('/profile'),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '账号设置',
      onClick: () => navigate('/settings'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  const notificationMenuItems: MenuProps['items'] = notifications.slice(0, 5).map((n) => ({
    key: n.id,
    label: (
      <div onClick={() => markAsRead(n.id)}>
        <div className={`font-medium ${!n.is_read ? 'text-primary-600' : ''}`}>{n.title}</div>
        <div className="text-sm text-gray-500 truncate">{n.content}</div>
      </div>
    ),
  }));

  return (
    <Layout className="min-h-screen">
      <Header className="bg-white shadow-sm px-6 flex items-center justify-between">
        <div className="flex items-center gap-8">
          <div
            className="text-xl font-bold text-primary-600 cursor-pointer"
            onClick={() => navigate('/')}
          >
            3D打印平台
          </div>
        </div>

        <div className="flex items-center gap-4">
          {userRole && (
            <div className="text-sm text-gray-600 px-3 py-1 bg-primary-50 rounded-full">
              {getRoleText(userRole)}
            </div>
          )}

          <Dropdown
            menu={{ items: notificationMenuItems }}
            placement="bottomRight"
            trigger={['click']}
          >
            <Badge count={unreadCount} size="small">
              <Button type="text" icon={<NotificationOutlined />} size="large" />
            </Badge>
          </Dropdown>

          <Dropdown
            menu={{ items: userMenuItems }}
            placement="bottomRight"
            trigger={['click']}
          >
            <div className="flex items-center gap-2 cursor-pointer hover:bg-gray-50 px-3 py-2 rounded-lg">
              <Avatar size="small" icon={<UserOutlined />} src={user?.avatar} />
              <span className="text-sm font-medium">{user?.username || '用户'}</span>
            </div>
          </Dropdown>
        </div>
      </Header>

      <Layout>
        <Sider width={220} className="bg-white border-r">
          <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={getMenuItems()}
            className="h-full border-r-0"
          />
        </Sider>

        <Content className="bg-gray-50 p-6">
          <div className="max-w-7xl mx-auto">{children}</div>
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;
