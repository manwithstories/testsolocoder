import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { roastingApi } from '@/api/roasting'
import { productApi } from '@/api/product'
import { useAuthStore } from '@/store/auth'

export default function RoasterDashboard() {
  const { user } = useAuthStore()
  const [stats, setStats] = useState<any>(null)
  const [products, setProducts] = useState<any[]>([])
  const [recentRecords, setRecentRecords] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [user?.id])

  const loadData = async () => {
    if (!user?.id) return
    setLoading(true)
    try {
      const [statsRes, productsRes, recordsRes] = await Promise.all([
        roastingApi.getStats({ roaster_id: String(user.id) }),
        productApi.list({ roaster_id: String(user.id), page_size: 5 }),
        roastingApi.list({ roaster_id: String(user.id), page_size: 5 }),
      ])
      setStats(statsRes.data)
      setProducts(productsRes.data?.items || [])
      setRecentRecords(recordsRes.data?.items || [])
    } catch {
      setStats(null)
      setProducts([])
      setRecentRecords([])
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">烘焙师工作台</h1>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="card p-4">
          <p className="text-sm text-gray-500">烘焙次数</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.total_count || 0}</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">平均入豆温</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.avg_input_temp?.toFixed(1) || 0}°C</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">平均出豆温</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.avg_drop_temp?.toFixed(1) || 0}°C</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">平均烘焙时间</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.avg_total_time?.toFixed(0) || 0}s</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-bold">我的商品</h2>
            <Link to="/roaster/products" className="text-coffee-600 hover:underline text-sm">管理全部 →</Link>
          </div>
          {products.length === 0 ? (
            <p className="text-gray-500 text-sm">暂无商品</p>
          ) : (
            <div className="space-y-2">
              {products.map((product) => (
                <div key={product.id} className="flex items-center justify-between py-2 border-b last:border-0">
                  <span className="text-sm">{product.name}</span>
                  <div className="flex items-center gap-2 text-sm">
                    <span className="text-coffee-600">¥{product.price}</span>
                    <span className="text-gray-400">库存: {product.stock}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="card p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-bold">最近烘焙记录</h2>
            <Link to="/roaster/roasting" className="text-coffee-600 hover:underline text-sm">查看全部 →</Link>
          </div>
          {recentRecords.length === 0 ? (
            <p className="text-gray-500 text-sm">暂无烘焙记录</p>
          ) : (
            <div className="space-y-2">
              {recentRecords.map((record) => (
                <div key={record.id} className="flex items-center justify-between py-2 border-b last:border-0">
                  <div>
                    <p className="text-sm font-medium">{record.batch_number}</p>
                    <p className="text-xs text-gray-500">{record.roasted_at?.split('T')[0]}</p>
                  </div>
                  <div className="text-right">
                    <p className="text-sm text-coffee-600">{record.total_roast_time}s</p>
                    <p className="text-xs text-gray-500">{record.drop_temp}°C</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Link to="/roaster/products" className="card p-6 hover:shadow-md transition-shadow">
          <h3 className="font-bold text-lg mb-2">商品管理</h3>
          <p className="text-gray-500 text-sm">添加、编辑和管理咖啡豆商品</p>
        </Link>
        <Link to="/roaster/roasting" className="card p-6 hover:shadow-md transition-shadow">
          <h3 className="font-bold text-lg mb-2">烘焙记录</h3>
          <p className="text-gray-500 text-sm">记录和管理烘焙参数曲线</p>
        </Link>
      </div>
    </div>
  )
}
