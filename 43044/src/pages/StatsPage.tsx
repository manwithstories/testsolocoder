import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts'
import { useAppStore } from '@/store'
import { getOverallStats, getTagStats, getWeeklyActivity, getPriorityStats, getMonthlyStats } from '@/utils/stats'
import { formatDate } from '@/utils/date'
import { Badge } from '@/components/ui/Badge'



export function StatsPage() {
  const { goals, tasks, tags } = useAppStore()
  
  const overallStats = getOverallStats(goals, tasks)
  const tagStats = getTagStats(goals, tasks, tags)
  const weeklyActivity = getWeeklyActivity(tasks, 7)
  const priorityStats = getPriorityStats(tasks)
  const monthlyStats = getMonthlyStats(goals, tasks, 6)
  
  const weeklyChartData = weeklyActivity.map((item) => ({
    date: formatDate(item.date, 'MM/dd'),
    完成: item.completedTasks,
    创建: item.createdTasks,
  }))
  
  const tagChartData = tagStats.map((item) => ({
    name: item.tag.name,
    完成率: item.completionRate,
    完成: item.completedTasks,
    总数: item.totalTasks,
  }))
  
  const statusData = [
    { name: '已完成', value: overallStats.completedTasks, color: '#10B981' },
    { name: '进行中', value: overallStats.inProgressTasks, color: '#3B82F6' },
    { name: '待开始', value: overallStats.pendingTasks, color: '#9CA3AF' },
    { name: '已搁置', value: overallStats.onHoldTasks, color: '#F59E0B' },
  ]
  
  const pieData = statusData.filter((item) => item.value > 0)
  
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800">统计分析</h1>
        <p className="text-gray-500 mt-1">查看目标完成情况和活跃度统计</p>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <p className="text-sm text-gray-500">总目标数</p>
          <p className="text-3xl font-bold text-gray-800 mt-2">{overallStats.totalGoals}</p>
          <div className="flex items-center gap-2 mt-2">
            <Badge text={`${overallStats.completedGoals} 已完成`} color="bg-green-100 text-green-700" size="sm" />
            <Badge text={`${overallStats.goalsInProgress} 进行中`} color="bg-blue-100 text-blue-700" size="sm" />
          </div>
        </div>
        
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <p className="text-sm text-gray-500">总任务数</p>
          <p className="text-3xl font-bold text-gray-800 mt-2">{overallStats.totalTasks}</p>
          <div className="flex items-center gap-2 mt-2">
            <Badge text={`${overallStats.completedTasks} 已完成`} color="bg-green-100 text-green-700" size="sm" />
          </div>
        </div>
        
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <p className="text-sm text-gray-500">整体完成率</p>
          <p className="text-3xl font-bold text-blue-600 mt-2">{overallStats.overallCompletionRate}%</p>
          <div className="w-full bg-gray-200 rounded-full h-2 mt-3">
            <div
              className="bg-blue-500 h-2 rounded-full transition-all"
              style={{ width: `${overallStats.overallCompletionRate}%` }}
            />
          </div>
        </div>
        
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <p className="text-sm text-gray-500">逾期任务</p>
          <p className="text-3xl font-bold text-red-600 mt-2">{overallStats.overdueTasks}</p>
          {overallStats.overdueTasks > 0 && (
            <Badge text="需要关注" color="bg-red-100 text-red-700" size="sm" className="mt-2" />
          )}
        </div>
      </div>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">最近一周活跃度</h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={weeklyChartData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
                <XAxis dataKey="date" tick={{ fontSize: 12 }} stroke="#9CA3AF" />
                <YAxis tick={{ fontSize: 12 }} stroke="#9CA3AF" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'white',
                    border: '1px solid #E5E7EB',
                    borderRadius: '8px',
                  }}
                />
                <Legend />
                <Bar dataKey="创建" fill="#3B82F6" radius={[4, 4, 0, 0]} />
                <Bar dataKey="完成" fill="#10B981" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
        
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">任务状态分布</h3>
          <div className="h-64 flex items-center justify-center">
            {pieData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={90}
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {pieData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.color} />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{
                      backgroundColor: 'white',
                      border: '1px solid #E5E7EB',
                      borderRadius: '8px',
                    }}
                  />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            ) : (
              <p className="text-gray-500">暂无数据</p>
            )}
          </div>
        </div>
      </div>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">按标签分组完成率</h3>
          {tagChartData.length > 0 ? (
            <div className="h-64">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={tagChartData} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
                  <XAxis type="number" domain={[0, 100]} tick={{ fontSize: 12 }} stroke="#9CA3AF" />
                  <YAxis dataKey="name" type="category" width={80} tick={{ fontSize: 12 }} stroke="#9CA3AF" />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: 'white',
                      border: '1px solid #E5E7EB',
                      borderRadius: '8px',
                    }}
                    formatter={(value) => [`${value}%`, '完成率']}
                  />
                  <Bar dataKey="完成率" fill="#8B5CF6" radius={[0, 4, 4, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <div className="h-64 flex items-center justify-center">
              <p className="text-gray-500">暂无标签数据</p>
            </div>
          )}
        </div>
        
        <div className="bg-white rounded-lg border border-gray-200 p-5">
          <h3 className="text-lg font-semibold text-gray-800 mb-4">按月统计</h3>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={monthlyStats}>
                <CartesianGrid strokeDasharray="3 3" stroke="#E5E7EB" />
                <XAxis dataKey="label" tick={{ fontSize: 11 }} stroke="#9CA3AF" />
                <YAxis tick={{ fontSize: 12 }} stroke="#9CA3AF" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'white',
                    border: '1px solid #E5E7EB',
                    borderRadius: '8px',
                  }}
                />
                <Legend />
                <Bar dataKey="创建" fill="#6366F1" radius={[4, 4, 0, 0]} />
                <Bar dataKey="完成" fill="#22C55E" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
      
      <div className="bg-white rounded-lg border border-gray-200 p-5">
        <h3 className="text-lg font-semibold text-gray-800 mb-4">优先级完成情况</h3>
        <div className="space-y-4">
          {priorityStats.map((stat) => (
            <div key={stat.priority}>
              <div className="flex items-center justify-between mb-1">
                <span className="text-sm font-medium text-gray-700 capitalize">
                  {stat.priority === 'low' ? '低优先级' :
                   stat.priority === 'medium' ? '中优先级' :
                   stat.priority === 'high' ? '高优先级' : '紧急'}
                </span>
                <span className="text-sm text-gray-500">
                  {stat.completed} / {stat.total} ({stat.rate}%)
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className={`h-2 rounded-full transition-all ${
                    stat.priority === 'urgent' ? 'bg-red-500' :
                    stat.priority === 'high' ? 'bg-orange-500' :
                    stat.priority === 'medium' ? 'bg-blue-500' : 'bg-gray-500'
                  }`}
                  style={{ width: `${stat.rate}%` }}
                />
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
