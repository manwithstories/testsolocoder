import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { reminderAPI, planAPI, activityAPI } from '@/services/api'
import type { Reminder, TravelPlan, Activity } from '@/types'
import dayjs from 'dayjs'

const channelLabels: Record<string, string> = {
  email: '📧 邮件',
  app: '🔔 应用内通知',
}

export default function Reminders() {
  const { id } = useParams<{ id: string }>()
  const [reminders, setReminders] = useState<Reminder[]>([])
  const [plans, setPlans] = useState<TravelPlan[]>([])
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [editingReminder, setEditingReminder] = useState<Reminder | null>(null)
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    reminder_time: '',
    channel: 'app' as 'email' | 'app',
    plan_id: '',
    activity_id: '',
  })

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      setLoading(true)
      const [remindersData, plansData] = await Promise.all([
        reminderAPI.getReminders(),
        planAPI.getPlans({ page_size: 100 }),
      ])
      setReminders(remindersData)
      setPlans(plansData.data)
      if (plansData.data.length > 0 && !formData.plan_id) {
        setFormData((prev) => ({ ...prev, plan_id: plansData.data[0].id }))
        loadActivities(plansData.data[0].id)
      }
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const loadActivities = async (planId: string) => {
    try {
      const activitiesData = await activityAPI.getActivities(planId)
      setActivities(activitiesData as Activity[])
    } catch (error) {
      console.error('Failed to load activities:', error)
    }
  }

  const handlePlanChange = (planId: string) => {
    setFormData((prev) => ({ ...prev, plan_id: planId, activity_id: '' }))
    loadActivities(planId)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingReminder) {
        await reminderAPI.updateReminder(editingReminder.id, formData)
      } else {
        await reminderAPI.createReminder(formData.plan_id, formData)
      }
      setShowModal(false)
      resetForm()
      loadData()
    } catch (error: any) {
      alert(error.message || '保存失败')
    }
  }

  const handleDelete = async (reminderId: string) => {
    if (!confirm('确定要删除这个提醒吗？')) return
    try {
      await reminderAPI.deleteReminder(reminderId)
      loadData()
    } catch (error: any) {
      alert(error.message || '删除失败')
    }
  }

  const handleEdit = (reminder: Reminder) => {
    setEditingReminder(reminder)
    setFormData({
      title: reminder.title,
      description: reminder.description || '',
      reminder_time: dayjs(reminder.reminder_time).format('YYYY-MM-DDTHH:mm'),
      channel: reminder.channel,
      plan_id: reminder.plan_id,
      activity_id: reminder.activity_id || '',
    })
    if (reminder.plan_id) {
      loadActivities(reminder.plan_id)
    }
    setShowModal(true)
  }

  const resetForm = () => {
    setEditingReminder(null)
    setFormData({
      title: '',
      description: '',
      reminder_time: '',
      channel: 'app',
      plan_id: plans[0]?.id || '',
      activity_id: '',
    })
  }

  const getUpcomingReminders = () => {
    const now = dayjs()
    return reminders
      .filter((r) => dayjs(r.reminder_time).isAfter(now) && !r.is_sent)
      .sort((a, b) => dayjs(a.reminder_time).valueOf() - dayjs(b.reminder_time).valueOf())
  }

  const getPastReminders = () => {
    const now = dayjs()
    return reminders
      .filter((r) => dayjs(r.reminder_time).isBefore(now) || r.is_sent)
      .sort((a, b) => dayjs(b.reminder_time).valueOf() - dayjs(a.reminder_time).valueOf())
  }

  const formatReminderTime = (time: string) => {
    const reminderTime = dayjs(time)
    const now = dayjs()
    const diffDays = reminderTime.diff(now, 'day')

    if (diffDays === 0) {
      return `今天 ${reminderTime.format('HH:mm')}`
    } else if (diffDays === 1) {
      return `明天 ${reminderTime.format('HH:mm')}`
    } else if (diffDays < 7) {
      return `${diffDays}天后 ${reminderTime.format('HH:mm')}`
    } else {
      return reminderTime.format('YYYY/MM/DD HH:mm')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  const upcoming = getUpcomingReminders()
  const past = getPastReminders()

  return (
    <div className="p-6">
      <div className="mb-6">
        <div className="flex items-center justify-between">
          <div>
            {id && (
              <Link to={`/plans/${id}`} className="text-primary-600 hover:underline text-sm mb-2 inline-block">
                ← 返回计划详情
              </Link>
            )}
            <h1 className="text-2xl font-bold text-gray-900">提醒管理</h1>
          </div>
          <button
            onClick={() => setShowModal(true)}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
          >
            + 创建提醒
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div>
          <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <span className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
            即将到来 ({upcoming.length})
          </h2>
          {upcoming.length === 0 ? (
            <div className="bg-white rounded-xl shadow-sm p-8 text-center text-gray-500">
              <p className="text-4xl mb-2">⏰</p>
              <p>暂无即将到来的提醒</p>
            </div>
          ) : (
            <div className="space-y-3">
              {upcoming.map((reminder) => (
                <div
                  key={reminder.id}
                  className="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="font-medium text-gray-900">{reminder.title}</h3>
                        <span className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded">
                          {channelLabels[reminder.channel]}
                        </span>
                      </div>
                      <p className="text-sm text-gray-500 mb-2">
                        🕐 {formatReminderTime(reminder.reminder_time)}
                      </p>
                      {reminder.description && (
                        <p className="text-sm text-gray-600">{reminder.description}</p>
                      )}
                      {reminder.plan && (
                        <p className="text-xs text-gray-400 mt-2">
                          📋 关联计划: {reminder.plan.title}
                        </p>
                      )}
                      {reminder.activity && (
                        <p className="text-xs text-gray-400">
                          📅 关联活动: {reminder.activity.title}
                        </p>
                      )}
                    </div>
                    <div className="flex gap-1">
                      <button
                        onClick={() => handleEdit(reminder)}
                        className="p-2 text-gray-500 hover:text-primary-600 hover:bg-gray-100 rounded"
                      >
                        ✏️
                      </button>
                      <button
                        onClick={() => handleDelete(reminder.id)}
                        className="p-2 text-gray-500 hover:text-red-600 hover:bg-gray-100 rounded"
                      >
                        🗑️
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div>
          <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <span className="w-2 h-2 bg-gray-400 rounded-full"></span>
            已发送 / 已过期 ({past.length})
          </h2>
          {past.length === 0 ? (
            <div className="bg-white rounded-xl shadow-sm p-8 text-center text-gray-500">
              <p className="text-4xl mb-2">📭</p>
              <p>暂无历史提醒</p>
            </div>
          ) : (
            <div className="space-y-3">
              {past.map((reminder) => (
                <div
                  key={reminder.id}
                  className="bg-white rounded-xl shadow-sm p-4 opacity-70"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h3 className="font-medium text-gray-900">{reminder.title}</h3>
                        {reminder.is_sent && (
                          <span className="text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded">
                            ✓ 已发送
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-gray-500 mb-2">
                        🕐 {dayjs(reminder.reminder_time).format('YYYY/MM/DD HH:mm')}
                      </p>
                      {reminder.description && (
                        <p className="text-sm text-gray-600">{reminder.description}</p>
                      )}
                    </div>
                    <button
                      onClick={() => handleDelete(reminder.id)}
                      className="p-2 text-gray-500 hover:text-red-600 hover:bg-gray-100 rounded"
                    >
                      🗑️
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <h2 className="text-xl font-bold mb-6">
                {editingReminder ? '编辑提醒' : '创建提醒'}
              </h2>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">提醒标题 *</label>
                  <input
                    type="text"
                    required
                    value={formData.title}
                    onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    placeholder="输入提醒标题"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">提醒时间 *</label>
                  <input
                    type="datetime-local"
                    required
                    value={formData.reminder_time}
                    onChange={(e) => setFormData({ ...formData, reminder_time: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">通知方式</label>
                  <select
                    value={formData.channel}
                    onChange={(e) => setFormData({ ...formData, channel: e.target.value as 'email' | 'app' })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  >
                    <option value="app">🔔 应用内通知</option>
                    <option value="email">📧 邮件通知</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">关联计划</label>
                  <select
                    value={formData.plan_id}
                    onChange={(e) => handlePlanChange(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  >
                    <option value="">无</option>
                    {plans.map((plan) => (
                      <option key={plan.id} value={plan.id}>
                        {plan.title}
                      </option>
                    ))}
                  </select>
                </div>
                {formData.plan_id && activities.length > 0 && (
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">关联活动</label>
                    <select
                      value={formData.activity_id}
                      onChange={(e) => setFormData({ ...formData, activity_id: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    >
                      <option value="">无</option>
                      {activities.map((activity) => (
                        <option key={activity.id} value={activity.id}>
                          {activity.title} ({dayjs(activity.date).format('MM/DD')})
                        </option>
                      ))}
                    </select>
                  </div>
                )}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    rows={3}
                    placeholder="提醒描述（可选）"
                  />
                </div>
                <div className="flex gap-3 pt-4">
                  <button
                    type="button"
                    onClick={() => {
                      setShowModal(false)
                      resetForm()
                    }}
                    className="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    className="flex-1 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    {editingReminder ? '保存修改' : '创建提醒'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
