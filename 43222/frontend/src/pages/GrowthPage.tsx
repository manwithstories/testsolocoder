import { useQuery } from '@tanstack/react-query'
import { Camera, Plus, Leaf } from 'lucide-react'
import { Link } from 'react-router-dom'
import { plantingRecordApi } from '@/api'
import { useState } from 'react'

export default function GrowthPage() {
  const [filterStatus, setFilterStatus] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['planting-records', filterStatus],
    queryFn: () =>
      plantingRecordApi.getAll({
        status: filterStatus || undefined,
        page_size: 50,
      }),
  })

  const records = data?.data?.records || []

  const statusLabels: Record<string, { label: string; color: string }> = {
    planted: { label: '已种植', color: 'bg-gray-100 text-gray-700' },
    growing: { label: '生长中', color: 'bg-blue-100 text-blue-700' },
    flowering: { label: '开花中', color: 'bg-pink-100 text-pink-700' },
    fruiting: { label: '结果中', color: 'bg-amber-100 text-amber-700' },
    harvested: { label: '已收获', color: 'bg-green-100 text-green-700' },
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">生长追踪</h1>
          <p className="text-gray-500">记录和追踪植物生长过程</p>
        </div>
        <Link to="/plots" className="btn-primary">
          <Plus className="w-4 h-4 mr-2" />
          新建种植
        </Link>
      </div>

      {/* Filter */}
      <div className="flex gap-2">
        {['', 'planted', 'growing', 'flowering', 'fruiting', 'harvested'].map((status) => (
          <button
            key={status}
            onClick={() => setFilterStatus(status)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              filterStatus === status
                ? 'bg-garden-600 text-white'
                : 'bg-white border border-gray-200 text-gray-600 hover:bg-gray-50'
            }`}
          >
            {status === '' ? '全部' : statusLabels[status]?.label}
          </button>
        ))}
      </div>

      {/* Records */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : records.length === 0 ? (
        <div className="card text-center py-12">
          <Leaf className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500 mb-4">还没有种植记录</p>
          <Link to="/plots" className="btn-primary">
            开始种植
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {records.map((record: any) => (
            <div key={record.id} className="card">
              <div className="card-body">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-center gap-4">
                    <div className="w-14 h-14 bg-garden-100 rounded-xl flex items-center justify-center">
                      {record.plant?.image_url ? (
                        <img
                          src={record.plant.image_url}
                          alt={record.plant?.name}
                          className="w-full h-full object-cover rounded-xl"
                        />
                      ) : (
                        <Leaf className="w-7 h-7 text-garden-600" />
                      )}
                    </div>
                    <div>
                      <h3 className="font-semibold text-gray-900">{record.plant?.name}</h3>
                      <p className="text-sm text-gray-500">
                        {record.plot?.name} · 种植于 {new Date(record.planting_date).toLocaleDateString('zh-CN')}
                      </p>
                    </div>
                  </div>
                  <span className={`badge ${statusLabels[record.status]?.color || 'bg-gray-100 text-gray-700'}`}>
                    {statusLabels[record.status]?.label || record.status}
                  </span>
                </div>

                {/* Growth Timeline */}
                {record.growth_logs && record.growth_logs.length > 0 && (
                  <div className="mt-4 pt-4 border-t border-gray-100">
                    <p className="text-sm font-medium text-gray-700 mb-3">生长日志</p>
                    <div className="flex gap-2 overflow-x-auto pb-2">
                      {record.growth_logs.map((log: any) => (
                        <div key={log.id} className="flex-shrink-0 w-24">
                          <div className="aspect-square bg-gray-100 rounded-lg overflow-hidden mb-1">
                            {log.image_url ? (
                              <img
                                src={log.image_url}
                                alt={log.title}
                                className="w-full h-full object-cover"
                              />
                            ) : (
                              <div className="w-full h-full flex items-center justify-center">
                                <Camera className="w-6 h-6 text-gray-300" />
                              </div>
                            )}
                          </div>
                          <p className="text-xs text-gray-600 truncate">{log.title}</p>
                          <p className="text-xs text-gray-400">
                            {new Date(log.log_date).toLocaleDateString('zh-CN')}
                          </p>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Stats */}
                <div className="mt-4 pt-4 border-t border-gray-100 grid grid-cols-3 gap-4 text-center">
                  <div>
                    <p className="text-lg font-semibold text-gray-900">{record.quantity}</p>
                    <p className="text-xs text-gray-500">种植数量</p>
                  </div>
                  <div>
                    <p className="text-lg font-semibold text-gray-900">
                      {record.expected_harvest_date
                        ? new Date(record.expected_harvest_date).toLocaleDateString('zh-CN')
                        : '-'}
                    </p>
                    <p className="text-xs text-gray-500">预计收获</p>
                  </div>
                  <div>
                    <p className="text-lg font-semibold text-gray-900">{record.growth_logs?.length || 0}</p>
                    <p className="text-xs text-gray-500">日志数</p>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
