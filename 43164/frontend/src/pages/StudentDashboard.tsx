import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { studentApi, bookingApi } from '@/services/api'
import { LearningGoal, Booking, Milestone } from '@/types'
import { Target, BookOpen, Award, Calendar, TrendingUp, Users } from 'lucide-react'

export default function StudentDashboard() {
  const navigate = useNavigate()
  const [goals, setGoals] = useState<LearningGoal[]>([])
  const [milestones, setMilestones] = useState<Milestone[]>([])
  const [recentBookings, setRecentBookings] = useState<Booking[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [profileRes, milestonesRes, bookingsRes] = await Promise.all([
        studentApi.getProfile(),
        studentApi.getMilestones(),
        bookingApi.getAll(),
      ])
      setGoals(profileRes.data?.learningGoals || [])
      setMilestones(milestonesRes.data || [])
      setRecentBookings(bookingsRes.data?.slice(0, 5) || [])
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">学习工作台</h1>
        <p className="text-gray-500">追踪您的学习进度和成就</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
              <BookOpen className="h-6 w-6 text-blue-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">学习目标</div>
              <div className="text-2xl font-bold text-gray-900">{goals.length}</div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
              <Award className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">已达成里程碑</div>
              <div className="text-2xl font-bold text-gray-900">
                {milestones.filter(m => m.isAchieved).length}
              </div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center">
              <Calendar className="h-6 w-6 text-yellow-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">已完成课程</div>
              <div className="text-2xl font-bold text-gray-900">
                {recentBookings.filter(b => b.status === 'completed').length}
              </div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
              <TrendingUp className="h-6 w-6 text-purple-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">学习时长</div>
              <div className="text-2xl font-bold text-gray-900">
                {recentBookings.reduce((acc, b) => acc + b.duration, 0) / 60}小时
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900">学习目标</h2>
            <button
              onClick={() => navigate('/learning')}
              className="text-primary-600 hover:text-primary-700 text-sm font-medium"
            >
              查看全部
            </button>
          </div>
          <div className="space-y-3">
            {goals.slice(0, 5).map((goal) => (
              <div key={goal.id} className="p-3 bg-gray-50 rounded-lg">
                <div className="flex items-center justify-between">
                  <span className="font-medium text-gray-900">{goal.title}</span>
                  {goal.isAchieved && (
                    <span className="badge bg-green-100 text-green-800">已达成</span>
                  )}
                </div>
                <div className="mt-2 flex items-center gap-2">
                  <div className="flex-1 bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-primary-600 h-2 rounded-full"
                      style={{ width: `${Math.min((goal.currentScore / (goal.targetScore || 100)) * 100, 100)}%` }}
                    />
                  </div>
                  <span className="text-sm text-gray-500">
                    {goal.currentScore}/{goal.targetScore}
                  </span>
                </div>
              </div>
            ))}
            {goals.length === 0 && (
              <p className="text-gray-500 text-center py-4">暂无学习目标</p>
            )}
          </div>
        </div>

        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900">最近课程</h2>
            <button
              onClick={() => navigate('/my-bookings')}
              className="text-primary-600 hover:text-primary-700 text-sm font-medium"
            >
              查看全部
            </button>
          </div>
          <div className="space-y-3">
            {recentBookings.map((booking) => (
              <div key={booking.id} className="p-3 bg-gray-50 rounded-lg flex items-center justify-between">
                <div>
                  <div className="font-medium text-gray-900">{booking.subject?.name}</div>
                  <div className="text-sm text-gray-500">
                    {new Date(booking.startTime).toLocaleDateString('zh-CN', {
                      month: 'short',
                      day: 'numeric',
                    })}
                  </div>
                </div>
                <span className={`badge ${
                  booking.status === 'completed' ? 'bg-green-100 text-green-800' :
                  booking.status === 'confirmed' ? 'bg-blue-100 text-blue-800' :
                  booking.status === 'pending' ? 'bg-yellow-100 text-yellow-800' :
                  'bg-gray-100 text-gray-800'
                }`}>
                  {booking.status === 'completed' ? '已完成' :
                   booking.status === 'confirmed' ? '已确认' :
                   booking.status === 'pending' ? '待确认' : booking.status}
                </span>
              </div>
            ))}
            {recentBookings.length === 0 && (
              <p className="text-gray-500 text-center py-4">暂无课程记录</p>
            )}
          </div>
        </div>
      </div>

      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">快捷操作</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <button
            onClick={() => navigate('/teachers')}
            className="p-4 bg-blue-50 rounded-xl hover:bg-blue-100 transition-colors text-left"
          >
            <Users className="h-6 w-6 text-blue-600 mb-2" />
            <div className="font-medium text-gray-900">寻找老师</div>
            <div className="text-sm text-gray-500">浏览老师列表</div>
          </button>
          <button
            onClick={() => navigate('/learning')}
            className="p-4 bg-green-50 rounded-xl hover:bg-green-100 transition-colors text-left"
          >
            <Target className="h-6 w-6 text-green-600 mb-2" />
            <div className="font-medium text-gray-900">学习进度</div>
            <div className="text-sm text-gray-500">查看学习详情</div>
          </button>
          <button
            onClick={() => navigate('/my-bookings')}
            className="p-4 bg-yellow-50 rounded-xl hover:bg-yellow-100 transition-colors text-left"
          >
            <Calendar className="h-6 w-6 text-yellow-600 mb-2" />
            <div className="font-medium text-gray-900">我的课程</div>
            <div className="text-sm text-gray-500">管理课程预约</div>
          </button>
          <button
            onClick={() => navigate('/messages')}
            className="p-4 bg-purple-50 rounded-xl hover:bg-purple-100 transition-colors text-left"
          >
            <Target className="h-6 w-6 text-purple-600 mb-2" />
            <div className="font-medium text-gray-900">消息中心</div>
            <div className="text-sm text-gray-500">与老师沟通</div>
          </button>
        </div>
      </div>
    </div>
  )
}
