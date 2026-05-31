import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { roastingApi } from '@/api/roasting'
import { RoastingRecord, RoastingDataPoint } from '@/types'
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'

export default function RoastingDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [record, setRecord] = useState<RoastingRecord | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      loadRecord(Number(id))
    }
  }, [id])

  const loadRecord = async (recordId: number) => {
    try {
      const res = await roastingApi.get(recordId)
      setRecord(res.data || null)
    } catch {
      setRecord(null)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (!record) {
    return <div className="text-center py-12">烘焙记录不存在</div>
  }

  const chartData = record.data_points?.map((dp: RoastingDataPoint) => ({
    time: dp.time_elapsed,
    beanTemp: dp.bean_temp,
    envTemp: dp.env_temp,
  })) || []

  return (
    <div className="space-y-6">
      <button onClick={() => navigate(-1)} className="text-coffee-600 hover:underline">
        ← 返回
      </button>

      <div className="card p-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-2xl font-bold text-gray-800">{record.product?.name || `商品 #${record.product_id}`}</h1>
            <p className="text-gray-500">批次号: {record.batch_number}</p>
          </div>
          <div className="text-right">
            <p className="text-coffee-600 font-bold text-2xl">{record.total_roast_time}s</p>
            <p className="text-sm text-gray-500">{record.roasted_at?.split('T')[0]}</p>
          </div>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <div className="card p-4">
            <p className="text-sm text-gray-500">入豆温</p>
            <p className="text-xl font-bold">{record.input_temp}°C</p>
          </div>
          <div className="card p-4">
            <p className="text-sm text-gray-500">回温点</p>
            <p className="text-xl font-bold">{record.turning_point}°C</p>
            <p className="text-sm text-gray-400">{record.turning_time}s</p>
          </div>
          <div className="card p-4">
            <p className="text-sm text-gray-500">一爆</p>
            <p className="text-xl font-bold">{record.first_crack_temp}°C</p>
            <p className="text-sm text-gray-400">{record.first_crack_time}s</p>
          </div>
          <div className="card p-4">
            <p className="text-sm text-gray-500">出豆温</p>
            <p className="text-xl font-bold">{record.drop_temp}°C</p>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="card p-4">
            <p className="text-sm text-gray-500">生豆重量</p>
            <p className="text-xl font-bold">{record.green_bean_weight}g</p>
          </div>
          <div className="card p-4">
            <p className="text-sm text-gray-500">熟豆重量</p>
            <p className="text-xl font-bold">{record.roasted_weight}g</p>
            <p className="text-sm text-gray-400">
              失重率: {record.green_bean_weight > 0
                ? (((record.green_bean_weight - record.roasted_weight) / record.green_bean_weight) * 100).toFixed(1)
                : 0}%
            </p>
          </div>
        </div>

        {record.second_crack_temp > 0 && (
          <div className="card p-4 mb-6">
            <p className="text-sm text-gray-500">二爆</p>
            <p className="text-xl font-bold">{record.second_crack_temp}°C</p>
            <p className="text-sm text-gray-400">{record.second_crack_time}s</p>
          </div>
        )}

        {record.notes && (
          <div className="card p-4">
            <p className="text-sm text-gray-500 mb-2">备注</p>
            <p className="text-gray-700 whitespace-pre-wrap">{record.notes}</p>
          </div>
        )}
      </div>

      {chartData.length > 0 && (
        <div className="card p-6">
          <h2 className="text-xl font-bold text-gray-800 mb-4">烘焙曲线</h2>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="time" label={{ value: '时间(s)', position: 'bottom' }} />
                <YAxis label={{ value: '温度(°C)', angle: -90, position: 'insideLeft' }} />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="beanTemp" name="豆温" stroke="#c07a2a" strokeWidth={2} dot={false} />
                <Line type="monotone" dataKey="envTemp" name="环境温度" stroke="#6b341c" strokeWidth={2} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}
    </div>
  )
}
