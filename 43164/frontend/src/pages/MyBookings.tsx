import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { bookingApi, videoApi, homeworkApi } from '@/services/api'
import { Booking } from '@/types'
import { useAuthStore } from '@/store/auth'
import { Calendar, Clock, Video, Star, XCircle, CheckCircle, BookOpen } from 'lucide-react'

export default function MyBookings() {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [bookings, setBookings] = useState<Booking[]>([])
  const [activeTab, setActiveTab] = useState('upcoming')
  const [loading, setLoading] = useState(true)
  const [showHomeworkModal, setShowHomeworkModal] = useState(false)
  const [selectedBooking, setSelectedBooking] = useState<Booking | null>(null)
  const [homeworkForm, setHomeworkForm] = useState({
    title: '',
    description: '',
    dueDate: '',
    maxScore: 100,
  })

  useEffect(() => {
    loadBookings()
  }, [activeTab])

  const loadBookings = async () => {
    try {
      setLoading(true)
      const res = await bookingApi.getAll({ status: activeTab === 'upcoming' ? undefined : activeTab })
      setBookings(res.data)
    } catch (error) {
      console.error('Failed to load bookings:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = async (bookingId: string) => {
    if (confirm('确定要取消这个预约吗？')) {
      try {
        await bookingApi.cancel({ bookingId, reason: '用户取消' })
        toast.success('预约已取消')
        loadBookings()
      } catch (error: any) {
        toast.error(error.response?.data?.error || '取消失败')
      }
    }
  }

  const handleJoinSession = async (booking: Booking) => {
    try {
      const res = await videoApi.createSession({ bookingId: booking.id })
      const session = res.data.session
      window.open(session.joinUrl, '_blank')
    } catch (error) {
      toast.error('无法创建视频会话')
    }
  }

  const handleComplete = async (bookingId: string) => {
    if (confirm('确定要完成这个课程吗？完成后将自动结算费用。')) {
      try {
        await bookingApi.complete(bookingId)
        toast.success('课程已完成，费用已结算')
        loadBookings()
      } catch (error: any) {
        toast.error(error.response?.data?.error || '完成课程失败')
      }
    }
  }

  const handleOpenHomeworkModal = (booking: Booking) => {
    setSelectedBooking(booking)
    setShowHomeworkModal(true)
    setHomeworkForm({
      title: '',
      description: '',
      dueDate: '',
      maxScore: 100,
    })
  }

  const handleCreateHomework = async () => {
    if (!selectedBooking || !homeworkForm.title || !homeworkForm.description || !homeworkForm.dueDate) {
      toast.error('请填写完整的作业信息')
      return
    }

    try {
      await homeworkApi.create({
        bookingId: selectedBooking.id,
        subjectId: selectedBooking.subjectId,
        title: homeworkForm.title,
        description: homeworkForm.description,
        dueDate: new Date(homeworkForm.dueDate).toISOString(),
        maxScore: homeworkForm.maxScore,
      })
      toast.success('作业已布置成功')
      setShowHomeworkModal(false)
    } catch (error: any) {
      toast.error(error.response?.data?.error || '布置作业失败')
    }
  }

  const tabs = [
    { id: 'upcoming', label: '即将开始' },
    { id: 'confirmed', label: '已确认' },
    { id: 'completed', label: '已完成' },
    { id: 'cancelled', label: '已取消' },
  ]

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'bg-yellow-100 text-yellow-800'
      case 'confirmed': return 'bg-green-100 text-green-800'
      case 'completed': return 'bg-blue-100 text-blue-800'
      case 'cancelled': return 'bg-red-100 text-red-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  const getStatusText = (status: string) => {
    switch (status) {
      case 'pending': return '待确认'
      case 'confirmed': return '已确认'
      case 'completed': return '已完成'
      case 'cancelled': return '已取消'
      default: return status
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
        <h1 className="text-2xl font-bold text-gray-900">我的课程</h1>
        <p className="text-gray-500">查看和管理您的课程预约</p>
      </div>

      <div className="flex gap-2 border-b border-gray-200">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors ${
              activeTab === tab.id
                ? 'border-primary-600 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {bookings.length > 0 ? (
        <div className="space-y-4">
          {bookings.map((booking) => (
            <div key={booking.id} className="card">
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <span className={`badge ${getStatusColor(booking.status)}`}>
                      {getStatusText(booking.status)}
                    </span>
                    <span className="text-sm text-gray-500">
                      {new Date(booking.startTime).toLocaleDateString('zh-CN', {
                        year: 'numeric',
                        month: 'long',
                        day: 'numeric',
                        weekday: 'long',
                      })}
                    </span>
                  </div>

                  <h3 className="text-lg font-semibold text-gray-900 mb-2">
                    {booking.subject?.name}
                  </h3>

                  <div className="flex items-center gap-4 text-sm text-gray-500">
                    <div className="flex items-center gap-1">
                      <Clock className="h-4 w-4" />
                      {new Date(booking.startTime).toLocaleTimeString('zh-CN', {
                        hour: '2-digit',
                        minute: '2-digit',
                      })} - {new Date(booking.endTime).toLocaleTimeString('zh-CN', {
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                    </div>
                    <div className="flex items-center gap-1">
                      <Calendar className="h-4 w-4" />
                      {booking.duration}分钟
                    </div>
                    <div className="font-medium text-primary-600">
                      ${booking.totalAmount.toFixed(2)}
                    </div>
                  </div>

                  {booking.notes && (
                    <p className="text-sm text-gray-500 mt-2">备注: {booking.notes}</p>
                  )}
                </div>

                <div className="flex gap-2 flex-wrap">
                  {booking.status === 'pending' && (
                    <button
                      onClick={() => handleCancel(booking.id)}
                      className="btn-secondary flex items-center gap-1"
                    >
                      <XCircle className="h-4 w-4" />
                      取消
                    </button>
                  )}
                  {booking.status === 'confirmed' && (
                    <>
                      <button
                        onClick={() => handleJoinSession(booking)}
                        className="btn-primary flex items-center gap-1"
                      >
                        <Video className="h-4 w-4" />
                        进入教室
                      </button>
                      <button
                        onClick={() => handleComplete(booking.id)}
                        className="btn-secondary flex items-center gap-1 bg-green-600 hover:bg-green-700 text-white border-green-600"
                      >
                        <CheckCircle className="h-4 w-4" />
                        完成课程
                      </button>
                      <button
                        onClick={() => handleCancel(booking.id)}
                        className="btn-secondary flex items-center gap-1"
                      >
                        <XCircle className="h-4 w-4" />
                        取消
                      </button>
                    </>
                  )}
                  {booking.status === 'completed' && (
                    <>
                      <button
                        onClick={() => navigate('/reviews')}
                        className="btn-primary flex items-center gap-1"
                      >
                        <Star className="h-4 w-4" />
                        评价
                      </button>
                      {user?.role === 'teacher' && (
                        <button
                          onClick={() => handleOpenHomeworkModal(booking)}
                          className="btn-secondary flex items-center gap-1"
                        >
                          <BookOpen className="h-4 w-4" />
                          布置作业
                        </button>
                      )}
                    </>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-12">
          <Calendar className="h-12 w-12 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">暂无课程记录</p>
          <button
            onClick={() => navigate('/teachers')}
            className="mt-4 text-primary-600 hover:text-primary-700 font-medium"
          >
            去预约课程 →
          </button>
        </div>
      )}

      {showHomeworkModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">布置作业</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">作业标题</label>
                <input
                  type="text"
                  value={homeworkForm.title}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, title: e.target.value })}
                  className="input-field"
                  placeholder="请输入作业标题"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">作业描述</label>
                <textarea
                  value={homeworkForm.description}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, description: e.target.value })}
                  className="input-field min-h-[100px]"
                  placeholder="请输入作业描述"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">截止日期</label>
                <input
                  type="datetime-local"
                  value={homeworkForm.dueDate}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, dueDate: e.target.value })}
                  className="input-field"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">满分</label>
                <input
                  type="number"
                  value={homeworkForm.maxScore}
                  onChange={(e) => setHomeworkForm({ ...homeworkForm, maxScore: Number(e.target.value) })}
                  className="input-field"
                  min="1"
                />
              </div>
            </div>
            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowHomeworkModal(false)}
                className="btn-secondary flex-1"
              >
                取消
              </button>
              <button
                onClick={handleCreateHomework}
                className="btn-primary flex-1"
              >
                布置作业
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
