import { NavLink } from 'react-router-dom'
import { useEffect, useRef } from 'react'
import {
  Home,
  Clock,
  StickyNote,
  ListMusic,
  BarChart3,
  Download,
  Radio,
  Bell,
  BellOff,
} from 'lucide-react'
import { useGetNewEpisodesCountQuery } from '@/app/services/api'
import { requestNotificationPermission, showNotification } from '@/utils/notifications'

export default function Sidebar() {
  const { data: newEpisodes } = useGetNewEpisodesCountQuery(undefined, {
    pollingInterval: 60000,
  })
  const prevCountRef = useRef<number>(0)
  const notificationEnabledRef = useRef<boolean>(false)

  useEffect(() => {
    const initNotifications = async () => {
      const enabled = await requestNotificationPermission()
      notificationEnabledRef.current = enabled
    }
    initNotifications()
  }, [])

  useEffect(() => {
    const currentCount = newEpisodes?.new_episodes_count || 0
    const prevCount = prevCountRef.current

    if (currentCount > prevCount && prevCount > 0 && notificationEnabledRef.current) {
      const newEpisodesCount = currentCount - prevCount
      showNotification(
        '播客更新提醒',
        {
          body: `有 ${newEpisodesCount} 个新剧集可供收听`,
          tag: 'new-episodes',
        } as NotificationOptions
      )
    }

    prevCountRef.current = currentCount
  }, [newEpisodes?.new_episodes_count])

  const navItems = [
    { path: '/', icon: Home, label: '播客订阅', badge: newEpisodes?.new_episodes_count },
    { path: '/history', icon: Clock, label: '收听历史' },
    { path: '/notes', icon: StickyNote, label: '笔记' },
    { path: '/playlists', icon: ListMusic, label: '播放列表' },
    { path: '/stats', icon: BarChart3, label: '数据统计' },
    { path: '/import-export', icon: Download, label: '导入导出' },
  ]

  return (
    <aside className="fixed left-0 top-0 h-full w-64 bg-white border-r border-gray-200 shadow-sm">
      <div className="p-6 border-b border-gray-200">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-indigo-600 rounded-xl flex items-center justify-center">
            <Radio className="w-6 h-6 text-white" />
          </div>
          <h1 className="text-xl font-bold text-gray-900">播客管理器</h1>
        </div>
      </div>

      <nav className="p-4 space-y-1">
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            end={item.path === '/'}
            className={({ isActive }) =>
              `flex items-center gap-3 px-4 py-3 rounded-lg transition-colors duration-200 ${
                isActive
                  ? 'bg-indigo-50 text-indigo-600'
                  : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
              }`
            }
          >
            <item.icon className="w-5 h-5" />
            <span className="font-medium">{item.label}</span>
            {item.badge && item.badge > 0 && (
              <span className="ml-auto bg-red-500 text-white text-xs px-2 py-0.5 rounded-full">
                {item.badge}
              </span>
            )}
          </NavLink>
        ))}
      </nav>

      <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-200">
        <button
          onClick={requestNotificationPermission}
          className="w-full flex items-center gap-3 px-4 py-3 text-gray-600 hover:bg-gray-50 rounded-lg transition-colors"
          title="桌面通知设置"
        >
          {notificationEnabledRef.current ? (
            <Bell className="w-5 h-5 text-green-500" />
          ) : (
            <BellOff className="w-5 h-5 text-gray-400" />
          )}
          <span className="font-medium text-sm">
            {notificationEnabledRef.current ? '通知已开启' : '开启桌面通知'}
          </span>
        </button>
      </div>
    </aside>
  )
}
