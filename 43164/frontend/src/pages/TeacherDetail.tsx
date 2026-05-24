import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { teacherApi, bookingApi } from '@/services/api'
import { TeacherProfile } from '@/types'
import { Star, Clock, DollarSign, Calendar, CheckCircle } from 'lucide-react'

export default function TeacherDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [teacher, setTeacher] = useState<TeacherProfile | null>(null)
  const [reviews, setReviews] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showBookingModal, setShowBookingModal] = useState(false)
  const [selectedDate, setSelectedDate] = useState('')
  const [selectedTime, setSelectedTime] = useState('')
  const [duration, setDuration] = useState(60)
  const [notes, setNotes] = useState('')

  useEffect(() => {
    if (id) {
      loadTeacherData()
    }
  }, [id])

  const loadTeacherData = async () => {
    try {
      const [teacherRes, reviewsRes] = await Promise.all([
        teacherApi.getById(id!),
        teacherApi.getReviews(id!),
      ])
      setTeacher(teacherRes.data)
      setReviews(reviewsRes.data?.reviews || [])
    } catch (error) {
      console.error('Failed to load teacher data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleBooking = async () => {
    if (!selectedDate || !selectedTime) {
      toast.error('请选择日期和时间')
      return
    }

    try {
      const startTime = new Date(`${selectedDate}T${selectedTime}:00`)
      const endTime = new Date(startTime.getTime() + duration * 60000)

      await bookingApi.create({
        teacherId: id,
        subjectId: teacher?.subjects?.[0]?.subjectId,
        startTime: startTime.toISOString(),
        endTime: endTime.toISOString(),
        notes,
      })

      toast.success('预约成功！等待老师确认')
      setShowBookingModal(false)
      navigate('/my-bookings')
    } catch (error: any) {
      toast.error(error.response?.data?.error || '预约失败')
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  if (!teacher) {
    return <div className="text-center py-12 text-gray-500">老师不存在</div>
  }

  const user = teacher.user

  return (
    <div className="space-y-6">
      <button
        onClick={() => navigate(-1)}
        className="text-primary-600 hover:text-primary-700 flex items-center gap-1"
      >
        ← 返回
      </button>

      <div className="card">
        <div className="flex flex-col md:flex-row gap-6">
          <div className="flex-shrink-0">
            {user?.avatarUrl ? (
              <img src={user.avatarUrl} alt="" className="h-32 w-32 rounded-full object-cover" />
            ) : (
              <div className="h-32 w-32 rounded-full bg-primary-600 flex items-center justify-center text-white text-4xl font-medium">
                {user?.firstName[0]}{user?.lastName[0]}
              </div>
            )}
          </div>

          <div className="flex-1">
            <div className="flex items-center gap-2 mb-2">
              <h1 className="text-2xl font-bold text-gray-900">{user?.firstName} {user?.lastName}</h1>
              {teacher.isVerified && (
                <span className="badge bg-green-100 text-green-800">
                  <CheckCircle className="h-3 w-3 mr-1" />
                  已认证
                </span>
              )}
            </div>

            <div className="flex items-center gap-4 text-sm text-gray-500 mb-4">
              <div className="flex items-center gap-1">
                <Star className="h-4 w-4 text-yellow-500 fill-current" />
                <span className="font-medium text-gray-900">{teacher.rating?.toFixed(1)}</span>
                <span>({teacher.reviewCount}条评价)</span>
              </div>
              <div className="flex items-center gap-1">
                <DollarSign className="h-4 w-4" />
                <span>${teacher.hourlyRate?.toFixed(2)}/小时</span>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="h-4 w-4" />
                <span>{teacher.totalHours}小时授课</span>
              </div>
            </div>

            <div className="flex flex-wrap gap-2 mb-4">
              {teacher.subjects?.map((s) => (
                <span key={s.id} className="badge bg-primary-100 text-primary-800 px-3 py-1">
                  {s.subject?.name} - {s.level}
                </span>
              ))}
            </div>

            {teacher.bio && (
              <div className="mb-4">
                <h3 className="font-medium text-gray-900 mb-1">简介</h3>
                <p className="text-gray-600">{teacher.bio}</p>
              </div>
            )}

            <button
              onClick={() => setShowBookingModal(true)}
              className="btn-primary"
            >
              <Calendar className="h-4 w-4 mr-2 inline" />
              立即预约
            </button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">教育背景</h2>
          <p className="text-gray-600">{teacher.education || '暂无信息'}</p>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">教学经验</h2>
          <p className="text-gray-600">{teacher.experience || '暂无信息'}</p>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">资质认证</h2>
          <p className="text-gray-600">{teacher.certifications || '暂无信息'}</p>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">可预约时段</h2>
          <div className="space-y-2">
            {teacher.availabilities?.length ? (
              teacher.availabilities.map((slot) => (
                <div key={slot.id} className="flex items-center justify-between py-2 border-b border-gray-100">
                  <span className="text-gray-600">
                    {['周日', '周一', '周二', '周三', '周四', '周五', '周六'][slot.dayOfWeek]}
                  </span>
                  <span className="text-primary-600 font-medium">{slot.startTime} - {slot.endTime}</span>
                </div>
              ))
            ) : (
              <p className="text-gray-500">暂无可预约时段</p>
            )}
          </div>
        </div>
      </div>

      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">学生评价</h2>
        {reviews.length > 0 ? (
          <div className="space-y-4">
            {reviews.map((review) => (
              <div key={review.id} className="pb-4 border-b border-gray-100 last:border-0">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <div className="h-8 w-8 rounded-full bg-primary-600 flex items-center justify-center text-white text-sm">
                      {review.reviewer?.firstName?.[0]}
                    </div>
                    <span className="font-medium text-gray-900">
                      {review.isAnonymous ? '匿名用户' : `${review.reviewer?.firstName} ${review.reviewer?.lastName}`}
                    </span>
                  </div>
                  <div className="flex items-center gap-1">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={`h-4 w-4 ${i < review.rating ? 'text-yellow-500 fill-current' : 'text-gray-300'}`}
                      />
                    ))}
                  </div>
                </div>
                <p className="text-gray-600">{review.content}</p>
                {review.teacherReply && (
                  <div className="mt-2 p-3 bg-gray-50 rounded-lg">
                    <div className="text-sm font-medium text-gray-900 mb-1">老师回复:</div>
                    <p className="text-sm text-gray-600">{review.teacherReply}</p>
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 text-center py-8">暂无评价</p>
        )}
      </div>

      {showBookingModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl p-6 w-full max-w-md mx-4">
            <h3 className="text-xl font-bold text-gray-900 mb-4">预约课程</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">选择日期</label>
                <input
                  type="date"
                  value={selectedDate}
                  onChange={(e) => setSelectedDate(e.target.value)}
                  className="input-field"
                  min={new Date().toISOString().split('T')[0]}
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">选择时间</label>
                <select
                  value={selectedTime}
                  onChange={(e) => setSelectedTime(e.target.value)}
                  className="input-field"
                >
                  <option value="">请选择时间</option>
                  {['09:00', '10:00', '11:00', '14:00', '15:00', '16:00', '19:00', '20:00'].map((time) => (
                    <option key={time} value={time}>{time}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">课程时长</label>
                <select
                  value={duration}
                  onChange={(e) => setDuration(Number(e.target.value))}
                  className="input-field"
                >
                  <option value={30}>30分钟</option>
                  <option value={60}>60分钟</option>
                  <option value={90}>90分钟</option>
                  <option value={120}>120分钟</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">备注</label>
                <textarea
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  className="input-field h-24 resize-none"
                  placeholder="请输入您的学习需求..."
                />
              </div>

              <div className="bg-gray-50 rounded-lg p-4">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">课程费用:</span>
                  <span className="font-medium text-gray-900">
                    ${(teacher.hourlyRate * duration / 60).toFixed(2)}
                  </span>
                </div>
              </div>
            </div>

            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowBookingModal(false)}
                className="btn-secondary flex-1"
              >
                取消
              </button>
              <button
                onClick={handleBooking}
                className="btn-primary flex-1"
              >
                确认预约
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
