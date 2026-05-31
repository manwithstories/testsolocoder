import { useEffect, useState } from 'react'
import { certificationApi } from '@/api/certification'
import { RoasterCertification } from '@/types'

export default function AdminCertifications() {
  const [certs, setCerts] = useState<RoasterCertification[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [status, setStatus] = useState('')

  useEffect(() => {
    loadCerts()
  }, [page, status])

  const loadCerts = async () => {
    setLoading(true)
    try {
      const res = await certificationApi.list({ page, page_size: 10, status: status || undefined })
      setCerts(res.data?.items || [])
      setTotal(res.data?.total || 0)
    } catch {
      setCerts([])
    } finally {
      setLoading(false)
    }
  }

  const handleReview = async (id: number, reviewStatus: string, comment?: string) => {
    try {
      await certificationApi.review(id, { status: reviewStatus, review_comment: comment })
      alert('审核完成')
      loadCerts()
    } catch (err: any) {
      alert(err.message || '审核失败')
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">认证审核</h1>

      <div className="card p-4">
        <div className="flex gap-4">
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="input w-48"
          >
            <option value="">全部状态</option>
            <option value="pending">待审核</option>
            <option value="approved">已通过</option>
            <option value="rejected">已拒绝</option>
          </select>
        </div>
      </div>

      <div className="space-y-4">
        {certs.length === 0 ? (
          <p className="text-center py-12 text-gray-500">暂无认证申请</p>
        ) : (
          certs.map((cert) => (
            <div key={cert.id} className="card p-6">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-coffee-200 rounded-full flex items-center justify-center">
                    {cert.user?.nickname?.[0] || cert.user?.username[0] || 'U'}
                  </div>
                  <div>
                    <p className="font-bold">{cert.user?.nickname || cert.user?.username}</p>
                    <p className="text-sm text-gray-500">@{cert.user?.username}</p>
                  </div>
                </div>
                <span className={`badge ${
                  cert.status === 'pending' ? 'bg-yellow-100 text-yellow-700' :
                  cert.status === 'approved' ? 'bg-green-100 text-green-700' :
                  'bg-red-100 text-red-700'
                }`}>
                  {cert.status === 'pending' ? '待审核' :
                   cert.status === 'approved' ? '已通过' : '已拒绝'}
                </span>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
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
                {cert.specialty && (
                  <div>
                    <p className="text-sm text-gray-500">专长</p>
                    <p className="font-semibold">{cert.specialty}</p>
                  </div>
                )}
              </div>

              <div className="mb-4">
                <p className="text-sm text-gray-500">烘焙经验</p>
                <p className="text-sm whitespace-pre-wrap">{cert.experience}</p>
              </div>

              {cert.review_comment && (
                <div className="mb-4 p-3 bg-gray-50 rounded">
                  <p className="text-sm text-gray-500">审核意见</p>
                  <p className="text-sm">{cert.review_comment}</p>
                </div>
              )}

              {cert.status === 'pending' && (
                <div className="flex gap-2">
                  <button
                    onClick={() => handleReview(cert.id, 'approved')}
                    className="btn btn-primary"
                  >
                    通过
                  </button>
                  <button
                    onClick={() => {
                      const comment = prompt('请输入拒绝原因')
                      if (comment !== null) {
                        handleReview(cert.id, 'rejected', comment)
                      }
                    }}
                    className="btn btn-danger"
                  >
                    拒绝
                  </button>
                </div>
              )}
            </div>
          ))
        )}
      </div>

      {total > 10 && (
        <div className="flex justify-center items-center gap-2">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="btn btn-secondary disabled:opacity-50"
          >
            上一页
          </button>
          <span className="text-gray-600">{page} / {Math.ceil(total / 10)}</span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={page >= Math.ceil(total / 10)}
            className="btn btn-secondary disabled:opacity-50"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
