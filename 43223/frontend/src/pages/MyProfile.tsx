import { useState } from 'react'
import { useAuthStore } from '@/store/auth'
import { authApi } from '@/api/auth'

export default function MyProfile() {
  const { user, setUser } = useAuthStore()
  const [formData, setFormData] = useState({
    nickname: user?.nickname || '',
    phone: user?.phone || '',
    avatar: user?.avatar || '',
    bio: user?.bio || '',
    address: user?.address || '',
  })
  const [saving, setSaving] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      const res = await authApi.updateProfile(formData)
      if (res.data) {
        setUser(res.data)
        alert('保存成功')
      }
    } catch (err: any) {
      alert(err.message || '保存失败')
    } finally {
      setSaving(false)
    }
  }

  if (!user) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">个人中心</h1>

      <div className="card p-6">
        <div className="flex items-center gap-4 mb-6">
          <div className="w-20 h-20 bg-coffee-200 rounded-full flex items-center justify-center text-3xl">
            {user.avatar ? (
              <img src={user.avatar} alt="" className="w-full h-full rounded-full object-cover" />
            ) : (
              user.nickname?.[0] || user.username[0]
            )}
          </div>
          <div>
            <h2 className="text-xl font-bold">{user.nickname || user.username}</h2>
            <p className="text-gray-500">@{user.username}</p>
            <div className="flex gap-2 mt-1">
              <span className="badge bg-coffee-100 text-coffee-700">
                {user.role === 'admin' ? '管理员' : user.role === 'roaster' ? '烘焙师' : '普通用户'}
              </span>
              {user.is_certified && (
                <span className="badge bg-green-100 text-green-700">已认证</span>
              )}
            </div>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="label">昵称</label>
            <input
              type="text"
              value={formData.nickname}
              onChange={(e) => setFormData({ ...formData, nickname: e.target.value })}
              className="input"
            />
          </div>
          <div>
            <label className="label">手机号</label>
            <input
              type="tel"
              value={formData.phone}
              onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              className="input"
            />
          </div>
          <div>
            <label className="label">头像URL</label>
            <input
              type="text"
              value={formData.avatar}
              onChange={(e) => setFormData({ ...formData, avatar: e.target.value })}
              className="input"
            />
          </div>
          <div>
            <label className="label">个人简介</label>
            <textarea
              value={formData.bio}
              onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
              className="input"
              rows={3}
            />
          </div>
          <div>
            <label className="label">地址</label>
            <textarea
              value={formData.address}
              onChange={(e) => setFormData({ ...formData, address: e.target.value })}
              className="input"
              rows={2}
            />
          </div>
          <button type="submit" disabled={saving} className="btn btn-primary w-full">
            {saving ? '保存中...' : '保存'}
          </button>
        </form>
      </div>

      <div className="card p-6">
        <h3 className="font-bold mb-4">账户信息</h3>
        <div className="space-y-2 text-sm">
          <p><span className="text-gray-500">邮箱:</span> {user.email}</p>
          <p><span className="text-gray-500">注册时间:</span> {user.created_at?.split('T')[0]}</p>
        </div>
      </div>
    </div>
  )
}
