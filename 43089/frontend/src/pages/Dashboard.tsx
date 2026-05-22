import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { planAPI, activityAPI, reminderAPI } from '@/services/api'
import type { TravelPlan, Reminder } from '@/types'
import dayjs from 'dayjs'

export default function Dashboard() {
  const [plans, setPlans] = useState<TravelPlan[]>([])
  const [reminders, setReminders] = useState<Reminder[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [plansRes, remindersRes] = await Promise.all([
        planAPI.getPlans({ page_size: 5 }),
        reminderAPI.getReminders(),
      ])
      setPlans(plansRes.data)
      setReminders(remindersRes.slice(0, 5))
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
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

  const stats = [
    { label: '旅行计划', value: plans.length, icon: '✈️', color: 'bg-blue-500' },
    { label: '待办提醒', value: reminders.filter((r) => !r.is_sent).length, icon: '🔔', color: 'bg-yellow-500' },
    { label: '进行中', value: plans.filter((p) => p.status === 'active').length, icon: '🚀', color: 'bg-green-500' },
    { label: '已完成', value: plans.filter((p) => p.status === 'completed').length, icon: '✅', color: 'bg-purple-500' },
  ]

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">仪表盘</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat) => (
          <div key={stat.label} className="bg-white rounded-xl shadow-sm p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-500 text-sm">{stat.label}</p>
                <p className="text-3xl font-bold text-gray-900 mt-1">{stat.value}</p>
              </div>
              <div className={`${stat.color} w-12 h-12 rounded-lg flex items-center justify-center text-2xl`}>
                {stat.icon}
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-xl shadow-sm p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900">最近的旅行计划</h2>
            <Link
              to="/plans"
              className="text-primary-600 hover:text-primary-700 text-sm font-medium"
            >
              查看全部
            </Link>
          </div>
          {plans.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <p className="text-4xl mb-2">📝</p>
              <p>还没有旅行计划</p>
              <Link
                to="/plans"
                className="mt-4 inline-block bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700 transition-colors"
              >
                创建第一个计划
              </Link>
            </div>
          ) : (
            <div className="space-y-4">
              {plans.map((plan) => (
                <Link
                  key={plan.id}
                  to={`/plans/${plan.id}`}
                  className="block p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <h3 className="font-medium text-gray-900">{plan.title}</h3>
                      <p className="text-sm text-gray-500 mt-1">
                        📍 {plan.destination}
                      </p>
                      <p className="text-sm text-gray-500">
                        📅 {dayjs(plan.start_date).format('YYYY/MM/DD')} - {dayjs(plan.end_date).format('YYYY/MM/DD')}
                      </p>
                    </div>
                    <span
                      className={`px-2 py-1 rounded-full text-xs font-medium ${
                        plan.status === 'active'
                          ? 'bg-green-100 text-green-700'
                          : plan.status === 'completed'
                          ? 'bg-gray-100 text-gray-700'
                          : 'bg-yellow-100 text-yellow-700'
                      }`}
                    >
                      {plan.status === 'draft'
                        ? '草稿'
                        : plan.status === 'active'
                        ? '进行中'
                        : '已完成'}
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>

        <div className="bg-white rounded-xl shadow-sm p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold text-gray-900">即将到来的提醒</h2>
            <Link
              to="/reminders"
              className="text-primary-600 hover:text-primary-700 text-sm font-medium"
            >
              查看全部
            </Link>
          </div>
          {reminders.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <p className="text-4xl mb-2">🔔</p>
              <p>暂无提醒</p>
            </div>
          ) : (
            <div className="space-y-3">
              {reminders.map((reminder) => (
                <div
                  key={reminder.id}
                  className="flex items-center p-3 bg-gray-50 rounded-lg"
                >
                  <div
                    className={`w-2 h-2 rounded-full mr-3 ${
                      reminder.is_sent ? 'bg-gray-400' : 'bg-yellow-500'
                    }`}
                  />
                  <div className="flex-1">
                    <p className="font-medium text-gray-900 text-sm">{reminder.title}</p>
                    <p className="text-xs text-gray-500">
                      {dayjs(reminder.reminder_time).format('YYYY/MM/DD HH:mm')}
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
