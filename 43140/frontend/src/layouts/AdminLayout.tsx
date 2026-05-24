import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { Shield, Users, LogOut } from 'lucide-react'
import { useAppDispatch } from '@/hooks/redux'
import { logout } from '@/store/slices/authSlice'

const navItems = [
  { path: '/admin/dashboard', label: 'Dashboard' },
  { path: '/admin/users', label: 'Users' },
]

export default function AdminLayout() {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  return (
    <div className="flex h-screen bg-gray-50">
      <aside className="w-64 bg-white shadow-lg">
        <div className="p-4 border-b">
          <div className="flex items-center gap-2">
            <Shield className="w-8 h-8 text-red-600" />
            <span className="text-xl font-bold">Admin Panel</span>
          </div>
        </div>

        <nav className="p-4 space-y-1">
          {navItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `flex items-center gap-3 px-4 py-2 rounded-lg transition-colors ${
                  isActive
                    ? 'bg-red-100 text-red-700'
                    : 'text-gray-600 hover:bg-gray-100'
                }`
              }
            >
              <Users className="w-5 h-5" />
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="bg-white shadow-sm px-6 py-4 flex items-center justify-between">
          <h1 className="text-xl font-semibold text-gray-800">Admin Dashboard</h1>
          <button onClick={handleLogout} className="text-gray-600 hover:text-red-600">
            <LogOut className="w-5 h-5" />
          </button>
        </header>

        <main className="flex-1 overflow-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
