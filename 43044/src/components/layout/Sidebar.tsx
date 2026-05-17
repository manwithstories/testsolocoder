import { NavLink } from 'react-router-dom'
import { Target, Calendar, BarChart3, Tag } from 'lucide-react'

export function Sidebar() {
  const navItems = [
    { to: '/', label: '目标管理', icon: Target },
    { to: '/calendar', label: '日历视图', icon: Calendar },
    { to: '/stats', label: '统计分析', icon: BarChart3 },
    { to: '/tags', label: '标签管理', icon: Tag },
  ]
  
  return (
    <aside className="w-64 bg-white border-r border-gray-200 h-screen flex-shrink-0">
      <div className="p-6 border-b border-gray-200">
        <h1 className="text-xl font-bold text-gray-800 flex items-center gap-2">
          <Target className="w-6 h-6 text-blue-600" />
          目标追踪
        </h1>
        <p className="text-sm text-gray-500 mt-1">Goal Tracker</p>
      </div>
      <nav className="p-4">
        <ul className="space-y-2">
          {navItems.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                end={item.to === '/'}
                className={({ isActive }) =>
                  `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                    isActive
                      ? 'bg-blue-50 text-blue-600'
                      : 'text-gray-600 hover:bg-gray-50'
                  }`
                }
              >
                <item.icon className="w-5 h-5" />
                <span className="font-medium">{item.label}</span>
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  )
}
