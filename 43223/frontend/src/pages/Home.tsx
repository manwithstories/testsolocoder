import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { productApi } from '@/api/product'
import { Product } from '@/types'

export default function Home() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadProducts()
  }, [])

  const loadProducts = async () => {
    try {
      const res = await productApi.list({ page: 1, page_size: 8, sort_by: 'cupping_score', sort_order: 'desc' })
      setProducts(res.data?.items || [])
    } catch {
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-12">
      <section className="bg-gradient-to-r from-coffee-600 to-coffee-800 text-white rounded-2xl p-12 text-center">
        <h1 className="text-4xl font-bold mb-4">探索精品咖啡豆的世界</h1>
        <p className="text-coffee-200 text-lg mb-8">从产地到杯中，每一颗咖啡豆都值得被认真对待</p>
        <Link to="/products" className="btn bg-white text-coffee-700 hover:bg-coffee-50 px-8 py-3 text-lg">
          浏览咖啡豆 →
        </Link>
      </section>

      <section>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-gray-800">热门推荐</h2>
          <Link to="/products" className="text-coffee-600 hover:underline">查看更多 →</Link>
        </div>

        {loading ? (
          <div className="text-center text-gray-500 py-12">加载中...</div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {products.map((product) => (
              <Link
                key={product.id}
                to={`/products/${product.id}`}
                className="card hover:shadow-md transition-shadow overflow-hidden"
              >
                <div className="aspect-square bg-gray-100 relative">
                  {product.images?.[0] ? (
                    <img src={product.images[0].url} alt={product.name} className="w-full h-full object-cover" />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-coffee-300 text-6xl">☕</div>
                  )}
                  {product.cupping_score > 0 && (
                    <div className="absolute top-2 right-2 bg-coffee-600 text-white text-xs px-2 py-1 rounded-full">
                      {product.cupping_score.toFixed(1)}分
                    </div>
                  )}
                </div>
                <div className="p-4">
                  <h3 className="font-semibold text-gray-800 truncate">{product.name}</h3>
                  <p className="text-sm text-gray-500 mb-2">{product.origin}</p>
                  <div className="flex items-center justify-between">
                    <span className="text-coffee-600 font-bold">¥{product.price}</span>
                    <span className="text-xs text-gray-400">{product.weight}g</span>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </section>

      <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card p-6 text-center">
          <div className="text-4xl mb-3">🌍</div>
          <h3 className="font-bold text-gray-800 mb-2">全球产地</h3>
          <p className="text-gray-500 text-sm">精选来自世界各地的优质咖啡豆</p>
        </div>
        <div className="card p-6 text-center">
          <div className="text-4xl mb-3">🔥</div>
          <h3 className="font-bold text-gray-800 mb-2">专业烘焙</h3>
          <p className="text-gray-500 text-sm">认证烘焙师精心烘焙每一批豆子</p>
        </div>
        <div className="card p-6 text-center">
          <div className="text-4xl mb-3">📊</div>
          <h3 className="font-bold text-gray-800 mb-2">杯测评分</h3>
          <p className="text-gray-500 text-sm">专业杯测师为每款咖啡豆打分</p>
        </div>
      </section>
    </div>
  )
}
