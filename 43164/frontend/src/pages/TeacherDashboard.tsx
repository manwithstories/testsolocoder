import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { teacherApi } from '@/services/api'
import { TeacherProfile } from '@/types'
import { Star, Clock, DollarSign, Calendar, BarChart3, TrendingUp } from 'lucide-react'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

export default function TeacherDashboard() {
  const navigate = useNavigate()
  const [profile, setProfile] = useState<TeacherProfile | null>(null)
  const [earnings, setEarnings] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [profileRes, earningsRes] = await Promise.all([
        teacherApi.getProfile(),
        teacherApi.getEarnings(),
      ])
      setProfile(profileRes.data)
      setEarnings(earningsRes.data)
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

  const earningsData = [
    { month: '1月', earnings: 1200 },
    { month: '2月', earnings: 1800 },
    { month: '3月', earnings: 1500 },
    { month: '4月', earnings: 2200 },
    { month: '5月', earnings: 2800 },
    { month: '6月', earnings: 2400 },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">教师工作台</h1>
        <p className="text-gray-500">管理您的教学业务和收入</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
              <DollarSign className="h-6 w-6 text-blue-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">账户余额</div>
              <div className="text-2xl font-bold text-gray-900">
                ${earnings?.balance?.toFixed(2) || '0.00'}
              </div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
              <TrendingUp className="h-6 w-6 text-green-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">本月收入</div>
              <div className="text-2xl font-bold text-gray-900">
                ${earnings?.monthEarnings?.toFixed(2) || '0.00'}
              </div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center">
              <Star className="h-6 w-6 text-yellow-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">评分</div>
              <div className="text-2xl font-bold text-gray-900">
                {profile?.rating?.toFixed(1) || '5.0'}
              </div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
              <Clock className="h-6 w-6 text-purple-600" />
            </div>
            <div>
              <div className="text-sm text-gray-500">总课时</div>
              <div className="text-2xl font-bold text-gray-900">
                {profile?.totalHours || 0}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">收入趋势</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={earningsData}>
                <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                <XAxis dataKey="month" stroke="#9ca3af" fontSize={12} />
                <YAxis stroke="#9ca3af" fontSize={12} />
                <Tooltip />
                <Line
                  type="monotone"
                  dataKey="earnings"
                  stroke="#3b82f6"
                  strokeWidth={2}
                  dot={{ fill: '#3b82f6' }}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">快捷操作</h2>
          <div className="grid grid-cols-2 gap-4">
            <button
              onClick={() => navigate('/my-bookings')}
              className="p-4 bg-blue-50 rounded-xl hover:bg-blue-100 transition-colors text-left"
            >
              <Calendar className="h-6 w-6 text-blue-600 mb-2" />
              <div className="font-medium text-gray-900">课程管理</div>
              <div className="text-sm text-gray-500">查看和管理课程</div>
            </button>
            <button
              onClick={() => navigate('/profile')}
              className="p-4 bg-purple-50 rounded-xl hover:bg-purple-100 transition-colors text-left"
            >
              <BarChart3 className="h-6 w-6 text-purple-600 mb-2" />
              <div className="font-medium text-gray-900">个人资料</div>
              <div className="text-sm text-gray-500">完善教师信息</div>
            </button>
            <button
              onClick={() => navigate('/wallet')}
              className="p-4 bg-green-50 rounded-xl hover:bg-green-100 transition-colors text-left"
            >
              <DollarSign className="h-6 w-6 text-green-600 mb-2" />
              <div className="font-medium text-gray-900">收入管理</div>
              <div className="text-sm text-gray-500">查看收入和提现</div>
            </button>
            <button
              onClick={() => navigate('/messages')}
              className="p-4 bg-yellow-50 rounded-xl hover:bg-yellow-100 transition-colors text-left"
            >
              <Star className="h-6 w-6 text-yellow-600 mb-2" />
              <div className="font-medium text-gray-900">消息中心</div>
              <div className="text-sm text-gray-500">与学生沟通</div>
            </button>
          </div>
        </div>
      </div>

      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">个人信息</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">简介</label>
            <p className="text-gray-600">{profile?.bio || '暂无简介'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">教育背景</label>
            <p className="text-gray-600">{profile?.education || '暂无信息'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">教学经验</label>
            <p className="text-gray-600">{profile?.experience || '暂无信息'}</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">资质认证</label>
            <p className="text-gray-600">{profile?.certifications || '暂无信息'}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
