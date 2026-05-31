import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { productApi } from '@/api/product'
import { Product, RoastLevel, ProcessMethod } from '@/types'

export default function Products() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [filters, setFilters] = useState({
    keyword: '',
    origin: '',
    roast_level: '',
    process_method: '',
    min_price: '',
    max_price: '',
    sort_by: 'created_at',
    sort_order: 'desc',
  })

  useEffect(() => {
    loadProducts()
  }, [page, filters])

  const loadProducts = async () => {
    setLoading(true)
    try {
      const res = await productApi.list({
        page,
        page_size: 12,
        ...filters,
      })
      setProducts(res.data?.items || [])
      setTotal(res.data?.total || 0)
    } catch {
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  const totalPages = Math.ceil(total / 12)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-gray-800">咖啡豆</h1>
      </div>

      <div className="card p-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="搜索名称/产地..."
            value={filters.keyword}
            onChange={(e) => setFilters({ ...filters, keyword: e.target.value })}
            className="input"
          />
          <select
            value={filters.roast_level}
            onChange={(e) => setFilters({ ...filters, roast_level: e.target.value })}
            className="input"
          >
            <option value="">全部烘焙度</option>
            <option value="light">浅烘</option>
            <option value="medium">中烘</option>
            <option value="medium_dark">中深烘</option>
            <option value="dark">深烘</option>
          </select>
          <select
            value={filters.process_method}
            onChange={(e) => setFilters({ ...filters, process_method: e.target.value })}
            className="input"
          >
            <option value="">全部处理法</option>
            <option value="washed">水洗</option>
            <option value="natural">日晒</option>
            <option value="honey">蜜处理</option>
            <option value="anaerobic">厌氧</option>
            <option value="wet_hulled">湿刨</option>
          </select>
          <select
            value={`${filters.sort_by}_${filters.sort_order}`}
            onChange={(e) => {
              const [sortBy, sortOrder] = e.target.value.split('_')
              setFilters({ ...filters, sort_by: sortBy, sort_order: sortOrder })
            }}
            className="input"
          >
            <option value="created_at_desc">最新上架</option>
            <option value="price_asc">价格从低到高</option>
            <option value="price_desc">价格从高到低</option>
            <option value="cupping_score_desc">评分从高到低</option>
          </select>
        </div>
        <div className="flex items-center gap-4 mt-4">
          <input
            type="number"
            placeholder="最低价"
            value={filters.min_price}
            onChange={(e) => setFilters({ ...filters, min_price: e.target.value })}
            className="input w-32"
          />
          <span className="text-gray-500">-</span>
          <input
            type="number"
            placeholder="最高价"
            value={filters.max_price}
            onChange={(e) => setFilters({ ...filters, max_price: e.target.value })}
            className="input w-32"
          />
          <button
            onClick={() => setPage(1)}
            className="btn btn-outline"
          >
            筛选
          </button>
        </div>
      </div>

      {loading ? (
        <div className="text-center text-gray-500 py-12">加载中...</div>
      ) : products.length === 0 ? (
        <div className="text-center text-gray-500 py-12">暂无商品</div>
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
                <div className="flex items-center gap-2 mb-2">
                  <span className="badge bg-coffee-100 text-coffee-700">
                    {product.roast_level === 'light' ? '浅烘' :
                     product.roast_level === 'medium' ? '中烘' :
                     product.roast_level === 'medium_dark' ? '中深烘' : '深烘'}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-coffee-600 font-bold">¥{product.price}</span>
                  <span className="text-xs text-gray-400">{product.weight}g</span>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex justify-center items-center gap-2">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="btn btn-secondary disabled:opacity-50"
          >
            上一页
          </button>
          <span className="text-gray-600">{page} / {totalPages}</span>
          <button
            onClick={() => setPage(Math.min(totalPages, page + 1))}
            disabled={page >= totalPages}
            className="btn btn-secondary disabled:opacity-50"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  )
}
