import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { productApi } from '@/api/product'
import { cartApi } from '@/api/order'
import { cuppingApi } from '@/api/cupping'
import { Product, CuppingScore } from '@/types'
import { useAuthStore } from '@/store/auth'

export default function ProductDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { isAuthenticated } = useAuthStore()
  const [product, setProduct] = useState<Product | null>(null)
  const [scores, setScores] = useState<CuppingScore[]>([])
  const [quantity, setQuantity] = useState(1)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      loadProduct(Number(id))
      loadScores(Number(id))
    }
  }, [id])

  const loadProduct = async (productId: number) => {
    try {
      const res = await productApi.get(productId)
      setProduct(res.data || null)
    } catch {
      setProduct(null)
    } finally {
      setLoading(false)
    }
  }

  const loadScores = async (productId: number) => {
    try {
      const res = await cuppingApi.list({ product_id: String(productId), page_size: 10 })
      setScores(res.data?.items || [])
    } catch {
      setScores([])
    }
  }

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      navigate('/login')
      return
    }
    try {
      await cartApi.add(product!.id, quantity)
      alert('已添加到购物车')
    } catch (err: any) {
      alert(err.message || '添加失败')
    }
  }

  if (loading) {
    return <div className="text-center py-12">加载中...</div>
  }

  if (!product) {
    return <div className="text-center py-12">商品不存在</div>
  }

  return (
    <div className="space-y-8">
      <button onClick={() => navigate(-1)} className="text-coffee-600 hover:underline">
        ← 返回
      </button>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="card overflow-hidden">
          <div className="aspect-square bg-gray-100 relative">
            {product.images?.[0] ? (
              <img src={product.images[0].url} alt={product.name} className="w-full h-full object-cover" />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-coffee-300 text-9xl">☕</div>
            )}
          </div>
          {product.images && product.images.length > 1 && (
            <div className="flex gap-2 p-4 overflow-x-auto">
              {product.images.map((img) => (
                <img
                  key={img.id}
                  src={img.url}
                  alt=""
                  className="w-20 h-20 object-cover rounded-lg cursor-pointer"
                />
              ))}
            </div>
          )}
        </div>

        <div className="space-y-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-800 mb-2">{product.name}</h1>
            <p className="text-gray-500">{product.origin} {product.farm && `· ${product.farm}`}</p>
          </div>

          <div className="text-4xl font-bold text-coffee-600">
            ¥{product.price}
            <span className="text-lg text-gray-400 ml-2">/ {product.weight}g</span>
          </div>

          {product.cupping_score > 0 && (
            <div className="flex items-center gap-2">
              <span className="text-2xl">⭐</span>
              <span className="text-2xl font-bold">{product.cupping_score.toFixed(1)}</span>
              <span className="text-gray-500">杯测评分</span>
            </div>
          )}

          <div className="grid grid-cols-2 gap-4">
            <div className="card p-4">
              <p className="text-sm text-gray-500">处理法</p>
              <p className="font-semibold">
                {product.process_method === 'washed' ? '水洗' :
                 product.process_method === 'natural' ? '日晒' :
                 product.process_method === 'honey' ? '蜜处理' :
                 product.process_method === 'anaerobic' ? '厌氧' : '湿刨'}
              </p>
            </div>
            <div className="card p-4">
              <p className="text-sm text-gray-500">烘焙度</p>
              <p className="font-semibold">
                {product.roast_level === 'light' ? '浅烘' :
                 product.roast_level === 'medium' ? '中烘' :
                 product.roast_level === 'medium_dark' ? '中深烘' : '深烘'}
              </p>
            </div>
            <div className="card p-4">
              <p className="text-sm text-gray-500">海拔</p>
              <p className="font-semibold">{product.altitude || '-'}</p>
            </div>
            <div className="card p-4">
              <p className="text-sm text-gray-500">品种</p>
              <p className="font-semibold">{product.variety || '-'}</p>
            </div>
          </div>

          <div className="card p-4">
            <p className="text-sm text-gray-500 mb-2">风味描述</p>
            <p className="text-gray-700">{product.flavor_notes || '暂无描述'}</p>
          </div>

          {product.description && (
            <div className="card p-4">
              <p className="text-sm text-gray-500 mb-2">商品描述</p>
              <p className="text-gray-700 whitespace-pre-wrap">{product.description}</p>
            </div>
          )}

          <div className="flex items-center gap-4">
            <div className="flex items-center border rounded-lg">
              <button
                onClick={() => setQuantity(Math.max(1, quantity - 1))}
                className="px-4 py-2 text-gray-600 hover:bg-gray-100"
              >
                -
              </button>
              <span className="px-4 py-2 border-x min-w-16 text-center">{quantity}</span>
              <button
                onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                className="px-4 py-2 text-gray-600 hover:bg-gray-100"
              >
                +
              </button>
            </div>
            <button
              onClick={handleAddToCart}
              disabled={product.stock === 0}
              className="btn btn-primary flex-1 py-3 disabled:opacity-50"
            >
              {product.stock === 0 ? '缺货' : '加入购物车'}
            </button>
          </div>

          <p className="text-sm text-gray-500">库存: {product.stock}</p>

          {product.roaster && (
            <div className="card p-4 flex items-center gap-4">
              <div className="w-12 h-12 bg-coffee-200 rounded-full flex items-center justify-center text-xl">
                {product.roaster.nickname?.[0] || product.roaster.username[0]}
              </div>
              <div>
                <p className="text-sm text-gray-500">烘焙师</p>
                <p className="font-semibold">{product.roaster.nickname || product.roaster.username}</p>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="space-y-4">
        <h2 className="text-xl font-bold text-gray-800">杯测评分</h2>
        {scores.length === 0 ? (
          <p className="text-gray-500">暂无评分</p>
        ) : (
          <div className="space-y-4">
            {scores.map((score) => (
              <div key={score.id} className="card p-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="font-semibold">
                    {score.user?.nickname || score.user?.username || '匿名用户'}
                  </span>
                  <span className="text-coffee-600 font-bold">{score.overall_score.toFixed(2)}分</span>
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
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
