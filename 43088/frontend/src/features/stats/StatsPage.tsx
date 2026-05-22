import { useState } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
} from 'recharts'
import {
  useGetListeningStatsQuery,
  useGetPodcastDistributionQuery,
  useGetCompletionStatsQuery,
  useGetListeningHabitsQuery,
} from '@/app/services/api'
import { formatDuration } from '@/utils/format'

const COLORS = ['#6366f1', '#8b5cf6', '#a855f7', '#d946ef', '#ec4899', '#f43f5e', '#f97316', '#eab308']

export default function StatsPage() {
  const [days, setDays] = useState(30)

  const { data: listeningStats, isLoading: loadingStats } = useGetListeningStatsQuery(days)
  const { data: distribution, isLoading: loadingDistribution } = useGetPodcastDistributionQuery()
  const { data: completionStats, isLoading: loadingCompletion } = useGetCompletionStatsQuery()
  const { data: habits, isLoading: loadingHabits } = useGetListeningHabitsQuery()

  const chartData = listeningStats
    ? [...listeningStats].reverse().map((item) => ({
        ...item,
        duration: item.total_duration / 60,
      }))
    : []

  const pieData = distribution?.map((item) => ({
    name: item.title,
    value: item.total_duration,
  })) || []

  const hourData = habits?.hour_distribution?.map((item: any) => ({
    hour: `${item.hour}:00`,
    count: Number(item.count),
  })) || []

  if (loadingStats || loadingDistribution || loadingCompletion || loadingHabits) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-900">数据统计</h1>
        <select
          value={days}
          onChange={(e) => setDays(Number(e.target.value))}
          className="input w-auto"
        >
          <option value={7}>最近7天</option>
          <option value={30}>最近30天</option>
          <option value={90}>最近90天</option>
          <option value={365}>最近一年</option>
        </select>
      </div>

      {completionStats && (
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="card p-6">
            <p className="text-sm text-gray-500 mb-1">总剧集数</p>
            <p className="text-3xl font-bold text-gray-900">{completionStats.total_episodes}</p>
          </div>
          <div className="card p-6">
            <p className="text-sm text-gray-500 mb-1">已完成</p>
            <p className="text-3xl font-bold text-green-600">{completionStats.completed_episodes}</p>
          </div>
          <div className="card p-6">
            <p className="text-sm text-gray-500 mb-1">进行中</p>
            <p className="text-3xl font-bold text-orange-600">{completionStats.in_progress_episodes}</p>
          </div>
          <div className="card p-6">
            <p className="text-sm text-gray-500 mb-1">完成率</p>
            <p className="text-3xl font-bold text-indigo-600">
              {completionStats.completion_rate?.toFixed(1)}%
            </p>
          </div>
        </div>
      )}

      <div className="card p-6">
        <h3 className="text-lg font-semibold mb-4">收听时长趋势（分钟）</h3>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
              <XAxis dataKey="date" stroke="#6b7280" fontSize={12} />
              <YAxis stroke="#6b7280" fontSize={12} />
              <Tooltip
                contentStyle={{
                  backgroundColor: '#fff',
                  border: '1px solid #e5e7eb',
                  borderRadius: '8px',
                }}
              />
              <Line
                type="monotone"
                dataKey="duration"
                stroke="#6366f1"
                strokeWidth={2}
                dot={{ fill: '#6366f1', r: 4 }}
                activeDot={{ r: 6 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card p-6">
          <h3 className="text-lg font-semibold mb-4">播客收听分布</h3>
          <div className="h-80">
            {pieData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    label={({ name, percent }) =>
                      `${name.length > 10 ? name.slice(0, 10) + '...' : name} ${(percent * 100).toFixed(0)}%`
                    }
                    outerRadius={100}
                    fill="#8884d8"
                    dataKey="value"
                  >
                    {pieData.map((_entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(value: number) => formatDuration(value)}
                  />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-gray-500">
                暂无数据
              </div>
            )}
          </div>
        </div>

        <div className="card p-6">
          <h3 className="text-lg font-semibold mb-4">收听时段分布</h3>
          <div className="h-80">
            {hourData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={hourData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                  <XAxis dataKey="hour" stroke="#6b7280" fontSize={10} />
                  <YAxis stroke="#6b7280" fontSize={12} />
                  <Tooltip />
                  <Bar dataKey="count" fill="#6366f1" radius={[4, 4, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex items-center justify-center h-full text-gray-500">
                暂无数据
              </div>
            )}
          </div>
        </div>
      </div>

      {completionStats && (
        <div className="card p-6">
          <h3 className="text-lg font-semibold mb-4">总收听时长</h3>
          <p className="text-4xl font-bold text-indigo-600">
            {formatDuration(completionStats.total_listened_time || 0)}
          </p>
        </div>
      )}
    </div>
  )
}
