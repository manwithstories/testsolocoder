import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { roastingApi } from '@/api/roasting'
import { productApi } from '@/api/product'
import { RoastingRecord, CreateRoastingRecordRequest } from '@/types'
import { useAuthStore } from '@/store/auth'

export default function RoasterRoastingRecords() {
  const { user } = useAuthStore()
  const [records, setRecords] = useState<RoastingRecord[]>([])
  const [products, setProducts] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingRecord, setEditingRecord] = useState<RoastingRecord | null>(null)
  const [formData, setFormData] = useState<CreateRoastingRecordRequest>({
    product_id: 0,
    batch_number: '',
    green_bean_weight: 0,
    roasted_weight: 0,
    input_temp: 0,
    turning_point: 0,
    turning_time: 0,
    first_crack_temp: 0,
    first_crack_time: 0,
    second_crack_temp: 0,
    second_crack_time: 0,
    drop_temp: 0,
    total_roast_time: 0,
    notes: '',
    data_points: [],
    roasted_at: new Date().toISOString().split('T')[0],
  })

  useEffect(() => {
    loadData()
  }, [user?.id])

  const loadData = async () => {
    if (!user?.id) return
    setLoading(true)
    try {
      const [recordsRes, productsRes] = await Promise.all([
        roastingApi.list({ roaster_id: String(user.id), page_size: 50 }),
        productApi.list({ roaster_id: String(user.id), page_size: 50 }),
      ])
      setRecords(recordsRes.data?.items || [])
      setProducts(productsRes.data?.items || [])
    } catch {
      setRecords([])
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingRecord) {
        await roastingApi.update(editingRecord.id, formData)
        alert('更新成功')
      } else {
        await roastingApi.create(formData)
        alert('创建成功')
      }
      setShowForm(false)
      setEditingRecord(null)
      loadData()
    } catch (err: any) {
      alert(err.message || '保存失败')
    }
  }

  const handleEdit = (record: RoastingRecord) => {
    setEditingRecord(record)
    setFormData({
      product_id: record.product_id,
      batch_number: record.batch_number,
      green_bean_weight: record.green_bean_weight,
      roasted_weight: record.roasted_weight,
      input_temp: record.input_temp,
      turning_point: record.turning_point,
      turning_time: record.turning_time,
      first_crack_temp: record.first_crack_temp,
      first_crack_time: record.first_crack_time,
      second_crack_temp: record.second_crack_temp,
      second_crack_time: record.second_crack_time,
      drop_temp: record.drop_temp,
      total_roast_time: record.total_roast_time,
      notes: record.notes || '',
      data_points: record.data_points?.map(dp => ({
        time_elapsed: dp.time_elapsed,
        bean_temp: dp.bean_temp,
        env_temp: dp.env_temp,
        rate_of_rise: dp.rate_of_rise,
      })) || [],
      roasted_at: record.roasted_at.split('T')[0],
    })
    setShowForm(true)
  }

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除该记录吗？')) return
    try {
      await roastingApi.delete(id)
      alert('删除成功')
      loadData()
    } catch (err: any) {
      alert(err.message || '删除失败')
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (showForm) {
    return (
      <div className="max-w-2xl mx-auto space-y-6">
        <h1 className="text-2xl font-bold text-gray-800">
          {editingRecord ? '编辑烘焙记录' : '新建烘焙记录'}
        </h1>
        <div className="card p-6">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="label">商品 *</label>
                <select
                  value={formData.product_id}
                  onChange={(e) => setFormData({ ...formData, product_id: Number(e.target.value) })}
                  className="input"
                  required
                >
                  <option value={0}>请选择商品</option>
                  {products.map((p) => (
                    <option key={p.id} value={p.id}>{p.name}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="label">批次号 *</label>
                <input
                  type="text"
                  value={formData.batch_number}
                  onChange={(e) => setFormData({ ...formData, batch_number: e.target.value })}
                  className="input"
                  required
                />
              </div>
              <div>
                <label className="label">烘焙日期</label>
                <input
                  type="date"
                  value={formData.roasted_at.split('T')[0]}
                  onChange={(e) => setFormData({ ...formData, roasted_at: e.target.value })}
                  className="input"
                />
              </div>
              <div>
                <label className="label">总烘焙时间 (秒)</label>
                <input
                  type="number"
                  value={formData.total_roast_time}
                  onChange={(e) => setFormData({ ...formData, total_roast_time: Number(e.target.value) })}
                  className="input"
                />
              </div>
            </div>

            <div className="border-t pt-4">
              <h3 className="font-bold mb-4">重量</h3>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">生豆重量 (g)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={formData.green_bean_weight}
                    onChange={(e) => setFormData({ ...formData, green_bean_weight: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">熟豆重量 (g)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={formData.roasted_weight}
                    onChange={(e) => setFormData({ ...formData, roasted_weight: Number(e.target.value) })}
                    className="input"
                  />
                </div>
              </div>
            </div>

            <div className="border-t pt-4">
              <h3 className="font-bold mb-4">温度参数</h3>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="label">入豆温 (°C)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={formData.input_temp}
                    onChange={(e) => setFormData({ ...formData, input_temp: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">回温点 (°C)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={formData.turning_point}
                    onChange={(e) => setFormData({ ...formData, turning_point: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">回温时间 (秒)</label>
                  <input
                    type="number"
                    value={formData.turning_time}
                    onChange={(e) => setFormData({ ...formData, turning_time: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">一爆温度 (°C)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={formData.first_crack_temp}
                    onChange={(e) => setFormData({ ...formData, first_crack_temp: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">一爆时间 (秒)</label>
                  <input
                    type="number"
                    value={formData.first_crack_time}
                    onChange={(e) => setFormData({ ...formData, first_crack_time: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">二爆温度 (°C)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={formData.second_crack_temp}
                    onChange={(e) => setFormData({ ...formData, second_crack_temp: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">二爆时间 (秒)</label>
                  <input
                    type="number"
                    value={formData.second_crack_time}
                    onChange={(e) => setFormData({ ...formData, second_crack_time: Number(e.target.value) })}
                    className="input"
                  />
                </div>
                <div>
                  <label className="label">出豆温 (°C)</label>
                  <input
                    type="number"
                    step="0.1"
                    value={formData.drop_temp}
                    onChange={(e) => setFormData({ ...formData, drop_temp: Number(e.target.value) })}
                    className="input"
                  />
                </div>
              </div>
            </div>

            <div>
              <label className="label">备注</label>
              <textarea
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                className="input"
                rows={3}
              />
            </div>

            <div className="flex gap-4">
              <button type="button" onClick={() => { setShowForm(false); setEditingRecord(null) }} className="btn btn-secondary flex-1">
                取消
              </button>
              <button type="submit" className="btn btn-primary flex-1">
                {editingRecord ? '更新' : '创建'}
              </button>
            </div>
          </form>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">烘焙记录</h1>
        <button onClick={() => { setEditingRecord(null); setShowForm(true) }} className="btn btn-primary">
          + 新建记录
        </button>
      </div>

      {records.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">🔥</div>
          <p className="text-gray-500 mb-4">暂无烘焙记录</p>
          <button onClick={() => setShowForm(true)} className="btn btn-primary">
            创建第一条记录
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {records.map((record) => (
            <div key={record.id} className="card p-4">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h3 className="font-bold">{record.product?.name || `商品 #${record.product_id}`}</h3>
                  <p className="text-sm text-gray-500">批次号: {record.batch_number} | {record.roasted_at.split('T')[0]}</p>
                </div>
                <div className="flex items-center gap-4">
                  <span className="text-coffee-600 font-bold">{record.total_roast_time}s</span>
                  <Link to={`/roasting/${record.id}`} className="text-coffee-600 hover:underline text-sm">
                    查看详情
                  </Link>
                  <button onClick={() => handleEdit(record)} className="text-coffee-600 hover:underline text-sm">
                    编辑
                  </button>
                  <button onClick={() => handleDelete(record.id)} className="text-red-500 hover:underline text-sm">
                    删除
                  </button>
                </div>
              </div>
              <div className="grid grid-cols-4 gap-4 text-sm">
                <div><span className="text-gray-500">入豆温</span> {record.input_temp}°C</div>
                <div><span className="text-gray-500">回温</span> {record.turning_point}°C / {record.turning_time}s</div>
                <div><span className="text-gray-500">一爆</span> {record.first_crack_temp}°C / {record.first_crack_time}s</div>
                <div><span className="text-gray-500">出豆</span> {record.drop_temp}°C</div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
