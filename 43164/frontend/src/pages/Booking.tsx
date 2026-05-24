import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { useForm } from 'react-hook-form'
import { teacherApi, subjectApi, bookingApi } from '@/services/api'
import { TeacherProfile, Subject } from '@/types'
import { Calendar, Clock, DollarSign } from 'lucide-react'

interface BookingForm {
  teacherId: string
  subjectId: string
  startTime: string
  endTime: string
  notes: string
}

export default function BookingPage() {
  const navigate = useNavigate()
  const [teachers, setTeachers] = useState<TeacherProfile[]>([])
  const [subjects, setSubjects] = useState<Subject[]>([])
  const [filteredTeachers, setFilteredTeachers] = useState<TeacherProfile[]>([])
  const [selectedSubject, setSelectedSubject] = useState('')
  const [selectedDate, setSelectedDate] = useState('')
  const [selectedTime, setSelectedTime] = useState('')
  const [duration, setDuration] = useState(60)
  const [loading, setLoading] = useState(true)

  const { register, handleSubmit, formState: { errors } } = useForm<BookingForm>()

  useEffect(() => {
    loadData()
  }, [])

  useEffect(() => {
    if (selectedSubject) {
      const filtered = teachers.filter(teacher =>
        teacher.subjects?.some(s => s.subjectId === selectedSubject)
      )
      setFilteredTeachers(filtered)
    } else {
      setFilteredTeachers(teachers)
    }
  }, [selectedSubject, teachers])

  const loadData = async () => {
    try {
      setLoading(true)
      const [teachersRes, subjectsRes] = await Promise.all([
        teacherApi.list(),
        subjectApi.getAll(),
      ])
      setTeachers(teachersRes.data)
      setSubjects(subjectsRes.data)
      setFilteredTeachers(teachersRes.data)
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const onSubmit = async (data: BookingForm) => {
    if (!selectedDate || !selectedTime) {
      toast.error('请选择日期和时间')
      return
    }

    try {
      const startDateTime = new Date(`${selectedDate}T${selectedTime}:00`)
      const endDateTime = new Date(startDateTime.getTime() + duration * 60000)

      await bookingApi.create({
        ...data,
        startTime: startDateTime.toISOString(),
        endTime: endDateTime.toISOString(),
      })

      toast.success('预约成功！等待老师确认')
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

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">预约课程</h1>
        <p className="text-gray-500">选择老师和时间进行预约</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">选择科目</h2>
          <select
            value={selectedSubject}
            onChange={(e) => setSelectedSubject(e.target.value)}
            className="input-field"
          >
            <option value="">全部科目</option>
            {subjects.map((subject) => (
              <option key={subject.id} value={subject.id}>{subject.name}</option>
            ))}
          </select>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">选择老师</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 max-h-96 overflow-y-auto">
            {filteredTeachers.map((teacher) => {
              const user = teacher.user
              if (!user) return null
              return (
                <label
                  key={teacher.id}
                  className={`p-4 border-2 rounded-lg cursor-pointer transition-colors ${
                    // @ts-ignore
                    errors.teacherId ? 'border-red-300' : 'border-gray-200 hover:border-primary-300'
                  }`}
                >
                  <input
                    type="radio"
                    value={user.id}
                    {...register('teacherId', { required: '请选择老师' })}
                    className="hidden"
                  />
                  <div className="flex items-center gap-3">
                    <div className="w-12 h-12 rounded-full bg-primary-600 flex items-center justify-center text-white">
                      {user.firstName[0]}{user.lastName[0]}
                    </div>
                    <div className="flex-1">
                      <div className="font-medium text-gray-900">
                        {user.firstName} {user.lastName}
                      </div>
                      <div className="flex items-center gap-2 text-sm text-gray-500">
                        <span className="text-yellow-500">★ {teacher.rating?.toFixed(1)}</span>
                        <span>${teacher.hourlyRate?.toFixed(2)}/小时</span>
                      </div>
                    </div>
                  </div>
                </label>
              )
            })}
          </div>
          {errors.teacherId && (
            <p className="text-red-500 text-sm mt-2">{errors.teacherId.message}</p>
          )}
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">选择时间</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                <Calendar className="h-4 w-4 inline mr-1" />
                日期
              </label>
              <input
                type="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                className="input-field"
                min={new Date().toISOString().split('T')[0]}
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                <Clock className="h-4 w-4 inline mr-1" />
                时间
              </label>
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
              <label className="block text-sm font-medium text-gray-700 mb-1">
                <DollarSign className="h-4 w-4 inline mr-1" />
                时长
              </label>
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
          </div>
        </div>

        <div className="card">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">备注</h2>
          <textarea
            {...register('notes')}
            className="input-field h-32 resize-none"
            placeholder="请输入您的学习需求或其他备注..."
          />
        </div>

        <div className="flex gap-4">
          <button
            type="button"
            onClick={() => navigate(-1)}
            className="btn-secondary flex-1"
          >
            取消
          </button>
          <button
            type="submit"
            className="btn-primary flex-1"
          >
            确认预约
          </button>
        </div>
      </form>
    </div>
  )
}
