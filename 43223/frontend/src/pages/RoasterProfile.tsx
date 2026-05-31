import { useEffect, useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { certificationApi } from '@/api/certification'
import { User, RoasterCertification, Product, RoastingRecord } from '@/types'

interface ProfileData {
  user: User
  certification?: RoasterCertification
  products: Product[]
  roasting_records: RoastingRecord[]
  total_products: number
  total_roasts: number
  avg_score: number
}

export default function RoasterProfile() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [profile, setProfile] = useState<ProfileData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      loadProfile(Number(id))
    }
  }, [id])

  const loadProfile = async (roasterId: number) => {
    try {
      const res = await certificationApi.getRoasterProfile(roasterId)
      setProfile(res.data || null)
    } catch {
      setProfile(null)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (!profile) {
    return <div className="text-center py-12">用户不存在或非认证烘焙师</div>
  }

  return (
    <div className="space-y-8">
      <button onClick={() => navigate(-1)} className="text-coffee-600 hover:underline">
        ← 返回
      </button>

      <div className="card p-8">
        <div className="flex items-center gap-6">
          <div className="w-24 h-24 bg-coffee-200 rounded-full flex items-center justify-center text-4xl">
            {profile.user.avatar ? (
              <img src={profile.user.avatar} alt="" className="w-full h-full rounded-full object-cover" />
            ) : (
              profile.user.nickname?.[0] || profile.user.username[0]
            )}
          </div>
          <div className="flex-1">
            <h1 className="text-3xl font-bold text-gray-800">
              {profile.user.nickname || profile.user.username}
              {profile.user.is_certified && (
                <span className="ml-2 badge bg-coffee-100 text-coffee-700">认证烘焙师</span>
              )}
            </h1>
            {profile.user.bio && (
              <p className="text-gray-600 mt-2">{profile.user.bio}</p>
            )}
            <div className="flex gap-6 mt-4">
              <div>
                <p className="text-2xl font-bold text-coffee-600">{profile.total_products}</p>
                <p className="text-sm text-gray-500">在售商品</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-coffee-600">{profile.total_roasts}</p>
                <p className="text-sm text-gray-500">烘焙记录</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-coffee-600">{profile.avg_score?.toFixed(1) || 'N/A'}</p>
                <p className="text-sm text-gray-500">平均评分</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {profile.certification && (
        <div className="card p-6">
          <h2 className="text-xl font-bold text-gray-800 mb-4">认证信息</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-sm text-gray-500">证书名称</p>
              <p className="font-semibold">{profile.certification.cert_name}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">证书编号</p>
              <p className="font-semibold">{profile.certification.cert_number}</p>
            </div>
            <div>
              <p className="text-sm text-gray-500">发证机构</p>
              <p className="font-semibold">{profile.certification.org_name}</p>
            </div>
            {profile.certification.specialty && (
              <div>
                <p className="text-sm text-gray-500">专长</p>
                <p className="font-semibold">{profile.certification.specialty}</p>
              </div>
            )}
          </div>
        </div>
      )}

      <div>
        <h2 className="text-xl font-bold text-gray-800 mb-4">在售商品</h2>
        {profile.products.length === 0 ? (
          <p className="text-gray-500">暂无在售商品</p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {profile.products.map((product) => (
              <Link
                key={product.id}
                to={`/products/${product.id}`}
                className="card overflow-hidden hover:shadow-md transition-shadow"
              >
                <div className="aspect-square bg-gray-100">
                  {product.images?.[0] ? (
                    <img src={product.images[0].url} alt={product.name} className="w-full h-full object-cover" />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-coffee-300 text-4xl">☕</div>
                  )}
                </div>
                <div className="p-3">
                  <h3 className="font-semibold text-sm truncate">{product.name}</h3>
                  <p className="text-coffee-600 font-bold">¥{product.price}</p>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>

      <div>
        <h2 className="text-xl font-bold text-gray-800 mb-4">最近烘焙记录</h2>
        {profile.roasting_records.length === 0 ? (
          <p className="text-gray-500">暂无烘焙记录</p>
        ) : (
          <div className="space-y-2">
            {profile.roasting_records.map((record) => (
              <Link
                key={record.id}
                to={`/roasting/${record.id}`}
                className="card p-4 hover:shadow-md transition-shadow flex items-center justify-between"
              >
                <div>
                  <p className="font-semibold">{record.product?.name || `批次 #${record.batch_number}`}</p>
                  <p className="text-sm text-gray-500">{record.roasted_at?.split('T')[0]}</p>
                </div>
                <div className="text-right">
                  <p className="font-bold text-coffee-600">{record.total_roast_time}s</p>
                  <p className="text-sm text-gray-500">{record.drop_temp}°C</p>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
