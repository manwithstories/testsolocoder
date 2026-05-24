import { useState, useEffect } from 'react';
import { Layout as AntLayout, Menu, Avatar, Dropdown, Badge, theme } from 'antd';
import type { MenuProps } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  UnorderedListOutlined,
  CalendarOutlined,
  ScanOutlined,
  DollarOutlined,
  StarOutlined,
  BarChartOutlined,
  UserOutlined,
  LogoutOutlined,
  BellOutlined,
} from '@ant-design/icons';
import { useAuthStore } from '../context/AuthContext';

const { Header, Sider, Content } = AntLayout;

type MenuItem = Required<MenuProps>['items'][number];

export const Layout = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout, isAuthenticated } = useAuthStore();
  const [collapsed, setCollapsed] = useState(false);
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, navigate]);

  const getMenuItems = (): MenuItem[] => {
    const items: MenuItem[] = [
      {
        key: '/dashboard',
        icon: <DashboardOutlined />,
        label: '仪表盘',
      },
    ];

    if (user?.role === 'employer' || user?.role === 'agent') {
      items.push({
        key: '/jobs',
        icon: <UnorderedListOutlined />,
        label: '招聘管理',
        children: [
          { key: '/jobs', label: '岗位列表' },
          { key: '/jobs/create', label: '发布岗位' },
          { key: '/jobs/my-applications', label: '申请记录' },
        ],
      });
    }

    if (user?.role === 'temporary') {
      items.push({
        key: '/jobs',
        icon: <UnorderedListOutlined />,
        label: '岗位浏览',
      });
    }

    items.push({
      key: '/schedules',
      icon: <CalendarOutlined />,
      label: '排班管理',
      children: [
        { key: '/schedules', label: '排班列表' },
        { key: '/schedules/board', label: '排班看板' },
      ],
    });

    items.push({
      key: '/checkins',
      icon: <ScanOutlined />,
      label: '签到管理',
      children: [
        { key: '/checkins', label: '签到打卡' },
        { key: '/checkins/records', label: '签到记录' },
      ],
    });

    items.push({
      key: '/salaries',
      icon: <DollarOutlined />,
      label: '薪资管理',
    });

    items.push({
      key: '/evaluations',
      icon: <StarOutlined />,
      label: '评价管理',
      children: [
        { key: '/evaluations', label: '评价列表' },
        { key: '/evaluations/mine', label: '我的评价' },
      ],
    });

    items.push({
      key: '/stats',
      icon: <BarChartOutlined />,
      label: '数据统计',
    });

    return items;
  };

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const userMenu = {
    items: [
      {
        key: 'profile',
        icon: <UserOutlined />,
        label: '个人中心',
        onClick: () => navigate('/profile'),
      },
      {
        key: 'logout',
        icon: <LogoutOutlined />,
        label: '退出登录',
        onClick: handleLogout,
      },
    ],
  };

  const getSelectedKeys = () => {
    const path = location.pathname;
    if (path.startsWith('/jobs')) return ['/jobs'];
    if (path.startsWith('/schedules')) return ['/schedules'];
    if (path.startsWith('/checkins')) return ['/checkins'];
    if (path.startsWith('/salaries')) return ['/salaries'];
    if (path.startsWith('/evaluations')) return ['/evaluations'];
    return [path];
  };

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
        theme="dark"
      >
        <div className="flex items-center justify-center h-16 text-white text-xl font-bold">
          {collapsed ? 'TSP' : '临时工管理平台'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={getSelectedKeys()}
          items={getMenuItems()}
          onClick={handleMenuClick}
        />
      </Sider>
      <AntLayout>
        <Header
          style={{ padding: '0 16px', background: colorBgContainer, display: 'flex', justifyContent: 'flex-end', alignItems: 'center', gap: '16px' }}
        >
          <Badge count={0} size="small">
            <BellOutlined style={{ fontSize: '20px', cursor: 'pointer' }} />
          </Badge>
          <Dropdown menu={userMenu} placement="bottomRight">
            <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: '8px' }}>
              <Avatar
                src={user?.avatar}
                icon={!user?.avatar && <UserOutlined />}
              />
              <span>{user?.real_name || user?.username}</span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{ margin: '16px' }}>
          <div style={{ background: colorBgContainer, padding: 24, minHeight: 'calc(100vh - 112px)', borderRadius: 8 }}>
            <Outlet />
          </div>
        </Content>
      </AntLayout>
    </AntLayout>
  );
};
