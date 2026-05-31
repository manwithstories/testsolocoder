import { useEffect, useState } from 'react'
import { certificationApi } from '@/api/certification'
import { RoasterCertification } from '@/types'
import { useAuthStore } from '@/store/auth'

export default function MyCertification() {
  const { user } = useAuthStore()
  const [cert, setCert] = useState<RoasterCertification | null>(null)
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({
    cert_name: '',
    cert_number: '',
    org_name: '',
    cert_file: '',
    experience: '',
    specialty: '',
  })

  useEffect(() => {
    loadCertification()
  }, [])

  const loadCertification = async () => {
    try {
      const res = await certificationApi.getMyCertification()
      setCert(res.data || null)
      if (res.data) {
        setFormData({
          cert_name: res.data.cert_name,
          cert_number: res.data.cert_number,
          org_name: res.data.org_name,
          cert_file: res.data.cert_file || '',
          experience: res.data.experience,
          specialty: res.data.specialty || '',
        })
      }
    } catch {
      setCert(null)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (cert) {
        await certificationApi.updateApplication(formData)
        alert('更新成功')
      } else {
        await certificationApi.apply(formData)
        alert('申请已提交')
      }
      setShowForm(false)
      loadCertification()
    } catch (err: any) {
      alert(err.message || '提交失败')
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (!cert && !showForm) {
    return (
      <div className="max-w-2xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold text-gray-800">烘焙师认证</h1>
        <div className="card p-8 text-center">
          <div className="text-6xl mb-4">📜</div>
          <p className="text-gray-600 mb-4">您还没有提交认证申请</p>
          <p className="text-sm text-gray-500 mb-6">
            成为认证烘焙师，展示您的作品和烘焙记录
          </p>
          <button onClick={() => setShowForm(true)} className="btn btn-primary">
            申请认证
          </button>
        </div>
      </div>
    )
  }

  if (cert && !showForm) {
    const statusColors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-700',
      approved: 'bg-green-100 text-green-700',
      rejected: 'bg-red-100 text-red-700',
    }
    const statusLabels: Record<string, string> = {
      pending: '待审核',
      approved: '已通过',
      rejected: '已拒绝',
    }

    return (
      <div className="max-w-2xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold text-gray-800">烘焙师认证</h1>
        <div className="card p-6">
          <div className="flex items-center justify-between mb-6">
            <span className={`badge ${statusColors[cert.status]}`}>{statusLabels[cert.status]}</span>
            {cert.status === 'pending' && (
              <button onClick={() => setShowForm(true)} className="btn btn-outline text-sm">
                修改申请
              </button>
            )}
          </div>

          <div className="space-y-4">
            <div>
              <p className="text-sm text-gray-500">证书名称</p>
              <p className="font-semibold">{cert.cert_name}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">证书编号</p>
              <p className="font-semibold">{cert.cert_number}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">发证机构</p>
              <p className="font-semibold">{cert.org_name}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">烘焙经验</p>
              <p className="font-semibold whitespace-pre-wrap">{cert.experience}</p>
            </div>
            {cert.specialty && (
              <div>
                <p className="text-sm text-gray-500">专长</p>
                <p className="font-semibold">{cert.specialty}</p>
              </div>
            )}
            {cert.review_comment && (
              <div>
                <p className="text-sm text-gray-500">审核意见</p>
                <p className="font-semibold">{cert.review_comment}</p>
              </div>
            )}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">
        {cert ? '修改认证申请' : '申请烘焙师认证'}
      </h1>
      <div className="card p-6">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="label">证书名称 *</label>
            <input
              type="text"
              value={formData.cert_name}
              onChange={(e) => setFormData({ ...formData, cert_name: e.target.value })}
              className="input"
              required
            />
          </div>
          <div>
            <label className="label">证书编号 *</label>
            <input
              type="text"
              value={formData.cert_number}
              onChange={(e) => setFormData({ ...formData, cert_number: e.target.value })}
              className="input"
              required
            />
          </div>
          <div>
            <label className="label">发证机构 *</label>
            <input
              type="text"
              value={formData.org_name}
              onChange={(e) => setFormData({ ...formData, org_name: e.target.value })}
              className="input"
              required
            />
          </div>
          <div>
            <label className="label">证书文件URL</label>
            <input
              type="text"
              value={formData.cert_file}
              onChange={(e) => setFormData({ ...formData, cert_file: e.target.value })}
              className="input"
            />
          </div>
          <div>
            <label className="label">烘焙经验 *</label>
            <textarea
              value={formData.experience}
              onChange={(e) => setFormData({ ...formData, experience: e.target.value })}
              className="input"
              rows={4}
              required
            />
          </div>
          <div>
            <label className="label">专长</label>
            <input
              type="text"
              value={formData.specialty}
              onChange={(e) => setFormData({ ...formData, specialty: e.target.value })}
              className="input"
            />
          </div>
          <div className="flex gap-4">
            <button type="button" onClick={() => setShowForm(false)} className="btn btn-secondary flex-1">
              取消
            </button>
            <button type="submit" className="btn btn-primary flex-1">
              {cert ? '更新' : '提交'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
