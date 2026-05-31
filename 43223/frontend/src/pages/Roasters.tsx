import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { certificationApi } from '@/api/certification'
import { User } from '@/types'

export default function Roasters() {
  const [roasters, setRoasters] = useState<User[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadRoasters()
  }, [])

  const loadRoasters = async () => {
    try {
      const res = await certificationApi.listCertifiedRoasters({ page: 1, page_size: 50 })
      setRoasters(res.data?.items || [])
    } catch {
      setRoasters([])
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">认证烘焙师</h1>

      {roasters.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">👨‍🍳</div>
          <p className="text-gray-500">暂无认证烘焙师</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {roasters.map((roaster) => (
            <Link
              key={roaster.id}
              to={`/roasters/${roaster.id}`}
              className="card p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex items-center gap-4 mb-4">
                <div className="w-16 h-16 bg-coffee-200 rounded-full flex items-center justify-center text-2xl">
                  {roaster.avatar ? (
                    <img src={roaster.avatar} alt="" className="w-full h-full rounded-full object-cover" />
                  ) : (
                    roaster.nickname?.[0] || roaster.username[0]
                  )}
                </div>
                <div>
                  <h3 className="font-bold text-gray-800">{roaster.nickname || roaster.username}</h3>
                  <p className="text-sm text-gray-500">@{roaster.username}</p>
                </div>
              </div>
              {roaster.bio && (
                <p className="text-sm text-gray-600 mb-4 line-clamp-2">{roaster.bio}</p>
              )}
              {roaster.certification && (
                <div className="text-sm">
                  <span className="text-coffee-600">📜 {roaster.certification.cert_name}</span>
                </div>
              )}
            </Link>
          ))}
        </div>
      )}
    </div>
  )
}
