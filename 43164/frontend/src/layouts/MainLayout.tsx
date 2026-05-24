import { Outlet, useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'
import {
  Home,
  Users,
  Calendar,
  Wallet,
  MessageSquare,
  User,
  BarChart3,
  LogOut,
  Bell,
  GraduationCap,
  BookOpen,
  Settings,
} from 'lucide-react'
import { useState } from 'react'
import NotificationDropdown from '@/components/NotificationDropdown'

export default function MainLayout() {
  const { user, logout, isAuthenticated } = useAuthStore()
  const navigate = useNavigate()
  const [showNotifications, setShowNotifications] = useState(false)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  if (!isAuthenticated) {
    return <Outlet />
  }

  const navItems = [
    { icon: Home, label: '首页', path: '/' },
    { icon: Users, label: '找老师', path: '/teachers' },
    { icon: Calendar, label: '我的课程', path: '/my-bookings' },
    { icon: MessageSquare, label: '消息', path: '/messages' },
    { icon: Wallet, label: '钱包', path: '/wallet' },
  ]

  if (user?.role === 'teacher') {
    navItems.push({ icon: BarChart3, label: '教师中心', path: '/teacher/dashboard' })
  } else if (user?.role === 'student') {
    navItems.push({ icon: BookOpen, label: '学习进度', path: '/learning' })
  } else if (user?.role === 'admin') {
    navItems.push({ icon: Settings, label: '管理后台', path: '/admin/dashboard' })
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200 sticky top-0 z-40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center space-x-2 cursor-pointer" onClick={() => navigate('/')}>
              <GraduationCap className="h-8 w-8 text-primary-600" />
              <span className="text-xl font-bold text-gray-900">在线家教平台</span>
            </div>

            <nav className="hidden md:flex space-x-1">
              {navItems.map((item) => (
                <button
                  key={item.path}
                  onClick={() => navigate(item.path)}
                  className="flex items-center px-3 py-2 text-sm font-medium text-gray-600 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
                >
                  <item.icon className="h-4 w-4 mr-2" />
                  {item.label}
                </button>
              ))}
            </nav>

            <div className="flex items-center space-x-4">
              <div className="relative">
                <button
                  onClick={() => setShowNotifications(!showNotifications)}
                  className="p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg"
                >
                  <Bell className="h-5 w-5" />
                </button>
                {showNotifications && (
                  <NotificationDropdown onClose={() => setShowNotifications(false)} />
                )}
              </div>

              <div className="relative group">
                <button className="flex items-center space-x-2 p-1 hover:bg-gray-100 rounded-lg">
                  {user?.avatarUrl ? (
                    <img src={user.avatarUrl} alt="" className="h-8 w-8 rounded-full" />
                  ) : (
                    <div className="h-8 w-8 rounded-full bg-primary-600 flex items-center justify-center text-white text-sm font-medium">
                      {user?.firstName?.[0]}{user?.lastName?.[0]}
                    </div>
                  )}
                </button>

                <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all">
                  <button
                    onClick={() => navigate('/profile')}
                    className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center"
                  >
                    <User className="h-4 w-4 mr-2" />
                    个人资料
                  </button>
                  <hr className="my-1" />
                  <button
                    onClick={handleLogout}
                    className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 flex items-center"
                  >
                    <LogOut className="h-4 w-4 mr-2" />
                    退出登录
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <Outlet />
      </main>
    </div>
  )
}
