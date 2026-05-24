import { useAuthStore } from '@/store/auth'
import { useNavigate } from 'react-router-dom'
import { GraduationCap, Users, Calendar, Star, Clock, TrendingUp } from 'lucide-react'

export default function Home() {
  const { user } = useAuthStore()
  const navigate = useNavigate()

  const features = [
    {
      icon: Users,
      title: '优质教师',
      description: '数千名经过认证的专业教师',
      color: 'bg-blue-500',
    },
    {
      icon: Calendar,
      title: '灵活预约',
      description: '根据您的时间安排课程',
      color: 'bg-green-500',
    },
    {
      icon: Star,
      title: '实时互动',
      description: '高清视频一对一授课',
      color: 'bg-yellow-500',
    },
    {
      icon: Clock,
      title: '进度追踪',
      description: '实时了解学习进度',
      color: 'bg-purple-500',
    },
  ]

  return (
    <div className="space-y-12">
      <section className="bg-gradient-to-r from-primary-600 to-primary-800 rounded-2xl p-8 md:p-12 text-white">
        <div className="max-w-3xl">
          <h1 className="text-4xl md:text-5xl font-bold mb-4">
            找到最适合您的
            <br />
            一对一在线家教
          </h1>
          <p className="text-xl text-primary-100 mb-8">
            连接全球优质教师，随时随地开启高效学习之旅
          </p>
          <div className="flex flex-wrap gap-4">
            <button
              onClick={() => navigate('/teachers')}
              className="bg-white text-primary-600 px-6 py-3 rounded-lg font-medium hover:bg-primary-50 transition-colors"
            >
              寻找老师
            </button>
            <button
              onClick={() => navigate('/register')}
              className="bg-primary-700 text-white px-6 py-3 rounded-lg font-medium hover:bg-primary-800 transition-colors border border-primary-500"
            >
              成为老师
            </button>
          </div>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-bold text-gray-900 mb-6">为什么选择我们</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {features.map((feature) => (
            <div key={feature.title} className="card hover:shadow-md transition-shadow">
              <div className={`w-12 h-12 ${feature.color} rounded-lg flex items-center justify-center mb-4`}>
                <feature.icon className="h-6 w-6 text-white" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">{feature.title}</h3>
              <p className="text-gray-500">{feature.description}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card bg-blue-50">
          <Users className="h-8 w-8 text-blue-600 mb-3" />
          <div className="text-3xl font-bold text-gray-900">10,000+</div>
          <div className="text-gray-500">认证教师</div>
        </div>
        <div className="card bg-green-50">
          <TrendingUp className="h-8 w-8 text-green-600 mb-3" />
          <div className="text-3xl font-bold text-gray-900">98%</div>
          <div className="text-gray-500">满意度</div>
        </div>
        <div className="card bg-purple-50">
          <GraduationCap className="h-8 w-8 text-purple-600 mb-3" />
          <div className="text-3xl font-bold text-gray-900">50万+</div>
          <div className="text-gray-500">已授课程</div>
        </div>
      </section>

      {user && (
        <section className="card">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            {user.role === 'teacher' ? '教师工作台' : user.role === 'student' ? '学生工作台' : '管理工作台'}
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {user.role === 'student' && (
              <>
                <button
                  onClick={() => navigate('/teachers')}
                  className="p-6 bg-primary-50 rounded-xl hover:bg-primary-100 transition-colors text-left"
                >
                  <Users className="h-8 w-8 text-primary-600 mb-2" />
                  <div className="font-medium text-gray-900">寻找老师</div>
                  <div className="text-sm text-gray-500">浏览并预约合适的老师</div>
                </button>
                <button
                  onClick={() => navigate('/my-bookings')}
                  className="p-6 bg-green-50 rounded-xl hover:bg-green-100 transition-colors text-left"
                >
                  <Calendar className="h-8 w-8 text-green-600 mb-2" />
                  <div className="font-medium text-gray-900">我的课程</div>
                  <div className="text-sm text-gray-500">查看已预约的课程</div>
                </button>
                <button
                  onClick={() => navigate('/learning')}
                  className="p-6 bg-purple-50 rounded-xl hover:bg-purple-100 transition-colors text-left"
                >
                  <TrendingUp className="h-8 w-8 text-purple-600 mb-2" />
                  <div className="font-medium text-gray-900">学习进度</div>
                  <div className="text-sm text-gray-500">追踪您的学习成果</div>
                </button>
              </>
            )}
            {user.role === 'teacher' && (
              <>
                <button
                  onClick={() => navigate('/teacher/dashboard')}
                  className="p-6 bg-blue-50 rounded-xl hover:bg-blue-100 transition-colors text-left"
                >
                  <Calendar className="h-8 w-8 text-blue-600 mb-2" />
                  <div className="font-medium text-gray-900">课程管理</div>
                  <div className="text-sm text-gray-500">查看和管理您的课程</div>
                </button>
                <button
                  onClick={() => navigate('/wallet')}
                  className="p-6 bg-green-50 rounded-xl hover:bg-green-100 transition-colors text-left"
                >
                  <GraduationCap className="h-8 w-8 text-green-600 mb-2" />
                  <div className="font-medium text-gray-900">收入统计</div>
                  <div className="text-sm text-gray-500">查看您的收入详情</div>
                </button>
                <button
                  onClick={() => navigate('/profile')}
                  className="p-6 bg-purple-50 rounded-xl hover:bg-purple-100 transition-colors text-left"
                >
                  <Users className="h-8 w-8 text-purple-600 mb-2" />
                  <div className="font-medium text-gray-900">个人资料</div>
                  <div className="text-sm text-gray-500">完善您的教师信息</div>
                </button>
              </>
            )}
          </div>
        </section>
      )}
    </div>
  )
}
