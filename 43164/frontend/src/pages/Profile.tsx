import { useState, useEffect } from 'react'
import { toast } from 'sonner'
import { userApi } from '@/services/api'
import { useAuthStore } from '@/store/auth'
import { User } from '@/types'
import { Mail, Phone, MapPin, Edit2, Save, X } from 'lucide-react'

export default function Profile() {
  const { user, setUser } = useAuthStore()
  const [profile, setProfile] = useState<User | null>(null)
  const [editing, setEditing] = useState(false)
  const [loading, setLoading] = useState(true)
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    phone: '',
    avatarUrl: '',
    timezone: 'UTC',
  })

  useEffect(() => {
    loadProfile()
  }, [])

  const loadProfile = async () => {
    try {
      setLoading(true)
      const res = await userApi.getProfile()
      setProfile(res.data)
      setFormData({
        firstName: res.data.firstName,
        lastName: res.data.lastName,
        phone: res.data.phone || '',
        avatarUrl: res.data.avatarUrl || '',
        timezone: res.data.timezone || 'UTC',
      })
    } catch (error) {
      console.error('Failed to load profile:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    try {
      await userApi.updateProfile(formData)
      toast.success('个人资料已更新')
      setEditing(false)
      if (user) {
        setUser({
          ...user,
          firstName: formData.firstName,
          lastName: formData.lastName,
          phone: formData.phone,
          avatarUrl: formData.avatarUrl,
          timezone: formData.timezone,
        })
      }
      loadProfile()
    } catch (error: any) {
      toast.error(error.response?.data?.error || '更新失败')
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
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">个人资料</h1>
          <p className="text-gray-500">查看和管理您的个人信息</p>
        </div>
        {editing ? (
          <div className="flex gap-2">
            <button onClick={() => setEditing(false)} className="btn-secondary">
              <X className="h-4 w-4 mr-1 inline" />
              取消
            </button>
            <button onClick={handleSave} className="btn-primary">
              <Save className="h-4 w-4 mr-1 inline" />
              保存
            </button>
          </div>
        ) : (
          <button onClick={() => setEditing(true)} className="btn-primary">
            <Edit2 className="h-4 w-4 mr-1 inline" />
            编辑
          </button>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="card text-center">
          <div className="mb-4">
            {profile?.avatarUrl ? (
              <img
                src={profile.avatarUrl}
                alt=""
                className="h-32 w-32 rounded-full mx-auto object-cover"
              />
            ) : (
              <div className="h-32 w-32 rounded-full bg-primary-600 flex items-center justify-center mx-auto text-white text-4xl font-medium">
                {profile?.firstName?.[0]}{profile?.lastName?.[0]}
              </div>
            )}
          </div>
          <h2 className="text-xl font-semibold text-gray-900">
            {profile?.firstName} {profile?.lastName}
          </h2>
          <p className="text-gray-500">{profile?.email}</p>
          <div className="mt-4">
            <span className={`badge ${
              profile?.role === 'teacher' ? 'bg-blue-100 text-blue-800' :
              profile?.role === 'student' ? 'bg-green-100 text-green-800' :
              'bg-purple-100 text-purple-800'
            }`}>
              {profile?.role === 'teacher' ? '老师' :
               profile?.role === 'student' ? '学生' : '管理员'}
            </span>
          </div>
        </div>

        <div className="lg:col-span-2 space-y-6">
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">基本信息</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">名</label>
                {editing ? (
                  <input
                    type="text"
                    value={formData.firstName}
                    onChange={(e) => setFormData({ ...formData, firstName: e.target.value })}
                    className="input-field"
                  />
                ) : (
                  <p className="text-gray-900">{profile?.firstName}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">姓</label>
                {editing ? (
                  <input
                    type="text"
                    value={formData.lastName}
                    onChange={(e) => setFormData({ ...formData, lastName: e.target.value })}
                    className="input-field"
                  />
                ) : (
                  <p className="text-gray-900">{profile?.lastName}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  <Mail className="h-4 w-4 inline mr-1" />
                  邮箱
                </label>
                <p className="text-gray-900">{profile?.email}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  <Phone className="h-4 w-4 inline mr-1" />
                  电话
                </label>
                {editing ? (
                  <input
                    type="text"
                    value={formData.phone}
                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                    className="input-field"
                  />
                ) : (
                  <p className="text-gray-900">{profile?.phone || '未设置'}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  <MapPin className="h-4 w-4 inline mr-1" />
                  时区
                </label>
                {editing ? (
                  <select
                    value={formData.timezone}
                    onChange={(e) => setFormData({ ...formData, timezone: e.target.value })}
                    className="input-field"
                  >
                    <option value="UTC">UTC</option>
                    <option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</option>
                    <option value="America/New_York">America/New_York (UTC-5)</option>
                    <option value="Europe/London">Europe/London (UTC+0)</option>
                    <option value="Australia/Sydney">Australia/Sydney (UTC+11)</option>
                  </select>
                ) : (
                  <p className="text-gray-900">{profile?.timezone}</p>
                )}
              </div>
            </div>
          </div>

          {profile?.role === 'teacher' && profile?.teacherProfile && (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">教师信息</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">时薪</label>
                  <p className="text-gray-900">${profile.teacherProfile.hourlyRate?.toFixed(2)}/小时</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">评分</label>
                  <p className="text-gray-900">
                    {profile.teacherProfile.rating?.toFixed(1)} ({profile.teacherProfile.reviewCount}条评价)
                  </p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">总课时</label>
                  <p className="text-gray-900">{profile.teacherProfile.totalHours}小时</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">认证状态</label>
                  <span className={`badge ${
                    profile.teacherProfile.approvalStatus === 'approved'
                      ? 'bg-green-100 text-green-800'
                      : profile.teacherProfile.approvalStatus === 'pending'
                      ? 'bg-yellow-100 text-yellow-800'
                      : 'bg-red-100 text-red-800'
                  }`}>
                    {profile.teacherProfile.approvalStatus === 'approved' ? '已认证' :
                     profile.teacherProfile.approvalStatus === 'pending' ? '审核中' : '已拒绝'}
                  </span>
                </div>
              </div>
            </div>
          )}

          {profile?.role === 'student' && profile?.studentProfile && (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">学生信息</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">年级</label>
                  <p className="text-gray-900">{profile.studentProfile.gradeLevel || '未设置'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">学校</label>
                  <p className="text-gray-900">{profile.studentProfile.school || '未设置'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">学习风格</label>
                  <p className="text-gray-900">{profile.studentProfile.learningStyle || '未设置'}</p>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">评估状态</label>
                  <p className="text-gray-900">{profile.studentProfile.assessmentStatus}</p>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
