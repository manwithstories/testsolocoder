import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { UserCog, FileText, Search, Briefcase, Calendar, Star, Bell, LogOut, Menu, X, User } from 'lucide-react'
import { useState } from 'react'
import { useAppDispatch, useAppSelector } from '@/hooks/redux'
import { logout, selectAuth } from '@/store/slices/authSlice'
import { selectUnreadCount } from '@/store/slices/notificationSlice'

const navItems = [
  { path: 'dashboard', label: 'Dashboard', icon: UserCog },
  { path: 'jobs', label: 'Search Jobs', icon: Search },
  { path: 'resumes', label: 'My Resumes', icon: FileText },
  { path: 'applications', label: 'Applications', icon: Briefcase },
  { path: 'interviews', label: 'Interviews', icon: Calendar },
  { path: 'reviews', label: 'Reviews', icon: Star },
]

export default function JobSeekerLayout() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()
  const { user } = useAppSelector(selectAuth)
  const unreadCount = useAppSelector(selectUnreadCount)
  const [sidebarOpen, setSidebarOpen] = useState(false)

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  return (
    <div className="flex h-screen bg-gray-50">
      <aside className={`fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-lg transform transition-transform lg:translate-x-0 lg:relative ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}`}>
        <div className="flex items-center justify-between p-4 border-b">
          <div className="flex items-center gap-2">
            <UserCog className="w-8 h-8 text-primary-600" />
            <span className="text-xl font-bold">JobFinder</span>
          </div>
          <button className="lg:hidden" onClick={() => setSidebarOpen(false)}>
            <X className="w-6 h-6" />
          </button>
        </div>

        <nav className="p-4 space-y-1">
          {navItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `flex items-center gap-3 px-4 py-2 rounded-lg transition-colors ${
                  isActive
                    ? 'bg-primary-100 text-primary-700'
                    : 'text-gray-600 hover:bg-gray-100'
                }`
              }
            >
              <item.icon className="w-5 h-5" />
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="bg-white shadow-sm px-6 py-4 flex items-center justify-between">
          <button className="lg:hidden" onClick={() => setSidebarOpen(true)}>
            <Menu className="w-6 h-6" />
          </button>
          <h1 className="text-xl font-semibold text-gray-800">{user?.name}</h1>
          <div className="flex items-center gap-4">
            <NavLink to="profile" className="flex items-center gap-2 text-gray-600 hover:text-primary-600">
              <User className="w-5 h-5" />
              <span className="hidden sm:inline">Profile</span>
            </NavLink>
            <button className="relative p-2 text-gray-600 hover:text-primary-600">
              <Bell className="w-5 h-5" />
              {unreadCount > 0 && (
                <span className="absolute -top-1 -right-1 w-5 h-5 bg-red-500 text-white text-xs rounded-full flex items-center justify-center">
                  {unreadCount}
                </span>
              )}
            </button>
            <button onClick={handleLogout} className="text-gray-600 hover:text-red-600">
              <LogOut className="w-5 h-5" />
            </button>
          </div>
        </header>

        <main className="flex-1 overflow-auto p-6">
          <Outlet />
          {sidebarOpen && (
            <div className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden" onClick={() => setSidebarOpen(false)} />
          )}
        </main>
      </div>
    </div>
  )
}
