import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { roastingApi } from '@/api/roasting'
import { RoastingRecord } from '@/types'

export default function RoastingRecords() {
  const [records, setRecords] = useState<RoastingRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)

  useEffect(() => {
    loadRecords()
  }, [page])

  const loadRecords = async () => {
    setLoading(true)
    try {
      const res = await roastingApi.list({ page, page_size: 10 })
      setRecords(res.data?.items || [])
      setTotal(res.data?.total || 0)
    } catch {
      setRecords([])
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">烘焙记录</h1>

      {records.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">🔥</div>
          <p className="text-gray-500">暂无烘焙记录</p>
        </div>
      ) : (
        <div className="space-y-4">
          {records.map((record) => (
            <Link
              key={record.id}
              to={`/roasting/${record.id}`}
              className="card p-4 hover:shadow-md transition-shadow block"
            >
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="font-semibold text-gray-800">{record.product?.name || `商品 #${record.product_id}`}</h3>
                  <p className="text-sm text-gray-500">批次号: {record.batch_number}</p>
                  <p className="text-sm text-gray-500">
                    {record.roaster?.nickname || record.roaster?.username || '未知烘焙师'}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-coffee-600 font-bold">{record.total_roast_time}s</p>
                  <p className="text-sm text-gray-500">{record.roasted_at?.split('T')[0]}</p>
                </div>
              </div>
              <div className="grid grid-cols-4 gap-4 mt-4 text-sm">
                <div>
                  <span className="text-gray-500">入豆温</span>
                  <p className="font-semibold">{record.input_temp}°C</p>
                </div>
                <div>
                  <span className="text-gray-500">一爆</span>
                  <p className="font-semibold">{record.first_crack_temp}°C / {record.first_crack_time}s</p>
                </div>
                <div>
                  <span className="text-gray-500">出豆温</span>
                  <p className="font-semibold">{record.drop_temp}°C</p>
                </div>
                <div>
                  <span className="text-gray-500">重量</span>
                  <p className="font-semibold">{record.green_bean_weight}g → {record.roasted_weight}g</p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

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
