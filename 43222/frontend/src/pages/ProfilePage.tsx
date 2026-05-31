import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { User, Camera, Edit2, Key, Save, X } from 'lucide-react'
import { userApi, uploadApi } from '@/api'
import { useAuthStore } from '@/store/auth'

export default function ProfilePage() {
  const queryClient = useQueryClient()
  const user = useAuthStore((state) => state.user)
  const setUser = useAuthStore((state) => state.setUser)
  const [isEditing, setIsEditing] = useState(false)
  const [showPasswordModal, setShowPasswordModal] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [formData, setFormData] = useState({
    nickname: '',
    avatar: '',
    phone: '',
    region: '',
    climate_zone: '',
  })
  const [passwordData, setPasswordData] = useState({
    old_password: '',
    new_password: '',
  })

  const { data, isLoading } = useQuery({
    queryKey: ['profile'],
    queryFn: () => userApi.getProfile(),
  })

  const updateMutation = useMutation({
    mutationFn: (data: object) => userApi.updateProfile(data),
    onSuccess: (response) => {
      setUser(response.data.user)
      setIsEditing(false)
    },
  })

  const passwordMutation = useMutation({
    mutationFn: (data: object) => userApi.changePassword(data),
    onSuccess: () => {
      setShowPasswordModal(false)
      setPasswordData({ old_password: '', new_password: '' })
      alert('密码修改成功')
    },
  })

  const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return

    setUploading(true)
    try {
      const response = await uploadApi.upload(file)
      setFormData({ ...formData, avatar: response.data.url })
    } catch (err) {
      alert('上传失败')
    } finally {
      setUploading(false)
    }
  }

  const startEditing = () => {
    setFormData({
      nickname: user?.nickname || '',
      avatar: user?.avatar || '',
      phone: user?.phone || '',
      region: user?.region || '',
      climate_zone: user?.climate_zone || '',
    })
    setIsEditing(true)
  }

  const handleSave = () => {
    updateMutation.mutate(formData)
  }

  const handlePasswordSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    passwordMutation.mutate(passwordData)
  }

  if (isLoading) {
    return <div className="text-center py-12 text-gray-500">加载中...</div>
  }

  const profile = data?.data?.user || user

  const userTypeLabels: Record<string, string> = {
    hobbyist: '园艺爱好者',
    expert: '园艺专家',
    merchant: '种子商家',
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">个人中心</h1>
          <p className="text-gray-500">管理您的个人信息</p>
        </div>
        {!isEditing && (
          <button onClick={startEditing} className="btn-primary">
            <Edit2 className="w-4 h-4 mr-2" />
            编辑资料
          </button>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Profile Card */}
        <div className="card">
          <div className="card-body text-center">
            <div className="relative w-24 h-24 mx-auto mb-4">
              <div className="w-24 h-24 rounded-full bg-garden-100 flex items-center justify-center overflow-hidden">
                {profile?.avatar ? (
                  <img src={profile.avatar} alt="头像" className="w-full h-full object-cover" />
                ) : (
                  <User className="w-12 h-12 text-garden-600" />
                )}
              </div>
              {isEditing && (
                <label className="absolute bottom-0 right-0 p-2 bg-garden-600 text-white rounded-full cursor-pointer hover:bg-garden-700">
                  <Camera className="w-4 h-4" />
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleAvatarUpload}
                    className="hidden"
                  />
                </label>
              )}
            </div>
            {uploading && <p className="text-sm text-garden-600">上传中...</p>}

            <h2 className="text-xl font-semibold text-gray-900">
              {profile?.nickname || profile?.username}
            </h2>
            <p className="text-gray-500">{profile?.email}</p>

            <div className="mt-4 flex justify-center gap-2">
              <span className="badge bg-garden-100 text-garden-700">
                {userTypeLabels[profile?.user_type] || profile?.user_type}
              </span>
              {profile?.is_verified && (
                <span className="badge bg-blue-100 text-blue-700">已认证</span>
              )}
            </div>

            <div className="mt-6 pt-6 border-t border-gray-100 space-y-3 text-left">
              <div className="flex justify-between">
                <span className="text-gray-500">信用评分</span>
                <span className="font-medium text-garden-600">{profile?.credit_score}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500">注册时间</span>
                <span className="font-medium">
                  {profile?.created_at && new Date(profile.created_at).toLocaleDateString('zh-CN')}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Profile Details */}
        <div className="lg:col-span-2 card">
          <div className="card-header flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-900">个人信息</h3>
            {!isEditing && (
              <button
                onClick={() => setShowPasswordModal(true)}
                className="text-sm text-garden-600 hover:text-garden-700"
              >
                修改密码
              </button>
            )}
          </div>
          <div className="card-body">
            {isEditing ? (
              <div className="space-y-4">
                <div>
                  <label className="label">昵称</label>
                  <input
                    type="text"
                    className="input"
                    value={formData.nickname}
                    onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
                  />
                </div>
                <div>
                  <label className="label">手机号</label>
                  <input
                    type="tel"
                    className="input"
                    value={formData.phone}
                    onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                  />
                </div>
                <div>
                  <label className="label">所在地区</label>
                  <input
                    type="text"
                    className="input"
                    value={formData.region}
                    onChange={(e) => setFormData({ ...formData, region: e.target.value })}
                  />
                </div>
                <div>
                  <label className="label">气候带</label>
                  <select
                    className="input"
                    value={formData.climate_zone}
                    onChange={(e) => setFormData({ ...formData, climate_zone: e.target.value })}
                  >
                    <option value="">选择气候带</option>
                    <option value="tropical">热带</option>
                    <option value="subtropical">亚热带</option>
                    <option value="temperate">温带</option>
                    <option value="cold">寒带</option>
                  </select>
                </div>
                <div className="flex gap-3 pt-4">
                  <button
                    onClick={() => setIsEditing(false)}
                    className="btn-outline flex-1"
                  >
                    取消
                  </button>
                  <button
                    onClick={handleSave}
                    disabled={updateMutation.isPending}
                    className="btn-primary flex-1"
                  >
                    <Save className="w-4 h-4 mr-2" />
                    {updateMutation.isPending ? '保存中...' : '保存'}
                  </button>
                </div>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-gray-500">用户名</p>
                    <p className="font-medium text-gray-900">{profile?.username}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">邮箱</p>
                    <p className="font-medium text-gray-900">{profile?.email}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">昵称</p>
                    <p className="font-medium text-gray-900">{profile?.nickname || '-'}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">手机号</p>
                    <p className="font-medium text-gray-900">{profile?.phone || '-'}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">所在地区</p>
                    <p className="font-medium text-gray-900">{profile?.region || '-'}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">气候带</p>
                    <p className="font-medium text-gray-900">{profile?.climate_zone || '-'}</p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Password Modal */}
      {showPasswordModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-2xl w-full max-w-md">
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h2 className="text-lg font-semibold">修改密码</h2>
              <button
                onClick={() => setShowPasswordModal(false)}
                className="p-2 hover:bg-gray-100 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <form onSubmit={handlePasswordSubmit} className="p-6 space-y-4">
              <div>
                <label className="label">当前密码</label>
                <input
                  type="password"
                  className="input"
                  required
                  value={passwordData.old_password}
                  onChange={(e) =>
                    setPasswordData({ ...passwordData, old_password: e.target.value })
                  }
                />
              </div>
              <div>
                <label className="label">新密码</label>
                <input
                  type="password"
                  className="input"
                  required
                  minLength={6}
                  value={passwordData.new_password}
                  onChange={(e) =>
                    setPasswordData({ ...passwordData, new_password: e.target.value })
                  }
                />
              </div>
              <div className="flex gap-3 pt-4">
                <button
                  type="button"
                  onClick={() => setShowPasswordModal(false)}
                  className="btn-outline flex-1"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={passwordMutation.isPending}
                  className="btn-primary flex-1"
                >
                  <Key className="w-4 h-4 mr-2" />
                  {passwordMutation.isPending ? '修改中...' : '确认修改'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
