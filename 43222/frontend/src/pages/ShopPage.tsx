import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { ShoppingCart, Plus, Filter, Search } from 'lucide-react'
import { productApi, cartApi } from '@/api'
import { useCartStore } from '@/store/cart'

export default function ShopPage() {
  const queryClient = useQueryClient()
  const addItem = useCartStore((state) => state.addItem)
  const [search, setSearch] = useState('')
  const [category, setCategory] = useState('')
  const [sortBy, setSortBy] = useState('created_at')

  const { data, isLoading } = useQuery({
    queryKey: ['products', search, category, sortBy],
    queryFn: () =>
      productApi.getAll({
        search: search || undefined,
        category: category || undefined,
        sort_by: sortBy,
        page_size: 50,
      }),
  })

  const addToCartMutation = useMutation({
    mutationFn: (productId: string) => cartApi.add({ product_id: productId, quantity: 1 }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] })
    },
  })

  const handleAddToCart = (product: any) => {
    addToCartMutation.mutate(product.id)
    addItem({
      id: '',
      user_id: '',
      product_id: product.id,
      quantity: 1,
      created_at: '',
      updated_at: '',
      product: product,
    })
  }

  const products = data?.data?.products || []

  const categories = ['种子', '肥料', '工具', '容器', '其他']

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">园艺商城</h1>
        <p className="text-gray-500">购买种子、肥料和园艺工具</p>
      </div>

      {/* Search and Filter */}
      <div className="card">
        <div className="card-body space-y-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="搜索商品..."
              className="input pl-10"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
          <div className="flex flex-wrap gap-4">
            <div className="flex items-center gap-2">
              <Filter className="w-4 h-4 text-gray-400" />
              <select
                className="input w-auto py-1 text-sm"
                value={category}
                onChange={(e) => setCategory(e.target.value)}
              >
                <option value="">全部分类</option>
                {categories.map((cat) => (
                  <option key={cat} value={cat}>
                    {cat}
                  </option>
                ))}
              </select>
            </div>
            <select
              className="input w-auto py-1 text-sm"
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
            >
              <option value="created_at">最新上架</option>
              <option value="price">价格：低到高</option>
              <option value="-price">价格：高到低</option>
            </select>
          </div>
        </div>
      </div>

      {/* Product List */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">加载中...</div>
      ) : products.length === 0 ? (
        <div className="card text-center py-12">
          <ShoppingCart className="w-16 h-16 text-gray-300 mx-auto mb-4" />
          <p className="text-gray-500">没有找到商品</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {products.map((product: any) => (
            <div key={product.id} className="card hover:shadow-md transition-shadow group">
              <div className="aspect-square bg-gray-50 flex items-center justify-center overflow-hidden">
                {product.image_urls ? (
                  <img
                    src={product.image_urls.split(',')[0]}
                    alt={product.name}
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform"
                  />
                ) : (
                  <ShoppingCart className="w-12 h-12 text-gray-300" />
                )}
              </div>
              <div className="card-body">
                <h3 className="font-semibold text-gray-900 mb-1 truncate">{product.name}</h3>
                {product.category && (
                  <span className="badge bg-garden-100 text-garden-700 mb-2">
                    {product.category}
                  </span>
                )}
                <p className="text-sm text-gray-500 line-clamp-2 h-10">
                  {product.description}
                </p>
                <div className="flex items-center justify-between mt-4">
                  <span className="text-xl font-bold text-garden-600">
                    ¥{product.price?.toFixed(2)}
                  </span>
                  <button
                    onClick={() => handleAddToCart(product)}
                    disabled={product.stock <= 0}
                    className="btn-primary text-sm"
                  >
                    <Plus className="w-4 h-4 mr-1" />
                    {product.stock <= 0 ? '缺货' : '加入购物车'}
                  </button>
                </div>
                {product.stock <= 10 && product.stock > 0 && (
                  <p className="text-xs text-amber-600 mt-2">仅剩 {product.stock} 件</p>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
