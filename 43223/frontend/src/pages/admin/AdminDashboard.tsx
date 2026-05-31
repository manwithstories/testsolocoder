import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { statsApi } from '@/api/search'

export default function AdminDashboard() {
  const [stats, setStats] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const res = await statsApi.getSalesStats()
      setStats(res.data)
    } catch {
      setStats(null)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">管理后台</h1>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="card p-4">
          <p className="text-sm text-gray-500">总订单数</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.total_orders || 0}</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">总销售额</p>
          <p className="text-2xl font-bold text-coffee-600">¥{(stats?.total_amount || 0).toFixed(2)}</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">商品总数</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.total_products || 0}</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">用户总数</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.total_users || 0}</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">认证烘焙师</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.total_roasters || 0}</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">今日订单</p>
          <p className="text-2xl font-bold text-coffee-600">{stats?.today_orders || 0}</p>
        </div>
        <div className="card p-4">
          <p className="text-sm text-gray-500">今日销售额</p>
          <p className="text-2xl font-bold text-coffee-600">¥{(stats?.today_amount || 0).toFixed(2)}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Link to="/admin/users" className="card p-6 hover:shadow-md transition-shadow">
          <h3 className="font-bold text-lg mb-2">用户管理</h3>
          <p className="text-gray-500">管理用户账户、角色和权限</p>
        </Link>
        <Link to="/admin/products" className="card p-6 hover:shadow-md transition-shadow">
          <h3 className="font-bold text-lg mb-2">商品管理</h3>
          <p className="text-gray-500">管理咖啡豆商品和上下架</p>
        </Link>
        <Link to="/admin/certifications" className="card p-6 hover:shadow-md transition-shadow">
          <h3 className="font-bold text-lg mb-2">认证审核</h3>
          <p className="text-gray-500">审核烘焙师认证申请</p>
        </Link>
        <Link to="/admin/stats" className="card p-6 hover:shadow-md transition-shadow">
          <h3 className="font-bold text-lg mb-2">数据统计</h3>
          <p className="text-gray-500">查看销售统计和报表导出</p>
        </Link>
      </div>
    </div>
  )
}
