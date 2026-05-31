import { useEffect, useState } from 'react'
import { statsApi } from '@/api/search'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, BarChart, Bar, PieChart, Pie, Cell } from 'recharts'

export default function AdminStats() {
  const [salesTrend, setSalesTrend] = useState<any[]>([])
  const [originDistribution, setOriginDistribution] = useState<any[]>([])
  const [userActivity, setUserActivity] = useState<any[]>([])
  const [topProducts, setTopProducts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const [trend, origins, activity, top] = await Promise.all([
        statsApi.getSalesTrend({ days: 30 }),
        statsApi.getOriginDistribution(),
        statsApi.getUserActivity({ days: 30 }),
        statsApi.getTopProducts({ limit: 10 }),
      ])
      setSalesTrend(trend.data || [])
      setOriginDistribution(origins.data || [])
      setUserActivity(activity.data || [])
      setTopProducts(top.data || [])
    } catch {
      setSalesTrend([])
      setOriginDistribution([])
      setUserActivity([])
      setTopProducts([])
    } finally {
      setLoading(false)
    }
  }

  const handleExport = async (type: string) => {
    try {
      const res = await statsApi.exportExcel(type)
      const blob = new Blob([res.data as any])
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `${type}_report_${Date.now()}.xlsx`
      a.click()
      URL.revokeObjectURL(url)
    } catch (err: any) {
      alert(err.message || '导出失败')
    }
  }

  const COLORS = ['#c07a2a', '#a85f22', '#8b461f', '#6b341c', '#4a2316', '#deb573', '#cf9647']

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">数据统计</h1>
        <div className="flex gap-2">
          <button onClick={() => handleExport('sales')} className="btn btn-outline text-sm">
            导出销售报表
          </button>
          <button onClick={() => handleExport('orders')} className="btn btn-outline text-sm">
            导出订单列表
          </button>
          <button onClick={() => handleExport('products')} className="btn btn-outline text-sm">
            导出商品列表
          </button>
          <button onClick={() => handleExport('users')} className="btn btn-outline text-sm">
            导出用户列表
          </button>
        </div>
      </div>

      <div className="card p-6">
        <h2 className="text-lg font-bold mb-4">销售趋势 (近30天)</h2>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={salesTrend}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis yAxisId="left" />
              <YAxis yAxisId="right" orientation="right" />
              <Tooltip />
              <Legend />
              <Line yAxisId="left" type="monotone" dataKey="order_count" name="订单数" stroke="#c07a2a" strokeWidth={2} />
              <Line yAxisId="right" type="monotone" dataKey="amount" name="销售额" stroke="#6b341c" strokeWidth={2} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card p-6">
          <h2 className="text-lg font-bold mb-4">产地分布</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={originDistribution}
                  dataKey="count"
                  nameKey="origin"
                  cx="50%"
                  cy="50%"
                  outerRadius={80}
                  label={(entry: any) => entry.origin}
                >
                  {originDistribution.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="card p-6">
          <h2 className="text-lg font-bold mb-4">用户活跃度 (近30天)</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={userActivity}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="user_count" name="活跃用户" fill="#c07a2a" />
                <Bar dataKey="login_count" name="操作数" fill="#6b341c" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      <div className="card p-6">
        <h2 className="text-lg font-bold mb-4">热销商品 Top 10</h2>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={topProducts} layout="vertical">
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis type="number" />
              <YAxis type="category" dataKey="product_name" width={150} />
              <Tooltip />
              <Legend />
              <Bar dataKey="total_sold" name="销量" fill="#c07a2a" />
              <Bar dataKey="total_amount" name="销售额" fill="#6b341c" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  )
}
