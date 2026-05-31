import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { Plus, Check, Calendar, X } from 'lucide-react'
import { calendarApi } from '@/api'
import dayjs from 'dayjs'

export default function CalendarPage() {
  const queryClient = useQueryClient()
  const [showModal, setShowModal] = useState(false)
  const [currentMonth, setCurrentMonth] = useState(dayjs())
  const [formData, setFormData] = useState({
    title: '',
    event_type: '',
    event_date: dayjs().format('YYYY-MM-DD'),
    description: '',
  })

  const { data, isLoading } = useQuery({
    queryKey: ['calendar-events', currentMonth.format('YYYY-MM')],
    queryFn: () =>
      calendarApi.getEvents({
        start_date: currentMonth.startOf('month').format('YYYY-MM-DD'),
        end_date: currentMonth.endOf('month').format('YYYY-MM-DD'),
      }),
  })

  const { data: recommendations } = useQuery({
    queryKey: ['recommendations'],
    queryFn: () => calendarApi.getRecommendations(),
  })

  const createMutation = useMutation({
    mutationFn: (data: object) => calendarApi.createEvent(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
      setShowModal(false)
      setFormData({
        title: '',
        event_type: '',
        event_date: dayjs().format('YYYY-MM-DD'),
        description: '',
      })
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: object }) => calendarApi.updateEvent(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => calendarApi.deleteEvent(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['calendar-events'] })
    },
  })

  const events = data?.data?.events || []

  const getDaysInMonth = () => {
    const start = currentMonth.startOf('month')
    const end = currentMonth.endOf('month')
    const days = []
    let current = start

    const startDay = start.day()
    for (let i = 0; i < startDay; i++) {
      days.push(null)
    }

    while (current.isBefore(end) || current.isSame(end, 'day')) {
      days.push(current)
      current = current.add(1, 'day')
    }

    return days
  }

  const getEventsForDate = (date: dayjs.Dayjs | null) => {
    if (!date) return []
    return events.filter((event: any) => dayjs(event.event_date).isSame(date, 'day'))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    createMutation.mutate(formData)
  }

  const handleToggleComplete = (event: any) => {
    updateMutation.mutate({
      id: event.id,
      data: { is_completed: !event.is_completed },
    })
  }

  const handleDelete = (id: string) => {
    if (confirm('确定要删除这个事件吗？')) {
      deleteMutation.mutate(id)
    }
  }

  const weekDays = ['日', '一', '二', '三', '四', '五', '六']
  const days = getDaysInMonth()

  const eventTypeColors: Record<string, string> = {
    planting: 'bg-green-100 text-green-700',
    watering: 'bg-blue-100 text-blue-700',
    fertilizing: 'bg-amber-100 text-amber-700',
    harvesting: 'bg-purple-100 text-purple-700',
    reminder: 'bg-gray-100 text-gray-700',
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">种植日历</h1>
          <p className="text-gray-500">管理您的种植计划和提醒</p>
        </div>
        <button onClick={() => setShowModal(true)} className="btn-primary">
          <Plus className="w-4 h-4 mr-2" />
          添加事件
        </button>
      </div>

      {/* Month Navigation */}
      <div className="card">
        <div className="card-header flex items-center justify-between">
          <button
            onClick={() => setCurrentMonth(currentMonth.subtract(1, 'month'))}
            className="p-2 hover:bg-gray-100 rounded-lg"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h2 className="text-lg font-semibold text-gray-900">
            {currentMonth.format('YYYY年MM月')}
          </h2>
          <button
            onClick={() => setCurrentMonth(currentMonth.add(1, 'month'))}
            className="p-2 hover:bg-gray-100 rounded-lg"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
        <div className="p-4">
          {/* Week Day Headers */}
          <div className="grid grid-cols-7 gap-1 mb-2">
            {weekDays.map((day) => (
              <div key={day} className="text-center text-sm font-medium text-gray-500 py-2">
                {day}
              </div>
            ))}
          </div>

          {/* Calendar Grid */}
          <div className="grid grid-cols-7 gap-1">
            {days.map((date, index) => (
              <div
                key={index}
                className={`min-h-24 p-2 rounded-lg border ${
                  date
                    ? dayjs().isSame(date, 'day')
                      ? 'border-garden-400 bg-garden-50'
                      : 'border-gray-200 hover:border-garden-300'
                    : 'border-transparent'
                }`}
              >
                {date && (
                  <>
                    <p className="text-sm font-medium text-gray-900">{date.date()}</p>
                    <div className="mt-1 space-y-1">
                      {getEventsForDate(date).map((event: any) => (
                        <div
                          key={event.id}
                          className={`text-xs p-1 rounded ${eventTypeColors[event.event_type] || 'bg-gray-100'} ${
                            event.is_completed ? 'opacity-50' : ''
                          }`}
                        >
                          <div className="flex items-center justify-between">
                            <span className="truncate">{event.title}</span>
                            <button
                              onClick={() => handleToggleComplete(event)}
                              className="ml-1"
                            >
                              {event.is_completed ? (
                                <Check className="w-3 h-3" />
                              ) : (
                                <div className="w-3 h-3 rounded-full border border-current" />
                              )}
                            </button>
                          </div>
                        </div>
                      ))}
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Recommendations */}
      {recommendations?.data?.recommendations && recommendations.data.recommendations.length > 0 && (
        <div className="card">
          <div className="card-header">
            <h2 className="text-lg font-semibold text-gray-900">本月推荐种植</h2>
          </div>
          <div className="card-body">
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
              {recommendations.data.recommendations.map((plant: any) => (
                <div key={plant.id} className="text-center">
                  <div className="aspect-square bg-garden-50 rounded-xl flex items-center justify-center mb-2">
                    {plant.image_url ? (
                      <img src={plant.image_url} alt={plant.name} className="w-full h-full object-cover rounded-xl" />
                    ) : (
                      <Calendar className="w-8 h-8 text-garden-400" />
                    )}
                  </div>
                  <p className="text-sm font-medium text-gray-900 truncate">{plant.name}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Event List */}
      <div className="card">
        <div className="card-header">
          <h2 className="text-lg font-semibold text-gray-900">本月事件</h2>
        </div>
        <div className="divide-y divide-gray-100">
          {events.length === 0 ? (
            <div className="p-6 text-center text-gray-500">本月暂无事件</div>
          ) : (
            events.map((event: any) => (
              <div key={event.id} className="px-6 py-4 flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <button onClick={() => handleToggleComplete(event)}>
                    {event.is_completed ? (
                      <div className="w-5 h-5 rounded-full bg-garden-500 flex items-center justify-center">
                        <Check className="w-3 h-3 text-white" />
                      </div>
                    ) : (
                      <div className="w-5 h-5 rounded-full border-2 border-gray-300" />
                    )}
                  </button>
                  <div>
                    <p className={`font-medium ${event.is_completed ? 'text-gray-400 line-through' : 'text-gray-900'}`}>
                      {event.title}
                    </p>
                    <p className="text-sm text-gray-500">
                      {dayjs(event.event_date).format('MM月DD日')}
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => handleDelete(event.id)}
                  className="p-2 hover:bg-red-50 rounded-lg"
                >
                  <X className="w-4 h-4 text-red-400" />
                </button>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Modal */}
      {showModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl w-full max-w-md">
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h2 className="text-lg font-semibold">添加事件</h2>
              <button
                onClick={() => setShowModal(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div>
                <label className="label">事件标题 *</label>
                <input
                  type="text"
                  className="input"
                  required
                  value={formData.title}
                  onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                />
              </div>
              <div>
                <label className="label">事件类型</label>
                <select
                  className="input"
                  value={formData.event_type}
                  onChange={(e) => setFormData({ ...formData, event_type: e.target.value })}
                >
                  <option value="">选择类型</option>
                  <option value="planting">种植</option>
                  <option value="watering">浇水</option>
                  <option value="fertilizing">施肥</option>
                  <option value="harvesting">收获</option>
                  <option value="reminder">提醒</option>
                </select>
              </div>
              <div>
                <label className="label">日期 *</label>
                <input
                  type="date"
                  className="input"
                  required
                  value={formData.event_date}
                  onChange={(e) => setFormData({ ...formData, event_date: e.target.value })}
                />
              </div>
              <div>
                <label className="label">描述</label>
                <textarea
                  className="input h-24"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowModal(false)}
                  className="btn-outline flex-1"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={createMutation.isPending}
                  className="btn-primary flex-1"
                >
                  {createMutation.isPending ? '添加中...' : '添加'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
