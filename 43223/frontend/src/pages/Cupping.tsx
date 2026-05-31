import { useEffect, useState } from 'react'
import { cuppingApi } from '@/api/cupping'
import { CuppingScore, Product } from '@/types'
import { useAuthStore } from '@/store/auth'

export default function Cupping() {
  const { user } = useAuthStore()
  const [scores, setScores] = useState<CuppingScore[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({
    product_id: 0,
    dry_fragrance: 7,
    wet_aroma: 7,
    body: 7,
    acidity: 7,
    sweetness: 7,
    aftertaste: 7,
    balance: 7,
    notes: '',
  })

  useEffect(() => {
    loadMyScores()
  }, [])

  const loadMyScores = async () => {
    try {
      const res = await cuppingApi.getMyScores()
      setScores(res.data || [])
    } catch {
      setScores([])
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (formData.product_id === 0) {
      alert('请选择商品')
      return
    }
    try {
      await cuppingApi.create(formData)
      alert('评分成功')
      setShowForm(false)
      loadMyScores()
    } catch (err: any) {
      alert(err.message || '评分失败')
    }
  }

  const calculateOverall = () => {
    return ((formData.dry_fragrance + formData.wet_aroma + formData.body +
      formData.acidity + formData.sweetness + formData.aftertaste + formData.balance) / 7).toFixed(2)
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">我的杯测评分</h1>
        <button onClick={() => setShowForm(!showForm)} className="btn btn-primary">
          {showForm ? '取消' : '+ 新建评分'}
        </button>
      </div>

      {showForm && (
        <div className="card p-6">
          <h2 className="text-xl font-bold mb-4">新建杯测评分</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
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
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[
                { key: 'dry_fragrance', label: '干香' },
                { key: 'wet_aroma', label: '湿香' },
                { key: 'body', label: '醇厚' },
                { key: 'acidity', label: '酸质' },
                { key: 'sweetness', label: '甜度' },
                { key: 'aftertaste', label: '余韵' },
                { key: 'balance', label: '平衡' },
              ].map((item) => (
                <div key={item.key}>
                  <label className="label">{item.label}</label>
                  <input
                    type="number"
                    min="0"
                    max="10"
                    step="0.5"
                    value={(formData as any)[item.key]}
                    onChange={(e) => setFormData({ ...formData, [item.key]: Number(e.target.value) })}
                    className="input"
                  />
                </div>
              ))}
            </div>
            <div className="card p-4 bg-coffee-50">
              <p className="text-sm text-gray-500">综合评分</p>
              <p className="text-3xl font-bold text-coffee-600">{calculateOverall()}</p>
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
            <button type="submit" className="btn btn-primary w-full">提交评分</button>
          </form>
        </div>
      )}

      {scores.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">📝</div>
          <p className="text-gray-500">暂无评分记录</p>
        </div>
      ) : (
        <div className="space-y-4">
          {scores.map((score) => (
            <div key={score.id} className="card p-4">
              <div className="flex items-center justify-between mb-2">
                <span className="font-semibold">
                  {score.product?.name || `商品 #${score.product_id}`}
                </span>
                <span className="text-coffee-600 font-bold text-xl">{score.overall_score.toFixed(2)}分</span>
              </div>
              <div className="grid grid-cols-7 gap-2 text-sm">
                <div><span className="text-gray-500">干香</span> {score.dry_fragrance}</div>
                <div><span className="text-gray-500">湿香</span> {score.wet_aroma}</div>
                <div><span className="text-gray-500">醇厚</span> {score.body}</div>
                <div><span className="text-gray-500">酸质</span> {score.acidity}</div>
                <div><span className="text-gray-500">甜度</span> {score.sweetness}</div>
                <div><span className="text-gray-500">余韵</span> {score.aftertaste}</div>
                <div><span className="text-gray-500">平衡</span> {score.balance}</div>
              </div>
              {score.notes && <p className="text-gray-600 mt-2 text-sm">{score.notes}</p>}
              <p className="text-xs text-gray-400 mt-2">{score.created_at?.split('T')[0]}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
