import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { activityAPI, planAPI } from '@/services/api'
import type { Activity, PlanDetail } from '@/types'
import dayjs from 'dayjs'

const activityTypes = [
  { value: 'sightseeing', label: '景点', icon: '🏛️', color: 'bg-blue-100 text-blue-700' },
  { value: 'transport', label: '交通', icon: '🚗', color: 'bg-yellow-100 text-yellow-700' },
  { value: 'accommodation', label: '住宿', icon: '🏨', color: 'bg-purple-100 text-purple-700' },
  { value: 'food', label: '餐饮', icon: '🍽️', color: 'bg-orange-100 text-orange-700' },
  { value: 'other', label: '其他', icon: '📌', color: 'bg-gray-100 text-gray-700' },
]

export default function Activities() {
  const { id } = useParams<{ id: string }>()
  const [plan, setPlan] = useState<PlanDetail | null>(null)
  const [activities, setActivities] = useState<Activity[]>([])
  const [loading, setLoading] = useState(true)
  const [showModal, setShowModal] = useState(false)
  const [selectedDate, setSelectedDate] = useState<string>('')
  const [editingActivity, setEditingActivity] = useState<Activity | null>(null)
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    type: 'sightseeing',
    date: '',
    start_time: '',
    end_time: '',
    location: '',
    cost: 0,
    notes: '',
  })

  useEffect(() => {
    if (id) {
      loadData()
    }
  }, [id])

  const loadData = async () => {
    try {
      setLoading(true)
      const [planData, activitiesData] = await Promise.all([
        planAPI.getPlan(id!),
        activityAPI.getActivities(id!),
      ])
      setPlan(planData)
      setActivities(activitiesData as Activity[])
      if (planData.start_date && !selectedDate) {
        setSelectedDate(planData.start_date)
        setFormData((prev) => ({ ...prev, date: planData.start_date }))
      }
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const getDateRange = () => {
    if (!plan) return []
    const start = dayjs(plan.start_date)
    const end = dayjs(plan.end_date)
    const dates: string[] = []
    let current = start
    while (current.isBefore(end) || current.isSame(end, 'day')) {
      dates.push(current.format('YYYY-MM-DD'))
      current = current.add(1, 'day')
    }
    return dates
  }

  const getActivitiesByDate = (date: string) => {
    return activities
      .filter((a) => a.date === date)
      .sort((a, b) => (a.start_time || '').localeCompare(b.start_time || ''))
  }

  const getActivityTypeInfo = (type: string) => {
    return activityTypes.find((t) => t.value === type) || activityTypes[4]
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingActivity) {
        await activityAPI.updateActivity(id!, editingActivity.id, formData)
      } else {
        await activityAPI.createActivity(id!, formData)
      }
      setShowModal(false)
      resetForm()
      loadData()
    } catch (error: any) {
      alert(error.message || '保存失败')
    }
  }

  const handleDelete = async (activityId: string) => {
    if (!confirm('确定要删除这个活动吗？')) return
    try {
      await activityAPI.deleteActivity(id!, activityId)
      loadData()
    } catch (error: any) {
      alert(error.message || '删除失败')
    }
  }

  const handleEdit = (activity: Activity) => {
    setEditingActivity(activity)
    setFormData({
      title: activity.title,
      description: activity.description || '',
      type: activity.type,
      date: activity.date,
      start_time: activity.start_time || '',
      end_time: activity.end_time || '',
      location: activity.location || '',
      cost: activity.cost || 0,
      notes: activity.notes || '',
    })
    setShowModal(true)
  }

  const resetForm = () => {
    setEditingActivity(null)
    setFormData({
      title: '',
      description: '',
      type: 'sightseeing',
      date: selectedDate,
      start_time: '',
      end_time: '',
      location: '',
      cost: 0,
      notes: '',
    })
  }

  const openAddModal = (date: string) => {
    setSelectedDate(date)
    setFormData((prev) => ({ ...prev, date }))
    setShowModal(true)
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  const dateRange = getDateRange()

  return (
    <div className="p-6">
      <div className="mb-6">
        <Link to={`/plans/${id}`} className="text-primary-600 hover:underline text-sm mb-2 inline-block">
          ← 返回计划详情
        </Link>
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-bold text-gray-900">行程安排</h1>
          <button
            onClick={() => openAddModal(selectedDate)}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
          >
            + 添加活动
          </button>
        </div>
      </div>

      <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
        {dateRange.map((date, index) => (
          <button
            key={date}
            onClick={() => setSelectedDate(date)}
            className={`flex-shrink-0 px-4 py-3 rounded-lg text-sm font-medium transition-colors ${
              selectedDate === date
                ? 'bg-primary-600 text-white'
                : 'bg-white border border-gray-200 hover:bg-gray-50'
            }`}
          >
            <div>第 {index + 1} 天</div>
            <div className="text-xs opacity-80">{dayjs(date).format('MM/DD')}</div>
          </button>
        ))}
      </div>

      <div className="bg-white rounded-xl shadow-sm">
        <div className="p-6">
          <h2 className="text-lg font-semibold mb-4">
            {dayjs(selectedDate).format('YYYY年MM月DD日')} 的行程
          </h2>
          {getActivitiesByDate(selectedDate).length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              <p className="text-4xl mb-2">📅</p>
              <p>这一天还没有安排活动</p>
              <button
                onClick={() => openAddModal(selectedDate)}
                className="mt-4 text-primary-600 hover:underline"
              >
                + 添加第一个活动
              </button>
            </div>
          ) : (
            <div className="space-y-4">
              {getActivitiesByDate(selectedDate).map((activity) => {
                const typeInfo = getActivityTypeInfo(activity.type)
                return (
                  <div
                    key={activity.id}
                    className="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <span className={`px-2 py-1 rounded text-xs font-medium ${typeInfo.color}`}>
                            {typeInfo.icon} {typeInfo.label}
                          </span>
                          {(activity.start_time || activity.end_time) && (
                            <span className="text-sm text-gray-500">
                              {activity.start_time} {activity.end_time ? `- ${activity.end_time}` : ''}
                            </span>
                          )}
                        </div>
                        <h3 className="font-medium text-gray-900">{activity.title}</h3>
                        {activity.description && (
                          <p className="text-sm text-gray-600 mt-1">{activity.description}</p>
                        )}
                        <div className="flex items-center gap-4 mt-3 text-sm text-gray-500">
                          {activity.location && <span>📍 {activity.location}</span>}
                          {activity.cost > 0 && (
                            <span>💰 {activity.cost.toLocaleString()} {activity.currency || plan?.currency}</span>
                          )}
                        </div>
                        {activity.notes && (
                          <p className="text-sm text-gray-500 mt-2 italic">📝 {activity.notes}</p>
                        )}
                      </div>
                      <div className="flex gap-2">
                        <button
                          onClick={() => handleEdit(activity)}
                          className="p-2 text-gray-500 hover:text-primary-600 hover:bg-gray-100 rounded"
                        >
                          ✏️
                        </button>
                        <button
                          onClick={() => handleDelete(activity.id)}
                          className="p-2 text-gray-500 hover:text-red-600 hover:bg-gray-100 rounded"
                        >
                          🗑️
                        </button>
                      </div>
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      </div>

      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
            <div className="p-6">
              <h2 className="text-xl font-bold mb-6">
                {editingActivity ? '编辑活动' : '添加活动'}
              </h2>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">活动名称 *</label>
                  <input
                    type="text"
                    required
                    value={formData.title}
                    onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    placeholder="输入活动名称"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">活动类型 *</label>
                  <select
                    required
                    value={formData.type}
                    onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                  >
                    {activityTypes.map((type) => (
                      <option key={type.value} value={type.value}>
                        {type.icon} {type.label}
                      </option>
                    ))}
                  </select>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">日期 *</label>
                    <input
                      type="date"
                      required
                      value={formData.date}
                      onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">费用</label>
                    <input
                      type="number"
                      value={formData.cost}
                      onChange={(e) => setFormData({ ...formData, cost: Number(e.target.value) })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                      placeholder="0"
                    />
                  </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">开始时间</label>
                    <input
                      type="time"
                      value={formData.start_time}
                      onChange={(e) => setFormData({ ...formData, start_time: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">结束时间</label>
                    <input
                      type="time"
                      value={formData.end_time}
                      onChange={(e) => setFormData({ ...formData, end_time: e.target.value })}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">地点</label>
                  <input
                    type="text"
                    value={formData.location}
                    onChange={(e) => setFormData({ ...formData, location: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    placeholder="输入活动地点"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">描述</label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    rows={3}
                    placeholder="活动描述"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">备注</label>
                  <textarea
                    value={formData.notes}
                    onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                    rows={2}
                    placeholder="其他备注信息"
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
                    {editingActivity ? '保存修改' : '添加活动'}
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
