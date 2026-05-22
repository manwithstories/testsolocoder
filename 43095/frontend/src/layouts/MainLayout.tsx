import { useState } from 'react'
import { Layout, Menu, Avatar, Dropdown, Badge, Button } from 'antd'
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  UserOutlined,
  CalendarOutlined,
  FileTextOutlined,
  BellOutlined,
  PayCircleOutlined,
  SettingOutlined,
  LogoutOutlined,
  TeamOutlined,
  MedicineBoxOutlined,
  ScheduleOutlined
} from '@ant-design/icons'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { useAuthStore, useAppStore } from '@/store'
import type { UserRole } from '@/types'

const { Header, Sider, Content } = Layout

const getMenuItems = (role: UserRole) => {
  const patientMenu = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '首页'
    },
    {
      key: '/doctors',
      icon: <UserOutlined />,
      label: '医生列表'
    },
    {
      key: '/appointments',
      icon: <CalendarOutlined />,
      label: '预约管理'
    },
    {
      key: '/health-records',
      icon: <FileTextOutlined />,
      label: '健康档案'
    },
    {
      key: '/notifications',
      icon: <BellOutlined />,
      label: '通知中心'
    },
    {
      key: '/payments',
      icon: <PayCircleOutlined />,
      label: '费用管理'
    }
  ]

  const doctorMenu = [
    {
      key: '/doctor',
      icon: <DashboardOutlined />,
      label: '工作台'
    },
    {
      key: '/doctor/schedule',
      icon: <ScheduleOutlined />,
      label: '我的排班'
    },
    {
      key: '/doctor/appointments',
      icon: <CalendarOutlined />,
      label: '我的预约'
    },
    {
      key: '/doctor/patients',
      icon: <TeamOutlined />,
      label: '我的患者'
    }
  ]

  const adminMenu = [
    {
      key: '/admin',
      icon: <DashboardOutlined />,
      label: '管理控制台'
    },
    {
      key: '/admin/departments',
      icon: <MedicineBoxOutlined />,
      label: '科室管理'
    },
    {
      key: '/admin/doctors',
      icon: <UserOutlined />,
      label: '医生管理'
    },
    {
      key: '/admin/patients',
      icon: <TeamOutlined />,
      label: '患者管理'
    }
  ]

  switch (role) {
    case 'admin':
      return adminMenu
    case 'doctor':
      return doctorMenu
    case 'patient':
    default:
      return patientMenu
  }
}

const MainLayout = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuthStore()
  const { collapsed, toggleCollapsed } = useAppStore()
  const [selectedKey, setSelectedKey] = useState(location.pathname)

  const handleMenuClick = ({ key }: { key: string }) => {
    setSelectedKey(key)
    navigate(key)
  }

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人中心'
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置'
    },
    {
      type: 'divider' as const
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout
    }
  ]

  return (
    <Layout className="min-h-screen">
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme="light"
        className="shadow-lg"
        width={240}
      >
        <div className="flex h-16 items-center justify-center border-b">
          <h1 className={`font-bold text-blue-600 transition-all ${collapsed ? 'text-lg' : 'text-xl'}`}>
            {collapsed ? 'HMS' : '医疗管理系统'}
          </h1>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[selectedKey]}
          items={user ? getMenuItems(user.role) : []}
          onClick={handleMenuClick}
          className="border-none"
        />
      </Sider>
      <Layout>
        <Header className="flex items-center justify-between bg-white px-6 shadow-sm">
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={toggleCollapsed}
            className="text-lg"
          />
          <div className="flex items-center gap-4">
            <Badge count={5} size="small">
              <Button type="text" icon={<BellOutlined />} className="text-lg" />
            </Badge>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <div className="flex cursor-pointer items-center gap-2 hover:bg-gray-50 rounded-lg px-2 py-1">
                <Avatar src={user?.avatar_url} icon={<UserOutlined />} />
                <span className="hidden md:inline">{user?.full_name || user?.username}</span>
              </div>
            </Dropdown>
          </div>
        </Header>
        <Content className="m-6 overflow-auto">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
