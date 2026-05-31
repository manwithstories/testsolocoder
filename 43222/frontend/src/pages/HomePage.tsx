import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { Sprout, Leaf, Calendar, Bug, MessageSquare, RefreshCw, ShoppingBag } from 'lucide-react'
import { plotApi, plantingRecordApi, calendarApi, diseaseApi } from '@/api'
import { useAuthStore } from '@/store/auth'

export default function HomePage() {
  const user = useAuthStore((state) => state.user)

  const { data: plotsData } = useQuery({
    queryKey: ['plots'],
    queryFn: () => plotApi.getAll({ page_size: 5 }),
  })

  const { data: recordsData } = useQuery({
    queryKey: ['planting-records'],
    queryFn: () => plantingRecordApi.getAll({ page_size: 5 }),
  })

  const { data: calendarData } = useQuery({
    queryKey: ['calendar-events'],
    queryFn: () => calendarApi.getEvents(),
  })

  const { data: recommendations } = useQuery({
    queryKey: ['recommendations'],
    queryFn: () => calendarApi.getRecommendations(),
  })

  const stats = [
    { label: '我的菜园', value: plotsData?.data?.total || 0, icon: Sprout, href: '/plots', color: 'bg-garden-100 text-garden-600' },
    { label: '种植记录', value: recordsData?.data?.total || 0, icon: Leaf, href: '/growth', color: 'bg-blue-100 text-blue-600' },
    { label: '日历事件', value: calendarData?.data?.events?.length || 0, icon: Calendar, href: '/calendar', color: 'bg-amber-100 text-amber-600' },
    { label: '诊断记录', value: 0, icon: Bug, href: '/disease', color: 'bg-red-100 text-red-600' },
  ]

  const quickLinks = [
    { name: '社区交流', icon: MessageSquare, href: '/community', color: 'bg-purple-100 text-purple-600' },
    { name: '种子交换', icon: RefreshCw, href: '/exchange', color: 'bg-pink-100 text-pink-600' },
    { name: '园艺商城', icon: ShoppingBag, href: '/shop', color: 'bg-orange-100 text-orange-600' },
  ]

  return (
    <div className="space-y-6">
      {/* Welcome */}
      <div className="bg-gradient-to-r from-garden-600 to-garden-500 rounded-2xl p-6 text-white">
        <h1 className="text-2xl font-bold mb-2">
          欢迎回来，{user?.nickname || user?.username}！
        </h1>
        <p className="text-garden-100">今天是种植的好日子，让我们开始吧 🌱</p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <Link key={stat.label} to={stat.href} className="card hover:shadow-md transition-shadow">
            <div className="card-body flex items-center gap-4">
              <div className={`p-3 rounded-xl ${stat.color}`}>
                <stat.icon className="w-6 h-6" />
              </div>
              <div>
                <p className="text-2xl font-bold text-gray-900">{stat.value}</p>
                <p className="text-sm text-gray-500">{stat.label}</p>
              </div>
            </div>
          </Link>
        ))}
      </div>

      {/* Quick Links */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {quickLinks.map((link) => (
          <Link key={link.name} to={link.href} className="card hover:shadow-md transition-shadow">
            <div className="card-body flex items-center gap-4">
              <div className={`p-3 rounded-xl ${link.color}`}>
                <link.icon className="w-6 h-6" />
              </div>
              <span className="font-medium text-gray-900">{link.name}</span>
            </div>
          </Link>
        ))}
      </div>

      {/* Recommendations */}
      {recommendations?.data?.recommendations && recommendations.data.recommendations.length > 0 && (
        <div className="card">
          <div className="card-header flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">本月推荐种植</h2>
            <Link to="/plants" className="text-sm text-garden-600 hover:text-garden-700">
              查看全部
            </Link>
          </div>
          <div className="card-body">
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
              {recommendations.data.recommendations.slice(0, 5).map((plant: any) => (
                <Link
                  key={plant.id}
                  to={`/plants/${plant.id}`}
                  className="group text-center"
                >
                  <div className="aspect-square bg-garden-50 rounded-xl flex items-center justify-center mb-2 group-hover:bg-garden-100 transition-colors">
                    {plant.image_url ? (
                      <img
                        src={plant.image_url}
                        alt={plant.name}
                        className="w-full h-full object-cover rounded-xl"
                      />
                    ) : (
                      <Leaf className="w-12 h-12 text-garden-400" />
                    )}
                  </div>
                  <p className="text-sm font-medium text-gray-900 truncate">{plant.name}</p>
                  <p className="text-xs text-gray-500">{plant.category}</p>
                </Link>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Recent Records */}
      {recordsData?.data?.records && recordsData.data.records.length > 0 && (
        <div className="card">
          <div className="card-header flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">最近种植</h2>
            <Link to="/growth" className="text-sm text-garden-600 hover:text-garden-700">
              查看全部
            </Link>
          </div>
          <div className="divide-y divide-gray-100">
            {recordsData.data.records.slice(0, 5).map((record: any) => (
              <div key={record.id} className="px-6 py-4 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-garden-100 rounded-lg flex items-center justify-center">
                    <Sprout className="w-5 h-5 text-garden-600" />
                  </div>
                  <div>
                    <p className="font-medium text-gray-900">{record.plant?.name}</p>
                    <p className="text-sm text-gray-500">
                      {new Date(record.planting_date).toLocaleDateString('zh-CN')}
                    </p>
                  </div>
                </div>
                <span className={`badge ${
                  record.status === 'harvested'
                    ? 'bg-green-100 text-green-700'
                    : record.status === 'growing'
                    ? 'bg-blue-100 text-blue-700'
                    : 'bg-gray-100 text-gray-700'
                }`}>
                  {record.status === 'planted' ? '已种植' : record.status === 'growing' ? '生长中' : record.status === 'harvested' ? '已收获' : record.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
