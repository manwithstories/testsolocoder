import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { searchApi } from '@/api/search'
import { Product } from '@/types'

export default function Search() {
  const [keyword, setKeyword] = useState('')
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [filters, setFilters] = useState({
    origin: '',
    roast_level: '',
    process_method: '',
    min_price: '',
    max_price: '',
    min_score: '',
    max_score: '',
  })
  const [availableFilters, setAvailableFilters] = useState<any>(null)
  const [suggestions, setSuggestions] = useState<any[]>([])

  useEffect(() => {
    if (keyword.trim()) {
      const timer = setTimeout(() => {
        loadSuggestions(keyword)
      }, 300)
      return () => clearTimeout(timer)
    } else {
      setSuggestions([])
    }
  }, [keyword])

  useEffect(() => {
    loadProducts()
  }, [page, filters, keyword])

  const loadSuggestions = async (q: string) => {
    try {
      const res = await searchApi.suggest(q)
      setSuggestions(res.data || [])
    } catch {
      setSuggestions([])
    }
  }

  const loadProducts = async () => {
    setLoading(true)
    try {
      const res = await searchApi.searchProducts({
        q: keyword,
        page,
        page_size: 12,
        ...filters,
      })
      setProducts(res.data?.items || [])
      setTotal(res.data?.total || 0)
      setAvailableFilters(res.data?.filters || null)
    } catch {
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  const totalPages = Math.ceil(total / 12)

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">搜索咖啡豆</h1>

      <div className="relative">
        <input
          type="text"
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          className="input text-lg py-3"
          placeholder="搜索咖啡豆名称、产地、风味..."
          onKeyDown={(e) => e.key === 'Enter' && loadProducts()}
        />
        {suggestions.length > 0 && (
          <div className="absolute top-full left-0 right-0 bg-white rounded-lg shadow-lg mt-1 z-10">
            {suggestions.map((s, i) => (
              <div
                key={i}
                className="px-4 py-2 hover:bg-gray-100 cursor-pointer"
                onClick={() => {
                  if (s.type === 'product' && s.id) {
                    // navigate to product
                  } else {
                    setKeyword(s.name)
                  }
                  setSuggestions([])
                }}
              >
                {s.name}
                <span className="text-xs text-gray-400 ml-2">{s.type}</span>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="card p-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <select
            value={filters.origin}
            onChange={(e) => setFilters({ ...filters, origin: e.target.value })}
            className="input"
          >
            <option value="">全部产地</option>
            {availableFilters?.origins?.map((o: string) => (
              <option key={o} value={o}>{o}</option>
            ))}
          </select>
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
          <div className="flex items-center gap-2">
            <input
              type="number"
              placeholder="最低价"
              value={filters.min_price}
              onChange={(e) => setFilters({ ...filters, min_price: e.target.value })}
              className="input"
            />
            <input
              type="number"
              placeholder="最高价"
              value={filters.max_price}
              onChange={(e) => setFilters({ ...filters, max_price: e.target.value })}
              className="input"
            />
          </div>
        </div>
        <div className="flex items-center gap-4 mt-4">
          <input
            type="number"
            placeholder="最低评分"
            step="0.1"
            value={filters.min_score}
            onChange={(e) => setFilters({ ...filters, min_score: e.target.value })}
            className="input w-32"
          />
          <input
            type="number"
            placeholder="最高评分"
            step="0.1"
            value={filters.max_score}
            onChange={(e) => setFilters({ ...filters, max_score: e.target.value })}
            className="input w-32"
          />
          <button onClick={() => setPage(1)} className="btn btn-primary">
            搜索
          </button>
          <button
            onClick={() => {
              setFilters({ origin: '', roast_level: '', process_method: '', min_price: '', max_price: '', min_score: '', max_score: '' })
              setKeyword('')
              setPage(1)
            }}
            className="btn btn-secondary"
          >
            重置
          </button>
        </div>
      </div>

      {loading ? (
        <div className="text-center py-12">搜索中...</div>
      ) : products.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-6xl mb-4">🔍</div>
          <p className="text-gray-500">没有找到匹配的商品</p>
        </div>
      ) : (
        <>
          <p className="text-gray-600">共找到 {total} 个结果</p>
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
        </>
      )}
    </div>
  )
}
