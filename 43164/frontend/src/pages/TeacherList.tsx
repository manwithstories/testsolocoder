import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { teacherApi, subjectApi } from '@/services/api'
import { TeacherProfile, Subject } from '@/types'
import { Search, Star, Clock, DollarSign, Filter, ChevronRight } from 'lucide-react'

export default function TeacherList() {
  const navigate = useNavigate()
  const [teachers, setTeachers] = useState<TeacherProfile[]>([])
  const [subjects, setSubjects] = useState<Subject[]>([])
  const [loading, setLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedSubject, setSelectedSubject] = useState('')
  const [minRating, setMinRating] = useState(0)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const [teachersRes, subjectsRes] = await Promise.all([
        teacherApi.list(),
        subjectApi.getAll(),
      ])
      setTeachers(teachersRes.data)
      setSubjects(subjectsRes.data)
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }

  const filteredTeachers = teachers.filter((teacher) => {
    const user = teacher.user
    if (!user) return false
    const nameMatch = `${user.firstName} ${user.lastName}`.toLowerCase().includes(searchQuery.toLowerCase())
    const ratingMatch = teacher.rating >= minRating
    const subjectMatch = !selectedSubject || teacher.subjects?.some(s => s.subjectId === selectedSubject)
    return nameMatch && ratingMatch && subjectMatch
  })

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">寻找老师</h1>
          <p className="text-gray-500">浏览并选择最适合您的老师</p>
        </div>
      </div>

      <div className="card">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-1">搜索</label>
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="input-field pl-10"
                placeholder="搜索老师姓名..."
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">科目</label>
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

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">最低评分</label>
            <select
              value={minRating}
              onChange={(e) => setMinRating(Number(e.target.value))}
              className="input-field"
            >
              <option value={0}>全部</option>
              <option value={4}>4星及以上</option>
              <option value={4.5}>4.5星及以上</option>
              <option value={5}>5星</option>
            </select>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {filteredTeachers.map((teacher) => {
          const user = teacher.user
          if (!user) return null
          return (
            <div
              key={teacher.id}
              className="card hover:shadow-md transition-shadow cursor-pointer"
              onClick={() => navigate(`/teachers/${user.id}`)}
            >
              <div className="flex items-start gap-4">
                {user.avatarUrl ? (
                  <img src={user.avatarUrl} alt="" className="h-16 w-16 rounded-full object-cover" />
                ) : (
                  <div className="h-16 w-16 rounded-full bg-primary-600 flex items-center justify-center text-white text-xl font-medium">
                    {user.firstName[0]}{user.lastName[0]}
                  </div>
                )}
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900">{user.firstName} {user.lastName}</h3>
                  <div className="flex items-center gap-1 text-sm text-yellow-500 mt-1">
                    <Star className="h-4 w-4 fill-current" />
                    <span>{teacher.rating?.toFixed(1) || '5.0'}</span>
                    <span className="text-gray-400">({teacher.reviewCount}条评价)</span>
                  </div>
                  <div className="flex flex-wrap gap-1 mt-2">
                    {teacher.subjects?.slice(0, 3).map((s) => (
                      <span key={s.id} className="badge bg-primary-100 text-primary-800">
                        {s.subject?.name}
                      </span>
                    ))}
                  </div>
                </div>
                <ChevronRight className="h-5 w-5 text-gray-400" />
              </div>

              <div className="mt-4 pt-4 border-t border-gray-100 flex items-center justify-between text-sm">
                <div className="flex items-center gap-1 text-gray-500">
                  <DollarSign className="h-4 w-4" />
                  <span>${teacher.hourlyRate?.toFixed(2) || '0.00'}/小时</span>
                </div>
                <div className="flex items-center gap-1 text-gray-500">
                  <Clock className="h-4 w-4" />
                  <span>{teacher.totalHours || 0}小时授课</span>
                </div>
              </div>
            </div>
          )
        })}
      </div>

      {filteredTeachers.length === 0 && (
        <div className="text-center py-12">
          <Filter className="h-12 w-12 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">没有找到符合条件的老师</p>
        </div>
      )}
    </div>
  )
}
